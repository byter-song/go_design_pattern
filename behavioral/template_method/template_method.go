// Package template_method 展示了 Go 语言中实现模板方法模式的惯用方法。
//
// 模板方法模式定义算法骨架，把具体步骤延迟到子类型或外部实现。
// Go 没有继承，因此常见写法不是“父类模板方法 + 子类重写”，而是：
//   1. 组合 + 接口：骨架对象调用外部实现的步骤
//   2. 闭包函数：骨架流程接收可替换步骤
package template_method

import (
	"fmt"
	"strings"
)

// ============================================================================
// 方案一：组合 + 接口
// ============================================================================

// ReportSteps 定义模板步骤。
type ReportSteps interface {
	LoadData() string
	Analyze(string) string
	Format(string) string
}

// ReportGenerator 封装算法骨架。
type ReportGenerator struct {
	steps ReportSteps
}

// NewReportGenerator 创建生成器。
func NewReportGenerator(steps ReportSteps) *ReportGenerator {
	return &ReportGenerator{steps: steps}
}

// Generate 定义固定流程。
func (g *ReportGenerator) Generate() string {
	raw := g.steps.LoadData()
	analyzed := g.steps.Analyze(raw)
	return g.steps.Format(analyzed)
}

// SalesReport 是一种具体实现。
type SalesReport struct{}

func (s SalesReport) LoadData() string     { return "sales:100,200,300" }
func (s SalesReport) Analyze(data string) string { return "sales-total=600 from " + data }
func (s SalesReport) Format(data string) string  { return "[SALES] " + strings.ToUpper(data) }

// ============================================================================
// 方案二：闭包函数
// ============================================================================

// Pipeline 是基于函数注入的模板方法。
type Pipeline struct {
	Load   func() string
	Handle func(string) string
	Save   func(string) string
}

// Execute 执行骨架流程。
func (p Pipeline) Execute() string {
	data := p.Load()
	handled := p.Handle(data)
	return p.Save(handled)
}

// DefaultCSVImport 提供一个基于闭包的模板实例。
func DefaultCSVImport() Pipeline {
	return Pipeline{
		Load: func() string {
			return "id,name\n1,go"
		},
		Handle: func(data string) string {
			return strings.ReplaceAll(data, ",", "|")
		},
		Save: func(data string) string {
			return fmt.Sprintf("saved:%s", data)
		},
	}
}
