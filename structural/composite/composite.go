// Package composite 展示了 Go 语言中实现组合模式的惯用方法。
//
// 组合模式将对象组织成树形结构，使客户端能够以统一方式处理单个对象和对象集合。
// 在 Go 中，这种模式通常由接口 + 切片 + 递归组合完成，非常自然。
//
// Go 语言实现特点：
//   1. 使用接口统一叶子节点与容器节点
//   2. 使用切片保存子节点，天然适合树形结构
//   3. 通过递归遍历完成聚合操作
//   4. 零值切片即可表示“没有子节点”
package composite

import (
	"fmt"
	"strings"
)

// ============================================================================
// 场景：文件系统树
// ============================================================================

// Node 是文件系统节点的统一接口。
type Node interface {
	Name() string
	Size() int
	Count() int
	Render(indent int) string
}

// File 表示叶子节点。
type File struct {
	name string
	size int
}

// NewFile 创建文件。
func NewFile(name string, size int) *File {
	return &File{name: name, size: size}
}

// Name 返回名称。
func (f *File) Name() string {
	return f.name
}

// Size 返回文件大小。
func (f *File) Size() int {
	return f.size
}

// Count 返回节点数量。
func (f *File) Count() int {
	return 1
}

// Render 输出树形文本。
func (f *File) Render(indent int) string {
	return fmt.Sprintf("%s- %s (%dKB)", strings.Repeat("  ", indent), f.name, f.size)
}

// Directory 表示组合节点。
type Directory struct {
	name     string
	children []Node
}

// NewDirectory 创建目录。
func NewDirectory(name string) *Directory {
	return &Directory{
		name:     name,
		children: make([]Node, 0),
	}
}

// Name 返回目录名。
func (d *Directory) Name() string {
	return d.name
}

// Add 添加子节点。
func (d *Directory) Add(nodes ...Node) {
	d.children = append(d.children, nodes...)
}

// Remove 删除指定名称的子节点。
func (d *Directory) Remove(name string) bool {
	for i, child := range d.children {
		if child.Name() == name {
			d.children = append(d.children[:i], d.children[i+1:]...)
			return true
		}
	}
	return false
}

// Children 返回当前子节点快照。
func (d *Directory) Children() []Node {
	result := make([]Node, len(d.children))
	copy(result, d.children)
	return result
}

// Size 返回目录总大小。
func (d *Directory) Size() int {
	total := 0
	for _, child := range d.children {
		total += child.Size()
	}
	return total
}

// Count 返回包含自身在内的节点数量。
func (d *Directory) Count() int {
	total := 1
	for _, child := range d.children {
		total += child.Count()
	}
	return total
}

// Render 输出目录树。
func (d *Directory) Render(indent int) string {
	lines := []string{
		fmt.Sprintf("%s+ %s/", strings.Repeat("  ", indent), d.name),
	}
	for _, child := range d.children {
		lines = append(lines, child.Render(indent+1))
	}
	return strings.Join(lines, "\n")
}

// Find 在树中查找节点。
func (d *Directory) Find(name string) Node {
	if d.name == name {
		return d
	}
	for _, child := range d.children {
		if child.Name() == name {
			return child
		}
		dir, ok := child.(*Directory)
		if ok {
			if found := dir.Find(name); found != nil {
				return found
			}
		}
	}
	return nil
}

// Walk 深度优先遍历树。
func (d *Directory) Walk(visit func(Node)) {
	visit(d)
	for _, child := range d.children {
		if dir, ok := child.(*Directory); ok {
			dir.Walk(visit)
			continue
		}
		visit(child)
	}
}
