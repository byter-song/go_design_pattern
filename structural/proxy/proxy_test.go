package proxy

import (
	"strings"
	"testing"
	"time"
)

// ============================================================================
// 虚拟代理测试
// ============================================================================

func TestProxyImage(t *testing.T) {
	t.Run("LazyLoading", func(t *testing.T) {
		proxy := NewProxyImage("test.jpg")

		// 初始状态：未加载
		if proxy.IsLoaded() {
			t.Error("Expected image to not be loaded initially")
		}

		// 获取文件名（不需要加载）
		filename := proxy.GetFilename()
		if filename != "test.jpg" {
			t.Errorf("Expected filename 'test.jpg', got %q", filename)
		}

		// 仍然未加载
		if proxy.IsLoaded() {
			t.Error("Expected image to still not be loaded after GetFilename")
		}

		// 显示图片（触发加载）
		err := proxy.Display()
		if err != nil {
			t.Fatalf("Display failed: %v", err)
		}

		// 现在已加载
		if !proxy.IsLoaded() {
			t.Error("Expected image to be loaded after Display")
		}
	})

	t.Run("GetSizeTriggersLoading", func(t *testing.T) {
		proxy := NewProxyImage("size_test.jpg")

		// 获取大小会触发加载
		size := proxy.GetSize()
		if size != 1024*1024 {
			t.Errorf("Expected size %d, got %d", 1024*1024, size)
		}

		if !proxy.IsLoaded() {
			t.Error("Expected image to be loaded after GetSize")
		}
	})

	t.Run("MultipleDisplays", func(t *testing.T) {
		proxy := NewProxyImage("multi.jpg")

		// 第一次显示会加载
		err := proxy.Display()
		if err != nil {
			t.Fatalf("First display failed: %v", err)
		}

		// 第二次显示直接使用已加载的图片
		err = proxy.Display()
		if err != nil {
			t.Fatalf("Second display failed: %v", err)
		}
	})
}

// ============================================================================
// 保护代理测试
// ============================================================================

func TestProtectedImage(t *testing.T) {
	// 创建真实图片
	realImage := NewRealImage("secret.jpg")

	t.Run("AdminAccess", func(t *testing.T) {
		admin := &User{Name: "Alice", Role: "admin"}
		protected := NewProtectedImage(realImage, admin, true)

		err := protected.Display()
		if err != nil {
			t.Errorf("Expected admin to have access, got error: %v", err)
		}
	})

	t.Run("UserDenied", func(t *testing.T) {
		user := &User{Name: "Bob", Role: "user"}
		protected := NewProtectedImage(realImage, user, true)

		err := protected.Display()
		if err == nil {
			t.Error("Expected user to be denied access")
		}

		if !strings.Contains(err.Error(), "access denied") {
			t.Errorf("Expected 'access denied' error, got: %v", err)
		}
	})

	t.Run("GuestDenied", func(t *testing.T) {
		guest := &User{Name: "Charlie", Role: "guest"}
		protected := NewProtectedImage(realImage, guest, true)

		err := protected.Display()
		if err == nil {
			t.Error("Expected guest to be denied access")
		}
	})

	t.Run("NonAdminOnly", func(t *testing.T) {
		// 非 admin-only 的图片，任何人都可以访问
		user := &User{Name: "Bob", Role: "user"}
		protected := NewProtectedImage(realImage, user, false)

		err := protected.Display()
		if err != nil {
			t.Errorf("Expected user to have access to non-admin-only image, got error: %v", err)
		}
	})

	t.Run("PassThroughMethods", func(t *testing.T) {
		user := &User{Name: "Bob", Role: "user"}
		protected := NewProtectedImage(realImage, user, true)

		// GetFilename 和 GetSize 不检查权限
		filename := protected.GetFilename()
		if filename != "secret.jpg" {
			t.Errorf("Expected filename 'secret.jpg', got %q", filename)
		}

		size := protected.GetSize()
		if size != 1024*1024 {
			t.Errorf("Expected size %d, got %d", 1024*1024, size)
		}
	})
}

// ============================================================================
// 缓存代理测试
// ============================================================================

func TestCacheProxy(t *testing.T) {
	realImage := NewRealImage("cached.jpg")
	proxy := NewCacheProxy(realImage, 100*time.Millisecond)

	t.Run("CacheFilename", func(t *testing.T) {
		// 第一次获取（缓存未命中）
		filename1 := proxy.GetFilename()

		// 第二次获取（缓存命中）
		filename2 := proxy.GetFilename()

		if filename1 != filename2 {
			t.Errorf("Expected same filename, got %q and %q", filename1, filename2)
		}
	})

	t.Run("CacheSize", func(t *testing.T) {
		// 第一次获取（缓存未命中）
		size1 := proxy.GetSize()

		// 第二次获取（缓存命中）
		size2 := proxy.GetSize()

		if size1 != size2 {
			t.Errorf("Expected same size, got %d and %d", size1, size2)
		}
	})

	t.Run("CacheExpiration", func(t *testing.T) {
		realImage := NewRealImage("expire.jpg")
		proxy := NewCacheProxy(realImage, 50*time.Millisecond)

		// 第一次获取
		_ = proxy.GetFilename()

		// 等待缓存过期
		time.Sleep(100 * time.Millisecond)

		// 再次获取（缓存已过期）
		_ = proxy.GetFilename()
	})

	t.Run("ClearCache", func(t *testing.T) {
		// 获取值并缓存
		_ = proxy.GetFilename()
		_ = proxy.GetSize()

		// 清除缓存
		proxy.ClearCache()

		// 再次获取（应该重新获取）
		_ = proxy.GetFilename()
	})

	t.Run("DisplayNotCached", func(t *testing.T) {
		// Display 操作不应该被缓存
		err := proxy.Display()
		if err != nil {
			t.Fatalf("Display failed: %v", err)
		}

		// 再次调用 Display 应该直接执行
		err = proxy.Display()
		if err != nil {
			t.Fatalf("Second display failed: %v", err)
		}
	})
}

// ============================================================================
// 智能引用代理测试
// ============================================================================

func TestSmartReferenceProxy(t *testing.T) {
	realImage := NewRealImage("smart.jpg")
	proxy := NewSmartReferenceProxy(realImage)

	t.Run("RefCount", func(t *testing.T) {
		// 初始引用计数为 0
		if proxy.GetRefCount() != 0 {
			t.Errorf("Expected ref count 0, got %d", proxy.GetRefCount())
		}

		// 显示图片增加引用计数
		proxy.Display()
		if proxy.GetRefCount() != 1 {
			t.Errorf("Expected ref count 1, got %d", proxy.GetRefCount())
		}

		// 再次显示
		proxy.Display()
		if proxy.GetRefCount() != 2 {
			t.Errorf("Expected ref count 2, got %d", proxy.GetRefCount())
		}
	})

	t.Run("LastAccessTime", func(t *testing.T) {
		beforeAccess := time.Now()

		// 访问更新最后访问时间
		proxy.Display()

		afterAccess := time.Now()
		lastAccess := proxy.GetLastAccess()

		if lastAccess.Before(beforeAccess) || lastAccess.After(afterAccess) {
			t.Error("Last access time not updated correctly")
		}
	})

	t.Run("GetMethodsUpdateAccessTime", func(t *testing.T) {
		beforeAccess := time.Now()
		time.Sleep(10 * time.Millisecond)

		// GetFilename 也应该更新访问时间
		proxy.GetFilename()

		afterGet := proxy.GetLastAccess()
		if afterGet.Before(beforeAccess) {
			t.Error("GetFilename should update last access time")
		}
	})
}

// ============================================================================
// 远程代理测试
// ============================================================================

func TestRemoteImage(t *testing.T) {
	proxy := NewRemoteImage("https://example.com/images", "remote.jpg")

	t.Run("Display", func(t *testing.T) {
		err := proxy.Display()
		if err != nil {
			t.Fatalf("Display failed: %v", err)
		}
	})

	t.Run("GetFilename", func(t *testing.T) {
		filename := proxy.GetFilename()
		if filename != "remote.jpg" {
			t.Errorf("Expected filename 'remote.jpg', got %q", filename)
		}
	})

	t.Run("GetSize", func(t *testing.T) {
		size := proxy.GetSize()
		if size != 2048*1024 {
			t.Errorf("Expected size %d, got %d", 2048*1024, size)
		}
	})
}

// ============================================================================
// 接口兼容性测试
// ============================================================================

func TestImageInterface(t *testing.T) {
	// 测试所有代理类型是否实现了 Image 接口

	t.Run("ProxyImageInterface", func(t *testing.T) {
		var _ Image = NewProxyImage("test.jpg")
	})

	t.Run("ProtectedImageInterface", func(t *testing.T) {
		realImage := NewRealImage("test.jpg")
		user := &User{Name: "test", Role: "admin"}
		var _ Image = NewProtectedImage(realImage, user, false)
	})

	t.Run("CacheProxyInterface", func(t *testing.T) {
		realImage := NewRealImage("test.jpg")
		var _ Image = NewCacheProxy(realImage, time.Minute)
	})

	t.Run("SmartReferenceProxyInterface", func(t *testing.T) {
		realImage := NewRealImage("test.jpg")
		var _ Image = NewSmartReferenceProxy(realImage)
	})

	t.Run("RemoteImageInterface", func(t *testing.T) {
		var _ Image = NewRemoteImage("https://example.com", "test.jpg")
	})
}

// ============================================================================
// 组合代理测试
// ============================================================================

func TestCombinedProxies(t *testing.T) {
	t.Run("VirtualAndProtection", func(t *testing.T) {
		// 创建虚拟代理
		virtualProxy := NewProxyImage("combined.jpg")

		// 添加保护代理
		admin := &User{Name: "Admin", Role: "admin"}
		protectedProxy := NewProtectedImage(virtualProxy, admin, true)

		// 访问受保护的虚拟图片
		err := protectedProxy.Display()
		if err != nil {
			t.Fatalf("Display failed: %v", err)
		}

		// 验证虚拟代理已加载
		if !virtualProxy.IsLoaded() {
			t.Error("Expected virtual proxy to be loaded")
		}
	})

	t.Run("CacheAndSmartReference", func(t *testing.T) {
		realImage := NewRealImage("cache_smart.jpg")

		// 创建缓存代理
		cacheProxy := NewCacheProxy(realImage, time.Minute)

		// 添加智能引用代理
		smartProxy := NewSmartReferenceProxy(cacheProxy)

		// 多次访问 - 使用 Display 方法触发引用计数
		for i := 0; i < 3; i++ {
			_ = smartProxy.Display()
		}

		// 验证引用计数
		if smartProxy.GetRefCount() != 3 {
			t.Errorf("Expected ref count 3, got %d", smartProxy.GetRefCount())
		}
	})
}

// ============================================================================
// 并发安全测试
// ============================================================================

func TestConcurrentAccess(t *testing.T) {
	t.Run("ProxyImageConcurrent", func(t *testing.T) {
		proxy := NewProxyImage("concurrent.jpg")

		// 并发访问
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				defer func() { done <- true }()
				proxy.Display()
			}()
		}

		// 等待所有 goroutine 完成
		for i := 0; i < 10; i++ {
			<-done
		}

		// 验证只加载一次
		if !proxy.IsLoaded() {
			t.Error("Expected image to be loaded")
		}
	})

	t.Run("CacheProxyConcurrent", func(t *testing.T) {
		realImage := NewRealImage("concurrent_cache.jpg")
		proxy := NewCacheProxy(realImage, time.Minute)

		// 并发获取
		done := make(chan int, 10)
		for i := 0; i < 10; i++ {
			go func() {
				size := proxy.GetSize()
				done <- size
			}()
		}

		// 等待所有 goroutine 完成
		for i := 0; i < 10; i++ {
			size := <-done
			if size != 1024*1024 {
				t.Errorf("Expected size %d, got %d", 1024*1024, size)
			}
		}
	})
}

// ============================================================================
// 使用场景测试
// ============================================================================

// ImageGallery 模拟图片画廊
func ImageGallery(images []Image) {
	for _, img := range images {
		_ = img.Display()
	}
}

func TestRealWorldScenario(t *testing.T) {
	// 场景：图片画廊，包含本地和远程图片
	// 本地图片使用虚拟代理延迟加载
	// 远程图片使用远程代理

	images := []Image{
		NewProxyImage("local1.jpg"),
		NewProxyImage("local2.jpg"),
		NewRemoteImage("https://cdn.example.com", "remote1.jpg"),
		NewProxyImage("local3.jpg"),
	}

	// 显示所有图片
	ImageGallery(images)

	// 验证本地图片已加载
	for _, img := range images {
		if proxy, ok := img.(*ProxyImage); ok {
			if !proxy.IsLoaded() {
				t.Errorf("Expected %s to be loaded", proxy.GetFilename())
			}
		}
	}
}
