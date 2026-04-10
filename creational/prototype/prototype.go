// Package prototype 展示了 Go 语言中实现原型模式的惯用方法。
//
// 原型模式通过复制现有对象来创建新对象，适用于对象初始化成本高或需要保留某个模板状态的场景。
// 在 Go 中，这通常表现为：
//  1. 直接值拷贝得到浅拷贝
//  2. 手写 Clone/DeepClone 方法完成深拷贝
//  3. 特别关注切片、map、指针等引用语义字段
package prototype

// Document 表示一个包含引用类型字段的原型对象。
type Document struct {
	Title    string
	Tags     []string
	Metadata map[string]string
	Sections []Section
}

// Section 表示文档的一个段落。
type Section struct {
	Heading string
	Notes   []string
}

// ShallowClone 执行浅拷贝。
func (d *Document) ShallowClone() *Document {
	if d == nil {
		return nil
	}
	copy := *d
	return &copy
}

// DeepClone 执行深拷贝。
func (d *Document) DeepClone() *Document {
	if d == nil {
		return nil
	}

	cloned := &Document{
		Title: d.Title,
	}

	if d.Tags != nil {
		cloned.Tags = append([]string(nil), d.Tags...)
	}

	if d.Metadata != nil {
		cloned.Metadata = make(map[string]string, len(d.Metadata))
		for k, v := range d.Metadata {
			cloned.Metadata[k] = v
		}
	}

	if d.Sections != nil {
		cloned.Sections = make([]Section, len(d.Sections))
		for i, section := range d.Sections {
			cloned.Sections[i] = Section{
				Heading: section.Heading,
			}
			if section.Notes != nil {
				cloned.Sections[i].Notes = append([]string(nil), section.Notes...)
			}
		}
	}

	return cloned
}
