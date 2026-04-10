// Package adapter 展示了 Go 语言中实现适配器模式的惯用方法。
//
// 适配器模式允许将一个类的接口转换成客户希望的另外一个接口。
// 在 Go 中，我们利用隐式接口（Implicit Interface）实现来实现这一模式，
// 这是 Go 语言最独特和强大的特性之一。
//
// 关键概念：
//   - 隐式接口：不需要显式声明实现某个接口
//   - 接口组合：通过嵌入接口实现功能扩展
//   - 类型断言：在适配时进行类型转换
//
// Go 语言特色：
//   - 隐式接口让适配器实现非常自然
//   - 可以在不修改原代码的情况下进行适配
//   - 适配器可以是临时性的（匿名适配）
package adapter

import (
	"fmt"
	"io"
	"time"
)

// ============================================================================
// 场景：第三方支付适配
// ============================================================================

// PaymentProcessor 是我们期望的统一支付接口。
// 这是客户端代码依赖的接口。
type PaymentProcessor interface {
	ProcessPayment(amount float64, currency string) (string, error)
	RefundPayment(transactionID string) error
}

// ============================================================================
// 旧支付系统（需要被适配）
// ============================================================================

// LegacyPaymentSystem 是旧的支付系统，接口不兼容。
// 假设这是一个我们无法修改的第三方库。
type LegacyPaymentSystem struct {
	merchantID string
}

// NewLegacyPaymentSystem 创建旧支付系统实例
func NewLegacyPaymentSystem(merchantID string) *LegacyPaymentSystem {
	return &LegacyPaymentSystem{merchantID: merchantID}
}

// MakeCharge 旧系统的支付方法（名称和签名都不同）
func (l *LegacyPaymentSystem) MakeCharge(amountInCents int64, currencyCode string) (int64, error) {
	// 模拟处理
	transactionID := time.Now().UnixNano()
	fmt.Printf("[Legacy] Charging %d cents in %s\n", amountInCents, currencyCode)
	return transactionID, nil
}

// DoRefund 旧系统的退款方法
func (l *LegacyPaymentSystem) DoRefund(transactionID int64) error {
	fmt.Printf("[Legacy] Refunding transaction %d\n", transactionID)
	return nil
}

// ============================================================================
// 适配器实现
// ============================================================================

// LegacyPaymentAdapter 将 LegacyPaymentSystem 适配到 PaymentProcessor 接口。
//
// 设计决策：
//  1. 使用组合而非继承（Go 没有继承）
//  2. 嵌入被适配的对象
//  3. 在适配方法中处理参数转换
type LegacyPaymentAdapter struct {
	legacy *LegacyPaymentSystem
}

// NewLegacyPaymentAdapter 创建适配器实例
func NewLegacyPaymentAdapter(legacy *LegacyPaymentSystem) *LegacyPaymentAdapter {
	return &LegacyPaymentAdapter{legacy: legacy}
}

// ProcessPayment 实现 PaymentProcessor 接口
func (a *LegacyPaymentAdapter) ProcessPayment(amount float64, currency string) (string, error) {
	// 转换参数：元 -> 分
	amountInCents := int64(amount * 100)

	// 调用旧系统方法
	txID, err := a.legacy.MakeCharge(amountInCents, currency)
	if err != nil {
		return "", err
	}

	// 转换返回值：int64 -> string
	return fmt.Sprintf("TX-%d", txID), nil
}

// RefundPayment 实现 PaymentProcessor 接口
func (a *LegacyPaymentAdapter) RefundPayment(transactionID string) error {
	// 解析 transactionID（去掉前缀 TX-）
	var txID int64
	_, err := fmt.Sscanf(transactionID, "TX-%d", &txID)
	if err != nil {
		return fmt.Errorf("invalid transaction ID: %s", transactionID)
	}

	return a.legacy.DoRefund(txID)
}

// ============================================================================
// 场景：IO 适配（io.Reader 适配器）
// ============================================================================

// StringReader 是一个简单的字符串读取器（模拟旧接口）
type StringReader struct {
	data string
	pos  int
}

// NewStringReader 创建字符串读取器
func NewStringReader(data string) *StringReader {
	return &StringReader{data: data}
}

// ReadString 读取字符串（返回 string，不符合 io.Reader）
func (s *StringReader) ReadString(n int) (string, error) {
	if s.pos >= len(s.data) {
		return "", io.EOF
	}
	end := s.pos + n
	if end > len(s.data) {
		end = len(s.data)
	}
	result := s.data[s.pos:end]
	s.pos = end
	return result, nil
}

// StringReaderAdapter 将 StringReader 适配为 io.Reader
type StringReaderAdapter struct {
	reader *StringReader
	buffer []byte
}

// NewStringReaderAdapter 创建适配器
func NewStringReaderAdapter(reader *StringReader) *StringReaderAdapter {
	return &StringReaderAdapter{
		reader: reader,
		buffer: make([]byte, 0, 1024),
	}
}

// Read 实现 io.Reader 接口
func (a *StringReaderAdapter) Read(p []byte) (n int, err error) {
	// 如果 buffer 为空，从 StringReader 读取
	if len(a.buffer) == 0 {
		str, err := a.reader.ReadString(1024)
		if err != nil && err != io.EOF {
			return 0, err
		}
		a.buffer = []byte(str)
		if err == io.EOF && len(a.buffer) == 0 {
			return 0, io.EOF
		}
	}

	// 复制数据到 p
	n = copy(p, a.buffer)
	a.buffer = a.buffer[n:]

	// 检查是否还有数据
	if len(a.buffer) == 0 && len(a.reader.data) == a.reader.pos {
		return n, io.EOF
	}

	return n, nil
}

// ============================================================================
// 隐式接口适配（Go 特有）
// ============================================================================

// Cache 是缓存接口
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}) error
	Delete(key string)
}

// ExternalCache 是外部缓存库（第三方，无法修改）
type ExternalCache struct {
	data map[string]interface{}
}

// NewExternalCache 创建外部缓存
func NewExternalCache() *ExternalCache {
	return &ExternalCache{data: make(map[string]interface{})}
}

// Fetch 获取值（方法名不同）
func (e *ExternalCache) Fetch(key string) (interface{}, bool) {
	val, ok := e.data[key]
	return val, ok
}

// Store 存储值（方法名不同）
func (e *ExternalCache) Store(key string, value interface{}) {
	e.data[key] = value
}

// Remove 删除值（方法名不同）
func (e *ExternalCache) Remove(key string) {
	delete(e.data, key)
}

// ExternalCacheAdapter 适配器
// 注意：由于 Go 的隐式接口，这个适配器非常简单
type ExternalCacheAdapter struct {
	cache *ExternalCache
}

// NewExternalCacheAdapter 创建适配器
func NewExternalCacheAdapter(cache *ExternalCache) *ExternalCacheAdapter {
	return &ExternalCacheAdapter{cache: cache}
}

// Get 实现 Cache 接口
func (a *ExternalCacheAdapter) Get(key string) (interface{}, bool) {
	return a.cache.Fetch(key)
}

// Set 实现 Cache 接口
func (a *ExternalCacheAdapter) Set(key string, value interface{}) error {
	a.cache.Store(key, value)
	return nil
}

// Delete 实现 Cache 接口
func (a *ExternalCacheAdapter) Delete(key string) {
	a.cache.Remove(key)
}

// ============================================================================
// 快速适配（Go 特有技巧）
// ============================================================================

// QuickCacheAdapter 是一个轻量级适配器，使用内嵌类型和闭包实现快速适配。
// 这是 Go 特有的简洁适配方式，不需要定义完整的适配器类型。
type QuickCacheAdapter struct {
	external *ExternalCache
}

// NewQuickCacheAdapter 创建快速适配器
func NewQuickCacheAdapter(external *ExternalCache) *QuickCacheAdapter {
	return &QuickCacheAdapter{external: external}
}

// Get 实现 Cache 接口
func (q *QuickCacheAdapter) Get(key string) (interface{}, bool) {
	return q.external.Fetch(key)
}

// Set 实现 Cache 接口
func (q *QuickCacheAdapter) Set(key string, value interface{}) error {
	q.external.Store(key, value)
	return nil
}

// Delete 实现 Cache 接口
func (q *QuickCacheAdapter) Delete(key string) {
	q.external.Remove(key)
}

// ============================================================================
// 多接口适配
// ============================================================================

// Storage 是存储接口
type Storage interface {
	Save(data []byte) (string, error)
	Load(id string) ([]byte, error)
}

// FileStorage 是文件存储实现
type FileStorage struct {
	basePath string
}

// NewFileStorage 创建文件存储
func NewFileStorage(basePath string) *FileStorage {
	return &FileStorage{basePath: basePath}
}

// Save 实现 Storage 接口
func (f *FileStorage) Save(data []byte) (string, error) {
	id := fmt.Sprintf("file-%d", time.Now().UnixNano())
	// 实际项目中写入文件...
	return id, nil
}

// Load 实现 Storage 接口
func (f *FileStorage) Load(id string) ([]byte, error) {
	// 实际项目中读取文件...
	return []byte("file data"), nil
}

// CloudStorage 是云存储（第三方库）
type CloudStorage struct {
	bucket string
}

// NewCloudStorage 创建云存储
func NewCloudStorage(bucket string) *CloudStorage {
	return &CloudStorage{bucket: bucket}
}

// Upload 上传（方法名不同）
func (c *CloudStorage) Upload(content []byte) (string, error) {
	id := fmt.Sprintf("cloud-%d", time.Now().UnixNano())
	fmt.Printf("[Cloud] Uploading to bucket %s\n", c.bucket)
	return id, nil
}

// Download 下载（方法名不同）
func (c *CloudStorage) Download(objectID string) ([]byte, error) {
	fmt.Printf("[Cloud] Downloading %s from bucket %s\n", objectID, c.bucket)
	return []byte("cloud data"), nil
}

// CloudStorageAdapter 云存储适配器
type CloudStorageAdapter struct {
	cloud *CloudStorage
}

// NewCloudStorageAdapter 创建适配器
func NewCloudStorageAdapter(cloud *CloudStorage) *CloudStorageAdapter {
	return &CloudStorageAdapter{cloud: cloud}
}

// Save 实现 Storage 接口
func (a *CloudStorageAdapter) Save(data []byte) (string, error) {
	return a.cloud.Upload(data)
}

// Load 实现 Storage 接口
func (a *CloudStorageAdapter) Load(id string) ([]byte, error) {
	return a.cloud.Download(id)
}

// ============================================================================
// 双向适配器
// ============================================================================

// EuropeanPlug 是欧洲插头接口
type EuropeanPlug interface {
	ConnectToEuropeanSocket() string
}

// USPlug 是美国插头接口
type USPlug interface {
	ConnectToUSSocket() string
}

// EuropeanDevice 是欧洲设备
type EuropeanDevice struct{}

// ConnectToEuropeanSocket 实现 EuropeanPlug
func (e *EuropeanDevice) ConnectToEuropeanSocket() string {
	return "Connected to European socket (220V)"
}

// USDevice 是美国设备
type USDevice struct{}

// ConnectToUSSocket 实现 USPlug
func (u *USDevice) ConnectToUSSocket() string {
	return "Connected to US socket (110V)"
}

// EuropeanToUSAdapter 欧洲插头转美国插座
type EuropeanToUSAdapter struct {
	device EuropeanPlug
}

// NewEuropeanToUSAdapter 创建适配器
func NewEuropeanToUSAdapter(device EuropeanPlug) *EuropeanToUSAdapter {
	return &EuropeanToUSAdapter{device: device}
}

// ConnectToUSSocket 实现 USPlug
func (a *EuropeanToUSAdapter) ConnectToUSSocket() string {
	return "Adapter converting US 110V to European 220V -> " + a.device.ConnectToEuropeanSocket()
}

// USToEuropeanAdapter 美国插头转欧洲插座
type USToEuropeanAdapter struct {
	device USPlug
}

// NewUSToEuropeanAdapter 创建适配器
func NewUSToEuropeanAdapter(device USPlug) *USToEuropeanAdapter {
	return &USToEuropeanAdapter{device: device}
}

// ConnectToEuropeanSocket 实现 EuropeanPlug
func (a *USToEuropeanAdapter) ConnectToEuropeanSocket() string {
	return "Adapter converting European 220V to US 110V -> " + a.device.ConnectToUSSocket()
}
