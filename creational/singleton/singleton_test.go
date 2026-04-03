package singleton

import (
	"sync"
	"testing"
)

// TestBasicSingleton 测试最基本的单例功能
// 验证：多次调用 GetInstance 返回的是同一个实例
func TestBasicSingleton(t *testing.T) {
	// 每个测试开始前重置，保证测试隔离性
	ResetInstance()

	// 第一次调用，创建实例
	db1 := GetInstance("postgres://localhost/db1")

	// 第二次调用，应该返回同一个实例
	db2 := GetInstance("postgres://localhost/db2")

	// 验证是同一个实例（指针相同）
	if db1 != db2 {
		t.Error("GetInstance 应该返回同一个实例")
	}

	// 验证使用的是第一次传入的参数
	if db1.GetDSN() != "postgres://localhost/db1" {
		t.Errorf("期望 DSN 为 postgres://localhost/db1，实际为 %s", db1.GetDSN())
	}

	t.Logf("✓ 基本单例测试通过：两次调用返回同一实例，DSN = %s", db1.GetDSN())
}

// TestConcurrentSingleton 测试并发环境下的单例
// 验证：sync.Once 能保证在并发情况下实例只被创建一次
func TestConcurrentSingleton(t *testing.T) {
	ResetInstance()

	const numGoroutines = 100
	var wg sync.WaitGroup
	instances := make(chan *Database, numGoroutines)

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			// 每个 goroutine 都尝试获取实例
			db := GetInstance("concurrent-test-dsn")
			instances <- db
		}(i)
	}

	wg.Wait()
	close(instances)

	// 收集所有返回的实例
	var firstInstance *Database
	count := 0
	for inst := range instances {
		if firstInstance == nil {
			firstInstance = inst
		}
		if inst != firstInstance {
			t.Error("并发环境下返回了不同的实例！")
		}
		count++
	}

	if count != numGoroutines {
		t.Errorf("期望收到 %d 个实例，实际收到 %d 个", numGoroutines, count)
	}

	t.Logf("✓ 并发测试通过：%d 个 goroutine 都获取到同一实例", numGoroutines)
}

// TestCounterConcurrency 测试实例内部状态在并发下的正确性
// 使用原子操作保证计数器的准确性
func TestCounterConcurrency(t *testing.T) {
	ResetInstance()

	db := GetInstance("counter-test-dsn")

	const numOperations = 1000
	var wg sync.WaitGroup

	wg.Add(numOperations)
	for i := 0; i < numOperations; i++ {
		go func() {
			defer wg.Done()
			db.IncrementCounter()
		}()
	}

	wg.Wait()

	finalCount := db.GetCounter()
	if finalCount != numOperations {
		t.Errorf("期望计数器值为 %d，实际为 %d", numOperations, finalCount)
	}

	t.Logf("✓ 原子计数器测试通过：%d 次并发递增后值为 %d", numOperations, finalCount)
}

// TestEagerSingleton 测试饿汉式单例
func TestEagerSingleton(t *testing.T) {
	// 饿汉式不需要重置，因为实例在包加载时就已创建
	db1 := GetEagerInstance()
	db2 := GetEagerInstance()

	if db1 != db2 {
		t.Error("饿汉式单例应该返回同一个实例")
	}

	if db1.GetDSN() != "eager-default-dsn" {
		t.Errorf("期望 DSN 为 eager-default-dsn，实际为 %s", db1.GetDSN())
	}

	t.Logf("✓ 饿汉式单例测试通过")
}

// TestSingletonFactory 测试带参数的单例工厂
func TestSingletonFactory(t *testing.T) {
	factory := NewSingletonFactory()

	// 获取主库实例
	master := factory.GetOrCreate("master", "postgres://master/db")

	// 获取从库实例
	slave := factory.GetOrCreate("slave", "postgres://slave/db")

	// 验证两个实例不同
	if master == slave {
		t.Error("主库和从库应该是不同的实例")
	}

	// 再次获取主库，应该是同一个实例
	master2 := factory.GetOrCreate("master", "postgres://master-new/db")
	if master != master2 {
		t.Error("相同 key 应该返回同一实例")
	}

	// 验证 DSN 没有被覆盖（使用第一次的值）
	if master2.GetDSN() != "postgres://master/db" {
		t.Errorf("DSN 不应该被覆盖，期望 postgres://master/db，实际为 %s", master2.GetDSN())
	}

	t.Logf("✓ 单例工厂测试通过：master=%s, slave=%s", master.GetDSN(), slave.GetDSN())
}

// TestSingletonFactoryConcurrent 测试工厂的并发安全性
func TestSingletonFactoryConcurrent(t *testing.T) {
	factory := NewSingletonFactory()

	const numGoroutines = 50
	var wg sync.WaitGroup

	wg.Add(numGoroutines * 2)

	// 并发获取 master
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			factory.GetOrCreate("master", "postgres://master/db")
		}()
	}

	// 并发获取 slave
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			factory.GetOrCreate("slave", "postgres://slave/db")
		}()
	}

	wg.Wait()

	// 验证最终只有两个实例
	master := factory.GetOrCreate("master", "ignored")
	slave := factory.GetOrCreate("slave", "ignored")

	if master == slave {
		t.Error("master 和 slave 应该是不同实例")
	}

	t.Logf("✓ 工厂并发测试通过：%d 个 goroutine 安全地创建了 2 个实例", numGoroutines*2)
}

// BenchmarkGetInstance 基准测试：单例获取性能
// 验证：创建完成后，获取实例的开销极小
func BenchmarkGetInstance(b *testing.B) {
	ResetInstance()
	// 先创建实例
	GetInstance("benchmark-dsn")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = GetInstance("ignored")
		}
	})
}

// BenchmarkEagerGetInstance 基准测试：饿汉式性能
func BenchmarkEagerGetInstance(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = GetEagerInstance()
		}
	})
}

// BenchmarkFactoryGetOrCreate 基准测试：工厂模式性能
func BenchmarkFactoryGetOrCreate(b *testing.B) {
	factory := NewSingletonFactory()
	// 预热
	factory.GetOrCreate("key", "value")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = factory.GetOrCreate("key", "ignored")
		}
	})
}

// ExampleGetInstance 示例函数：展示基本用法
func ExampleGetInstance() {
	// 实际使用时不需要 ResetInstance，这里是为了示例可重复运行
	ResetInstance()

	// 获取数据库连接单例
	db := GetInstance("postgres://user:pass@localhost/mydb")

	// 使用实例
	_ = db.GetDSN()

	// Output:
}
