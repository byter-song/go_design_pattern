// Package chain_of_responsibility 展示了 Go 语言中实现责任链模式的惯用方法。
//
// 责任链模式允许多个对象有机会处理请求，从而避免请求的发送者和接收者之间的耦合。
// 将这些对象连成一条链，并沿着这条链传递请求，直到有一个对象处理它为止。
//
// Go 语言实现特点：
//   1. 使用 interface 定义处理器契约
//   2. 每个处理器持有下一个处理器的引用（链式结构）
//   3. 利用函数类型实现轻量级处理器
//   4. 支持动态构建和修改责任链
//
// 与装饰器模式的区别：
//   - 装饰器模式：所有装饰器都会处理请求，用于增强功能
//   - 责任链模式：只有一个处理器处理请求，用于决策分发
//
// 适用场景：
//   - 多个对象可以处理同一请求，具体由运行时决定
//   - 不想显式指定处理者
//   - 需要动态指定处理者集合
//   - 审批流程、中间件链、日志级别处理等
package chain_of_responsibility

import (
	"fmt"
	"strings"
	"time"
)

// ============================================================================
// 场景：技术支持工单处理系统
// ============================================================================

// SupportTicket 是支持工单
type SupportTicket struct {
	ID          string
	Level       int           // 1=一线, 2=二线, 3=三线
	Category    string        // "technical", "billing", "general"
	Description string
	CreatedAt   time.Time
}

// NewSupportTicket 创建支持工单
func NewSupportTicket(id, category, description string, level int) *SupportTicket {
	return &SupportTicket{
		ID:          id,
		Level:       level,
		Category:    category,
		Description: description,
		CreatedAt:   time.Now(),
	}
}

// ============================================================================
// 处理器接口
// ============================================================================

// SupportHandler 是支持处理器接口
type SupportHandler interface {
	// Handle 处理工单，返回是否处理成功
	Handle(ticket *SupportTicket) bool
	// SetNext 设置下一个处理器
	SetNext(handler SupportHandler) SupportHandler
	// GetName 返回处理器名称
	GetName() string
}

// ============================================================================
// 基础处理器（抽象层）
// ============================================================================

// BaseHandler 是处理器的基础实现，提供链式结构
type BaseHandler struct {
	next   SupportHandler
	name   string
}

// SetNext 设置下一个处理器，返回下一个处理器以便链式调用
func (b *BaseHandler) SetNext(handler SupportHandler) SupportHandler {
	b.next = handler
	return handler
}

// GetName 返回处理器名称
func (b *BaseHandler) GetName() string {
	return b.name
}

// passToNext 将请求传递给下一个处理器
func (b *BaseHandler) passToNext(ticket *SupportTicket) bool {
	if b.next != nil {
		return b.next.Handle(ticket)
	}
	return false
}

// ============================================================================
// 具体处理器实现
// ============================================================================

// FirstLineSupport 一线支持（处理简单问题）
type FirstLineSupport struct {
	BaseHandler
}

// NewFirstLineSupport 创建一线支持处理器
func NewFirstLineSupport() *FirstLineSupport {
	return &FirstLineSupport{
		BaseHandler: BaseHandler{name: "First Line Support"},
	}
}

// Handle 实现 SupportHandler 接口
func (f *FirstLineSupport) Handle(ticket *SupportTicket) bool {
	// 一线支持只处理级别 1 的一般问题
	if ticket.Level == 1 && ticket.Category == "general" {
		fmt.Printf("[%s] Handling ticket %s: %s\n", f.name, ticket.ID, ticket.Description)
		return true
	}

	fmt.Printf("[%s] Cannot handle ticket %s (level %d, %s), passing to next...\n",
		f.name, ticket.ID, ticket.Level, ticket.Category)
	return f.passToNext(ticket)
}

// SetNext 设置下一个处理器
func (f *FirstLineSupport) SetNext(handler SupportHandler) SupportHandler {
	f.next = handler
	return handler
}

// TechnicalSupport 技术支持（处理技术问题）
type TechnicalSupport struct {
	BaseHandler
}

// NewTechnicalSupport 创建技术支持处理器
func NewTechnicalSupport() *TechnicalSupport {
	return &TechnicalSupport{
		BaseHandler: BaseHandler{name: "Technical Support"},
	}
}

// Handle 实现 SupportHandler 接口
func (t *TechnicalSupport) Handle(ticket *SupportTicket) bool {
	// 技术支持处理技术问题（级别 1-2）
	if ticket.Category == "technical" && ticket.Level <= 2 {
		fmt.Printf("[%s] Handling technical ticket %s: %s\n", t.name, ticket.ID, ticket.Description)
		return true
	}

	fmt.Printf("[%s] Cannot handle ticket %s, passing to next...\n", t.name, ticket.ID)
	return t.passToNext(ticket)
}

// SetNext 设置下一个处理器
func (t *TechnicalSupport) SetNext(handler SupportHandler) SupportHandler {
	t.next = handler
	return handler
}

// BillingSupport 账单支持（处理账单问题）
type BillingSupport struct {
	BaseHandler
}

// NewBillingSupport 创建账单支持处理器
func NewBillingSupport() *BillingSupport {
	return &BillingSupport{
		BaseHandler: BaseHandler{name: "Billing Support"},
	}
}

// Handle 实现 SupportHandler 接口
func (b *BillingSupport) Handle(ticket *SupportTicket) bool {
	// 账单支持处理账单问题（级别 1-2）
	if ticket.Category == "billing" && b.canHandleBilling(ticket) {
		fmt.Printf("[%s] Handling billing ticket %s: %s\n", b.name, ticket.ID, ticket.Description)
		return true
	}

	fmt.Printf("[%s] Cannot handle ticket %s, passing to next...\n", b.name, ticket.ID)
	return b.passToNext(ticket)
}

// canHandleBilling 判断是否能处理账单问题
func (b *BillingSupport) canHandleBilling(ticket *SupportTicket) bool {
	// 简单账单问题可以处理
	return ticket.Level <= 2 && !strings.Contains(ticket.Description, "fraud")
}

// SetNext 设置下一个处理器
func (b *BillingSupport) SetNext(handler SupportHandler) SupportHandler {
	b.next = handler
	return handler
}

// SeniorSupport 高级支持（处理复杂问题）
type SeniorSupport struct {
	BaseHandler
}

// NewSeniorSupport 创建高级支持处理器
func NewSeniorSupport() *SeniorSupport {
	return &SeniorSupport{
		BaseHandler: BaseHandler{name: "Senior Support"},
	}
}

// Handle 实现 SupportHandler 接口
func (s *SeniorSupport) Handle(ticket *SupportTicket) bool {
	// 高级支持处理级别 2 的所有问题
	if ticket.Level == 2 {
		fmt.Printf("[%s] Handling complex ticket %s: %s\n", s.name, ticket.ID, ticket.Description)
		return true
	}

	fmt.Printf("[%s] Cannot handle ticket %s, passing to next...\n", s.name, ticket.ID)
	return s.passToNext(ticket)
}

// SetNext 设置下一个处理器
func (s *SeniorSupport) SetNext(handler SupportHandler) SupportHandler {
	s.next = handler
	return handler
}

// ManagerSupport 经理支持（处理级别 3 的紧急问题）
type ManagerSupport struct {
	BaseHandler
}

// NewManagerSupport 创建经理支持处理器
func NewManagerSupport() *ManagerSupport {
	return &ManagerSupport{
		BaseHandler: BaseHandler{name: "Manager Support"},
	}
}

// Handle 实现 SupportHandler 接口
func (m *ManagerSupport) Handle(ticket *SupportTicket) bool {
	// 经理处理所有级别 3 的问题
	if ticket.Level == 3 {
		fmt.Printf("[%s] Handling escalated ticket %s: %s\n", m.name, ticket.ID, ticket.Description)
		return true
	}

	fmt.Printf("[%s] Cannot handle ticket %s, passing to next...\n", m.name, ticket.ID)
	return m.passToNext(ticket)
}

// SetNext 设置下一个处理器
func (m *ManagerSupport) SetNext(handler SupportHandler) SupportHandler {
	m.next = handler
	return handler
}

// DefaultHandler 默认处理器（处理链末端）
type DefaultHandler struct {
	BaseHandler
}

// NewDefaultHandler 创建默认处理器
func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{
		BaseHandler: BaseHandler{name: "Default Handler"},
	}
}

// Handle 实现 SupportHandler 接口
func (d *DefaultHandler) Handle(ticket *SupportTicket) bool {
	// 默认处理器记录无法处理的工单
	fmt.Printf("[%s] No handler found for ticket %s (level %d, %s)\n",
		d.name, ticket.ID, ticket.Level, ticket.Category)
	return false
}

// SetNext 设置下一个处理器
func (d *DefaultHandler) SetNext(handler SupportHandler) SupportHandler {
	d.next = handler
	return handler
}

// ============================================================================
// 责任链构建器
// ============================================================================

// ChainBuilder 责任链构建器
type ChainBuilder struct {
	head SupportHandler
	tail SupportHandler
}

// NewChainBuilder 创建责任链构建器
func NewChainBuilder() *ChainBuilder {
	return &ChainBuilder{}
}

// AddHandler 添加处理器到链尾
func (c *ChainBuilder) AddHandler(handler SupportHandler) *ChainBuilder {
	if c.head == nil {
		c.head = handler
		c.tail = handler
	} else {
		c.tail.SetNext(handler)
		c.tail = handler
	}
	return c
}

// Build 构建责任链，返回链头
func (c *ChainBuilder) Build() SupportHandler {
	return c.head
}

// ============================================================================
// 函数类型处理器（Go 特有轻量级实现）
// ============================================================================

// Request 是请求结构
type Request struct {
	Content string
	Type    string
}

// Response 是响应结构
type Response struct {
	Handled bool
	Result  string
}

// HandlerFunc 是处理器函数类型
type HandlerFunc func(*Request) *Response

// Chain 责任链
type Chain struct {
	handlers []HandlerFunc
}

// NewChain 创建责任链
func NewChain() *Chain {
	return &Chain{
		handlers: make([]HandlerFunc, 0),
	}
}

// Add 添加处理器到链
func (c *Chain) Add(handler HandlerFunc) *Chain {
	c.handlers = append(c.handlers, handler)
	return c
}

// Execute 执行责任链
// 返回第一个处理成功的响应，如果没有处理器能处理，返回 nil
func (c *Chain) Execute(req *Request) *Response {
	for _, handler := range c.handlers {
		resp := handler(req)
		if resp != nil && resp.Handled {
			return resp
		}
	}
	return &Response{Handled: false, Result: "No handler could process the request"}
}

// ============================================================================
// 中间件风格的责任链
// ============================================================================

// Middleware 中间件函数类型
type Middleware func(ctx *Context) bool

// Context 请求上下文
type Context struct {
	Request  *Request
	Response *Response
	index    int
	handlers []Middleware
}

// NewContext 创建上下文
func NewContext(req *Request) *Context {
	return &Context{
		Request:  req,
		Response: &Response{Handled: false},
		handlers: make([]Middleware, 0),
		index:    -1,
	}
}

// Use 添加中间件
func (c *Context) Use(middleware ...Middleware) {
	c.handlers = append(c.handlers, middleware...)
}

// Next 调用下一个中间件
func (c *Context) Next() bool {
	c.index++
	if c.index < len(c.handlers) {
		return c.handlers[c.index](c)
	}
	return false
}

// Execute 执行中间件链
func (c *Context) Execute() bool {
	return c.Next()
}

// ============================================================================
// 审批流程示例
// ============================================================================

// ExpenseRequest 费用申请
type ExpenseRequest struct {
	ID     string
	Amount float64
	Reason string
}

// Approver 审批者接口
type Approver interface {
	Approve(request *ExpenseRequest) (bool, string)
	SetNext(approver Approver) Approver
}

// TeamLeader 团队负责人
type TeamLeader struct {
	next Approver
	name string
}

// NewTeamLeader 创建团队负责人
func NewTeamLeader(name string) *TeamLeader {
	return &TeamLeader{name: name}
}

// Approve 审批
func (t *TeamLeader) Approve(request *ExpenseRequest) (bool, string) {
	if request.Amount <= 1000 {
		return true, fmt.Sprintf("Approved by Team Leader %s", t.name)
	}
	if t.next != nil {
		return t.next.Approve(request)
	}
	return false, "No approver available"
}

// SetNext 设置下一个审批者
func (t *TeamLeader) SetNext(approver Approver) Approver {
	t.next = approver
	return approver
}

// DepartmentManager 部门经理
type DepartmentManager struct {
	next Approver
	name string
}

// NewDepartmentManager 创建部门经理
func NewDepartmentManager(name string) *DepartmentManager {
	return &DepartmentManager{name: name}
}

// Approve 审批
func (d *DepartmentManager) Approve(request *ExpenseRequest) (bool, string) {
	if request.Amount <= 5000 {
		return true, fmt.Sprintf("Approved by Department Manager %s", d.name)
	}
	if d.next != nil {
		return d.next.Approve(request)
	}
	return false, "No approver available"
}

// SetNext 设置下一个审批者
func (d *DepartmentManager) SetNext(approver Approver) Approver {
	d.next = approver
	return approver
}

// CFO 首席财务官
type CFO struct {
	name string
}

// NewCFO 创建 CFO
func NewCFO(name string) *CFO {
	return &CFO{name: name}
}

// Approve 审批
func (c *CFO) Approve(request *ExpenseRequest) (bool, string) {
	// CFO 可以审批任何金额
	return true, fmt.Sprintf("Approved by CFO %s", c.name)
}

// SetNext 设置下一个审批者
func (c *CFO) SetNext(approver Approver) Approver {
	// CFO 是最终审批者，不需要下一个
	return approver
}
