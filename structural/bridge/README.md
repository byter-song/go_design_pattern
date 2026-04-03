# 桥接模式 (Bridge Pattern)

## 概述

桥接模式把抽象部分与实现部分拆开，让两者可以独立变化。客户端依赖抽象，抽象再委托给实现接口完成实际工作。

在 Go 中，这种模式非常自然：结构体持有接口就是抽象与实现解耦的直接方式；再结合接口嵌入，可以把设备能力拆成多个小接口并按需组合。

## Go 语言实现特点

- 抽象层持有接口，而不是具体实现
- 通过接口嵌入组合不同能力边界
- 可以独立扩展新的抽象与新的实现
- 不需要复杂继承体系也能表达“抽象 × 实现”的正交扩展

## 核心概念

```go
type Device interface {
    Enable()
    Disable()
    IsEnabled() bool
}

type EntertainmentDevice interface {
    Device
    VolumeControl
    ChannelControl
}

type BasicRemote struct {
    device Device
}
```

## 代码示例

当前实现使用“遥控器 × 设备”来说明桥接模式：

### 实现侧

- `TV`
- `StreamingBox`

二者都实现 `EntertainmentDevice`，但各自有不同的默认值与约束。

### 抽象侧

- `BasicRemote`：只负责开关机
- `SmartRemote`：在基础上增加音量与频道控制

示例：

```go
tv := NewTV("Bedroom TV")
remote := NewSmartRemote(tv)

remote.TogglePower()
remote.VolumeUp(15)
remote.ChannelNext()
```

## 为什么这里强调接口组合

桥接模式的关键是“把能力边界拆开”。在 Go 中，最优雅的方式往往不是继承，而是：

```go
type EntertainmentDevice interface {
    Device
    VolumeControl
    ChannelControl
}
```

这样一来：

- 基础抽象只依赖 `Device`
- 高级抽象依赖 `EntertainmentDevice`
- 新设备只需实现对应能力集合即可接入

## 使用场景

1. 抽象与实现都可能独立扩展
2. 不希望产生子类爆炸
3. 希望不同抽象共享同一组实现
4. 需要通过接口组合定义不同能力等级

## 优缺点

### 优点

- 解耦抽象与实现
- 避免继承层次膨胀
- 更适合 Go 的组合风格
- 方便测试和替换实现

### 缺点

- 增加了一层间接调用
- 如果接口设计过细，理解成本会提高

## 与其他模式的关系

- 与 **适配器模式** 不同：桥接是事先设计好的解耦结构，适配器是事后兼容已有接口
- 与 **策略模式** 相似：都依赖接口；但桥接强调“抽象与实现的两个维度独立演化”
- 与 **外观模式** 不同：桥接解耦结构，外观简化入口

## 参考

- 实现代码：[`bridge.go`](./bridge.go)
- 测试代码：[`bridge_test.go`](./bridge_test.go)
