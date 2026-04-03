// Package flyweight 展示了 Go 语言中实现享元模式的惯用方法。
//
// 享元模式通过共享可复用的内部状态来减少对象创建成本。
// 在 Go 中，除经典的“工厂缓存共享对象”外，sync.Pool 也是非常重要的实现手段，
// 适合复用生命周期短、创建频繁的临时对象。
//
// Go 语言实现特点：
//   1. 使用 map + 互斥锁缓存共享享元
//   2. 使用 sync.Pool 复用高频创建对象
//   3. 将内部状态与外部状态显式拆分
//   4. 让并发环境下的对象复用更自然
package flyweight

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// ============================================================================
// 场景：子弹系统
// ============================================================================

// BulletStyle 是共享的内部状态。
type BulletStyle struct {
	Kind   string
	Color  string
	Sprite string
}

func styleKey(kind, color, sprite string) string {
	return kind + "|" + color + "|" + sprite
}

// StyleFactory 负责缓存 BulletStyle。
type StyleFactory struct {
	mu     sync.RWMutex
	styles map[string]*BulletStyle
}

// NewStyleFactory 创建样式工厂。
func NewStyleFactory() *StyleFactory {
	return &StyleFactory{
		styles: make(map[string]*BulletStyle),
	}
}

// GetStyle 获取共享样式。
func (f *StyleFactory) GetStyle(kind, color, sprite string) *BulletStyle {
	key := styleKey(kind, color, sprite)

	f.mu.RLock()
	style, ok := f.styles[key]
	f.mu.RUnlock()
	if ok {
		return style
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	if style, ok = f.styles[key]; ok {
		return style
	}

	style = &BulletStyle{
		Kind:   kind,
		Color:  color,
		Sprite: sprite,
	}
	f.styles[key] = style
	return style
}

// Count 返回共享样式数量。
func (f *StyleFactory) Count() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.styles)
}

// Bullet 是带外部状态的可复用对象。
type Bullet struct {
	Style  *BulletStyle
	Owner  string
	X      float64
	Y      float64
	Speed  float64
	Active bool
}

// Move 更新子弹位置。
func (b *Bullet) Move(delta float64) {
	if !b.Active {
		return
	}
	b.X += b.Speed * delta
}

// String 返回调试信息。
func (b *Bullet) String() string {
	if b.Style == nil {
		return "<released bullet>"
	}
	return fmt.Sprintf("%s bullet(%s) at %.1f,%.1f", b.Style.Kind, b.Owner, b.X, b.Y)
}

// BulletPool 使用 sync.Pool 复用 Bullet。
type BulletPool struct {
	pool    sync.Pool
	created atomic.Int64
}

// NewBulletPool 创建对象池。
func NewBulletPool() *BulletPool {
	bp := &BulletPool{}
	bp.pool.New = func() any {
		bp.created.Add(1)
		return &Bullet{}
	}
	return bp
}

// Acquire 获取一个可用子弹对象。
func (p *BulletPool) Acquire(style *BulletStyle, owner string, x, y, speed float64) *Bullet {
	bullet := p.pool.Get().(*Bullet)
	bullet.Style = style
	bullet.Owner = owner
	bullet.X = x
	bullet.Y = y
	bullet.Speed = speed
	bullet.Active = true
	return bullet
}

// Release 归还对象到池中。
func (p *BulletPool) Release(bullet *Bullet) {
	if bullet == nil {
		return
	}
	*bullet = Bullet{}
	p.pool.Put(bullet)
}

// CreatedCount 返回对象池创建过的 Bullet 数量。
func (p *BulletPool) CreatedCount() int64 {
	return p.created.Load()
}
