# 访问者模式 (Visitor Pattern)

## 概述

访问者模式把作用于对象结构上的操作抽离出来，使新增操作时不必修改元素类型本身。

## Go 语言实现特点

- 通过接口定义访问者
- 元素通过 `Accept` 分发给访问者
- 适合“对象结构稳定、操作种类经常增加”的场景

## 代码示例

当前实现包含：

- `File`
- `Directory`
- `SizeVisitor`
- `NameVisitor`

分别演示统计总大小与收集名称两类操作。

## 使用场景

1. 编译器 AST
2. 文件树分析
3. 报表导出

## 参考

- 实现代码：[`visitor.go`](./visitor.go)
- 测试代码：[`visitor_test.go`](./visitor_test.go)
