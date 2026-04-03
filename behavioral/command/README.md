# 命令模式 (Command Pattern)

## 概述

命令模式把请求封装成独立对象，从而实现请求发送者与请求执行者解耦。它常用于：

- UI 事件绑定
- 任务调度
- 撤销重做
- 队列化执行

## Go 语言实现特点

- 传统写法可以使用接口
- 更轻量的写法可以直接使用 `func()`
- 函数闭包能自然捕获接收者与参数

## 传统接口实现

```go
type Command interface {
    Execute() string
}
```

这种方式更显式，适合：

- 命令需要命名类型
- 命令对象有额外状态
- 命令体系较稳定

## Go 中更轻量的实现：func 作为命令

```go
type FuncCommand func() string

func (c FuncCommand) Execute() string {
    return c()
}
```

这是 Go 里很自然的做法，因为函数本身就是一等公民。

优势在于：

- 样板代码更少
- 容易就地定义命令
- 闭包可以直接捕获上下文

## 代码示例

当前实现同时展示两种方式：

1. 传统接口式命令：`TurnOnCommand`、`TurnOffCommand`
2. 函数式命令：`FuncCommand`

还额外提供了一个 `Scheduler`，演示如何把函数命令排队执行。

## 使用场景

1. 事件处理
2. 任务队列
3. 批处理
4. 撤销/重做

## 优缺点

### 优点

- 调用者和执行者解耦
- 易于排队、记录、重放
- Go 中函数式命令非常简洁

### 缺点

- 过度对象化会增加样板代码
- 简单场景下可能不如直接函数调用直观

## 参考

- 实现代码：[`command.go`](./command.go)
- 测试代码：[`command_test.go`](./command_test.go)
