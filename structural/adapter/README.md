# 适配器模式 (Adapter Pattern)

## 概述

适配器模式允许将一个类的接口转换成客户希望的另外一个接口。适配器模式使得原本由于接口不兼容而不能一起工作的那些类可以一起工作。

## Go 语言实现特点

Go 语言的**隐式接口（Implicit Interface）**是适配器模式最强大的武器：

1. **无需显式声明**：类型不需要声明它实现了哪个接口
2. **解耦合**：接口定义和使用可以分离
3. **适配灵活**：可以在不修改原有代码的情况下进行适配

### 核心概念

```go
// 目标接口
type Target interface {
    Request()
}

// 被适配者（已有代码，无法修改）
type Adaptee struct{}
func (a *Adaptee) SpecificRequest() {}

// 适配器
type Adapter struct {
    adaptee *Adaptee
}

func (a *Adapter) Request() {
    // 转换调用
    a.adaptee.SpecificRequest()
}
```

## 代码示例

### 对象适配器（组合方式）

```go
// 目标接口
type PaymentProcessor interface {
    ProcessPayment(amount float64, currency string) (string, error)
}

// 旧系统（需要被适配）
type LegacyPaymentSystem struct{}
func (l *LegacyPaymentSystem) MakeCharge(amountInCents int64, currencyCode string) (int64, error) {
    // ...
}

// 适配器
type LegacyPaymentAdapter struct {
    legacy *LegacyPaymentSystem
}

func (a *LegacyPaymentAdapter) ProcessPayment(amount float64, currency string) (string, error) {
    // 参数转换：元 -> 分
    amountInCents := int64(amount * 100)
    txID, err := a.legacy.MakeCharge(amountInCents, currency)
    return fmt.Sprintf("TX-%d", txID), err
}
```

### 隐式接口适配

Go 的独特优势：

```go
// 缓存接口
type Cache interface {
    Get(key string) (interface{}, bool)
    Set(key string, value interface{}) error
}

// 外部缓存库（方法名不同）
type ExternalCache struct{}
func (e *ExternalCache) Fetch(key string) (interface{}, bool)
func (e *ExternalCache) Store(key string, value interface{})

// 适配器 - 利用隐式接口
type ExternalCacheAdapter struct {
    cache *ExternalCache
}

func (a *ExternalCacheAdapter) Get(key string) (interface{}, bool) {
    return a.cache.Fetch(key)  // 方法名映射
}

func (a *ExternalCacheAdapter) Set(key string, value interface{}) error {
    a.cache.Store(key, value)
    return nil
}
```

### 匿名适配（Go 特有技巧）

```go
// 使用闭包和匿名结构体快速适配
func AdaptExternalCache(external *ExternalCache) Cache {
    return &struct {
        Get func(string) (interface{}, bool)
        Set func(string, interface{}) error
        Del func(string)
    }{
        Get: external.Fetch,
        Set: func(k string, v interface{}) error {
            external.Store(k, v)
            return nil
        },
        Del: external.Remove,
    }
}
```

### 双向适配器

```go
// 欧洲插头接口
type EuropeanPlug interface {
    ConnectToEuropeanSocket() string
}

// 美国插头接口
type USPlug interface {
    ConnectToUSSocket() string
}

// 欧洲 -> 美国适配器
type EuropeanToUSAdapter struct {
    device EuropeanPlug
}

func (a *EuropeanToUSAdapter) ConnectToUSSocket() string {
    return "Adapter: 110V -> 220V -> " + a.device.ConnectToEuropeanSocket()
}

// 美国 -> 欧洲适配器
type USToEuropeanAdapter struct {
    device USPlug
}

func (a *USToEuropeanAdapter) ConnectToEuropeanSocket() string {
    return "Adapter: 220V -> 110V -> " + a.device.ConnectToUSSocket()
}
```

## 使用场景

1. **第三方库集成**：适配第三方库的接口到项目标准接口
2. **遗留系统改造**：将旧系统接口适配到新系统
3. **多版本兼容**：适配不同版本的 API
4. **跨平台适配**：适配不同平台的实现
5. **测试 Mock**：用适配器创建测试替身

## 优缺点

### 优点

- **复用性**：复用已有的类，无需修改
- **灵活性**：通过适配器灵活转换接口
- **解耦合**：客户端与具体实现解耦
- **Go 特色**：隐式接口让适配更加自然

### 缺点

- **代码量增加**：需要编写适配器类
- **调用链变长**：增加了一层间接调用

## Go 隐式接口的优势

```go
// 标准库 io.Reader 接口
type Reader interface {
    Read(p []byte) (n int, err error)
}

// 任何实现了 Read 方法的类型都自动实现了 io.Reader
// 无需显式声明：type MyReader struct{} implements io.Reader

// 这意味着：
// 1. 可以在不修改原有代码的情况下适配
// 2. 接口定义可以放在使用方
// 3. 更容易进行依赖注入和测试
```

## 与其他模式的关系

- **桥接模式**：适配器改变接口，桥接模式分离抽象和实现
- **装饰器模式**：适配器改变接口，装饰器保持接口并添加功能
- **代理模式**：适配器改变接口，代理模式保持接口并控制访问

## 参考

- [Go 接口详解](https://golang.org/doc/effective_go#interfaces)
- [Go 隐式接口的威力](https://www.alexedwards.net/blog/interfaces-explained)
