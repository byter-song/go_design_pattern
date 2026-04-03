# 责任链模式 (Chain of Responsibility Pattern)

## 概述

责任链模式允许多个对象有机会处理请求，从而避免请求的发送者和接收者之间的耦合。将这些对象连成一条链，并沿着这条链传递请求，直到有一个对象处理它为止。

Go 语言通过 **接口（interface）**、**函数类型（func type）** 和 **闭包** 提供了多种灵活的责任链实现方式。

## 实现方式

### 1. 传统接口实现

使用接口定义处理器，每个处理器持有下一个处理器的引用。

```
┌─────────────────────────────────────────────────────────────┐
│                      SupportHandler                         │
│                      (处理器接口)                            │
│                                                              │
│  type SupportHandler interface {                             │
│      Handle(ticket *SupportTicket) bool  // 处理请求         │
│      SetNext(handler SupportHandler) SupportHandler          │
│  }                                                           │
└──────────────────────┬──────────────────────────────────────┘
                       │
         ┌─────────────┼─────────────┐
         ▼             ▼             ▼
┌────────────────┬────────────┬──────────────┐
│ FirstLine      │ Technical  │ Senior       │
│ Support        │ Support    │ Support      │
│ (一线客服)      │ (技术支持)  │ (高级支持)    │
└───────┬────────┴─────┬──────┴──────┬───────┘
        │              │             │
        ▼              ▼             ▼
   简单问题        技术问题       复杂问题
   (直接解决)      (技术解决)      (升级处理)
```

```go
// SupportHandler 支持处理器接口
type SupportHandler interface {
    Handle(ticket *SupportTicket) bool
    SetNext(handler SupportHandler) SupportHandler
}

// BaseHandler 基础处理器
type BaseHandler struct {
    next SupportHandler
}

func (h *BaseHandler) SetNext(handler SupportHandler) SupportHandler {
    h.next = handler
    return handler
}

func (h *BaseHandler) PassToNext(ticket *SupportTicket) bool {
    if h.next != nil {
        return h.next.Handle(ticket)
    }
    return false
}
```

### 2. 函数类型实现

使用函数类型构建轻量级责任链。

```go
// HandlerFunc 处理器函数类型
type HandlerFunc func(*Request) *Response

// Chain 责任链
type Chain struct {
    handlers []HandlerFunc
    index    int
}

func (c *Chain) Next(req *Request) *Response {
    if c.index < len(c.handlers) {
        handler := c.handlers[c.index]
        c.index++
        return handler(req)
    }
    return nil
}
```

### 3. 中间件风格实现

参考 HTTP 中间件的设计模式。

```go
// Middleware 中间件类型
type Middleware func(ctx *Context) bool

// Context 上下文
type Context struct {
    Request  *Request
    Response *Response
    index    int
    chain    []Middleware
}

func (ctx *Context) Next() bool {
    if ctx.index < len(ctx.chain) {
        middleware := ctx.chain[ctx.index]
        ctx.index++
        return middleware(ctx)
    }
    return true
}
```

### 4. 构建器模式构建链

使用构建器模式（Builder Pattern）优雅地构建责任链。

```go
// ChainBuilder 责任链构建器
type ChainBuilder struct {
    first SupportHandler
    last  SupportHandler
}

func (cb *ChainBuilder) AddHandler(handler SupportHandler) *ChainBuilder {
    if cb.first == nil {
        cb.first = handler
        cb.last = handler
    } else {
        cb.last.SetNext(handler)
        cb.last = handler
    }
    return cb
}

func (cb *ChainBuilder) Build() SupportHandler {
    return cb.first
}

// 使用
chain := NewChainBuilder().
    AddHandler(NewFirstLineSupport()).
    AddHandler(NewTechnicalSupport()).
    AddHandler(NewSeniorSupport()).
    Build()
```

## 架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                        Client                                │
│                   (请求发送者)                                │
└──────────────────────┬──────────────────────────────────────┘
                       │ 发送请求
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                      Handler Chain                           │
│                                                              │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐     │
│  │  Handler 1  │───▶│  Handler 2  │───▶│  Handler 3  │     │
│  │             │    │             │    │             │     │
│  │ 能处理?     │    │ 能处理?     │    │ 能处理?     │     │
│  │  Yes:处理   │    │  Yes:处理   │    │  Yes:处理   │     │
│  │  No:传递    │    │  No:传递    │    │  No:传递    │     │
│  └─────────────┘    └─────────────┘    └─────────────┘     │
│        │                  │                  │             │
│        ▼                  ▼                  ▼             │
│     请求结束            请求结束            请求结束         │
└─────────────────────────────────────────────────────────────┘
```

## 核心代码解析

### 工单处理系统

```go
// SupportTicket 支持工单
type SupportTicket struct {
    ID          string
    Level       int           // 1:简单, 2:技术, 3:账单, 4:高级, 5:管理
    Description string
    Customer    string
}

// FirstLineSupport 一线支持
type FirstLineSupport struct {
    BaseHandler
}

func (f *FirstLineSupport) Handle(ticket *SupportTicket) bool {
    if ticket.Level <= 1 {
        fmt.Printf("一线客服处理工单: %s\n", ticket.ID)
        return true
    }
    return f.PassToNext(ticket)
}
```

### 审批流程系统

```go
// ApprovalRequest 审批请求
type ApprovalRequest struct {
    ID       string
    Amount   float64
    Requester string
}

// TeamLeader 团队负责人审批
type TeamLeader struct {
    BaseHandler
    limit float64
}

func (t *TeamLeader) Handle(request *ApprovalRequest) bool {
    if request.Amount <= t.limit {
        fmt.Printf("团队负责人审批通过: ¥%.2f\n", request.Amount)
        return true
    }
    return t.PassToNext(request)
}
```

## 使用示例

### 客服工单系统

```go
// 构建责任链
firstLine := NewFirstLineSupport()
technical := NewTechnicalSupport()
billing := NewBillingSupport()
senior := NewSeniorSupport()
manager := NewManagerSupport()

// 设置链式关系
firstLine.SetNext(technical).SetNext(billing).SetNext(senior).SetNext(manager)

// 处理不同级别的工单
ticket1 := &SupportTicket{ID: "T001", Level: 1, Description: "密码重置"}
firstLine.Handle(ticket1)
// 输出: 一线客服处理工单: T001

ticket2 := &SupportTicket{ID: "T002", Level: 3, Description: "退款申请"}
firstLine.Handle(ticket2)
// 输出: 账单部门处理工单: T002

ticket3 := &SupportTicket{ID: "T003", Level: 5, Description: "重大投诉"}
firstLine.Handle(ticket3)
// 输出: 经理处理工单: T003
```

### 使用构建器模式

```go
// 使用构建器优雅地构建责任链
chain := NewChainBuilder().
    AddHandler(NewFirstLineSupport()).
    AddHandler(NewTechnicalSupport()).
    AddHandler(NewBillingSupport()).
    AddHandler(NewSeniorSupport()).
    AddHandler(NewManagerSupport()).
    Build()

// 处理请求
ticket := &SupportTicket{ID: "T004", Level: 2, Description: "系统错误"}
chain.Handle(ticket)
// 输出: 技术支持处理工单: T004
```

### 审批流程

```go
// 构建审批链
approvalChain := NewApprovalChainBuilder().
    AddHandler(NewTeamLeader(1000)).      // 团队负责人: <= 1000
    AddHandler(NewDepartmentManager(5000)). // 部门经理: <= 5000
    AddHandler(NewCFO(50000)).             // CFO: <= 50000
    AddHandler(NewCEO()).                  // CEO: 无限制
    Build()

// 不同金额的审批请求
requests := []*ApprovalRequest{
    {ID: "R1", Amount: 500, Requester: "张三"},
    {ID: "R2", Amount: 3000, Requester: "李四"},
    {ID: "R3", Amount: 30000, Requester: "王五"},
    {ID: "R4", Amount: 100000, Requester: "赵六"},
}

for _, req := range requests {
    approvalChain.Handle(req)
}
// 输出:
// 团队负责人审批通过: ¥500.00
// 部门经理审批通过: ¥3000.00
// CFO审批通过: ¥30000.00
// CEO审批通过: ¥100000.00
```

### 中间件风格链

```go
// 创建中间件链
middlewareChain := NewMiddlewareChain()

// 添加中间件
middlewareChain.Use(LoggingMiddleware)      // 日志记录
middlewareChain.Use(AuthMiddleware)          // 认证
middlewareChain.Use(RateLimitMiddleware)     // 限流
middlewareChain.Use(BusinessLogicMiddleware) // 业务逻辑

// 执行
ctx := &Context{
    Request: &Request{Path: "/api/users", Token: "valid_token"},
}
middlewareChain.Execute(ctx)
```

## 优势

1. **降低耦合**：发送者不需要知道谁处理请求
2. **动态组合**：可以运行时调整链的结构
3. **单一职责**：每个处理器只负责自己的逻辑
4. **开闭原则**：新增处理器无需修改现有代码
5. **Go 语言特色**：
   - 接口实现灵活
   - 函数类型支持轻量级实现
   - 闭包支持状态保持

## 适用场景

- 多级审批流程
- 客服工单升级
- HTTP 中间件链
- 日志过滤器链
- 权限检查链
- 数据处理管道

## 对比：实现方式选择

| 实现方式 | 复杂度 | 灵活性 | 适用场景 |
|---------|--------|--------|---------|
| 传统接口 | 中 | 高 | 复杂业务逻辑，需要维护状态 |
| 函数类型 | 低 | 中 | 简单处理逻辑，无状态 |
| 中间件风格 | 中 | 高 | HTTP 处理，Web 中间件 |
| 构建器模式 | 中 | 高 | 需要动态构建链的场景 |

## 注意事项

1. **链断裂**：确保链的末端有默认处理或返回未处理标记
2. **循环引用**：避免在链中形成循环
3. **性能考虑**：链过长可能影响性能
4. **错误处理**：考虑在链中传递错误信息

## 参考

- [Go 接口最佳实践](https://golang.org/doc/effective_go.html#interfaces)
- [HTTP Middleware Pattern](https://drstearns.github.io/tutorials/gomiddleware/)
- [责任链模式 - 设计模式](https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern)
