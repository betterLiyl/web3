// Rust 引用和借用学习文件
// 演示 Rust 独特的所有权、借用和生命周期系统

use std::rc::Rc;
use std::cell::RefCell;
use std::sync::{Arc, Mutex};
use std::collections::HashMap;

// 定义结构体用于演示
#[derive(Debug, Clone)]
struct Person {
    name: String,
    age: u32,
}

impl Person {
    fn new(name: String, age: u32) -> Self {
        Person { name, age }
    }
    
    // 不可变引用方法
    fn get_info(&self) -> String {
        format!("姓名: {}, 年龄: {}", self.name, self.age)
    }
    
    // 可变引用方法
    fn set_age(&mut self, age: u32) {
        self.age = age;
    }
    
    fn birthday(&mut self) {
        self.age += 1;
    }
    
    // 获取名字的引用
    fn get_name(&self) -> &str {
        &self.name
    }
}

// 基本借用演示
fn basic_borrowing_demo() {
    println!("\n=== 基本借用演示 ===");
    
    let x = 42;
    let y = &x;  // 不可变引用
    
    println!("x 的值: {}", x);
    println!("y 指向的值: {}", *y);  // 解引用
    println!("y 指向的值 (自动解引用): {}", y);
    
    // 多个不可变引用是允许的
    let z = &x;
    println!("多个不可变引用: x={}, y={}, z={}", x, y, z);
}

// 可变借用演示
fn mutable_borrowing_demo() {
    println!("\n=== 可变借用演示 ===");
    
    let mut x = 42;
    println!("原始值: {}", x);
    
    {
        let y = &mut x;  // 可变引用
        *y += 10;        // 通过可变引用修改值
        println!("通过可变引用修改后: {}", y);
    } // y 的作用域结束
    
    println!("修改后的 x: {}", x);
    
    // 借用规则演示
    let r1 = &x;      // 不可变引用
    let r2 = &x;      // 可以有多个不可变引用
    println!("多个不可变引用: r1={}, r2={}", r1, r2);
    
    // 注意：在 r1 和 r2 仍在使用时，不能创建可变引用
    // let r3 = &mut x;  // 这会编译错误
    
    // 当不可变引用不再使用后，可以创建可变引用
    let r3 = &mut x;
    *r3 += 5;
    println!("可变引用修改后: {}", r3);
}

// 函数参数借用
fn function_borrowing_demo() {
    println!("\n=== 函数参数借用演示 ===");
    
    // 不可变借用函数
    fn calculate_length(s: &String) -> usize {
        s.len()  // 可以读取，但不能修改
    }
    
    // 可变借用函数
    fn change_string(s: &mut String) {
        s.push_str(" - 已修改");
    }
    
    // 获取所有权的函数
    fn take_ownership(s: String) -> String {
        println!("获取所有权: {}", s);
        s  // 返回所有权
    }
    
    let s1 = String::from("Hello");
    let len = calculate_length(&s1);  // 借用，不转移所有权
    println!("字符串 '{}' 的长度是 {}", s1, len);
    
    let mut s2 = String::from("Hello");
    change_string(&mut s2);  // 可变借用
    println!("修改后的字符串: {}", s2);
    
    let s3 = String::from("World");
    let s3 = take_ownership(s3);  // 转移所有权并返回
    println!("返回的字符串: {}", s3);
}

// 结构体借用演示
fn struct_borrowing_demo() {
    println!("\n=== 结构体借用演示 ===");
    
    let mut person = Person::new(String::from("张三"), 25);
    println!("原始信息: {}", person.get_info());
    
    // 不可变借用
    let name_ref = person.get_name();
    println!("姓名引用: {}", name_ref);
    
    // 可变借用
    person.set_age(30);
    println!("修改年龄后: {}", person.get_info());
    
    // 部分借用
    let name = &person.name;
    let age = &person.age;
    println!("部分借用 - 姓名: {}, 年龄: {}", name, age);
}

// 生命周期演示
fn lifetime_demo() {
    println!("\n=== 生命周期演示 ===");
    
    // 生命周期注解函数
    fn longest<'a>(x: &'a str, y: &'a str) -> &'a str {
        if x.len() > y.len() {
            x
        } else {
            y
        }
    }
    
    let string1 = String::from("long string is long");
    {
        let string2 = String::from("xyz");
        let result = longest(string1.as_str(), string2.as_str());
        println!("最长的字符串是: {}", result);
    }
    
    // 结构体生命周期
    #[derive(Debug)]
    struct ImportantExcerpt<'a> {
        part: &'a str,
    }
    
    impl<'a> ImportantExcerpt<'a> {
        fn level(&self) -> i32 {
            3
        }
        
        fn announce_and_return_part(&self, announcement: &str) -> &str {
            println!("注意！{}", announcement);
            self.part
        }
    }
    
    let novel = String::from("Call me Ishmael. Some years ago...");
    let first_sentence = novel.split('.').next().expect("Could not find a '.'");
    let i = ImportantExcerpt {
        part: first_sentence,
    };
    println!("重要摘录: {:?}", i);
    println!("等级: {}", i.level());
    let part = i.announce_and_return_part("这是一个重要的摘录");
    println!("返回的部分: {}", part);
}

// 智能指针演示
fn smart_pointers_demo() {
    println!("\n=== 智能指针演示 ===");
    
    // Box<T> - 堆分配
    println!("--- Box<T> 演示 ---");
    let b = Box::new(5);
    println!("Box 中的值: {}", b);
    
    // 递归类型需要 Box
    #[derive(Debug)]
    enum BoxList {
        Cons(i32, Box<BoxList>),
        Nil,
    }
    
    use BoxList::{Cons as BoxCons, Nil as BoxNil};
    
    let list = BoxCons(1, Box::new(BoxCons(2, Box::new(BoxCons(3, Box::new(BoxNil))))));
    println!("Box 链表: {:?}", list);
    
    // Rc<T> - 引用计数智能指针
    println!("\n--- Rc<T> 演示 ---");
    
    // 为 Rc 定义单独的 List 类型
    #[derive(Debug)]
    enum List {
        Cons(i32, Rc<List>),
        Nil,
    }
    
    use List::{Cons, Nil};
    
    let a = Rc::new(Cons(5, Rc::new(Cons(10, Rc::new(Nil)))));
    println!("创建 a 后，a 的引用计数: {}", Rc::strong_count(&a));
    
    let _b = Cons(3, Rc::clone(&a));
    println!("创建 b 后，a 的引用计数: {}", Rc::strong_count(&a));
    
    {
        let _c = Cons(4, Rc::clone(&a));
        println!("创建 c 后，a 的引用计数: {}", Rc::strong_count(&a));
    }
    println!("c 离开作用域后，a 的引用计数: {}", Rc::strong_count(&a));
    
    // RefCell<T> - 内部可变性
    println!("\n--- RefCell<T> 演示 ---");
    let value = RefCell::new(5);
    
    {
        let mut borrowed = value.borrow_mut();
        *borrowed += 10;
    }
    
    println!("RefCell 中的值: {}", value.borrow());
    
    // Rc<RefCell<T>> 组合使用
    println!("\n--- Rc<RefCell<T>> 组合演示 ---");
    let shared_data = Rc::new(RefCell::new(vec![1, 2, 3]));
    let shared_data1 = Rc::clone(&shared_data);
    let shared_data2 = Rc::clone(&shared_data);
    
    shared_data1.borrow_mut().push(4);
    shared_data2.borrow_mut().push(5);
    
    println!("共享数据: {:?}", shared_data.borrow());
}

// 并发安全的智能指针
fn concurrent_smart_pointers_demo() {
    println!("\n=== 并发安全智能指针演示 ===");
    
    // Arc<T> - 原子引用计数
    let counter = Arc::new(Mutex::new(0));
    let mut handles = vec![];
    
    for _ in 0..10 {
        let counter = Arc::clone(&counter);
        let handle = std::thread::spawn(move || {
            let mut num = counter.lock().unwrap();
            *num += 1;
        });
        handles.push(handle);
    }
    
    for handle in handles {
        handle.join().unwrap();
    }
    
    println!("并发计数结果: {}", *counter.lock().unwrap());
}

// 借用检查器规避技巧
fn borrowing_patterns_demo() {
    println!("\n=== 借用模式演示 ===");
    
    // 使用作用域分离借用
    let mut data = vec![1, 2, 3, 4, 5];
    
    {
        let first = &data[0];
        println!("第一个元素: {}", first);
    } // first 的借用在这里结束
    
    data.push(6); // 现在可以修改 data
    println!("添加元素后: {:?}", data);
    
    // 使用索引而不是引用
    let index = 0;
    println!("通过索引访问: {}", data[index]);
    data.push(7);
    println!("再次修改后: {:?}", data);
    
    // 分割借用
    let (left, right) = data.split_at_mut(3);
    left[0] = 100;
    right[0] = 200;
    println!("分割借用修改后: {:?}", data);
}

// 高级引用模式
fn advanced_reference_patterns() {
    println!("\n=== 高级引用模式演示 ===");
    
    // 引用的引用
    let x = 5;
    let y = &x;
    let z = &y;
    println!("引用的引用: x={}, y={}, z={}", x, y, **z);
    
    // 模式匹配中的引用
    let point = (3, 5);
    match point {
        (x, y) => println!("坐标: ({}, {})", x, y),
    }
    
    match &point {
        &(x, y) => println!("引用模式匹配: ({}, {})", x, y),
    }
    
    // 解构引用
    let person = Person::new(String::from("李四"), 28);
    let Person { name, age } = &person;
    println!("解构引用: 姓名={}, 年龄={}", name, age);
    
    // 引用和 Option
    let maybe_person: Option<Person> = Some(person);
    match &maybe_person {
        Some(Person { name, age }) => {
            println!("Option 中的引用: 姓名={}, 年龄={}", name, age);
        }
        None => println!("没有人员信息"),
    }
}

// 内存安全演示
fn memory_safety_demo() {
    println!("\n=== 内存安全演示 ===");
    
    // 悬垂引用预防
    fn no_dangling_reference() -> String {
        let s = String::from("hello");
        s  // 返回所有权而不是引用
        // &s  // 这会导致编译错误：悬垂引用
    }
    
    let result = no_dangling_reference();
    println!("安全返回的字符串: {}", result);
    
    // 数据竞争预防
    let mut data = vec![1, 2, 3];
    
    // 这些操作是安全的，因为借用检查器确保没有数据竞争
    let len = data.len();
    println!("数据长度: {}", len);
    
    data.push(4);
    println!("添加元素后: {:?}", data);
    
    // 迭代器和借用
    println!("迭代器演示:");
    for (index, value) in data.iter().enumerate() {
        println!("  索引 {}: 值 {}", index, value);
    }
    
    // 可变迭代器
    for value in data.iter_mut() {
        *value *= 2;
    }
    println!("乘以2后: {:?}", data);
}

// 实际应用示例：简单的缓存系统
fn cache_system_demo() {
    println!("\n=== 缓存系统演示 ===");
    
    struct Cache {
        data: HashMap<String, String>,
    }
    
    impl Cache {
        fn new() -> Self {
            Cache {
                data: HashMap::new(),
            }
        }
        
        fn get(&self, key: &str) -> Option<&String> {
            self.data.get(key)
        }
        
        fn set(&mut self, key: String, value: String) {
            self.data.insert(key, value);
        }
        
        fn contains_key(&self, key: &str) -> bool {
            self.data.contains_key(key)
        }
    }
    
    let mut cache = Cache::new();
    
    // 设置缓存
    cache.set("user:1".to_string(), "张三".to_string());
    cache.set("user:2".to_string(), "李四".to_string());
    
    // 获取缓存
    if let Some(user) = cache.get("user:1") {
        println!("缓存命中: {}", user);
    }
    
    // 检查键是否存在
    if cache.contains_key("user:3") {
        println!("用户3存在");
    } else {
        println!("用户3不存在");
    }
}

fn main() {
    println!("=== Rust 引用和借用学习 ===");
    
    // 1. 基本借用
    basic_borrowing_demo();
    
    // 2. 可变借用
    mutable_borrowing_demo();
    
    // 3. 函数参数借用
    function_borrowing_demo();
    
    // 4. 结构体借用
    struct_borrowing_demo();
    
    // 5. 生命周期
    lifetime_demo();
    
    // 6. 智能指针
    smart_pointers_demo();
    
    // 7. 并发安全智能指针
    concurrent_smart_pointers_demo();
    
    // 8. 借用模式
    borrowing_patterns_demo();
    
    // 9. 高级引用模式
    advanced_reference_patterns();
    
    // 10. 内存安全
    memory_safety_demo();
    
    // 11. 实际应用示例
    cache_system_demo();
    
    println!("\n=== Rust 借用规则总结 ===");
    println!("1. 在任意给定时间，要么只能有一个可变引用，要么只能有多个不可变引用");
    println!("2. 引用必须总是有效的（不能有悬垂引用）");
    println!("3. 所有权系统确保内存安全，无需垃圾回收器");
    println!("4. 借用检查器在编译时验证这些规则");
    println!("5. 智能指针提供了额外的内存管理模式");
    println!("6. 生命周期确保引用的有效性");
    
    println!("\n学习完成！");
}