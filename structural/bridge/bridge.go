// Package bridge 展示了 Go 语言中实现桥接模式的惯用方法。
//
// 桥接模式将抽象部分与实现部分解耦，使二者可以独立变化。
// 在 Go 中，这通常通过“抽象对象持有接口”来完成，而接口组合（Embedding）
// 又能让我们把实现能力拆成更小的部件，再按需拼装成更丰富的能力集合。
//
// Go 语言实现特点：
//   1. 抽象依赖接口，而不是依赖具体类型
//   2. 使用接口嵌入组合能力边界
//   3. 通过结构体组合实现“抽象层次”
//   4. 新增抽象或实现时互不影响
package bridge

import "fmt"

// ============================================================================
// 实现部分：设备能力接口
// ============================================================================

// Device 表示最基础的设备能力。
type Device interface {
	Name() string
	IsEnabled() bool
	Enable()
	Disable()
}

// VolumeControl 表示音量控制能力。
type VolumeControl interface {
	SetVolume(level int)
	Volume() int
}

// ChannelControl 表示频道控制能力。
type ChannelControl interface {
	SetChannel(channel int)
	Channel() int
}

// EntertainmentDevice 通过接口嵌入组合出更强能力。
type EntertainmentDevice interface {
	Device
	VolumeControl
	ChannelControl
}

// ============================================================================
// 具体实现：电视与机顶盒
// ============================================================================

// TV 是具体设备实现。
type TV struct {
	name    string
	enabled bool
	volume  int
	channel int
}

// NewTV 创建电视设备。
func NewTV(name string) *TV {
	return &TV{
		name:    name,
		volume:  20,
		channel: 1,
	}
}

// Name 返回设备名。
func (t *TV) Name() string {
	return t.name
}

// IsEnabled 返回电源状态。
func (t *TV) IsEnabled() bool {
	return t.enabled
}

// Enable 打开设备。
func (t *TV) Enable() {
	t.enabled = true
}

// Disable 关闭设备。
func (t *TV) Disable() {
	t.enabled = false
}

// SetVolume 设置音量。
func (t *TV) SetVolume(level int) {
	if level < 0 {
		level = 0
	}
	if level > 100 {
		level = 100
	}
	t.volume = level
}

// Volume 返回音量。
func (t *TV) Volume() int {
	return t.volume
}

// SetChannel 设置频道。
func (t *TV) SetChannel(channel int) {
	if channel < 1 {
		channel = 1
	}
	t.channel = channel
}

// Channel 返回频道。
func (t *TV) Channel() int {
	return t.channel
}

// StreamingBox 是另一种具体设备实现。
type StreamingBox struct {
	name    string
	enabled bool
	volume  int
	channel int
}

// NewStreamingBox 创建流媒体设备。
func NewStreamingBox(name string) *StreamingBox {
	return &StreamingBox{
		name:    name,
		volume:  10,
		channel: 100,
	}
}

// Name 返回设备名。
func (s *StreamingBox) Name() string {
	return s.name
}

// IsEnabled 返回电源状态。
func (s *StreamingBox) IsEnabled() bool {
	return s.enabled
}

// Enable 打开设备。
func (s *StreamingBox) Enable() {
	s.enabled = true
}

// Disable 关闭设备。
func (s *StreamingBox) Disable() {
	s.enabled = false
}

// SetVolume 设置音量。
func (s *StreamingBox) SetVolume(level int) {
	if level < 0 {
		level = 0
	}
	if level > 50 {
		level = 50
	}
	s.volume = level
}

// Volume 返回音量。
func (s *StreamingBox) Volume() int {
	return s.volume
}

// SetChannel 设置频道。
func (s *StreamingBox) SetChannel(channel int) {
	if channel < 100 {
		channel = 100
	}
	s.channel = channel
}

// Channel 返回频道。
func (s *StreamingBox) Channel() int {
	return s.channel
}

// ============================================================================
// 抽象部分：遥控器
// ============================================================================

// BasicRemote 是基础抽象。
type BasicRemote struct {
	device Device
}

// NewBasicRemote 创建基础遥控器。
func NewBasicRemote(device Device) *BasicRemote {
	return &BasicRemote{device: device}
}

// TogglePower 切换电源。
func (r *BasicRemote) TogglePower() {
	if r.device.IsEnabled() {
		r.device.Disable()
		return
	}
	r.device.Enable()
}

// DeviceName 返回设备名。
func (r *BasicRemote) DeviceName() string {
	return r.device.Name()
}

// SmartRemote 是扩展抽象。
type SmartRemote struct {
	*BasicRemote
	device EntertainmentDevice
}

// NewSmartRemote 创建智能遥控器。
func NewSmartRemote(device EntertainmentDevice) *SmartRemote {
	return &SmartRemote{
		BasicRemote: NewBasicRemote(device),
		device:      device,
	}
}

// VolumeUp 调高音量。
func (r *SmartRemote) VolumeUp(step int) {
	r.device.SetVolume(r.device.Volume() + step)
}

// ChannelNext 切换到下一个频道。
func (r *SmartRemote) ChannelNext() {
	r.device.SetChannel(r.device.Channel() + 1)
}

// Status 返回当前状态。
func (r *SmartRemote) Status() string {
	return fmt.Sprintf("%s power=%t volume=%d channel=%d",
		r.device.Name(),
		r.device.IsEnabled(),
		r.device.Volume(),
		r.device.Channel(),
	)
}
