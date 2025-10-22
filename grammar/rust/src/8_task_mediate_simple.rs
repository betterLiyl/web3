// ============================================================================
// 依赖包导入说明
// ============================================================================

// 标准库导入
use std::collections::HashMap;    // 标准库的哈希映射，用于存储键值对（如HTTP头部）
use std::time::Duration;          // 标准库的时间间隔类型，用于设置超时时间
use std::sync::Arc;               // 标准库的原子引用计数智能指针，用于多线程间安全共享数据

// Tokio异步运行时相关导入
use tokio::sync::Semaphore;       // Tokio提供的信号量，用于控制并发连接数（连接池实现）

// Serde序列化框架导入
use serde::{Deserialize, Serialize}; // Serde库的序列化和反序列化trait，用于JSON数据处理

// 错误处理库导入
use thiserror::Error;             // thiserror库提供的Error derive宏，简化错误类型定义

// 日志库导入
use log::{info, warn, error};     // log库提供的日志宏，用于记录不同级别的日志信息

// ============================================================================
// 错误类型定义
// ============================================================================

// #[derive(Error, Debug)] 是属性宏的组合：
// - Error: 来自thiserror库，自动为枚举实现std::error::Error trait
// - Debug: 标准库trait，允许使用{:?}格式化输出，用于调试
#[derive(Error, Debug)]
pub enum HttpClientError {
    // #[error("...")] 是thiserror提供的属性宏，定义错误的显示信息
    // {0} 表示元组结构体的第一个字段，用于格式化字符串
    #[error("Request failed: {0}")]
    RequestFailed(String),                    // 请求失败错误，包含失败原因
    
    #[error("Connection pool exhausted")]
    PoolExhausted,                           // 连接池耗尽错误
    
    #[error("Timeout error: {0}")]
    Timeout(String),                         // 超时错误，包含超时详情
    
    #[error("Authentication failed")]
    AuthenticationFailed,                    // 认证失败错误
    
    #[error("URL parse error: {0}")]
    UrlParseError(String),                   // URL解析错误，包含解析失败的原因
    
    #[error("Serialization error: {0}")]
    SerializationError(String),              // 序列化/反序列化错误，包含错误详情
}

// 类型别名：简化Result类型的使用，T是成功时的类型，错误类型固定为HttpClientError
pub type Result<T> = std::result::Result<T, HttpClientError>;

// ============================================================================
// HTTP方法枚举定义
// ============================================================================

// #[derive(Debug, Clone)] 属性宏说明：
// - Debug: 允许使用{:?}格式化输出，便于调试
// - Clone: 允许克隆枚举值，因为HTTP方法需要在多处使用
#[derive(Debug, Clone)]
pub enum HttpMethod {
    GET,        // HTTP GET方法，用于获取资源
    POST,       // HTTP POST方法，用于创建资源
    PUT,        // HTTP PUT方法，用于更新资源
    DELETE,     // HTTP DELETE方法，用于删除资源
    PATCH,      // HTTP PATCH方法，用于部分更新资源
    HEAD,       // HTTP HEAD方法，只获取响应头
    OPTIONS,    // HTTP OPTIONS方法，用于获取服务器支持的方法
}

// 为HttpMethod实现Display trait，来自std::fmt模块
// 这允许使用{}格式化输出，将枚举转换为字符串
impl std::fmt::Display for HttpMethod {
    // fmt函数是Display trait的必需方法
    // &self: 不可变引用自身
    // f: 格式化器的可变引用，用于写入输出
    // -> std::fmt::Result: 返回格式化结果，成功或失败
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        // match表达式匹配枚举的每个变体
        match self {
            HttpMethod::GET => write!(f, "GET"),           // write!宏将字符串写入格式化器
            HttpMethod::POST => write!(f, "POST"),
            HttpMethod::PUT => write!(f, "PUT"),
            HttpMethod::DELETE => write!(f, "DELETE"),
            HttpMethod::PATCH => write!(f, "PATCH"),
            HttpMethod::HEAD => write!(f, "HEAD"),
            HttpMethod::OPTIONS => write!(f, "OPTIONS"),
        }
    }
}

// ============================================================================
// HTTP请求结构体定义
// ============================================================================

// #[derive(Debug, Clone)] 属性宏说明：
// - Debug: 允许调试输出
// - Clone: 允许克隆请求对象，便于重试等操作
#[derive(Debug, Clone)]
pub struct HttpRequest {
    pub method: HttpMethod,                    // HTTP方法（GET、POST等）
    pub url: String,                          // 请求的URL地址
    pub headers: HashMap<String, String>,     // HTTP头部，键值对形式存储
    pub body: Option<String>,                 // 请求体，Option表示可能为空
    pub timeout: Option<Duration>,            // 超时时间，Option表示可选配置
}

impl HttpRequest {
    // 构造函数：创建新的HTTP请求
    // method: HTTP方法枚举
    // url: URL字符串切片（&str是字符串的借用）
    // -> Result<Self>: 返回Result类型，成功时是Self（即HttpRequest），失败时是错误
    pub fn new(method: HttpMethod, url: &str) -> Result<Self> {
        // 简单的URL验证：检查URL是否以http://或https://开头
        if !url.starts_with("http://") && !url.starts_with("https://") {
            // 如果URL格式不正确，返回错误
            // to_string()将&str转换为String
            return Err(HttpClientError::UrlParseError(
                "URL must start with http:// or https://".to_string()
            ));
        }

        // 创建并返回HttpRequest实例
        Ok(HttpRequest {
            method,                               // 使用传入的HTTP方法
            url: url.to_string(),                // 将&str转换为String
            headers: HashMap::new(),             // 创建空的HashMap存储头部
            body: None,                          // 初始时没有请求体
            timeout: None,                       // 初始时没有设置超时
        })
    }

    // 构建器模式方法：添加HTTP头部
    // mut self: 获取self的所有权并允许修改
    // key, value: 头部的键和值，都是字符串切片
    // -> Self: 返回修改后的自身，支持链式调用
    pub fn header(mut self, key: &str, value: &str) -> Self {
        // 将键值对插入到headers HashMap中
        self.headers.insert(key.to_string(), value.to_string());
        self  // 返回修改后的自身
    }

    // 构建器模式方法：设置请求体
    pub fn body(mut self, body: String) -> Self {
        self.body = Some(body);  // 将请求体包装在Some中
        self
    }

    // 构建器模式方法：设置JSON请求体
    // <T: Serialize>: 泛型约束，T必须实现Serialize trait
    // data: 要序列化的数据的引用
    pub fn json<T: Serialize>(mut self, data: &T) -> Result<Self> {
        // 使用serde_json将数据序列化为JSON字符串
        // map_err将serde_json的错误转换为我们的错误类型
        let json_body = serde_json::to_string(data)
            .map_err(|e| HttpClientError::SerializationError(e.to_string()))?;
        
        self.body = Some(json_body);  // 设置JSON字符串为请求体
        // 自动设置Content-Type头部为application/json
        self.headers.insert("Content-Type".to_string(), "application/json".to_string());
        Ok(self)  // 返回Result包装的自身
    }

    // 构建器模式方法：设置超时时间
    pub fn timeout(mut self, timeout: Duration) -> Self {
        self.timeout = Some(timeout);
        self
    }
}

// ============================================================================
// HTTP响应结构体定义
// ============================================================================

// HTTP响应结构体，存储服务器返回的响应信息
#[derive(Debug, Clone)]
pub struct HttpResponse {
    pub status: u16,                          // HTTP状态码（如200、404、500等）
    pub headers: HashMap<String, String>,     // 响应头部信息
    pub body: String,                         // 响应体内容
}

impl HttpResponse {
    // 构造函数：创建新的HTTP响应
    pub fn new(status: u16, body: String) -> Self {
        HttpResponse {
            status,                           // HTTP状态码
            headers: HashMap::new(),          // 初始化空的头部映射
            body,                            // 响应体内容
        }
    }

    // 判断响应是否成功
    // HTTP状态码在200-299范围内表示成功
    pub fn is_success(&self) -> bool {
        self.status >= 200 && self.status < 300
    }

    // 将响应体解析为JSON
    // <T: for<'de> Deserialize<'de>>: 高阶trait约束
    // - for<'de>: 高阶生命周期，表示对任何生命周期'de都成立
    // - Deserialize<'de>: serde的反序列化trait，'de是反序列化生命周期
    pub fn json<T: for<'de> Deserialize<'de>>(&self) -> Result<T> {
        // 使用serde_json从字符串反序列化为类型T
        serde_json::from_str(&self.body)
            .map_err(|e| HttpClientError::SerializationError(e.to_string()))
    }
}

// ============================================================================
// 连接池实现
// ============================================================================

// 连接池结构体 - 展示Rust所有权和生命周期的核心概念
// 用于管理并发连接数
pub struct ConnectionPool {
    // Arc<Semaphore>: 原子引用计数智能指针，体现Rust所有权系统的核心特性
    // - Arc: Atomic Reference Counting，允许多个所有者共享同一数据
    // - 解决了Rust单一所有权规则在多线程环境下的限制
    // - 当最后一个Arc被drop时，内部数据才会被释放
    // - 线程安全的引用计数，可以在多线程间安全传递
    semaphore: Arc<Semaphore>,               // Arc包装的信号量，用于多线程共享
    max_connections: usize,                  // 最大连接数限制（栈上数据，实现Copy trait）
}

impl ConnectionPool {
    // 创建新的连接池
    // max_connections: 最大并发连接数
    pub fn new(max_connections: usize) -> Self {
        ConnectionPool {
            // Arc::new创建原子引用计数智能指针，允许多线程安全共享
            // Semaphore::new创建信号量，初始许可数为max_connections
            semaphore: Arc::new(Semaphore::new(max_connections)),
            max_connections,
        }
    }

    // 异步方法：获取连接许可 - 展示Rust借用检查器和RAII模式
    // &self: 不可变借用，允许多个线程同时调用此方法
    // async: 异步函数，返回Future<Output = Result<ConnectionGuard>>
    // -> Result<ConnectionGuard>: 返回连接守卫，体现RAII（资源获取即初始化）模式
    pub async fn acquire(&self) -> Result<ConnectionGuard> {
        // self.semaphore: 通过&self访问Arc<Semaphore>
        // .clone(): 克隆Arc智能指针（只复制引用计数，不复制数据）
        // .acquire_owned(): 获取拥有所有权的许可证
        // - 返回OwnedSemaphorePermit，它拥有许可证的所有权
        // - 当许可证被drop时，会自动释放回信号量
        // .await: 等待异步操作完成
        // ?: 错误传播操作符，如果获取失败则提前返回错误
        let permit = self.semaphore.clone().acquire_owned().await
            .map_err(|_| HttpClientError::PoolExhausted)?;  // 将错误转换为我们的错误类型
        
        // 创建ConnectionGuard，将许可证的所有权转移给它
        // 这体现了Rust的移动语义：permit的所有权被转移，不能再使用
        Ok(ConnectionGuard { _permit: permit })
    }

    // 获取当前可用连接数
    pub fn available_connections(&self) -> usize {
        self.semaphore.available_permits()   // 返回信号量的可用许可数
    }
}

// 连接守卫结构体 - 展示Rust的RAII（资源获取即初始化）模式和自动内存管理
// 这是Rust所有权系统的经典应用：通过类型系统确保资源的正确释放
pub struct ConnectionGuard {
    // _permit: 下划线前缀表示这个字段不会被直接使用，但它的存在很重要
    // OwnedSemaphorePermit: 拥有所有权的信号量许可证
    // - 当ConnectionGuard被创建时，许可证被"获取"
    // - 当ConnectionGuard被drop时，许可证自动"释放"回信号量
    // - 这确保了连接资源的自动管理，防止资源泄漏
    // - 体现了Rust的"零成本抽象"：编译时保证，运行时无额外开销
    _permit: tokio::sync::OwnedSemaphorePermit,  // 拥有所有权的信号量许可
    // 注意：这里没有显式实现Drop trait，但OwnedSemaphorePermit自己实现了Drop
    // 当ConnectionGuard超出作用域时，_permit会被自动drop，从而释放许可证
}

// ============================================================================
// 中间件系统定义
// ============================================================================

// 中间件trait定义：定义了中间件必须实现的行为
// trait类似于其他语言的接口，定义了一组方法签名
// Send + Sync: trait约束，表示实现者必须是线程安全的
// - Send: 可以在线程间转移所有权
// - Sync: 可以在多线程间安全共享引用
pub trait Middleware: Send + Sync {
    // 返回中间件名称，&str是字符串切片的引用
    fn name(&self) -> &str;
    
    // 处理请求的方法
    // &self: 不可变引用自身
    // request: 可变引用HTTP请求，允许中间件修改请求
    // -> Result<()>: 返回空的Result，()表示成功时无返回值
    fn process_request(&self, request: &mut HttpRequest) -> Result<()>;
    
    // 处理响应的方法
    // response: 可变引用HTTP响应，允许中间件修改响应
    fn process_response(&self, response: &mut HttpResponse) -> Result<()>;
}

// ============================================================================
// 日志中间件实现
// ============================================================================

// 日志中间件结构体，用于记录HTTP请求和响应
pub struct LoggingMiddleware {
    pub log_requests: bool,      // 是否记录请求日志
    pub log_responses: bool,     // 是否记录响应日志
}

impl LoggingMiddleware {
    // 构造函数：创建新的日志中间件，默认记录请求和响应
    pub fn new() -> Self {
        LoggingMiddleware {
            log_requests: true,
            log_responses: true,
        }
    }
}

// 为LoggingMiddleware实现Middleware trait
impl Middleware for LoggingMiddleware {
    // 实现trait的name方法
    fn name(&self) -> &str {
        "LoggingMiddleware"  // 返回中间件名称
    }

    // 实现请求处理方法
    fn process_request(&self, request: &mut HttpRequest) -> Result<()> {
        if self.log_requests {
            // info!是log库提供的宏，记录信息级别日志
            // {}是格式化占位符，类似于printf的%s
            info!("HTTP Request: {} {}", request.method, request.url);
            
            // 遍历请求头部，&表示借用，避免所有权转移
            for (key, value) in &request.headers {
                info!("  Header: {}: {}", key, value);
            }
            
            // if let是模式匹配，只在Some时执行
            if let Some(body) = &request.body {
                info!("  Body: {}", body);
            }
        }
        Ok(())  // 返回成功结果
    }

    // 实现响应处理方法
    fn process_response(&self, response: &mut HttpResponse) -> Result<()> {
        if self.log_responses {
            info!("HTTP Response: Status {}", response.status);
            for (key, value) in &response.headers {
                info!("  Header: {}: {}", key, value);
            }
            info!("  Body: {}", response.body);
        }
        Ok(())
    }
}

// ============================================================================
// 认证中间件实现
// ============================================================================

// 认证中间件结构体，用于自动添加认证信息到请求头
pub struct AuthMiddleware {
    pub token: String,           // 认证令牌
    pub auth_type: AuthType,     // 认证类型枚举
}

// 认证类型枚举，支持多种认证方式
#[derive(Debug, Clone)]
pub enum AuthType {
    Bearer,                      // Bearer Token认证（常用于JWT）
    Basic,                       // Basic认证（用户名密码的Base64编码）
    ApiKey(String),             // API Key认证，String是头部名称（如X-API-Key）
}

impl AuthMiddleware {
    // 创建Bearer Token认证中间件
    pub fn bearer(token: String) -> Self {
        AuthMiddleware {
            token,
            auth_type: AuthType::Bearer,
        }
    }

    // 创建Basic认证中间件
    pub fn basic(token: String) -> Self {
        AuthMiddleware {
            token,
            auth_type: AuthType::Basic,
        }
    }

    // 创建API Key认证中间件
    // header_name: 自定义头部名称（如"X-API-Key"、"Authorization"等）
    pub fn api_key(header_name: String, token: String) -> Self {
        AuthMiddleware {
            token,
            auth_type: AuthType::ApiKey(header_name),
        }
    }
}

// 为AuthMiddleware实现Middleware trait
impl Middleware for AuthMiddleware {
    fn name(&self) -> &str {
        "AuthMiddleware"
    }

    // 在请求处理中添加认证头部
    fn process_request(&self, request: &mut HttpRequest) -> Result<()> {
        // match表达式根据认证类型执行不同逻辑
        match &self.auth_type {
            AuthType::Bearer => {
                // format!宏用于字符串格式化，类似于sprintf
                request.headers.insert("Authorization".to_string(), format!("Bearer {}", self.token));
            }
            AuthType::Basic => {
                request.headers.insert("Authorization".to_string(), format!("Basic {}", self.token));
            }
            AuthType::ApiKey(header_name) => {
                // clone()创建header_name的副本，避免所有权问题
                request.headers.insert(header_name.clone(), self.token.clone());
            }
        }
        Ok(())
    }

    // 认证中间件通常不需要处理响应
    fn process_response(&self, _response: &mut HttpResponse) -> Result<()> {
        // _response前缀表示参数未使用，避免编译器警告
        Ok(())
    }
}

// ============================================================================
// HTTP客户端主体实现
// ============================================================================

// HTTP客户端结构体，整合了连接池、中间件、超时和重试功能
pub struct HttpClient {
    pool: ConnectionPool,                        // 连接池，管理并发连接数
    middlewares: Vec<Box<dyn Middleware>>,      // 中间件列表，Box<dyn Trait>是trait对象
    default_timeout: Duration,                   // 默认超时时间
    retry_config: RetryConfig,                  // 重试配置
}

// 重试配置结构体
#[derive(Debug, Clone)]
pub struct RetryConfig {
    pub max_retries: usize,                     // 最大重试次数
    pub retry_delay: Duration,                  // 重试间隔时间
    pub retry_on_status: Vec<u16>,             // 需要重试的HTTP状态码列表
}

// 为RetryConfig实现Default trait，提供默认配置
impl Default for RetryConfig {
    fn default() -> Self {
        RetryConfig {
            max_retries: 3,                                    // 默认重试3次
            retry_delay: Duration::from_millis(1000),         // 默认重试间隔1秒
            retry_on_status: vec![500, 502, 503, 504],        // 服务器错误时重试
        }
    }
}

impl HttpClient {
    // 构造函数：创建默认配置的HTTP客户端
    pub fn new() -> Self {
        HttpClient {
            pool: ConnectionPool::new(10),           // 默认10个连接
            middlewares: Vec::new(),                 // 空的中间件列表
            default_timeout: Duration::from_secs(30), // 默认30秒超时
            retry_config: RetryConfig::default(),   // 默认重试配置
        }
    }

    // 构建器模式：设置连接池大小
    pub fn with_pool_size(mut self, size: usize) -> Self {
        self.pool = ConnectionPool::new(size);
        self
    }

    // 构建器模式：设置默认超时时间
    pub fn with_timeout(mut self, timeout: Duration) -> Self {
        self.default_timeout = timeout;
        self
    }

    // 构建器模式：设置重试配置
    pub fn with_retry_config(mut self, config: RetryConfig) -> Self {
        self.retry_config = config;
        self
    }

    // 构建器模式：添加中间件 - 展示Rust泛型约束和生命周期参数
    // <M: Middleware + 'static>: 泛型约束，体现Rust类型系统的强大之处
    // - M: 泛型类型参数，可以是任何实现了Middleware trait的类型
    // - Middleware: trait约束，确保M类型具有中间件的行为
    // - 'static: 生命周期约束，要求M类型在整个程序运行期间都有效
    //   - 'static不意味着永远存在，而是意味着没有非'static的引用
    //   - 这确保中间件可以安全地存储在HttpClient中，不会出现悬垂引用
    // mut self: 获取self的可变所有权，允许修改并返回
    pub fn add_middleware<M: Middleware + 'static>(mut self, middleware: M) -> Self {
        // Box::new(middleware): 将中间件装箱到堆上
        // - 这是所有权转移：middleware的所有权被转移给Box
        // - Box<dyn Middleware>: trait对象，允许存储不同类型的中间件
        // - dyn关键字表示动态分发，运行时确定具体类型
        self.middlewares.push(Box::new(middleware));
        self // 返回self，支持链式调用（移动语义）
    }

    // 私有方法：模拟HTTP请求执行 - 展示异步编程和借用检查器
    // async fn: 异步函数，返回Future<Output = Result<HttpResponse>>
    // &self: 不可变借用，多个异步任务可以同时调用此方法
    // request: &HttpRequest: 借用HttpRequest，不获取所有权
    //   - 使用引用避免不必要的克隆，提高性能
    //   - 借用检查器确保request在函数执行期间保持有效
    async fn execute_request(&self, request: &HttpRequest) -> Result<HttpResponse> {
        // 获取连接许可，?操作符用于错误传播
        // _guard: 连接守卫，变量名前的下划线表示我们不直接使用它
        // 但它的存在确保了连接资源的RAII管理
        // 当_guard超出作用域时，连接会自动释放
        let _guard = self.pool.acquire().await?;
        
        // tokio::time::sleep: Tokio提供的异步睡眠函数
        // Duration::from_millis(100): 创建100毫秒的时间间隔
        // .await: 等待异步操作完成，让出CPU给其他任务
        // 这模拟了网络延迟，展示了异步编程的非阻塞特性
        tokio::time::sleep(Duration::from_millis(100)).await;
        
        // 根据HTTP方法模拟不同的响应
        let response = match request.method {
            HttpMethod::GET => {
                // 根据URL内容模拟不同响应
                if request.url.contains("error") {
                    HttpResponse::new(500, "Internal Server Error".to_string())
                } else if request.url.contains("notfound") {
                    HttpResponse::new(404, "Not Found".to_string())
                } else {
                    // r#"..."#是原始字符串字面量，避免转义引号
                    let mut resp = HttpResponse::new(200, r#"{"message": "GET request successful", "data": {"id": 1, "name": "test"}}"#.to_string());
                    resp.headers.insert("Content-Type".to_string(), "application/json".to_string());
                    resp
                }
            }
            HttpMethod::POST => {
                let mut resp = HttpResponse::new(201, r#"{"message": "POST request successful", "id": 123}"#.to_string());
                resp.headers.insert("Content-Type".to_string(), "application/json".to_string());
                resp
            }
            _ => {
                // 其他HTTP方法的默认响应
                HttpResponse::new(200, "Request successful".to_string())
            }
        };

        Ok(response)
    }

    // 公共方法：发送HTTP请求 - 展示可变借用和所有权转移
    // mut request: HttpRequest: 获取request的所有权，允许修改
    //   - 不是借用(&HttpRequest)，而是移动(HttpRequest)
    //   - 调用者失去对request的所有权，避免了克隆的开销
    //   - mut关键字允许我们修改request（如添加中间件处理的头部）
    pub async fn send(&self, mut request: HttpRequest) -> Result<HttpResponse> {
        // 设置默认超时时间（如果请求没有指定）
        if request.timeout.is_none() {
            request.timeout = Some(self.default_timeout);
        }

        // 处理中间件：遍历所有中间件并处理请求
        // self.middlewares.iter(): 创建中间件的不可变迭代器
        // &mut request: 可变借用request，允许中间件修改请求
        for middleware in &self.middlewares {
            // middleware.process_request(): 调用中间件的处理方法
            // ?操作符：如果中间件处理失败，立即返回错误
            middleware.process_request(&mut request)?;
        }

        // 重试循环 - 展示Rust的模式匹配和错误处理
        // for attempt in 0..=self.retry_config.max_retries: 范围迭代器
        //   - 0..=n: 包含端点的范围，从0到n（包括n）
        //   - 如果max_retries=3，则尝试0,1,2,3共4次
        for attempt in 0..=self.retry_config.max_retries {
            // 执行请求：&request借用，不转移所有权
            // match表达式：Rust的模式匹配，必须处理所有可能的情况
            match self.execute_request(&request).await {
                // Ok(mut response): 请求成功，获取响应的可变所有权
                // mut关键字允许中间件修改响应
                Ok(mut response) => {
                    // 处理响应中间件：遍历所有中间件
                    // &self.middlewares: 借用中间件列表，不获取所有权
                    for middleware in &self.middlewares {
                        // &mut response: 可变借用响应，允许中间件修改
                        // ?操作符：如果中间件处理失败，立即返回错误
                        middleware.process_response(&mut response)?;
                    }
                    
                    // return Ok(response): 成功时立即返回响应
                    // response的所有权被转移给调用者
                    return Ok(response);
                }
                // Err(e): 请求失败，e是错误值
                Err(e) => {
                    // 检查是否应该重试：最后一次尝试或不可重试的错误
                    if attempt == self.retry_config.max_retries {
                        // 最后一次尝试失败，返回错误
                        // e的所有权被转移给调用者
                        return Err(e);
                    }
                    
                    // 记录重试警告：warn!宏用于输出警告日志
                    // attempt + 1: 显示人类友好的尝试次数（从1开始）
                    warn!("Request failed on attempt {}, retrying...", attempt + 1);
                    
                    // 等待重试延迟：tokio::time::sleep异步睡眠
                    // self.retry_config.retry_delay: 借用重试延迟配置
                    // .await: 等待睡眠完成，让出CPU给其他任务
                    tokio::time::sleep(self.retry_config.retry_delay).await;
                }
            }
        }

        // 理论上不会到达这里，因为上面的循环总是会return
        // 但为了满足Rust的类型检查，提供一个默认错误
        Err(HttpClientError::RequestFailed("Unexpected error: retry loop completed without return".to_string()))
    }

    // ============================================================================
    // 便捷方法：简化常用HTTP操作
    // ============================================================================
    
    // 便捷方法：发送GET请求 - 展示字符串切片和所有权
    // url: &str: 字符串切片，借用字符串数据而不获取所有权
    //   - &str是对字符串的不可变引用，指向内存中的字符串数据
    //   - 相比String，&str更轻量，避免了不必要的内存分配
    //   - 调用者保持对原始字符串的所有权
    pub async fn get(&self, url: &str) -> Result<HttpResponse> {
        // HttpRequest::new(): 创建新的HTTP请求
        // HttpMethod::GET: 枚举值，实现了Copy trait，可以直接复制
        // url: &str被传递给new方法，在内部会被转换为String
        let request = HttpRequest::new(HttpMethod::GET, url)?;
        // self.send(request): 调用send方法，request的所有权被转移
        self.send(request).await
    }

    // 异步POST请求便捷方法
    pub async fn post(&self, url: &str) -> Result<HttpResponse> {
        let request = HttpRequest::new(HttpMethod::POST, url)?;
        self.send(request).await
    }

    // 便捷方法：发送POST JSON请求 - 展示泛型约束和引用传递
    // <T: Serialize>: 泛型约束，T必须实现Serialize trait
    //   - 这是编译时约束，确保T可以被序列化为JSON
    //   - Serialize来自serde库，是序列化的标准trait
    // data: &T: 借用泛型类型T的数据
    //   - 使用引用避免获取data的所有权，调用者可以继续使用data
    //   - &T表示对任何实现Serialize的类型的不可变引用
    pub async fn post_json<T: Serialize>(&self, url: &str, data: &T) -> Result<HttpResponse> {
        // HttpRequest::new(): 创建POST请求
        // .json(data): 调用json方法，将data序列化为JSON字符串
        //   - data: &T被传递给json方法，在内部进行序列化
        //   - ?操作符处理序列化可能的错误
        let request = HttpRequest::new(HttpMethod::POST, url)?
            .json(data)?; // 链式调用设置JSON数据
        // self.send(request): 转移request的所有权并发送
        self.send(request).await
    }
}

// ============================================================================
// 构建器模式实现：提供更灵活的客户端配置方式
// ============================================================================

// HTTP客户端构建器结构体
pub struct HttpClientBuilder {
    pool_size: usize,                           // 连接池大小
    timeout: Duration,                          // 超时时间
    retry_config: RetryConfig,                  // 重试配置
    middlewares: Vec<Box<dyn Middleware>>,      // 中间件列表
}

impl HttpClientBuilder {
    // 创建默认配置的构建器
    pub fn new() -> Self {
        HttpClientBuilder {
            pool_size: 10,                          // 默认10个连接
            timeout: Duration::from_secs(30),       // 默认30秒超时
            retry_config: RetryConfig::default(),   // 默认重试配置
            middlewares: Vec::new(),                // 空中间件列表
        }
    }

    // 构建器模式：设置连接池大小
    pub fn pool_size(mut self, size: usize) -> Self {
        self.pool_size = size;
        self // 返回self支持链式调用
    }

    // 构建器模式：设置超时时间
    pub fn timeout(mut self, timeout: Duration) -> Self {
        self.timeout = timeout;
        self
    }

    // 构建器模式：设置重试配置
    pub fn retry_config(mut self, config: RetryConfig) -> Self {
        self.retry_config = config;
        self
    }
    }

    // 构建器模式：添加中间件
    // <M: Middleware + 'static>: 泛型约束，M必须实现Middleware trait且具有'static生命周期
    pub fn add_middleware<M: Middleware + 'static>(mut self, middleware: M) -> Self {
        // 将中间件装箱并添加到列表中
        self.middlewares.push(Box::new(middleware));
        self // 返回self支持链式调用
    }

    // 构建器模式：构建最终的HttpClient实例
    pub fn build(self) -> HttpClient {
        // 使用构建器模式创建客户端
        let mut client = HttpClient::new()
            .with_pool_size(self.pool_size)      // 设置连接池大小
            .with_timeout(self.timeout)          // 设置超时时间
            .with_retry_config(self.retry_config); // 设置重试配置

        // 将构建器中的中间件转移到客户端
        for middleware in self.middlewares {
            client.middlewares.push(middleware);
        }

        client // 返回配置完成的客户端
    }


// ============================================================================
// 示例数据结构：用于演示JSON序列化和反序列化
// ============================================================================

// 用户数据结构体
// #[derive(Serialize, Deserialize, Debug)]: 自动实现序列化、反序列化和调试打印
#[derive(Serialize, Deserialize, Debug)]
struct User {
    id: u32,        // 用户ID
    name: String,   // 用户姓名
    email: String,  // 用户邮箱
}

// 通用API响应结构体
// <T>: 泛型参数，允许data字段存储不同类型的数据
#[derive(Serialize, Deserialize, Debug)]
struct ApiResponse<T> {
    message: String,    // 响应消息
    data: Option<T>,    // 响应数据，Option表示可能为空
}

// ============================================================================
// 主函数：演示HTTP客户端的使用
// ============================================================================

// #[tokio::main]: Tokio提供的宏，将async main函数转换为同步入口点
// 这个宏会创建Tokio运行时并执行异步main函数
#[tokio::main]
async fn main() -> Result<()> {
    // 初始化环境变量日志记录器（env_logger库）
    // 可以通过RUST_LOG环境变量控制日志级别
    env_logger::init();

    // 使用println!宏输出程序标题
    println!("🚀 HTTP客户端库演示");
    println!("==================");

    // 使用构建器模式创建HTTP客户端
    let client = HttpClientBuilder::new()
        .pool_size(5)                           // 设置连接池大小为5
        .timeout(Duration::from_secs(10))       // 设置超时时间为10秒
        .retry_config(RetryConfig {             // 自定义重试配置
            max_retries: 2,                     // 最大重试2次
            retry_delay: Duration::from_millis(500), // 重试间隔500毫秒
            retry_on_status: vec![500, 502, 503], // 在这些状态码时重试
        })
        .add_middleware(LoggingMiddleware::new()) // 添加日志中间件
        .add_middleware(AuthMiddleware::bearer("your-api-token".to_string())) // 添加Bearer认证中间件
        .build(); // 构建客户端实例

    // ============================================================================
    // 1. 基本GET请求演示
    // ============================================================================
    println!("\n📡 执行GET请求...");
    // match表达式处理异步请求结果
    match client.get("https://api.example.com/users").await {
        Ok(response) => {
            // 请求成功，打印状态码和响应内容
            println!("✅ GET请求成功: Status {}", response.status);
            println!("📄 响应内容: {}", response.body);
        }
        Err(e) => {
            // 请求失败，打印错误信息
            // {:?}是Debug格式化，显示详细的错误信息
            println!("❌ GET请求失败: {:?}", e);
        }
    }

    // ============================================================================
    // 2. POST JSON请求演示
    // ============================================================================
    println!("\n📤 执行POST JSON请求...");
    // 创建用户数据结构体实例
    let new_user = User {
        id: 0,                                  // 新用户ID为0
        name: "张三".to_string(),               // 用户姓名
        email: "zhangsan@example.com".to_string(), // 用户邮箱
    };

    // 发送POST JSON请求，&new_user是借用引用
    match client.post_json("https://api.example.com/users", &new_user).await {
        Ok(response) => {
            println!("✅ POST请求成功: Status {}", response.status);
            println!("📄 响应内容: {}", response.body);
        }
        Err(e) => {
            println!("❌ POST请求失败: {:?}", e);
        }
    }

    // ============================================================================
    // 3. 自定义请求演示
    // ============================================================================
    println!("\n🔧 执行自定义请求...");
    // 使用构建器模式创建自定义请求
    let custom_request = HttpRequest::new(HttpMethod::GET, "https://api.example.com/data")?
        .header("X-Custom-Header", "custom-value") // 添加自定义头部
        .timeout(Duration::from_secs(5));          // 设置5秒超时

    // 发送自定义请求
    match client.send(custom_request).await {
        Ok(response) => {
            println!("✅ 自定义请求成功: Status {}", response.status);
            // 检查响应是否成功（状态码200-299）
            if response.is_success() {
                println!("🎉 请求处理成功!");
            }
        }
        Err(e) => {
            println!("❌ 自定义请求失败: {:?}", e);
        }
    }

    // ============================================================================
    // 4. 测试错误处理和重试机制
    // ============================================================================
    println!("\n🔄 测试错误处理和重试...");
    // 故意请求一个会返回错误的URL来测试重试机制
    match client.get("https://api.example.com/error").await {
        Ok(response) => {
            println!("📊 错误请求响应: Status {}", response.status);
        }
        Err(e) => {
            println!("❌ 错误请求最终失败: {:?}", e);
        }
    }

    // ============================================================================
    // 5. 连接池状态检查
    // ============================================================================
    println!("\n🏊 连接池状态:");
    // 显示当前可用的连接数
    println!("可用连接数: {}", client.pool.available_connections());

    println!("\n✨ HTTP客户端库演示完成!");
    Ok(()) // 返回成功结果
}

// ============================================================================
// 测试模块：使用Rust内置的测试框架
// ============================================================================

// #[cfg(test)]: 条件编译属性，只在测试时编译这个模块
// cfg = configuration，test表示测试配置
#[cfg(test)]
mod tests {
    use super::*; // 导入父模块的所有公共项

    // ============================================================================
    // 测试HTTP请求创建功能
    // ============================================================================
    
    // #[tokio::test]: Tokio提供的异步测试宏
    // 将异步测试函数转换为可以在测试环境中运行的同步函数
    #[tokio::test]
    async fn test_http_request_creation() {
        // 测试正常的HTTP请求创建
        let request = HttpRequest::new(HttpMethod::GET, "https://example.com").unwrap();
        // assert_eq!宏：断言两个值相等，如果不相等测试失败
        assert_eq!(request.method.to_string(), "GET");
        assert_eq!(request.url, "https://example.com");
    }

    // ============================================================================
    // 测试无效URL处理
    // ============================================================================
    
    #[tokio::test]
    async fn test_invalid_url() {
        // 测试无效URL是否正确返回错误
        let result = HttpRequest::new(HttpMethod::GET, "invalid-url");
        // assert!宏：断言表达式为true，这里检查结果是否为错误
        assert!(result.is_err());
    }

    // ============================================================================
    // 测试请求构建器模式
    // ============================================================================
    
    #[tokio::test]
    async fn test_request_builder() {
        // 测试链式调用构建请求
        let request = HttpRequest::new(HttpMethod::POST, "https://example.com")
            .unwrap()                                    // 解包Result，如果是错误会panic
            .header("Content-Type", "application/json") // 添加头部
            .body("test body".to_string());             // 设置请求体
        
        // 验证头部是否正确设置
        // .get()方法返回Option<&String>，.unwrap()解包获取值
        assert_eq!(request.headers.get("Content-Type").unwrap(), "application/json");
        // 验证请求体是否正确设置
        assert_eq!(request.body.unwrap(), "test body");
    }

    // ============================================================================
    // 测试JSON序列化功能
    // ============================================================================
    
    #[tokio::test]
    async fn test_json_serialization() {
        // 创建测试用户数据
        let user = User {
            id: 1,
            name: "Test".to_string(),
            email: "test@example.com".to_string(),
        };

        // 测试JSON序列化
        let request = HttpRequest::new(HttpMethod::POST, "https://example.com")
            .unwrap()
            .json(&user)    // 将用户数据序列化为JSON
            .unwrap();

        // 验证请求体不为空（包含序列化的JSON数据）
        assert!(request.body.is_some());
        // 验证Content-Type头部被自动设置
        assert_eq!(request.headers.get("Content-Type").unwrap(), "application/json");
    }

    // ============================================================================
    // 测试连接池功能
    // ============================================================================
    
    #[tokio::test]
    async fn test_connection_pool() {
        // 创建一个最大连接数为2的连接池
        let pool = ConnectionPool::new(2);
        // 验证初始可用连接数
        assert_eq!(pool.available_connections(), 2);

        // 获取第一个连接，_guard1变量名前的下划线表示我们不直接使用这个变量
        // 但需要保持它的生命周期以维持连接的占用状态
        let _guard1 = pool.acquire().await.unwrap();
        // 验证获取一个连接后，可用连接数减1
        assert_eq!(pool.available_connections(), 1);

        // 获取第二个连接
        let _guard2 = pool.acquire().await.unwrap();
        // 验证获取两个连接后，可用连接数为0
        assert_eq!(pool.available_connections(), 0);
        
        // 当_guard1和_guard2离开作用域时，连接会自动释放（RAII模式）
    }

    // ============================================================================
    // 测试中间件功能
    // ============================================================================
    
    #[tokio::test]
    async fn test_middleware() {
        // 创建一个HTTP请求用于测试
        let mut request = HttpRequest::new(HttpMethod::GET, "https://example.com").unwrap();
        // 创建Bearer Token认证中间件
        let auth_middleware = AuthMiddleware::bearer("test-token".to_string());
        
        // 处理请求，中间件会自动添加Authorization头部
        auth_middleware.process_request(&mut request).unwrap();
        
        // 验证Authorization头部是否正确设置
        assert_eq!(
            request.headers.get("Authorization").unwrap(),
            "Bearer test-token"
        );
    }

    // ============================================================================
    // 测试客户端构建器模式
    // ============================================================================
    
    #[tokio::test]
    async fn test_client_builder() {
        // 使用构建器模式创建HTTP客户端
        let client = HttpClientBuilder::new()
            .pool_size(5)                                    // 设置连接池大小为5
            .timeout(Duration::from_secs(10))               // 设置超时时间为10秒
            .add_middleware(LoggingMiddleware::new())       // 添加日志中间件
            .build();                                       // 构建最终的客户端实例

        // 验证连接池大小设置是否正确
        assert_eq!(client.pool.available_connections(), 5);
        // 验证超时时间设置是否正确
        assert_eq!(client.default_timeout, Duration::from_secs(10));
        // 验证中间件数量是否正确
        assert_eq!(client.middlewares.len(), 1);
    }
}