// Package factory 展示了 Go 语言中实现工厂方法模式的惯用方法。
//
// 工厂方法模式定义了一个创建对象的接口，但由子类决定实例化哪个类。
// 工厂方法让类的实例化推迟到子类。
//
// 关键概念：
//   - 将对象创建逻辑从使用逻辑中分离
//   - 通过接口而非具体类型编程
//   - 符合开闭原则：新增产品类型无需修改现有代码
//
// Go 语言特色：
//   - 利用隐式接口实现松耦合
//   - 使用函数类型作为"工厂"，无需复杂的类层次结构
//   - 通过闭包实现配置化的工厂
package factory

import (
	"fmt"
)

// ============================================================================
// 核心接口定义
// ============================================================================

// Notification 是通知接口，定义了所有通知类型必须实现的方法。
//
// 设计决策：
//   1. 使用小写接口名（idiomatic Go），因为接口通常不导出
//   2. 接口定义在工厂包中，而不是产品包中
//   3. 方法名简洁明了，符合 Go 命名习惯
type Notification interface {
	// Send 发送通知到指定接收者
	// 参数：receiver - 接收者标识（手机号、邮箱等）
	// 参数：message - 消息内容
	// 返回：错误信息（如果有）
	Send(receiver, message string) error

	// Channel 返回通知渠道名称
	Channel() string
}

// ============================================================================
// 具体产品实现
// ============================================================================

// SMSNotification 是短信通知的具体实现。
//
// 设计特点：
//   - 结构体小写开头，包外不可直接实例化
//   - 强制通过工厂方法创建，保证创建逻辑统一
type SMSNotification struct {
	// provider 短信服务提供商
	provider string

	// apiKey API 密钥
	apiKey string
}

// Send 实现 Notification 接口
func (s *SMSNotification) Send(receiver, message string) error {
	// 模拟发送短信
	fmt.Printf("[SMS via %s] To: %s, Message: %s\n", s.provider, receiver, message)
	return nil
}

// Channel 返回渠道名称
func (s *SMSNotification) Channel() string {
	return "sms"
}

// EmailNotification 是邮件通知的具体实现
type EmailNotification struct {
	smtpServer string
	port       int
	username   string
	password   string
}

// Send 实现 Notification 接口
func (e *EmailNotification) Send(receiver, message string) error {
	fmt.Printf("[Email via %s:%d] To: %s, Message: %s\n", e.smtpServer, e.port, receiver, message)
	return nil
}

// Channel 返回渠道名称
func (e *EmailNotification) Channel() string {
	return "email"
}

// PushNotification 是推送通知的具体实现
type PushNotification struct {
	platform string // iOS, Android, etc.
	appID    string
}

// Send 实现 Notification 接口
func (p *PushNotification) Send(receiver, message string) error {
	fmt.Printf("[Push via %s] To: %s, Message: %s\n", p.platform, receiver, message)
	return nil
}

// Channel 返回渠道名称
func (p *PushNotification) Channel() string {
	return "push"
}

// ============================================================================
// 简单工厂（Simple Factory）- 最常用
// ============================================================================

// NotificationType 是通知类型的枚举
type NotificationType string

const (
	TypeSMS   NotificationType = "sms"
	TypeEmail NotificationType = "email"
	TypePush  NotificationType = "push"
)

// CreateNotification 是简单工厂函数，根据类型创建对应的通知实例。
//
// 这是 Go 中最常用的工厂模式形式，简单直接。
//
// 优点：
//   - 集中管理对象创建逻辑
//   - 调用者无需知道具体类名
//   - 易于添加新类型
//
// 缺点：
//   - 增加新类型时需要修改此函数（违反开闭原则）
//   - 如果类型很多，函数会变得很长
//
// 参数：
//   - notifyType: 通知类型
//   - config: 配置参数（根据类型不同解释不同）
//
// 使用示例：
//
//	sms, err := factory.CreateNotification(factory.TypeSMS, map[string]string{
//	    "provider": "twilio",
//	    "apiKey":   "xxx",
//	})
func CreateNotification(notifyType NotificationType, config map[string]string) (Notification, error) {
	switch notifyType {
	case TypeSMS:
		return &SMSNotification{
			provider: config["provider"],
			apiKey:   config["apiKey"],
		}, nil

	case TypeEmail:
		// 将字符串 port 转换为 int（简化处理）
		port := 587
		if p := config["port"]; p != "" {
			// 实际项目中应该使用 strconv.Atoi
			port = 587
		}
		return &EmailNotification{
			smtpServer: config["smtpServer"],
			port:       port,
			username:   config["username"],
			password:   config["password"],
		}, nil

	case TypePush:
		return &PushNotification{
			platform: config["platform"],
			appID:    config["appID"],
		}, nil

	default:
		return nil, fmt.Errorf("unsupported notification type: %s", notifyType)
	}
}

// ============================================================================
// 工厂方法模式（Factory Method）- 更灵活
// ============================================================================

// NotificationFactory 是工厂接口，定义创建通知的方法。
//
// 当需要不同的创建策略（如 Mock 工厂、缓存工厂）时，实现此接口。
type NotificationFactory interface {
	Create(config map[string]string) (Notification, error)
}

// SMSFactory 是短信通知的工厂
type SMSFactory struct{}

// Create 实现 NotificationFactory 接口
func (f *SMSFactory) Create(config map[string]string) (Notification, error) {
	return &SMSNotification{
		provider: config["provider"],
		apiKey:   config["apiKey"],
	}, nil
}

// EmailFactory 是邮件通知的工厂
type EmailFactory struct{}

// Create 实现 NotificationFactory 接口
func (f *EmailFactory) Create(config map[string]string) (Notification, error) {
	port := 587
	return &EmailNotification{
		smtpServer: config["smtpServer"],
		port:       port,
		username:   config["username"],
		password:   config["password"],
	}, nil
}

// PushFactory 是推送通知的工厂
type PushFactory struct{}

// Create 实现 NotificationFactory 接口
func (f *PushFactory) Create(config map[string]string) (Notification, error) {
	return &PushNotification{
		platform: config["platform"],
		appID:    config["appID"],
	}, nil
}

// FactoryRegistry 是工厂注册表，管理所有工厂。
//
// 这是简单工厂和工厂方法的结合，既保持了简单性，又支持扩展。
type FactoryRegistry struct {
	factories map[NotificationType]NotificationFactory
}

// NewFactoryRegistry 创建工厂注册表
func NewFactoryRegistry() *FactoryRegistry {
	return &FactoryRegistry{
		factories: make(map[NotificationType]NotificationFactory),
	}
}

// Register 注册工厂
func (r *FactoryRegistry) Register(notifyType NotificationType, factory NotificationFactory) {
	r.factories[notifyType] = factory
}

// Create 使用注册的工厂创建通知实例
func (r *FactoryRegistry) Create(notifyType NotificationType, config map[string]string) (Notification, error) {
	factory, ok := r.factories[notifyType]
	if !ok {
		return nil, fmt.Errorf("no factory registered for type: %s", notifyType)
	}
	return factory.Create(config)
}

// ============================================================================
// 函数式工厂 - Go 惯用法
// ============================================================================

// NotificationFunc 是函数类型的工厂，这是 Go 特有的简洁实现。
//
// 相比接口方式，函数类型更轻量，适合简单的创建逻辑。
// 这是 Go 中"函数是一等公民"理念的体现。
type NotificationFunc func(config map[string]string) (Notification, error)

// FuncRegistry 是基于函数的工厂注册表
type FuncRegistry struct {
	factories map[NotificationType]NotificationFunc
}

// NewFuncRegistry 创建函数式工厂注册表
func NewFuncRegistry() *FuncRegistry {
	return &FuncRegistry{
		factories: make(map[NotificationType]NotificationFunc),
	}
}

// Register 注册工厂函数
func (r *FuncRegistry) Register(notifyType NotificationType, fn NotificationFunc) {
	r.factories[notifyType] = fn
}

// Create 使用工厂函数创建实例
func (r *FuncRegistry) Create(notifyType NotificationType, config map[string]string) (Notification, error) {
	fn, ok := r.factories[notifyType]
	if !ok {
		return nil, fmt.Errorf("no factory function for type: %s", notifyType)
	}
	return fn(config)
}

// ============================================================================
// 预配置的工厂函数（闭包应用）
// ============================================================================

// CreateConfiguredSMSFactory 返回一个预配置了默认参数的 SMS 工厂函数。
//
// 这是闭包的经典应用，工厂函数"记住"了默认配置。
//
// 使用场景：
//   - 大部分 SMS 使用相同的 provider
//   - 只需要在特定情况下覆盖某些参数
//
// 使用示例：
//
//	factory := CreateConfiguredSMSFactory("twilio", "default-api-key")
//	sms1, _ := factory(map[string]string{})  // 使用默认配置
//	sms2, _ := factory(map[string]string{"apiKey": "special-key"})  // 覆盖 apiKey
func CreateConfiguredSMSFactory(defaultProvider, defaultAPIKey string) NotificationFunc {
	return func(config map[string]string) (Notification, error) {
		// 使用传入的配置，如果不存在则使用默认值
		provider := config["provider"]
		if provider == "" {
			provider = defaultProvider
		}

		apiKey := config["apiKey"]
		if apiKey == "" {
			apiKey = defaultAPIKey
		}

		return &SMSNotification{
			provider: provider,
			apiKey:   apiKey,
		}, nil
	}
}

// ============================================================================
// 抽象工厂模式（Abstract Factory）- 产品族
// ============================================================================

// NotificationKit 是通知工具包接口（抽象工厂）。
//
// 当需要创建一组相关的产品时使用，如：
//   - 开发环境工具包（所有通知都打印到控制台）
//   - 生产环境工具包（所有通知都调用真实服务）
type NotificationKit interface {
	CreateSMS() Notification
	CreateEmail() Notification
	CreatePush() Notification
}

// ProductionKit 是生产环境的通知工具包
type ProductionKit struct {
	smsConfig   map[string]string
	emailConfig map[string]string
	pushConfig  map[string]string
}

// NewProductionKit 创建生产环境工具包
func NewProductionKit(sms, email, push map[string]string) *ProductionKit {
	return &ProductionKit{
		smsConfig:   sms,
		emailConfig: email,
		pushConfig:  push,
	}
}

// CreateSMS 创建生产环境的 SMS 通知
func (k *ProductionKit) CreateSMS() Notification {
	sms, _ := CreateNotification(TypeSMS, k.smsConfig)
	return sms
}

// CreateEmail 创建生产环境的 Email 通知
func (k *ProductionKit) CreateEmail() Notification {
	email, _ := CreateNotification(TypeEmail, k.emailConfig)
	return email
}

// CreatePush 创建生产环境的 Push 通知
func (k *ProductionKit) CreatePush() Notification {
	push, _ := CreateNotification(TypePush, k.pushConfig)
	return push
}

// MockKit 是测试环境的通知工具包（Mock 实现）
type MockKit struct{}

// CreateSMS 创建 Mock SMS 通知
func (k *MockKit) CreateSMS() Notification {
	return &MockNotification{channel: "sms"}
}

// CreateEmail 创建 Mock Email 通知
func (k *MockKit) CreateEmail() Notification {
	return &MockNotification{channel: "email"}
}

// CreatePush 创建 Mock Push 通知
func (k *MockKit) CreatePush() Notification {
	return &MockNotification{channel: "push"}
}

// MockNotification 是用于测试的 Mock 实现
type MockNotification struct {
	channel      string
	lastReceiver string
	lastMessage  string
}

// Send 记录调用参数但不实际发送
func (m *MockNotification) Send(receiver, message string) error {
	m.lastReceiver = receiver
	m.lastMessage = message
	return nil
}

// Channel 返回渠道名称
func (m *MockNotification) Channel() string {
	return m.channel
}

// GetLastCall 返回最后一次调用的参数（测试辅助方法）
func (m *MockNotification) GetLastCall() (receiver, message string) {
	return m.lastReceiver, m.lastMessage
}
