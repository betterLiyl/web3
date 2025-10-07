package main

import (
	"fmt"
	"strings"
)

// 全局变量声明
var globalVar int = 100
var globalString string = "全局字符串"
var number float64 = 3.14

// 常量声明
const PI = 3.14159
const MaxSize = 1000
const DefaultAge = 18


// 结构体定义
type BasicPerson struct {
	Name string
	Age  int
	City string
}

// 结构体方法
func (p BasicPerson) Introduce() string {
	return fmt.Sprintf("我叫%s，今年%d岁，来自%s", p.Name, p.Age, p.City)
}

// 结构体指针方法
func (p *BasicPerson) SetAge(age int) {
	p.Age = age
}

// 普通函数
func add(a, b int) int {
	return a + b
}
func substract(a,b float64) float64{
	return a - b;
}

// 多返回值函数
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("除数不能为零")
	}
	return a / b, nil
}

// 可变参数函数
func sum(numbers ...int) int {
	total := 0
	for _, num := range numbers {
		total += num
	}
	return total
}

// 演示变量声明的函数
func variableDemo() {
	fmt.Println("=== 变量声明演示 ===")
	
	// 局部变量声明的几种方式
	var localVar1 int = 10
	var localVar2 = 20 // 类型推断
	localVar3 := 30    // 短变量声明
	
	fmt.Printf("局部变量: %d, %d, %d\n", localVar1, localVar2, localVar3)
	fmt.Printf("全局变量: %d, %s\n", globalVar, globalString)
	fmt.Printf("常量: PI=%.5f, MaxSize=%d\n", PI, MaxSize)
	
	// 多变量声明
	var x, y, z int = 1, 2, 3
	fmt.Printf("多变量声明: x=%d, y=%d, z=%d\n", x, y, z)
	
	// 不同类型变量
	var (
		name    string = "张三"
		age     int    = 25
		height  float64 = 175.5
		married bool   = false
	)
	fmt.Printf("个人信息: %s, %d岁, %.1fcm, 已婚:%t\n", name, age, height, married)
}

// 演示函数的函数
func functionDemo() {
	fmt.Println("\n=== 函数演示 ===")
	
	// 调用普通函数
	result := add(5, 3)
	fmt.Printf("5 + 3 = %d\n", result)
	res2 := substract(10.5, 3.2)
	fmt.Printf("10.5 - 3.2 = %.2f\n", res2)
	// 调用多返回值函数
	quotient, err := divide(10, 3)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("10 / 3 = %.2f\n", quotient)
	}
	
	// 调用可变参数函数
	total := sum(1, 2, 3, 4, 5)
	fmt.Printf("1+2+3+4+5 = %d\n", total)
	
	// 匿名函数
	multiply := func(a, b int) int {
		return a * b
	}
	fmt.Printf("匿名函数: 4 * 6 = %d\n", multiply(4, 6))
}

// 演示结构体的函数
func structDemo() {
	fmt.Println("\n=== 结构体演示 ===")
	
	// 创建结构体实例
	person1 := BasicPerson{
		Name: "李四",
		Age:  28,
		City: "北京",
	}
	
	// 另一种创建方式
	person2 := BasicPerson{"王五", 32, "上海"}
	
	// 使用new创建
	person3 := new(BasicPerson)
	person3.Name = "赵六"
	person3.Age = 24
	person3.City = "广州"
	
	fmt.Println(person1.Introduce())
	fmt.Println(person2.Introduce())
	fmt.Println(person3.Introduce())
	
	// 调用指针方法
	person1.SetAge(29)
	fmt.Printf("%s 更新年龄后: %d岁\n", person1.Name, person1.Age)
}

// 演示循环的函数
func loopDemo() {
	fmt.Println("\n=== 循环演示 ===")
	
	// for 循环 - 传统形式
	fmt.Print("传统for循环: ")
	for i := 1; i <= 5; i++ {
		fmt.Printf("%d ", i)
	}
	fmt.Println()
	
	// for 循环 - while形式
	fmt.Print("while形式循环: ")
	j := 1
	for j <= 3 {
		fmt.Printf("%d ", j*2)
		j++
	}
	fmt.Println()
	
	// for range 循环 - 数组/切片
	numbers := []int{10, 20, 30, 40, 50}
	fmt.Print("range遍历切片: ")
	for index, value := range numbers {
		fmt.Printf("[%d]=%d ", index, value)
	}
	fmt.Println()
	
	// for range 循环 - 字符串
	text := "Hello"
	fmt.Print("range遍历字符串: ")
	for i, char := range text {
		fmt.Printf("%d:%c ", i, char)
	}
	fmt.Println()
	
	// for range 循环 - map
	colors := map[string]string{
		"red":   "红色",
		"green": "绿色",
		"blue":  "蓝色",
	}
	fmt.Println("range遍历map:")
	for key, value := range colors {
		fmt.Printf("  %s: %s\n", key, value)
	}
	
	// 无限循环示例（带break）
	fmt.Print("无限循环示例: ")
	count := 0
	for {
		if count >= 3 {
			break
		}
		fmt.Printf("第%d次 ", count+1)
		count++
	}
	fmt.Println()
}

// 演示判断的函数
func conditionDemo() {
	fmt.Println("\n=== 条件判断演示 ===")
	
	age := 18
	score := 85
	
	// if-else 基本用法
	if age >= 18 {
		fmt.Printf("年龄%d岁，已成年\n", age)
	} else {
		fmt.Printf("年龄%d岁，未成年\n", age)
	}
	
	// if-else if-else
	if score >= 90 {
		fmt.Printf("分数%d，等级：优秀\n", score)
	} else if score >= 80 {
		fmt.Printf("分数%d，等级：良好\n", score)
	} else if score >= 60 {
		fmt.Printf("分数%d，等级：及格\n", score)
	} else {
		fmt.Printf("分数%d，等级：不及格\n", score)
	}
	
	// if 初始化语句
	if num := 42; num%2 == 0 {
		fmt.Printf("数字%d是偶数\n", num)
	} else {
		fmt.Printf("数字%d是奇数\n", num)
	}
	
	// 逻辑运算符
	temperature := 25
	humidity := 60
	if temperature > 20 && humidity < 70 {
		fmt.Printf("温度%d°C，湿度%d%%，天气舒适\n", temperature, humidity)
	}
	
	// 字符串判断
	name := "Go语言"
	if strings.Contains(name, "Go") || strings.Contains(name, "语言") {
		fmt.Printf("'%s' 包含关键词\n", name)
	}
}

// 演示switch的函数
func switchDemo() {
	fmt.Println("\n=== Switch演示 ===")
	
	// 基本switch
	day := 3
	switch day {
	case 1:
		fmt.Println("星期一")
	case 2:
		fmt.Println("星期二")
	case 3:
		fmt.Println("星期三")
	case 4:
		fmt.Println("星期四")
	case 5:
		fmt.Println("星期五")
	case 6, 7:
		fmt.Println("周末")
	default:
		fmt.Println("无效的日期")
	}
	
	// switch 初始化语句
	switch hour := 14; {
	case hour < 6:
		fmt.Printf("%d点，深夜时间\n", hour)
	case hour < 12:
		fmt.Printf("%d点，上午时间\n", hour)
	case hour < 18:
		fmt.Printf("%d点，下午时间\n", hour)
	default:
		fmt.Printf("%d点，晚上时间\n", hour)
	}
	
	// switch 表达式
	grade := 'B'
	switch grade {
	case 'A':
		fmt.Println("优秀")
	case 'B':
		fmt.Println("良好")
	case 'C':
		fmt.Println("及格")
	case 'D', 'F':
		fmt.Println("不及格")
	default:
		fmt.Println("无效等级")
	}
	
	// switch 类型判断
	var value interface{} = "Hello"
	switch v := value.(type) {
	case string:
		fmt.Printf("字符串类型: %s\n", v)
	case int:
		fmt.Printf("整数类型: %d\n", v)
	case bool:
		fmt.Printf("布尔类型: %t\n", v)
	default:
		fmt.Printf("未知类型: %T\n", v)
	}
	
	// fallthrough 示例
	number := 2
	fmt.Printf("数字%d的特性: ", number)
	switch number {
	case 2:
		fmt.Print("是偶数 ")
		fallthrough
	case 1, 3, 5, 7:
		fmt.Print("是质数 ")
	default:
		fmt.Print("是合数")
	}
	fmt.Println()
}

func main3() {
	fmt.Println("Go语言基础语法学习")
	fmt.Println("==================")
	fmt.Println("test-----------")
	// 演示各个语法特性
	variableDemo()
	functionDemo()
	structDemo()
	loopDemo()
	conditionDemo()
	switchDemo()
	
	fmt.Println("\n学习完成！")
}