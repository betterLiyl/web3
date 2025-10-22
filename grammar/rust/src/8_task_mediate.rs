// ### 任务4：HTTP客户端库 - 已完成 ✅
// **目标**：掌握异步编程、Future、错误处理
// **描述**：实现一个功能完整的HTTP客户端库，支持连接池和中间件

// **完成状态**：
// ✅ 1. 设计HTTP请求和响应结构
// ✅ 2. 实现基于tokio的异步HTTP客户端
// ✅ 3. 添加连接池和连接复用
// ✅ 4. 实现中间件系统（重试、日志、认证）
// ✅ 5. 支持不同的序列化格式
// ✅ 6. 添加超时和取消机制

// **注意**：由于依赖版本兼容性问题，完整实现请参考 8_task_mediate_simple.rs
// 该文件包含了一个功能完整的纯Rust实现版本，避免了外部依赖的版本冲突问题。

// **简化版实现特性**：
// - 完整的HTTP客户端结构设计
// - 异步请求处理和连接池管理
// - 中间件系统（日志、重试、认证、超时）
// - 错误处理和重试机制
// - JSON序列化支持
// - 构建器模式配置
// - 完整的单元测试覆盖

// **运行方式**：
// cargo run --bin 8_task_mediate_simple
// cargo test --bin 8_task_mediate_simple

use std::time::Duration;
use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::{Mutex, Semaphore};
use reqwest::{Client as ReqwestClient, Method, Response, Url};
use serde::{Deserialize, Serialize};
use thiserror::Error;
use async_trait::async_trait;
use futures::future::BoxFuture;
use log::{info, warn, error};

// 错误处理
#[derive(Error, Debug)]
pub enum HttpClientError {
    #[error("Request failed: {0}")]
    RequestFailed(#[from] reqwest::Error),
    #[error("Timeout error")]
    Timeout,
    #[error("Connection pool exhausted")]
    PoolExhausted,
    #[error("Serialization error: {0}")]
    SerializationError(#[from] serde_json::Error),
    #[error("URL parse error: {0}")]
    UrlParseError(String),
    #[error("Middleware error: {0}")]
    MiddlewareError(String),
}

pub type Result<T> = std::result::Result<T, HttpClientError>;

// HTTP请求结构
#[derive(Debug, Clone)]
pub struct HttpRequest {
    pub method: Method,
    pub url: Url,
    pub headers: HashMap<String, String>,
    pub body: Option<Vec<u8>>,
    pub timeout: Option<Duration>,
}

impl HttpRequest {
    pub fn new(method: Method, url: &str) -> Result<Self> {
        let parsed_url = Url::parse(url)
            .map_err(|e| HttpClientError::UrlParseError(e.to_string()))?;
        
        Ok(Self {
            method,
            url: parsed_url,
            headers: HashMap::new(),
            body: None,
            timeout: None,
        })
    }

    pub fn header(mut self, key: &str, value: &str) -> Self {
        self.headers.insert(key.to_string(), value.to_string());
        self
    }

    pub fn json<T: Serialize>(mut self, data: &T) -> Result<Self> {
        self.body = Some(serde_json::to_vec(data)?);
        self.headers.insert("Content-Type".to_string(), "application/json".to_string());
        Ok(self)
    }

    pub fn timeout(mut self, timeout: Duration) -> Self {
        self.timeout = Some(timeout);
        self
    }
}

// HTTP响应结构
#[derive(Debug, Clone)]
pub struct HttpResponse {
    pub status: u16,
    pub headers: HashMap<String, String>,
    pub body: Vec<u8>,
}

impl HttpResponse {
    pub fn json<T: for<'de> Deserialize<'de>>(&self) -> Result<T> {
        Ok(serde_json::from_slice(&self.body)?)
    }

    pub fn text(&self) -> String {
        String::from_utf8_lossy(&self.body).to_string()
    }

    pub fn is_success(&self) -> bool {
        self.status >= 200 && self.status < 300
    }
}

// 中间件trait
#[async_trait]
pub trait Middleware: Send + Sync {
    async fn handle(&self, request: &mut HttpRequest, next: Next<'_>) -> Result<HttpResponse>;
}
#[derive(Clone)]
pub struct Next<'a> {
    middlewares: &'a [Arc<dyn Middleware>],
    index: usize,
    client: &'a HttpClient,
}

impl<'a> Next<'a> {
    pub async fn run(mut self, request: &mut HttpRequest) -> Result<HttpResponse> {
        if self.index < self.middlewares.len() {
            let middleware = self.middlewares[self.index].clone();
            self.index += 1;
            middleware.handle(request, self).await
        } else {
            self.client.execute_request(request).await
        }
    }
}

// 日志中间件
pub struct LoggingMiddleware;

#[async_trait]
impl Middleware for LoggingMiddleware {
    async fn handle(&self, request: &mut HttpRequest, next: Next<'_>) -> Result<HttpResponse> {
        info!("Sending {} request to {}", request.method, request.url);
        let start = std::time::Instant::now();
        
        let response = next.run(request).await;
        
        let duration = start.elapsed();
        match &response {
            Ok(resp) => info!("Request completed in {:?} with status {}", duration, resp.status),
            Err(e) => error!("Request failed in {:?}: {}", duration, e),
        }
        
        response
    }
}

// 重试中间件
pub struct RetryMiddleware {
    max_retries: usize,
    retry_delay: Duration,
}

impl RetryMiddleware {
    pub fn new(max_retries: usize, retry_delay: Duration) -> Self {
        Self { max_retries, retry_delay }
    }
}

#[async_trait]
impl Middleware for RetryMiddleware {
    async fn handle(&self, request: &mut HttpRequest, next: Next<'_>) -> Result<HttpResponse> {
        let mut attempts = 0;
        
        loop {
            let response = next.clone().run(request).await;
            
            match response {
                Ok(resp) if resp.is_success() => return Ok(resp),
                Ok(resp) if attempts < self.max_retries => {
                    warn!("Request failed with status {}, retrying... (attempt {}/{})", 
                          resp.status, attempts + 1, self.max_retries);
                    attempts += 1;
                    tokio::time::sleep(self.retry_delay).await;
                    continue;
                }
                Ok(resp) => return Ok(resp),
                Err(e) if attempts < self.max_retries => {
                    warn!("Request failed with error {}, retrying... (attempt {}/{})", 
                          e, attempts + 1, self.max_retries);
                    attempts += 1;
                    tokio::time::sleep(self.retry_delay).await;
                    continue;
                }
                Err(e) => return Err(e),
            }
        }
    }
}

// 超时中间件
pub struct TimeoutMiddleware {
    timeout: Duration,
}

impl TimeoutMiddleware {
    pub fn new(timeout: Duration) -> Self {
        Self { timeout }
    }
}

#[async_trait]
impl Middleware for TimeoutMiddleware {
    async fn handle(&self, request: &mut HttpRequest, next: Next<'_>) -> Result<HttpResponse> {
        let timeout = request.timeout.unwrap_or(self.timeout);
        
        match tokio::time::timeout(timeout, next.run(request)).await {
            Ok(result) => result,
            Err(_) => Err(HttpClientError::Timeout),
        }
    }
}

// 认证中间件
pub struct AuthMiddleware {
    token: String,
}

impl AuthMiddleware {
    pub fn bearer(token: &str) -> Self {
        Self { token: token.to_string() }
    }
}

#[async_trait]
impl Middleware for AuthMiddleware {
    async fn handle(&self, request: &mut HttpRequest, next: Next<'_>) -> Result<HttpResponse> {
        request.headers.insert(
            "Authorization".to_string(),
            format!("Bearer {}", self.token)
        );
        next.run(request).await
    }
}

// 连接池
pub struct ConnectionPool {
    semaphore: Arc<Semaphore>,
    client: ReqwestClient,
}

impl ConnectionPool {
    pub fn new(max_connections: usize) -> Self {
        let client = ReqwestClient::builder()
            .pool_max_idle_per_host(max_connections)
            .pool_idle_timeout(Duration::from_secs(30))
            .build()
            .expect("Failed to create HTTP client");

        Self {
            semaphore: Arc::new(Semaphore::new(max_connections)),
            client,
        }
    }

    pub async fn execute(&self, request: &HttpRequest) -> Result<HttpResponse> {
        let _permit = self.semaphore.acquire().await
            .map_err(|_| HttpClientError::PoolExhausted)?;

        let mut req_builder = self.client.request(request.method.clone(), request.url.clone());

        // 添加请求头
        for (key, value) in &request.headers {
            req_builder = req_builder.header(key, value);
        }

        // 添加请求体
        if let Some(body) = &request.body {
            req_builder = req_builder.body(body.clone());
        }

        // 设置超时
        if let Some(timeout) = request.timeout {
            req_builder = req_builder.timeout(timeout);
        }

        let response = req_builder.send().await?;
        
        let status = response.status().as_u16();
        let headers = response.headers()
            .iter()
            .map(|(k, v)| (k.to_string(), v.to_str().unwrap_or("").to_string()))
            .collect();
        let body = response.bytes().await?.to_vec();

        Ok(HttpResponse {
            status,
            headers,
            body,
        })
    }
}

// HTTP客户端
pub struct HttpClient {
    pool: Arc<ConnectionPool>,
    middlewares: Vec<Arc<dyn Middleware>>,
}

impl HttpClient {
    pub fn new() -> Self {
        Self {
            pool: Arc::new(ConnectionPool::new(10)),
            middlewares: Vec::new(),
        }
    }

    pub fn with_pool(max_connections: usize) -> Self {
        Self {
            pool: Arc::new(ConnectionPool::new(max_connections)),
            middlewares: Vec::new(),
        }
    }

    pub fn with_middleware(mut self, middleware: Arc<dyn Middleware>) -> Self {
        self.middlewares.push(middleware);
        self
    }

    pub async fn request(&self, mut request: HttpRequest) -> Result<HttpResponse> {
        let next = Next {
            middlewares: &self.middlewares,
            index: 0,
            client: self,
        };
        
        next.run(&mut request).await
    }

    async fn execute_request(&self, request: &HttpRequest) -> Result<HttpResponse> {
        self.pool.execute(request).await
    }

    // 便捷方法
    pub async fn get(&self, url: &str) -> Result<HttpResponse> {
        let request = HttpRequest::new(Method::GET, url)?;
        self.request(request).await
    }

    pub async fn post<T: Serialize>(&self, url: &str, data: &T) -> Result<HttpResponse> {
        let request = HttpRequest::new(Method::POST, url)?.json(data)?;
        self.request(request).await
    }

    pub async fn put<T: Serialize>(&self, url: &str, data: &T) -> Result<HttpResponse> {
        let request = HttpRequest::new(Method::PUT, url)?.json(data)?;
        self.request(request).await
    }

    pub async fn delete(&self, url: &str) -> Result<HttpResponse> {
        let request = HttpRequest::new(Method::DELETE, url)?;
        self.request(request).await
    }
}

impl Default for HttpClient {
    fn default() -> Self {
        Self::new()
    }
}

// 客户端构建器
pub struct HttpClientBuilder {
    max_connections: usize,
    middlewares: Vec<Arc<dyn Middleware>>,
}

impl HttpClientBuilder {
    pub fn new() -> Self {
        Self {
            max_connections: 10,
            middlewares: Vec::new(),
        }
    }

    pub fn max_connections(mut self, max: usize) -> Self {
        self.max_connections = max;
        self
    }

    pub fn with_logging(mut self) -> Self {
        self.middlewares.push(Arc::new(LoggingMiddleware));
        self
    }

    pub fn with_retry(mut self, max_retries: usize, delay: Duration) -> Self {
        self.middlewares.push(Arc::new(RetryMiddleware::new(max_retries, delay)));
        self
    }

    pub fn with_timeout(mut self, timeout: Duration) -> Self {
        self.middlewares.push(Arc::new(TimeoutMiddleware::new(timeout)));
        self
    }

    pub fn with_auth(mut self, token: &str) -> Self {
        self.middlewares.push(Arc::new(AuthMiddleware::bearer(token)));
        self
    }

    pub fn build(self) -> HttpClient {
        let mut client = HttpClient::with_pool(self.max_connections);
        for middleware in self.middlewares {
            client = client.with_middleware(middleware);
        }
        client
    }
}

impl Default for HttpClientBuilder {
    fn default() -> Self {
        Self::new()
    }
}

#[tokio::main]
async fn main() -> Result<()> {
    env_logger::init();

    // 示例1: 基本用法
    println!("=== 基本HTTP客户端示例 ===");
    let client = HttpClient::new();
    
    match client.get("https://httpbin.org/get").await {
        Ok(response) => {
            println!("Status: {}", response.status);
            println!("Response: {}", response.text());
        }
        Err(e) => println!("Error: {}", e),
    }

    // 示例2: 带中间件的客户端
    println!("\n=== 带中间件的HTTP客户端示例 ===");
    let client_with_middleware = HttpClientBuilder::new()
        .max_connections(5)
        .with_logging()
        .with_retry(3, Duration::from_millis(500))
        .with_timeout(Duration::from_secs(10))
        .build();

    // POST请求示例
    #[derive(Serialize)]
    struct PostData {
        name: String,
        age: u32,
    }

    let data = PostData {
        name: "Alice".to_string(),
        age: 30,
    };

    match client_with_middleware.post("https://httpbin.org/post", &data).await {
        Ok(response) => {
            println!("POST Status: {}", response.status);
            if response.is_success() {
                println!("POST successful!");
            }
        }
        Err(e) => println!("POST Error: {}", e),
    }

    // 示例3: 自定义请求
    println!("\n=== 自定义请求示例 ===");
    let custom_request = HttpRequest::new(Method::GET, "https://httpbin.org/headers")?
        .header("User-Agent", "Custom-HTTP-Client/1.0")
        .header("X-Custom-Header", "test-value")
        .timeout(Duration::from_secs(5));

    match client_with_middleware.request(custom_request).await {
        Ok(response) => {
            println!("Custom request status: {}", response.status);
            println!("Response body: {}", response.text());
        }
        Err(e) => println!("Custom request error: {}", e),
    }

    // 示例4: JSON响应解析
    println!("\n=== JSON响应解析示例 ===");
    #[derive(Deserialize, Debug)]
    struct IpInfo {
        origin: String,
    }

    match client.get("https://httpbin.org/ip").await {
        Ok(response) => {
            if let Ok(ip_info) = response.json::<IpInfo>() {
                println!("Your IP: {}", ip_info.origin);
            }
        }
        Err(e) => println!("IP request error: {}", e),
    }

    Ok(())
}
