// Package iterator 展示了 Go 语言中实现迭代器模式的惯用方法。
//
// Go 传统上更倾向使用 for-range，而不是显式定义复杂迭代器对象。
// 但在需要隐藏遍历细节、支持自定义生成逻辑时，闭包迭代器和 Go 1.22 的 iter.Seq 都很合适。
package iterator

import "iter"

// IntCollection 是一个简单集合。
type IntCollection struct {
	items []int
}

// NewIntCollection 创建集合。
func NewIntCollection(items ...int) *IntCollection {
	copied := append([]int(nil), items...)
	return &IntCollection{items: copied}
}

// Iterator 返回传统闭包迭代器。
func (c *IntCollection) Iterator() func() (int, bool) {
	index := 0
	return func() (int, bool) {
		if index >= len(c.items) {
			return 0, false
		}
		value := c.items[index]
		index++
		return value, true
	}
}

// Seq 返回 Go 1.22+ 的标准迭代序列。
func (c *IntCollection) Seq() iter.Seq[int] {
	return func(yield func(int) bool) {
		for _, item := range c.items {
			if !yield(item) {
				return
			}
		}
	}
}

// FilterSeq 过滤一个序列。
func FilterSeq(seq iter.Seq[int], keep func(int) bool) iter.Seq[int] {
	return func(yield func(int) bool) {
		for item := range seq {
			if keep(item) && !yield(item) {
				return
			}
		}
	}
}
