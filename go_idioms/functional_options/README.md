# 函数式选项模式 (Functional Options Pattern)

## 概述

函数式选项模式是 Go 语言社区广泛认可的初始化复杂结构体的最佳实践，由 Rob Pike 推广，被广泛应用于标准库（如 gRPC、Zap、Viper 等）和第三方库中。

## 设计思想

1. 使用函数作为配置参数，提供类型安全、可扩展的 API
2. 避免使用大量的构造函数或复杂的配置结构体
3. 支持默认值的优雅设置
4. 允许调用者只设置关心的选项

## 核心优势

| 优势 | 说明 |
|------|------|
| **向后兼容** | 新增选项不会破坏现有代码 |
| **自文档化** | 调用代码清晰展示配置意图 |
| **类型安全** | 编译时检查，避免使用 `interface{}` |
| **可扩展** | 轻松添加新选项而不影响现有代码 |
| **默认值友好** | 未设置的选项使用合理的默认值 |

## 对比其他方案

### 方案1: 大量构造函数（不推荐 ❌）

```go
NewServer(addr string) *Server
NewServerWithPort(addr string, port int) *Server
NewServerWithPortAndTimeout(addr string, port int, timeout time.Duration) *Server
// 组合爆炸问题！
```

### 方案2: 配置结构体（不够优雅 ⚠️）

```go
config := &Config{Addr: "localhost", Port: 8080}
server := NewServer(config)
// 需要引入额外的 Config 类型，且无法强制设置必填参数
```

### 方案3: 函数式选项（推荐 ✅）

```go
server, err := NewServer("localhost",
    WithPort(8080),
    WithTimeout(30*time.Second),
    WithMaxConnections(1000),
)
// 清晰、类型安全、可扩展
```

## 实现模式

### 核心类型定义

```go
// Option 定义函数式选项类型
type Option func(*Server)

// Server HTTP 服务器结构体
type Server struct {
    host string
    port int
    // ... 其他字段
}
```

### 构造函数

```go
func NewServer(host string, opts ...Option) (*Server, error) {
    // 1. 设置默认值
    s := &Server{
        host: host,
        port: 8080,
        readTimeout: 15 * time.Second,
        // ...
    }

    // 2. 应用所有选项
    for _, opt := range opts {
        opt(s)
    }

    // 3. 验证配置
    if err := s.validate(); err != nil {
        return nil, err
    }

    return s, nil
}
```

### 选项函数

```go
// WithPort 设置端口
func WithPort(port int) Option {
    return func(s *Server) {
        s.port = port
    }
}

// WithTimeout 设置所有超时
func WithTimeout(timeout time.Duration) Option {
    return func(s *Server) {
        s.readTimeout = timeout
        s.writeTimeout = timeout
        s.idleTimeout = timeout * 4
    }
}

// WithTLS 启用 TLS
func WithTLS(certFile, keyFile string) Option {
    return func(s *Server) {
        s.enableTLS = true
        s.certFile = certFile
        s.keyFile = keyFile
    }
}
```

## 使用示例

### 最小化配置

```go
// 使用所有默认值
server, err := NewServer("localhost")
if err != nil {
    log.Fatal(err)
}
```

### 完整配置

```go
server, err := NewServer("0.0.0.0",
    WithPort(8080),
    WithReadTimeout(30*time.Second),
    WithWriteTimeout(30*time.Second),
    WithMaxConnections(1000),
    WithLogging(true),
    WithMetrics(true),
    WithCORS(true),
    WithTLS("cert.pem", "key.pem"),
)
```

### 预设配置组合

```go
// 生产环境预设
server, err := NewServer("0.0.0.0", ProductionDefaults()...)

// 开发环境预设 + 自定义
server, err := NewServer("localhost",
    append(DevelopmentDefaults(), WithPort(3000))...,
)
```

## 高级技巧

### 1. 预设配置

```go
// ProductionDefaults 生产环境推荐配置
func ProductionDefaults() []Option {
    return []Option{
        WithTimeout(30 * time.Second),
        WithMaxConnections(1000),
        WithLogging(true),
        WithMetrics(true),
    }
}

// DevelopmentDefaults 开发环境推荐配置
func DevelopmentDefaults() []Option {
    return []Option{
        WithTimeout(5 * time.Minute),
        WithCORS(true),
        WithLogging(true),
    }
}
```

### 2. 功能性选项

```go
// WithSecureHeaders 添加安全头中间件
func WithSecureHeaders() Option {
    return WithMiddleware(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("X-Content-Type-Options", "nosniff")
            w.Header().Set("X-Frame-Options", "DENY")
            next.ServeHTTP(w, r)
        })
    })
}
```

### 3. 选项验证

```go
func (s *Server) validate() error {
    if s.port < 1 || s.port > 65535 {
        return fmt.Errorf("invalid port: %d", s.port)
    }
    if s.enableTLS && (s.certFile == "" || s.keyFile == "") {
        return fmt.Errorf("TLS enabled but cert/key not provided")
    }
    return nil
}
```

## 真实案例

### gRPC 客户端配置

```go
conn, err := grpc.Dial("localhost:50051",
    grpc.WithInsecure(),
    grpc.WithBlock(),
    grpc.WithTimeout(10*time.Second),
    grpc.WithUnaryInterceptor(interceptor),
)
```

### Zap 日志配置

```go
logger, err := zap.NewProduction(
    zap.Fields(zap.String("service", "my-service")),
    zap.AddCaller(),
)
```

## 最佳实践

1. **必填参数作为函数参数**：如 `host` 是必填的，直接作为 `NewServer` 的参数
2. **可选参数使用 Option**：所有配置项都通过 `...Option` 变长参数传递
3. **提供合理的默认值**：让最小化配置也能工作
4. **支持配置验证**：在构造函数中进行配置有效性检查
5. **使用预设配置**：为常见场景提供预设配置组合

## 参考

- [Self-referential functions and the design of options](https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html) - Rob Pike
- [Functional Options for Friendly APIs](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis) - Dave Cheney
