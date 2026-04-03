# 享元模式 (Flyweight Pattern)

## 概述

享元模式通过共享内部状态来降低内存占用和对象创建成本。经典定义强调把“不会频繁变化的内部状态”抽离出来共享，而把“每次使用都不同的外部状态”留给调用方传入。

在 Go 中，除了经典的缓存工厂写法，`sync.Pool` 也是非常地道的享元实践：它适合复用生命周期短、创建频繁的对象，减少临时分配压力。

## Go 语言实现特点

- 使用 `map` 缓存共享对象
- 使用 `sync.RWMutex` 保证并发安全
- 使用 `sync.Pool` 复用临时对象
- 显式区分内部状态与外部状态

## 核心概念

```go
type BulletStyle struct {
    Kind   string
    Color  string
    Sprite string
}

type Bullet struct {
    Style *BulletStyle
    X, Y  float64
}
```

- `BulletStyle` 是共享的内部状态
- `Bullet` 是带外部状态的运行时对象

## 代码示例

当前实现模拟游戏中的子弹系统：

### 1. 共享样式

```go
factory := NewStyleFactory()
laser1 := factory.GetStyle("laser", "green", ">>>")
laser2 := factory.GetStyle("laser", "green", ">>>")

fmt.Println(laser1 == laser2) // true
```

相同样式只创建一次，后续直接复用。

### 2. 使用 sync.Pool 复用对象

```go
pool := NewBulletPool()
style := factory.GetStyle("laser", "green", ">>>")

bullet := pool.Acquire(style, "player-1", 10, 5, 20)
bullet.Move(0.5)
pool.Release(bullet)
```

这里的 `BulletPool` 基于 `sync.Pool` 实现，是 Go 中非常常见的对象复用方式。

## 为什么这里使用 sync.Pool

在很多 Go 项目里，享元模式不会严格照搬面向对象教材里的写法，而是直接落在这些实践上：

1. 缓存共享的不可变对象
2. 用 `sync.Pool` 复用高频临时对象

这正符合 Go 的工程风格：更关注分配成本、GC 压力和并发下的实际收益。

## 使用场景

1. 日志缓冲区、序列化缓冲区复用
2. 游戏对象、粒子对象的短生命周期复用
3. 大量格式相同、状态少量变化的对象
4. 需要降低 GC 压力的热点路径

## 优缺点

### 优点

- 降低内存占用
- 减少重复分配
- 适合热点路径优化
- 共享状态更易统一管理

### 缺点

- 需要清晰划分内部状态与外部状态
- `sync.Pool` 不保证对象一定被复用
- 过早优化会增加理解成本

## 与其他模式的关系

- 与 **单例模式** 不同：享元共享的是大量细粒度对象，不是单个全局实例
- 与 **对象池模式** 接近：Go 中常常与 `sync.Pool` 配合使用
- 与 **工厂模式** 常一起出现：工厂负责返回共享享元

## 参考

- 实现代码：[`flyweight.go`](./flyweight.go)
- 测试代码：[`flyweight_test.go`](./flyweight_test.go)
