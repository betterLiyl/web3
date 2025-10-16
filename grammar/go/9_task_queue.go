package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net"
    "time"
    "bytes"
    "encoding/gob"
    "encoding/binary"
    "sync"
    "context"
)
// 任务6：分布式任务队列
// **目标**：掌握网络编程、序列化、设计模式
// **描述**：实现一个简单的分布式任务队列系统，支持任务的提交、执行和结果获取

// **流程提示**：
// 1. 设计任务结构和接口
// 2. 实现任务序列化（JSON/gob）
// 3. 实现基于TCP的网络通信
// 4. 设计客户端-服务器协议
// 5. 实现任务调度和执行
// 6. 添加任务持久化和故障恢复

// 客户端连接信息
type ClientInfo struct {
    ID         string
    Conn       *net.TCPConn
    ConnectedAt time.Time
    LastPing   time.Time
    ClientType string // "producer" or "consumer"
    Topics     []string // 消费者订阅的主题
}

type Task struct {
    ID        string
    Payload   []byte
    CreatedAt time.Time
    Topic     string
    SerializedType string
}

type TaskClient struct{
    addr string   //serverAddr
    conn *net.TCPConn
    clientID string
}

func (c *TaskClient) Connect() error {
    addr, err := net.ResolveTCPAddr("tcp", c.addr)
    if err != nil {
        return err
    }
    conn, err := net.DialTCP("tcp", nil, addr)
    if err != nil {
        return err
    }
    c.conn = conn
    return nil
}

// 注册客户端
func (c *TaskClient) Register(clientType string) error {
    if c.clientID == "" {
        c.clientID = fmt.Sprintf("%s-%d", clientType, time.Now().UnixNano())
    }
    
    header := []byte{opRegister, 0} // 注册不需要序列化类型
    if _, err := c.conn.Write(header); err != nil {
        return err
    }
    
    registerData := []byte(fmt.Sprintf("%s:%s", c.clientID, clientType))
    if err := writeFrame(c.conn, registerData); err != nil {
        return err
    }
    
    // 读取注册响应
    response, err := readFrame(c.conn)
    if err != nil {
        return err
    }
    
    if string(response) != "OK" {
        return fmt.Errorf("registration failed: %s", string(response))
    }
    
    return nil
}

// 发送心跳
func (c *TaskClient) Ping() error {
    header := []byte{opPing, 0}
    if _, err := c.conn.Write(header); err != nil {
        return err
    }
    
    // 读取pong响应
    response, err := readFrame(c.conn)
    if err != nil {
        return err
    }
    
    if string(response) != "PONG" {
        return fmt.Errorf("ping failed: %s", string(response))
    }
    
    return nil
}

// 查询服务器状态
func (c *TaskClient) GetStatus() (string, error) {
    header := []byte{opStatus, 0}
    if _, err := c.conn.Write(header); err != nil {
        return "", err
    }
    
    response, err := readFrame(c.conn)
    if err != nil {
        return "", err
    }
    
    return string(response), nil
}

// 启动心跳goroutine
func (c *TaskClient) StartHeartbeat(ctx context.Context) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            if err := c.Ping(); err != nil {
                log.Printf("Heartbeat failed for client %s: %v", c.clientID, err)
                return
            }
        }
    }
}

type TaskProducer struct{
    client *TaskClient
    serializer TaskSerializer
    serializerName string
	selfAddr string //clientAddr
}
func (p *TaskProducer) Submit(task *Task) error {
    data, err := p.serializer.Serialize(task)
    if err != nil {
        return err
    }
    st := serializerNameToID(p.serializerName)
    header := []byte{opSubmit, st}
    if _, err := p.client.conn.Write(header); err != nil {
        return err
    }
    return writeFrame(p.client.conn, data)
}
type TaskConsumer struct{
    client *TaskClient
    serializer TaskSerializer
    serverAddr string
    serializerName string
	selfAddr string //clientAddr
}
func (c *TaskConsumer) Consume(topic string) (*Task, error) {
    st := serializerNameToID(c.serializerName)
    header := []byte{opConsume, st}
    if _, err := c.client.conn.Write(header); err != nil {
        return nil, err
    }
    topicBytes := []byte(topic)
    if err := writeFrame(c.client.conn, topicBytes); err != nil {
        return nil, err
    }
    resp, err := readFrame(c.client.conn)
    if err != nil {
        return nil, err
    }
    task, err := c.serializer.Deserialize(resp)
    if err != nil {
        return nil, err
    }
    return task, nil
}


type TaskSerializer interface{
    Serialize(task *Task) ([]byte, error)
    Deserialize(data []byte) (*Task, error)
}
type JSONTaskSerializer struct{
}
func (s *JSONTaskSerializer) Serialize(task *Task) ([]byte, error) {
    return json.Marshal(task)
}
func (s *JSONTaskSerializer) Deserialize(data []byte) (*Task, error) {
    task := &Task{}
    err := json.Unmarshal(data, task)
    if err != nil {
        return nil, err
    }
    return task, nil
}
type GobTaskSerializer struct{
}
func (s *GobTaskSerializer) Serialize(task *Task) ([]byte, error) {
    var buf bytes.Buffer
    enc := gob.NewEncoder(&buf)
    err := enc.Encode(task)
    if err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}
func (s *GobTaskSerializer) Deserialize(data []byte) (*Task, error) {
    task := &Task{}
    dec := gob.NewDecoder(bytes.NewBuffer(data))
    err := dec.Decode(task)
    if err != nil {
        return nil, err
    }
    return task, nil
}

type TaskServer struct{
    addr string
    queue map[string][]*Task  //topic -> task[]
    mu sync.Mutex
    cond *sync.Cond
    clients map[string]*ClientInfo // clientID -> ClientInfo
    clientsMu sync.RWMutex
}

func (t *TaskServer) Submit(task *Task) error {
    t.mu.Lock()
    defer t.mu.Unlock()
    t.queue[task.Topic] = append(t.queue[task.Topic], task)
    t.cond.Broadcast()
    return nil
}

// 注册客户端连接
func (t *TaskServer) RegisterClient(clientID string, conn *net.TCPConn, clientType string) {
    t.clientsMu.Lock()
    defer t.clientsMu.Unlock()
    
    client := &ClientInfo{
        ID:         clientID,
        Conn:       conn,
        ConnectedAt: time.Now(),
        LastPing:   time.Now(),
        ClientType: clientType,
        Topics:     make([]string, 0),
    }
    t.clients[clientID] = client
    log.Printf("Client registered: %s (%s) from %s", clientID, clientType, conn.RemoteAddr())
}

// 注销客户端连接
func (t *TaskServer) UnregisterClient(clientID string) {
    t.clientsMu.Lock()
    defer t.clientsMu.Unlock()
    
    if client, exists := t.clients[clientID]; exists {
        log.Printf("Client unregistered: %s (%s)", clientID, client.ClientType)
        delete(t.clients, clientID)
    }
}

// 更新客户端心跳
func (t *TaskServer) UpdateClientPing(clientID string) {
    t.clientsMu.Lock()
    defer t.clientsMu.Unlock()
    
    if client, exists := t.clients[clientID]; exists {
        client.LastPing = time.Now()
    }
}

// 获取所有连接的客户端
func (t *TaskServer) GetConnectedClients() map[string]*ClientInfo {
    t.clientsMu.RLock()
    defer t.clientsMu.RUnlock()
    
    result := make(map[string]*ClientInfo)
    for id, client := range t.clients {
        // 创建副本避免并发问题
        clientCopy := *client
        result[id] = &clientCopy
    }
    return result
}

// 检查并清理超时的客户端连接
func (t *TaskServer) cleanupTimeoutClients() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        t.clientsMu.Lock()
        now := time.Now()
        var toRemove []string
        
        for clientID, client := range t.clients {
            // 如果超过60秒没有心跳，认为连接超时
            if now.Sub(client.LastPing) > 60*time.Second {
                toRemove = append(toRemove, clientID)
                client.Conn.Close()
            }
        }
        
        for _, clientID := range toRemove {
            log.Printf("Client timeout removed: %s", clientID)
            delete(t.clients, clientID)
        }
        t.clientsMu.Unlock()
    }
}

func NewTaskServer(addr string) *TaskServer{
    server := &TaskServer{
        addr: addr,
        queue: make(map[string][]*Task),
        clients: make(map[string]*ClientInfo),
    }
    server.cond = sync.NewCond(&server.mu)
    
    // 启动客户端超时清理goroutine
    go server.cleanupTimeoutClients()
    
    return server
}

const (
    opSubmit  = 1
    opConsume = 2
    opRegister = 3  // 客户端注册
    opPing    = 4   // 心跳
    opStatus  = 5   // 查询服务器状态
)

const (
    stJSON = 1
    stGOB  = 2
)

func serializerNameToID(name string) byte {
    switch name {
    case "json":
        return stJSON
    case "gob":
        return stGOB
    default:
        return stJSON
    }
}

func serializerByID(id byte) TaskSerializer {
    switch id {
    case stJSON:
        return &JSONTaskSerializer{}
    case stGOB:
        return &GobTaskSerializer{}
    default:
        return &JSONTaskSerializer{}
    }
}

func writeFrame(conn *net.TCPConn, payload []byte) error {
    var hdr [4]byte
    binary.BigEndian.PutUint32(hdr[:], uint32(len(payload)))
    if _, err := conn.Write(hdr[:]); err != nil {
        return err
    }
    _, err := conn.Write(payload)
    return err
}

func readFrame(conn *net.TCPConn) ([]byte, error) {
    var hdr [4]byte
    if _, err := io.ReadFull(conn, hdr[:]); err != nil {
        return nil, err
    }
    n := binary.BigEndian.Uint32(hdr[:])
    buf := make([]byte, n)
    if _, err := io.ReadFull(conn, buf); err != nil {
        return nil, err
    }
    return buf, nil
}

func (t *TaskServer) Start() error {
    addr, err := net.ResolveTCPAddr("tcp", t.addr)
    if err != nil {
        return err
    }
    ln, err := net.ListenTCP("tcp", addr)
    if err != nil {
        return err
    }
    log.Printf("TaskServer listening on %s", t.addr)
    for {
        conn, err := ln.AcceptTCP()
        if err != nil {
            log.Printf("accept error: %v", err)
            continue
        }
        go t.handleConn(conn)
    }
}

func (t *TaskServer) handleConn(conn *net.TCPConn) {
    defer conn.Close()
    
    var clientID string
    var registered bool
    
    for {
        // 读取操作码和序列化类型
        header := make([]byte, 2)
        if _, err := io.ReadFull(conn, header); err != nil {
            if registered {
                t.UnregisterClient(clientID)
            }
            return
        }
        
        op := header[0]
        serializerType := header[1]
        
        switch op {
        case opRegister:
            // 读取客户端注册信息
            data, err := readFrame(conn)
            if err != nil {
                log.Printf("Failed to read register data: %v", err)
                return
            }
            
            // 解析注册信息 (clientID:clientType)
            parts := bytes.Split(data, []byte(":"))
            if len(parts) != 2 {
                log.Printf("Invalid register format")
                return
            }
            
            clientID = string(parts[0])
            clientType := string(parts[1])
            
            t.RegisterClient(clientID, conn, clientType)
            registered = true
            
            // 发送注册成功响应
            response := []byte("OK")
            if err := writeFrame(conn, response); err != nil {
                log.Printf("Failed to send register response: %v", err)
                return
            }
            
        case opPing:
            if registered {
                t.UpdateClientPing(clientID)
                // 发送pong响应
                response := []byte("PONG")
                if err := writeFrame(conn, response); err != nil {
                    log.Printf("Failed to send ping response: %v", err)
                    return
                }
            }
            
        case opStatus:
            if registered {
                clients := t.GetConnectedClients()
                statusInfo := fmt.Sprintf("Connected clients: %d", len(clients))
                for id, client := range clients {
                    statusInfo += fmt.Sprintf("\n- %s (%s) connected at %s", 
                        id, client.ClientType, client.ConnectedAt.Format("15:04:05"))
                }
                
                if err := writeFrame(conn, []byte(statusInfo)); err != nil {
                    log.Printf("Failed to send status response: %v", err)
                    return
                }
            }
            
        case opSubmit:
            if !registered {
                log.Printf("Client not registered for submit operation")
                return
            }
            
            data, err := readFrame(conn)
            if err != nil {
                log.Printf("Failed to read submit data: %v", err)
                return
            }
            
            serializer := serializerByID(serializerType)
            if serializer == nil {
                log.Printf("Unknown serializer type: %d", serializerType)
                return
            }
            
            task, err := serializer.Deserialize(data)
            if err != nil {
                log.Printf("Failed to deserialize task: %v", err)
                return
            }
            
            if err := t.Submit(task); err != nil {
                log.Printf("Failed to submit task: %v", err)
                return
            }
            
            log.Printf("Task submitted: %s to topic %s by client %s", task.ID, task.Topic, clientID)
            
        case opConsume:
            if !registered {
                log.Printf("Client not registered for consume operation")
                return
            }
            
            topicData, err := readFrame(conn)
            if err != nil {
                log.Printf("Failed to read topic: %v", err)
                return
            }
            
            topic := string(topicData)
            
            // 等待任务
            t.mu.Lock()
            for len(t.queue[topic]) == 0 {
                t.cond.Wait()
            }
            
            task := t.queue[topic][0]
            t.queue[topic] = t.queue[topic][1:]
            t.mu.Unlock()
            
            serializer := serializerByID(serializerType)
            if serializer == nil {
                log.Printf("Unknown serializer type: %d", serializerType)
                return
            }
            
            data, err := serializer.Serialize(task)
            if err != nil {
                log.Printf("Failed to serialize task: %v", err)
                return
            }
            
            if err := writeFrame(conn, data); err != nil {
                log.Printf("Failed to send task: %v", err)
                return
            }
            
            log.Printf("Task consumed: %s from topic %s by client %s", task.ID, topic, clientID)
            
        default:
            log.Printf("Unknown operation: %d", op)
            return
        }
    }
}

func NewTaskProducer(client *TaskClient, serializerName string) *TaskProducer {
    return &TaskProducer{
        client: client,
        serializer: func() TaskSerializer {
            switch serializerName {
            case "json":
                return &JSONTaskSerializer{}
            case "gob":
                return &GobTaskSerializer{}
            default:
                return &JSONTaskSerializer{}
            }
        }(),
        serializerName: serializerName,
    }
}

func NewTaskConsumer(client *TaskClient, serializerName string) *TaskConsumer {
    return &TaskConsumer{
        client: client,
        serializer: func() TaskSerializer {
            switch serializerName {
            case "json":
                return &JSONTaskSerializer{}
            case "gob":
                return &GobTaskSerializer{}
            default:
                return &JSONTaskSerializer{}
            }
        }(),
        serializerName: serializerName,
    }
}

func init() {
    gob.Register(Task{})
}

func main9() {
    server := NewTaskServer("127.0.0.1:9000")
    go func() {
        if err := server.Start(); err != nil {
            log.Fatalf("server error: %v", err)
        }
    }()
    time.Sleep(300 * time.Millisecond)

    // 创建生产者客户端
    prodClient := &TaskClient{addr: "127.0.0.1:9000"}
    if err := prodClient.Connect(); err != nil {
        log.Fatalf("producer connect error: %v", err)
    }
    
    // 注册生产者
    if err := prodClient.Register("producer"); err != nil {
        log.Fatalf("producer register error: %v", err)
    }
    
    // 启动生产者心跳
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    go prodClient.StartHeartbeat(ctx)
    
    producer := NewTaskProducer(prodClient, "json")
    
    // 创建消费者客户端
    consClient := &TaskClient{addr: "127.0.0.1:9000"}
    if err := consClient.Connect(); err != nil {
        log.Fatalf("consumer connect error: %v", err)
    }
    
    // 注册消费者
    if err := consClient.Register("consumer"); err != nil {
        log.Fatalf("consumer register error: %v", err)
    }
    
    // 启动消费者心跳
    go consClient.StartHeartbeat(ctx)
    
    consumer := NewTaskConsumer(consClient, "json")
    
    // 查询服务器状态
    status, err := prodClient.GetStatus()
    if err != nil {
        log.Printf("Failed to get status: %v", err)
    } else {
        fmt.Printf("Server status:\n%s\n\n", status)
    }
    
    // 提交任务
    task := &Task{
        ID: "task-001", 
        Payload: []byte("hello world"), 
        CreatedAt: time.Now(), 
        Topic: "alpha", 
        SerializedType: "json",
    }
    if err := producer.Submit(task); err != nil {
        log.Fatalf("submit error: %v", err)
    }
    
    // 消费任务
    got, err := consumer.Consume("alpha")
    if err != nil {
        log.Fatalf("consume error: %v", err)
    }
    fmt.Printf("Consumed task: id=%s topic=%s payload=%s\n", got.ID, got.Topic, string(got.Payload))
    
    // 再次查询服务器状态
    time.Sleep(100 * time.Millisecond)
    status, err = consClient.GetStatus()
    if err != nil {
        log.Printf("Failed to get status: %v", err)
    } else {
        fmt.Printf("\nFinal server status:\n%s\n", status)
    }
}
