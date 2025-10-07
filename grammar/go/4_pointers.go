package main

import (
	"fmt"
	"unsafe"
)

// 定义结构体用于演示
type PointerPerson struct {
	Name string
	Age  int
}

// 结构体方法 - 值接收者
func (p PointerPerson) GetInfo() string {
	return fmt.Sprintf("姓名: %s, 年龄: %d", p.Name, p.Age)
}

// 结构体方法 - 指针接收者
func (p *PointerPerson) SetAge(age int) {
	p.Age = age
}

func (p *PointerPerson) Birthday() {
	p.Age++
}

// 函数参数 - 值传递
func updatePersonByValue(p PointerPerson, newAge int) {
	p.Age = newAge
	fmt.Printf("函数内部 (值传递): %s\n", p.GetInfo())
}

// 函数参数 - 指针传递
func updatePersonByPointer(p *PointerPerson, newAge int) {
	p.Age = newAge
	fmt.Printf("函数内部 (指针传递): %s\n", p.GetInfo())
}

// 返回指针的函数
func createPerson(name string, age int) *PointerPerson {
	// 局部变量的地址可以安全返回，Go 会自动处理内存分配
	person := PointerPerson{Name: name, Age: age}
	return &person
}

// 指针作为函数参数进行交换
func swap(a, b *int) {
	*a, *b = *b, *a
}

// 指针数组和数组指针
func arrayPointerDemo() {
	fmt.Println("\n=== 数组和指针演示 ===")

	// 数组
	arr := [5]int{1, 2, 3, 4, 5}
	fmt.Printf("原始数组: %v\n", arr)

	// 数组指针 - 指向整个数组的指针
	var arrPtr *[5]int = &arr
	fmt.Printf("数组指针指向的数组: %v\n", *arrPtr)
	fmt.Printf("通过数组指针访问元素: %d\n", (*arrPtr)[2])

	// 修改数组元素
	(*arrPtr)[2] = 30
	fmt.Printf("修改后的数组: %v\n", arr)

	// 指针数组 - 存储指针的数组
	var ptrArr [3]*int
	x, y, z := 10, 20, 30
	ptrArr[0] = &x
	ptrArr[1] = &y
	ptrArr[2] = &z

	fmt.Printf("指针数组中的值: ")
	for i, ptr := range ptrArr {
		fmt.Printf("ptrArr[%d] = %d ", i, *ptr)
	}
	fmt.Println()
}

// 切片和指针
func slicePointerDemo() {
	fmt.Println("\n=== 切片和指针演示 ===")

	// 切片本身就包含指向底层数组的指针
	//切片的长度和容量 是不同的概念，切片前面的[]表示长度，后面的[]表示容量，[]里面是空着的，数组的[]表示元素数量
	// 长度是切片当前包含的元素数量
	// 容量是切片底层数组从开始到结束的元素数量
	slice := []int{1, 2, 3, 4, 5}
	fmt.Printf("原始切片: %v\n", slice)

	// 获取切片元素的指针
	ptr := &slice[2]
	fmt.Printf("slice[2] 的地址: %p, 值: %d\n", ptr, *ptr)

	// 通过指针修改切片元素
	*ptr = 30
	fmt.Printf("修改后的切片: %v\n", slice)

	// 切片指针
	slicePtr := &slice
	fmt.Printf("切片指针指向的切片: %v\n", *slicePtr)

	// 通过切片指针修改元素
	(*slicePtr)[0] = 100
	fmt.Printf("通过切片指针修改后: %v\n", slice)
}

// 指针运算和内存地址
func pointerArithmeticDemo() {
	fmt.Println("\n=== 指针和内存地址演示 ===")

	var x int = 42
	var ptr *int = &x

	fmt.Printf("变量 x 的值: %d\n", x)
	fmt.Printf("变量 x 的地址: %p\n", &x)
	fmt.Printf("指针 ptr 的值 (x的地址): %p\n", ptr)
	fmt.Printf("指针 ptr 指向的值: %d\n", *ptr)
	fmt.Printf("指针 ptr 自身的地址: %p\n", &ptr)
	fmt.Printf("指针 ptr引用 指向的: %p\n", *&ptr) // 指针 ptr 引用的是指针 ptr 指向的地址，而不是指针 ptr 自身的地址

	// 指针的大小
	fmt.Printf("指针的大小: %d 字节\n", unsafe.Sizeof(ptr))

	// 不同类型的指针
	var str string = "Hello"
	var strPtr *string = &str
	fmt.Printf("字符串指针指向的值: %s\n", *strPtr)

	// 空指针
	var nilPtr *int
	fmt.Printf("空指针的值: %v\n", nilPtr)
	if nilPtr == nil {
		fmt.Println("指针为空")
	}
}

// 多级指针
func multiLevelPointerDemo() {
	fmt.Println("\n=== 多级指针演示 ===")

	var x int = 100
	var ptr *int = &x              // 一级指针
	var ptrPtr **int = &ptr        // 二级指针
	var ptrPtrPtr ***int = &ptrPtr // 三级指针

	fmt.Printf("x 的值: %d\n", x)
	fmt.Printf("*ptr 的值: %d\n", *ptr)
	fmt.Printf("**ptrPtr 的值: %d\n", **ptrPtr)
	fmt.Printf("***ptrPtrPtr 的值: %d\n", ***ptrPtrPtr)

	// 通过三级指针修改原始值
	***ptrPtrPtr = 200
	fmt.Printf("通过三级指针修改后，x 的值: %d\n", x)
}

// 函数指针 (Go中的函数类型)
func functionPointerDemo() {
	fmt.Println("\n=== 函数指针演示 ===")

	// 定义函数类型
	type MathFunc func(int, int) int

	// 定义具体函数
	add := func(a, b int) int {
		return a + b
	}

	multiply := func(a, b int) int {
		return a * b
	}

	// 函数指针变量
	var operation MathFunc

	// 赋值函数指针
	operation = add
	fmt.Printf("加法运算: 5 + 3 = %d\n", operation(5, 3))

	operation = multiply
	fmt.Printf("乘法运算: 5 * 3 = %d\n", operation(5, 3))

	// 函数指针数组
	operations := []MathFunc{add, multiply}
	operationNames := []string{"加法", "乘法"}

	for i, op := range operations {
		result := op(4, 6)
		fmt.Printf("%s运算: 4 和 6 = %d\n", operationNames[i], result)
	}
}

// 指针和接口
type PointerShape interface {
	Area() float64
	Perimeter() float64
}

type PointerRectangle struct {
	Width  float64
	Height float64
}

func (r PointerRectangle) Area() float64 {
	return r.Width * r.Height
}

func (r PointerRectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func interfacePointerDemo() {
	fmt.Println("\n=== 接口和指针演示 ===")

	rect := PointerRectangle{Width: 5.0, Height: 3.0}

	// 值类型实现接口
	var shape1 PointerShape = rect
	fmt.Printf("矩形面积 (值类型): %.2f\n", shape1.Area())

	// 指针类型实现接口
	var shape2 PointerShape = &rect
	fmt.Printf("矩形周长 (指针类型): %.2f\n", shape2.Perimeter())

	// 接口指针
	var shapePtr *PointerShape = &shape1
	fmt.Printf("通过接口指针访问面积: %.2f\n", (*shapePtr).Area())
}

func main4() {
	fmt.Println("=== Go 语言指针和引用学习 ===")

	// 1. 基本指针操作
	fmt.Println("\n=== 基本指针操作 ===")
	var x int = 42
	var ptr *int = &x // 获取 x 的地址

	fmt.Printf("x 的值: %d\n", x)
	fmt.Printf("x 的地址: %p\n", &x)
	fmt.Printf("ptr 指向的地址: %p\n", ptr)
	fmt.Printf("ptr 指向的值: %d\n", *ptr)

	// 通过指针修改值
	*ptr = 100
	fmt.Printf("通过指针修改后，x 的值: %d\n", x)

	// 2. 结构体指针
	fmt.Println("\n=== 结构体指针演示 ===")
	person := PointerPerson{Name: "张三", Age: 25}
	fmt.Printf("原始信息: %s\n", person.GetInfo())

	// 获取结构体指针
	personPtr := &person
	fmt.Printf("通过指针访问: %s\n", (*personPtr).GetInfo())
	fmt.Printf("简化语法访问: %s\n", personPtr.GetInfo()) // Go 自动解引用

	// 通过指针修改结构体
	personPtr.SetAge(30)
	fmt.Printf("修改年龄后: %s\n", person.GetInfo())

	// 3. 函数参数传递
	fmt.Println("\n=== 函数参数传递演示 ===")
	person2 := PointerPerson{Name: "李四", Age: 20}
	fmt.Printf("调用前: %s\n", person2.GetInfo())

	// 值传递 - 不会修改原始值
	updatePersonByValue(person2, 25)
	fmt.Printf("值传递后: %s\n", person2.GetInfo())

	// 指针传递 - 会修改原始值
	updatePersonByPointer(&person2, 25)
	fmt.Printf("指针传递后: %s\n", person2.GetInfo())

	// 4. 返回指针的函数
	fmt.Println("\n=== 返回指针的函数演示 ===")
	newPerson := createPerson("王五", 35)
	fmt.Printf("新创建的人员: %s\n", newPerson.GetInfo())
	newPerson.Birthday()
	fmt.Printf("生日后: %s\n", newPerson.GetInfo())

	// 5. 指针交换
	fmt.Println("\n=== 指针交换演示 ===")
	a, b := 10, 20
	fmt.Printf("交换前: a = %d, b = %d\n", a, b)
	swap(&a, &b)
	fmt.Printf("交换后: a = %d, b = %d\n", a, b)

	// 6. 数组和指针
	arrayPointerDemo()

	// 7. 切片和指针
	slicePointerDemo()

	// 8. 指针和内存地址
	pointerArithmeticDemo()

	// 9. 多级指针
	multiLevelPointerDemo()

	// 10. 函数指针
	functionPointerDemo()

	// 11. 接口和指针
	interfacePointerDemo()

	fmt.Println("\n=== 指针使用注意事项 ===")
	fmt.Println("1. 空指针检查: 使用指针前要检查是否为 nil")
	fmt.Println("2. 内存安全: Go 有垃圾回收，不需要手动释放内存")
	fmt.Println("3. 指针传递: 大结构体建议使用指针传递以提高性能")
	fmt.Println("4. 方法接收者: 需要修改接收者时使用指针接收者")
	fmt.Println("5. 接口实现: 值类型和指针类型都可以实现接口")

	fmt.Println("\n学习完成！")
}
