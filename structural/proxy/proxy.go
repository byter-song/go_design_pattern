// Package proxy 展示了 Go 语言中实现代理模式的惯用方法。
//
// 代理模式为其他对象提供一种代理以控制对这个对象的访问。
// 在 Go 中，代理模式常用于：
//   - 延迟加载（Virtual Proxy）
//   - 访问控制（Protection Proxy）
//   - 远程代理（Remote Proxy）
//   - 智能引用（Smart Reference）
//   - 缓存代理（Caching Proxy）
//
// 关键概念：
//   - 代理对象和被代理对象实现相同接口
//   - 客户端通过代理间接访问真实对象
//   - 代理可以在访问前后添加额外逻辑
//
// Go 语言特色：
//   - 使用接口实现透明代理
//   - 利用闭包实现简单代理
//   - 结合 channel 实现并发控制
package proxy

import (
	"fmt"
	"sync"
	"time"
)

// ============================================================================
// 核心接口定义
// ============================================================================

// Image 是图片接口，作为被代理的接口
type Image interface {
	Display() error
	GetFilename() string
	GetSize() int
}

// ============================================================================
// 真实对象：真实图片
// ============================================================================

// RealImage 是真实的图片对象，加载成本高
type RealImage struct {
	filename string
	size     int
	data     []byte // 模拟图片数据
}

// NewRealImage 创建真实图片（模拟从磁盘加载）
func NewRealImage(filename string) *RealImage {
	fmt.Printf("[RealImage] Loading image from disk: %s\n", filename)
	// 模拟耗时加载
	time.Sleep(100 * time.Millisecond)
	return &RealImage{
		filename: filename,
		size:     1024 * 1024, // 1MB
		data:     make([]byte, 1024*1024),
	}
}

// Display 显示图片
func (r *RealImage) Display() error {
	fmt.Printf("[RealImage] Displaying: %s\n", r.filename)
	return nil
}

// GetFilename 获取文件名
func (r *RealImage) GetFilename() string {
	return r.filename
}

// GetSize 获取大小
func (r *RealImage) GetSize() int {
	return r.size
}

// ============================================================================
// 虚拟代理：延迟加载
// ============================================================================

// ProxyImage 是图片代理，实现延迟加载
type ProxyImage struct {
	filename  string
	realImage *RealImage
	mu        sync.Mutex
}

// NewProxyImage 创建图片代理
func NewProxyImage(filename string) *ProxyImage {
	return &ProxyImage{filename: filename}
}

// Display 显示图片（延迟加载）
func (p *ProxyImage) Display() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 延迟加载：第一次访问时才创建真实对象
	if p.realImage == nil {
		p.realImage = NewRealImage(p.filename)
	}
	return p.realImage.Display()
}

// GetFilename 获取文件名（无需加载真实对象）
func (p *ProxyImage) GetFilename() string {
	return p.filename
}

// GetSize 获取大小（延迟加载）
func (p *ProxyImage) GetSize() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.realImage == nil {
		p.realImage = NewRealImage(p.filename)
	}
	return p.realImage.GetSize()
}

// IsLoaded 检查是否已加载（代理特有方法）
func (p *ProxyImage) IsLoaded() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.realImage != nil
}

// ============================================================================
// 保护代理：访问控制
// ============================================================================

// User 表示用户
type User struct {
	Name string
	Role string // admin, user, guest
}

// IsAdmin 检查是否是管理员
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// ProtectedImage 是受保护的图片代理
type ProtectedImage struct {
	image    Image
	user     *User
	adminOnly bool
}

// NewProtectedImage 创建受保护图片代理
func NewProtectedImage(image Image, user *User, adminOnly bool) *ProtectedImage {
	return &ProtectedImage{
		image:     image,
		user:      user,
		adminOnly: adminOnly,
	}
}

// Display 显示图片（带权限检查）
func (p *ProtectedImage) Display() error {
	if p.adminOnly && !p.user.IsAdmin() {
		return fmt.Errorf("access denied: %s is not admin", p.user.Name)
	}
	fmt.Printf("[ProtectedImage] Access granted to %s\n", p.user.Name)
	return p.image.Display()
}

// GetFilename 获取文件名
func (p *ProtectedImage) GetFilename() string {
	return p.image.GetFilename()
}

// GetSize 获取大小
func (p *ProtectedImage) GetSize() int {
	return p.image.GetSize()
}

// ============================================================================
// 缓存代理
// ============================================================================

// CacheProxy 是缓存代理
type CacheProxy struct {
	image      Image
	cache      map[string]interface{}
	cacheTime  map[string]time.Time
	ttl        time.Duration
	mu         sync.RWMutex
}

// NewCacheProxy 创建缓存代理
func NewCacheProxy(image Image, ttl time.Duration) *CacheProxy {
	return &CacheProxy{
		image:     image,
		cache:     make(map[string]interface{}),
		cacheTime: make(map[string]time.Time),
		ttl:       ttl,
	}
}

// Display 显示图片（带缓存）
func (c *CacheProxy) Display() error {
	// Display 操作不缓存，直接执行
	return c.image.Display()
}

// GetFilename 获取文件名（缓存）
func (c *CacheProxy) GetFilename() string {
	return c.getCached("filename", func() interface{} {
		return c.image.GetFilename()
	}).(string)
}

// GetSize 获取大小（缓存）
func (c *CacheProxy) GetSize() int {
	return c.getCached("size", func() interface{} {
		return c.image.GetSize()
	}).(int)
}

// getCached 获取缓存值
func (c *CacheProxy) getCached(key string, fetch func() interface{}) interface{} {
	c.mu.RLock()
	if val, ok := c.cache[key]; ok {
		if time.Since(c.cacheTime[key]) < c.ttl {
			c.mu.RUnlock()
			fmt.Printf("[CacheProxy] Cache hit for %s\n", key)
			return val
		}
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// 双重检查
	if val, ok := c.cache[key]; ok {
		if time.Since(c.cacheTime[key]) < c.ttl {
			return val
		}
	}

	fmt.Printf("[CacheProxy] Cache miss for %s, fetching...\n", key)
	val := fetch()
	c.cache[key] = val
	c.cacheTime[key] = time.Now()
	return val
}

// ClearCache 清除缓存
func (c *CacheProxy) ClearCache() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = make(map[string]interface{})
	c.cacheTime = make(map[string]time.Time)
}

// ============================================================================
// 智能引用代理
// ============================================================================

// SmartReferenceProxy 是智能引用代理
type SmartReferenceProxy struct {
	image       Image
	refCount    int
	lastAccess  time.Time
	mu          sync.Mutex
}

// NewSmartReferenceProxy 创建智能引用代理
func NewSmartReferenceProxy(image Image) *SmartReferenceProxy {
	return &SmartReferenceProxy{
		image:      image,
		lastAccess: time.Now(),
	}
}

// Display 显示图片（记录访问）
func (s *SmartReferenceProxy) Display() error {
	s.mu.Lock()
	s.refCount++
	s.lastAccess = time.Now()
	s.mu.Unlock()
	return s.image.Display()
}

// GetFilename 获取文件名
func (s *SmartReferenceProxy) GetFilename() string {
	s.mu.Lock()
	s.lastAccess = time.Now()
	s.mu.Unlock()
	return s.image.GetFilename()
}

// GetSize 获取大小
func (s *SmartReferenceProxy) GetSize() int {
	s.mu.Lock()
	s.lastAccess = time.Now()
	s.mu.Unlock()
	return s.image.GetSize()
}

// GetRefCount 获取引用计数
func (s *SmartReferenceProxy) GetRefCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.refCount
}

// GetLastAccess 获取最后访问时间
func (s *SmartReferenceProxy) GetLastAccess() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lastAccess
}

// ============================================================================
// 远程代理（模拟）
// ============================================================================

// RemoteImage 是远程图片代理
type RemoteImage struct {
	serverURL string
	filename  string
}

// NewRemoteImage 创建远程图片代理
func NewRemoteImage(serverURL, filename string) *RemoteImage {
	return &RemoteImage{
		serverURL: serverURL,
		filename:  filename,
	}
}

// Display 从远程服务器获取并显示图片
func (r *RemoteImage) Display() error {
	fmt.Printf("[RemoteImage] Fetching %s from %s...\n", r.filename, r.serverURL)
	// 模拟网络延迟
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("[RemoteImage] Displaying remote image: %s\n", r.filename)
	return nil
}

// GetFilename 获取文件名
func (r *RemoteImage) GetFilename() string {
	return r.filename
}

// GetSize 获取大小（模拟）
func (r *RemoteImage) GetSize() int {
	// 模拟获取远程图片大小
	return 2048 * 1024 // 2MB
}
