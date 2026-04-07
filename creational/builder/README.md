# 建造者模式 (Builder Pattern)

## 模式定义

建造者模式将一个复杂对象的构建与其表示分离，使得同样的构建过程可以创建不同的表示。

> **核心思想**：将复杂对象的构造步骤分解，让构造过程更清晰、更可控。

***

## 适用场景

### 应该使用建造者模式的场景

1. **对象有很多可选参数**
   ```go
   // 不好的做法：构造函数参数太多
   NewServer(host, port, timeout, maxConn, tlsCert, tlsKey, debug, compress, ...)

   // 好的做法：使用建造者
   NewServerBuilder().WithHost(host).WithPort(port).WithDebug().Build()
   ```
2. **需要分步骤构建对象**
   - 某些参数必须在其他参数之前设置
   - 需要验证步骤间的依赖关系
3. **需要创建不可变对象**
   - 对象创建后不应该被修改
   - 所有字段在构建时确定
4. **需要不同的表示**
   - 同样的构建过程，不同的配置组合
   - 如：开发环境配置 vs 生产环境配置

### 具体示例

- HTTP 服务器配置（大量可选参数）
- SQL 查询构建器
- 测试数据构建器
- 配置文件解析器

***

## Go 语言实现

Go 语言中有三种主要的建造者模式实现：

### 1. 传统建造者模式（Traditional Builder）

```go
type ServerBuilder struct {
    server *Server
}

func NewServerBuilder() *ServerBuilder {
    return &ServerBuilder{server: &Server{/* 默认值 */}}
}

func (b *ServerBuilder) WithHost(host string) *ServerBuilder {
    b.server.Host = host
    return b  // 返回自身，支持链式调用
}

func (b *ServerBuilder) Build() (*Server, error) {
    // 验证并返回
    return b.server, nil
}
```

**使用方式**：

```go
server, err := NewServerBuilder().
    WithHost("localhost").
    WithPort(8080).
    WithTimeout(30 * time.Second).
    Build()
```

**优缺点**：

- ✅ 清晰的链式 API
- ✅ 可以在 Build 时进行验证
- ❌ 需要额外的 Builder 结构体
- ❌ 代码量较多

### 2. 函数式选项模式（Functional Options）- **Go 推荐**

这是 Go 社区广泛采用的模式，由 Rob Pike 推广，被 `grpc`、`zap`、`go-redis` 等库使用。

```go
// 定义选项类型
type Option func(*Server)

// 构造函数接收变长选项参数
func NewServer(host string, port int, opts ...Option) (*Server, error) {
    s := &Server{
        Host: host,
        Port: port,
        // 设置默认值
    }
    
    // 应用所有选项
    for _, opt := range opts {
        opt(s)
    }
    
    return s, nil
}

// 定义具体选项
func WithTimeout(timeout time.Duration) Option {
    return func(s *Server) {
        s.Timeout = timeout
    }
}
```

**使用方式**：

```go
server, err := NewServer("localhost", 8080,
    WithTimeout(30*time.Second),
    WithMaxConnections(100),
    WithDebug(),
)
```

**为什么这是 Go 的最佳实践？**

| 特性      | 传统建造者        | 函数式选项        |
| ------- | ------------ | ------------ |
| API 简洁性 | 较好           | **优秀**       |
| 扩展性     | 需要修改 Builder | **新增函数即可**   |
| 代码量     | 较多           | **较少**       |
| 选项组合    | 复杂           | **简单（切片展开）** |
| 必填参数处理  | 清晰           | **通过构造函数参数** |

**函数式选项的优势**：

1. **简洁的 API**
   ```go
   // 必填参数放前面，清晰明确
   NewServer(host, port, opts...)
   ```
2. **易于扩展**
   ```go
   // 新增选项只需添加函数，不修改现有代码
   func WithNewFeature(value string) Option { ... }
   ```
3. **选项组合**
   ```go
   // 预设配置
   prodOpts := []Option{WithTimeout(60*time.Second), WithMaxConnections(1000)}

   // 组合使用
   server, _ := NewServer("localhost", 8080, append(prodOpts, WithDebug())...)
   ```
4. **符合 Go 风格**
   - 利用变长参数
   - 利用闭包
   - 函数是一等公民

### 3. 分步建造者（Step Builder）

强制按照特定顺序设置参数，编译时检查。

```go
// 第一步：必须设置 Host
type StepBuilderHost struct{}
func (b *StepBuilderHost) WithHost(host string) *StepBuilderPort

// 第二步：必须设置 Port
type StepBuilderPort struct{}
func (b *StepBuilderPort) WithPort(port int) *StepBuilderFinal

// 最后：可选参数和 Build
type StepBuilderFinal struct{}
func (b *StepBuilderFinal) WithTimeout(d time.Duration) *StepBuilderFinal
func (b *StepBuilderFinal) Build() (*Server, error)
```

**使用方式**：

```go
server, err := NewStepBuilder().
    WithHost("localhost").     // 返回 StepBuilderPort
    WithPort(8080).            // 返回 StepBuilderFinal
    WithTimeout(30*time.Second).
    Build()
```

**优缺点**：

- ✅ 编译时强制顺序
- ✅ 不可能遗漏必填参数
- ❌ 代码复杂
- ❌ 可读性稍差

***

## Go 语言特殊点

### 1. 变长参数 + 切片展开

```go
// 定义选项切片
opts := []Option{
    WithTimeout(30*time.Second),
    WithDebug(),
}

// 使用 ... 展开
server, _ := NewServer("localhost", 8080, opts...)
```

这是 Go 独特的语法，让选项组合非常灵活。

### 2. 闭包捕获配置

```go
func WithTLS(certFile, keyFile string) Option {
    // 闭包捕获 certFile 和 keyFile
    return func(s *Server) {
        s.TLSConfig = &TLSConfig{
            CertFile: certFile,
            KeyFile:  keyFile,
        }
    }
}
```

闭包让选项函数可以携带上下文。

### 3. 零值和默认值

```go
func NewServer(host string, port int, opts ...Option) (*Server, error) {
    s := &Server{
        // 合理的默认值
        Timeout:    30 * time.Second,
        MaxConn:    100,
        Compress:   true,
    }
    // ...
}
```

Go 的零值特性让默认值处理很自然。

### 4. 命名返回值用于验证

```go
func (s *Server) Validate() error {
    if s.Host == "" {
        return fmt.Errorf("host is required")
    }
    // ...
}
```

***

## 优缺点分析

### 优点

1. **清晰性**
   - 参数名即文档
   - 链式调用可读性好
2. **灵活性**
   - 可选参数自由组合
   - 易于添加新参数
3. **不可变性**
   - 对象创建后不可修改
   - 线程安全
4. **验证集中**
   - 在 Build 时统一验证
   - 可以验证参数间的关系

### 缺点

1. **代码量增加**
   - 需要额外的 Builder/Option 代码
   - 对于简单对象是过度设计
2. **性能开销**
   - 函数调用开销（极小）
   - Builder 对象分配（可忽略）
3. **学习成本**
   - 团队需要理解模式
   - 函数式选项对新手不太直观

***

## 最佳实践

### ✅ 推荐做法

1. **优先使用函数式选项**
   ```go
   // 这是 Go 社区的首选
   func NewServer(host string, port int, opts ...Option)
   ```
2. **合理的默认值**
   ```go
   func NewServer(host string, port int, opts ...Option) (*Server, error) {
       s := &Server{
           Timeout:  30 * time.Second,  // 合理的默认值
           MaxConn:  100,
       }
       // ...
   }
   ```
3. **选项组合**
   ```go
   // 预设配置
   func ProductionOptions() []Option { ... }
   func DevelopmentOptions() []Option { ... }

   // 使用
   server, _ := NewServer("localhost", 8080, ProductionOptions()...)
   ```
4. **验证逻辑**
   ```go
   func (s *Server) Validate() error {
       if s.Host == "" {
           return fmt.Errorf("host is required")
       }
       // ...
   }
   ```

### ❌ 避免的做法

1. **不要为简单对象使用建造者**
   ```go
   // 过度设计
   user := NewUserBuilder().WithName("John").WithAge(30).Build()

   // 直接创建更清晰
   user := &User{Name: "John", Age: 30}
   ```
2. **避免选项冲突**
   ```go
   // 不好：WithReadTimeout 和 WithTimeout 可能冲突
   NewServer("localhost", 8080,
       WithReadTimeout(10*time.Second),
       WithTimeout(30*time.Second),  // 覆盖了上面的设置
   )
   ```
3. **不要忽视错误处理**
   ```go
   // 不好
   server := NewServerBuilder().WithHost("").Build()  // 应该返回错误

   // 好
   server, err := NewServerBuilder().WithHost("").Build()
   if err != nil { ... }
   ```

***

## 实际案例

### uber-go/zap

```go
logger, err := zap.NewProduction(
    zap.Fields(zap.String("service", "my-service")),
    zap.AddCaller(),
)
```

### grpc

```go
conn, err := grpc.Dial("localhost:50051",
    grpc.WithInsecure(),
    grpc.WithBlock(),
    grpc.WithTimeout(10*time.Second),
)
```

### go-redis

```go
client := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
})
```

***

## 与其他模式的关系

| 模式       | 关系   | 说明                   |
| -------- | ---- | -------------------- |
| **工厂**   | 配合使用 | 先用工厂选择类型，再用建造者配置     |
| **单例**   | 不同目的 | 单例控制实例数量，建造者控制实例配置   |
| **模板方法** | 相似   | 都定义了算法骨架，但建造者专注于对象构造 |

***

## 总结

Go 语言中实现建造者模式有三种方式：

| 方式        | 适用场景     | 推荐度   |
| --------- | -------- | ----- |
| **函数式选项** | 大多数场景    | ⭐⭐⭐⭐⭐ |
| **传统建造者** | 需要复杂验证逻辑 | ⭐⭐⭐   |
| **分步建造者** | 必须强制设置顺序 | ⭐⭐    |

> 💡 **核心建议**：在 Go 中，**函数式选项模式**是首选。它简洁、灵活、符合 Go 风格，被广泛应用于标准库和主流第三方库中。只有在确实需要强制顺序或复杂验证时，才考虑使用传统建造者。

