// Package singleton 展示了 Go 语言中实现单例模式的惯用方法。
//
// 单例模式确保一个类只有一个实例，并提供一个全局访问点。
// 在 Go 中，我们利用 sync.Once 来保证线程安全的延迟初始化。
//
// 关键概念：
//   - sync.Once 保证函数只执行一次，即使在并发环境下
//   - 零值可用的特性使得 once 不需要显式初始化
//   - 懒加载（Lazy Initialization）避免程序启动时的不必要开销
package singleton

import (
	"sync"
	"sync/atomic"
)

// Database 是一个模拟数据库连接的单例类型。
// 在实际项目中，这可能是数据库连接池、配置管理器、缓存客户端等。
//
// 设计决策：
//   1. 使用小写字母开头（database）表示包私有，强制通过 GetInstance 访问
//   2. 包含 counter 字段用于演示实例确实只被创建一次
type Database struct {
	// dsn 存储数据库连接字符串
	dsn string

	// counter 用于测试验证，记录该实例被"使用"的次数
	// 使用 atomic 保证并发安全
	counter int32
}

// GetDSN 返回数据库连接字符串
func (d *Database) GetDSN() string {
	return d.dsn
}

// IncrementCounter 原子性地增加计数器
// 用于验证在并发环境下是否只有一个实例
func (d *Database) IncrementCounter() int32 {
	return atomic.AddInt32(&d.counter, 1)
}

// GetCounter 原子性地获取当前计数器值
func (d *Database) GetCounter() int32 {
	return atomic.LoadInt32(&d.counter)
}

// ============================================================================
// 核心实现：使用 sync.Once 的单例模式
// ============================================================================

// instance 是单例实例的存储变量。
// 注意：这里声明为包级别的变量，但小写开头，包外无法直接访问。
var instance *Database

// once 是 sync.Once 类型的变量。
// sync.Once 的零值就是可用的状态，不需要显式初始化。
// 这是 Go "零值有用" 原则的典型体现。
var once sync.Once

// GetInstance 返回 Database 的单例实例。
//
// 这是获取单例的唯一入口，保证：
//   1. 线程安全：sync.Once 内部使用互斥锁和原子操作
//   2. 延迟加载：只在第一次调用时创建实例
//   3. 高性能：创建完成后，后续调用几乎没有开销
//
// 使用示例：
//
//	db := singleton.GetInstance("postgres://localhost/mydb")
//	fmt.Println(db.GetDSN())
//
// 参数说明：
//   - dsn: 数据库连接字符串，只在第一次调用时生效
//
// 注意事项：
//   - 如果需要在创建后修改配置，应该提供专门的方法，而不是重新创建实例
//   - 如果业务需要多个不同配置的实例，说明不应该使用单例模式
func GetInstance(dsn string) *Database {
	// sync.Once.Do 接收一个无参无返回值的函数
	// 该函数只会被执行一次，即使在多个 goroutine 同时调用时
	once.Do(func() {
		// 这段代码在程序生命周期中只执行一次
		instance = &Database{
			dsn:     dsn,
			counter: 0,
		}
	})

	// 返回已创建的实例
	// 注意：如果第一次调用时传入了 dsn，后续调用传入的 dsn 会被忽略
	return instance
}

// ============================================================================
// 扩展实现：支持重置的单例（主要用于测试）
// ============================================================================

// ResetInstance 重置单例实例，主要用于单元测试。
//
// ⚠️ 警告：这个方法不是单例模式的常规用法！
// 它破坏了单例的核心约束（实例唯一性），仅在测试场景中使用。
//
// 为什么需要它：
//   - 单元测试需要隔离性，每个测试用例应该从零开始
//   - 如果不重置，测试执行顺序会影响结果
//
// 生产环境不要使用！
func ResetInstance() {
	// 使用互斥锁保护重置操作
	resetMutex.Lock()
	defer resetMutex.Unlock()

	instance = nil
	// 必须创建新的 sync.Once，因为 sync.Once 不能被重置
	once = sync.Once{}
}

// resetMutex 用于保护 ResetInstance 操作
var resetMutex sync.Mutex

// ============================================================================
// 替代实现：饿汉式单例（Eager Initialization）
// ============================================================================

// eagerInstance 在包初始化时就创建好的单例实例。
// 这种方式称为"饿汉式"，与"懒汉式"（Lazy Initialization）相对。
//
// 适用场景：
//   - 实例创建开销很小
//   - 程序启动时肯定需要该实例
//   - 不需要考虑启动顺序问题
//
// 优点：
//   - 实现简单，没有并发问题
//   - 获取实例时无需检查，性能最优
//
// 缺点：
//   - 增加程序启动时间
//   - 如果初始化失败，程序无法启动
//   - 即使不使用也会占用资源
var eagerInstance = &Database{
	dsn:     "eager-default-dsn",
	counter: 0,
}

// GetEagerInstance 返回饿汉式单例实例
func GetEagerInstance() *Database {
	return eagerInstance
}

// ============================================================================
// 进阶实现：带参数的单例工厂
// ============================================================================

// SingletonFactory 是一个带参数的单例工厂类型。
// 当需要根据不同的"类型"获取不同单例时使用。
//
// 使用场景示例：
//   - 多个数据库连接（主库、从库）
//   - 不同环境的配置管理器
//   - 多租户系统的租户级单例
//
// 注意：这已经不是严格意义上的"单例"，而是"有限实例池"
type SingletonFactory struct {
	// instances 存储不同类型的单例实例
	instances map[string]*Database

	// onceMap 为每个类型存储一个 sync.Once
	onceMap map[string]*sync.Once

	// mu 保护 map 的并发访问
	mu sync.RWMutex
}

// NewSingletonFactory 创建一个新的单例工厂
func NewSingletonFactory() *SingletonFactory {
	return &SingletonFactory{
		instances: make(map[string]*Database),
		onceMap:   make(map[string]*sync.Once),
	}
}

// GetOrCreate 根据类型获取或创建单例实例
//
// 参数：
//   - key: 实例类型的唯一标识
//   - dsn: 创建实例时使用的连接字符串
//
// 返回值：
//   - 对应类型的单例实例
func (f *SingletonFactory) GetOrCreate(key, dsn string) *Database {
	// 先使用读锁检查是否已存在（优化读性能）
	f.mu.RLock()
	if inst, ok := f.instances[key]; ok {
		f.mu.RUnlock()
		return inst
	}
	f.mu.RUnlock()

	// 不存在，需要创建
	f.mu.Lock()
	defer f.mu.Unlock()

	// 双重检查：获取锁后再次确认
	if inst, ok := f.instances[key]; ok {
		return inst
	}

	// 创建该类型的 sync.Once（如果不存在）
	if f.onceMap[key] == nil {
		f.onceMap[key] = &sync.Once{}
	}

	var inst *Database
	f.onceMap[key].Do(func() {
		inst = &Database{
			dsn:     dsn,
			counter: 0,
		}
		f.instances[key] = inst
	})

	return f.instances[key]
}
