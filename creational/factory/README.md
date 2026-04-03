# 工厂方法模式 (Factory Method Pattern)

## 模式定义

工厂方法模式定义了一个创建对象的接口，但由子类决定实例化哪个类。工厂方法让类的实例化推迟到子类。

> **核心思想**：将对象的创建逻辑与使用逻辑分离，通过接口编程而非具体类型编程。

---

## 适用场景

### 应该使用工厂模式的场景

1. **对象创建逻辑复杂**
   - 需要根据不同条件创建不同对象
   - 创建过程涉及多个步骤或依赖其他对象

2. **需要解耦对象创建与使用**
   - 调用者不需要知道具体类名
   - 便于替换实现（如切换数据库、消息队列）

3. **需要统一管理对象创建**
   - 连接池管理
   - 资源限制（如最大实例数）
   - 创建日志记录

4. **框架和库设计**
   - 提供扩展点，让用户提供具体实现
   - 如 database/sql 的 Driver 接口

### 具体示例

```go
// 根据配置创建不同的数据库连接
db, err := factory.CreateDatabase("postgres", config)
db, err := factory.CreateDatabase("mysql", config)

// 根据用户偏好创建不同的通知方式
notifier := factory.CreateNotifier(user.Preference)
```

---

## Go 语言实现

Go 语言实现工厂模式有三种主要形式，从简单到复杂：

### 1. 简单工厂（Simple Factory）- 最常用

```go
// 根据类型创建对应的通知实例
func CreateNotification(notifyType NotificationType, config map[string]string) (Notification, error) {
    switch notifyType {
    case TypeSMS:
        return &SMSNotification{...}, nil
    case TypeEmail:
        return &EmailNotification{...}, nil
    // ...
    }
}
```

**适用场景**：
- 类型数量不多（< 10 个）
- 创建逻辑相对简单
- 不需要运行时动态注册新类型

**优缺点**：
- ✅ 最简单直接
- ✅ 集中管理创建逻辑
- ❌ 增加新类型需要修改函数（违反开闭原则）

### 2. 工厂方法（Factory Method）- 更灵活

```go
// 工厂接口
type NotificationFactory interface {
    Create(config map[string]string) (Notification, error)
}

// 具体工厂
type SMSFactory struct{}
func (f *SMSFactory) Create(config map[string]string) (Notification, error) {
    return &SMSNotification{...}, nil
}

// 工厂注册表
type FactoryRegistry struct {
    factories map[NotificationType]NotificationFactory
}
```

**适用场景**：
- 需要不同的创建策略（Mock、缓存、池化）
- 需要运行时动态注册工厂
- 创建逻辑复杂，需要独立管理

**优缺点**：
- ✅ 符合开闭原则
- ✅ 支持运行时扩展
- ❌ 代码量稍多

### 3. 函数式工厂 - Go 惯用法

Go 中可以用函数类型替代接口，更简洁：

```go
// 函数类型定义
type NotificationFunc func(config map[string]string) (Notification, error)

// 注册和使用
type FuncRegistry struct {
    factories map[NotificationType]NotificationFunc
}

// 注册时直接传入函数
registry.Register(TypeSMS, func(config map[string]string) (Notification, error) {
    return &SMSNotification{...}, nil
})
```

**适用场景**：
- 创建逻辑简单，不需要状态
- 追求代码简洁
- 配合闭包实现配置化

**优缺点**：
- ✅ 最简洁的 Go 风格
- ✅ 配合闭包非常灵活
- ❌ 不适合复杂的工厂逻辑

### 4. 抽象工厂（Abstract Factory）- 产品族

当需要创建一组相关的产品时使用：

```go
// 抽象工厂接口
type NotificationKit interface {
    CreateSMS() Notification
    CreateEmail() Notification
    CreatePush() Notification
}

// 生产环境实现
type ProductionKit struct{...}

// Mock 测试实现
type MockKit struct{...}
```

**适用场景**：
- 需要保证产品之间的兼容性
- 不同环境需要不同的产品族（开发/测试/生产）
- 跨平台支持（Windows/Mac/Linux）

---

## Go 语言特殊点

### 1. 隐式接口 = 松耦合

```go
// 接口定义在工厂包中
type Notification interface {
    Send(receiver, message string) error
}

// 具体类型在另一个包中，甚至可以是第三方库
// 只要实现了 Send 方法，就能被工厂使用
```

不需要显式声明 `implements`，这带来了极大的灵活性。

### 2. 函数是一等公民

```go
// 函数类型作为工厂
type NotificationFunc func(config map[string]string) (Notification, error)

// 闭包实现配置化工厂
func CreateConfiguredFactory(defaultConfig Config) NotificationFunc {
    return func(override Config) (Notification, error) {
        // 合并 defaultConfig 和 override
        // ...
    }
}
```

这是 Go 相比 Java/C++ 的独特优势。

### 3. 小写 = 封装

```go
// 小写开头 = 包私有，强制通过工厂创建
type smsNotification struct{...}

// 大写开头 = 公开 API
func CreateSMS(config Config) Notification {...}
```

通过命名约定实现封装，不需要 `private`/`public` 关键字。

### 4. map + 函数 = 注册表模式

```go
var factories = map[string]func() Notification{
    "sms":   func() Notification { return &SMSNotification{} },
    "email": func() Notification { return &EmailNotification{} },
}

func Create(t string) Notification {
    if f, ok := factories[t]; ok {
        return f()
    }
    return nil
}
```

这是 Go 中非常常见的模式。

---

## 优缺点分析

### 优点

1. **解耦**
   - 调用者不需要知道具体类名
   - 便于替换实现

2. **集中管理**
   - 创建逻辑统一维护
   - 便于添加日志、监控、限制

3. **符合开闭原则**
   - 新增产品类型无需修改现有代码（工厂方法模式）

4. **便于测试**
   - 可以轻松注入 Mock 实现

### 缺点

1. **增加复杂度**
   - 需要额外的工厂类/函数
   - 代码量增加

2. **简单工厂违反开闭原则**
   - 新增类型需要修改 switch 语句

3. **过度设计风险**
   - 简单场景不需要工厂模式
   - 直接 `&Struct{}` 更清晰

---

## 最佳实践

### ✅ 推荐做法

1. **从简单工厂开始**
   ```go
   // 先写简单工厂，需要时再重构
   func Create(t Type) (Interface, error) {
       switch t { ... }
   }
   ```

2. **利用 Go 的函数类型**
   ```go
   // 优先使用函数而非接口
   type Factory func(Config) (Product, error)
   ```

3. **返回接口而非具体类型**
   ```go
   // 好
   func Create() (Notification, error)
   
   // 不好
   func Create() (*SMSNotification, error)
   ```

4. **配合错误处理**
   ```go
   func Create(t Type) (Product, error) {
       switch t {
       case TypeA:
           return &ProductA{}, nil
       default:
           return nil, fmt.Errorf("unknown type: %s", t)
       }
   }
   ```

### ❌ 避免的做法

1. **不要为简单对象使用工厂**
   ```go
   // 过度设计
   user := factory.CreateUser()
   
   // 直接创建更清晰
   user := &User{}
   ```

2. **避免工厂嵌套工厂**
   ```go
   // 难以理解
   factory := factoryFactory.CreateFactory()
   product := factory.Create()
   ```

3. **不要忽视零值**
   ```go
   // 很多时候零值就够用了
   var cfg Config  // 零值就是有效配置
   ```

---

## 与其他模式的关系

| 模式 | 关系 | 说明 |
|------|------|------|
| **单例** | 配合使用 | 工厂本身可以是单例 |
| **建造者** | 配合使用 | 复杂对象先用工厂选择类型，再用建造者构建 |
| **依赖注入** | 替代方案 | DI 容器可以替代工厂模式 |
| **策略** | 相似 | 工厂创建对象，策略切换算法 |

---

## 实际案例

### database/sql 包

```go
// 标准库的工厂模式
import "database/sql"
import _ "github.com/lib/pq"  // 注册 postgres driver

db, err := sql.Open("postgres", dsn)  // 工厂方法
```

Driver 通过 `init()` 函数自注册：

```go
func init() {
    sql.Register("postgres", &Driver{})
}
```

### net/http 的 Handler

```go
// Handler 是接口
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}

// HandlerFunc 是函数类型工厂
type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}
```

---

## 总结

工厂模式在 Go 中有多种实现形式：

1. **简单工厂**：最常用，适合大多数场景
2. **工厂方法**：需要运行时扩展时使用
3. **函数式工厂**：最 Go 风格，简洁灵活
4. **抽象工厂**：处理产品族时使用

> 💡 **核心建议**：从简单工厂开始，只有当确实需要更灵活的结构时，再重构为更复杂的模式。Go 的隐式接口和函数类型让工厂模式实现得非常优雅。
