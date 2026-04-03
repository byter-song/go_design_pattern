# 迭代器模式 (Iterator Pattern)

## 概述

迭代器模式用于顺序访问聚合对象中的元素，而不暴露其内部表示。

在 Go 里，很多时候直接 `for range` 就够了，因此传统面向对象风格的显式迭代器没有那么常见。但在下面这些情况下，迭代器依然很有价值：

- 想隐藏底层存储结构
- 想支持惰性生成
- 想对遍历流程做包装、过滤、组合

## Go 语言实现特点

- 传统写法可以用闭包返回 `next()`
- Go 1.22 起可以使用 `iter.Seq`
- 结合 `range` 使用时非常自然

## 方案一：传统闭包迭代器

```go
next := collection.Iterator()
for {
    v, ok := next()
    if !ok {
        break
    }
}
```

这种方式兼容性强，理解成本低。

## 方案二：Go 1.22 的 iter.Seq

```go
for v := range collection.Seq() {
    fmt.Println(v)
}
```

这是新标准写法，优点是：

- 与 `range` 语法直接集成
- 更适合惰性遍历
- 容易封装过滤、映射等操作

## 代码示例

当前实现展示了：

1. `Iterator()`：返回闭包形式的传统迭代器
2. `Seq()`：返回 `iter.Seq[int]`
3. `FilterSeq()`：对序列做过滤包装

## 使用场景

1. 自定义集合遍历
2. 惰性生成数据
3. 封装查询、过滤、管道式遍历

## 优缺点

### 优点

- 隐藏底层存储细节
- 支持惰性遍历
- `iter.Seq` 与 Go 新语法配合自然

### 缺点

- 对简单切片而言，直接 `range` 更直观
- 新的 `iter.Seq` 对旧版本 Go 不兼容

## 参考

- 实现代码：[`iterator.go`](./iterator.go)
- 测试代码：[`iterator_test.go`](./iterator_test.go)
