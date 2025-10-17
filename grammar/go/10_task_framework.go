package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"reflect"
	"sync"
	"time"
	"math/rand"
	"encoding/json"
	"net/http"
)

// ============================================================================
// 任务7：微服务框架
// **目标**：掌握反射、代码生成、中间件模式
// **描述**：实现一个简单的微服务框架，支持服务注册、发现、负载均衡
//
// **流程提示**：
// 1. 设计服务接口和注册机制
// 2. 实现基于反射的RPC调用
// 3. 实现服务发现和注册中心
// 4. 添加负载均衡算法
// 5. 实现中间件机制（日志、认证、限流）
// 6. 添加健康检查和故障转移
// ============================================================================

// ============================================================================
// 1. 核心接口和数据结构
// ============================================================================

// ServiceInfo 服务信息
type ServiceInfo struct {
	Name     string            `json:"name"`
	Version  string            `json:"version"`
	Address  string            `json:"address"`
	Port     int               `json:"port"`
	Metadata map[string]string `json:"metadata"`
	Health   HealthStatus      `json:"health"`
}

// HealthStatus 健康状态
type HealthStatus struct {
	Status    string    `json:"status"` // "healthy", "unhealthy", "unknown"
	LastCheck time.Time `json:"last_check"`
	Message   string    `json:"message"`
}

// ServiceRegistry 服务注册接口
type ServiceRegistry interface {
	
	Register(ctx context.Context, service *ServiceInfo) error
	
	
	Unregister(ctx context.Context, serviceName, serviceID string) error
	
	
	Discover(ctx context.Context, serviceName string) ([]*ServiceInfo, error)
	
	// TODO: 实现服务监听逻辑
	Watch(ctx context.Context, serviceName string) (<-chan []*ServiceInfo, error)
}

// LoadBalancer 负载均衡接口
type LoadBalancer interface {
	// TODO: 实现负载均衡选择逻辑
	Select(services []*ServiceInfo) (*ServiceInfo, error)
	
	// TODO: 实现权重更新逻辑
	UpdateWeights(serviceName string, weights map[string]int) error
}

// RPCClient RPC客户端接口
type RPCClient interface {
	// TODO: 实现RPC调用逻辑
	Call(ctx context.Context, serviceName, method string, args interface{}, reply interface{}) error
	
	// TODO: 实现连接关闭逻辑
	Close() error
}

// RPCServer RPC服务端接口
type RPCServer interface {
	// TODO: 实现服务注册逻辑
	RegisterService(serviceName string, service interface{}) error
	
	// TODO: 实现服务启动逻辑
	Start(address string) error
	
	// TODO: 实现服务停止逻辑
	Stop() error
}

// ============================================================================
// 2. 服务注册中心实现框架
// ============================================================================

// InMemoryRegistry 内存服务注册中心
type InMemoryRegistry struct {
	mu       sync.RWMutex
	services map[string][]*ServiceInfo // serviceName -> []*ServiceInfo
	watchers map[string][]chan []*ServiceInfo
}

func NewInMemoryRegistry() *InMemoryRegistry {
	return &InMemoryRegistry{
		services: make(map[string][]*ServiceInfo),
		watchers: make(map[string][]chan []*ServiceInfo),
	}
}

func (r *InMemoryRegistry) Register(ctx context.Context, service *ServiceInfo) error {
	// 1. 验证服务信息
	if service == nil {
		return fmt.Errorf("service cannot be nil")
	}
	if service.Name == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	if service.Address == "" {
		return fmt.Errorf("service address cannot be empty")
	}
	if service.Port <= 0 || service.Port > 65535 {
		return fmt.Errorf("service port must be between 1 and 65535, got %d", service.Port)
	}
	
	// 设置默认值
	if service.Version == "" {
		service.Version = "1.0.0"
	}
	if service.Metadata == nil {
		service.Metadata = make(map[string]string)
	}
	if service.Health.Status == "" {
		service.Health = HealthStatus{
			Status:    "healthy",
			LastCheck: time.Now(),
			Message:   "newly registered",
		}
	}
	
	// 生成服务ID（如果没有提供）
	serviceID := fmt.Sprintf("%s-%s:%d", service.Name, service.Address, service.Port)
	service.Metadata["service_id"] = serviceID
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// 2. 添加到服务列表
	if r.services[service.Name] == nil {
		r.services[service.Name] = make([]*ServiceInfo, 0)
	}
	
	// 检查是否已存在相同的服务实例
	for i, existing := range r.services[service.Name] {
		if existing.Address == service.Address && existing.Port == service.Port {
			// 更新现有服务信息
			r.services[service.Name][i] = service
			log.Printf("Updated existing service %s at %s:%d", service.Name, service.Address, service.Port)
			
			// 3. 通知所有监听者
			r.notifyWatchers(service.Name)
			return nil
		}
	}
	
	// 添加新服务实例
	r.services[service.Name] = append(r.services[service.Name], service)
	log.Printf("Registered new service %s at %s:%d", service.Name, service.Address, service.Port)
	
	// 3. 通知所有监听者
	r.notifyWatchers(service.Name)
	
	return nil
}

// notifyWatchers 通知所有监听指定服务的watchers
func (r *InMemoryRegistry) notifyWatchers(serviceName string) {
	if watchers, exists := r.watchers[serviceName]; exists {
		services := r.services[serviceName]
		for _, watcher := range watchers {
			select {
			case watcher <- services:
				// 成功发送通知
			default:
				// 通道已满，跳过此watcher
				log.Printf("Warning: watcher channel full for service %s", serviceName)
			}
		}
	}
}

func (r *InMemoryRegistry) Unregister(ctx context.Context, serviceName, serviceID string) error {
	// 1. 验证参数
	if serviceName == "" || serviceID == "" {
		return fmt.Errorf("serviceName and serviceID cannot be empty")
	}
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// 2. 检查服务是否存在
	services, exists := r.services[serviceName]
	if !exists {
		return fmt.Errorf("service %s not found", serviceName)
	}
	
	// 3. 查找并删除指定ID的服务实例
	found := false
	for i, service := range services {
		if service.Metadata["service_id"] == serviceID {
			// 从切片中删除该服务实例
			r.services[serviceName] = append(services[:i], services[i+1:]...)
			found = true
			log.Printf("Unregistered service instance %s from service %s", serviceID, serviceName)
			break
		}
	}
	if !found {
		return fmt.Errorf("service instance with ID %s not found in service %s", serviceID, serviceName)
	}
	
	// 4. 如果没有服务实例了，删除整个服务条目
	if len(r.services[serviceName]) == 0 {
		delete(r.services, serviceName)
		log.Printf("Removed empty service entry for %s", serviceName)
	}
	// 5. 通知所有监听者
	r.notifyWatchers(serviceName)
	
	return nil
}

func (r *InMemoryRegistry) Discover(ctx context.Context, serviceName string) ([]*ServiceInfo, error) {
	// 1. 从服务列表中查找
	services, exists := r.services[serviceName]
	if !exists {
		return nil, fmt.Errorf("service %s not found", serviceName)
	}
	healthyServices := make([]*ServiceInfo, 0)
	// 2. 过滤健康的服务
	for _, service := range services {
		if service.Health.Status == "healthy" {
			healthyServices = append(healthyServices, service)
		}
	}
	if len(healthyServices) == 0 {
		return nil, fmt.Errorf("no healthy services found for %s", serviceName)
	}
	
	log.Printf(": Discover services for %s", serviceName)
	return healthyServices, nil
}

func (r *InMemoryRegistry) Watch(ctx context.Context, serviceName string) (<-chan []*ServiceInfo, error) {
	// 1. 参数验证
	if serviceName == "" {
		return nil, fmt.Errorf("serviceName cannot be empty")
	}
	
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}
	
	// 2. 创建监听通道
	ch := make(chan []*ServiceInfo, 10)
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// 3. 添加到监听者列表
	if r.watchers[serviceName] == nil {
		r.watchers[serviceName] = make([]chan []*ServiceInfo, 0)
	}
	r.watchers[serviceName] = append(r.watchers[serviceName], ch)
	
	// 4. 发送当前服务状态（如果存在）
	if services, exists := r.services[serviceName]; exists {
		// 创建服务副本，避免并发修改
		servicesCopy := make([]*ServiceInfo, len(services))
		copy(servicesCopy, services)
		
		// 非阻塞发送初始数据
		select {
		case ch <- servicesCopy:
			log.Printf("Sent initial services for %s: %d instances", serviceName, len(servicesCopy))
		default:
			log.Printf("Warning: failed to send initial data for %s", serviceName)
		}
	}
	
	// 5. 启动goroutine处理context取消
	go func() {
		<-ctx.Done() // 阻塞等待context被取消
		r.removeWatcher(serviceName, ch)
		close(ch)
		log.Printf("Watcher for service %s removed due to context cancellation", serviceName)
	}()
	
	log.Printf("Started watching service: %s", serviceName)
	return ch, nil
}

// removeWatcher 从监听者列表中移除指定的watcher
func (r *InMemoryRegistry) removeWatcher(serviceName string, targetCh chan []*ServiceInfo) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	watchers := r.watchers[serviceName]
	for i, ch := range watchers {
		if ch == targetCh {
			// 从切片中移除该watcher
			r.watchers[serviceName] = append(watchers[:i], watchers[i+1:]...)
			
			// 如果没有监听者了，删除该服务的监听者列表
			if len(r.watchers[serviceName]) == 0 {
				delete(r.watchers, serviceName)
			}
			break
		}
	}
}

// ============================================================================
// 3. 负载均衡实现框架
// ============================================================================

// RoundRobinBalancer 轮询负载均衡器
type RoundRobinBalancer struct {
	mu      sync.Mutex
	counter map[string]int // serviceName -> counter
}

func NewRoundRobinBalancer() *RoundRobinBalancer {
	return &RoundRobinBalancer{
		counter: make(map[string]int),
	}
}

func (b *RoundRobinBalancer) Select(services []*ServiceInfo) (*ServiceInfo, error) {
	// 1. 过滤健康的服务
	healthyServices := make([]*ServiceInfo, 0)
	for _, service := range services {
		if service.Health.Status == "healthy" {
			healthyServices = append(healthyServices, service)
		}
	}
	if len(healthyServices) == 0 {
		return nil, fmt.Errorf("no healthy services available")
	}
	
	// 2. 轮询选择
	b.mu.Lock()
	defer b.mu.Unlock()
	
	serviceName := healthyServices[0].Name
	counter := b.counter[serviceName]
	selected := healthyServices[counter % len(healthyServices)]
	b.counter[serviceName] = counter + 1
	
	log.Printf("Selected service %s at %s:%d from %d available services", 
		selected.Name, selected.Address, selected.Port, len(healthyServices))
	
	return selected, nil
}

func (b *RoundRobinBalancer) UpdateWeights(serviceName string, weights map[string]int) error {
	// RoundRobinBalancer 不使用权重，所以这个方法是空实现
	// 但需要实现以满足 LoadBalancer 接口
	return nil
}


// WeightedBalancer 加权负载均衡器
type WeightedBalancer struct {
	mu      sync.RWMutex
	weights map[string]map[string]int // serviceName -> serviceID -> weight
}

func NewWeightedBalancer() *WeightedBalancer {
	return &WeightedBalancer{
		weights: make(map[string]map[string]int),
	}
}

func (b *WeightedBalancer) Select(services []*ServiceInfo) (*ServiceInfo, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	
	// 1. 过滤健康的服务
	healthyServices := make([]*ServiceInfo, 0)
	for _, service := range services {
		if service.Health.Status == "healthy" {
			healthyServices = append(healthyServices, service)
		}
	}
	if len(healthyServices) == 0 {
		return nil, fmt.Errorf("no healthy services available")
	}
	
	// 2. 计算总权重
	totalWeight := 0
	for _, service := range healthyServices {
		serviceName := service.Name
		serviceID := service.Metadata["service_id"]
		
		// 检查权重配置是否存在
		if serviceWeights, exists := b.weights[serviceName]; exists {
			if weight, exists := serviceWeights[serviceID]; exists && weight > 0 {
				totalWeight += weight
			} else {
				// 如果没有配置权重或权重为0，使用默认权重1
				totalWeight += 1
			}
		} else {
			// 如果服务没有权重配置，使用默认权重1
			totalWeight += 1
		}
	}
	
	if totalWeight == 0 {
		// 如果所有权重都为0，使用轮询方式选择第一个
		return healthyServices[0], nil
	}
	
	// 3. 加权随机选择
	rand.Seed(time.Now().UnixNano())
	threshold := rand.Intn(totalWeight)
	
	for _, service := range healthyServices {
		serviceName := service.Name
		serviceID := service.Metadata["service_id"]
		
		weight := 1 // 默认权重
		if serviceWeights, exists := b.weights[serviceName]; exists {
			if w, exists := serviceWeights[serviceID]; exists && w > 0 {
				weight = w
			}
		}
		
		threshold -= weight
		if threshold < 0 {
			return service, nil
		}
	}
	
	// 兜底返回第一个服务
	return healthyServices[0], nil
}

func (b *WeightedBalancer) UpdateWeights(serviceName string, weights map[string]int) error {
	// TODO: 实现权重更新逻辑
	log.Printf("TODO: Update weights for service %s: %v", serviceName, weights)
	// 1. 检查服务是否存在
	if _, exists := b.weights[serviceName]; !exists {
		b.weights[serviceName] = make(map[string]int)
	}
	// 2. 更新权重
	for serviceID, weight := range weights {
		b.weights[serviceName][serviceID] = weight
	}
	return nil
}

// SetHealthChecker 设置健康检查器
func (s *DefaultRPCServer) SetHealthChecker(checker HealthChecker) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.healthChecker = checker
}

// ============================================================================
// 4. RPC客户端实现框架
// ============================================================================

// DefaultRPCClient 默认RPC客户端
type DefaultRPCClient struct {
	registry     ServiceRegistry
	loadBalancer LoadBalancer
	connections  map[string]net.Conn // serviceAddress -> connection
	mu           sync.RWMutex
}

func NewRPCClient(registry ServiceRegistry, loadBalancer LoadBalancer) *DefaultRPCClient {
	return &DefaultRPCClient{
		registry:     registry,
		loadBalancer: loadBalancer,
		connections:  make(map[string]net.Conn),
	}
}

func (c *DefaultRPCClient) Call(ctx context.Context, serviceName, method string, args interface{}, reply interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 1. 服务发现
	services, err := c.registry.Discover(ctx, serviceName)
	if err != nil {
		log.Printf("Service discovery failed for %s: %v", serviceName, err)
		return fmt.Errorf("service discovery failed: %w", err)
	}
	
	if len(services) == 0 {
		log.Printf("No services found for %s", serviceName)
		return fmt.Errorf("no services found for %s", serviceName)
	}
	
	log.Printf("Found %d services for %s", len(services), serviceName)
	
	// 2. 负载均衡选择
	selected, err := c.loadBalancer.Select(services)
	if err != nil {
		log.Printf("Load balancing failed for %s: %v", serviceName, err)
		return fmt.Errorf("load balancing failed: %w", err)
	}
	
	log.Printf("Selected service %s at %s:%d", selected.Name, selected.Address, selected.Port)
	
	// 3. 建立连接
	address := fmt.Sprintf("%s:%d", selected.Address, selected.Port)
	conn := c.connections[address]
	if conn == nil {
		log.Printf("Establishing new connection to %s", address)
		conn, err = net.Dial("tcp", address)
		if err != nil {
			log.Printf("Connection failed to %s: %v", address, err)
			return fmt.Errorf("connection failed to %s: %w", address, err)
		}
		c.connections[address] = conn
		log.Printf("Successfully connected to %s", address)
	} else {
		log.Printf("Reusing existing connection to %s", address)
	}
	
	// 4. 序列化请求
	req := &RPCRequest{
		ServiceName: serviceName,
		Method:      method,
		Args:        args,
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Request serialization failed: %v", err)
		return fmt.Errorf("request serialization failed: %w", err)
	}
	
	// 5. 发送请求
	log.Printf("Sending RPC request: %s.%s", serviceName, method)
	_, err = conn.Write(reqBytes)
	if err != nil {
		// 连接失败，删除并重试
		log.Printf("Request sending failed to %s: %v", address, err)
		delete(c.connections, address)
		conn.Close()
		return fmt.Errorf("request sending failed to %s: %w", address, err)
	}
	
	// 6. 接收响应
	respBytes := make([]byte, 1024)
	n, err := conn.Read(respBytes)
	if err != nil {
		// 连接失败，删除并重试
		log.Printf("Response reception failed from %s: %v", address, err)
		delete(c.connections, address)
		conn.Close()
		return fmt.Errorf("response reception failed from %s: %w", address, err)
	}
	respBytes = respBytes[:n]
	
	// 7. 反序列化响应
	resp := &RPCResponse{}
	if err := json.Unmarshal(respBytes, resp); err != nil {
		log.Printf("Response deserialization failed: %v", err)
		return fmt.Errorf("response deserialization failed: %w", err)
	}
	
	// 8. 处理响应
	if resp.Error != "" {
		log.Printf("Remote error from %s.%s: %s", serviceName, method, resp.Error)
		return fmt.Errorf("remote error: %s", resp.Error)
	}
	
	// 9. 将结果复制到reply
	if reply != nil {
		replyBytes, err := json.Marshal(resp.Result)
		if err != nil {
			log.Printf("Result serialization failed: %v", err)
			return fmt.Errorf("result serialization failed: %w", err)
		}
		if err := json.Unmarshal(replyBytes, reply); err != nil {
			log.Printf("Result deserialization failed: %v", err)
			return fmt.Errorf("result deserialization failed: %w", err)
		}
	}
	
	log.Printf("RPC call %s.%s completed successfully", serviceName, method)
	return nil
}

func (c *DefaultRPCClient) Close() error {
	// TODO: 实现连接关闭逻辑
	log.Printf("TODO: Close all connections")
	for _, conn := range c.connections {
		conn.Close()
	}
	// 8. 清除连接映射
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connections = make(map[string]net.Conn)
	return nil
}

// ============================================================================
// 5. RPC服务端实现框架
// ============================================================================

// DefaultRPCServer RPC服务端的默认实现
type DefaultRPCServer struct {
	services         map[string]interface{} // serviceName -> service instance
	listener         net.Listener
	mu               sync.RWMutex
	registry         ServiceRegistry        // 服务注册中心
	serverInfo       *ServiceInfo          // 服务端自身信息
	registeredServices map[string]string   // serviceName -> serviceID (已注册的服务)
	healthChecker    HealthChecker         // 健康检查器
}

func NewRPCServer() *DefaultRPCServer {
	return &DefaultRPCServer{
		services:           make(map[string]interface{}),
		registeredServices: make(map[string]string),
	}
}

// NewRPCServerWithRegistry 创建带注册中心的RPC服务端
func NewRPCServerWithRegistry(registry ServiceRegistry) *DefaultRPCServer {
	return &DefaultRPCServer{
		services:           make(map[string]interface{}),
		registeredServices: make(map[string]string),
		registry:           registry,
		healthChecker:      NewTCPHealthChecker(5 * time.Second), // 默认使用TCP健康检查
	}
}

func (s *DefaultRPCServer) RegisterService(serviceName string, service interface{}) error {
	if serviceName == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	if service == nil {
		return fmt.Errorf("service cannot be nil")
	}

	// 验证服务接口
	serviceType := reflect.TypeOf(service)
	if serviceType.Kind() != reflect.Ptr {
		return fmt.Errorf("service must be a pointer")
	}

	// 验证服务方法
	numMethods := serviceType.NumMethod()
	if numMethods == 0 {
		return fmt.Errorf("service has no exported methods")
	}

	// 检查每个方法的签名
	for i := 0; i < numMethods; i++ {
		method := serviceType.Method(i)
		methodType := method.Type

		// 检查方法签名: func(receiver, context.Context, *ArgsType) (*ReplyType, error)
		if methodType.NumIn() != 3 {
			log.Printf("Warning: Method %s.%s has %d parameters, expected 3 (receiver, context, args)", 
				serviceName, method.Name, methodType.NumIn())
			continue
		}

		// 检查第二个参数是否为context.Context
		if methodType.In(1) != reflect.TypeOf((*context.Context)(nil)).Elem() {
			log.Printf("Warning: Method %s.%s second parameter should be context.Context", 
				serviceName, method.Name)
			continue
		}

		// 检查返回值数量
		if methodType.NumOut() != 2 {
			log.Printf("Warning: Method %s.%s has %d return values, expected 2 (result, error)", 
				serviceName, method.Name, methodType.NumOut())
			continue
		}

		// 检查最后一个返回值是否为error
		errorType := reflect.TypeOf((*error)(nil)).Elem()
		if !methodType.Out(1).Implements(errorType) {
			log.Printf("Warning: Method %s.%s second return value should implement error interface", 
				serviceName, method.Name)
			continue
		}

		log.Printf("Validated method: %s.%s", serviceName, method.Name)
	}

	s.mu.Lock()
	s.services[serviceName] = service
	s.mu.Unlock()

	log.Printf("Registered service: %s with %d methods", serviceName, numMethods)
	return nil
}

func (s *DefaultRPCServer) Start(address string) error {
	// 1. 监听端口
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("listen failed: %w", err)
	}
	s.listener = listener
	
	log.Printf("RPC server started on %s", address)
	
	// 2. 解析地址信息
	host, portStr, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("invalid address format: %w", err)
	}
	
	// 转换端口为整数
	port := 0
	if _, err := fmt.Sscanf(portStr, "%d", &port); err != nil {
		return fmt.Errorf("invalid port: %w", err)
	}
	
	// 3. 如果配置了注册中心，自动注册已注册的服务
	if s.registry != nil {
		s.mu.RLock()
		for serviceName := range s.services {
			s.mu.RUnlock()
			
			// 创建服务信息
			serviceInfo := &ServiceInfo{
				Name:    serviceName,
				Version: "1.0.0",
				Address: host,
				Port:    port,
				Metadata: map[string]string{
					"service_id": fmt.Sprintf("%s-%s:%d", serviceName, host, port),
					"protocol":   "rpc",
					"started_at": time.Now().Format(time.RFC3339),
				},
				Health: HealthStatus{
					Status:    "unknown",
					LastCheck: time.Now(),
					Message:   "Service starting, health check pending",
				},
			}
			
			// 执行健康检查
			if s.healthChecker != nil {
				healthStatus := s.healthChecker.Check(context.Background(), serviceInfo)
				serviceInfo.Health = healthStatus
				log.Printf("Health check for service %s: %s - %s", serviceName, healthStatus.Status, healthStatus.Message)
			}
			
			// 注册到注册中心
			ctx := context.Background()
			if err := s.registry.Register(ctx, serviceInfo); err != nil {
				log.Printf("Failed to register service %s to registry: %v", serviceName, err)
			} else {
				s.mu.Lock()
				s.registeredServices[serviceName] = serviceInfo.Metadata["service_id"]
				s.mu.Unlock()
				log.Printf("Auto-registered service %s to registry at %s:%d", serviceName, host, port)
			}
			
			s.mu.RLock()
		}
		s.mu.RUnlock()
	}
	
	// 4. 接受连接并处理请求
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Accept error: %v", err)
				continue
			}
			// 为每个连接启动一个goroutine处理请求
			go s.handleConnection(conn)
		}
	}()

	return nil
}

// handleConnection 处理单个连接的请求
func (s *DefaultRPCServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)
	
	for {
		// 读取请求
		var req RPCRequest
		if err := decoder.Decode(&req); err != nil {
			log.Printf("Request decode error: %v", err)
			return
		}
		
		// 处理请求
		resp := s.processRequest(context.Background(), &req)
		
		// 发送响应
		if err := encoder.Encode(resp); err != nil {
			log.Printf("Response encode error: %v", err)
			return
		}
	}
}

// processRequest 处理RPC请求
func (s *DefaultRPCServer) processRequest(ctx context.Context, req *RPCRequest) *RPCResponse {
	s.mu.RLock()
	service, exists := s.services[req.ServiceName]
	s.mu.RUnlock()
	
	if !exists {
		return &RPCResponse{
			Error: fmt.Sprintf("service '%s' not found", req.ServiceName),
		}
	}
	
	// 使用反射调用方法
	serviceValue := reflect.ValueOf(service)
	serviceType := serviceValue.Type()
	
	// 查找方法
	method, exists := serviceType.MethodByName(req.Method)
	if !exists {
		return &RPCResponse{
			Error: fmt.Sprintf("method '%s' not found in service '%s'", req.Method, req.ServiceName),
		}
	}
	
	// 验证方法签名: func(receiver, ctx context.Context, args *ArgsType) (*ResultType, error)
	methodType := method.Type
	if methodType.NumIn() != 3 || methodType.NumOut() != 2 {
		return &RPCResponse{
			Error: fmt.Sprintf("invalid method signature for '%s'", req.Method),
		}
	}
	
	// 检查参数类型
	if !methodType.In(1).Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
		return &RPCResponse{
			Error: fmt.Sprintf("method '%s' second parameter must be context.Context", req.Method),
		}
	}
	
	// 检查返回值类型
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	if !methodType.Out(1).Implements(errorInterface) {
		return &RPCResponse{
			Error: fmt.Sprintf("method '%s' second return value must be error", req.Method),
		}
	}
	
	// 准备参数
	argsType := methodType.In(2)
	var argsValue reflect.Value
	
	if req.Args != nil {
		// 将args转换为正确的类型
		argsBytes, err := json.Marshal(req.Args)
		if err != nil {
			return &RPCResponse{
				Error: fmt.Sprintf("args serialization failed: %v", err),
			}
		}
		
		// 创建参数实例
		if argsType.Kind() == reflect.Ptr {
			argsValue = reflect.New(argsType.Elem())
		} else {
			argsValue = reflect.New(argsType).Elem()
		}
		
		// 反序列化参数
		if err := json.Unmarshal(argsBytes, argsValue.Interface()); err != nil {
			return &RPCResponse{
				Error: fmt.Sprintf("args deserialization failed: %v", err),
			}
		}
	} else {
		// 如果没有参数，创建零值
		if argsType.Kind() == reflect.Ptr {
			argsValue = reflect.New(argsType.Elem())
		} else {
			argsValue = reflect.Zero(argsType)
		}
	}
	
	// 调用方法
	results := method.Func.Call([]reflect.Value{
		serviceValue,
		reflect.ValueOf(ctx),
		argsValue,
	})
	
	// 处理返回值
	result := results[0].Interface()
	errValue := results[1]
	
	if !errValue.IsNil() {
		err := errValue.Interface().(error)
		return &RPCResponse{
			Error: err.Error(),
		}
	}
	
	return &RPCResponse{
		Result: result,
	}
}

func (s *DefaultRPCServer) Stop() error {
	// 1. 从注册中心注销服务
	if s.registry != nil {
		s.mu.RLock()
		for serviceName, serviceID := range s.registeredServices {
			ctx := context.Background()
			if err := s.registry.Unregister(ctx, serviceName, serviceID); err != nil {
				log.Printf("Failed to unregister service %s from registry: %v", serviceName, err)
			} else {
				log.Printf("Auto-unregistered service %s from registry", serviceName)
			}
		}
		s.mu.RUnlock()
		
		// 清空已注册服务记录
		s.mu.Lock()
		s.registeredServices = make(map[string]string)
		s.mu.Unlock()
	}
	
	// 2. 关闭监听器
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

// ============================================================================
// 6. 中间件系统框架
// ============================================================================

// Middleware 中间件接口
type Middleware interface {
	// TODO: 实现中间件处理逻辑
	Process(ctx context.Context, req *RPCRequest, next MiddlewareFunc) (*RPCResponse, error)
}

// MiddlewareFunc 中间件函数类型
type MiddlewareFunc func(ctx context.Context, req *RPCRequest) (*RPCResponse, error)

// RPCRequest RPC请求
type RPCRequest struct {
	ServiceName string      `json:"service_name"`
	Method      string      `json:"method"`
	Args        interface{} `json:"args"`
	Headers     map[string]string `json:"headers"`
}

// RPCResponse RPC响应
type RPCResponse struct {
	Result interface{} `json:"result"`
	Error  string      `json:"error"`
	Headers map[string]string `json:"headers"`
}

// LoggingMiddleware 日志中间件
type LoggingMiddleware struct{}

func (m *LoggingMiddleware) Process(ctx context.Context, req *RPCRequest, next MiddlewareFunc) (*RPCResponse, error) {
	// TODO: 实现日志记录逻辑
	// 1. 记录请求开始
	startTime := time.Now()

	// 2. 调用下一个中间件
	resp, err := next(ctx, req)
	// 3. 记录请求结束和耗时
	log.Printf("TODO: Log request %s.%s, duration: %v, error: %v", req.ServiceName, req.Method, time.Since(startTime), err)
	return resp, err
}

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	tokenValidator func(token string) bool
}

func NewAuthMiddleware(validator func(string) bool) *AuthMiddleware {
	return &AuthMiddleware{tokenValidator: validator}
}

func (m *AuthMiddleware) Process(ctx context.Context, req *RPCRequest, next MiddlewareFunc) (*RPCResponse, error) {
	// TODO: 实现认证逻辑
	// 1. 从请求头获取token
	token, ok := req.Headers["Authorization"]
	if !ok || !m.tokenValidator(token) {
		return nil, fmt.Errorf("invalid or missing token")
	}
	// 2. 验证token

	// 3. 设置用户信息到context
	log.Printf("TODO: Authenticate request %s.%s", req.ServiceName, req.Method)
	return next(ctx, req)
}

// RateLimitMiddleware 限流中间件
type RateLimitMiddleware struct {
	limiter map[string]*TokenBucket // clientID -> bucket
	mu      sync.RWMutex
}

// TokenBucket 令牌桶
type TokenBucket struct {
	capacity int
	tokens   int
	lastRefill time.Time
	mu       sync.Mutex
}

func NewRateLimitMiddleware() *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limiter: make(map[string]*TokenBucket),
	}
}

func (m *RateLimitMiddleware) Process(ctx context.Context, req *RPCRequest, next MiddlewareFunc) (*RPCResponse, error) {
	// TODO: 实现限流逻辑
	// 1. 获取客户端ID
	clientID := req.Headers["Client-ID"]
	if clientID == "" {
		return nil, fmt.Errorf("missing Client-ID header")
	}
	// 2. 检查令牌桶
	m.mu.RLock()
	bucket, ok := m.limiter[clientID]
	m.mu.RUnlock()
	if !ok {
		// 初始化令牌桶
		bucket = &TokenBucket{
			capacity: 100, // 100个令牌
			tokens:   100, // 初始填充100个令牌
			lastRefill: time.Now(),
		}
		m.mu.Lock()
		m.limiter[clientID] = bucket
		m.mu.Unlock()
	}
	// 3. 消费令牌或拒绝请求
	log.Printf("TODO: Rate limit check for %s.%s", req.ServiceName, req.Method)
	return next(ctx, req)
}

// ============================================================================
// 7. 健康检查和故障转移框架
// ============================================================================

// HealthChecker 健康检查器接口
type HealthChecker interface {
	// TODO: 实现健康检查逻辑
	Check(ctx context.Context, service *ServiceInfo) HealthStatus
}

// HTTPHealthChecker HTTP健康检查器
type HTTPHealthChecker struct {
	timeout time.Duration
}

func NewHTTPHealthChecker(timeout time.Duration) *HTTPHealthChecker {
	return &HTTPHealthChecker{timeout: timeout}
}

func (h *HTTPHealthChecker) Check(ctx context.Context, service *ServiceInfo) HealthStatus {
	// TODO: 实现HTTP健康检查逻辑
	// 1. 发送HTTP请求到健康检查端点
	url := fmt.Sprintf("http://%s:%d/health", service.Address, service.Port)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return HealthStatus{
			Status:    "unknown",
			LastCheck: time.Now(),
			Message:   fmt.Sprintf("create request failed: %v", err),
		}
	}
	// 2. 检查响应状态
	client := &http.Client{Timeout: h.timeout}
	resp, err := client.Do(req)
	if err != nil {
		return HealthStatus{
			Status:    "unknown",
			LastCheck: time.Now(),
			Message:   fmt.Sprintf("request failed: %v", err),
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return HealthStatus{
			Status:    "unhealthy",
			LastCheck: time.Now(),
			Message:   fmt.Sprintf("unexpected status code: %d", resp.StatusCode),
		}
	}
	// 3. 返回健康状态
	log.Printf("TODO: HTTP health check for %s at %s:%d", service.Name, service.Address, service.Port)
	return HealthStatus{
		Status:    "healthy",
		LastCheck: time.Now(),
		Message:   "OK",
	}
}

// TCPHealthChecker TCP健康检查器
type TCPHealthChecker struct {
	timeout time.Duration
}

func NewTCPHealthChecker(timeout time.Duration) *TCPHealthChecker {
	return &TCPHealthChecker{timeout: timeout}
}

func (h *TCPHealthChecker) Check(ctx context.Context, service *ServiceInfo) HealthStatus {
	// TODO: 实现TCP健康检查逻辑
	// 1. 尝试建立TCP连接
	addr := fmt.Sprintf("%s:%d", service.Address, service.Port)
	conn, err := net.DialTimeout("tcp", addr, h.timeout)
	if err != nil {
		return HealthStatus{
			Status:    "unhealthy",
			LastCheck: time.Now(),
			Message:   fmt.Sprintf("connect failed: %v", err),
		}
	}
	defer conn.Close()
	// 2. 检查连接状态
	if err := conn.SetReadDeadline(time.Now().Add(h.timeout)); err != nil {
		return HealthStatus{
			Status:    "unhealthy",
			LastCheck: time.Now(),
			Message:   fmt.Sprintf("set read deadline failed: %v", err),
		}
	}
	// 3. 返回健康状态
	log.Printf("TODO: TCP health check for %s at %s:%d", service.Name, service.Address, service.Port)
	return HealthStatus{
		Status:    "healthy",
		LastCheck: time.Now(),
		Message:   "OK",
	}
}

// FailoverManager 故障转移管理器
type FailoverManager struct {
	registry      ServiceRegistry
	healthChecker HealthChecker
	checkInterval time.Duration
	mu            sync.RWMutex
}

func NewFailoverManager(registry ServiceRegistry, checker HealthChecker, interval time.Duration) *FailoverManager {
	return &FailoverManager{
		registry:      registry,
		healthChecker: checker,
		checkInterval: interval,
	}
}

func (f *FailoverManager) Start(ctx context.Context) {
	// TODO: 实现故障转移逻辑
	// 1. 定期检查所有服务健康状态
	ticker := time.NewTicker(f.checkInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("Failover manager stopped")
			return
		case <-ticker.C:
			// f.healthChecker.Check(ctx, service)
		}
	}
}

// ============================================================================
// 8. 示例服务接口
// ============================================================================

// UserService 用户服务接口示例
type UserService interface {
	GetUser(ctx context.Context, userID string) (*User, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	DeleteUser(ctx context.Context, userID string) (*User, error)
}

// User 用户模型
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserServiceImpl 用户服务实现示例
type UserServiceImpl struct {
	users map[string]*User
	mu    sync.RWMutex
}

func NewUserService() *UserServiceImpl {
	return &UserServiceImpl{
		users: make(map[string]*User),
	}
}

func (s *UserServiceImpl) GetUser(ctx context.Context, userID string) (*User, error) {
	// TODO: 实现获取用户逻辑
	log.Printf("TODO: Get user %s", userID)
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[userID]
	if !ok {
		return nil, fmt.Errorf("user %s not found", userID)
	}
	return user, nil
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, user *User) (*User, error) {
	// TODO: 实现创建用户逻辑
	log.Printf("TODO: Create user %v", user)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[user.ID] = user
	return user, nil
}

func (s *UserServiceImpl) UpdateUser(ctx context.Context, user *User) (*User, error) {
	// TODO: 实现更新用户逻辑
	log.Printf("TODO: Update user %v", user)
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[user.ID]; !ok {
		return nil, fmt.Errorf("user %s not found", user.ID)
	}
	s.users[user.ID] = user
	return user, nil
}

func (s *UserServiceImpl) DeleteUser(ctx context.Context, userID string) (*User, error) {
	// TODO: 实现删除用户逻辑
	log.Printf("TODO: Delete user %s", userID)
	s.mu.Lock()
	defer s.mu.Unlock()
	user, ok := s.users[userID]
	if !ok {
		return nil, fmt.Errorf("user %s not found", userID)
	}
	delete(s.users, userID)
	return user, nil
}

// ============================================================================
// 9. 主函数和测试框架
// ============================================================================

func main10() {
	fmt.Println("=== 微服务框架演示 ===")
	
	// 创建服务注册中心
	registry := NewInMemoryRegistry()
	
	// 测试Register方法
	fmt.Println("\n--- 测试服务注册功能 ---")
	
	// 测试1: 正常注册服务
	service1 := &ServiceInfo{
		Name:    "UserService",
		Version: "1.0.0",
		Address: "127.0.0.1",
		Port:    8001,
		Metadata: map[string]string{
			"region": "us-west",
		},
	}
	
	ctx := context.Background()
	err := registry.Register(ctx, service1)
	if err != nil {
		log.Printf("注册服务失败: %v", err)
	} else {
		fmt.Printf("✓ 成功注册服务: %s at %s:%d\n", service1.Name, service1.Address, service1.Port)
	}
	
	// 测试2: 注册相同地址的服务（应该更新）
	service2 := &ServiceInfo{
		Name:    "UserService",
		Version: "1.1.0",
		Address: "127.0.0.1",
		Port:    8001,
		Metadata: map[string]string{
			"region": "us-west",
			"updated": "true",
		},
	}
	
	err = registry.Register(ctx, service2)
	if err != nil {
		log.Printf("更新服务失败: %v", err)
	} else {
		fmt.Printf("✓ 成功更新服务: %s at %s:%d (版本: %s)\n", service2.Name, service2.Address, service2.Port, service2.Version)
	}
	
	// 测试3: 注册不同端口的服务实例
	service3 := &ServiceInfo{
		Name:    "UserService",
		Address: "127.0.0.1",
		Port:    8002,
	}
	
	err = registry.Register(ctx, service3)
	if err != nil {
		log.Printf("注册第二个实例失败: %v", err)
	} else {
		fmt.Printf("✓ 成功注册第二个服务实例: %s at %s:%d\n", service3.Name, service3.Address, service3.Port)
	}
	
	// 测试4: 验证错误处理
	fmt.Println("\n--- 测试错误处理 ---")
	
	// 测试空服务
	err = registry.Register(ctx, nil)
	if err != nil {
		fmt.Printf("✓ 正确处理空服务错误: %v\n", err)
	}
	
	// 测试无效端口
	invalidService := &ServiceInfo{
		Name:    "InvalidService",
		Address: "127.0.0.1",
		Port:    -1,
	}
	err = registry.Register(ctx, invalidService)
	if err != nil {
		fmt.Printf("✓ 正确处理无效端口错误: %v\n", err)
	}
	
	// 测试空服务名
	emptyNameService := &ServiceInfo{
		Name:    "",
		Address: "127.0.0.1",
		Port:    8003,
	}
	err = registry.Register(ctx, emptyNameService)
	if err != nil {
		fmt.Printf("✓ 正确处理空服务名错误: %v\n", err)
	}
	
	// 显示当前注册的服务
	fmt.Println("\n--- 当前注册的服务 ---")
	registry.mu.RLock()
	for serviceName, services := range registry.services {
		fmt.Printf("服务名: %s\n", serviceName)
		for i, svc := range services {
			fmt.Printf("  实例 %d: %s:%d (版本: %s, 状态: %s, ID: %s)\n", 
				i+1, svc.Address, svc.Port, svc.Version, svc.Health.Status, svc.Metadata["service_id"])
		}
	}
	registry.mu.RUnlock()
	
	// 测试Unregister方法
	fmt.Println("\n--- 测试服务注销功能 ---")
	
	// 获取第一个服务实例的ID进行注销测试
	registry.mu.RLock()
	var testServiceID string
	if services, exists := registry.services["UserService"]; exists && len(services) > 0 {
		testServiceID = services[0].Metadata["service_id"]
	}
	registry.mu.RUnlock()
	
	if testServiceID != "" {
		// 测试正常注销
		err = registry.Unregister(ctx, "UserService", testServiceID)
		if err != nil {
			log.Printf("注销服务失败: %v", err)
		} else {
			fmt.Printf("✓ 成功注销服务实例: %s\n", testServiceID)
		}
		
		// 测试注销不存在的服务ID
		err = registry.Unregister(ctx, "UserService", "non-existent-id")
		if err != nil {
			fmt.Printf("✓ 正确处理不存在的服务ID错误: %v\n", err)
		}
		
		// 测试注销不存在的服务名
		err = registry.Unregister(ctx, "NonExistentService", "some-id")
		if err != nil {
			fmt.Printf("✓ 正确处理不存在的服务名错误: %v\n", err)
		}
		
		// 测试空参数
		err = registry.Unregister(ctx, "", "some-id")
		if err != nil {
			fmt.Printf("✓ 正确处理空服务名错误: %v\n", err)
		}
		
		err = registry.Unregister(ctx, "UserService", "")
		if err != nil {
			fmt.Printf("✓ 正确处理空服务ID错误: %v\n", err)
		}
	}
	
	// 显示注销后的服务状态
	fmt.Println("\n--- 注销后的服务状态 ---")
	registry.mu.RLock()
	for serviceName, services := range registry.services {
		fmt.Printf("服务名: %s\n", serviceName)
		for i, svc := range services {
			fmt.Printf("  实例 %d: %s:%d (版本: %s, 状态: %s, ID: %s)\n", 
				i+1, svc.Address, svc.Port, svc.Version, svc.Health.Status, svc.Metadata["service_id"])
		}
	}
	registry.mu.RUnlock()
	
	// 测试Watch方法
	fmt.Println("\n--- 测试服务监听功能 ---")
	
	// 创建带超时的context
	watchCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	// 测试正常监听
	ch, err := registry.Watch(watchCtx, "UserService")
	if err != nil {
		log.Printf("创建监听失败: %v", err)
	} else {
		fmt.Printf("✓ 成功创建UserService监听器\n")
		
		// 启动goroutine接收监听数据
		go func() {
			for services := range ch {
				fmt.Printf("📡 监听到UserService变化: %d个实例\n", len(services))
				for i, svc := range services {
					fmt.Printf("  实例 %d: %s:%d (状态: %s)\n", 
						i+1, svc.Address, svc.Port, svc.Health.Status)
				}
			}
			fmt.Println("📡 监听器已关闭")
		}()
	}
	
	// 测试错误处理
	_, err = registry.Watch(watchCtx, "")
	if err != nil {
		fmt.Printf("✓ 正确处理空服务名错误: %v\n", err)
	}
	
	_, err = registry.Watch(nil, "TestService")
	if err != nil {
		fmt.Printf("✓ 正确处理空context错误: %v\n", err)
	}
	
	// 注册新服务触发监听通知
	fmt.Println("\n--- 触发监听通知 ---")
	time.Sleep(100 * time.Millisecond) // 确保监听器已就绪
	
	newService := &ServiceInfo{
		Name:    "UserService",
		Version: "2.0.0",
		Address: "127.0.0.1",
		Port:    8003,
		Metadata: make(map[string]string),
		Health: HealthStatus{
			Status:    "healthy",
			LastCheck: time.Now(),
			Message:   "OK",
		},
	}
	
	err = registry.Register(ctx, newService)
	if err != nil {
		log.Printf("注册新服务失败: %v", err)
	} else {
		fmt.Printf("✓ 注册新服务实例: UserService at 127.0.0.1:8003\n")
	}
	
	// 等待监听通知
	time.Sleep(200 * time.Millisecond)
	
	// 注销服务触发监听通知
	registry.mu.RLock()
	var testServiceID2 string
	if services, exists := registry.services["UserService"]; exists && len(services) > 1 {
		testServiceID2 = services[1].Metadata["service_id"]
	}
	registry.mu.RUnlock()
	
	if testServiceID2 != "" {
		err = registry.Unregister(ctx, "UserService", testServiceID2)
		if err != nil {
			log.Printf("注销服务失败: %v", err)
		} else {
			fmt.Printf("✓ 注销服务实例: %s\n", testServiceID2)
		}
	}
	
	// 等待监听通知
	time.Sleep(200 * time.Millisecond)
	
	// 测试监听不存在的服务
	ch2, err := registry.Watch(watchCtx, "NonExistentService")
	if err != nil {
		log.Printf("监听不存在服务失败: %v", err)
	} else {
		fmt.Printf("✓ 成功创建不存在服务的监听器\n")
		
		// 启动goroutine接收数据
		go func() {
			select {
			case services := <-ch2:
				fmt.Printf("📡 监听到NonExistentService初始数据: %d个实例\n", len(services))
			case <-time.After(500 * time.Millisecond):
				fmt.Println("📡 NonExistentService监听器未收到初始数据（正常）")
			}
		}()
	}
	
	// 等待context超时，测试自动清理
	fmt.Println("\n--- 等待context超时测试自动清理 ---")
	time.Sleep(1 * time.Second)
	
	fmt.Println("\n--- Register方法测试完成 ---")
	
	// ============================================================================
	// 测试RPC服务端功能
	// ============================================================================
	fmt.Println("\n--- 测试RPC服务端功能 ---")
	
	// 清理注册中心中的旧服务实例，避免客户端连接到不存在的服务
	fmt.Println("\n--- 清理注册中心中的旧服务实例 ---")
	registry.mu.Lock()
	// 清空所有服务实例
	registry.services = make(map[string][]*ServiceInfo)
	registry.mu.Unlock()
	fmt.Printf("✓ 已清理注册中心中的所有旧服务实例\n")
	
	// 创建RPC服务端（使用带注册中心的构造函数）
	rpcServer := NewRPCServerWithRegistry(registry)
	
	// 创建用户服务实例
	userService := NewUserService()
	
	// 注册服务到RPC服务端
	err = rpcServer.RegisterService("UserService", userService)
	if err != nil {
		log.Printf("注册RPC服务失败: %v", err)
	} else {
		fmt.Printf("✓ 成功注册UserService到RPC服务端\n")
	}
	
	// 启动RPC服务端
	serverAddr := "127.0.0.1:9001"
	err = rpcServer.Start(serverAddr)
	if err != nil {
		log.Printf("启动RPC服务端失败: %v", err)
	} else {
		fmt.Printf("✓ RPC服务端已启动在 %s\n", serverAddr)
	}
	
	// 等待服务端启动
	time.Sleep(500 * time.Millisecond)
	
	// 测试RPC客户端调用
	fmt.Println("\n--- 测试RPC客户端调用 ---")
	
	// 创建负载均衡器
	balancer := NewRoundRobinBalancer()
	
	// 创建RPC客户端
	rpcClient := NewRPCClient(registry, balancer)

	// 注意：现在不需要手动注册RPC服务到注册中心，
	// 因为RPC服务端在启动时会自动注册所有已注册的服务

	// 等待服务端启动并完成自动注册
	time.Sleep(500 * time.Millisecond)
	
	// 测试创建用户
	fmt.Println("\n--- 测试RPC调用: CreateUser ---")
	testUser := &User{
		ID:    "user001",
		Name:  "张三",
		Email: "zhangsan@example.com",
	}
	
	var createResult interface{}
	err = rpcClient.Call(ctx, "UserService", "CreateUser", testUser, &createResult)
	if err != nil {
		fmt.Printf("✗ CreateUser调用失败: %v\n", err)
	} else {
		fmt.Printf("✓ CreateUser调用成功\n")
	}
	
	// 测试获取用户
	fmt.Println("\n--- 测试RPC调用: GetUser ---")
	var getResult *User
	err = rpcClient.Call(ctx, "UserService", "GetUser", "user001", &getResult)
	if err != nil {
		fmt.Printf("✗ GetUser调用失败: %v\n", err)
	} else {
		fmt.Printf("✓ GetUser调用成功: %+v\n", getResult)
	}
	
	// 测试更新用户
	fmt.Println("\n--- 测试RPC调用: UpdateUser ---")
	testUser.Name = "张三（已更新）"
	testUser.Email = "zhangsan.updated@example.com"
	
	var updateResult interface{}
	err = rpcClient.Call(ctx, "UserService", "UpdateUser", testUser, &updateResult)
	if err != nil {
		fmt.Printf("✗ UpdateUser调用失败: %v\n", err)
	} else {
		fmt.Printf("✓ UpdateUser调用成功\n")
	}
	
	// 再次获取用户验证更新
	fmt.Println("\n--- 验证用户更新 ---")
	var updatedUser *User
	err = rpcClient.Call(ctx, "UserService", "GetUser", "user001", &updatedUser)
	if err != nil {
		fmt.Printf("✗ 验证更新失败: %v\n", err)
	} else {
		fmt.Printf("✓ 用户更新验证成功: %+v\n", updatedUser)
	}
	
	// 测试删除用户
	fmt.Println("\n--- 测试RPC调用: DeleteUser ---")
	var deleteResult interface{}
	err = rpcClient.Call(ctx, "UserService", "DeleteUser", "user001", &deleteResult)
	if err != nil {
		fmt.Printf("✗ DeleteUser调用失败: %v\n", err)
	} else {
		fmt.Printf("✓ DeleteUser调用成功\n")
	}
	
	// 验证用户已删除
	fmt.Println("\n--- 验证用户删除 ---")
	var deletedUser *User
	err = rpcClient.Call(ctx, "UserService", "GetUser", "user001", &deletedUser)
	if err != nil {
		fmt.Printf("✓ 用户删除验证成功: %v\n", err)
	} else {
		fmt.Printf("✗ 用户删除验证失败，用户仍存在: %+v\n", deletedUser)
	}
	
	// 测试错误处理
	fmt.Println("\n--- 测试RPC错误处理 ---")
	
	// 调用不存在的服务
	var errorResult interface{}
	err = rpcClient.Call(ctx, "NonExistentService", "SomeMethod", nil, &errorResult)
	if err != nil {
		fmt.Printf("✓ 正确处理不存在服务错误: %v\n", err)
	}
	
	// 调用不存在的方法
	err = rpcClient.Call(ctx, "UserService", "NonExistentMethod", nil, &errorResult)
	if err != nil {
		fmt.Printf("✓ 正确处理不存在方法错误: %v\n", err)
	}
	
	// 关闭RPC客户端
	err = rpcClient.Close()
	if err != nil {
		log.Printf("关闭RPC客户端失败: %v", err)
	} else {
		fmt.Printf("✓ RPC客户端已关闭\n")
	}
	
	// 停止RPC服务端
	err = rpcServer.Stop()
	if err != nil {
		log.Printf("停止RPC服务端失败: %v", err)
	} else {
		fmt.Printf("✓ RPC服务端已停止\n")
	}
	
	fmt.Println("\n--- RPC服务端测试完成 ---")
	
	// 测试健康检查机制集成
	fmt.Println("\n=== 健康检查机制集成测试 ===")
	testHealthCheckIntegration()
	
	fmt.Println("演示完成，程序将在3秒后退出...")
	
	// 等待3秒后退出
	time.Sleep(3 * time.Second)
	fmt.Println("程序退出")
}

// testHealthCheckIntegration 测试健康检查机制集成
func testHealthCheckIntegration() {
	// 创建注册中心
	registry := NewInMemoryRegistry()
	fmt.Println("✓ 创建内存注册中心")
	
	// 创建带注册中心的RPC服务端
	rpcServer := NewRPCServerWithRegistry(registry)
	fmt.Println("✓ 创建RPC服务端（已集成健康检查）")
	
	// 验证默认健康检查器
	if rpcServer.healthChecker != nil {
		fmt.Println("✓ 默认健康检查器已设置（TCP健康检查）")
	} else {
		fmt.Println("✗ 默认健康检查器未设置")
	}
	
	// 创建用户服务
	userService := NewUserService()
	
	// 注册服务
	err := rpcServer.RegisterService("TestUserService", userService)
	if err != nil {
		log.Printf("注册服务失败: %v", err)
		return
	}
	fmt.Println("✓ 成功注册TestUserService")
	
	// 启动服务端（这会触发自动注册和健康检查）
	go func() {
		err := rpcServer.Start("127.0.0.1:9002")
		if err != nil {
			log.Printf("启动RPC服务端失败: %v", err)
		}
	}()
	
	// 等待服务启动
	time.Sleep(2 * time.Second)
	fmt.Println("✓ RPC服务端已启动在 127.0.0.1:9002")
	
	// 检查服务是否已注册到注册中心
	services, err := registry.Discover(context.Background(), "TestUserService")
	if err != nil {
		log.Printf("发现服务失败: %v", err)
		return
	}
	
	if len(services) > 0 {
		fmt.Printf("✓ 发现 %d 个TestUserService实例\n", len(services))
		for i, service := range services {
			fmt.Printf("  实例 %d: %s:%d, 健康状态: %s, 消息: %s\n", 
				i+1, service.Address, service.Port, service.Health.Status, service.Health.Message)
		}
	} else {
		fmt.Println("✗ 未发现TestUserService实例")
	}
	
	// 测试自定义健康检查器
	fmt.Println("\n--- 测试自定义健康检查器 ---")
	httpChecker := NewHTTPHealthChecker(3 * time.Second)
	rpcServer.SetHealthChecker(httpChecker)
	fmt.Println("✓ 设置HTTP健康检查器")
	
	// 停止服务端
	err = rpcServer.Stop()
	if err != nil {
		log.Printf("停止RPC服务端失败: %v", err)
	} else {
		fmt.Println("✓ RPC服务端已停止")
	}
	
	fmt.Println("✓ 健康检查机制集成测试完成")
}