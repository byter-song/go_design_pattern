# 解释器模式 (Interpreter Pattern)

## 概述

解释器模式为某种简单语言定义语法结构，并提供解释执行逻辑。它常用于规则表达式、筛选条件、权限判断等小型 DSL。

## Go 语言实现特点

- 用接口表示表达式节点
- 用组合表达 AST
- 通过上下文传递变量环境

## 代码示例

当前实现支持布尔表达式：

- 变量：`is_admin`
- 与：`And`
- 或：`Or`
- 非：`Not`

例如：

```go
rule := And(a, And(b, Not(c)))
```

可以表达：

`is_admin AND has_token AND NOT is_banned`

## 使用场景

1. 权限规则
2. 过滤表达式
3. 简单 DSL

## 参考

- 实现代码：[`interpreter.go`](./interpreter.go)
- 测试代码：[`interpreter_test.go`](./interpreter_test.go)
