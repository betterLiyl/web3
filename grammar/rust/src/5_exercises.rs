// Rust 语言练习题
// 请完成以下练习题，每个函数都有详细的说明和示例
// 重点练习 Rust 的所有权、借用、生命周期等核心概念

use std::collections::HashMap;

// ==================== 练习题 ====================

// 练习1: 所有权和移动
// 完成函数，计算字符串的长度但不获取所有权
fn calculate_length(s: &String) -> usize {
    // TODO: 返回字符串的长度
    // 注意：这里使用引用，不会获取所有权
    s.len()
}

// 练习2: 可变借用
// 完成函数，修改字符串内容（添加后缀）
fn add_suffix(s: &mut String, suffix: &str) {
    // TODO: 在字符串末尾添加后缀
    // 提示：使用 push_str 方法
    s.push_str(suffix);
}

// 练习3: 结构体和方法
// 定义一个 Book 结构体
#[derive(Debug, Clone)]
struct Book {
    // TODO: 定义字段
    title: String,
    author: String,
    pages: u32,
    available: bool,
}

impl Book {
    // TODO: 实现构造函数
    fn new(title: String, author: String, pages: u32) -> Self {
        // TODO: 创建新的 Book 实例，默认 available 为 true
        Book {
            title: String::new(),
            author: String::new(),
            pages: 0,
            available: false,
        }
    }
    
    // TODO: 实现获取书籍信息的方法
    fn get_info(&self) -> String {
        // TODO: 返回格式化的书籍信息
        format!("书名: {}, 作者: {}, 页数: {}, 状态: {}", self.title, self.author, self.pages, if self.available {"可用"} else {"已借"})
    }
    
    // TODO: 实现借书方法
    fn borrow_book(&mut self) -> Result<(), String> {
        // TODO: 如果书可用，设置为不可用并返回 Ok(())
        if self.available{
            return Ok(())
        }
        // 如果不可用，返回错误信息
        Err("未实现".to_string())
    }
    
    // TODO: 实现还书方法
    fn return_book(&mut self) {
        // TODO: 设置书为可用状态
        self.available = true;
    }
}

// 练习4: 枚举和模式匹配
// 定义一个表示几何形状的枚举
#[derive(Debug)]
enum Shape {
    // TODO: 定义不同的形状
    Circle(f64),                    // 半径
    Rectangle(f64, f64),           // 宽度, 高度
    Triangle(f64, f64, f64),       // 三边长
}

impl Shape {
    // TODO: 实现计算面积的方法
    fn area(&self) -> f64 {
        // TODO: 使用 match 表达式计算不同形状的面积
        // 圆形: π * r²
        // 矩形: 宽 * 高
        // 三角形: 使用海伦公式
        match self {
            Shape::Circle(r) => std::f64::consts::PI * r * r,
            Shape::Rectangle(w, h) => w * h,
            Shape::Triangle(a, b, c) => {
                let s = (a + b + c) / 2.0;
                (s * (s - a) * (s - b) * (s - c)).sqrt()
            }
        }
    }
    
    // TODO: 实现计算周长的方法
    fn perimeter(&self) -> f64 {
        // TODO: 使用 match 表达式计算不同形状的周长
        match self{
            Shape::Circle(r) => 2.0 * std::f64::consts::PI * r,
            Shape::Rectangle(w, h) => 2.0 * (w + h),
            Shape::Triangle(a, b, c) => a + b + c,
        }
    }
}

// 练习5: 错误处理
// 完成函数，安全地解析字符串为数字
fn parse_number(s: &str) -> Result<i32, String> {
    // TODO: 尝试将字符串解析为 i32
    // 成功时返回 Ok(number)
    // 失败时返回 Err(错误信息)
    // 提示：使用 s.parse::<i32>() 方法
    match s.parse::<i32>() {
        Ok(num) => Ok(num),
        Err(_) => Err("未实现".to_string()),
    }
}

// 练习6: 向量操作
// 完成函数，过滤出偶数并返回新向量
fn filter_even_numbers(numbers: Vec<i32>) -> Vec<i32> {
    // TODO: 使用迭代器和闭包过滤出偶数
    // 提示：使用 into_iter(), filter(), collect()
    numbers.into_iter().filter(|i| i % 2 == 0).collect()
}

// 练习7: HashMap 操作
// 完成函数，统计字符串中每个字符的出现次数
fn count_characters(text: &str) -> HashMap<char, usize> {
    // TODO: 创建 HashMap 统计字符出现次数
    // 提示：遍历字符，使用 entry().or_insert() 或 get_mut()
    let mut res = HashMap::new();
    for c in text.chars(){
        *res.entry(c).or_insert(0) += 1;
    }
    res
}

// 练习8: 生命周期
// 完成函数，返回两个字符串切片中较长的那个
fn longest<'a>(x: &'a str, y: &'a str) -> &'a str {
    // TODO: 比较两个字符串的长度，返回较长的那个
    // 注意生命周期参数的使用
    if x.len() > y.len(){
        x
    }else{
        y
    }
}

// 练习9: 结构体中的引用（生命周期）
// 定义一个包含字符串引用的结构体
struct TextAnalyzer<'a> {
    // TODO: 定义一个包含字符串引用的字段
     text: &'a str,
}

impl<'a> TextAnalyzer<'a> {
    // TODO: 实现构造函数
    fn new(text: &'a str) -> Self {
        TextAnalyzer {
            // TODO: 初始化字段
            text,
        }
    }
    
    // TODO: 实现统计单词数量的方法
    fn word_count(&self) -> usize {
        // TODO: 统计文本中的单词数量
        // 提示：使用 split_whitespace()
        self.text.split_whitespace().count()
    }
    
    // TODO: 实现查找最长单词的方法
    fn longest_word(&self) -> Option<&str> {
        // TODO: 找到最长的单词
        // 提示：使用 split_whitespace(), max_by_key()
        self.text.split_whitespace().max_by_key(|s| s.len())
    }
}

// 练习10: 智能指针（可选，较难）
use std::rc::Rc;
use std::cell::RefCell;

// 定义一个简单的链表节点
#[derive(Debug)]
struct ListNode {
    value: i32,
    next: Option<Rc<RefCell<ListNode>>>,
}

impl ListNode {
    fn new(value: i32) -> Rc<RefCell<Self>> {
        Rc::new(RefCell::new(ListNode { value, next: None }))
    }
    
    // TODO: 实现添加下一个节点的方法
    fn add_next(&mut self, value: i32) {
        // TODO: 创建新节点并设置为当前节点的下一个节点
        // 提示：使用 Rc::new(), RefCell::new()
        self.next = Some(ListNode::new(value));
    }
}

// ==================== 测试函数 ====================

fn test_ownership() {
    println!("=== 测试所有权和借用 ===");
    let s = String::from("Hello, Rust!");
    let len = calculate_length(&s);
    println!("字符串 '{}' 的长度是: {}", s, len);
    
    let mut s2 = String::from("Hello");
    add_suffix(&mut s2, ", World!");
    println!("添加后缀后: {}", s2);
}

fn test_book() {
    println!("\n=== 测试书籍结构体 ===");
    // TODO: 创建书籍实例并测试方法
    let mut book = Book::new("Rust程序设计语言".to_string(), "Steve Klabnik".to_string(), 500);
    println!("书籍信息: {}", book.get_info());
    
    match book.borrow_book() {
        Ok(()) => println!("成功借阅书籍"),
        Err(e) => println!("借阅失败: {}", e),
    }
    
    book.return_book();
    println!("书籍已归还");
}

fn test_shapes() {
    println!("\n=== 测试几何形状 ===");
    // TODO: 创建不同形状并测试面积和周长计算
    let circle = Shape::Circle(5.0);
    let rectangle = Shape::Rectangle(4.0, 6.0);
    let triangle = Shape::Triangle(3.0, 4.0, 5.0);
    
    println!("圆形面积: {:.2}", circle.area());
    println!("矩形面积: {:.2}", rectangle.area());
    println!("三角形面积: {:.2}", triangle.area());
}

fn test_error_handling() {
    println!("\n=== 测试错误处理 ===");
    let test_strings = vec!["42", "hello", "123", ""];
    
    for s in test_strings {
        match parse_number(s) {
            Ok(num) => println!("'{}' 解析为: {}", s, num),
            Err(e) => println!("'{}' 解析失败: {}", s, e),
        }
    }
}

fn test_vector_operations() {
    println!("\n=== 测试向量操作 ===");
    let numbers = vec![1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
    let even_numbers = filter_even_numbers(numbers.clone());
    println!("原始数组: {:?}", numbers);
    println!("偶数: {:?}", even_numbers);
}

fn test_hashmap() {
    println!("\n=== 测试 HashMap ===");
    let text = "hello world";
    let char_count = count_characters(text);
    println!("字符串 '{}' 中各字符出现次数:", text);
    for (ch, count) in char_count {
        println!("  '{}': {}", ch, count);
    }
}

fn test_lifetimes() {
    println!("\n=== 测试生命周期 ===");
    let string1 = "short";
    let string2 = "longer string";
    let result = longest(string1, string2);
    println!("较长的字符串是: '{}'", result);
    
    // TODO: 测试 TextAnalyzer
    let text = "Hello world! This is a test string with multiple words.";
    let analyzer = TextAnalyzer::new(text);
    println!("文本: '{}'", text);
    println!("单词数量: {}", analyzer.word_count());
    if let Some(longest) = analyzer.longest_word() {
        println!("最长单词: '{}'", longest);
    }
}


fn main() {
    println!("Rust 语言练习题");
    println!("请完成上面标有 TODO 的函数实现");
    println!("然后运行测试函数查看结果");
    println!("=====================================");
    
    // 运行测试函数
    test_ownership();
    test_book();
    test_shapes();
    test_error_handling();
    test_vector_operations();
    test_hashmap();
    test_lifetimes();
    
    println!("\n练习完成后，你可以尝试以下进阶挑战：");
    println!("1. 实现一个简单的图书管理系统");
    println!("2. 创建一个命令行计算器");
    println!("3. 使用智能指针实现数据结构（如链表、树）");
    println!("4. 实现多线程程序处理数据");
    println!("5. 创建一个简单的 Web 服务器");
}