package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"net/http"
	"time"
	"context"
)

// ============= 命令行计算器伪代码实现 =============

// 栈结构定义
type Stack struct {
	items []interface{}
}

func (s *Stack) Push(item interface{}) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() interface{} {
	if len(s.items) == 0 {
		return nil
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}

func (s *Stack) Peek() interface{} {
	if len(s.items) == 0 {
		return nil
	}
	return s.items[len(s.items)-1]
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

// Token 类型定义
type Token struct {
	Type  string // "number", "operator", "parenthesis"
	Value string
}

// 命令行计算器主函数
func LineCalculator() {
	fmt.Println("=== 命令行计算器 ===")
	fmt.Println("支持 +, -, *, /, () 运算")
	fmt.Println("输入 'quit' 退出")

	for {
		// 1. 使用 fmt.Scan() 读取用户输入
		fmt.Print("请输入表达式: ")
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Println("输入错误，请重试")
			continue
		}

		// 检查退出条件
		if strings.ToLower(input) == "quit" {
			fmt.Println("再见！")
			break
		}

		// 2. 实现字符串解析函数，处理数字和操作符
		tokens, parseErr := tokenize(input)
		if parseErr != nil {
			fmt.Printf("解析错误: %v\n", parseErr)
			continue
		}

		// 3. 使用栈数据结构处理括号和运算符优先级
		// 4. 实现基本的四则运算函数
		result, calcErr := evaluateExpression(tokens)
		if calcErr != nil {
			fmt.Printf("计算错误: %v\n", calcErr)
			continue
		}

		fmt.Printf("结果: %.2f\n\n", result)
	}
}

// 2. 字符串解析函数 - 将输入转换为token序列
func tokenize(input string) ([]Token, error) {
	var tokens []Token
	input = strings.ReplaceAll(input, " ", "") // 去除空格

	for i := 0; i < len(input); i++ {
		ch := input[i]

		// 处理数字（包括小数）
		if unicode.IsDigit(rune(ch)) || ch == '.' {
			numStr := ""
			for i < len(input) && (unicode.IsDigit(rune(input[i])) || input[i] == '.') {
				numStr += string(input[i])
				i++
			}
			i-- // 回退一位，因为for循环会自增

			// 验证数字格式
			if _, err := strconv.ParseFloat(numStr, 64); err != nil {
				return nil, fmt.Errorf("无效的数字格式: %s", numStr)
			}

			tokens = append(tokens, Token{Type: "number", Value: numStr})
		} else if ch == '+' || ch == '-' || ch == '*' || ch == '/' {
			// 处理运算符
			tokens = append(tokens, Token{Type: "operator", Value: string(ch)})
		} else if ch == '(' || ch == ')' {
			// 处理括号
			tokens = append(tokens, Token{Type: "parenthesis", Value: string(ch)})
		} else {
			// 5. 错误处理 - 非法字符
			return nil, fmt.Errorf("非法字符: %c", ch)
		}
	}

	return tokens, nil
}

// 3. 表达式求值函数 - 使用调度场算法（Shunting Yard）
func evaluateExpression(tokens []Token) (float64, error) {
	// 运算符优先级
	precedence := map[string]int{
		"+": 1, "-": 1,
		"*": 2, "/": 2,
	}

	var outputQueue []Token // 输出队列（后缀表达式）
	var operatorStack Stack // 运算符栈

	// 转换为后缀表达式
	for _, token := range tokens {
		switch token.Type {
		case "number":
			outputQueue = append(outputQueue, token)

		case "operator":
			// 处理运算符优先级
			for !operatorStack.IsEmpty() {
				top := operatorStack.Peek().(Token)
				if top.Value == "(" {
					break
				}
				if precedence[top.Value] >= precedence[token.Value] {
					outputQueue = append(outputQueue, operatorStack.Pop().(Token))
				} else {
					break
				}
			}
			operatorStack.Push(token)

		case "parenthesis":
			if token.Value == "(" {
				operatorStack.Push(token)
			} else { // ")"
				// 弹出直到遇到左括号
				for !operatorStack.IsEmpty() {
					top := operatorStack.Pop().(Token)
					if top.Value == "(" {
						break
					}
					outputQueue = append(outputQueue, top)
				}
			}
		}
	}

	// 弹出剩余的运算符
	for !operatorStack.IsEmpty() {
		outputQueue = append(outputQueue, operatorStack.Pop().(Token))
	}

	// 4. 计算后缀表达式
	return evaluatePostfix(outputQueue)
}

// 4. 后缀表达式求值
func evaluatePostfix(tokens []Token) (float64, error) {
	var stack Stack

	for _, token := range tokens {
		if token.Type == "number" {
			num, _ := strconv.ParseFloat(token.Value, 64)
			stack.Push(num)
		} else if token.Type == "operator" {
			// 5. 错误处理 - 检查栈中是否有足够的操作数
			if len(stack.items) < 2 {
				return 0, fmt.Errorf("表达式格式错误")
			}

			b := stack.Pop().(float64)
			a := stack.Pop().(float64)

			result, err := performOperation(a, b, token.Value)
			if err != nil {
				return 0, err
			}

			stack.Push(result)
		}
	}

	// 5. 错误处理 - 检查最终结果
	if len(stack.items) != 1 {
		return 0, fmt.Errorf("表达式格式错误")
	}

	return stack.Pop().(float64), nil
}

// 4. 基本四则运算函数
func performOperation(a, b float64, operator string) (float64, error) {
	switch operator {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		// 5. 错误处理 - 除零检查
		if b == 0 {
			return 0, fmt.Errorf("除零错误")
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("未知运算符: %s", operator)
	}
}

// ============= 示例和测试 =============

func demonstrateCalculator() {
	fmt.Println("\n=== 计算器功能演示 ===")

	testCases := []string{
		"2 + 3",
		"10 - 4 * 2",
		"(1 + 2) * 3",
		"10 / 2 + 3",
		"2.5 * 4",
		"(10 + 2) / (4 - 1)",
	}

	for _, expr := range testCases {
		fmt.Printf("表达式: %s\n", expr)

		tokens, err := tokenize(expr)
		if err != nil {
			fmt.Printf("  解析错误: %v\n", err)
			continue
		}

		result, err := evaluateExpression(tokens)
		if err != nil {
			fmt.Printf("  计算错误: %v\n", err)
			continue
		}

		fmt.Printf("  结果: %.2f\n", result)
	}
}

// 文件处理工具
// **目标**：掌握文件I/O、字符串处理、命令行参数
// **描述**：创建一个文件处理工具，可以统计文件中的行数、字数、字符数

// **流程提示**：
// 1. 使用 `os.Args` 获取命令行参数
// 2. 使用 `os.Open()` 打开文件
// 3. 使用 `bufio.Scanner` 逐行读取文件
// 4. 实现统计函数（行数、单词数、字符数）
// 5. 格式化输出结果

func FileCountCommand() {
	args := os.Args
	fmt.Println("%s", args)

	index := sort.SearchStrings(args, "count")
	if index < 0 {
		fmt.Println("请输入count filepath")
		return
	}
	filepath := args[index+1]
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("打开文件失败: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lineCount, wordCount, charCount int
	for scanner.Scan() {
		lineCount++
		words := strings.Fields(scanner.Text())
		wordCount += len(words)
		charCount += len(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("读取文件失败: %v\n", err)
		return
	}
	fmt.Printf("文件 %s 统计结果:\n", filepath)
	fmt.Printf("  行数: %d\n", lineCount)
	fmt.Printf("  单词数: %d\n", wordCount)
	fmt.Printf("  字符数: %d\n", charCount)
}

// ============= 任务3：简单HTTP服务器 =============
// **目标**：掌握结构体、方法、接口、HTTP处理
// **描述**：实现一个简单的HTTP服务器，提供静态文件服务和简单的API

// **流程提示**：
// 1. 定义服务器结构体，包含端口、根目录等配置
// 2. 实现静态文件服务处理器
// 3. 实现简单的JSON API（如获取服务器状态）
// 4. 使用 `http.HandleFunc` 注册路由
// 5. 启动服务器并处理优雅关闭

// Server 结构体定义 - HTTP服务器的核心配置
type Server struct {
	Port       int                    // 服务器监听端口
	RootDir    string                 // 静态文件根目录
	Status     string                 // 服务器当前状态 (running/stopped)
	StartTime  time.Time              // 服务器启动时间
	Files      map[string]string      // 注册的静态文件映射 (路径 -> 文件名)
	HttpServer *http.Server           // Go标准库HTTP服务器实例
	Router     *http.ServeMux         // HTTP路由多路复用器
}

// NewServer 创建新的服务器实例
// 参数：port - 监听端口，rootDir - 静态文件根目录
func NewServer(port int, rootDir string) *Server {
	return &Server{
		Port:    port,
		RootDir: rootDir,
		Status:  "initialized", // 初始状态
		Files:   make(map[string]string), // 初始化文件映射
		Router:  http.NewServeMux(),      // 创建新的路由器
	}
}

// Start 启动HTTP服务器
// 1. 设置服务器状态和启动时间
// 2. 创建HTTP服务器实例
// 3. 注册默认API路由
// 4. 在goroutine中启动服务器（非阻塞）
func (s *Server) Start() {
	fmt.Printf("正在启动HTTP服务器，端口: %d\n", s.Port)
	
	// 1. 设置服务器运行状态
	s.Status = "running"
	s.StartTime = time.Now()
	
	// 2. 创建HTTP服务器实例
	s.HttpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Port), // 监听地址格式：":8080"
		Handler: s.Router,                   // 使用自定义路由器
	}
	
	// 3. 注册默认API路由
	s.registerDefaultRoutes()
	
	// 4. 在goroutine中启动服务器（避免阻塞主线程）
	go func() {
		fmt.Printf("HTTP服务器已启动: http://localhost:%d\n", s.Port)
		if err := s.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP 服务器启动失败: %v\n", err)
		}
	}()
}

// registerDefaultRoutes 注册默认的API路由
func (s *Server) registerDefaultRoutes() {
	// API路由1：获取服务器状态
	s.RegisterAPI("/api/status", s.handleStatus)
	
	// API路由2：获取注册的文件列表
	s.RegisterAPI("/api/files", s.handleFiles)
	
	// API路由3：根路径欢迎页面
	s.RegisterAPI("/", s.handleHome)
	
	// 静态文件服务器（处理 /static/ 路径下的文件）
	s.Router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(s.RootDir))))
}

// handleStatus 处理服务器状态API请求
// 返回JSON格式的服务器状态信息
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	// 设置响应头为JSON格式
	w.Header().Set("Content-Type", "application/json")
	
	// 计算服务器运行时间
	uptime := time.Since(s.StartTime).Seconds()
	
	// 构造JSON响应
	response := fmt.Sprintf(`{
		"status": "%s",
		"uptime": %.2f,
		"port": %d,
		"start_time": "%s",
		"registered_files": %d
	}`, s.Status, uptime, s.Port, s.StartTime.Format("2006-01-02 15:04:05"), len(s.Files))
	
	// 写入响应
	w.Write([]byte(response))
}

// handleFiles 处理文件列表API请求
// 返回已注册的静态文件列表
func (s *Server) handleFiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// 构造文件列表JSON
	filesJSON := "{"
	count := 0
	for path, filename := range s.Files {
		if count > 0 {
			filesJSON += ","
		}
		filesJSON += fmt.Sprintf(`"%s": "%s"`, path, filename)
		count++
	}
	filesJSON += "}"
	
	w.Write([]byte(filesJSON))
}

// handleHome 处理根路径请求
// 返回简单的HTML欢迎页面
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Go HTTP 服务器</title>
    <meta charset="utf-8">
</head>
<body>
    <h1>🚀 Go HTTP 服务器</h1>
    <p>服务器状态: <strong>` + s.Status + `</strong></p>
    <p>启动时间: <strong>` + s.StartTime.Format("2006-01-02 15:04:05") + `</strong></p>
    <h2>可用的API端点:</h2>
    <ul>
        <li><a href="/api/status">/api/status</a> - 服务器状态信息</li>
        <li><a href="/api/files">/api/files</a> - 注册的文件列表</li>
        <li><a href="/static/">/static/</a> - 静态文件服务</li>
    </ul>
</body>
</html>`
	
	w.Write([]byte(html))
}

// Stop 优雅关闭HTTP服务器
// 使用context.Background()进行优雅关闭，等待现有连接处理完成
func (s *Server) Stop() {
	fmt.Println("正在关闭HTTP服务器...")
	s.Status = "stopped"
	
	// 创建5秒超时的context用于优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() //会在函数结束时执行
	
	// 优雅关闭服务器
	if err := s.HttpServer.Shutdown(ctx); err != nil {
		fmt.Printf("服务器关闭失败: %v\n", err)
	} else {
		fmt.Println("HTTP服务器已成功关闭")
	}
}

// RegisterFile 注册静态文件映射
// 参数：path - URL路径，filename - 对应的文件名
func (s *Server) RegisterFile(path, filename string) {
	s.Files[path] = filename
	fmt.Printf("已注册静态文件: %s -> %s\n", path, filename)
}

// RegisterAPI 注册自定义API路由
// 参数：path - URL路径，handler - 处理函数
func (s *Server) RegisterAPI(path string, handler func(http.ResponseWriter, *http.Request)) {
	s.Router.HandleFunc(path, handler)
	fmt.Printf("已注册API路由: %s\n", path)
}

// DemoHTTPServer 演示HTTP服务器功能
func DemoHTTPServer() {
	fmt.Println("\n=== HTTP服务器演示 ===")
	
	// 1. 创建服务器实例
	server := NewServer(8080, "./static")
	
	// 2. 注册一些示例文件
	server.RegisterFile("/example.txt", "example.txt")
	server.RegisterFile("/readme.md", "README.md")
	
	// 3. 注册自定义API
	server.RegisterAPI("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := `{"message": "Hello from Go HTTP Server!", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`
		w.Write([]byte(response))
	})
	
	// 4. 启动服务器
	server.Start()
	
	// 5. 等待用户输入后关闭服务器
	fmt.Println("\n服务器已启动，访问 http://localhost:8080")
	fmt.Println("按回车键关闭服务器...")
	fmt.Scanln()
	
	// 6. 优雅关闭服务器
	server.Stop()
}

// **目标**：掌握goroutine、channel、sync包
// **描述**：实现一个支持并发的文件下载器，可以同时下载多个文件

// **流程提示**：
// 1. 定义下载任务结构体
// 2. 实现单个文件下载函数
// 3. 使用goroutine实现并发下载
// 4. 使用channel传递下载结果和进度
// 5. 实现下载进度显示
// 6. 添加超时和重试机制
// ============= 任务4：并发下载器 =============
// **目标**：掌握goroutine、channel、并发编程
// **描述**：实现一个支持并发下载的下载器，使用channel进行goroutine间通信

// DownloadTask 下载任务结构体
type DownloadTask struct {
	URL         string        // 下载URL
	Path        string        // 保存路径
	Progress    int           // 下载进度 (0-100)
	Status      string        // 任务状态 (pending/downloading/completed/failed)
	MaxRetries  int           // 最大重试次数
	RetryCount  int           // 当前重试次数
	Timeout     time.Duration // 超时时间
	LastError   error         // 最后一次错误信息
}

// NewDownloadTask 创建新的下载任务
// NewDownloadTask 创建新的下载任务
func NewDownloadTask(url, path string) *DownloadTask {
	return &DownloadTask{
		URL:        url,
		Path:       path,
		Progress:   0,
		Status:     "pending",
		MaxRetries: 3,                    // 默认最大重试3次
		RetryCount: 0,                    // 初始重试次数为0
		Timeout:    30 * time.Second,     // 默认超时30秒
		LastError:  nil,                  // 初始无错误
	}
}

// NewDownloadTaskWithRetry 创建带自定义重试配置的下载任务
func NewDownloadTaskWithRetry(url, path string, maxRetries int, timeout time.Duration) *DownloadTask {
	return &DownloadTask{
		URL:        url,
		Path:       path,
		Progress:   0,
		Status:     "pending",
		MaxRetries: maxRetries,
		RetryCount: 0,
		Timeout:    timeout,
		LastError:  nil,
	}
}

// SetProgress 设置下载进度
func (t *DownloadTask) SetProgress(progress int) {
	t.Progress = progress
}

// SetStatus 设置任务状态
func (t *DownloadTask) SetStatus(status string) {
	t.Status = status
}

// DoDownloadTaskWithRetry 执行带超时和重试机制的下载任务
// 参数：task - 下载任务，results - 结果通道，progress - 进度通道（可选）
func DoDownloadTaskWithRetry(task *DownloadTask, results chan<- *DownloadTask, progress chan<- *DownloadTask) {
	fmt.Printf("🚀 开始下载: %s (最大重试: %d次, 超时: %v)\n", task.URL, task.MaxRetries, task.Timeout)
	
	for task.RetryCount <= task.MaxRetries {
		// 创建带超时的上下文
		ctx, cancel := context.WithTimeout(context.Background(), task.Timeout)
		
		// 尝试下载
		success := attemptDownload(ctx, task, progress)
		cancel() // 释放上下文资源
		
		if success {
			// 下载成功
			task.SetStatus("completed")
			fmt.Printf("✅ 下载成功: %s (重试次数: %d)\n", task.URL, task.RetryCount)
			results <- task
			return
		}
		
		// 下载失败，检查是否还能重试
		task.RetryCount++
		if task.RetryCount <= task.MaxRetries {
			waitTime := time.Duration(task.RetryCount) * 2 * time.Second // 指数退避
			fmt.Printf("⚠️  下载失败，%v后进行第%d次重试: %s\n", waitTime, task.RetryCount, task.URL)
			time.Sleep(waitTime)
		}
	}
	
	// 所有重试都失败了
	task.SetStatus("failed")
	fmt.Printf("❌ 下载最终失败: %s (已重试%d次)\n", task.URL, task.RetryCount-1)
	results <- task
}

// attemptDownload 尝试执行单次下载
// 参数：ctx - 上下文（用于超时控制），task - 下载任务，progress - 进度通道
// 返回：bool - 是否下载成功
func attemptDownload(ctx context.Context, task *DownloadTask, progress chan<- *DownloadTask) bool {
	task.SetStatus("downloading")
	task.SetProgress(0)
	
	// 模拟下载过程，支持超时取消
	for i := 0; i <= 100; i += 20 {
		select {
		case <-ctx.Done():
			// 超时或取消
			task.LastError = ctx.Err()
			fmt.Printf("⏰ 下载超时: %s (进度: %d%%)\n", task.URL, task.Progress)
			return false
		default:
			// 继续下载
			task.SetProgress(i)
			
			// 发送进度更新（如果有进度通道）
			if progress != nil {
				select {
				case progress <- task:
				default: // 非阻塞发送
				}
			}
			
			// 模拟下载耗时
			time.Sleep(300 * time.Millisecond)
			
			// 模拟随机失败（用于测试重试机制）
			if i == 60 && task.RetryCount < 2 {
				task.LastError = fmt.Errorf("模拟网络错误")
				fmt.Printf("🔌 网络错误: %s (进度: %d%%)\n", task.URL, task.Progress)
				return false
			}
		}
	}
	
	task.SetProgress(100)
	task.LastError = nil
	return true
}

// DownLoadWithRetry 带重试机制的并发下载多个任务
// 演示Go语言的channel通信机制和超时重试
func DownLoadWithRetry(tasks []*DownloadTask) {
	fmt.Printf("\n=== 带重试机制的并发下载 ===\n")
	fmt.Printf("任务总数: %d\n", len(tasks))
	
	// 创建带缓冲的通道，避免阻塞
	results := make(chan *DownloadTask, len(tasks))
	
	// 启动并发下载，使用带重试的版本
	for _, task := range tasks {
		go DoDownloadTaskWithRetry(task, results, nil) // 不使用进度通道
	}
	
	// 接收所有完成的任务
	successCount := 0
	failedCount := 0
	
	for i := 0; i < len(tasks); i++ {
		completedTask := <-results
		if completedTask.Status == "completed" {
			successCount++
			fmt.Printf("✅ 任务完成 [%d/%d]: %s (重试次数: %d)\n", 
				i+1, len(tasks), completedTask.URL, completedTask.RetryCount)
		} else {
			failedCount++
			fmt.Printf("❌ 任务失败 [%d/%d]: %s (重试次数: %d, 错误: %v)\n", 
				i+1, len(tasks), completedTask.URL, completedTask.RetryCount-1, completedTask.LastError)
		}
	}
	
	fmt.Printf("\n=== 下载统计 ===\n")
	fmt.Printf("总任务数: %d\n", len(tasks))
	fmt.Printf("成功完成: %d\n", successCount)
	fmt.Printf("失败任务: %d\n", failedCount)
}
func DownLoad(tasks []*DownloadTask) {
	fmt.Printf("\n=== 开始并发下载 %d 个任务 ===\n", len(tasks))
	
	// 1. 创建一个channel用于接收下载结果
	// 使用缓冲channel，容量等于任务数量，避免goroutine阻塞
	results := make(chan *DownloadTask, len(tasks))
	
	// 2. 启动多个goroutine并发下载
	fmt.Println("启动并发下载goroutines...")
	for i, task := range tasks {
		fmt.Printf("启动goroutine %d: %s\n", i+1, task.URL)
		// 🔑 关键：将results通道传递给DoDownloadTaskWithRetry
		go DoDownloadTaskWithRetry(task, results, nil)
	}
	
	// 3. 从channel接收下载结果
	fmt.Println("\n等待下载完成...")
	completedTasks := make([]*DownloadTask, 0, len(tasks))
	
	// 接收所有任务的完成通知
	for i := 0; i < len(tasks); i++ {
		// 🔑 关键：从results通道接收完成的任务
		completedTask := <-results
		completedTasks = append(completedTasks, completedTask)
		
		fmt.Printf("✅ 任务完成 [%d/%d]: %s (进度: %d%%, 状态: %s)\n", 
			i+1, len(tasks), completedTask.URL, completedTask.Progress, completedTask.Status)
	}
	
	// 4. 关闭channel（可选，因为没有其他goroutine在等待）
	close(results)
	
	// 5. 显示下载统计
	fmt.Printf("\n=== 下载统计 ===\n")
	successful := 0
	for _, task := range completedTasks {
		if task.Status == "completed" {
			successful++
		}
	}
	fmt.Printf("总任务数: %d\n", len(tasks))
	fmt.Printf("成功完成: %d\n", successful)
	fmt.Printf("失败任务: %d\n", len(tasks)-successful)
}

// DownLoadWithProgressAndRetry 带进度显示和重试机制的并发下载
// 演示更复杂的channel通信模式和超时重试
func DownLoadWithProgressAndRetry(tasks []*DownloadTask) {
	fmt.Printf("\n=== 带进度显示和重试机制的并发下载 ===\n")
	
	// 创建两个通道：进度通道和完成通道
	progress := make(chan *DownloadTask, 10)  // 缓冲进度更新
	completed := make(chan *DownloadTask, len(tasks)) // 完成通知
	
	// 启动所有下载任务
	for _, task := range tasks {
		go DoDownloadTaskWithRetry(task, completed, progress)
	}
	
	// 使用select语句同时处理进度更新和任务完成
	active := len(tasks)
	successCount := 0
	failedCount := 0
	
	for active > 0 {
		select {
		case task := <-progress:
			// 处理进度更新
			fmt.Printf("📊 进度更新: %s -> %d%% (重试: %d)\n", 
				task.URL, task.Progress, task.RetryCount)
			
		case task := <-completed:
			// 处理任务完成
			if task.Status == "completed" {
				successCount++
				fmt.Printf("✅ 任务完成: %s (重试次数: %d)\n", task.URL, task.RetryCount)
			} else {
				failedCount++
				fmt.Printf("❌ 任务失败: %s (重试次数: %d, 错误: %v)\n", 
					task.URL, task.RetryCount-1, task.LastError)
			}
			active--
		}
	}
	
	fmt.Printf("\n=== 最终统计 ===\n")
	fmt.Printf("总任务数: %d\n", len(tasks))
	fmt.Printf("成功完成: %d\n", successCount)
	fmt.Printf("失败任务: %d\n", failedCount)
}

// DemoDownloaderWithRetry 演示带超时重试机制的下载器功能
func DemoDownloaderWithRetry() {
	fmt.Println("\n=== 带超时重试机制的并发下载器演示 ===")
	
	// 创建不同配置的下载任务来演示重试机制
	tasks := []*DownloadTask{
		// 正常任务 - 默认配置
		NewDownloadTask("https://example.com/file1.zip", "./downloads/file1.zip"),
		
		// 短超时任务 - 容易超时，需要重试
		NewDownloadTaskWithRetry("https://example.com/slow-file.pdf", "./downloads/slow-file.pdf", 2, 1*time.Second),
		
		// 高重试任务 - 最多重试5次
		NewDownloadTaskWithRetry("https://example.com/unstable-file.mp4", "./downloads/unstable-file.mp4", 5, 5*time.Second),
		
		// 零重试任务 - 不重试
		NewDownloadTaskWithRetry("https://example.com/fail-file.jpg", "./downloads/fail-file.jpg", 0, 2*time.Second),
	}
	
	fmt.Println("\n📋 任务配置:")
	for i, task := range tasks {
		fmt.Printf("任务%d: %s (最大重试: %d次, 超时: %v)\n", 
			i+1, task.URL, task.MaxRetries, task.Timeout)
	}
	
	// 方式1：基本重试下载
	fmt.Println("\n🔄 方式1：基本重试下载")
	DownLoadWithRetry(tasks)
	
	// 等待用户确认
	fmt.Println("\n按回车键继续演示带进度显示的重试下载...")
	fmt.Scanln()
	
	// 重置任务状态用于第二次演示
	for _, task := range tasks {
		task.Progress = 0
		task.Status = "pending"
		task.RetryCount = 0
		task.LastError = nil
	}
	
	// 方式2：带进度显示的重试下载
	fmt.Println("\n📊 方式2：带进度显示和重试机制的下载")
	DownLoadWithProgressAndRetry(tasks)
	
	// 演示总结
	fmt.Println("\n=== 超时重试机制特性总结 ===")
	fmt.Println("✅ 支持自定义超时时间")
	fmt.Println("✅ 支持自定义最大重试次数") 
	fmt.Println("✅ 指数退避重试策略 (2s, 4s, 6s...)")
	fmt.Println("✅ 上下文超时控制")
	fmt.Println("✅ 详细的错误信息记录")
	fmt.Println("✅ 实时进度更新")
	fmt.Println("✅ 并发安全的通道通信")
}

func main7() {
	fmt.Println("Go 语言学习任务")
	fmt.Println("==================")

	// 演示计算器功能
	// demonstrateCalculator()

	// // 启动交互式计算器
	// fmt.Println("\n按回车键启动交互式计算器...")
	// fmt.Scanln()
	// LineCalculator()

	// 演示文件处理工具
	// FileCountCommand()
	
	// 演示HTTP服务器功能
	// DemoHTTPServer()
	
	// 演示并发下载器功能（带超时重试机制）
	DemoDownloaderWithRetry()
}
