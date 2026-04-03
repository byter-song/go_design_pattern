// Package observer 展示了 Go 语言中实现观察者模式的惯用方法。
//
// 观察者模式定义了对象之间的一对多依赖关系，当一个对象状态改变时，
// 所有依赖于它的对象都会得到通知并自动更新。
//
// Go 语言实现特点：
//   1. 使用 interface 定义观察者和被观察者契约
//   2. 使用切片（slice）存储观察者列表
//   3. 利用 sync.RWMutex 保证并发安全
//   4. 支持移除观察者（取消订阅）
//
// Channel 异步通知思路（在注释中讨论）：
//   在 Go 中，可以使用 Channel 实现异步事件通知，但需要考虑以下几点：
//
//   优点：
//   - 真正的异步处理，不会阻塞被观察者
//   - 可以利用 Go 的调度器高效处理并发
//   - 支持超时控制和取消操作（通过 context）
//
//   缺点：
//   - 增加了复杂度，需要管理 goroutine 生命周期
//   - 事件顺序可能无法保证（如果多个 goroutine 处理）
//   - 需要处理背压（backpressure）问题
//
//   实现思路：
//   ```go
//   type AsyncSubject struct {
//       observers []Observer
//       eventChan chan Event
//       mu        sync.RWMutex
//   }
//
//   func (s *AsyncSubject) Notify(event Event) {
//       select {
//       case s.eventChan <- event:
//           // 成功放入队列
//       default:
//           // 队列满，丢弃或处理
//       }
//   }
//
//   func (s *AsyncSubject) Start() {
//       go func() {
//           for event := range s.eventChan {
//               s.mu.RLock()
//               observers := s.observers
//               s.mu.RUnlock()
//               for _, observer := range observers {
//                   go observer.Update(event) // 每个观察者一个 goroutine
//               }
//           }
//       }()
//   }
//   ```
//
// 适用场景：
//   - 事件驱动系统
//   - 数据变更通知（MVC 中的模型更新通知视图）
//   - 消息队列消费者
//   - 状态监控和告警
package observer

import (
	"fmt"
	"sync"
	"time"
)

// ============================================================================
// 核心接口定义
// ============================================================================

// Observer 是观察者接口。
// 任何实现了 Update 方法的类型都可以作为观察者。
type Observer interface {
	// Update 接收事件通知
	Update(event Event)
	// GetID 返回观察者唯一标识
	GetID() string
}

// Subject 是被观察者接口（主题）。
// 维护观察者列表，并在状态变化时通知它们。
type Subject interface {
	// Attach 添加观察者
	Attach(observer Observer)
	// Detach 移除观察者
	Detach(observerID string)
	// Notify 通知所有观察者
	Notify(event Event)
}

// Event 是事件数据结构。
// 可以扩展以包含更多上下文信息。
type Event struct {
	Type      string
	Data      interface{}
	Timestamp time.Time
	Source    string
}

// ============================================================================
// 具体实现：新闻发布系统
// ============================================================================

// NewsPublisher 是新闻发布者（被观察者）。
// 使用切片存储观察者，并使用读写锁保证并发安全。
type NewsPublisher struct {
	observers []Observer
	mu        sync.RWMutex
	name      string
}

// NewNewsPublisher 创建新闻发布者
func NewNewsPublisher(name string) *NewsPublisher {
	return &NewsPublisher{
		observers: make([]Observer, 0),
		name:      name,
	}
}

// Attach 添加观察者。
// 使用写锁保证并发安全。
func (n *NewsPublisher) Attach(observer Observer) {
	n.mu.Lock()
	defer n.mu.Unlock()

	// 检查是否已存在
	for _, o := range n.observers {
		if o.GetID() == observer.GetID() {
			fmt.Printf("[Publisher] Observer %s already attached\n", observer.GetID())
			return
		}
	}

	n.observers = append(n.observers, observer)
	fmt.Printf("[Publisher] Observer %s attached. Total: %d\n", observer.GetID(), len(n.observers))
}

// Detach 移除观察者。
// 使用写锁保证并发安全，使用切片技巧删除元素。
func (n *NewsPublisher) Detach(observerID string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	for i, observer := range n.observers {
		if observer.GetID() == observerID {
			// 使用切片技巧删除元素
			n.observers = append(n.observers[:i], n.observers[i+1:]...)
			fmt.Printf("[Publisher] Observer %s detached. Total: %d\n", observerID, len(n.observers))
			return
		}
	}

	fmt.Printf("[Publisher] Observer %s not found\n", observerID)
}

// Notify 通知所有观察者。
// 使用读锁保证并发安全，遍历观察者列表。
func (n *NewsPublisher) Notify(event Event) {
	n.mu.RLock()
	observers := make([]Observer, len(n.observers))
	copy(observers, n.observers)
	n.mu.RUnlock()

	fmt.Printf("[Publisher] Notifying %d observers about: %s\n", len(observers), event.Type)

	for _, observer := range observers {
		observer.Update(event)
	}
}

// PublishNews 发布新闻
func (n *NewsPublisher) PublishNews(category, content string) {
	event := Event{
		Type:      "NEWS_" + category,
		Data:      content,
		Timestamp: time.Now(),
		Source:    n.name,
	}
	n.Notify(event)
}

// GetObserverCount 获取观察者数量
func (n *NewsPublisher) GetObserverCount() int {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return len(n.observers)
}

// ============================================================================
// 具体观察者实现
// ============================================================================

// EmailSubscriber 邮件订阅者。
// 实现了 Observer 接口。
type EmailSubscriber struct {
	id    string
	email string
}

// NewEmailSubscriber 创建邮件订阅者
func NewEmailSubscriber(id, email string) *EmailSubscriber {
	return &EmailSubscriber{
		id:    id,
		email: email,
	}
}

// Update 实现 Observer 接口
func (e *EmailSubscriber) Update(event Event) {
	fmt.Printf("[EmailSubscriber %s] Sending email to %s: %s - %v\n",
		e.id, e.email, event.Type, event.Data)
}

// GetID 实现 Observer 接口
func (e *EmailSubscriber) GetID() string {
	return e.id
}

// SMSSubscriber 短信订阅者。
// 实现了 Observer 接口。
type SMSSubscriber struct {
	id      string
	phoneNo string
}

// NewSMSSubscriber 创建短信订阅者
func NewSMSSubscriber(id, phoneNo string) *SMSSubscriber {
	return &SMSSubscriber{
		id:      id,
		phoneNo: phoneNo,
	}
}

// Update 实现 Observer 接口
func (s *SMSSubscriber) Update(event Event) {
	// 只处理紧急新闻
	if event.Type == "NEWS_URGENT" {
		fmt.Printf("[SMSSubscriber %s] Sending SMS to %s: %v\n",
			s.id, s.phoneNo, event.Data)
	} else {
		fmt.Printf("[SMSSubscriber %s] Ignoring non-urgent news: %s\n",
			s.id, event.Type)
	}
}

// GetID 实现 Observer 接口
func (s *SMSSubscriber) GetID() string {
	return s.id
}

// PushSubscriber 推送订阅者。
// 实现了 Observer 接口。
type PushSubscriber struct {
	id       string
	deviceID string
}

// NewPushSubscriber 创建推送订阅者
func NewPushSubscriber(id, deviceID string) *PushSubscriber {
	return &PushSubscriber{
		id:       id,
		deviceID: deviceID,
	}
}

// Update 实现 Observer 接口
func (p *PushSubscriber) Update(event Event) {
	fmt.Printf("[PushSubscriber %s] Pushing to device %s: %s - %v\n",
		p.id, p.deviceID, event.Type, event.Data)
}

// GetID 实现 Observer 接口
func (p *PushSubscriber) GetID() string {
	return p.id
}

// ============================================================================
// 高级实现：带过滤器的观察者
// ============================================================================

// FilteredObserver 是带过滤器的观察者。
// 只有满足条件的事件才会被处理。
type FilteredObserver struct {
	id        string
	observer  Observer
	filter    func(Event) bool
}

// NewFilteredObserver 创建带过滤器的观察者
func NewFilteredObserver(id string, observer Observer, filter func(Event) bool) *FilteredObserver {
	return &FilteredObserver{
		id:       id,
		observer: observer,
		filter:   filter,
	}
}

// Update 实现 Observer 接口，添加过滤逻辑
func (f *FilteredObserver) Update(event Event) {
	if f.filter(event) {
		f.observer.Update(event)
	}
}

// GetID 实现 Observer 接口
func (f *FilteredObserver) GetID() string {
	return f.id
}

// ============================================================================
// 高级实现：批量通知
// ============================================================================

// BatchPublisher 是批量发布者。
// 收集事件后批量通知，减少通知频率。
type BatchPublisher struct {
	observers []Observer
	mu        sync.RWMutex
	events    []Event
	batchSize int
	ticker    *time.Ticker
	stopChan  chan struct{}
}

// NewBatchPublisher 创建批量发布者
func NewBatchPublisher(batchSize int, interval time.Duration) *BatchPublisher {
	bp := &BatchPublisher{
		observers: make([]Observer, 0),
		events:    make([]Event, 0),
		batchSize: batchSize,
		ticker:    time.NewTicker(interval),
		stopChan:  make(chan struct{}),
	}

	// 启动后台 goroutine 定期刷新
	go bp.flushLoop()

	return bp
}

// Attach 添加观察者
func (b *BatchPublisher) Attach(observer Observer) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, o := range b.observers {
		if o.GetID() == observer.GetID() {
			return
		}
	}

	b.observers = append(b.observers, observer)
}

// Detach 移除观察者
func (b *BatchPublisher) Detach(observerID string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i, observer := range b.observers {
		if observer.GetID() == observerID {
			b.observers = append(b.observers[:i], b.observers[i+1:]...)
			return
		}
	}
}

// Notify 收集事件
func (b *BatchPublisher) Notify(event Event) {
	b.mu.Lock()
	b.events = append(b.events, event)
	shouldFlush := len(b.events) >= b.batchSize
	b.mu.Unlock()

	if shouldFlush {
		b.Flush()
	}
}

// Flush 立即刷新所有事件
func (b *BatchPublisher) Flush() {
	b.mu.Lock()
	if len(b.events) == 0 {
		b.mu.Unlock()
		return
	}

	events := make([]Event, len(b.events))
	copy(events, b.events)
	b.events = b.events[:0] // 清空切片
	observers := make([]Observer, len(b.observers))
	copy(observers, b.observers)
	b.mu.Unlock()

	fmt.Printf("[BatchPublisher] Flushing %d events to %d observers\n", len(events), len(observers))

	batchEvent := Event{
		Type:      "BATCH",
		Data:      events,
		Timestamp: time.Now(),
	}

	for _, observer := range observers {
		observer.Update(batchEvent)
	}
}

// flushLoop 定期刷新循环
func (b *BatchPublisher) flushLoop() {
	for {
		select {
		case <-b.ticker.C:
			b.Flush()
		case <-b.stopChan:
			return
		}
	}
}

// Stop 停止批量发布者
func (b *BatchPublisher) Stop() {
	close(b.stopChan)
	b.ticker.Stop()
	b.Flush() // 刷新剩余事件
}

// ============================================================================
// 高级实现：事件总线（全局观察者模式）
// ============================================================================

// EventBus 是全局事件总线。
// 支持按主题订阅，实现一对多、多对多的消息传递。
type EventBus struct {
	subscribers map[string][]Observer // topic -> observers
	mu          sync.RWMutex
}

// NewEventBus 创建事件总线
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]Observer),
	}
}

// Subscribe 订阅主题
func (e *EventBus) Subscribe(topic string, observer Observer) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// 检查是否已订阅
	observers := e.subscribers[topic]
	for _, o := range observers {
		if o.GetID() == observer.GetID() {
			return
		}
	}

	e.subscribers[topic] = append(observers, observer)
	fmt.Printf("[EventBus] Observer %s subscribed to topic: %s\n", observer.GetID(), topic)
}

// Unsubscribe 取消订阅
func (e *EventBus) Unsubscribe(topic, observerID string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	observers, ok := e.subscribers[topic]
	if !ok {
		return
	}

	for i, o := range observers {
		if o.GetID() == observerID {
			e.subscribers[topic] = append(observers[:i], observers[i+1:]...)
			fmt.Printf("[EventBus] Observer %s unsubscribed from topic: %s\n", observerID, topic)
			return
		}
	}
}

// Publish 发布事件到指定主题
func (e *EventBus) Publish(topic string, event Event) {
	e.mu.RLock()
	observers, ok := e.subscribers[topic]
	if !ok || len(observers) == 0 {
		e.mu.RUnlock()
		return
	}

	// 复制观察者列表，避免在遍历时被修改
	observersCopy := make([]Observer, len(observers))
	copy(observersCopy, observers)
	e.mu.RUnlock()

	fmt.Printf("[EventBus] Publishing to topic '%s', %d subscribers\n", topic, len(observersCopy))

	for _, observer := range observersCopy {
		observer.Update(event)
	}
}

// GetSubscriberCount 获取主题的订阅者数量
func (e *EventBus) GetSubscriberCount(topic string) int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.subscribers[topic])
}
