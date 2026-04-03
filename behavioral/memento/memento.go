// Package memento 展示了 Go 语言中实现备忘录模式的惯用方法。
//
// 备忘录模式在不破坏封装的前提下保存对象状态，方便后续恢复。
// 在 Go 中，这通常表现为“Originator 导出快照对象 + Caretaker 管理历史记录”。
package memento

// Editor 是原发器。
type Editor struct {
	content string
}

// NewEditor 创建编辑器。
func NewEditor() *Editor {
	return &Editor{}
}

// Write 写入内容。
func (e *Editor) Write(content string) {
	e.content = content
}

// Content 返回当前内容。
func (e *Editor) Content() string {
	return e.content
}

// Save 保存快照。
func (e *Editor) Save() Snapshot {
	return Snapshot{content: e.content}
}

// Restore 恢复快照。
func (e *Editor) Restore(snapshot Snapshot) {
	e.content = snapshot.content
}

// Snapshot 是备忘录对象。
type Snapshot struct {
	content string
}

// History 管理快照历史。
type History struct {
	snapshots []Snapshot
}

// Push 保存一份快照。
func (h *History) Push(snapshot Snapshot) {
	h.snapshots = append(h.snapshots, snapshot)
}

// Pop 弹出最近一次快照。
func (h *History) Pop() (Snapshot, bool) {
	if len(h.snapshots) == 0 {
		return Snapshot{}, false
	}
	last := h.snapshots[len(h.snapshots)-1]
	h.snapshots = h.snapshots[:len(h.snapshots)-1]
	return last, true
}
