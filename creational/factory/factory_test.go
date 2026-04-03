package factory

import (
	"testing"
)

// TestSimpleFactory 测试简单工厂
func TestSimpleFactory(t *testing.T) {
	tests := []struct {
		name        string
		notifyType  NotificationType
		config      map[string]string
		wantChannel string
		wantErr     bool
	}{
		{
			name:        "创建 SMS 通知",
			notifyType:  TypeSMS,
			config:      map[string]string{"provider": "twilio", "apiKey": "key123"},
			wantChannel: "sms",
			wantErr:     false,
		},
		{
			name:        "创建 Email 通知",
			notifyType:  TypeEmail,
			config:      map[string]string{"smtpServer": "smtp.example.com", "username": "user"},
			wantChannel: "email",
			wantErr:     false,
		},
		{
			name:        "创建 Push 通知",
			notifyType:  TypePush,
			config:      map[string]string{"platform": "iOS", "appID": "com.example.app"},
			wantChannel: "push",
			wantErr:     false,
		},
		{
			name:        "不支持的通知类型",
			notifyType:  "unknown",
			config:      map[string]string{},
			wantChannel: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notification, err := CreateNotification(tt.notifyType, tt.config)

			if tt.wantErr {
				if err == nil {
					t.Errorf("期望返回错误，但没有")
				}
				return
			}

			if err != nil {
				t.Errorf("不期望错误，但得到: %v", err)
				return
			}

			if notification.Channel() != tt.wantChannel {
				t.Errorf("期望渠道 %s，实际 %s", tt.wantChannel, notification.Channel())
			}

			t.Logf("✓ 成功创建 %s 通知", tt.wantChannel)
		})
	}
}

// TestFactoryMethod 测试工厂方法模式
func TestFactoryMethod(t *testing.T) {
	// 创建工厂注册表
	registry := NewFactoryRegistry()

	// 注册各种工厂
	registry.Register(TypeSMS, &SMSFactory{})
	registry.Register(TypeEmail, &EmailFactory{})
	registry.Register(TypePush, &PushFactory{})

	// 测试创建 SMS
	sms, err := registry.Create(TypeSMS, map[string]string{
		"provider": "twilio",
		"apiKey":   "test-key",
	})
	if err != nil {
		t.Errorf("创建 SMS 失败: %v", err)
	}
	if sms.Channel() != "sms" {
		t.Errorf("期望 sms，实际 %s", sms.Channel())
	}

	// 测试创建 Email
	email, err := registry.Create(TypeEmail, map[string]string{
		"smtpServer": "smtp.example.com",
		"username":   "test@example.com",
	})
	if err != nil {
		t.Errorf("创建 Email 失败: %v", err)
	}
	if email.Channel() != "email" {
		t.Errorf("期望 email，实际 %s", email.Channel())
	}

	// 测试未注册的工厂
	_, err = registry.Create("unknown", map[string]string{})
	if err == nil {
		t.Error("未注册的工厂应该返回错误")
	}

	t.Log("✓ 工厂方法模式测试通过")
}

// TestFunctionalFactory 测试函数式工厂
func TestFunctionalFactory(t *testing.T) {
	registry := NewFuncRegistry()

	// 注册工厂函数
	registry.Register(TypeSMS, func(config map[string]string) (Notification, error) {
		return &SMSNotification{
			provider: config["provider"],
			apiKey:   config["apiKey"],
		}, nil
	})

	registry.Register(TypeEmail, func(config map[string]string) (Notification, error) {
		return &EmailNotification{
			smtpServer: config["smtpServer"],
			port:       587,
			username:   config["username"],
			password:   config["password"],
		}, nil
	})

	// 测试创建
	sms, err := registry.Create(TypeSMS, map[string]string{
		"provider": "twilio",
		"apiKey":   "func-test",
	})
	if err != nil {
		t.Errorf("函数式工厂创建 SMS 失败: %v", err)
	}
	if sms.Channel() != "sms" {
		t.Errorf("期望 sms，实际 %s", sms.Channel())
	}

	t.Log("✓ 函数式工厂测试通过")
}

// TestConfiguredFactory 测试预配置的工厂函数（闭包）
func TestConfiguredFactory(t *testing.T) {
	// 创建一个预配置了默认参数的工厂
	factory := CreateConfiguredSMSFactory("default-provider", "default-api-key")

	// 测试 1：使用默认配置
	sms1, err := factory(map[string]string{})
	if err != nil {
		t.Fatalf("创建失败: %v", err)
	}

	// 验证使用了默认值（通过 Send 方法输出间接验证）
	err = sms1.Send("123456", "test")
	if err != nil {
		t.Errorf("发送失败: %v", err)
	}

	// 测试 2：覆盖部分配置
	sms2, err := factory(map[string]string{
		"apiKey": "special-key",
	})
	if err != nil {
		t.Fatalf("创建失败: %v", err)
	}
	err = sms2.Send("789012", "test2")
	if err != nil {
		t.Errorf("发送失败: %v", err)
	}

	// 测试 3：完全覆盖配置
	sms3, err := factory(map[string]string{
		"provider": "custom-provider",
		"apiKey":   "custom-key",
	})
	if err != nil {
		t.Fatalf("创建失败: %v", err)
	}
	err = sms3.Send("345678", "test3")
	if err != nil {
		t.Errorf("发送失败: %v", err)
	}

	t.Log("✓ 预配置工厂（闭包）测试通过")
}

// TestAbstractFactory 测试抽象工厂模式
func TestAbstractFactory(t *testing.T) {
	// 测试生产环境工具包
	t.Run("ProductionKit", func(t *testing.T) {
		kit := NewProductionKit(
			map[string]string{"provider": "twilio", "apiKey": "prod-key"},
			map[string]string{"smtpServer": "smtp.prod.com", "username": "prod@example.com"},
			map[string]string{"platform": "iOS", "appID": "com.prod.app"},
		)

		sms := kit.CreateSMS()
		if sms.Channel() != "sms" {
			t.Errorf("期望 sms，实际 %s", sms.Channel())
		}

		email := kit.CreateEmail()
		if email.Channel() != "email" {
			t.Errorf("期望 email，实际 %s", email.Channel())
		}

		push := kit.CreatePush()
		if push.Channel() != "push" {
			t.Errorf("期望 push，实际 %s", push.Channel())
		}

		t.Log("✓ ProductionKit 测试通过")
	})

	// 测试 Mock 工具包
	t.Run("MockKit", func(t *testing.T) {
		kit := &MockKit{}

		sms := kit.CreateSMS()
		if sms.Channel() != "sms" {
			t.Errorf("期望 sms，实际 %s", sms.Channel())
		}

		// 测试 Mock 功能
		sms.Send("123456", "test message")

		// 验证可以类型断言获取 Mock 特有方法
		if mockSMS, ok := sms.(*MockNotification); ok {
			receiver, message := mockSMS.GetLastCall()
			if receiver != "123456" || message != "test message" {
				t.Error("Mock 记录调用参数失败")
			}
		} else {
			t.Error("类型断言失败")
		}

		t.Log("✓ MockKit 测试通过")
	})
}

// TestNotificationSend 测试通知发送功能
func TestNotificationSend(t *testing.T) {
	tests := []struct {
		name      string
		notify    Notification
		receiver  string
		message   string
		wantError bool
	}{
		{
			name: "发送 SMS",
			notify: &SMSNotification{
				provider: "twilio",
				apiKey:   "test",
			},
			receiver:  "+8613800138000",
			message:   "Hello SMS",
			wantError: false,
		},
		{
			name: "发送 Email",
			notify: &EmailNotification{
				smtpServer: "smtp.example.com",
				port:       587,
				username:   "sender@example.com",
				password:   "password",
			},
			receiver:  "receiver@example.com",
			message:   "Hello Email",
			wantError: false,
		},
		{
			name: "发送 Push",
			notify: &PushNotification{
				platform: "Android",
				appID:    "com.test.app",
			},
			receiver:  "device-token-123",
			message:   "Hello Push",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.notify.Send(tt.receiver, tt.message)
			if tt.wantError && err == nil {
				t.Error("期望错误但未得到")
			}
			if !tt.wantError && err != nil {
				t.Errorf("不期望错误但得到: %v", err)
			}
		})
	}
}

// TestInterfaceCompliance 验证所有产品都正确实现了接口
func TestInterfaceCompliance(t *testing.T) {
	// 编译时检查：确保所有具体类型都实现了 Notification 接口
	var _ Notification = (*SMSNotification)(nil)
	var _ Notification = (*EmailNotification)(nil)
	var _ Notification = (*PushNotification)(nil)
	var _ Notification = (*MockNotification)(nil)

	// 编译时检查：确保所有工厂都实现了 NotificationFactory 接口
	var _ NotificationFactory = (*SMSFactory)(nil)
	var _ NotificationFactory = (*EmailFactory)(nil)
	var _ NotificationFactory = (*PushFactory)(nil)

	// 编译时检查：确保所有工具包都实现了 NotificationKit 接口
	var _ NotificationKit = (*ProductionKit)(nil)
	var _ NotificationKit = (*MockKit)(nil)

	t.Log("✓ 所有接口实现检查通过")
}

// BenchmarkSimpleFactory 基准测试：简单工厂
func BenchmarkSimpleFactory(b *testing.B) {
	config := map[string]string{
		"provider": "twilio",
		"apiKey":   "bench-key",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CreateNotification(TypeSMS, config)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkFactoryRegistry 基准测试：工厂注册表
func BenchmarkFactoryRegistry(b *testing.B) {
	registry := NewFactoryRegistry()
	registry.Register(TypeSMS, &SMSFactory{})

	config := map[string]string{
		"provider": "twilio",
		"apiKey":   "bench-key",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := registry.Create(TypeSMS, config)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkFuncRegistry 基准测试：函数式工厂注册表
func BenchmarkFuncRegistry(b *testing.B) {
	registry := NewFuncRegistry()
	registry.Register(TypeSMS, func(config map[string]string) (Notification, error) {
		return &SMSNotification{
			provider: config["provider"],
			apiKey:   config["apiKey"],
		}, nil
	})

	config := map[string]string{
		"provider": "twilio",
		"apiKey":   "bench-key",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := registry.Create(TypeSMS, config)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ExampleCreateNotification 示例：简单工厂用法
func ExampleCreateNotification() {
	// 创建 SMS 通知
	sms, _ := CreateNotification(TypeSMS, map[string]string{
		"provider": "twilio",
		"apiKey":   "my-api-key",
	})

	// 使用实例
	_ = sms.Channel()

	// Output:
}

// ExampleFactoryRegistry 示例：工厂注册表用法
func ExampleFactoryRegistry() {
	// 创建注册表并注册工厂
	registry := NewFactoryRegistry()
	registry.Register(TypeSMS, &SMSFactory{})
	registry.Register(TypeEmail, &EmailFactory{})

	// 使用注册表创建通知
	sms, _ := registry.Create(TypeSMS, map[string]string{
		"provider": "twilio",
		"apiKey":   "key",
	})

	_ = sms.Channel()

	// Output:
}

// ExampleCreateConfiguredSMSFactory 示例：预配置工厂用法
func ExampleCreateConfiguredSMSFactory() {
	// 创建预配置了默认值的工厂
	factory := CreateConfiguredSMSFactory("twilio", "default-key")

	// 使用默认配置创建
	sms1, _ := factory(map[string]string{})
	_ = sms1.Channel()

	// 覆盖部分配置
	sms2, _ := factory(map[string]string{"apiKey": "special-key"})
	_ = sms2.Channel()

	// Output:
}
