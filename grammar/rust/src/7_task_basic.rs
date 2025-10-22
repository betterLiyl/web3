// ### 任务1：安全的内存管理器
// **目标**：掌握所有权、借用、生命周期
// **描述**：实现一个简单的内存池管理器，演示Rust的内存安全特性

// **流程提示**：
// 1. 定义内存块结构体，包含大小和状态信息
// 2. 实现内存池结构，使用Vec管理内存块
// 3. 实现分配和释放方法，注意所有权转移 ✓
// 4. 使用生命周期参数确保引用安全 ✓
// 5. 添加统计功能（已用内存、碎片等）

use std::time;
use std::rc::Rc;
use std::cell::RefCell;

// 内存块结构体
#[derive(Debug)]
struct MemoryBlock {
    id: usize,
    size: usize,
    is_free: bool,
    data: Vec<u8>,
}

// 并发日志系统测试演示
fn test_concurrent_logging() {
    println!("\n=== 并发日志系统演示 ===");
    
    // 1. 创建默认配置的日志系统
    println!("1. 创建日志系统");
    let logger = Logger::new(LogConfig::default()).expect("创建日志系统失败");
    
    // 2. 测试不同级别的日志
    println!("2. 测试不同级别的日志");
    logger.log(LogLevel::Debug, "这是一条调试信息".to_string()).unwrap();
    logger.log(LogLevel::Info, "应用程序启动成功".to_string()).unwrap();
    logger.log(LogLevel::Warn, "这是一个警告信息".to_string()).unwrap();
    logger.log(LogLevel::Error, "发生了一个错误".to_string()).unwrap();
    
    // 3. 测试并发日志写入
    println!("3. 测试并发日志写入");
    use std::sync::Arc;
    let logger_arc = Arc::new(logger);
    let mut handles = vec![];
    
    for i in 0..5 {
        let logger_clone = Arc::clone(&logger_arc);
        let handle = thread::spawn(move || {
            for j in 0..3 {
                let message = format!("线程 {} 的第 {} 条消息", i, j + 1);
                logger_clone.log(LogLevel::Info, message).unwrap();
                thread::sleep(std::time::Duration::from_millis(10));
            }
        });
        handles.push(handle);
    }
    
    // 等待所有线程完成
    for handle in handles {
        handle.join().unwrap();
    }
    
    // 4. 测试配置更新
    println!("4. 测试配置更新（设置为只记录错误级别）");
    let mut new_config = LogConfig::default();
    new_config.min_level = LogLevel::Error;
    logger_arc.update_config(new_config);
    
    // 这些日志不会被记录
    logger_arc.log(LogLevel::Debug, "这条调试信息不会被记录".to_string()).unwrap();
    logger_arc.log(LogLevel::Info, "这条信息不会被记录".to_string()).unwrap();
    logger_arc.log(LogLevel::Warn, "这条警告不会被记录".to_string()).unwrap();
    
    // 只有错误级别会被记录
    logger_arc.log(LogLevel::Error, "只有这条错误会被记录".to_string()).unwrap();
    
    // 5. 测试日志轮转（创建一个小文件大小的配置）
    println!("5. 测试日志轮转");
    let mut rotation_config = LogConfig::default();
    rotation_config.max_file_size = 100; // 很小的文件大小来触发轮转
    rotation_config.log_dir = "test_logs".to_string();
    
    let rotation_logger = Logger::new(rotation_config).expect("创建轮转日志系统失败");
    
    // 写入足够多的日志来触发轮转
    for i in 0..10 {
        let long_message = format!("这是一条很长的日志消息用来测试日志轮转功能 - 消息编号: {}", i);
        rotation_logger.log(LogLevel::Info, long_message).unwrap();
    }
    
    // 等待一下让日志处理完成
    thread::sleep(std::time::Duration::from_millis(100));
    
    println!("=== 并发日志系统演示完成 ===");
    println!("日志文件已写入到 logs/ 和 test_logs/ 目录");
}

// 内存句柄 - 体现所有权转移和RAII
pub struct MemoryHandle {
    id: usize,
    data: Vec<u8>,
    pool: Rc<RefCell<MemoryPool>>,
}

impl MemoryHandle {
    // 获取内存数据的可变引用 - 体现借用规则
    pub fn data_mut(&mut self) -> &mut [u8] {
        &mut self.data
    }
    
    // 获取内存数据的不可变引用
    pub fn data(&self) -> &[u8] {
        &self.data
    }
    
    // 获取内存块大小
    pub fn size(&self) -> usize {
        self.data.len()
    }
    
    // 获取内存块ID
    pub fn id(&self) -> usize {
        self.id
    }
}

// 当MemoryHandle被drop时，自动释放内存 - 体现RAII和所有权
impl Drop for MemoryHandle {
    fn drop(&mut self) {
        println!("内存句柄 {} 被释放，内存块自动回收", self.id);
        if let Ok(mut pool) = self.pool.try_borrow_mut() {
            pool.free_by_id(self.id);
        }
    }
}

// 内存池结构体
#[derive(Debug)]
pub struct MemoryPool {
    blocks: Vec<MemoryBlock>,
    next_id: usize,
    total_allocated: usize,
}

impl MemoryPool {
    pub fn new() -> Rc<RefCell<Self>> {
        Rc::new(RefCell::new(Self {
            blocks: Vec::new(),
            next_id: 1,
            total_allocated: 0,
        }))
    }
    
    // 添加新的内存块到池中
    pub fn add_block(&mut self, size: usize) {
        self.blocks.push(MemoryBlock {
            id: self.next_id,
            size,
            is_free: true,
            data: vec![0; size],
        });
        self.next_id += 1;
    }
    
    // 分配内存 - 返回拥有所有权的内存句柄
    pub fn allocate(pool: &Rc<RefCell<Self>>, size: usize) -> Option<MemoryHandle> {
        let mut pool_ref = pool.borrow_mut();
        
        // 寻找合适的空闲块
        let mut found_block_id = None;
        for block in pool_ref.blocks.iter() {
            if block.is_free && block.size >= size {
                found_block_id = Some(block.id);
                break;
            }
        }
        
        if let Some(block_id) = found_block_id {
            // 现在修改找到的块
            for block in pool_ref.blocks.iter_mut() {
                if block.id == block_id {
                    block.is_free = false;
                    pool_ref.total_allocated += size;
                    
                    // 创建数据副本 - 所有权转移给句柄
                    let data = vec![0; size];
                    
                    drop(pool_ref); // 释放借用
                    
                    return Some(MemoryHandle {
                        id: block_id,
                        data,
                        pool: Rc::clone(pool),
                    });
                }
            }
        }
        
        None
    }
    
    // 通过ID释放内存块
    fn free_by_id(&mut self, id: usize) {
        for block in self.blocks.iter_mut() {
            if block.id == id && !block.is_free {
                block.is_free = true;
                self.total_allocated = self.total_allocated.saturating_sub(block.size);
                println!("   内存块 {} 已释放", id);
                return;
            }
        }
    }
    
    // 手动释放内存块
    pub fn free(&mut self, id: usize) -> Result<(), &'static str> {
        for block in self.blocks.iter_mut() {
            if block.id == id {
                if !block.is_free {
                    block.is_free = true;
                    self.total_allocated = self.total_allocated.saturating_sub(block.size);
                    return Ok(());
                } else {
                    return Err("内存块已经是空闲状态");
                }
            }
        }
        Err("无效的内存块ID")
    }
    
    // 获取内存块信息 - 体现借用规则
    pub fn get_block_info(&self, id: usize) -> Option<&MemoryBlock> {
        self.blocks.iter().find(|block| block.id == id)
    }
    
    // 统计信息 - 体现借用规则（不可变借用）
    pub fn stats(&self) -> (usize, usize, usize, usize) {
        let (mut used, mut free, mut fragmented) = (0, 0, 0);
        
        for block in self.blocks.iter() {
            if block.is_free {
                free += 1;
            } else {
                used += 1;
                // 简化的碎片检测
                if block.data.len() > block.size {
                    fragmented += 1;
                }
            }
        }
        
        (used, free, fragmented, self.total_allocated)
    }
}

fn main() {
    println!("=== 内存管理器演示：所有权转移和生命周期 ===");
    
    // 创建内存池
    let pool = MemoryPool::new();
    
    // 添加一些内存块到池中
    {
        let mut pool_ref = pool.borrow_mut();
        pool_ref.add_block(1024);
        pool_ref.add_block(512);
        pool_ref.add_block(256);
    }
    
    // 1. 演示内存分配 - 所有权转移
    {
        println!("\n1. 分配内存块 (512 字节)");
        if let Some(mut handle) = MemoryPool::allocate(&pool, 512) {
            println!("   分配成功，内存块ID: {}, 大小: {} 字节", handle.id(), handle.size());
            
            // 使用内存句柄 - 体现借用规则
            {
                let data = handle.data_mut();
                data[0] = 42;
                data[1] = 100;
                println!("   写入数据: [{}, {}]", data[0], data[1]);
            }
            
            // 读取数据
            let read_data = handle.data();
            println!("   读取数据: [{}, {}]", read_data[0], read_data[1]);
            
            println!("   内存句柄即将被释放...");
        }
    } // handle在这里被drop，触发自动释放
    
    // 2. 演示多个内存分配的生命周期管理
    {
        println!("\n2. 分配多个内存块");
        
        let handle1 = MemoryPool::allocate(&pool, 256);
        let handle2 = MemoryPool::allocate(&pool, 128);
        
        if let (Some(mut h1), Some(mut h2)) = (handle1, handle2) {
            println!("   成功分配两个内存块: ID {} ({} 字节) 和 ID {} ({} 字节)", 
                    h1.id(), h1.size(), h2.id(), h2.size());
            
            // 分别操作两个内存块
            h1.data_mut()[0] = 1;
            h2.data_mut()[0] = 2;
            
            println!("   内存块{}数据: {}", h1.id(), h1.data()[0]);
            println!("   内存块{}数据: {}", h2.id(), h2.data()[0]);
        }
    } // 两个句柄都在这里被drop
    
    // 3. 演示手动释放和错误处理
    {
        println!("\n3. 手动释放内存块");
        if let Some(handle) = MemoryPool::allocate(&pool, 100) {
            let id = handle.id();
            drop(handle); // 显式释放句柄
            
            // 尝试手动释放（这会失败，因为已经自动释放了）
            let mut pool_ref = pool.borrow_mut();
            match pool_ref.free(id) {
                Ok(_) => println!("   手动释放成功"),
                Err(e) => println!("   手动释放失败: {}", e),
            }
        }
    }
    
    // 4. 演示借用检查 - 获取内存块信息
    {
        println!("\n4. 查看内存块信息（借用检查）");
        let pool_ref = pool.borrow();
        if let Some(block_info) = pool_ref.get_block_info(1) {
            println!("   内存块1 - ID: {}, 大小: {}, 空闲: {}", 
                    block_info.id, block_info.size, block_info.is_free);
        }
    }
    
    // 5. 显示统计信息
    {
        println!("\n5. 内存池统计信息");
        let pool_ref = pool.borrow();
        let (used, free, fragmented, total_allocated) = pool_ref.stats();
        println!("   已使用块: {}", used);
        println!("   空闲块: {}", free);
        println!("   碎片块: {}", fragmented);
        println!("   总分配内存: {} 字节", total_allocated);
    }
    
    // 6. 演示作用域和自动释放
    {
        println!("\n6. 作用域演示");
        let _handle = MemoryPool::allocate(&pool, 64);
        if let Some(h) = _handle {
            println!("   分配64字节内存块，ID: {}", h.id());
            // handle会在作用域结束时自动释放
        }
    }
    
    println!("\n=== 演示完成 ===");
    println!("所有内存句柄都已正确释放，展示了Rust的内存安全特性：");
    println!("- 所有权转移：内存数据的所有权从pool转移到handle");
    println!("- 借用规则：可变和不可变引用的正确使用");
    println!("- 生命周期：通过Rc<RefCell<>>管理共享所有权");
    println!("- RAII：资源在离开作用域时自动释放");
    
    // 调用配置解析器测试
    test_config_parser();
    
    // 调用并发日志系统测试
    test_concurrent_logging();
}



// ### 任务2：配置文件解析器
// **目标**：掌握模式匹配、枚举、错误处理
// **描述**：实现一个支持多种格式的配置文件解析器（JSON、TOML、YAML）

// **流程提示**：
// 1. 定义配置值枚举（String、Number、Boolean、Array、Object）
// 2. 实现不同格式的解析器trait
// 3. 使用模式匹配处理不同的配置值类型
// 4. 实现自定义错误类型和错误传播
// 5. 添加配置验证和默认值功能
use std::collections::HashMap;
#[derive(Debug, Clone, Copy)]
enum ConfigFormat {
    Json,
    Toml,
    Yaml,
}

#[derive(Debug)]
enum ConfigValueEnum<T> {
    String(String),
    Number(f64),
    Boolean(bool),
    Array(Vec<T>),
    Object(HashMap<String, T>),
}

trait ConfigParser<T> {
    fn parse(&self, input: &str) -> Result<T, ConfigError>;
}

#[derive(Debug)]
struct ConfigError {
    message: String,
}

impl std::fmt::Display for ConfigError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "配置错误: {}", self.message)
    }
}

impl std::error::Error for ConfigError {}

impl ConfigError {
    fn new(message: &str) -> Self {
        Self {
            message: message.to_string(),
        }
    }
}

struct JsonParser;
struct TomlParser;
struct YamlParser;

impl ConfigParser<ConfigValueEnum<String>> for JsonParser {
    fn parse(&self, input: &str) -> Result<ConfigValueEnum<String>, ConfigError> {
        // 实现JSON解析逻辑
        Ok(ConfigValueEnum::String(input.to_string()))
    }
}

impl ConfigParser<ConfigValueEnum<String>> for TomlParser {
    fn parse(&self, input: &str) -> Result<ConfigValueEnum<String>, ConfigError> {
        // 实现TOML解析逻辑
        Ok(ConfigValueEnum::String(input.to_string()))
    }
}
impl ConfigParser<ConfigValueEnum<String>> for YamlParser {
    fn parse(&self, input: &str) -> Result<ConfigValueEnum<String>, ConfigError> {
        // 实现YAML解析逻辑
        Ok(ConfigValueEnum::String(input.to_string()))
    }
}

fn parse_config<T>(format: ConfigFormat, input: &str) -> Result<T, ConfigError>
where
    T: for<'a> TryFrom<&'a ConfigValueEnum<String>, Error = ConfigError>,
{
    let result = match format {
        ConfigFormat::Json => JsonParser.parse(input)?,
        ConfigFormat::Toml => TomlParser.parse(input)?,
        ConfigFormat::Yaml => YamlParser.parse(input)?,
    };
    
    T::try_from(&result).map_err(|e| ConfigError::new(&format!("类型转换失败: {}", e)))
}

// 辅助函数：直接解析为ConfigValueEnum
fn parse_config_value(format: ConfigFormat, input: &str) -> Result<ConfigValueEnum<String>, ConfigError> {
    match format {
        ConfigFormat::Json => JsonParser.parse(input),
        ConfigFormat::Toml => TomlParser.parse(input),
        ConfigFormat::Yaml => YamlParser.parse(input),
    }
}

// 配置验证trait
trait ConfigValidator {
    fn validate(&self) -> Result<(), ConfigError>;
}

// 为ConfigValueEnum实现验证
impl ConfigValidator for ConfigValueEnum<String> {
    fn validate(&self) -> Result<(), ConfigError> {
        match self {
            ConfigValueEnum::String(s) if s.is_empty() => {
                Err(ConfigError::new("字符串值不能为空"))
            }
            ConfigValueEnum::Array(arr) if arr.is_empty() => {
                Err(ConfigError::new("数组不能为空"))
            }
            ConfigValueEnum::Object(obj) if obj.is_empty() => {
                Err(ConfigError::new("对象不能为空"))
            }
            _ => Ok(()),
        }
    }
}

// 默认值trait
trait DefaultConfig {
    fn default_value() -> Self;
}

impl DefaultConfig for ConfigValueEnum<String> {
    fn default_value() -> Self {
        ConfigValueEnum::String("default".to_string())
    }
}

// 示例：为String类型实现TryFrom，以便可以从ConfigValueEnum转换
impl TryFrom<&ConfigValueEnum<String>> for String {
    type Error = ConfigError;
    
    fn try_from(value: &ConfigValueEnum<String>) -> Result<Self, Self::Error> {
        match value {
            ConfigValueEnum::String(s) => Ok(s.clone()),
            _ => Err(ConfigError::new("无法将配置值转换为字符串")),
        }
    }
}

// 配置解析器的完整示例和测试
fn test_config_parser() {
    println!("\n=== 配置文件解析器演示 ===");
    
    // 1. 测试基本解析功能
    println!("\n1. 基本解析测试");
    let json_input = r#"{"name": "test", "version": "1.0"}"#;
    
    match parse_config_value(ConfigFormat::Json, json_input) {
        Ok(config) => {
            println!("   JSON解析成功: {:?}", config);
            
            // 验证配置
            match config.validate() {
                Ok(_) => println!("   配置验证通过"),
                Err(e) => println!("   配置验证失败: {}", e),
            }
        }
        Err(e) => println!("   JSON解析失败: {}", e),
    }
    
    // 2. 测试类型转换
    println!("\n2. 类型转换测试");
    let config_value = ConfigValueEnum::String("Hello World".to_string());
    
    match String::try_from(&config_value) {
        Ok(s) => println!("   转换为字符串成功: {}", s),
        Err(e) => println!("   转换失败: {}", e),
    }
    
    // 3. 测试错误处理
    println!("\n3. 错误处理测试");
    let empty_config = ConfigValueEnum::String("".to_string());
    
    match empty_config.validate() {
        Ok(_) => println!("   空字符串验证通过"),
        Err(e) => println!("   空字符串验证失败: {}", e),
    }
    
    // 4. 测试默认值
    println!("\n4. 默认值测试");
    let default_config = ConfigValueEnum::<String>::default_value();
    println!("   默认配置值: {:?}", default_config);
    
    // 5. 测试不同格式
    println!("\n5. 多格式解析测试");
    let formats = [
        (ConfigFormat::Json, "JSON"),
        (ConfigFormat::Toml, "TOML"),
        (ConfigFormat::Yaml, "YAML"),
    ];
    
    for (format, name) in formats.iter() {
        match parse_config_value(*format, "test_config") {
            Ok(config) => println!("   {} 解析成功: {:?}", name, config),
            Err(e) => println!("   {} 解析失败: {}", name, e),
        }
    }
    
    println!("\n=== 配置解析器演示完成 ===");
}


// ### 任务3：并发日志系统
// **目标**：掌握并发编程、通道、线程
// **描述**：实现一个线程安全的日志系统，支持不同级别和异步写入

// **流程提示**：
// 1. 定义日志级别枚举和日志条目结构
// 2. 实现日志格式化器trait
// 3. 使用通道实现异步日志写入
// 4. 创建后台线程处理日志写入
// 5. 实现日志轮转和文件管理
// 6. 添加配置和过滤功能
#[derive(Debug, Clone)]
enum LogLevel {
    Debug,
    Info,
    Warn,
    Error,
}

impl std::fmt::Display for LogLevel {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            LogLevel::Debug => write!(f, "DEBUG"),
            LogLevel::Info => write!(f, "INFO"),
            LogLevel::Warn => write!(f, "WARN"),
            LogLevel::Error => write!(f, "ERROR"),
        }
    }
}

struct LogEntry {
    level: LogLevel,
    message: String,
    timestamp: time::SystemTime,
    dest_path: String,
}

trait LogFormatter {
    fn format(&self, entry: &LogEntry) -> String;
}

struct DebugLogFormatter;
struct InfoLogFormatter;

struct WarnLogFormatter;
struct ErrorLogFormatter;

impl LogFormatter for DebugLogFormatter {
    fn format(&self, entry: &LogEntry) -> String {
        let timestamp = entry.timestamp.duration_since(time::UNIX_EPOCH)
            .expect("时间戳错误")
            .as_secs();
        format!("{} [{}] {}", timestamp, entry.level.to_string(), entry.message)
    }
}

impl LogFormatter for InfoLogFormatter {
    fn format(&self, entry: &LogEntry) -> String {
        let timestamp = entry.timestamp.duration_since(time::UNIX_EPOCH)
            .expect("时间戳错误")
            .as_secs();
        format!("{} [{}] {}", timestamp, entry.level.to_string(), entry.message)
    }
}

impl LogFormatter for WarnLogFormatter {
    fn format(&self, entry: &LogEntry) -> String {
        let timestamp = entry.timestamp.duration_since(time::UNIX_EPOCH)
            .expect("时间戳错误")
            .as_secs();
        format!("{} [{}] {}", timestamp, entry.level.to_string(), entry.message)
    }
}

impl LogFormatter for ErrorLogFormatter {
    fn format(&self, entry: &LogEntry) -> String {
        let timestamp = entry.timestamp.duration_since(time::UNIX_EPOCH)
            .expect("时间戳错误")
            .as_secs();
        format!("{} [{}] {}", timestamp, entry.level.to_string(), entry.message)
    }
}



use std::sync::mpsc::{self, Sender, Receiver};
use std::thread;
use std::sync::{Arc, Mutex};
use std::path::Path;

// 日志配置结构体
#[derive(Debug, Clone)]
struct LogConfig {
    min_level: LogLevel,
    max_file_size: u64,  // 字节
    max_files: usize,
    log_dir: String,
}

impl Default for LogConfig {
    fn default() -> Self {
        Self {
            min_level: LogLevel::Info,
            max_file_size: 1024 * 1024, // 1MB
            max_files: 5,
            log_dir: "logs".to_string(),
        }
    }
}

// 日志系统主结构体
struct Logger {
    sender: Sender<LogMessage>,
    config: Arc<Mutex<LogConfig>>,
    _worker_handle: thread::JoinHandle<()>,
}

// 日志消息枚举，支持关闭信号
enum LogMessage {
    Entry(LogEntry),
    Shutdown,
}

impl Logger {
    fn new(config: LogConfig) -> std::io::Result<Self> {
        let (sender, receiver) = mpsc::channel();
        let config_arc = Arc::new(Mutex::new(config));
        let worker_config = Arc::clone(&config_arc);
        
        // 创建日志目录
        let log_dir = {
            let config_guard = worker_config.lock().unwrap();
            config_guard.log_dir.clone()
        };
        std::fs::create_dir_all(&log_dir)?;
        
        // 启动后台工作线程
        let worker_handle = thread::spawn(move || {
            Self::log_worker(receiver, worker_config);
        });
        
        Ok(Logger {
            sender,
            config: config_arc,
            _worker_handle: worker_handle,
        })
    }
    
    // 异步写入日志
    fn log(&self, level: LogLevel, message: String) -> Result<(), Box<dyn std::error::Error>> {
        // 检查日志级别过滤
        let config = self.config.lock().unwrap();
        if !Self::should_log(&level, &config.min_level) {
            return Ok(());
        }
        
        let log_dir = config.log_dir.clone();
        drop(config); // 释放锁
        
        let entry = LogEntry {
            level,
            message,
            timestamp: time::SystemTime::now(),
            dest_path: format!("{}/app.log", log_dir),
        };
        
        self.sender.send(LogMessage::Entry(entry))?;
        Ok(())
    }
    
    // 检查是否应该记录日志
    fn should_log(level: &LogLevel, min_level: &LogLevel) -> bool {
        let level_priority = match level {
            LogLevel::Debug => 0,
            LogLevel::Info => 1,
            LogLevel::Warn => 2,
            LogLevel::Error => 3,
        };
        
        let min_priority = match min_level {
            LogLevel::Debug => 0,
            LogLevel::Info => 1,
            LogLevel::Warn => 2,
            LogLevel::Error => 3,
        };
        
        level_priority >= min_priority
    }
    
    // 后台工作线程
    fn log_worker(receiver: Receiver<LogMessage>, config: Arc<Mutex<LogConfig>>) {
        while let Ok(message) = receiver.recv() {
            match message {
                LogMessage::Entry(entry) => {
                    if let Err(e) = Self::process_log_entry(entry, &config) {
                        eprintln!("日志处理错误: {}", e);
                    }
                }
                LogMessage::Shutdown => {
                    println!("日志系统正在关闭...");
                    break;
                }
            }
        }
    }
    
    // 处理单个日志条目
    fn process_log_entry(entry: LogEntry, config: &Arc<Mutex<LogConfig>>) -> std::io::Result<()> {
        let formatter: Box<dyn LogFormatter> = match entry.level {
            LogLevel::Debug => Box::new(DebugLogFormatter),
            LogLevel::Info => Box::new(InfoLogFormatter),
            LogLevel::Warn => Box::new(WarnLogFormatter),
            LogLevel::Error => Box::new(ErrorLogFormatter),
        };
        
        let formatted = formatter.format(&entry);
        
        // 控制台输出
        println!("{}", formatted);
        
        // 检查文件大小并轮转
        let config_guard = config.lock().unwrap();
        Self::rotate_if_needed(&entry.dest_path, &config_guard)?;
        drop(config_guard);
        
        // 写入文件
        use std::fs::OpenOptions;
        use std::io::Write;
        
        let mut file = OpenOptions::new()
            .create(true)
            .append(true)
            .open(&entry.dest_path)?;
        
        writeln!(file, "{}", formatted)?;
        Ok(())
    }
    
    // 日志轮转
    fn rotate_if_needed(log_path: &str, config: &LogConfig) -> std::io::Result<()> {
        if let Ok(metadata) = std::fs::metadata(log_path) {
            if metadata.len() > config.max_file_size {
                Self::rotate_logs(log_path, config)?;
            }
        }
        Ok(())
    }
    
    // 执行日志轮转
    fn rotate_logs(log_path: &str, config: &LogConfig) -> std::io::Result<()> {
        let path = Path::new(log_path);
        let parent = path.parent().unwrap();
        let stem = path.file_stem().unwrap().to_str().unwrap();
        let extension = path.extension().unwrap_or_default().to_str().unwrap();
        
        // 轮转现有文件
        for i in (1..config.max_files).rev() {
            let old_file = parent.join(format!("{}.{}.{}", stem, i, extension));
            let new_file = parent.join(format!("{}.{}.{}", stem, i + 1, extension));
            
            if old_file.exists() {
                if i + 1 >= config.max_files {
                    std::fs::remove_file(&old_file)?; // 删除最老的文件
                } else {
                    std::fs::rename(&old_file, &new_file)?;
                }
            }
        }
        
        // 重命名当前文件
        let backup_file = parent.join(format!("{}.1.{}", stem, extension));
        std::fs::rename(log_path, backup_file)?;
        
        Ok(())
    }
    
    // 更新配置
    fn update_config(&self, new_config: LogConfig) {
        let mut config = self.config.lock().unwrap();
        *config = new_config;
    }
    
    // 优雅关闭
    fn shutdown(self) -> Result<(), Box<dyn std::error::Error>> {
        self.sender.send(LogMessage::Shutdown)?;
        // 注意：这里我们不能等待线程结束，因为会消费self
        // 在实际应用中，可能需要不同的设计来处理这个问题
        Ok(())
    }
}