package observer

import (
	"testing"
	"time"
)

// ============================================================================
// 新闻发布系统测试
// ============================================================================

func TestNewsPublisher(t *testing.T) {
	t.Run("AttachAndDetach", func(t *testing.T) {
		publisher := NewNewsPublisher("TestPublisher")
		subscriber := NewEmailSubscriber("sub1", "test@example.com")

		// 添加观察者
		publisher.Attach(subscriber)
		if publisher.GetObserverCount() != 1 {
			t.Errorf("Expected 1 observer, got %d", publisher.GetObserverCount())
		}

		// 重复添加同一观察者（应该被忽略）
		publisher.Attach(subscriber)
		if publisher.GetObserverCount() != 1 {
			t.Errorf("Expected 1 observer after duplicate attach, got %d", publisher.GetObserverCount())
		}

		// 移除观察者
		publisher.Detach("sub1")
		if publisher.GetObserverCount() != 0 {
			t.Errorf("Expected 0 observers, got %d", publisher.GetObserverCount())
		}

		// 移除不存在的观察者
		publisher.Detach("nonexistent")
	})

	t.Run("Notify", func(t *testing.T) {
		publisher := NewNewsPublisher("TestPublisher")
		subscriber1 := NewEmailSubscriber("sub1", "test1@example.com")
		subscriber2 := NewEmailSubscriber("sub2", "test2@example.com")

		publisher.Attach(subscriber1)
		publisher.Attach(subscriber2)

		event := Event{
			Type:      "NEWS_TEST",
			Data:      "Test news content",
			Timestamp: time.Now(),
			Source:    "Test",
		}

		// 通知所有观察者
		publisher.Notify(event)
	})

	t.Run("PublishNews", func(t *testing.T) {
		publisher := NewNewsPublisher("TestPublisher")
		subscriber := NewEmailSubscriber("sub1", "test@example.com")

		publisher.Attach(subscriber)
		publisher.PublishNews("TECH", "New Go version released!")
	})
}

// ============================================================================
// 订阅者测试
// ============================================================================

func TestEmailSubscriber(t *testing.T) {
	subscriber := NewEmailSubscriber("email1", "user@example.com")

	t.Run("GetID", func(t *testing.T) {
		if subscriber.GetID() != "email1" {
			t.Errorf("Expected ID 'email1', got %s", subscriber.GetID())
		}
	})

	t.Run("Update", func(t *testing.T) {
		event := Event{
			Type:      "NEWS_TEST",
			Data:      "Test content",
			Timestamp: time.Now(),
		}
		subscriber.Update(event)
	})
}

func TestSMSSubscriber(t *testing.T) {
	subscriber := NewSMSSubscriber("sms1", "+1234567890")

	t.Run("GetID", func(t *testing.T) {
		if subscriber.GetID() != "sms1" {
			t.Errorf("Expected ID 'sms1', got %s", subscriber.GetID())
		}
	})

	t.Run("UpdateUrgent", func(t *testing.T) {
		// 紧急新闻应该被处理
		event := Event{
			Type:      "NEWS_URGENT",
			Data:      "Breaking news!",
			Timestamp: time.Now(),
		}
		subscriber.Update(event)
	})

	t.Run("UpdateNormal", func(t *testing.T) {
		// 普通新闻应该被忽略
		event := Event{
			Type:      "NEWS_NORMAL",
			Data:      "Regular news",
			Timestamp: time.Now(),
		}
		subscriber.Update(event)
	})
}

func TestPushSubscriber(t *testing.T) {
	subscriber := NewPushSubscriber("push1", "device-123")

	t.Run("GetID", func(t *testing.T) {
		if subscriber.GetID() != "push1" {
			t.Errorf("Expected ID 'push1', got %s", subscriber.GetID())
		}
	})

	t.Run("Update", func(t *testing.T) {
		event := Event{
			Type:      "NEWS_TEST",
			Data:      "Push notification",
			Timestamp: time.Now(),
		}
		subscriber.Update(event)
	})
}

// ============================================================================
// 过滤器观察者测试
// ============================================================================

func TestFilteredObserver(t *testing.T) {
	innerObserver := NewEmailSubscriber("inner1", "test@example.com")
	
	// 创建只处理 TECH 类型事件的过滤器
	filter := func(event Event) bool {
		return event.Type == "NEWS_TECH"
	}
	
	filtered := NewFilteredObserver("filtered1", innerObserver, filter)

	t.Run("GetID", func(t *testing.T) {
		if filtered.GetID() != "filtered1" {
			t.Errorf("Expected ID 'filtered1', got %s", filtered.GetID())
		}
	})

	t.Run("FilteredUpdate", func(t *testing.T) {
		// TECH 新闻应该通过
		techEvent := Event{
			Type:      "NEWS_TECH",
			Data:      "Go 1.22 released",
			Timestamp: time.Now(),
		}
		filtered.Update(techEvent)

		// SPORTS 新闻应该被过滤
		sportsEvent := Event{
			Type:      "NEWS_SPORTS",
			Data:      "Match result",
			Timestamp: time.Now(),
		}
		filtered.Update(sportsEvent)
	})
}

// ============================================================================
// 批量发布者测试
// ============================================================================

func TestBatchPublisher(t *testing.T) {
	t.Run("AttachAndNotify", func(t *testing.T) {
		publisher := NewBatchPublisher(3, 100*time.Millisecond)
		defer publisher.Stop()

		subscriber := NewEmailSubscriber("batch_sub1", "test@example.com")
		publisher.Attach(subscriber)

		// 发送事件（未达到批量大小）
		publisher.Notify(Event{Type: "EVENT_1", Data: "data1"})
		publisher.Notify(Event{Type: "EVENT_2", Data: "data2"})

		// 等待定期刷新
		time.Sleep(150 * time.Millisecond)
	})

	t.Run("BatchFlush", func(t *testing.T) {
		publisher := NewBatchPublisher(3, 10*time.Second) // 很长的间隔，依赖批量大小触发
		defer publisher.Stop()

		subscriber := NewEmailSubscriber("batch_sub2", "test@example.com")
		publisher.Attach(subscriber)

		// 发送 3 个事件，触发批量刷新
		publisher.Notify(Event{Type: "EVENT_1", Data: "data1"})
		publisher.Notify(Event{Type: "EVENT_2", Data: "data2"})
		publisher.Notify(Event{Type: "EVENT_3", Data: "data3"})

		// 给一点时间让通知完成
		time.Sleep(50 * time.Millisecond)
	})

	t.Run("ManualFlush", func(t *testing.T) {
		publisher := NewBatchPublisher(10, 10*time.Second)
		defer publisher.Stop()

		subscriber := NewEmailSubscriber("batch_sub3", "test@example.com")
		publisher.Attach(subscriber)

		publisher.Notify(Event{Type: "EVENT_1", Data: "data1"})
		publisher.Notify(Event{Type: "EVENT_2", Data: "data2"})

		// 手动刷新
		publisher.Flush()
	})
}

// ============================================================================
// 事件总线测试
// ============================================================================

func TestEventBus(t *testing.T) {
	t.Run("SubscribeAndPublish", func(t *testing.T) {
		bus := NewEventBus()
		subscriber := NewEmailSubscriber("bus_sub1", "test@example.com")

		bus.Subscribe("topic1", subscriber)
		
		if bus.GetSubscriberCount("topic1") != 1 {
			t.Errorf("Expected 1 subscriber, got %d", bus.GetSubscriberCount("topic1"))
		}

		event := Event{
			Type:      "TEST_EVENT",
			Data:      "test data",
			Timestamp: time.Now(),
		}

		bus.Publish("topic1", event)
		bus.Publish("topic2", event) // 没有订阅者的主题
	})

	t.Run("Unsubscribe", func(t *testing.T) {
		bus := NewEventBus()
		subscriber := NewEmailSubscriber("bus_sub2", "test@example.com")

		bus.Subscribe("topic1", subscriber)
		bus.Unsubscribe("topic1", "bus_sub2")

		if bus.GetSubscriberCount("topic1") != 0 {
			t.Errorf("Expected 0 subscribers, got %d", bus.GetSubscriberCount("topic1"))
		}

		// 取消订阅不存在的主题
		bus.Unsubscribe("nonexistent", "bus_sub2")
	})

	t.Run("MultipleSubscribers", func(t *testing.T) {
		bus := NewEventBus()
		subscriber1 := NewEmailSubscriber("bus_sub3", "test1@example.com")
		subscriber2 := NewEmailSubscriber("bus_sub4", "test2@example.com")

		bus.Subscribe("topic1", subscriber1)
		bus.Subscribe("topic1", subscriber2)

		if bus.GetSubscriberCount("topic1") != 2 {
			t.Errorf("Expected 2 subscribers, got %d", bus.GetSubscriberCount("topic1"))
		}

		event := Event{
			Type:      "TEST_EVENT",
			Data:      "test data",
			Timestamp: time.Now(),
		}

		bus.Publish("topic1", event)
	})

	t.Run("DuplicateSubscribe", func(t *testing.T) {
		bus := NewEventBus()
		subscriber := NewEmailSubscriber("bus_sub5", "test@example.com")

		bus.Subscribe("topic1", subscriber)
		bus.Subscribe("topic1", subscriber) // 重复订阅

		if bus.GetSubscriberCount("topic1") != 1 {
			t.Errorf("Expected 1 subscriber after duplicate subscribe, got %d", bus.GetSubscriberCount("topic1"))
		}
	})
}

// ============================================================================
// 接口兼容性测试
// ============================================================================

func TestInterfaceCompatibility(t *testing.T) {
	t.Run("ObserverInterface", func(t *testing.T) {
		var _ Observer = NewEmailSubscriber("test", "test@example.com")
		var _ Observer = NewSMSSubscriber("test", "+1234567890")
		var _ Observer = NewPushSubscriber("test", "device-123")
	})

	t.Run("SubjectInterface", func(t *testing.T) {
		var _ Subject = NewNewsPublisher("test")
	})
}

// ============================================================================
// 并发安全测试
// ============================================================================

func TestConcurrentAccess(t *testing.T) {
	t.Run("ConcurrentAttachDetach", func(t *testing.T) {
		publisher := NewNewsPublisher("ConcurrentTest")
		done := make(chan bool, 20)

		// 并发添加
		for i := 0; i < 10; i++ {
			go func(id int) {
				subscriber := NewEmailSubscriber(string(rune('0'+id)), "test@example.com")
				publisher.Attach(subscriber)
				done <- true
			}(i)
		}

		// 等待添加完成
		for i := 0; i < 10; i++ {
			<-done
		}

		// 并发移除
		for i := 0; i < 10; i++ {
			go func(id int) {
				publisher.Detach(string(rune('0' + id)))
				done <- true
			}(i)
		}

		// 等待移除完成
		for i := 0; i < 10; i++ {
			<-done
		}
	})

	t.Run("ConcurrentNotify", func(t *testing.T) {
		publisher := NewNewsPublisher("ConcurrentNotify")
		
		// 添加多个观察者
		for i := 0; i < 5; i++ {
			subscriber := NewEmailSubscriber(string(rune('0'+i)), "test@example.com")
			publisher.Attach(subscriber)
		}

		done := make(chan bool, 10)

		// 并发通知
		for i := 0; i < 10; i++ {
			go func(id int) {
				event := Event{
					Type:      "CONCURRENT_EVENT",
					Data:      id,
					Timestamp: time.Now(),
				}
				publisher.Notify(event)
				done <- true
			}(i)
		}

		// 等待所有通知完成
		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

// ============================================================================
// 使用场景测试
// ============================================================================

func TestNewsSystemScenario(t *testing.T) {
	// 场景：新闻发布系统，不同类型的订阅者接收不同类型的通知

	publisher := NewNewsPublisher("GlobalNews")

	// 邮件订阅者 - 接收所有新闻
	emailSub := NewEmailSubscriber("email_user", "user@example.com")
	publisher.Attach(emailSub)

	// 短信订阅者 - 只接收紧急新闻
	smsSub := NewSMSSubscriber("sms_user", "+1234567890")
	publisher.Attach(smsSub)

	// 推送订阅者 - 接收所有新闻
	pushSub := NewPushSubscriber("push_user", "device-abc123")
	publisher.Attach(pushSub)

	// 发布普通新闻
	publisher.PublishNews("TECH", "New smartphone announced")

	// 发布紧急新闻
	publisher.PublishNews("URGENT", "Severe weather warning!")

	// 验证观察者数量
	if publisher.GetObserverCount() != 3 {
		t.Errorf("Expected 3 observers, got %d", publisher.GetObserverCount())
	}
}

func TestEventBusScenario(t *testing.T) {
	// 场景：微服务间的事件总线通信

	bus := NewEventBus()

	// 订单服务订阅订单事件
	orderService := NewEmailSubscriber("order_service", "orders@company.com")
	bus.Subscribe("orders", orderService)

	// 库存服务订阅订单和库存事件
	inventoryService := NewEmailSubscriber("inventory_service", "inventory@company.com")
	bus.Subscribe("orders", inventoryService)
	bus.Subscribe("inventory", inventoryService)

	// 物流服务订阅订单事件
	shippingService := NewEmailSubscriber("shipping_service", "shipping@company.com")
	bus.Subscribe("orders", shippingService)

	// 发布订单创建事件
	orderEvent := Event{
		Type:      "ORDER_CREATED",
		Data:      map[string]string{"order_id": "ORD-123", "amount": "99.99"},
		Timestamp: time.Now(),
	}
	bus.Publish("orders", orderEvent)

	// 发布库存更新事件
	inventoryEvent := Event{
		Type:      "STOCK_UPDATED",
		Data:      map[string]string{"product_id": "PROD-456", "quantity": "100"},
		Timestamp: time.Now(),
	}
	bus.Publish("inventory", inventoryEvent)

	// 验证订阅数量
	if bus.GetSubscriberCount("orders") != 3 {
		t.Errorf("Expected 3 subscribers for 'orders', got %d", bus.GetSubscriberCount("orders"))
	}
	if bus.GetSubscriberCount("inventory") != 1 {
		t.Errorf("Expected 1 subscriber for 'inventory', got %d", bus.GetSubscriberCount("inventory"))
	}
}

func TestFilteredNotificationScenario(t *testing.T) {
	// 场景：用户只关心特定类型的新闻

	publisher := NewNewsPublisher("TechNews")

	// 基础订阅者
	baseSubscriber := NewEmailSubscriber("tech_user", "tech@example.com")

	// 只接收 AI 相关新闻的过滤器
	aiFilter := func(event Event) bool {
		data, ok := event.Data.(string)
		return ok && len(data) > 0 && contains(data, "AI")
	}
	aiFiltered := NewFilteredObserver("ai_filtered", baseSubscriber, aiFilter)

	publisher.Attach(aiFiltered)

	// 发布 AI 新闻（应该被接收）
	publisher.PublishNews("TECH", "New AI model released")

	// 发布普通新闻（应该被过滤）
	publisher.PublishNews("TECH", "New smartphone released")
}

// 辅助函数
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
