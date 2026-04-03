package functional_options

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"
)

// TestNewServerWithDefaults 测试使用默认配置创建服务器
func TestNewServerWithDefaults(t *testing.T) {
	server, err := NewServer("localhost")
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// 验证默认值
	if server.host != "localhost" {
		t.Errorf("Expected host 'localhost', got '%s'", server.host)
	}
	if server.port != 8080 {
		t.Errorf("Expected default port 8080, got %d", server.port)
	}
	if server.readTimeout != 15*time.Second {
		t.Errorf("Expected default read timeout 15s, got %v", server.readTimeout)
	}
	if !server.enableLogging {
		t.Error("Expected logging to be enabled by default")
	}
	if server.enableTLS {
		t.Error("Expected TLS to be disabled by default")
	}
}

// TestNewServerWithOptions 测试使用选项创建服务器
func TestNewServerWithOptions(t *testing.T) {
	server, err := NewServer("0.0.0.0",
		WithPort(9090),
		WithReadTimeout(30*time.Second),
		WithWriteTimeout(45*time.Second),
		WithMaxConnections(500),
		WithLogging(false),
		WithMetrics(true),
	)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// 验证自定义值
	if server.port != 9090 {
		t.Errorf("Expected port 9090, got %d", server.port)
	}
	if server.readTimeout != 30*time.Second {
		t.Errorf("Expected read timeout 30s, got %v", server.readTimeout)
	}
	if server.writeTimeout != 45*time.Second {
		t.Errorf("Expected write timeout 45s, got %v", server.writeTimeout)
	}
	if server.maxConnections != 500 {
		t.Errorf("Expected max connections 500, got %d", server.maxConnections)
	}
	if server.enableLogging {
		t.Error("Expected logging to be disabled")
	}
	if !server.enableMetrics {
		t.Error("Expected metrics to be enabled")
	}
}

// TestWithTimeout 测试统一超时设置
func TestWithTimeout(t *testing.T) {
	server, err := NewServer("localhost", WithTimeout(20*time.Second))
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	if server.readTimeout != 20*time.Second {
		t.Errorf("Expected read timeout 20s, got %v", server.readTimeout)
	}
	if server.writeTimeout != 20*time.Second {
		t.Errorf("Expected write timeout 20s, got %v", server.writeTimeout)
	}
	// idle timeout 应该是 timeout 的 4 倍
	if server.idleTimeout != 80*time.Second {
		t.Errorf("Expected idle timeout 80s, got %v", server.idleTimeout)
	}
}

// TestWithTLS 测试 TLS 配置
func TestWithTLS(t *testing.T) {
	server, err := NewServer("localhost",
		WithPort(8443),
		WithTLS("cert.pem", "key.pem"),
	)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	if !server.enableTLS {
		t.Error("Expected TLS to be enabled")
	}
	if server.certFile != "cert.pem" {
		t.Errorf("Expected cert file 'cert.pem', got '%s'", server.certFile)
	}
	if server.keyFile != "key.pem" {
		t.Errorf("Expected key file 'key.pem', got '%s'", server.keyFile)
	}
	if server.port != 8443 {
		t.Errorf("Expected port 8443, got %d", server.port)
	}
}

// TestWithTLSValidation 测试 TLS 配置验证
func TestWithTLSValidation(t *testing.T) {
	// 启用 TLS 但不提供证书应该报错
	_, err := NewServer("localhost",
		WithTLS("", ""),
	)
	if err == nil {
		t.Error("Expected error when TLS enabled but no cert provided")
	}
}

// TestWithMiddleware 测试中间件添加
func TestWithMiddleware(t *testing.T) {
	customHeader := "X-Custom-Header"
	customValue := "test-value"

	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(customHeader, customValue)
			next.ServeHTTP(w, r)
		})
	}

	server, err := NewServer("localhost",
		WithMiddleware(middleware),
	)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	if len(server.middlewares) != 1 {
		t.Errorf("Expected 1 middleware, got %d", len(server.middlewares))
	}
}

// TestWithMultipleMiddlewares 测试多个中间件
func TestWithMultipleMiddlewares(t *testing.T) {
	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-MW1", "1")
			next.ServeHTTP(w, r)
		})
	}

	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-MW2", "2")
			next.ServeHTTP(w, r)
		})
	}

	server, err := NewServer("localhost",
		WithMiddleware(mw1),
		WithMiddleware(mw2),
	)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	if len(server.middlewares) != 2 {
		t.Errorf("Expected 2 middlewares, got %d", len(server.middlewares))
	}
}

// TestNewServerValidation 测试服务器配置验证
func TestNewServerValidation(t *testing.T) {
	// 空 host 应该报错
	_, err := NewServer("")
	if err == nil {
		t.Error("Expected error for empty host")
	}

	// 无效端口应该报错
	_, err = NewServer("localhost", WithPort(0))
	if err == nil {
		t.Error("Expected error for port 0")
	}

	_, err = NewServer("localhost", WithPort(70000))
	if err == nil {
		t.Error("Expected error for port > 65535")
	}
}

// TestProductionDefaults 测试生产环境预设
func TestProductionDefaults(t *testing.T) {
	server, err := NewServer("0.0.0.0", ProductionDefaults()...)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	if server.readTimeout != 30*time.Second {
		t.Errorf("Expected production read timeout 30s, got %v", server.readTimeout)
	}
	if server.maxConnections != 1000 {
		t.Errorf("Expected production max connections 1000, got %d", server.maxConnections)
	}
	if !server.enableMetrics {
		t.Error("Expected metrics to be enabled in production")
	}
}

// TestDevelopmentDefaults 测试开发环境预设
func TestDevelopmentDefaults(t *testing.T) {
	server, err := NewServer("localhost", DevelopmentDefaults()...)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// 开发环境应该有更长的超时
	if server.readTimeout != 5*time.Minute {
		t.Errorf("Expected dev read timeout 5m, got %v", server.readTimeout)
	}
	if !server.enableCORS {
		t.Error("Expected CORS to be enabled in development")
	}
}

// TestOptionCombination 测试选项组合
func TestOptionCombination(t *testing.T) {
	// 先应用生产环境预设，再覆盖特定选项
	server, err := NewServer("0.0.0.0",
		ProductionDefaults()...,
	)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// 验证生产环境默认值
	if server.port != 8080 {
		t.Errorf("Expected default port 8080, got %d", server.port)
	}

	// 现在创建另一个服务器，覆盖端口
	server2, err := NewServer("0.0.0.0",
		append(ProductionDefaults(), WithPort(3000))...,
	)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	if server2.port != 3000 {
		t.Errorf("Expected overridden port 3000, got %d", server2.port)
	}
	// 其他生产环境设置应该保留
	if server2.maxConnections != 1000 {
		t.Errorf("Expected max connections 1000, got %d", server2.maxConnections)
	}
}

// TestWithSecureHeaders 测试安全头中间件
func TestWithSecureHeaders(t *testing.T) {
	server, err := NewServer("localhost", WithSecureHeaders())
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// 安全头中间件应该被添加到中间件列表
	if len(server.middlewares) != 1 {
		t.Errorf("Expected 1 middleware for secure headers, got %d", len(server.middlewares))
	}
}

// TestServerAddr 测试地址生成
func TestServerAddr(t *testing.T) {
	server, err := NewServer("localhost", WithPort(8080))
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	addr := server.Addr()
	expected := "localhost:8080"
	if addr != expected {
		t.Errorf("Expected address '%s', got '%s'", expected, addr)
	}
}

// TestWithHandler 测试自定义处理器
func TestWithHandler(t *testing.T) {
	customHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from custom handler"))
	})

	server, err := NewServer("localhost", WithHandler(customHandler))
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	if server.handler == nil {
		t.Error("Expected custom handler to be set")
	}
}

// TestWithLogger 测试自定义 logger
func TestWithLogger(t *testing.T) {
	customLogger := log.New(io.Discard, "[TEST] ", log.LstdFlags)

	server, err := NewServer("localhost", WithLogger(customLogger))
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	if server.logger != customLogger {
		t.Error("Expected custom logger to be set")
	}
}

// ==================== Database Config Tests ====================

// TestNewDBConfigWithDefaults 测试数据库配置默认值
func TestNewDBConfigWithDefaults(t *testing.T) {
	cfg, err := NewDBConfig("localhost", "mydb")
	if err != nil {
		t.Fatalf("Failed to create DB config: %v", err)
	}

	if cfg.Host != "localhost" {
		t.Errorf("Expected host 'localhost', got '%s'", cfg.Host)
	}
	if cfg.Database != "mydb" {
		t.Errorf("Expected database 'mydb', got '%s'", cfg.Database)
	}
	if cfg.Port != 5432 {
		t.Errorf("Expected default port 5432, got %d", cfg.Port)
	}
	if cfg.MaxOpenConns != 25 {
		t.Errorf("Expected default max open conns 25, got %d", cfg.MaxOpenConns)
	}
	if cfg.SSLMode != "require" {
		t.Errorf("Expected default SSL mode 'require', got '%s'", cfg.SSLMode)
	}
}

// TestNewDBConfigWithOptions 测试数据库配置选项
func TestNewDBConfigWithOptions(t *testing.T) {
	cfg, err := NewDBConfig("db.example.com", "production",
		WithDBPort(5433),
		WithDBCredentials("admin", "secret123"),
		WithDBPoolSize(50, 10),
		WithDBSSLMode("disable"),
	)
	if err != nil {
		t.Fatalf("Failed to create DB config: %v", err)
	}

	if cfg.Port != 5433 {
		t.Errorf("Expected port 5433, got %d", cfg.Port)
	}
	if cfg.Username != "admin" {
		t.Errorf("Expected username 'admin', got '%s'", cfg.Username)
	}
	if cfg.Password != "secret123" {
		t.Errorf("Expected password 'secret123', got '%s'", cfg.Password)
	}
	if cfg.MaxOpenConns != 50 {
		t.Errorf("Expected max open conns 50, got %d", cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns != 10 {
		t.Errorf("Expected max idle conns 10, got %d", cfg.MaxIdleConns)
	}
	if cfg.SSLMode != "disable" {
		t.Errorf("Expected SSL mode 'disable', got '%s'", cfg.SSLMode)
	}
}

// TestDBConfigDSN 测试 DSN 生成
func TestDBConfigDSN(t *testing.T) {
	cfg, err := NewDBConfig("localhost", "testdb",
		WithDBCredentials("user", "pass"),
		WithDBPort(5432),
	)
	if err != nil {
		t.Fatalf("Failed to create DB config: %v", err)
	}

	dsn := cfg.DSN()
	if !strings.Contains(dsn, "host=localhost") {
		t.Errorf("DSN should contain host, got: %s", dsn)
	}
	if !strings.Contains(dsn, "dbname=testdb") {
		t.Errorf("DSN should contain database, got: %s", dsn)
	}
	if !strings.Contains(dsn, "user=user") {
		t.Errorf("DSN should contain username, got: %s", dsn)
	}
}

// TestNewDBConfigValidation 测试数据库配置验证
func TestNewDBConfigValidation(t *testing.T) {
	_, err := NewDBConfig("", "mydb")
	if err == nil {
		t.Error("Expected error for empty host")
	}

	_, err = NewDBConfig("localhost", "")
	if err == nil {
		t.Error("Expected error for empty database")
	}
}

// TestWithDBTimeout 测试数据库超时配置
func TestWithDBTimeout(t *testing.T) {
	cfg, err := NewDBConfig("localhost", "mydb",
		WithDBTimeout(5*time.Second, 10*time.Second, 15*time.Second),
	)
	if err != nil {
		t.Fatalf("Failed to create DB config: %v", err)
	}

	if cfg.DialTimeout != 5*time.Second {
		t.Errorf("Expected dial timeout 5s, got %v", cfg.DialTimeout)
	}
	if cfg.ReadTimeout != 10*time.Second {
		t.Errorf("Expected read timeout 10s, got %v", cfg.ReadTimeout)
	}
	if cfg.WriteTimeout != 15*time.Second {
		t.Errorf("Expected write timeout 15s, got %v", cfg.WriteTimeout)
	}
}

// ==================== Integration Tests ====================

// TestServerStartAndStop 测试服务器启动和停止
func TestServerStartAndStop(t *testing.T) {
	// 使用随机端口避免冲突（使用 0 让系统自动分配）
	// 注意：这里我们使用一个非零端口进行测试，因为 NewServer 会验证端口范围
	server, err := NewServer("127.0.0.1",
		WithPort(18080), // 使用一个不太可能冲突的端口
		WithLogging(false),
	)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// 异步启动
	err = server.StartAsync()
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// 给服务器一点启动时间
	time.Sleep(100 * time.Millisecond)

	// 优雅停止
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Stop(ctx)
	if err != nil {
		t.Errorf("Failed to stop server: %v", err)
	}
}

// TestRealWorldScenario 测试真实场景
func TestRealWorldScenario(t *testing.T) {
	// 场景：创建一个生产环境的 API 服务器
	apiHandler := http.NewServeMux()
	apiHandler.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	server, err := NewServer("0.0.0.0",
		WithPort(8080),
		WithHandler(apiHandler),
		WithTimeout(30*time.Second),
		WithMaxConnections(1000),
		WithLogging(true),
		WithMetrics(true),
		WithCORS(true),
		WithSecureHeaders(),
		WithMiddleware(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-API-Version", "v1")
				next.ServeHTTP(w, r)
			})
		}),
	)
	if err != nil {
		t.Fatalf("Failed to create production server: %v", err)
	}

	// 验证所有配置
	if server.port != 8080 {
		t.Errorf("Expected port 8080, got %d", server.port)
	}
	if server.readTimeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", server.readTimeout)
	}
	if server.maxConnections != 1000 {
		t.Errorf("Expected max connections 1000, got %d", server.maxConnections)
	}
	if !server.enableMetrics {
		t.Error("Expected metrics enabled")
	}
	if !server.enableCORS {
		t.Error("Expected CORS enabled")
	}
	if len(server.middlewares) != 2 { // SecureHeaders + custom
		t.Errorf("Expected 2 middlewares, got %d", len(server.middlewares))
	}
}

// TestFunctionalOptionsPatternBenefits 展示函数式选项模式的优势
func TestFunctionalOptionsPatternBenefits(t *testing.T) {
	// 优势 1: 向后兼容 - 添加新选项不会破坏现有代码
	// 优势 2: 自文档化 - 代码清晰展示意图
	// 优势 3: 类型安全 - 编译时检查
	// 优势 4: 默认值友好

	// 示例 1: 最小配置（使用所有默认值）
	_, err := NewServer("localhost")
	if err != nil {
		t.Fatalf("Minimal config failed: %v", err)
	}

	// 示例 2: 只覆盖关心的选项
	_, err = NewServer("localhost", WithPort(3000))
	if err != nil {
		t.Fatalf("Partial config failed: %v", err)
	}

	// 示例 3: 完整配置
	_, err = NewServer("0.0.0.0",
		WithPort(8080),
		WithReadTimeout(30*time.Second),
		WithWriteTimeout(30*time.Second),
		WithIdleTimeout(120*time.Second),
		WithMaxConnections(1000),
		WithMaxHeaderBytes(1<<20),
		WithLogging(true),
		WithMetrics(true),
		WithCORS(true),
	)
	if err != nil {
		t.Fatalf("Full config failed: %v", err)
	}

	// 示例 4: 使用预设 + 自定义
	_, err = NewServer("0.0.0.0",
		append(ProductionDefaults(), WithPort(9000))...,
	)
	if err != nil {
		t.Fatalf("Preset + custom config failed: %v", err)
	}
}

// ExampleNewServer 提供可运行的示例
func ExampleNewServer() {
	// 创建一个简单的服务器
	server, _ := NewServer("localhost",
		WithPort(8080),
		WithTimeout(30*time.Second),
		WithLogging(true),
	)

	fmt.Printf("Server created: %s\n", server.Addr())
	fmt.Printf("Read timeout: %v\n", server.readTimeout)

	// Output:
	// Server created: localhost:8080
	// Read timeout: 30s
}
