// Package abstract_factory 展示了 Go 语言中实现抽象工厂模式的惯用方法。
//
// 抽象工厂用于创建一组相关或相互依赖的对象族，确保它们在同一“主题/外观”下保持一致性。
// 在 Go 中，我们通过接口定义产品族与工厂契约，利用组合与依赖注入来获得灵活性。
//
// Go 语言实现特点：
//   1. 工厂返回接口而非具体类型，客户端仅依赖接口
//   2. 通过多个小接口表示“能力族”，避免臃肿类型
//   3. 倾向轻量化工厂函数，而非层级复杂的类结构
package abstract_factory

import "fmt"

// ============================================================================
// 产品族接口
// ============================================================================

// Button 是按钮接口。
type Button interface {
	Render() string
	Click() string
}

// Checkbox 是复选框接口。
type Checkbox interface {
	Check(on bool) string
	Render() string
}

// UIFactory 是抽象工厂，负责创建同一主题下的产品族。
type UIFactory interface {
	CreateButton() Button
	CreateCheckbox() Checkbox
	Theme() string
}

// ============================================================================
// 具体产品与工厂：Mac 主题
// ============================================================================

type macButton struct{}

func (b *macButton) Render() string { return "MacButton" }
func (b *macButton) Click() string  { return "MacButton clicked" }

type macCheckbox struct {
	on bool
}

func (c *macCheckbox) Check(on bool) string {
	c.on = on
	if on {
		return "MacCheckbox ON"
	}
	return "MacCheckbox OFF"
}
func (c *macCheckbox) Render() string { return "MacCheckbox" }

type MacFactory struct{}

func (f *MacFactory) CreateButton() Button   { return &macButton{} }
func (f *MacFactory) CreateCheckbox() Checkbox { return &macCheckbox{} }
func (f *MacFactory) Theme() string          { return "Mac" }

// ============================================================================
// 具体产品与工厂：Windows 主题
// ============================================================================

type winButton struct{}

func (b *winButton) Render() string { return "WinButton" }
func (b *winButton) Click() string  { return "WinButton clicked" }

type winCheckbox struct {
	on bool
}

func (c *winCheckbox) Check(on bool) string {
	c.on = on
	if on {
		return "WinCheckbox ON"
	}
	return "WinCheckbox OFF"
}
func (c *winCheckbox) Render() string { return "WinCheckbox" }

type WinFactory struct{}

func (f *WinFactory) CreateButton() Button   { return &winButton{} }
func (f *WinFactory) CreateCheckbox() Checkbox { return &winCheckbox{} }
func (f *WinFactory) Theme() string          { return "Windows" }

// ============================================================================
// 客户端：只依赖接口
// ============================================================================

// BuildUI 使用给定工厂构建一套 UI，并返回渲染摘要。
func BuildUI(factory UIFactory) string {
	btn := factory.CreateButton()
	cb := factory.CreateCheckbox()
	return fmt.Sprintf("[%s] %s | %s",
		factory.Theme(),
		btn.Render(),
		cb.Render(),
	)
}
