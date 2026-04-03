package adapter

import (
	"bytes"
	"io"
	"testing"
)

// ============================================================================
// 支付适配器测试
// ============================================================================

func TestLegacyPaymentAdapter(t *testing.T) {
	// 创建旧系统
	legacy := NewLegacyPaymentSystem("MERCHANT-123")

	// 创建适配器
	adapter := NewLegacyPaymentAdapter(legacy)

	// 测试支付
	t.Run("ProcessPayment", func(t *testing.T) {
		txID, err := adapter.ProcessPayment(99.99, "USD")
		if err != nil {
			t.Fatalf("ProcessPayment failed: %v", err)
		}

		if txID == "" {
			t.Error("Expected non-empty transaction ID")
		}

		// 验证 transaction ID 格式
		if len(txID) < 4 || txID[:3] != "TX-" {
			t.Errorf("Expected transaction ID to start with 'TX-', got: %s", txID)
		}
	})

	// 测试退款
	t.Run("RefundPayment", func(t *testing.T) {
		txID, _ := adapter.ProcessPayment(50.00, "EUR")

		err := adapter.RefundPayment(txID)
		if err != nil {
			t.Fatalf("RefundPayment failed: %v", err)
		}
	})

	// 测试无效 transaction ID
	t.Run("InvalidTransactionID", func(t *testing.T) {
		err := adapter.RefundPayment("INVALID-ID")
		if err == nil {
			t.Error("Expected error for invalid transaction ID")
		}
	})
}

// ============================================================================
// IO 适配器测试
// ============================================================================

func TestStringReaderAdapter(t *testing.T) {
	// 创建 StringReader
	strReader := NewStringReader("Hello, World!")

	// 创建适配器
	adapter := NewStringReaderAdapter(strReader)

	// 测试读取
	t.Run("ReadAll", func(t *testing.T) {
		data, err := io.ReadAll(adapter)
		if err != nil {
			t.Fatalf("ReadAll failed: %v", err)
		}

		expected := "Hello, World!"
		if string(data) != expected {
			t.Errorf("Expected %q, got %q", expected, string(data))
		}
	})

	// 测试分块读取
	t.Run("ReadChunks", func(t *testing.T) {
		strReader := NewStringReader("ABCDEFGHIJ")
		adapter := NewStringReaderAdapter(strReader)

		buf := make([]byte, 3)
		var result []byte

		for {
			n, err := adapter.Read(buf)
			if n > 0 {
				result = append(result, buf[:n]...)
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("Read failed: %v", err)
			}
		}

		expected := "ABCDEFGHIJ"
		if string(result) != expected {
			t.Errorf("Expected %q, got %q", expected, string(result))
		}
	})

	// 测试空字符串
	t.Run("EmptyString", func(t *testing.T) {
		strReader := NewStringReader("")
		adapter := NewStringReaderAdapter(strReader)

		buf := make([]byte, 1024)
		n, err := adapter.Read(buf)

		if err != io.EOF {
			t.Errorf("Expected EOF, got: %v", err)
		}
		if n != 0 {
			t.Errorf("Expected 0 bytes, got: %d", n)
		}
	})
}

// ============================================================================
// 缓存适配器测试
// ============================================================================

func TestExternalCacheAdapter(t *testing.T) {
	external := NewExternalCache()
	adapter := NewExternalCacheAdapter(external)

	t.Run("SetAndGet", func(t *testing.T) {
		// 设置值
		err := adapter.Set("key1", "value1")
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}

		// 获取值
		val, ok := adapter.Get("key1")
		if !ok {
			t.Error("Expected key to exist")
		}
		if val != "value1" {
			t.Errorf("Expected 'value1', got %v", val)
		}
	})

	t.Run("GetNonExistent", func(t *testing.T) {
		val, ok := adapter.Get("nonexistent")
		if ok {
			t.Error("Expected key to not exist")
		}
		if val != nil {
			t.Errorf("Expected nil, got %v", val)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		adapter.Set("key2", "value2")
		adapter.Delete("key2")

		_, ok := adapter.Get("key2")
		if ok {
			t.Error("Expected key to be deleted")
		}
	})
}

// ============================================================================
// 快速适配器测试
// ============================================================================

func TestQuickCacheAdapter(t *testing.T) {
	external := NewExternalCache()
	adapter := NewQuickCacheAdapter(external)

	t.Run("QuickAdapter", func(t *testing.T) {
		// 设置值
		err := adapter.Set("quick-key", "quick-value")
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}

		// 获取值
		val, ok := adapter.Get("quick-key")
		if !ok {
			t.Error("Expected key to exist")
		}
		if val != "quick-value" {
			t.Errorf("Expected 'quick-value', got %v", val)
		}

		// 删除值
		adapter.Delete("quick-key")

		_, ok = adapter.Get("quick-key")
		if ok {
			t.Error("Expected key to be deleted")
		}
	})
}

// ============================================================================
// 云存储适配器测试
// ============================================================================

func TestCloudStorageAdapter(t *testing.T) {
	cloud := NewCloudStorage("my-bucket")
	adapter := NewCloudStorageAdapter(cloud)

	t.Run("SaveAndLoad", func(t *testing.T) {
		// 保存数据
		data := []byte("test data")
		id, err := adapter.Save(data)
		if err != nil {
			t.Fatalf("Save failed: %v", err)
		}

		if id == "" {
			t.Error("Expected non-empty ID")
		}

		// 加载数据
		loaded, err := adapter.Load(id)
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		if !bytes.Equal(loaded, []byte("cloud data")) {
			t.Errorf("Expected 'cloud data', got %s", string(loaded))
		}
	})
}

// ============================================================================
// 双向适配器测试
// ============================================================================

func TestBidirectionalAdapter(t *testing.T) {
	t.Run("EuropeanToUS", func(t *testing.T) {
		europeanDevice := &EuropeanDevice{}
		adapter := NewEuropeanToUSAdapter(europeanDevice)

		result := adapter.ConnectToUSSocket()
		expected := "Adapter converting US 110V to European 220V -> Connected to European socket (220V)"

		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("USToEuropean", func(t *testing.T) {
		usDevice := &USDevice{}
		adapter := NewUSToEuropeanAdapter(usDevice)

		result := adapter.ConnectToEuropeanSocket()
		expected := "Adapter converting European 220V to US 110V -> Connected to US socket (110V)"

		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})
}

// ============================================================================
// 接口兼容性测试
// ============================================================================

func TestInterfaceCompatibility(t *testing.T) {
	// 测试所有适配器是否实现了目标接口

	t.Run("PaymentProcessorInterface", func(t *testing.T) {
		var _ PaymentProcessor = NewLegacyPaymentAdapter(NewLegacyPaymentSystem("test"))
	})

	t.Run("ReaderInterface", func(t *testing.T) {
		var _ io.Reader = NewStringReaderAdapter(NewStringReader("test"))
	})

	t.Run("CacheInterface", func(t *testing.T) {
		var _ Cache = NewExternalCacheAdapter(NewExternalCache())
		var _ Cache = NewQuickCacheAdapter(NewExternalCache())
	})

	t.Run("StorageInterface", func(t *testing.T) {
		var _ Storage = NewCloudStorageAdapter(NewCloudStorage("test"))
	})

	t.Run("PlugInterfaces", func(t *testing.T) {
		var _ USPlug = NewEuropeanToUSAdapter(&EuropeanDevice{})
		var _ EuropeanPlug = NewUSToEuropeanAdapter(&USDevice{})
	})
}

// ============================================================================
// 使用场景测试
// ============================================================================

// ProcessOrder 演示使用 PaymentProcessor 接口的业务逻辑
func ProcessOrder(processor PaymentProcessor, amount float64) (string, error) {
	return processor.ProcessPayment(amount, "USD")
}

func TestRealWorldScenario(t *testing.T) {
	// 场景：系统需要 PaymentProcessor，但我们只有 LegacyPaymentSystem
	legacy := NewLegacyPaymentSystem("MERCHANT-001")
	adapter := NewLegacyPaymentAdapter(legacy)

	// 使用适配器处理订单
	txID, err := ProcessOrder(adapter, 199.99)
	if err != nil {
		t.Fatalf("ProcessOrder failed: %v", err)
	}

	if txID == "" {
		t.Error("Expected non-empty transaction ID")
	}
}

// ReadFromReader 演示使用 io.Reader 的业务逻辑
func ReadFromReader(reader io.Reader) (string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func TestIOAdapterScenario(t *testing.T) {
	// 场景：系统需要 io.Reader，但我们只有 StringReader
	strReader := NewStringReader("adapter pattern in Go")
	adapter := NewStringReaderAdapter(strReader)

	// 使用适配器读取数据
	result, err := ReadFromReader(adapter)
	if err != nil {
		t.Fatalf("ReadFromReader failed: %v", err)
	}

	if result != "adapter pattern in Go" {
		t.Errorf("Expected 'adapter pattern in Go', got %q", result)
	}
}
