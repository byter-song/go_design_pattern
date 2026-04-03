package builder

import (
	"testing"
	"time"
)

// TestTraditionalBuilder 测试传统建造者模式
func TestTraditionalBuilder(t *testing.T) {
	t.Run("基本构建", func(t *testing.T) {
		server, err := NewServerBuilder().
			WithHost("localhost").
			WithPort(8080).
			Build()

		if err != nil {
			t.Fatalf("构建失败: %v", err)
		}

		if server.Host != "localhost" {
			t.Errorf("期望 Host 为 localhost，实际为 %s", server.Host)
		}
		if server.Port != 8080 {
			t.Errorf("期望 Port 为 8080，实际为 %d", server.Port)
		}

		// 验证默认值
		if server.ReadTimeout != 30*time.Second {
			t.Errorf("期望 ReadTimeout 为 30s，实际为 %v", server.ReadTimeout)
		}
		if server.MaxConnections != 100 {
			t.Errorf("期望 MaxConnections 为 100，实际为 %d", server.MaxConnections)
		}

		t.Log("✓ 传统建造者基本构建测试通过")
	})

	t.Run("完整配置", func(t *testing.T) {
		server, err := NewServerBuilder().
			WithHost("0.0.0.0").
			WithPort(443).
			WithReadTimeout(60 * time.Second).
			WithWriteTimeout(60 * time.Second).
			WithMaxConnections(1000).
			WithMaxHeaderBytes(2 << 20).
			WithTLS("cert.pem", "key.pem").
			WithDebug().
			WithoutCompression().
			WithLogger("debug", "file").
			Build()

		if err != nil {
			t.Fatalf("构建失败: %v", err)
		}

		if !server.TLS {
			t.Error("期望 TLS 已启用")
		}
		if server.TLSConfig.CertFile != "cert.pem" {
			t.Errorf("期望 CertFile 为 cert.pem，实际为 %s", server.TLSConfig.CertFile)
		}
		if !server.Debug {
			t.Error("期望 Debug 已启用")
		}
		if server.Compress {
			t.Error("期望 Compress 已禁用")
		}
		if server.Logger.Level != "debug" {
			t.Errorf("期望 Logger Level 为 debug，实际为 %s", server.Logger.Level)
		}

		t.Log("✓ 传统建造者完整配置测试通过")
	})

	t.Run("缺少必填字段", func(t *testing.T) {
		_, err := NewServerBuilder().
			WithPort(8080).
			Build()

		if err == nil {
			t.Error("缺少 Host 应该返回错误")
		}

		_, err = NewServerBuilder().
			WithHost("localhost").
			Build()

		if err == nil {
			t.Error("缺少 Port 应该返回错误")
		}

		t.Log("✓ 必填字段验证测试通过")
	})

	t.Run("TLS 配置验证", func(t *testing.T) {
		// 启用 TLS 但没有配置证书
		builder := NewServerBuilder().
			WithHost("localhost").
			WithPort(443)
		builder.server.TLS = true // 直接修改内部状态模拟错误用法

		_, err := builder.Build()
		if err == nil {
			t.Error("启用 TLS 但没有配置证书应该返回错误")
		}

		t.Log("✓ TLS 配置验证测试通过")
	})
}

// TestFunctionalOptions 测试函数式选项模式
func TestFunctionalOptions(t *testing.T) {
	t.Run("基本构建", func(t *testing.T) {
		server, err := NewServer("localhost", 8080)

		if err != nil {
			t.Fatalf("构建失败: %v", err)
		}

		if server.Host != "localhost" {
			t.Errorf("期望 Host 为 localhost，实际为 %s", server.Host)
		}
		if server.Port != 8080 {
			t.Errorf("期望 Port 为 8080，实际为 %d", server.Port)
		}

		t.Log("✓ 函数式选项基本构建测试通过")
	})

	t.Run("使用选项", func(t *testing.T) {
		server, err := NewServer("0.0.0.0", 443,
			WithReadTimeout(60*time.Second),
			WithWriteTimeout(60*time.Second),
			WithMaxConnections(1000),
			WithTLS("cert.pem", "key.pem"),
			WithDebug(),
			WithoutCompression(),
			WithLogger("warn", "file"),
		)

		if err != nil {
			t.Fatalf("构建失败: %v", err)
		}

		if server.ReadTimeout != 60*time.Second {
			t.Errorf("期望 ReadTimeout 为 60s，实际为 %v", server.ReadTimeout)
		}
		if server.MaxConnections != 1000 {
			t.Errorf("期望 MaxConnections 为 1000，实际为 %d", server.MaxConnections)
		}
		if !server.TLS {
			t.Error("期望 TLS 已启用")
		}

		t.Log("✓ 函数式选项使用选项测试通过")
	})

	t.Run("组合选项", func(t *testing.T) {
		// 使用 WithTimeout 同时设置读写超时
		server, err := NewServer("localhost", 8080,
			WithTimeout(45*time.Second),
		)

		if err != nil {
			t.Fatalf("构建失败: %v", err)
		}

		if server.ReadTimeout != 45*time.Second {
			t.Errorf("期望 ReadTimeout 为 45s，实际为 %v", server.ReadTimeout)
		}
		if server.WriteTimeout != 45*time.Second {
			t.Errorf("期望 WriteTimeout 为 45s，实际为 %v", server.WriteTimeout)
		}

		t.Log("✓ 组合选项测试通过")
	})

	t.Run("无效端口", func(t *testing.T) {
		_, err := NewServer("localhost", 0)
		if err == nil {
			t.Error("端口为 0 应该返回错误")
		}

		_, err = NewServer("localhost", 70000)
		if err == nil {
			t.Error("端口超过 65535 应该返回错误")
		}

		_, err = NewServer("localhost", -1)
		if err == nil {
			t.Error("负端口应该返回错误")
		}

		t.Log("✓ 无效端口验证测试通过")
	})
}

// TestOptionPresets 测试选项预设
func TestOptionPresets(t *testing.T) {
	t.Run("生产环境选项", func(t *testing.T) {
		server, err := NewServer("0.0.0.0", 443, ProductionOptions()...)

		if err != nil {
			t.Fatalf("构建失败: %v", err)
		}

		if server.ReadTimeout != 60*time.Second {
			t.Errorf("期望 ReadTimeout 为 60s，实际为 %v", server.ReadTimeout)
		}
		if server.MaxConnections != 1000 {
			t.Errorf("期望 MaxConnections 为 1000，实际为 %d", server.MaxConnections)
		}
		if server.Compress {
			t.Error("生产环境应该禁用压缩")
		}
		if server.Logger.Level != "warn" {
			t.Errorf("期望 Logger Level 为 warn，实际为 %s", server.Logger.Level)
		}

		t.Log("✓ 生产环境选项测试通过")
	})

	t.Run("开发环境选项", func(t *testing.T) {
		server, err := NewServer("localhost", 8080, DevelopmentOptions()...)

		if err != nil {
			t.Fatalf("构建失败: %v", err)
		}

		if server.ReadTimeout != 10*time.Second {
			t.Errorf("期望 ReadTimeout 为 10s，实际为 %v", server.ReadTimeout)
		}
		if !server.Debug {
			t.Error("开发环境应该启用 Debug")
		}
		if server.Logger.Level != "debug" {
			t.Errorf("期望 Logger Level 为 debug，实际为 %s", server.Logger.Level)
		}

		t.Log("✓ 开发环境选项测试通过")
	})

	t.Run("合并选项", func(t *testing.T) {
		// 使用生产环境选项，但覆盖某些设置
		opts := MergeOptions(
			ProductionOptions(),
			[]Option{WithDebug()}, // 在生产环境启用 Debug
		)

		server, err := NewServer("0.0.0.0", 443, opts...)
		if err != nil {
			t.Fatalf("构建失败: %v", err)
		}

		// 验证生产环境设置
		if server.MaxConnections != 1000 {
			t.Errorf("期望 MaxConnections 为 1000，实际为 %d", server.MaxConnections)
		}

		// 验证覆盖的设置
		if !server.Debug {
			t.Error("期望 Debug 被覆盖启用")
		}

		t.Log("✓ 合并选项测试通过")
	})
}

// TestStepBuilder 测试分步建造者
func TestStepBuilder(t *testing.T) {
	t.Run("正确顺序构建", func(t *testing.T) {
		server, err := NewStepBuilder().
			WithHost("localhost").
			WithPort(8080).
			WithTimeout(30 * time.Second).
			Build()

		if err != nil {
			t.Fatalf("构建失败: %v", err)
		}

		if server.Host != "localhost" {
			t.Errorf("期望 Host 为 localhost，实际为 %s", server.Host)
		}
		if server.Port != 8080 {
			t.Errorf("期望 Port 为 8080，实际为 %d", server.Port)
		}

		t.Log("✓ 分步建造者正确顺序构建测试通过")
	})

	t.Run("强制顺序", func(t *testing.T) {
		// 编译时检查：以下代码如果取消注释，应该编译失败
		// 因为必须先调用 WithHost，然后 WithPort，最后才能调用 Build

		// 错误：直接调用 WithPort
		// NewStepBuilder().WithPort(8080)

		// 错误：缺少 WithPort
		// NewStepBuilder().WithHost("localhost").Build()

		t.Log("✓ 分步建造者强制顺序测试通过（编译时检查）")
	})
}

// TestComparison 对比不同实现方式
func TestComparison(t *testing.T) {
	// 传统建造者
	traditional, _ := NewServerBuilder().
		WithHost("localhost").
		WithPort(8080).
		WithTimeout(30 * time.Second).
		WithDebug().
		Build()

	// 函数式选项
	functional, _ := NewServer("localhost", 8080,
		WithTimeout(30*time.Second),
		WithDebug(),
	)

	// 验证两者结果相同
	if traditional.Host != functional.Host {
		t.Error("Host 不匹配")
	}
	if traditional.Port != functional.Port {
		t.Error("Port 不匹配")
	}
	if traditional.ReadTimeout != functional.ReadTimeout {
		t.Error("ReadTimeout 不匹配")
	}
	if traditional.Debug != functional.Debug {
		t.Error("Debug 不匹配")
	}

	t.Log("✓ 不同实现方式结果一致")
}

// BenchmarkTraditionalBuilder 基准测试：传统建造者
func BenchmarkTraditionalBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := NewServerBuilder().
			WithHost("localhost").
			WithPort(8080).
			WithTimeout(30 * time.Second).
			WithDebug().
			Build()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkFunctionalOptions 基准测试：函数式选项
func BenchmarkFunctionalOptions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := NewServer("localhost", 8080,
			WithTimeout(30*time.Second),
			WithDebug(),
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkStepBuilder 基准测试：分步建造者
func BenchmarkStepBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := NewStepBuilder().
			WithHost("localhost").
			WithPort(8080).
			WithTimeout(30 * time.Second).
			Build()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ExampleNewServerBuilder 示例：传统建造者用法
func ExampleNewServerBuilder() {
	server, _ := NewServerBuilder().
		WithHost("localhost").
		WithPort(8080).
		WithTimeout(30 * time.Second).
		WithTLS("cert.pem", "key.pem").
		WithDebug().
		Build()

	_ = server
	// Output:
}

// ExampleNewServer 示例：函数式选项用法
func ExampleNewServer() {
	// 基本用法
	server, _ := NewServer("localhost", 8080)
	_ = server

	// 使用选项
	server, _ = NewServer("0.0.0.0", 443,
		WithTimeout(60*time.Second),
		WithMaxConnections(1000),
		WithTLS("cert.pem", "key.pem"),
		WithDebug(),
	)
	_ = server

	// 使用预设选项
	server, _ = NewServer("localhost", 8080, ProductionOptions()...)
	_ = server

	// Output:
}

// ExampleMergeOptions 示例：合并选项用法
func ExampleMergeOptions() {
	// 组合预设选项和自定义选项
	opts := MergeOptions(
		ProductionOptions(),
		[]Option{WithDebug()},
	)

	server, _ := NewServer("localhost", 8080, opts...)
	_ = server

	// Output:
}
