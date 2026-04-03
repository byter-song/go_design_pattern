# 代理模式 (Proxy Pattern)

## 概述

代理模式为其他对象提供一种代理以控制对这个对象的访问。在某些情况下，一个对象不适合或者不能直接引用另一个对象，而代理对象可以在客户端和目标对象之间起到中介的作用。

## Go 语言实现特点

Go 语言通过接口实现代理模式，具有类型安全和运行时可替换的优势。

### 代理类型

1. **虚拟代理 (Virtual Proxy)**：延迟加载，按需创建昂贵对象
2. **保护代理 (Protection Proxy)**：控制访问权限
3. **缓存代理 (Caching Proxy)**：缓存结果，提高性能
4. **智能引用 (Smart Reference)**：额外操作如引用计数
5. **远程代理 (Remote Proxy)**：访问远程对象

### 核心概念

```go
// 主题接口
type Subject interface {
    Request() error
}

// 真实主题
type RealSubject struct{}
func (r *RealSubject) Request() error {
    // 实际业务逻辑
    return nil
}

// 代理
type Proxy struct {
    realSubject *RealSubject
}

func (p *Proxy) Request() error {
    // 前置处理
    // ...
    
    // 调用真实主题
    err := p.realSubject.Request()
    
    // 后置处理
    // ...
    
    return err
}
```

## 代码示例

### 虚拟代理（延迟加载）

```go
type Image interface {
    Display() error
    GetFilename() string
}

// 真实图片（加载成本高）
type RealImage struct {
    filename string
    data     []byte
}

func NewRealImage(filename string) *RealImage {
    // 模拟耗时加载
    time.Sleep(100 * time.Millisecond)
    return &RealImage{filename: filename, data: make([]byte, 1024*1024)}
}

func (r *RealImage) Display() error {
    fmt.Printf("Displaying: %s\n", r.filename)
    return nil
}

// 虚拟代理
type ProxyImage struct {
    filename  string
    realImage *RealImage
    mu        sync.Mutex
}

func (p *ProxyImage) Display() error {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    // 延迟加载：第一次访问时才创建
    if p.realImage == nil {
        p.realImage = NewRealImage(p.filename)
    }
    return p.realImage.Display()
}
```

### 保护代理（访问控制）

```go
type User struct {
    Name string
    Role string
}

func (u *User) IsAdmin() bool {
    return u.Role == "admin"
}

type ProtectedImage struct {
    image     Image
    user      *User
    adminOnly bool
}

func (p *ProtectedImage) Display() error {
    if p.adminOnly && !p.user.IsAdmin() {
        return fmt.Errorf("access denied: %s is not admin", p.user.Name)
    }
    return p.image.Display()
}
```

### 缓存代理（带 TTL）

```go
type CacheProxy struct {
    image     Image
    cache     map[string]interface{}
    cacheTime map[string]time.Time
    ttl       time.Duration
    mu        sync.RWMutex
}

func (c *CacheProxy) GetSize() int {
    c.mu.RLock()
    if val, ok := c.cache["size"]; ok {
        if time.Since(c.cacheTime["size"]) < c.ttl {
            c.mu.RUnlock()
            return val.(int)
        }
    }
    c.mu.RUnlock()
    
    // 缓存未命中，获取并缓存
    size := c.image.GetSize()
    c.mu.Lock()
    c.cache["size"] = size
    c.cacheTime["size"] = time.Now()
    c.mu.Unlock()
    
    return size
}
```

### 智能引用代理

```go
type SmartReferenceProxy struct {
    image      Image
    refCount   int
    lastAccess time.Time
    mu         sync.Mutex
}

func (s *SmartReferenceProxy) Display() error {
    s.mu.Lock()
    s.refCount++
    s.lastAccess = time.Now()
    s.mu.Unlock()
    return s.image.Display()
}

func (s *SmartReferenceProxy) GetRefCount() int {
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.refCount
}
```

### 远程代理（模拟）

```go
type RemoteImage struct {
    serverURL string
    filename  string
}

func (r *RemoteImage) Display() error {
    fmt.Printf("Fetching %s from %s...\n", r.filename, r.serverURL)
    // 模拟网络请求
    time.Sleep(50 * time.Millisecond)
    fmt.Printf("Displaying remote image: %s\n", r.filename)
    return nil
}
```

## 使用场景

1. **延迟加载**：大对象按需加载，减少启动时间
2. **访问控制**：基于角色的权限管理
3. **缓存优化**：缓存昂贵操作的结果
4. **远程访问**：访问远程服务就像访问本地对象
5. **监控统计**：记录访问次数、频率等
6. **资源管理**：引用计数、连接池等

## 优缺点

### 优点

- **职责分离**：将访问控制与业务逻辑分离
- **延迟加载**：优化资源使用
- **透明性**：客户端无需知道代理的存在
- **扩展性**：易于添加新的代理类型

### 缺点

- **响应延迟**：增加了一层间接调用
- **代码复杂度**：需要维护代理和真实对象

## 代理 vs 装饰器

| 特性 | 代理模式 | 装饰器模式 |
|------|---------|-----------|
| 目的 | 控制访问 | 添加功能 |
| 接口 | 通常相同 | 通常相同 |
| 关注点 | 访问管理 | 功能增强 |
| 创建时机 | 控制对象创建 | 包装已有对象 |

## 与其他模式的关系

- **装饰器模式**：都实现相同接口，但目的不同
- **适配器模式**：代理控制访问，适配器转换接口
- **外观模式**：代理控制单个对象访问，外观简化子系统接口

## 参考

- [Go 设计模式 - 代理模式](https://golangbyexample.com/proxy-design-pattern-go/)
- [Proxy Pattern in Go](https://refactoring.guru/design-patterns/proxy/go/example)
