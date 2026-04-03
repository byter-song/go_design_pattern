// Package command 展示了 Go 语言中实现命令模式的惯用方法。
//
// 命令模式把请求封装成对象或可调用值，使请求发送者与执行者解耦。
// 在 Go 中，除了传统的接口式命令，func() 或 func() error 往往是更轻量、更实用的实现方式。
package command

import "fmt"

// ============================================================================
// 传统接口式命令
// ============================================================================

// Command 是命令接口。
type Command interface {
	Execute() string
}

// Light 是接收者。
type Light struct {
	on bool
}

// NewLight 创建灯对象。
func NewLight() *Light {
	return &Light{}
}

// TurnOn 打开灯。
func (l *Light) TurnOn() string {
	l.on = true
	return "light on"
}

// TurnOff 关闭灯。
func (l *Light) TurnOff() string {
	l.on = false
	return "light off"
}

// IsOn 返回状态。
func (l *Light) IsOn() bool {
	return l.on
}

// TurnOnCommand 是传统命令对象。
type TurnOnCommand struct {
	light *Light
}

// NewTurnOnCommand 创建命令。
func NewTurnOnCommand(light *Light) *TurnOnCommand {
	return &TurnOnCommand{light: light}
}

// Execute 执行命令。
func (c *TurnOnCommand) Execute() string {
	return c.light.TurnOn()
}

// TurnOffCommand 是传统命令对象。
type TurnOffCommand struct {
	light *Light
}

// NewTurnOffCommand 创建命令。
func NewTurnOffCommand(light *Light) *TurnOffCommand {
	return &TurnOffCommand{light: light}
}

// Execute 执行命令。
func (c *TurnOffCommand) Execute() string {
	return c.light.TurnOff()
}

// Button 是调用者。
type Button struct {
	command Command
}

// NewButton 创建按钮。
func NewButton(command Command) *Button {
	return &Button{command: command}
}

// Press 按下按钮。
func (b *Button) Press() string {
	return b.command.Execute()
}

// ============================================================================
// Go 风格：func 命令
// ============================================================================

// FuncCommand 是轻量级命令类型。
type FuncCommand func() string

// Execute 执行函数命令。
func (c FuncCommand) Execute() string {
	return c()
}

// Scheduler 可以顺序执行多条命令。
type Scheduler struct {
	commands []FuncCommand
}

// Add 添加命令。
func (s *Scheduler) Add(cmd FuncCommand) {
	s.commands = append(s.commands, cmd)
}

// RunAll 执行全部命令。
func (s *Scheduler) RunAll() []string {
	results := make([]string, 0, len(s.commands))
	for _, cmd := range s.commands {
		results = append(results, cmd())
	}
	return results
}

// LogCommand 返回一个简单命令。
func LogCommand(msg string) FuncCommand {
	return func() string {
		return fmt.Sprintf("log:%s", msg)
	}
}
