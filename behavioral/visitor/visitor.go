// Package visitor 展示了 Go 语言中实现访问者模式的惯用方法。
//
// 访问者模式把“操作”从对象结构中分离出来，使得在不修改元素类型的情况下新增操作成为可能。
package visitor

import "fmt"

// Visitor 定义访问行为。
type Visitor interface {
	VisitFile(*File)
	VisitDirectory(*Directory)
}

// Element 是可被访问的元素。
type Element interface {
	Accept(Visitor)
}

// File 是叶子元素。
type File struct {
	Name string
	Size int
}

// Accept 接受访问。
func (f *File) Accept(v Visitor) {
	v.VisitFile(f)
}

// Directory 是组合元素。
type Directory struct {
	Name     string
	Children []Element
}

// Accept 接受访问。
func (d *Directory) Accept(v Visitor) {
	v.VisitDirectory(d)
	for _, child := range d.Children {
		child.Accept(v)
	}
}

// SizeVisitor 统计大小。
type SizeVisitor struct {
	Total int
}

// VisitFile 处理文件。
func (v *SizeVisitor) VisitFile(file *File) {
	v.Total += file.Size
}

// VisitDirectory 处理目录。
func (v *SizeVisitor) VisitDirectory(dir *Directory) {}

// NameVisitor 收集名称。
type NameVisitor struct {
	Names []string
}

// VisitFile 处理文件。
func (v *NameVisitor) VisitFile(file *File) {
	v.Names = append(v.Names, "file:"+file.Name)
}

// VisitDirectory 处理目录。
func (v *NameVisitor) VisitDirectory(dir *Directory) {
	v.Names = append(v.Names, "dir:"+dir.Name)
}

// Summary 返回名称摘要。
func (v *NameVisitor) Summary() string {
	return fmt.Sprint(v.Names)
}
