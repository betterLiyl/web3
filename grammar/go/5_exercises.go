package main

import (
	"fmt"
	"math"
	"strings"
)

// ==================== 练习题 ====================
// 请完成以下练习题，每个函数都有详细的说明和示例

// 练习1: 变量和常量
// 完成以下函数，计算圆的面积和周长
// 提示：使用 math.Pi 常量
func calculateCircle(radius float64) (area float64, perimeter float64) {
	// TODO: 实现计算圆的面积和周长
	area = math.Pi * radius * radius
	// perimeter = 2 * π * r
	perimeter = 2 * math.Pi * radius
	return area, perimeter
}

// 练习2: 字符串操作
// 完成函数，将输入的字符串转换为标题格式（每个单词首字母大写）
// 例如："hello world" -> "Hello World"
func toTitleCase(input string) string {
	// TODO: 使用 strings 包的函数来实现
	// 提示：可以使用 strings.Fields() 和 strings.Title() 或 strings.ToTitle()
	words := strings.Fields(input)
	for i, word := range words {
		words[i] = strings.Title(word)
	}
	return strings.Join(words, " ")
}

// 练习3: 切片操作
// 完成函数，找出切片中的最大值和最小值
func findMinMax(numbers []int) (min int, max int) {
	if len(numbers) == 0 {
		return 0, 0
	}
	max = numbers[0]
	min = numbers[0]
	for _, v := range numbers {
		if v >= max {
			max = v
		}
		if v <= min {
			min = v
		}
	}

	return min, max // 请替换为正确的实现
}

// 练习4: 结构体和方法
// 定义一个学生结构体，包含姓名、年龄和成绩切片
// 实现计算平均成绩的方法和添加成绩的方法
type ExerciseStudent struct {
	Name   string
	Age    int
	Grades []float64
}

// 为 ExerciseStudent 实现一个方法，计算平均成绩
func (s ExerciseStudent) AverageGrade() float64 {
	// TODO: 计算并返回平均成绩
	// 如果没有成绩，返回 0
	if len(s.Grades) == 0 {
		return 0.0
	}
	var sum float64 = 0.0
	for _, grade := range s.Grades {
		sum += grade
	}
	return sum / float64(len(s.Grades))
}

// 为 ExerciseStudent 实现一个指针方法，添加新成绩
func (s *ExerciseStudent) AddGrade(grade float64) {
	s.Grades = append(s.Grades, grade)
}

// 练习5: 指针操作
// 完成函数，交换两个整数的值
func swapIntegers(a, b *int) {
	// TODO: 使用指针交换两个整数的值
	*a,*b = *b,*a

}

// 练习6: 错误处理
// 完成函数，安全地进行除法运算
func safeDivide(a, b float64) (result float64, err error) {
	// TODO: 实现安全除法
	// 如果 b 为 0，返回错误
	if b == 0 {
		return 0, fmt.Errorf("除数不能为0")
	}
	// 否则返回 a/b 的结果
	return a/b, nil
}

// 练习7: 接口
// 定义一个 Shape 接口，包含 Area() 和 Perimeter() 方法
// 实现一个 Rectangle 结构体来满足这个接口
type ExerciseShape interface {
	Area() float64
	Perimeter() float64
}

// 矩形结构体
type ExerciseRectangle struct {
	Width  float64
	Height float64
}

// TODO: 为 Rectangle 实现 Shape 接口的方法
func (r ExerciseRectangle) Area() float64 {
	return r.Width * r.Height
}

func (r ExerciseRectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// 练习8: 高阶函数
// 完成函数，对切片中的每个元素应用给定的函数
func mapSlice(numbers []int, fn func(int) int) []int {
	// TODO: 实现 map 函数
	// 对 numbers 中的每个元素应用 fn 函数
	newNumbers := make([]int, len(numbers))
	for i,v := range numbers{
		newNumbers[i] = fn(v)
	}
	// 返回新的切片
	return newNumbers
}

// 练习9: 并发（可选，较难）
// 使用 goroutine 和 channel 计算数字的平方
func calculateSquares(numbers []int) []int {
	// TODO: 使用 goroutine 并发计算每个数字的平方
	// 提示：创建一个 channel 来收集结果
	return nil
}

// ==================== 测试函数 ====================
// 以下是一些测试函数，用于验证你的实现

func testCalculateCircle() {
	fmt.Println("=== 测试圆形计算 ===")
	area, perimeter := calculateCircle(5.0)
	fmt.Printf("半径为5的圆：面积=%.2f, 周长=%.2f\n", area, perimeter)
	// 期望结果：面积≈78.54, 周长≈31.42
}

func testTitleCase() {
	fmt.Println("\n=== 测试标题格式转换 ===")
	result := toTitleCase("hello world go programming")
	fmt.Printf("转换结果: %s\n", result)
	// 期望结果："Hello World Go Programming"
}

func testMinMax() {
	fmt.Println("\n=== 测试最大最小值 ===")
	numbers := []int{3, 7, 2, 9, 1, 5, 8}
	min, max := findMinMax(numbers)
	fmt.Printf("数组 %v 的最小值: %d, 最大值: %d\n", numbers, min, max)
	// 期望结果：最小值=1, 最大值=9
}

func testStudent() {
	fmt.Println("\n=== 测试学生结构体 ===")
	// TODO: 创建学生实例并测试方法
	student := ExerciseStudent{Name: "张三", Age: 20, Grades: []float64{85.5, 92.0, 78.5}}
	fmt.Printf("学生信息: %+v\n", student)
	fmt.Printf("平均成绩: %.2f\n", student.AverageGrade())
	student.AddGrade(88.0)
	fmt.Printf("添加成绩后平均分: %.2f\n", student.AverageGrade())
}

func testSwap() {
	fmt.Println("\n=== 测试交换函数 ===")
	a, b := 10, 20
	fmt.Printf("交换前: a=%d, b=%d\n", a, b)
	swapIntegers(&a, &b)
	fmt.Printf("交换后: a=%d, b=%d\n", a, b)
	// 期望结果：a=20, b=10
}

func testSafeDivide() {
	fmt.Println("\n=== 测试安全除法 ===")
	result, err := safeDivide(10, 2)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("10 ÷ 2 = %.2f\n", result)
	}

	result, err = safeDivide(10, 0)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("10 ÷ 0 = %.2f\n", result)
	}
}

func testShape() {
	fmt.Println("\n=== 测试形状接口 ===")
	// TODO: 创建矩形实例并测试接口方法
	rect := ExerciseRectangle{Width: 5, Height: 3}
	fmt.Printf("矩形 %+v\n", rect)
	fmt.Printf("面积: %.2f\n", rect.Area())
	fmt.Printf("周长: %.2f\n", rect.Perimeter())
	
	// 接口使用
	var shape ExerciseShape = rect
	fmt.Printf("通过接口访问面积: %.2f\n", shape.Area())
}

func testMapSlice() {
	fmt.Println("\n=== 测试映射函数 ===")
	numbers := []int{1, 2, 3, 4, 5}

	// 测试平方函数
	fmt.Printf("原数组: %v\n", numbers)
	squares := mapSlice(numbers, func(x int) int { return x * x })
	
	fmt.Printf("平方后: %v\n", squares)

	// 测试加倍函数
	doubled := mapSlice(numbers, func(x int) int { return x * 2 })
	fmt.Printf("加倍后: %v\n", doubled)
}

func main5() {
	fmt.Println("Go 语言练习题")
	fmt.Println("请完成上面标有 TODO 的函数实现")
	fmt.Println("然后运行测试函数查看结果")
	fmt.Println("=====================================")

	// 运行测试函数
	testCalculateCircle()
	testTitleCase()
	testMinMax()
	testStudent()
	testSwap()
	testSafeDivide()
	testShape()
	testMapSlice()

	fmt.Println("\n练习完成后，你可以尝试以下进阶挑战：")
	fmt.Println("1. 实现一个简单的计算器")
	fmt.Println("2. 创建一个学生管理系统")
	fmt.Println("3. 使用 goroutine 实现并发计算")
	fmt.Println("4. 实现自定义排序算法")
}
