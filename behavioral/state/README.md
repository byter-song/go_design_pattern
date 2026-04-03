# 状态模式 (State Pattern)

## 概述

状态模式把对象在不同状态下的行为拆分到独立状态对象中，让上下文对象不必维护庞大的状态分支逻辑。

## Go 语言实现特点

- 使用接口描述“状态能做什么”
- 每个状态单独实现自己的行为和状态迁移
- 上下文只保存当前状态，不再堆满 switch-case

## 核心概念

```go
type OrderState interface {
    Pay(*Order) error
    Ship(*Order) error
    Complete(*Order) error
}
```

不同状态分别实现该接口，并在合适时机切换到下一个状态。

## 为什么它能消除复杂 switch-case

没有状态模式时，代码常常像这样：

```go
switch order.status {
case "pending":
case "paid":
case "shipped":
}
```

一旦状态变多、操作变多，这些分支会迅速膨胀。

使用状态模式后：

- `pendingState` 负责未支付逻辑
- `paidState` 负责已支付逻辑
- `shippedState` 负责已发货逻辑

行为分散回各状态对象中，代码会更清晰。

## 代码示例

当前实现是一个订单状态机：

`pending -> paid -> shipped -> completed`

非法状态迁移会直接返回错误，例如：

- 未支付不能发货
- 已支付不能再次支付
- 已完成不能再次操作

## 使用场景

1. 状态机复杂且分支众多
2. 同一操作在不同状态下行为差异明显
3. 希望让状态迁移规则可维护、可扩展

## 优缺点

### 优点

- 消除大型条件分支
- 状态职责清晰
- 扩展新状态更自然

### 缺点

- 状态类型增多
- 简单状态机用它可能显得过度设计

## 参考

- 实现代码：[`state.go`](./state.go)
- 测试代码：[`state_test.go`](./state_test.go)
