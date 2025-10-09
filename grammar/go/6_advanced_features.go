package main

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"
)

// ============= 1. Goroutines 和 Channels =============

// 基本的 goroutine 使用
func basicGoroutine() {
	fmt.Println("=== Goroutines 基础 ===")
	
	// 启动一个 goroutine
	go func() {
		fmt.Println("这是在 goroutine 中执行的")
	}()
	
	// 等待一下，让 goroutine 有时间执行
	time.Sleep(100 * time.Millisecond)
}

// Channel 基础使用
func basicChannels() {
	fmt.Println("\n=== Channels 基础 ===")
	
	// 创建一个无缓冲 channel
	ch := make(chan string)
	
	// 在 goroutine 中发送数据
	go func() {
		ch <- "Hello from goroutine!"
	}()
	
	// 接收数据
	message := <-ch
	fmt.Println("接收到:", message)
	
	// 带缓冲的 channel
	bufferedCh := make(chan int, 3)
	bufferedCh <- 1
	bufferedCh <- 2
	bufferedCh <- 3
	
	fmt.Println("缓冲 channel 中的值:", <-bufferedCh, <-bufferedCh, <-bufferedCh)
}

// Channel 方向和关闭
func channelDirections() {
	fmt.Println("\n=== Channel 方向和关闭 ===")
	
	// 只发送的 channel
	sendOnly := func(ch chan<- int) {
		for i := 0; i < 5; i++ {
			ch <- i
		}
		close(ch) // 关闭 channel
	}
	
	// 只接收的 channel
	receiveOnly := func(ch <-chan int) {
		for value := range ch { // 使用 range 遍历 channel
			fmt.Printf("接收到: %d\n", value)
		}
	}
	
	ch := make(chan int)
	go sendOnly(ch)
	receiveOnly(ch)
}

// Select 语句
func selectStatement() {
	fmt.Println("\n=== Select 语句 ===")
	
	ch1 := make(chan string)
	ch2 := make(chan string)
	
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "来自 channel 1"
	}()
	
	go func() {
		time.Sleep(200 * time.Millisecond)
		ch2 <- "来自 channel 2"
	}()
	
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("收到:", msg1)
		case msg2 := <-ch2:
			fmt.Println("收到:", msg2)
		case <-time.After(300 * time.Millisecond):
			fmt.Println("超时")
		}
	}
}

// ============= 2. Interfaces 接口 =============

// 定义接口
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Writer interface {
	Write([]byte) (int, error)
}

type ReadWriter interface {
	Reader
	Writer
}

type Reader interface {
	Read([]byte) (int, error)
}

// 实现接口的结构体
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return 3.14159 * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * 3.14159 * c.Radius
}

// 接口的使用
func interfaceUsage() {
	fmt.Println("\n=== 接口使用 ===")
	
	var shapes []Shape
	shapes = append(shapes, Rectangle{Width: 10, Height: 5})
	shapes = append(shapes, Circle{Radius: 3})
	
	for _, shape := range shapes {
		fmt.Printf("面积: %.2f, 周长: %.2f\n", shape.Area(), shape.Perimeter())
	}
}

// 空接口和类型断言
func emptyInterfaceAndTypeAssertion() {
	fmt.Println("\n=== 空接口和类型断言 ===")
	
	var i interface{} = "hello"
	
	// 类型断言
	s, ok := i.(string)
	if ok {
		fmt.Printf("字符串值: %s\n", s)
	}
	
	// 类型 switch
	describe := func(i interface{}) {
		switch v := i.(type) {
		case int:
			fmt.Printf("整数: %d\n", v)
		case string:
			fmt.Printf("字符串: %s\n", v)
		case bool:
			fmt.Printf("布尔值: %t\n", v)
		default:
			fmt.Printf("未知类型: %T\n", v)
		}
	}
	
	describe(42)
	describe("hello")
	describe(true)
	describe([]int{1, 2, 3})
}

// ============= 3. Reflection 反射 =============

type Person struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"min=0,max=150"`
}

func (p Person) SayHello() {
	fmt.Printf("Hello, I'm %s\n", p.Name)
}

func reflectionBasics() {
	fmt.Println("\n=== 反射基础 ===")
	
	p := Person{Name: "Alice", Age: 30}
	
	// 获取类型和值
	t := reflect.TypeOf(p)
	v := reflect.ValueOf(p)
	
	fmt.Printf("类型: %s\n", t.Name())
	fmt.Printf("包路径: %s\n", t.PkgPath())
	
	// 遍历字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		
		fmt.Printf("字段 %s: %v (标签: %s)\n", 
			field.Name, 
			value.Interface(), 
			field.Tag)
	}
	
	// 遍历方法
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		fmt.Printf("方法: %s\n", method.Name)
	}
}

func reflectionModification() {
	fmt.Println("\n=== 反射修改值 ===")
	
	p := Person{Name: "Bob", Age: 25}
	
	// 获取指针的反射值
	v := reflect.ValueOf(&p).Elem()
	
	// 修改字段值
	nameField := v.FieldByName("Name")
	if nameField.CanSet() {
		nameField.SetString("Charlie")
	}
	
	ageField := v.FieldByName("Age")
	if ageField.CanSet() {
		ageField.SetInt(35)
	}
	
	fmt.Printf("修改后: %+v\n", p)
}

// ============= 4. Context 上下文 =============

func contextBasics() {
	fmt.Println("\n=== Context 基础 ===")
	
	// 创建带超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	// 模拟长时间运行的操作
	go func() {
		select {
		case <-time.After(3 * time.Second):
			fmt.Println("操作完成")
		case <-ctx.Done():
			fmt.Println("操作被取消:", ctx.Err())
		}
	}()
	
	time.Sleep(3 * time.Second)
}

func contextWithValue() {
	fmt.Println("\n=== Context 传递值 ===")
	
	type key string
	const userKey key = "user"
	
	// 创建带值的 context
	ctx := context.WithValue(context.Background(), userKey, "Alice")
	
	// 从 context 中获取值
	processRequest := func(ctx context.Context) {
		if user, ok := ctx.Value(userKey).(string); ok {
			fmt.Printf("处理用户 %s 的请求\n", user)
		}
	}
	
	processRequest(ctx)
}

// ============= 5. 并发模式 =============

// Worker Pool 模式
func workerPool() {
	fmt.Println("\n=== Worker Pool 模式 ===")
	
	jobs := make(chan int, 100)
	results := make(chan int, 100)
	
	// 启动 3 个 worker
	for w := 1; w <= 3; w++ {
		go func(id int) {
			for job := range jobs {
				fmt.Printf("Worker %d 处理任务 %d\n", id, job)
				time.Sleep(100 * time.Millisecond)
				results <- job * 2
			}
		}(w)
	}
	
	// 发送任务
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)
	
	// 收集结果
	for r := 1; r <= 5; r++ {
		<-results
	}
}

// Fan-out/Fan-in 模式
func fanOutFanIn() {
	fmt.Println("\n=== Fan-out/Fan-in 模式 ===")
	
	// 输入 channel
	input := make(chan int)
	
	// Fan-out: 将工作分发给多个 goroutine
	worker1 := make(chan int)
	worker2 := make(chan int)
	
	go func() {
		for i := range input {
			select {
			case worker1 <- i:
			case worker2 <- i:
			}
		}
		close(worker1)
		close(worker2)
	}()
	
	// 处理工作
	process := func(input <-chan int, output chan<- int) {
		for i := range input {
			output <- i * i
		}
		close(output)
	}
	
	output1 := make(chan int)
	output2 := make(chan int)
	go process(worker1, output1)
	go process(worker2, output2)
	
	// Fan-in: 合并结果
	result := make(chan int)
	var wg sync.WaitGroup
	wg.Add(2)
	
	merge := func(ch <-chan int) {
		defer wg.Done()
		for v := range ch {
			result <- v
		}
	}
	
	go merge(output1)
	go merge(output2)
	
	go func() {
		wg.Wait()
		close(result)
	}()
	
	// 发送数据
	go func() {
		for i := 1; i <= 5; i++ {
			input <- i
		}
		close(input)
	}()
	
	// 接收结果
	for r := range result {
		fmt.Printf("结果: %d\n", r)
	}
}

// ============= 6. 错误处理模式 =============

// 自定义错误类型
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("验证错误 - %s: %s", e.Field, e.Message)
}

// 错误包装
func validatePerson(p Person) error {
	if p.Name == "" {
		return ValidationError{Field: "Name", Message: "不能为空"}
	}
	if p.Age < 0 || p.Age > 150 {
		return ValidationError{Field: "Age", Message: "必须在 0-150 之间"}
	}
	return nil
}

func errorHandlingPatterns() {
	fmt.Println("\n=== 错误处理模式 ===")
	
	p1 := Person{Name: "", Age: 30}
	if err := validatePerson(p1); err != nil {
		if ve, ok := err.(ValidationError); ok {
			fmt.Printf("验证失败: %s\n", ve.Error())
		}
	}
	
	p2 := Person{Name: "Alice", Age: 30}
	if err := validatePerson(p2); err == nil {
		fmt.Println("验证通过")
	}
}

// ============= 主函数 =============

func main() {
	fmt.Println("Go 语言高级特性学习")
	fmt.Println("==================")
	
	// 1. Goroutines 和 Channels
	basicGoroutine()
	basicChannels()
	channelDirections()
	selectStatement()
	
	// 2. Interfaces
	interfaceUsage()
	emptyInterfaceAndTypeAssertion()
	
	// 3. Reflection
	reflectionBasics()
	reflectionModification()
	
	// 4. Context
	contextBasics()
	contextWithValue()
	
	// 5. 并发模式
	workerPool()
	fanOutFanIn()
	
	// 6. 错误处理
	errorHandlingPatterns()
	
	fmt.Println("\n学习完成！")
}