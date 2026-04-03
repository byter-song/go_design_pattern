# 单例模式 (Singleton Pattern)

## 模式定义

单例模式确保一个类只有一个实例，并提供一个全局访问点来访问这个实例。

> **核心思想**：控制实例化过程，保证全局唯一性。

---

## 适用场景

### 应该使用单例的场景

1. **资源管理类**
   - 数据库连接池（创建成本高，需要复用）
   - 线程池、连接池
   - 缓存客户端（Redis、Memcached）

2. **配置管理**
   - 全局配置对象
   - 环境变量管理器
   - 特性开关（Feature Flags）

3. **硬件接口**
   - 打印机队列管理
   - 文件系统访问器

4. **日志记录器**
   - 统一的日志输出控制
   - 避免多个日志实例竞争文件句柄

### 不应该使用单例的场景

1. **需要多实例的场景**
   - 不同配置的数据库连接
   - 多租户系统的租户隔离

2. **测试困难的场景**
   - 单例使得单元测试难以隔离
   - 难以 Mock 依赖

3. **需要横向扩展的场景**
   - 微服务架构中，单例可能成为瓶颈

---

## Go 语言实现

### 核心实现：`sync.Once` 懒加载

```go
var instance *Database
var once sync.Once

func GetInstance(dsn string) *Database {
    once.Do(func() {
        instance = &Database{dsn: dsn}
    })
    return instance
}
```

#### 为什么选择 `sync.Once`？

| 方案 | 优点 | 缺点 |
|------|------|------|
| `sync.Once` | 简洁、标准库支持、性能优秀 | 无法重置 |
| `sync.Mutex` + 检查 | 灵活、可重置 | 代码复杂、易出错 |
| `atomic` 操作 | 性能最好 | 代码复杂、可读性差 |
| `init()` 函数 | 最简单 | 不支持懒加载、无法处理错误 |

**`sync.Once` 的优势**：
- **简洁性**：只需几行代码
- **线程安全**：内部使用原子操作 + 互斥锁，保证并发安全
- **性能**：实例创建后，后续调用几乎没有开销（原子读取标志位）
- **惯用法**：Go 社区标准做法

### 替代实现对比

#### 1. 饿汉式（Eager Initialization）

```go
var eagerInstance = &Database{dsn: "default"}

func GetEagerInstance() *Database {
    return eagerInstance
}
```

**适用场景**：
- 实例创建开销很小
- 程序启动时肯定需要
- 不需要处理初始化错误

**优缺点**：
- ✅ 实现最简单，无并发问题
- ✅ 获取实例性能最优
- ❌ 增加启动时间
- ❌ 无法懒加载
- ❌ 初始化失败会导致程序无法启动

#### 2. 带参数的单例工厂

当需要管理多个"单例"（如主库、从库）时使用：

```go
type SingletonFactory struct {
    instances map[string]*Database
    onceMap   map[string]*sync.Once
    mu        sync.RWMutex
}

func (f *SingletonFactory) GetOrCreate(key, dsn string) *Database {
    // 先读检查
    f.mu.RLock()
    if inst, ok := f.instances[key]; ok {
        f.mu.RUnlock()
        return inst
    }
    f.mu.RUnlock()
    
    // 加锁创建
    f.mu.Lock()
    defer f.mu.Unlock()
    
    // 双重检查
    if inst, ok := f.instances[key]; ok {
        return inst
    }
    
    // 创建实例
    f.onceMap[key].Do(func() {
        f.instances[key] = &Database{dsn: dsn}
    })
    
    return f.instances[key]
}
```

**注意**：这已经不是严格意义上的单例，而是"有限实例池"。

---

## Go 语言特殊点

### 1. 零值可用（Zero Value is Useful）

```go
var once sync.Once  // 零值就是可用的！
```

Go 的 `sync.Once` 零值已经处于可用状态，不需要显式初始化。这是 Go 与其他语言的重要区别。

### 2. 包级变量 + 小写 = 封装

```go
var instance *Database  // 小写开头，包外不可见

func GetInstance() *Database {  // 大写开头，公开 API
    // ...
}
```

Go 通过命名约定实现访问控制，不需要 `private`/`public` 关键字。

### 3. 隐式接口

单例返回的具体类型可以实现接口，便于测试时 Mock：

```go
type DatabaseInterface interface {
    Query(sql string) (*Rows, error)
}

// 单例返回具体类型
func GetInstance() *Database { ... }

// 但使用者可以依赖接口
func ProcessData(db DatabaseInterface) { ... }
```

### 4. 测试的挑战与解决

单例在测试中最大的问题是**状态污染**。解决方案：

```go
// 提供重置方法（仅用于测试）
func ResetInstance() {
    resetMutex.Lock()
    defer resetMutex.Unlock()
    
    instance = nil
    once = sync.Once{}  // 必须创建新的 sync.Once
}
```

**更好的做法**：避免使用单例，改用**依赖注入**。

---

## 优缺点分析

### 优点

1. **资源控制**
   - 确保昂贵资源只创建一次
   - 统一管理和复用

2. **全局访问**
   - 任何地方都能获取实例
   - 简化 API 设计

3. **延迟初始化**
   - 按需创建，节省启动时间
   - 避免不必要的资源占用

### 缺点

1. **违反单一职责原则**
   - 类既要负责业务逻辑，又要管理实例化

2. **隐藏依赖**
   - 代码中直接调用 `GetInstance()`，依赖关系不明显
   - 难以追踪数据流向

3. **测试困难**
   - 全局状态导致测试间相互影响
   - 难以 Mock 和隔离

4. **并发复杂性**
   - 需要处理线程安全问题
   - 不当实现可能导致性能问题

5. **扩展困难**
   - 一旦需要多个实例，改造成本高
   - 与微服务架构理念冲突

---

## 最佳实践

### ✅ 推荐做法

1. **优先使用 `sync.Once`**
   ```go
   var once sync.Once
   once.Do(func() { ... })
   ```

2. **考虑是否需要单例**
   - 是否真的需要全局唯一？
   - 是否可以用依赖注入替代？

3. **提供接口抽象**
   ```go
   type Logger interface {
       Log(msg string)
   }
   
   func GetLogger() Logger { ... }
   ```

4. **延迟错误处理**
   ```go
   type Database struct {
       err error  // 保存初始化错误
   }
   
   func (d *Database) Query(...) (*Rows, error) {
       if d.err != nil {
           return nil, d.err
       }
       // ...
   }
   ```

### ❌ 避免的做法

1. **不要滥用单例**
   ```go
   // 错误：简单的工具类不需要单例
   type StringUtils struct{}
   func GetStringUtils() *StringUtils { ... }
   ```

2. **不要在单例中持有可变状态**
   ```go
   // 危险：并发访问可能出问题
   func (d *Database) SetConfig(cfg Config) {
       d.config = cfg  // 不是线程安全的！
   }
   ```

3. **不要忽视测试性**
   ```go
   // 难以测试：直接调用单例
   func ProcessOrder(orderID string) {
       db := GetInstance("")
       db.Query(...)
   }
   
   // 推荐：依赖注入
   func ProcessOrder(db Database, orderID string) {
       db.Query(...)
   }
   ```

---

## 相关模式

| 模式 | 关系 | 说明 |
|------|------|------|
| **工厂模式** | 配合使用 | 单例可以使用工厂方法创建 |
| **依赖注入** | 替代方案 | 用 DI 容器管理实例生命周期 |
| **Monostate** | 变体 | 通过共享状态实现"逻辑单例" |

---

## 总结

单例模式在 Go 中的最佳实现是使用 `sync.Once`，它简洁、线程安全且性能优秀。但更重要的是**审慎地决定是否需要单例**——在现代 Go 开发中，依赖注入往往是更好的选择。

> 💡 **核心建议**：单例应该是最后的手段，而不是第一选择。在决定使用单例之前，先考虑是否可以通过函数参数传递依赖。
