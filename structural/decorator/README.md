# 装饰器模式 (Decorator Pattern)

## 概述

装饰器模式允许向一个现有的对象添加新的功能，同时又不改变其结构。这种模式创建了一个装饰类，用来包装原有的类，并在保持类方法签名完整性的前提下，提供了额外的功能。

## Go 语言实现特点

在 Go 中，装饰器模式最常见的应用是 **HTTP 中间件**。Go 的函数是一等公民，这使得基于高阶函数的装饰器实现非常自然和惯用。

### 核心概念

```go
// 处理器函数类型
type Handler func(ctx context.Context, req Request) (Response, error)

// 中间件类型：接收一个 Handler，返回一个新的 Handler
type Middleware func(Handler) Handler
```

### 关键特性

1. **高阶函数**：中间件是接收函数并返回函数的函数
2. **链式组合**：多个中间件可以链式组合
3. **洋葱模型**：请求和响应按照中间件顺序处理

## 代码示例

### 基本中间件

```go
// 日志中间件
func LoggingMiddleware(next Handler) Handler {
    return func(ctx context.Context, req Request) (Response, error) {
        start := time.Now()
        log.Printf("[Request] %s %s", req.Method, req.Path)
        
        resp, err := next(ctx, req)
        
        log.Printf("[Response] %s %s - %v (%v)", 
            req.Method, req.Path, resp.Status, time.Since(start))
        return resp, err
    }
}
```

### 中间件链

```go
// Chain 将多个中间件组合成一个
func Chain(middlewares ...Middleware) Middleware {
    return func(final Handler) Handler {
        // 从后向前包装
        for i := len(middlewares) - 1; i >= 0; i-- {
            final = middlewares[i](final)
        }
        return final
    }
}

// 使用
chain := Chain(
    RecoveryMiddleware,
    LoggingMiddleware,
    AuthMiddleware,
)
handler := chain(businessHandler)
```

### HTTP 标准库兼容

```go
// 适配标准库的 http.Handler
type HTTPMiddleware func(http.Handler) http.Handler

func LoggingHTTPMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}
```

## 使用场景

1. **日志记录**：记录请求和响应信息
2. **认证授权**：验证用户身份和权限
3. **错误恢复**：捕获 panic，防止服务崩溃
4. **超时控制**：限制请求处理时间
5. **指标收集**：收集性能指标和统计数据
6. **请求限流**：控制请求频率

## 优缺点

### 优点

- **单一职责**：每个中间件只负责一个功能
- **可组合**：灵活组合多个中间件
- **可复用**：中间件可以在不同处理器间复用
- **开闭原则**：无需修改原有代码即可添加功能

### 缺点

- **调用链过长**：过多中间件可能影响性能
- **调试困难**：调用栈变深，问题定位复杂

## 与其他模式的关系

- **代理模式**：装饰器关注添加功能，代理模式关注控制访问
- **责任链模式**：装饰器是静态组合，责任链是动态传递

## 参考

- [Go 语言 HTTP Middleware 最佳实践](https://drstearns.github.io/tutorials/gomiddleware/)
- [Alice - Go 中间件链库](https://github.com/justinas/alice)
