# 组合模式 (Composite Pattern)

## 概述

组合模式通过统一接口将对象组织成树形结构，使客户端可以一致地处理单个对象和对象集合。它特别适用于具有层次结构的数据，如文件系统、组织架构、UI 组件树等。

## Go 语言实现特点

- 接口定义统一行为，叶子和容器都实现该接口
- 切片天然适合表示子节点集合
- 通过递归轻松实现聚合操作（求和、计数、渲染）
- 零值切片即可表示“没有子节点”，无需额外初始化

## 核心概念

```go
type Node interface {
    Name() string
    Size() int
    Count() int
    Render(indent int) string
}

type File struct{}      // 叶子
type Directory struct{} // 组合（容器）
```

## 代码示例

当前实现模拟一个简单的文件系统：

- `File`：叶子节点，包含名称和大小
- `Directory`：组合节点，持有 `[]Node` 子节点
- 支持 `Add/Remove/Find/Walk` 等常用操作

渲染树形结构：

```go
output := root.Render(0)
// 示例片段：
// + root/
//   + docs/
//     - README.md (12KB)
//     - spec.md (20KB)
//   + images/
//     - logo.png (256KB)
//   - go.mod (1KB)
```

## 使用场景

1. 树形结构数据（文件、菜单、UI、组织）
2. 需要统一处理叶子与容器
3. 需要聚合计算（总大小、节点数量）
4. 需要支持树的遍历与查找

## 优缺点

### 优点

- 一致性 API，简化客户端逻辑
- 容易扩展新的节点类型
- 聚合操作自然、直观

### 缺点

- 容器与叶子间的行为差异可能导致接口臃肿
- 某些操作对叶子无意义（如 Add），需谨慎设计

## 与其他模式的关系

- 与 **装饰器模式** 可组合使用：在树上对节点添加增强行为
- 与 **迭代器模式** 结合：为树提供统一遍历接口
- 与 **责任链模式** 结合：在树上向上或向下传递请求

## 参考

- 实现代码：[`composite.go`](./composite.go)
- 测试代码：[`composite_test.go`](./composite_test.go)
