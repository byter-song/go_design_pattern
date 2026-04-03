# 观察者模式 (Observer Pattern)

## 概述

观察者模式定义了对象之间的一对多依赖关系，当一个对象（Subject）状态发生改变时，所有依赖于它的对象（Observers）都会收到通知并自动更新。

Go 语言通过 **接口（interface）** 和 **切片（slice）** 提供了简洁高效的观察者模式实现。

## 实现方式

### 核心组件

```
┌─────────────────────────────────────────────────────────────┐
│                    Subject (被观察者)                        │
│                  - NewsPublisher                             │
│                  - EventBus                                  │
│                  - BatchPublisher                            │
│                                                              │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  observers []Observer  // 观察者切片                 │    │
│  │  mu sync.RWMutex       // 并发安全锁                 │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                              │
│  - Register(observer)   // 注册观察者                        │
│  - Unregister(observer) // 注销观察者                        │
│  - Notify(event)        // 通知所有观察者                    │
└──────────────────────┬──────────────────────────────────────┘
                       │ 通知
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                   Observer (观察者接口)                      │
│                                                              │
│  type Observer interface {                                   │
│      Update(event Event)  // 接收通知                        │
│      GetID() string       // 唯一标识                        │
│  }                                                           │
│                                                              │
│  具体实现:                                                    │
│  - EmailSubscriber    // 邮件订阅者                          │
│  - SMSSubscriber      // 短信订阅者                          │
│  - PushSubscriber     // 推送订阅者                          │
│  - FilteredObserver   // 过滤观察者                          │
└─────────────────────────────────────────────────────────────┘
```

### 基础实现

```go
// Event 事件结构体
type Event struct {
    Type    string
    Data    interface{}
    Source  string
    Time    time.Time
}

// Observer 观察者接口
type Observer interface {
    Update(event Event)
    GetID() string
}

// Subject 被观察者接口
type Subject interface {
    Register(observer Observer)
    Unregister(observer Observer)
    Notify(event Event)
}
```

## 关于 Channel 异步事件通知的讨论

在实际项目中，观察者模式可以使用 Channel 实现异步事件通知。以下是设计思路和权衡：

### Channel 实现思路

```go
// AsyncPublisher 使用 Channel 的异步发布者
type AsyncPublisher struct {
    observers []Observer
    eventChan chan Event      // 事件 Channel
    quitChan  chan struct{}   // 退出信号
    mu        sync.RWMutex
}

func (ap *AsyncPublisher) Start() {
    go func() {
        for {
            select {
            case event := <-ap.eventChan:
                // 异步通知所有观察者
                ap.notifyObservers(event)
            case <-ap.quitChan:
                return
            }
        }
    }()
}

func (ap *AsyncPublisher) Publish(event Event) {
    ap.eventChan <- event  // 非阻塞发送
}
```

### 使用 Channel 的优缺点

**优点：**
1. **天然异步**：发布者无需等待观察者处理完成
2. **解耦彻底**：发布者和观察者完全解耦
3. **背压控制**：通过缓冲 Channel 控制事件堆积
4. **优雅关闭**：通过 quitChan 实现优雅关闭

**缺点：**
1. **复杂度增加**：需要管理 Goroutine 生命周期
2. **顺序保证**：Channel 不保证观察者接收顺序
3. **内存开销**：每个发布者需要维护 Channel
4. **调试困难**：异步逻辑难以追踪和调试

### 选择建议

| 场景 | 推荐方案 | 原因 |
|------|---------|------|
| 简单通知 | 同步切片 | 简单、直观、易调试 |
| 高并发 | Channel + 缓冲 | 避免阻塞发布者 |
| 事件顺序重要 | 同步切片 | 保证通知顺序 |
| 观察者耗时操作 | Channel + Worker Pool | 防止阻塞主流程 |
| 微服务通信 | 消息队列 | 跨进程、持久化 |

## 高级实现

### 1. 过滤观察者 (FilteredObserver)

只接收感兴趣的事件类型：

```go
type FilteredObserver struct {
    id          string
    eventTypes  map[string]bool
    handler     func(Event)
}

func (fo *FilteredObserver) Update(event Event) {
    if fo.eventTypes[event.Type] {
        fo.handler(event)
    }
}
```

### 2. 批量发布者 (BatchPublisher)

批量收集事件后统一通知：

```go
type BatchPublisher struct {
    observers []Observer
    batch     []Event
    ticker    *time.Ticker
    batchSize int
}

func (bp *BatchPublisher) startBatchProcessor() {
    go func() {
        for range bp.ticker.C {
            bp.flush()
        }
    }()
}
```

### 3. 事件总线 (EventBus)

全局事件分发中心：

```go
type EventBus struct {
    subscribers map[string][]EventHandler  // 按事件类型分组
    mu          sync.RWMutex
}

func (eb *EventBus) Subscribe(eventType string, handler EventHandler) {
    // 订阅特定类型事件
}

func (eb *EventBus) Publish(eventType string, data interface{}) {
    // 发布事件给所有订阅者
}
```

## 使用示例

### 新闻发布系统

```go
// 创建发布者
publisher := NewNewsPublisher()

// 创建观察者
emailSub := NewEmailSubscriber("user@example.com")
smsSub := NewSMSSubscriber("13800138000")
pushSub := NewPushSubscriber("device_token_123")

// 注册观察者
publisher.Register(emailSub)
publisher.Register(smsSub)
publisher.Register(pushSub)

// 发布新闻
event := Event{
    Type:   "breaking_news",
    Data:   "Go 2.0 正式发布！",
    Source: "TechNews",
}
publisher.Notify(event)

// 所有观察者都会收到通知
// EmailSubscriber: 发送邮件
// SMSSubscriber: 发送短信
// PushSubscriber: 发送推送
```

### 过滤观察者

```go
// 只接收特定类型事件
filteredSub := NewFilteredObserver(
    "price_alert",
    []string{"price_drop", "price_rise"},
    func(e Event) {
        fmt.Printf("价格变动: %v\n", e.Data)
    },
)

publisher.Register(filteredSub)

// 只有 price_drop 和 price_rise 事件会触发通知
```

### 事件总线

```go
// 全局事件总线
bus := NewEventBus()

// 订阅订单事件
bus.Subscribe("order_created", func(e Event) {
    fmt.Println("新订单创建:", e.Data)
})

bus.Subscribe("order_paid", func(e Event) {
    fmt.Println("订单已支付:", e.Data)
})

// 发布事件
bus.Publish("order_created", Order{ID: "123", Amount: 100})
bus.Publish("order_paid", Order{ID: "123"})
```

## 并发安全

使用 `sync.RWMutex` 保证并发安全：

```go
type NewsPublisher struct {
    observers []Observer
    mu        sync.RWMutex  // 读写锁
}

func (np *NewsPublisher) Register(observer Observer) {
    np.mu.Lock()
    defer np.mu.Unlock()
    np.observers = append(np.observers, observer)
}

func (np *NewsPublisher) Notify(event Event) {
    np.mu.RLock()
    defer np.mu.RUnlock()
    
    for _, observer := range np.observers {
        observer.Update(event)
    }
}
```

## 优势

1. **松耦合**：观察者和被观察者之间没有直接依赖
2. **可扩展**：可以动态添加/移除观察者
3. **广播通信**：一个事件可以通知多个接收者
4. **Go 语言特色**：
   - 接口实现灵活
   - 切片存储高效
   - 可选 Channel 实现异步

## 适用场景

- 事件驱动系统
- 消息发布/订阅
- 状态变更通知
- 日志收集系统
- 实时监控告警

## 参考

- [Go 并发模式](https://golang.org/doc/effective_go.html#concurrency)
- [Go Channel 最佳实践](https://golang.org/doc/effective_go.html#channels)
- [观察者模式 - 设计模式](https://en.wikipedia.org/wiki/Observer_pattern)
