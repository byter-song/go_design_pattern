// Package builder 展示了 Go 语言中实现建造者模式的惯用方法，
// 以及 Go 社区更推荐的"函数式选项模式(Functional Options Pattern)"。
//
// 建造者模式将一个复杂对象的构建与其表示分离，使得同样的构建过程可以创建不同的表示。
//
// 关键概念：
//   - 分步骤构建复杂对象
//   - 链式调用提供流畅 API
//   - 必填参数与可选参数的区分
//
// Go 语言特色：
//   - 函数式选项模式是 Go 社区的首选方案
//   - 利用闭包实现灵活的配置
//   - 零值和默认值简化 API 设计
package builder

import (
	"fmt"
	"time"
)

// ============================================================================
// 产品：HTTP 服务器配置（复杂对象示例）
// ============================================================================

// Server 是一个 HTTP 服务器配置结构体。
// 包含多个配置项，其中一些是必填的，一些是可选的。
//
// 设计决策：
//   1. 字段导出（大写），但推荐通过建造者创建
//   2. 包含合理的默认值
//   3. 必填字段使用指针或特殊值表示未设置
type Server struct {
	// 必填字段
	Host string
	Port int

	// 可选字段（有默认值）
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	MaxHeaderBytes  int
	MaxConnections  int

	// 可选字段（布尔标志）
	TLS      bool
	Debug    bool
	Compress bool

	// 复杂嵌套配置
	TLSConfig *TLSConfiguration
	Logger    *LoggerConfiguration
}

// TLSConfiguration TLS 配置
type TLSConfiguration struct {
	CertFile string
	KeyFile  string
}

// LoggerConfiguration 日志配置
type LoggerConfiguration struct {
	Level  string
	Output string
}

// Validate 验证服务器配置是否有效
func (s *Server) Validate() error {
	if s.Host == "" {
		return fmt.Errorf("host is required")
	}
	if s.Port <= 0 || s.Port > 65535 {
		return fmt.Errorf("invalid port: %d", s.Port)
	}
	if s.TLS {
		if s.TLSConfig == nil || s.TLSConfig.CertFile == "" || s.TLSConfig.KeyFile == "" {
			return fmt.Errorf("TLS config is required when TLS is enabled")
		}
	}
	return nil
}

// ============================================================================
// 传统建造者模式（Traditional Builder）
// ============================================================================

// ServerBuilder 是服务器配置的建造者。
//
// 传统建造者模式的特点：
//   - 持有产品对象的引用
//   - 提供链式调用的设置方法
//   - 提供 Build 方法返回最终产品
type ServerBuilder struct {
	server *Server
}

// NewServerBuilder 创建一个新的服务器建造者
//
// 初始化时设置默认值
func NewServerBuilder() *ServerBuilder {
	return &ServerBuilder{
		server: &Server{
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   30 * time.Second,
			MaxHeaderBytes: 1 << 20, // 1MB
			MaxConnections: 100,
			Debug:          false,
			Compress:       true,
			Logger: &LoggerConfiguration{
				Level:  "info",
				Output: "stdout",
			},
		},
	}
}

// WithHost 设置主机地址（必填）
func (b *ServerBuilder) WithHost(host string) *ServerBuilder {
	b.server.Host = host
	return b
}

// WithPort 设置端口（必填）
func (b *ServerBuilder) WithPort(port int) *ServerBuilder {
	b.server.Port = port
	return b
}

// WithReadTimeout 设置读取超时
func (b *ServerBuilder) WithReadTimeout(timeout time.Duration) *ServerBuilder {
	b.server.ReadTimeout = timeout
	return b
}

// WithWriteTimeout 设置写入超时
func (b *ServerBuilder) WithWriteTimeout(timeout time.Duration) *ServerBuilder {
	b.server.WriteTimeout = timeout
	return b
}

// WithMaxConnections 设置最大连接数
func (b *ServerBuilder) WithMaxConnections(max int) *ServerBuilder {
	b.server.MaxConnections = max
	return b
}

// WithMaxHeaderBytes 设置最大请求头大小
func (b *ServerBuilder) WithMaxHeaderBytes(max int) *ServerBuilder {
	b.server.MaxHeaderBytes = max
	return b
}

// WithTimeout 同时设置读取和写入超时
func (b *ServerBuilder) WithTimeout(timeout time.Duration) *ServerBuilder {
	b.server.ReadTimeout = timeout
	b.server.WriteTimeout = timeout
	return b
}

// WithTLS 启用 TLS 并配置证书
func (b *ServerBuilder) WithTLS(certFile, keyFile string) *ServerBuilder {
	b.server.TLS = true
	b.server.TLSConfig = &TLSConfiguration{
		CertFile: certFile,
		KeyFile:  keyFile,
	}
	return b
}

// WithDebug 启用调试模式
func (b *ServerBuilder) WithDebug() *ServerBuilder {
	b.server.Debug = true
	return b
}

// WithoutCompression 禁用压缩（默认启用）
func (b *ServerBuilder) WithoutCompression() *ServerBuilder {
	b.server.Compress = false
	return b
}

// WithLogger 配置日志
func (b *ServerBuilder) WithLogger(level, output string) *ServerBuilder {
	b.server.Logger = &LoggerConfiguration{
		Level:  level,
		Output: output,
	}
	return b
}

// Build 构建并返回服务器配置
//
// 在返回前进行验证，确保必填字段已设置
func (b *ServerBuilder) Build() (*Server, error) {
	if err := b.server.Validate(); err != nil {
		return nil, fmt.Errorf("build failed: %w", err)
	}
	return b.server, nil
}

// ============================================================================
// 函数式选项模式（Functional Options Pattern）- Go 推荐
// ============================================================================

// Option 是函数式选项类型，这是 Go 社区广泛采用的模式。
//
// 相比传统建造者，函数式选项模式的优势：
//   1. 更简洁的 API
//   2. 易于扩展，新增选项不需要修改结构体
//   3. 可以组合和复用选项
//   4. 更符合 Go 的函数式编程风格
//
// 这是 Go 标准库（如 grpc、uber-go/zap）采用的方式。
type Option func(*Server)

// NewServer 使用函数式选项模式创建服务器。
//
// 参数：
//   - host: 主机地址（必填）
//   - port: 端口（必填）
//   - opts: 可选配置选项
//
// 使用示例：
//
//	server, err := NewServer("localhost", 8080,
//	    WithTimeout(10*time.Second),
//	    WithTLS("cert.pem", "key.pem"),
//	    WithDebug(),
//	)
func NewServer(host string, port int, opts ...Option) (*Server, error) {
	// 创建带有默认值的服务器
	s := &Server{
		Host: host,
		Port: port,
		// 默认值
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
		MaxConnections: 100,
		Compress:       true,
		Logger: &LoggerConfiguration{
			Level:  "info",
			Output: "stdout",
		},
	}

	// 应用所有选项
	for _, opt := range opts {
		opt(s)
	}

	// 验证
	if err := s.Validate(); err != nil {
		return nil, err
	}

	return s, nil
}

// 以下是所有可用的选项函数

// WithReadTimeout 设置读取超时
func WithReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.ReadTimeout = timeout
	}
}

// WithWriteTimeout 设置写入超时
func WithWriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.WriteTimeout = timeout
	}
}

// WithTimeout 同时设置读取和写入超时
func WithTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.ReadTimeout = timeout
		s.WriteTimeout = timeout
	}
}

// WithMaxConnections 设置最大连接数
func WithMaxConnections(max int) Option {
	return func(s *Server) {
		s.MaxConnections = max
	}
}

// WithMaxHeaderBytes 设置最大请求头大小
func WithMaxHeaderBytes(max int) Option {
	return func(s *Server) {
		s.MaxHeaderBytes = max
	}
}

// WithTLS 启用 TLS
func WithTLS(certFile, keyFile string) Option {
	return func(s *Server) {
		s.TLS = true
		s.TLSConfig = &TLSConfiguration{
			CertFile: certFile,
			KeyFile:  keyFile,
		}
	}
}

// WithDebug 启用调试模式
func WithDebug() Option {
	return func(s *Server) {
		s.Debug = true
	}
}

// WithoutCompression 禁用压缩
func WithoutCompression() Option {
	return func(s *Server) {
		s.Compress = false
	}
}

// WithLogger 配置日志
func WithLogger(level, output string) Option {
	return func(s *Server) {
		s.Logger = &LoggerConfiguration{
			Level:  level,
			Output: output,
		}
	}
}

// ============================================================================
// 选项组合（高级用法）
// ============================================================================

// CommonOptions 返回常用的选项组合
//
// 使用场景：
//   - 多个服务共享相同的配置
//   - 预设配置模板（开发环境、生产环境）
func CommonOptions() []Option {
	return []Option{
		WithTimeout(30 * time.Second),
		WithMaxConnections(100),
		WithLogger("info", "stdout"),
	}
}

// ProductionOptions 返回生产环境推荐的选项
func ProductionOptions() []Option {
	return []Option{
		WithTimeout(60 * time.Second),
		WithMaxConnections(1000),
		WithMaxHeaderBytes(1 << 20),
		WithoutCompression(), // 生产环境可能由反向代理处理压缩
		WithLogger("warn", "file"),
	}
}

// DevelopmentOptions 返回开发环境推荐的选项
func DevelopmentOptions() []Option {
	return []Option{
		WithTimeout(10 * time.Second),
		WithMaxConnections(10),
		WithDebug(),
		WithLogger("debug", "stdout"),
	}
}

// MergeOptions 合并多个选项切片
//
// 使用场景：组合预设选项和自定义选项
//
//	opts := MergeOptions(ProductionOptions(), WithDebug())
//	server, _ := NewServer("localhost", 8080, opts...)
func MergeOptions(optSlices ...[]Option) []Option {
	var result []Option
	for _, opts := range optSlices {
		result = append(result, opts...)
	}
	return result
}

// ============================================================================
// 分步建造者（Step Builder）- 强制设置顺序
// ============================================================================

// StepBuilderHost 是第一步：必须设置 Host
type StepBuilderHost struct {
	server *Server
}

// StepBuilderPort 是第二步：必须设置 Port
type StepBuilderPort struct {
	server *Server
}

// StepBuilderFinal 是最后一步：可以设置可选参数
type StepBuilderFinal struct {
	server *Server
}

// NewStepBuilder 开始分步建造
func NewStepBuilder() *StepBuilderHost {
	return &StepBuilderHost{
		server: &Server{
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   30 * time.Second,
			MaxHeaderBytes: 1 << 20,
			MaxConnections: 100,
		},
	}
}

// WithHost 设置 Host，返回下一步建造者
func (b *StepBuilderHost) WithHost(host string) *StepBuilderPort {
	b.server.Host = host
	return &StepBuilderPort{server: b.server}
}

// WithPort 设置 Port，返回最终建造者
func (b *StepBuilderPort) WithPort(port int) *StepBuilderFinal {
	b.server.Port = port
	return &StepBuilderFinal{server: b.server}
}

// WithTimeout 设置超时（可选）
func (b *StepBuilderFinal) WithTimeout(timeout time.Duration) *StepBuilderFinal {
	b.server.ReadTimeout = timeout
	b.server.WriteTimeout = timeout
	return b
}

// WithTLS 启用 TLS（可选）
func (b *StepBuilderFinal) WithTLS(certFile, keyFile string) *StepBuilderFinal {
	b.server.TLS = true
	b.server.TLSConfig = &TLSConfiguration{
		CertFile: certFile,
		KeyFile:  keyFile,
	}
	return b
}

// Build 完成构建
func (b *StepBuilderFinal) Build() (*Server, error) {
	if err := b.server.Validate(); err != nil {
		return nil, err
	}
	return b.server, nil
}

// ============================================================================
// 验证：编译时检查接口实现
// ============================================================================

// 这些变量用于编译时检查，确保代码正确性
// 如果类型没有实现预期的方法，编译会失败
var (
	// 确保所有建造者都能构建出 Server
	_ = NewServerBuilder().Build
	_ = NewStepBuilder().WithHost
	_ = NewServer
)
