// Package state 展示了 Go 语言中实现状态模式的惯用方法。
//
// 状态模式把不同状态下的行为拆分到独立对象中，使上下文不再依赖大量 switch-case。
// 在 Go 中，这通常通过接口 + 上下文持有当前状态来实现。
//
// Go 语言实现特点：
//   1. 不同状态实现同一接口
//   2. 状态切换逻辑分散到各状态对象中
//   3. 上下文不需要维护复杂的分支判断
package state

import "fmt"

// OrderState 是订单状态接口。
type OrderState interface {
	Pay(*Order) error
	Ship(*Order) error
	Complete(*Order) error
	Name() string
}

// Order 是上下文对象。
type Order struct {
	state OrderState
}

// NewOrder 创建新订单。
func NewOrder() *Order {
	return &Order{state: pendingState{}}
}

// StateName 返回当前状态名。
func (o *Order) StateName() string {
	return o.state.Name()
}

// Pay 执行支付。
func (o *Order) Pay() error {
	return o.state.Pay(o)
}

// Ship 执行发货。
func (o *Order) Ship() error {
	return o.state.Ship(o)
}

// Complete 执行完成。
func (o *Order) Complete() error {
	return o.state.Complete(o)
}

func (o *Order) transitionTo(state OrderState) {
	o.state = state
}

type pendingState struct{}

func (s pendingState) Pay(o *Order) error {
	o.transitionTo(paidState{})
	return nil
}
func (s pendingState) Ship(o *Order) error     { return fmt.Errorf("cannot ship in %s", s.Name()) }
func (s pendingState) Complete(o *Order) error { return fmt.Errorf("cannot complete in %s", s.Name()) }
func (s pendingState) Name() string            { return "pending" }

type paidState struct{}

func (s paidState) Pay(o *Order) error         { return fmt.Errorf("already paid") }
func (s paidState) Ship(o *Order) error        { o.transitionTo(shippedState{}); return nil }
func (s paidState) Complete(o *Order) error    { return fmt.Errorf("cannot complete in %s", s.Name()) }
func (s paidState) Name() string               { return "paid" }

type shippedState struct{}

func (s shippedState) Pay(o *Order) error      { return fmt.Errorf("already paid") }
func (s shippedState) Ship(o *Order) error     { return fmt.Errorf("already shipped") }
func (s shippedState) Complete(o *Order) error { o.transitionTo(completedState{}); return nil }
func (s shippedState) Name() string            { return "shipped" }

type completedState struct{}

func (s completedState) Pay(o *Order) error      { return fmt.Errorf("order completed") }
func (s completedState) Ship(o *Order) error     { return fmt.Errorf("order completed") }
func (s completedState) Complete(o *Order) error { return fmt.Errorf("already completed") }
func (s completedState) Name() string            { return "completed" }
