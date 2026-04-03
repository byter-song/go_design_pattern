// Package interpreter 展示了 Go 语言中实现解释器模式的惯用方法。
//
// 解释器模式为某种简单语言定义语法表示，并提供解释执行逻辑。
// 在工程实践中，它更适合表达“小而稳定的领域规则”，而不是复杂通用语言。
package interpreter

import (
	"fmt"
	"strings"
)

// Context 保存变量环境。
type Context struct {
	vars map[string]bool
}

// NewContext 创建上下文。
func NewContext(vars map[string]bool) *Context {
	copied := make(map[string]bool, len(vars))
	for k, v := range vars {
		copied[k] = v
	}
	return &Context{vars: copied}
}

// Lookup 查询变量值。
func (c *Context) Lookup(name string) bool {
	return c.vars[name]
}

// Expression 是表达式接口。
type Expression interface {
	Interpret(*Context) bool
}

// VariableExpression 表示变量。
type VariableExpression struct {
	name string
}

// NewVariable 创建变量表达式。
func NewVariable(name string) VariableExpression {
	return VariableExpression{name: name}
}

// Interpret 解释变量。
func (e VariableExpression) Interpret(ctx *Context) bool {
	return ctx.Lookup(e.name)
}

// AndExpression 表示与运算。
type AndExpression struct {
	left  Expression
	right Expression
}

// And 创建与表达式。
func And(left, right Expression) AndExpression {
	return AndExpression{left: left, right: right}
}

// Interpret 执行与运算。
func (e AndExpression) Interpret(ctx *Context) bool {
	return e.left.Interpret(ctx) && e.right.Interpret(ctx)
}

// OrExpression 表示或运算。
type OrExpression struct {
	left  Expression
	right Expression
}

// Or 创建或表达式。
func Or(left, right Expression) OrExpression {
	return OrExpression{left: left, right: right}
}

// Interpret 执行或运算。
func (e OrExpression) Interpret(ctx *Context) bool {
	return e.left.Interpret(ctx) || e.right.Interpret(ctx)
}

// NotExpression 表示非运算。
type NotExpression struct {
	expr Expression
}

// Not 创建非表达式。
func Not(expr Expression) NotExpression {
	return NotExpression{expr: expr}
}

// Interpret 执行非运算。
func (e NotExpression) Interpret(ctx *Context) bool {
	return !e.expr.Interpret(ctx)
}

// Describe 用于输出表达式描述。
func Describe(parts ...string) string {
	return strings.Join(parts, " ")
}

// EvaluateRule 用于快速评估一个规则。
func EvaluateRule(expr Expression, ctx *Context) string {
	return fmt.Sprintf("result=%t", expr.Interpret(ctx))
}
