use std::collections::HashMap;

// 全局常量
const PI: f64 = 3.14159;
const MAX_SIZE: usize = 1000;


// 静态变量（全局可变变量，需要unsafe访问）
static mut GLOBAL_COUNTER: i32 = 0;
static GLOBAL_VAR: f64 = 3.14;
// 结构体定义
#[derive(Debug, Clone)]
struct Person {
    name: String,
    age: u32,
    city: String,
}

// 结构体实现
impl Person {
    // 关联函数（类似构造函数）
    fn new(name: String, age: u32, city: String) -> Self {
        Person { name, age, city }
    }
    
    // 方法
    fn introduce(&self) -> String {
        format!("我叫{}，今年{}岁，来自{}", self.name, self.age, self.city)
    }
    
    // 可变方法
    fn set_age(&mut self, age: u32) {
        self.age = age;
    }
    
    // 获取年龄
    fn get_age(&self) -> u32 {
        self.age
    }
}

// 枚举定义
#[derive(Debug)]
enum Color {
    Red,
    Green,
    Blue,
    RGB(u8, u8, u8),
    HSL { h: u16, s: u8, l: u8 },
}

// 普通函数
fn add(a: i32, b: i32) -> i32 {
    a + b
}

// 返回Result的函数
fn divide(a: f64, b: f64) -> Result<f64, String> {
    if b == 0.0 {
        Err("除数不能为零".to_string())
    } else {
        Ok(a / b)
    }
}

// 泛型函数
fn find_max<T: PartialOrd>(list: &[T]) -> Option<&T> {
    if list.is_empty() {
        None
    } else {
        let mut max = &list[0];
        for item in list.iter() {
            if item > max {
                max = item;
            }
        }
        Some(max)
    }
}

// 闭包示例函数
fn closure_demo() {
    println!("=== 闭包演示 ===");
    
    let numbers = vec![1, 2, 3, 4, 5];
    
    // 简单闭包
    let double = |x| x * 2;
    println!("双倍: {:?}", numbers.iter().map(|&x| double(x)).collect::<Vec<_>>());
    
    // 捕获环境变量的闭包
    let multiplier = 3;
    let triple: Vec<i32> = numbers.iter().map(|&x| x * multiplier).collect();
    println!("三倍: {:?}", triple);
    
    // 使用move关键字
    let text = String::from("Hello");
    let closure = move || println!("闭包中的文本: {}", text);
    closure();
}

// 演示变量声明
fn variable_demo() {
    println!("=== 变量声明演示 ===");
    
    // 不可变变量
    let immutable_var = 10;
    println!("不可变变量: {}", immutable_var);
    
    // 可变变量
    let mut mutable_var = 20;
    println!("可变变量初始值: {}", mutable_var);
    mutable_var = 30;
    println!("可变变量修改后: {}", mutable_var);
    
    // 变量遮蔽（shadowing）
    let x = 5;
    let x = x + 1;
    let x = x * 2;
    println!("变量遮蔽结果: {}", x);
    
    // 类型注解
    let typed_var: f64 = 3.14;
    println!("带类型注解的变量: {}", typed_var);
    
    // 元组解构
    let tuple = (1, 2.5, "hello");
    let (a, b, c) = tuple;
    println!("元组解构: a={}, b={}, c={}", a, b, c);
    
    // 数组
    let array = [1, 2, 3, 4, 5];
    println!("数组: {:?}", array);
    
    // 向量
    let mut vector = vec![1, 2, 3];
    vector.push(4);
    println!("向量: {:?}", vector);
    
    // 字符串
    let string_literal = "字符串字面量";
    let string_object = String::from("String对象");
    println!("字符串: '{}' 和 '{}'", string_literal, string_object);
    
    // 常量使用
    println!("常量: PI={}, MAX_SIZE={}", PI, MAX_SIZE);
}

// 演示函数
fn function_demo() {
    println!("\n=== 函数演示 ===");
    
    // 调用普通函数
    let result = add(5, 3);
    println!("5 + 3 = {}", result);
    
    // 调用返回Result的函数
    match divide(10.0, 3.0) {
        Ok(quotient) => println!("10.0 / 3.0 = {:.2}", quotient),
        Err(error) => println!("错误: {}", error),
    }
    
    match divide(10.0, 0.0) {
        Ok(quotient) => println!("10.0 / 0.0 = {:.2}", quotient),
        Err(error) => println!("错误: {}", error),
    }
    
    // 调用泛型函数
    let numbers = vec![3, 7, 1, 9, 2];
    if let Some(max) = find_max(&numbers) {
        println!("数组 {:?} 的最大值是: {}", numbers, max);
    }
    
    // 高阶函数示例
    let squared: Vec<i32> = numbers.iter().map(|x| x * x).collect();
    println!("平方: {:?}", squared);
    
    let even_numbers: Vec<&i32> = numbers.iter().filter(|&&x| x % 2 == 0).collect();
    println!("偶数: {:?}", even_numbers);
}

// 演示结构体
fn struct_demo() {
    println!("\n=== 结构体演示 ===");
    
    // 创建结构体实例
    let mut person1 = Person::new(
        "李四".to_string(),
        28,
        "北京".to_string(),
    );
    
    let person2 = Person {
        name: "王五".to_string(),
        age: 32,
        city: "上海".to_string(),
    };
    
    println!("{}", person1.introduce());
    println!("{}", person2.introduce());
    
    // 调用可变方法
    person1.set_age(29);
    println!("{} 更新年龄后: {}岁", person1.name, person1.get_age());
    
    // 结构体更新语法
    let person3 = Person {
        name: "赵六".to_string(),
        ..person2.clone()
    };
    println!("使用更新语法创建: {}", person3.introduce());
    
    // 元组结构体
    struct Point(i32, i32, i32);
    let origin = Point(0, 0, 0);
    println!("3D点坐标: ({}, {}, {})", origin.0, origin.1, origin.2);
}

// 演示循环
fn loop_demo() {
    println!("\n=== 循环演示 ===");
    
    // loop 无限循环
    print!("loop循环: ");
    let mut count = 0;
    loop {
        if count >= 3 {
            break;
        }
        print!("{} ", count + 1);
        count += 1;
    }
    println!();
    
    // while 循环
    print!("while循环: ");
    let mut number = 1;
    while number <= 5 {
        print!("{} ", number);
        number += 1;
    }
    println!();
    
    // for 循环 - 范围
    print!("for范围循环: ");
    for i in 1..=5 {
        print!("{} ", i);
    }
    println!();
    
    // for 循环 - 数组/向量
    let numbers = vec![10, 20, 30, 40, 50];
    print!("for遍历向量: ");
    for (index, value) in numbers.iter().enumerate() {
        print!("[{}]={} ", index, value);
    }
    println!();
    
    // for 循环 - 字符串
    let text = "Hello";
    print!("for遍历字符串: ");
    for (i, ch) in text.chars().enumerate() {
        print!("{}:{} ", i, ch);
    }
    println!();
    
    // for 循环 - HashMap
    let mut colors = HashMap::new();
    colors.insert("red", "红色");
    colors.insert("green", "绿色");
    colors.insert("blue", "蓝色");
    
    println!("for遍历HashMap:");
    for (key, value) in &colors {
        println!("  {}: {}", key, value);
    }
    
    // 带标签的循环
    'outer: for i in 1..=3 {
        for j in 1..=3 {
            if i == 2 && j == 2 {
                println!("在 i={}, j={} 时跳出外层循环", i, j);
                break 'outer;
            }
            print!("({},{}) ", i, j);
        }
    }
    println!();
}

// 演示条件判断
fn condition_demo() {
    println!("\n=== 条件判断演示 ===");
    
    let age = 18;
    let score = 85;
    
    // if-else 基本用法
    if age >= 18 {
        println!("年龄{}岁，已成年", age);
    } else {
        println!("年龄{}岁，未成年", age);
    }
    
    // if-else if-else
    if score >= 90 {
        println!("分数{}，等级：优秀", score);
    } else if score >= 80 {
        println!("分数{}，等级：良好", score);
    } else if score >= 60 {
        println!("分数{}，等级：及格", score);
    } else {
        println!("分数{}，等级：不及格", score);
    }
    
    // if 作为表达式
    let status = if age >= 18 { "成年人" } else { "未成年人" };
    println!("状态: {}", status);
    
    // 逻辑运算符
    let temperature = 25;
    let humidity = 60;
    if temperature > 20 && humidity < 70 {
        println!("温度{}°C，湿度{}%，天气舒适", temperature, humidity);
    }
    
    // Option 和 Result 的条件处理
    let maybe_number = Some(42);
    if let Some(number) = maybe_number {
        println!("Option中的数字: {}", number);
    }
    
    let result: Result<i32, &str> = Ok(100);
    if let Ok(value) = result {
        println!("Result中的值: {}", value);
    }
}

// 演示match表达式
fn match_demo() {
    println!("\n=== Match表达式演示 ===");
    
    // 基本match
    let day = 3;
    let day_name = match day {
        1 => "星期一",
        2 => "星期二",
        3 => "星期三",
        4 => "星期四",
        5 => "星期五",
        6 | 7 => "周末",
        _ => "无效的日期",
    };
    println!("第{}天是: {}", day, day_name);
    
    // match 范围
    let hour = 14;
    let time_period = match hour {
        0..=5 => "深夜",
        6..=11 => "上午",
        12..=17 => "下午",
        18..=23 => "晚上",
        _ => "无效时间",
    };
    println!("{}点属于: {}", hour, time_period);
    
    // match 枚举
    let color = Color::RGB(255, 0, 0);
    match color {
        Color::Red => println!("纯红色"),
        Color::Green => println!("纯绿色"),
        Color::Blue => println!("纯蓝色"),
        Color::RGB(r, g, b) => println!("RGB颜色: ({}, {}, {})", r, g, b),
        Color::HSL { h, s, l } => println!("HSL颜色: h={}, s={}, l={}", h, s, l),
    }
    
    // match Option
    let maybe_value = Some(42);
    match maybe_value {
        Some(value) => println!("Option包含值: {}", value),
        None => println!("Option为空"),
    }
    
    // match Result
    let result = divide(10.0, 2.0);
    match result {
        Ok(value) => println!("除法结果: {:.2}", value),
        Err(error) => println!("除法错误: {}", error),
    }
    
    // match 守卫
    let number = 4;
    match number {
        x if x < 0 => println!("{} 是负数", x),
        x if x == 0 => println!("{} 是零", x),
        x if x % 2 == 0 => println!("{} 是正偶数", x),
        x => println!("{} 是正奇数", x),
    }
    
    // match 解构
    let tuple = (1, 2, 3);
    match tuple {
        (0, y, z) => println!("第一个元素是0: y={}, z={}", y, z),
        (1, ..) => println!("第一个元素是1"),
        (.., 3) => println!("最后一个元素是3"),
        _ => println!("其他情况"),
    }
}

// 演示所有权和借用
fn ownership_demo() {
    println!("\n=== 所有权和借用演示 ===");
    
    // 所有权转移
    let s1 = String::from("hello");
    let s2 = s1; // s1的所有权转移给s2
    // println!("{}", s1); // 这行会编译错误
    println!("s2: {}", s2);
    
    // 克隆
    let s3 = s2.clone();
    println!("s2: {}, s3: {}", s2, s3);
    
    // 不可变借用
    let s4 = String::from("world");
    let len = calculate_length(&s4);
    println!("字符串 '{}' 的长度是 {}", s4, len);
    
    // 可变借用
    let mut s5 = String::from("hello");
    change_string(&mut s5);
    println!("修改后的字符串: {}", s5);
}

fn calculate_length(s: &String) -> usize {
    s.len()
}

fn change_string(s: &mut String) {
    s.push_str(", world!");
}

fn main() {
    println!("Rust语言基础语法学习");
    println!("===================");
    
    // 演示各个语法特性
    variable_demo();
    function_demo();
    struct_demo();
    loop_demo();
    condition_demo();
    match_demo();
    closure_demo();
    ownership_demo();
    
    // 安全地访问静态可变变量
    unsafe {
        GLOBAL_COUNTER += 1;
        println!("\n全局计数器: {}", GLOBAL_COUNTER);
        // GLOBAL_VAR = 2.71828;
        // println!("全局变量: {}", GLOBAL_VAR);
    }
    
    println!("\n学习完成！");
}