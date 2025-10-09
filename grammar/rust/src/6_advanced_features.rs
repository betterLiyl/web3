use std::collections::HashMap;
use std::fmt::{Display, Debug};
use std::sync::{Arc, Mutex};
use std::thread;
use std::time::Duration;

// ============= 1. 所有权 (Ownership) =============

fn ownership_basics() {
    println!("=== 所有权基础 ===");
    
    // 移动语义
    let s1 = String::from("hello");
    let s2 = s1; // s1 的所有权移动到 s2
    // println!("{}", s1); // 这会编译错误，因为 s1 已经失效
    println!("s2: {}", s2);
    
    // 克隆
    let s3 = String::from("world");
    let s4 = s3.clone(); // 深拷贝
    println!("s3: {}, s4: {}", s3, s4);
    
    // Copy trait 的类型
    let x = 5;
    let y = x; // 整数实现了 Copy，所以这是拷贝而不是移动
    println!("x: {}, y: {}", x, y);
}

fn ownership_functions() {
    println!("\n=== 函数中的所有权 ===");
    
    fn takes_ownership(some_string: String) {
        println!("函数接收: {}", some_string);
    } // some_string 在这里被丢弃
    
    fn makes_copy(some_integer: i32) {
        println!("函数接收: {}", some_integer);
    }
    
    fn gives_ownership() -> String {
        String::from("yours")
    }
    
    fn takes_and_gives_back(a_string: String) -> String {
        a_string
    }
    
    let s = String::from("hello");
    takes_ownership(s); // s 的所有权移动到函数中
    // println!("{}", s); // 这会编译错误
    
    let x = 5;
    makes_copy(x); // x 被拷贝到函数中
    println!("x 仍然有效: {}", x);
    
    let s1 = gives_ownership();
    let s2 = String::from("hello");
    let s3 = takes_and_gives_back(s2);
    println!("s1: {}, s3: {}", s1, s3);
}

// ============= 2. 借用 (Borrowing) 和引用 (References) =============

fn borrowing_basics() {
    println!("\n=== 借用基础 ===");
    
    fn calculate_length(s: &String) -> usize {
        s.len()
    } // s 是引用，不会获得所有权
    
    fn change(some_string: &mut String) {
        some_string.push_str(", world");
    }
    
    let s1 = String::from("hello");
    let len = calculate_length(&s1); // 借用 s1
    println!("'{}' 的长度是 {}", s1, len);
    
    let mut s2 = String::from("hello");
    change(&mut s2); // 可变借用
    println!("修改后: {}", s2);
}

fn borrowing_rules() {
    println!("\n=== 借用规则 ===");
    
    let mut s = String::from("hello");
    
    // 规则1: 可以有多个不可变引用
    let r1 = &s;
    let r2 = &s;
    println!("r1: {}, r2: {}", r1, r2);
    
    // 规则2: 只能有一个可变引用
    let r3 = &mut s;
    println!("r3: {}", r3);
    // let r4 = &mut s; // 这会编译错误
    
    // 规则3: 不能同时有可变和不可变引用
    let r5 = &s;
    // let r6 = &mut s; // 这会编译错误
    println!("r5: {}", r5);
}

// ============= 3. 生命周期 (Lifetimes) =============

fn lifetime_basics() {
    println!("\n=== 生命周期基础 ===");
    
    // 生命周期注解
    fn longest<'a>(x: &'a str, y: &'a str) -> &'a str {
        if x.len() > y.len() {
            x
        } else {
            y
        }
    }
    
    let string1 = String::from("abcd");
    let string2 = "xyz";
    
    let result = longest(string1.as_str(), string2);
    println!("最长的字符串是: {}", result);
}

// 结构体中的生命周期
struct ImportantExcerpt<'a> {
    part: &'a str,
}

impl<'a> ImportantExcerpt<'a> {
    fn level(&self) -> i32 {
        3
    }
    
    fn announce_and_return_part(&self, announcement: &str) -> &str {
        println!("注意: {}", announcement);
        self.part
    }
}

fn lifetime_in_structs() {
    println!("\n=== 结构体中的生命周期 ===");
    
    let novel = String::from("Call me Ishmael. Some years ago...");
    let first_sentence = novel.split('.').next().expect("Could not find a '.'");
    let i = ImportantExcerpt {
        part: first_sentence,
    };
    
    println!("摘录: {}", i.part);
    println!("级别: {}", i.level());
    let result = i.announce_and_return_part("这是重要信息");
    println!("返回: {}", result);
}

// 静态生命周期
static HELLO_WORLD: &str = "Hello, world!";

fn static_lifetime() {
    println!("\n=== 静态生命周期 ===");
    println!("静态字符串: {}", HELLO_WORLD);
}

// ============= 4. Traits 特征 =============

// 定义 trait
trait Summary {
    fn summarize(&self) -> String;
    
    // 默认实现
    fn summarize_author(&self) -> String {
        String::from("(Read more...)")
    }
}

struct NewsArticle {
    headline: String,
    location: String,
    author: String,
    content: String,
}

impl Summary for NewsArticle {
    fn summarize(&self) -> String {
        format!("{}, by {} ({})", self.headline, self.author, self.location)
    }
    
    fn summarize_author(&self) -> String {
        format!("@{}", self.author)
    }
}

struct Tweet {
    username: String,
    content: String,
    reply: bool,
    retweet: bool,
}

impl Summary for Tweet {
    fn summarize(&self) -> String {
        format!("{}: {}", self.username, self.content)
    }
}

// trait 作为参数
fn notify(item: &impl Summary) {
    println!("突发新闻! {}", item.summarize());
}

// trait bound 语法
fn notify_bound<T: Summary>(item: &T) {
    println!("突发新闻! {}", item.summarize());
}

// 多个 trait bound
fn notify_multiple<T: Summary + Display>(item: &T) {
    println!("突发新闻! {}", item.summarize());
}

// where 子句
fn some_function<T, U>(_t: &T, _u: &U) -> i32
where
    T: Display + Clone,
    U: Clone + Debug,
{
    42
}

// 返回实现了 trait 的类型
fn returns_summarizable() -> impl Summary {
    Tweet {
        username: String::from("horse_ebooks"),
        content: String::from("当然，你可能已经知道了"),
        reply: false,
        retweet: false,
    }
}

fn trait_examples() {
    println!("\n=== Traits 示例 ===");
    
    let article = NewsArticle {
        headline: String::from("企鹅队赢得斯坦利杯冠军！"),
        location: String::from("匹兹堡，宾夕法尼亚州，美国"),
        author: String::from("Iceburgh"),
        content: String::from("企鹅队再次成为NHL最好的曲棍球队。"),
    };
    
    let tweet = Tweet {
        username: String::from("horse_ebooks"),
        content: String::from("当然，你可能已经知道了"),
        reply: false,
        retweet: false,
    };
    
    println!("新文章: {}", article.summarize());
    println!("新推文: {}", tweet.summarize());
    
    notify(&article);
    notify_bound(&tweet);
    
    let summarizable = returns_summarizable();
    println!("返回的摘要: {}", summarizable.summarize());
}

// ============= 5. 泛型 (Generics) =============

// 泛型函数
fn largest<T: PartialOrd + Copy>(list: &[T]) -> T {
    let mut largest = list[0];
    
    for &item in list {
        if item > largest {
            largest = item;
        }
    }
    
    largest
}

// 泛型结构体
struct Point<T> {
    x: T,
    y: T,
}

impl<T> Point<T> {
    fn x(&self) -> &T {
        &self.x
    }
}

impl Point<f32> {
    fn distance_from_origin(&self) -> f32 {
        (self.x.powi(2) + self.y.powi(2)).sqrt()
    }
}

// 多个泛型参数
struct Point2<T, U> {
    x: T,
    y: U,
}

impl<T, U> Point2<T, U> {
    fn mixup<V, W>(self, other: Point2<V, W>) -> Point2<T, W> {
        Point2 {
            x: self.x,
            y: other.y,
        }
    }
}

// 泛型枚举
enum MyOption<T> {
    Some(T),
    None,
}

enum MyResult<T, E> {
    Ok(T),
    Err(E),
}

fn generics_examples() {
    println!("\n=== 泛型示例 ===");
    
    let number_list = vec![34, 50, 25, 100, 65];
    let result = largest(&number_list);
    println!("最大的数字是 {}", result);
    
    let char_list = vec!['y', 'm', 'a', 'q'];
    let result = largest(&char_list);
    println!("最大的字符是 {}", result);
    
    let integer_point = Point { x: 5, y: 10 };
    let float_point = Point { x: 1.0, y: 4.0 };
    
    println!("integer_point.x = {}", integer_point.x());
    println!("float_point 到原点的距离 = {}", float_point.distance_from_origin());
    
    let p1 = Point2 { x: 5, y: 10.4 };
    let p2 = Point2 { x: "Hello", y: 'c' };
    let p3 = p1.mixup(p2);
    println!("p3.x = {}, p3.y = {}", p3.x, p3.y);
}

// ============= 6. 错误处理 =============

use std::fs::File;
use std::io::{self, Read};

fn error_handling_basics() {
    println!("\n=== 错误处理基础 ===");
    
    // Result 类型
    let f = File::open("hello.txt");
    
    let _f = match f {
        Ok(file) => file,
        Err(error) => {
            println!("打开文件时出现问题: {:?}", error);
            return;
        }
    };
    
    // unwrap 和 expect
    // let f = File::open("hello.txt").unwrap(); // panic if error
    // let f = File::open("hello.txt").expect("无法打开 hello.txt"); // panic with message
}

// 传播错误
fn read_username_from_file() -> Result<String, io::Error> {
    let mut f = File::open("hello.txt")?;
    let mut s = String::new();
    f.read_to_string(&mut s)?;
    Ok(s)
}

// ? 操作符的简化版本
fn read_username_from_file_short() -> Result<String, io::Error> {
    let mut s = String::new();
    File::open("hello.txt")?.read_to_string(&mut s)?;
    Ok(s)
}

// 自定义错误类型
#[derive(Debug)]
enum MyError {
    Io(io::Error),
    Parse(std::num::ParseIntError),
}

impl From<io::Error> for MyError {
    fn from(error: io::Error) -> Self {
        MyError::Io(error)
    }
}

impl From<std::num::ParseIntError> for MyError {
    fn from(error: std::num::ParseIntError) -> Self {
        MyError::Parse(error)
    }
}

// ============= 7. 宏 (Macros) =============

// 声明式宏
macro_rules! vec {
    ( $( $x:expr ),* ) => {
        {
            let mut temp_vec = Vec::new();
            $(
                temp_vec.push($x);
            )*
            temp_vec
        }
    };
}

// 更复杂的宏
macro_rules! hashmap {
    ($( $key: expr => $val: expr ),*) => {{
        let mut map = HashMap::new();
        $( map.insert($key, $val); )*
        map
    }}
}

// 过程宏示例（需要单独的 crate）
// #[derive(Debug)]
// struct MyStruct;

fn macro_examples() {
    println!("\n=== 宏示例 ===");
    
    let v = vec![1, 2, 3];
    println!("使用自定义 vec! 宏: {:?}", v);
    
    let map = hashmap!{
        "one" => 1,
        "two" => 2,
        "three" => 3
    };
    println!("使用 hashmap! 宏: {:?}", map);
}

// ============= 8. 并发 (Concurrency) =============

fn concurrency_basics() {
    println!("\n=== 并发基础 ===");
    
    // 创建线程
    let handle = thread::spawn(|| {
        for i in 1..10 {
            println!("hi number {} from the spawned thread!", i);
            thread::sleep(Duration::from_millis(1));
        }
    });
    
    for i in 1..5 {
        println!("hi number {} from the main thread!", i);
        thread::sleep(Duration::from_millis(1));
    }
    
    handle.join().unwrap();
}

fn move_closures() {
    println!("\n=== move 闭包 ===");
    
    let v = vec![1, 2, 3];
    
    let handle = thread::spawn(move || {
        println!("Here's a vector: {:?}", v);
    });
    
    handle.join().unwrap();
}

fn message_passing() {
    println!("\n=== 消息传递 ===");
    
    use std::sync::mpsc;
    
    let (tx, rx) = mpsc::channel();
    
    thread::spawn(move || {
        let val = String::from("hi");
        tx.send(val).unwrap();
    });
    
    let received = rx.recv().unwrap();
    println!("Got: {}", received);
    
    // 多个发送者
    let (tx, rx) = mpsc::channel();
    
    let tx1 = tx.clone();
    thread::spawn(move || {
        let vals = vec![
            String::from("hi"),
            String::from("from"),
            String::from("the"),
            String::from("thread")
        ];
        
        for val in vals {
            tx1.send(val).unwrap();
            thread::sleep(Duration::from_secs(1));
        }
    });
    
    thread::spawn(move || {
        let vals = vec![
            String::from("more"),
            String::from("messages"),
            String::from("for"),
            String::from("you")
        ];
        
        for val in vals {
            tx.send(val).unwrap();
            thread::sleep(Duration::from_secs(1));
        }
    });
    
    for received in rx {
        println!("Got: {}", received);
    }
}

fn shared_state() {
    println!("\n=== 共享状态 ===");
    
    let counter = Arc::new(Mutex::new(0));
    let mut handles = vec![];
    
    for _ in 0..10 {
        let counter = Arc::clone(&counter);
        let handle = thread::spawn(move || {
            let mut num = counter.lock().unwrap();
            *num += 1;
        });
        handles.push(handle);
    }
    
    for handle in handles {
        handle.join().unwrap();
    }
    
    println!("Result: {}", *counter.lock().unwrap());
}

// ============= 9. 智能指针 =============

use std::rc::Rc;
use std::cell::RefCell;

fn smart_pointers() {
    println!("\n=== 智能指针 ===");
    
    // Box<T>
    let b = Box::new(5);
    println!("b = {}", b);
    
    // Rc<T> - 引用计数
    let a = Rc::new(5);
    let b = Rc::clone(&a);
    let c = Rc::clone(&a);
    println!("引用计数: {}", Rc::strong_count(&a));
    
    // RefCell<T> - 内部可变性
    let value = RefCell::new(5);
    
    let borrowed_value = value.borrow();
    println!("borrowed value: {}", *borrowed_value);
    drop(borrowed_value); // 必须先释放借用
    
    let mut mutable_borrow = value.borrow_mut();
    *mutable_borrow += 1;
    drop(mutable_borrow);
    
    println!("modified value: {}", *value.borrow());
}

// ============= 10. 模式匹配 =============

fn pattern_matching() {
    println!("\n=== 模式匹配 ===");
    
    // 基本匹配
    let x = 1;
    match x {
        1 => println!("one"),
        2 => println!("two"),
        3 => println!("three"),
        _ => println!("anything"),
    }
    
    // 匹配命名变量
    let x = Some(5);
    let y = 10;
    
    match x {
        Some(50) => println!("Got 50"),
        Some(y) => println!("Matched, y = {:?}", y), // 这里的 y 是新变量
        _ => println!("Default case, x = {:?}", x),
    }
    
    println!("at the end: x = {:?}, y = {:?}", x, y);
    
    // 多个模式
    let x = 1;
    match x {
        1 | 2 => println!("one or two"),
        3 => println!("three"),
        _ => println!("anything"),
    }
    
    // 范围匹配
    let x = 5;
    match x {
        1..=5 => println!("one through five"),
        _ => println!("something else"),
    }
    
    // 解构结构体
    struct Point {
        x: i32,
        y: i32,
    }
    
    let p = Point { x: 0, y: 7 };
    
    match p {
        Point { x, y: 0 } => println!("On the x axis at {}", x),
        Point { x: 0, y } => println!("On the y axis at {}", y),
        Point { x, y } => println!("On neither axis: ({}, {})", x, y),
    }
    
    // 守卫
    let num = Some(4);
    
    match num {
        Some(x) if x < 5 => println!("less than five: {}", x),
        Some(x) => println!("{}", x),
        None => (),
    }
}

// ============= 主函数 =============

fn main() {
    println!("Rust 语言高级特性学习");
    println!("====================");
    
    // 1. 所有权
    ownership_basics();
    ownership_functions();
    
    // 2. 借用和引用
    borrowing_basics();
    borrowing_rules();
    
    // 3. 生命周期
    lifetime_basics();
    lifetime_in_structs();
    static_lifetime();
    
    // 4. Traits
    trait_examples();
    
    // 5. 泛型
    generics_examples();
    
    // 6. 错误处理
    error_handling_basics();
    
    // 7. 宏
    macro_examples();
    
    // 8. 并发
    concurrency_basics();
    move_closures();
    message_passing();
    shared_state();
    
    // 9. 智能指针
    smart_pointers();
    
    // 10. 模式匹配
    pattern_matching();
    
    println!("\n学习完成！");
}