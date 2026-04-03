# 模板方法模式 (Template Method Pattern)

## 概述

模板方法模式定义一个算法骨架，把其中可变的步骤留给具体实现。传统面向对象语言通常依赖“父类模板方法 + 子类重写钩子方法”。

Go 没有继承，所以实现模板方法时通常会转向更符合语言风格的两种写法：

- 组合 + 接口
- 传递闭包函数

## Go 语言实现特点

- 不依赖继承
- 用接口表达“步骤”
- 用结构体表达“骨架”
- 用闭包快速注入变化步骤

## 方案一：组合 + 接口

```go
type ReportSteps interface {
    LoadData() string
    Analyze(string) string
    Format(string) string
}

type ReportGenerator struct {
    steps ReportSteps
}
```

`ReportGenerator.Generate()` 固定执行顺序：

1. LoadData
2. Analyze
3. Format

## 方案二：传递闭包函数

```go
type Pipeline struct {
    Load   func() string
    Handle func(string) string
    Save   func(string) string
}
```

这种方式更轻量，特别适合：

- 一次性流程
- 不值得专门定义多个类型的小场景
- 想快速表达“骨架 + 可变步骤”

## 为什么 Go 不直接照搬传统模板方法

因为 Go 更鼓励：

- 组合而不是继承
- 小接口而不是深层类层次
- 函数作为一等公民

所以在 Go 里，模板方法最地道的写法往往不是“子类覆盖父类方法”，而是“骨架对象调用被注入的步骤实现”。

## 使用场景

1. 多个流程整体相同，但局部步骤不同
2. 需要稳定的处理顺序
3. 希望把变化点单独隔离

## 优缺点

### 优点

- 骨架稳定，变化点清晰
- 很适合提炼重复流程
- 闭包版实现非常轻量

### 缺点

- 如果步骤很多，接口会变大
- 闭包版虽然轻量，但约束不如类型系统直观

## 参考

- 实现代码：[`template_method.go`](./template_method.go)
- 测试代码：[`template_method_test.go`](./template_method_test.go)
