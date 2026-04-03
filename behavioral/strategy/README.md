# 策略模式 (Strategy Pattern)

## 概述

策略模式定义了一系列算法，并将每个算法封装起来，使它们可以互相替换。策略模式让算法的变化独立于使用算法的客户。

Go 语言通过 **接口（interface）** 和 **函数类型（func type）** 提供了比传统面向对象语言更轻量、更灵活的策略模式实现方式。

## 实现方式

### 1. 接口策略 (Interface Strategy)

使用 Go 接口定义策略契约，不同的策略实现相同的接口。

```go
// PaymentStrategy 支付策略接口
type PaymentStrategy interface {
    Pay(amount float64) (string, error)
    GetName() string
}

// 具体策略实现
- AlipayPayment    // 支付宝支付
- WeChatPayment    // 微信支付
- CreditCardPayment // 信用卡支付
```

**适用场景**：策略需要维护状态，或策略之间有较大的行为差异。

### 2. 函数类型策略 (Function Type Strategy)

使用 Go 的函数类型直接定义策略，配合闭包实现轻量级策略。

```go
// SortStrategy 排序策略 - 函数类型
type SortStrategy func([]int) []int

// 具体策略实现
- BubbleSortStrategy    // 冒泡排序
- QuickSortStrategy     // 快速排序
- MergeSortStrategy     // 归并排序
```

**适用场景**：策略是无状态的纯函数，策略逻辑相对简单。

### 3. 闭包策略 (Closure Strategy)

使用闭包捕获外部变量，实现参数化的策略。

```go
// DiscountStrategy 折扣策略 - 闭包函数类型
type DiscountStrategy func(float64) float64

// 工厂函数创建参数化策略
func NewPercentageDiscount(percentage float64) DiscountStrategy {
    return func(price float64) float64 {
        return price * (1 - percentage/100)
    }
}
```

**适用场景**：策略需要参数化配置，运行时动态生成策略。

## 架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                     Context (上下文)                         │
│                    - ShoppingCart                            │
│                    - Sorter                                  │
│                    - PriceCalculator                         │
└──────────────────────┬──────────────────────────────────────┘
                       │ 使用
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                   Strategy (策略接口/类型)                    │
│         ┌─────────────────┬──────────────────┐              │
│         │  Interface      │  Function Type   │              │
│         │  Strategy       │  Strategy        │              │
│         └────────┬────────┴────────┬─────────┘              │
│                  │                 │                        │
│         ┌────────▼────────┐ ┌──────▼───────┐                │
│         │ Concrete        │ │ Concrete     │                │
│         │ Strategies      │ │ Strategies   │                │
│         │ (Alipay,        │ │ (BubbleSort, │                │
│         │  WeChat, etc)   │ │  QuickSort)  │                │
│         └─────────────────┘ └──────────────┘                │
└─────────────────────────────────────────────────────────────┘
```

## 核心代码解析

### 接口策略实现

```go
// PaymentStrategy 定义支付策略接口
type PaymentStrategy interface {
    Pay(amount float64) (string, error)
    GetName() string
}

// AlipayPayment 支付宝支付策略
type AlipayPayment struct {
    account string
}

func (a *AlipayPayment) Pay(amount float64) (string, error) {
    return fmt.Sprintf("支付宝支付 ¥%.2f，账号: %s", amount, a.account), nil
}
```

### 函数类型策略实现

```go
// SortStrategy 排序策略函数类型
type SortStrategy func([]int) []int

// QuickSortStrategy 快速排序策略
var QuickSortStrategy SortStrategy = func(data []int) []int {
    if len(data) <= 1 {
        return data
    }
    // 快速排序实现...
}
```

### 闭包策略实现

```go
// DiscountStrategy 折扣策略函数类型
type DiscountStrategy func(float64) float64

// NewPercentageDiscount 创建百分比折扣策略
func NewPercentageDiscount(percentage float64) DiscountStrategy {
    return func(price float64) float64 {
        return price * (1 - percentage/100)
    }
}

// 使用
strategy := NewPercentageDiscount(20) // 20% 折扣
finalPrice := strategy(100)           // 返回 80
```

## 策略注册表模式

```go
// StrategyRegistry 策略注册表 - 支持运行时策略选择
type StrategyRegistry struct {
    strategies map[string]PaymentStrategy
}

func (sr *StrategyRegistry) Register(name string, strategy PaymentStrategy) {
    sr.strategies[name] = strategy
}

func (sr *StrategyRegistry) Get(name string) (PaymentStrategy, error) {
    // 返回对应策略...
}
```

## 使用示例

### 支付系统

```go
// 创建购物车
cart := NewShoppingCart()

// 设置支付宝策略
cart.SetPaymentStrategy(NewAlipayPayment("user@example.com"))
result, _ := cart.Checkout(100.0)
// 输出: 支付宝支付 ¥100.00，账号: user@example.com

// 动态切换为微信支付
cart.SetPaymentStrategy(NewWeChatPayment("wx_user"))
result, _ = cart.Checkout(100.0)
// 输出: 微信支付 ¥100.00，OpenID: wx_user
```

### 排序系统

```go
sorter := NewSorter()

// 使用快速排序
data := []int{5, 2, 8, 1, 9}
sorter.SetStrategy(QuickSortStrategy)
sorted := sorter.Sort(data) // [1, 2, 5, 8, 9]

// 切换为归并排序
sorter.SetStrategy(MergeSortStrategy)
sorted = sorter.Sort(data) // [1, 2, 5, 8, 9]
```

### 价格计算

```go
// 创建固定金额折扣
calculator := NewPriceCalculator(NewFixedDiscount(50))
price := calculator.Calculate(200) // 150

// 切换为百分比折扣
calculator.SetStrategy(NewPercentageDiscount(20))
price = calculator.Calculate(200) // 160

// 切换为满减策略
calculator.SetStrategy(NewThresholdDiscount(100, 30))
price = calculator.Calculate(200) // 170 (满100减30)
```

## 优势

1. **开闭原则**：新增策略无需修改现有代码
2. **消除条件语句**：避免大量的 if-else 或 switch-case
3. **运行时切换**：可以在程序运行时动态改变策略
4. **Go 语言特色**：
   - 函数类型让简单策略更轻量
   - 闭包实现参数化策略
   - 隐式接口实现降低耦合

## 适用场景

- 多种算法或行为需要互相替换
- 算法需要独立于使用它的客户端
- 需要消除大量的条件判断语句
- 需要对客户端隐藏复杂的算法实现

## 对比：Interface vs Function Type

| 特性 | Interface Strategy | Function Type Strategy |
|------|-------------------|----------------------|
| 状态维护 | ✅ 可以维护状态 | ❌ 无状态（除非使用闭包） |
| 复杂度 | 适合复杂策略 | 适合简单策略 |
| 代码量 | 需要定义结构体 | 直接定义函数 |
| 扩展性 | 易于添加方法 | 需要修改函数签名 |
| 适用场景 | 复杂业务逻辑 | 纯函数算法 |

## 参考

- [Go 接口最佳实践](https://golang.org/doc/effective_go.html#interfaces)
- [Go 函数类型](https://golang.org/ref/spec#Function_types)
- [设计模式：可复用面向对象软件的基础](https://en.wikipedia.org/wiki/Design_Patterns)
