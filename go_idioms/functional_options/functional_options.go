// Package functional_options 展示了 Go 语言中函数式选项模式（Functional Options Pattern）的实现
//
// 这是 Go 社区广泛认可的初始化复杂结构体的最佳实践，由 Rob Pike 推广，
// 被广泛应用于标准库（如 grpc、zap、viper 等）和第三方库中。
//
// 设计思想：
// 1. 使用函数作为配置参数，提供类型安全、可扩展的 API
// 2. 避免使用大量的构造函数或复杂的配置结构体
// 3. 支持默认值的优雅设置
// 4. 允许调用者只设置关心的选项
//
// 核心优势：
// - 向后兼容：新增选项不会破坏现有代码
// - 自文档化：调用代码清晰展示配置意图
// - 类型安全：编译时检查，避免使用 interface{} 或 map[string]interface{}
// - 可扩展：轻松添加新选项而不影响现有代码
// - 默认值友好：未设置的选项使用合理的默认值
//
// 对比其他方案：
//
// 方案1: 大量构造函数（不推荐）
//   NewServer(addr string) *Server
//   NewServerWithPort(addr string, port int) *Server
//   NewServerWithPortAndTimeout(addr string, port int, timeout time.Duration) *Server
//   // 组合爆炸问题！
//
// 方案2: 配置结构体（不够优雅）
//   config := &Config{Addr: "localhost", Port: 8080}
//   server := NewServer(config)
//   // 需要引入额外的 Config 类型，且无法强制设置必填参数
//
// 方案3: 函数式选项（推荐 ✅）
//   server, err := NewServer("localhost", WithPort(8080), WithTimeout(30*time.Second))
//   // 清晰、类型安全、可扩展
package functional_options

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

// Option 定义函数式选项类型
// 这是模式的核心：一个接收 *Server 并修改其配置的函数类型
type Option func(*Server)

// Server 是一个 HTTP 服务器结构体
// 包含各种可配置选项，展示函数式选项模式的实际应用
type Server struct {
	// 网络配置
	host string
	port int

	// 超时配置
	readTimeout  time.Duration
	writeTimeout time.Duration
	idleTimeout  time.Duration

	// 性能配置
	maxHeaderBytes    int
	maxConnections    int
	readHeaderTimeout time.Duration

	// 功能开关
	enableLogging bool
	enableMetrics bool
	enableCORS    bool
	enableTLS     bool

	// TLS 配置
	certFile string
	keyFile  string

	// 处理器和中间件
	handler     http.Handler
	middlewares []func(http.Handler) http.Handler

	// 内部状态
	server   *http.Server
	listener net.Listener
	logger   *log.Logger
}

// NewServer 创建一个新的 Server 实例
// 必填参数直接作为函数参数，可选参数使用 ...Option 变长参数
//
// 设计决策：
// - host 是必填参数，因为服务器必须监听某个地址
// - 其他所有配置都通过 Option 函数设置
// - 使用合理的默认值，让最小化配置也能工作
//
// 使用示例：
//
//	// 最小化配置（使用所有默认值）
//	server, err := NewServer("localhost")
//
//	// 完整配置
//	server, err := NewServer("0.0.0.0",
//	    WithPort(8080),
//	    WithTimeout(30*time.Second),
//	    WithMaxConnections(1000),
//	    WithLogging(true),
//	    WithTLS("cert.pem", "key.pem"),
//	)
func NewServer(host string, opts ...Option) (*Server, error) {
	if host == "" {
		return nil, fmt.Errorf("host cannot be empty")
	}

	// 创建带有默认值的服务器实例
	// 这是函数式选项模式的另一个优点：默认值集中管理
	s := &Server{
		host: host,
		// 网络默认值
		port: 8080,

		// 超时默认值（合理的生产环境值）
		readTimeout:       15 * time.Second,
		writeTimeout:      15 * time.Second,
		idleTimeout:       60 * time.Second,
		readHeaderTimeout: 10 * time.Second,

		// 性能默认值
		maxHeaderBytes: 1 << 20, // 1 MB
		maxConnections: 100,

		// 功能开关默认值
		enableLogging: true,
		enableMetrics: false,
		enableCORS:    false,
		enableTLS:     false,

		// 默认处理器
		handler: http.DefaultServeMux,

		// 中间件切片
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}

	// 应用所有选项函数
	// 这是模式的核心逻辑：按顺序执行每个 Option 函数
	for _, opt := range opts {
		opt(s)
	}

	// 验证配置
	if err := s.validate(); err != nil {
		return nil, fmt.Errorf("invalid server configuration: %w", err)
	}

	// 初始化 logger
	if s.logger == nil {
		s.logger = log.Default()
	}

	// 构建最终的 HTTP 服务器
	s.buildHTTPServer()

	return s, nil
}

// validate 验证服务器配置的有效性
func (s *Server) validate() error {
	if s.port < 1 || s.port > 65535 {
		return fmt.Errorf("invalid port: %d", s.port)
	}

	if s.enableTLS {
		if s.certFile == "" || s.keyFile == "" {
			return fmt.Errorf("TLS enabled but certFile or keyFile not provided")
		}
	}

	return nil
}

// buildHTTPServer 构建底层的 http.Server
func (s *Server) buildHTTPServer() {
	// 包装处理器（添加中间件）
	handler := s.handler
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		handler = s.middlewares[i](handler)
	}

	// 如果需要日志，包装日志中间件
	if s.enableLogging {
		handler = s.loggingMiddleware(handler)
	}

	s.server = &http.Server{
		Addr:              fmt.Sprintf("%s:%d", s.host, s.port),
		Handler:           handler,
		ReadTimeout:       s.readTimeout,
		WriteTimeout:      s.writeTimeout,
		IdleTimeout:       s.idleTimeout,
		ReadHeaderTimeout: s.readHeaderTimeout,
		MaxHeaderBytes:    s.maxHeaderBytes,
	}
}

// ==================== Option 函数实现 ====================

// WithPort 设置服务器端口
// 这是最基本的选项函数示例：接收一个值，返回修改该值的函数
func WithPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

// WithTimeout 设置所有超时（简化配置）
// 展示如何用一个选项同时修改多个字段
func WithTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.readTimeout = timeout
		s.writeTimeout = timeout
		s.idleTimeout = timeout * 4 // idle 通常更长
	}
}

// WithReadTimeout 设置读取超时
// 展示如何提供细粒度的选项
func WithReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.readTimeout = timeout
	}
}

// WithWriteTimeout 设置写入超时
func WithWriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.writeTimeout = timeout
	}
}

// WithIdleTimeout 设置空闲连接超时
func WithIdleTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.idleTimeout = timeout
	}
}

// WithMaxConnections 设置最大连接数
func WithMaxConnections(max int) Option {
	return func(s *Server) {
		s.maxConnections = max
	}
}

// WithMaxHeaderBytes 设置最大请求头大小
func WithMaxHeaderBytes(max int) Option {
	return func(s *Server) {
		s.maxHeaderBytes = max
	}
}

// WithLogging 启用或禁用日志
func WithLogging(enable bool) Option {
	return func(s *Server) {
		s.enableLogging = enable
	}
}

// WithMetrics 启用指标收集
func WithMetrics(enable bool) Option {
	return func(s *Server) {
		s.enableMetrics = enable
	}
}

// WithCORS 启用 CORS 支持
func WithCORS(enable bool) Option {
	return func(s *Server) {
		s.enableCORS = enable
	}
}

// WithTLS 启用 TLS 并指定证书
// 展示选项之间的依赖关系处理
func WithTLS(certFile, keyFile string) Option {
	return func(s *Server) {
		s.enableTLS = true
		s.certFile = certFile
		s.keyFile = keyFile
	}
}

// WithHandler 设置自定义 HTTP 处理器
func WithHandler(handler http.Handler) Option {
	return func(s *Server) {
		s.handler = handler
	}
}

// WithMiddleware 添加中间件
// 展示如何处理切片类型的选项（追加而非替换）
func WithMiddleware(mw func(http.Handler) http.Handler) Option {
	return func(s *Server) {
		s.middlewares = append(s.middlewares, mw)
	}
}

// WithLogger 设置自定义 logger
func WithLogger(logger *log.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

// ==================== 服务器方法 ====================

// Start 启动服务器
func (s *Server) Start() error {
	if s.enableTLS {
		s.logger.Printf("Starting HTTPS server on %s:%d", s.host, s.port)
		return s.server.ListenAndServeTLS(s.certFile, s.keyFile)
	}
	s.logger.Printf("Starting HTTP server on %s:%d", s.host, s.port)
	return s.server.ListenAndServe()
}

// StartAsync 异步启动服务器
func (s *Server) StartAsync() error {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		return err
	}
	s.listener = ln

	go func() {
		if s.enableTLS {
			s.server.ServeTLS(ln, s.certFile, s.keyFile)
		} else {
			s.server.Serve(ln)
		}
	}()

	return nil
}

// Stop 优雅停止服务器
func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// Addr 返回服务器地址
func (s *Server) Addr() string {
	return fmt.Sprintf("%s:%d", s.host, s.port)
}

// loggingMiddleware 日志中间件
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		s.logger.Printf("[%s] %s %s - %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}

// ==================== 高级选项组合 ====================

// ProductionDefaults 返回生产环境推荐的默认配置
// 展示如何组合多个选项形成预设配置
func ProductionDefaults() []Option {
	return []Option{
		WithTimeout(30 * time.Second),
		WithMaxConnections(1000),
		WithMaxHeaderBytes(1 << 20), // 1MB
		WithLogging(true),
		WithMetrics(true),
	}
}

// DevelopmentDefaults 返回开发环境推荐的默认配置
func DevelopmentDefaults() []Option {
	return []Option{
		WithTimeout(5 * time.Minute), // 开发时更长的超时方便调试
		WithMaxConnections(100),
		WithLogging(true),
		WithCORS(true),
	}
}

// WithSecureHeaders 添加安全相关的 HTTP 头中间件
// 展示如何创建功能性的选项
func WithSecureHeaders() Option {
	return WithMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			next.ServeHTTP(w, r)
		})
	})
}

// ==================== 另一个示例：Database 连接配置 ====================

// DBConfig 数据库配置结构体
type DBConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string

	// 连接池配置
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration

	// 超时配置
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	// 其他选项
	SSLMode  string
	Timezone string
}

// DBOption 数据库配置选项类型
type DBOption func(*DBConfig)

// NewDBConfig 创建数据库配置
func NewDBConfig(host, database string, opts ...DBOption) (*DBConfig, error) {
	if host == "" || database == "" {
		return nil, fmt.Errorf("host and database are required")
	}

	cfg := &DBConfig{
		Host:            host,
		Port:            5432, // PostgreSQL 默认端口
		Database:        database,
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: time.Minute * 30,
		DialTimeout:     time.Second * 10,
		ReadTimeout:     time.Second * 30,
		WriteTimeout:    time.Second * 30,
		SSLMode:         "require",
		Timezone:        "UTC",
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg, nil
}

// DSN 生成数据源名称（Data Source Name）
func (c *DBConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		c.Host, c.Port, c.Database, c.Username, c.Password, c.SSLMode)
}

// DB 选项函数
func WithDBPort(port int) DBOption {
	return func(c *DBConfig) {
		c.Port = port
	}
}

func WithDBCredentials(username, password string) DBOption {
	return func(c *DBConfig) {
		c.Username = username
		c.Password = password
	}
}

func WithDBPoolSize(maxOpen, maxIdle int) DBOption {
	return func(c *DBConfig) {
		c.MaxOpenConns = maxOpen
		c.MaxIdleConns = maxIdle
	}
}

func WithDBTimeout(dial, read, write time.Duration) DBOption {
	return func(c *DBConfig) {
		c.DialTimeout = dial
		c.ReadTimeout = read
		c.WriteTimeout = write
	}
}

func WithDBSSLMode(mode string) DBOption {
	return func(c *DBConfig) {
		c.SSLMode = mode
	}
}
