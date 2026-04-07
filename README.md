# Go 设计模式学习项目

> 🎯 目标：通过 Go 语言惯用法（Idiomatic Go）深入理解设计模式，避免 Java 式的面向对象思维陷阱。

## 📚 项目介绍

本项目是一个系统性的 Go 语言设计模式学习指南。每个模式都包含：

- **模式实现代码**：符合 Go 语言惯用法的实现，附带详细注释
- **测试驱动示例**：通过 `*_test.go` 文件展示如何使用
- **深度解析文档**：包含适用场景、Go 语言特殊点、优缺点分析

### 为什么这个项目与众不同？

1. **Go 原生思维**：不照搬 Java/C++ 的类继承体系，充分利用 Go 的 `interface`、`embedding`、函数式特性
2. **并发安全优先**：每个模式都会考虑 Go 的并发模型（goroutine、channel、sync 包）
3. **实用主义**：每个模式都附带真实场景的使用建议

---

## 🗺️ 学习路径

建议按照以下顺序学习：

### 第一阶段：创建型模式 (Creational Patterns)
学习如何优雅地创建对象，这是理解其他模式的基础。

| 模式 | 文件路径 | 核心知识点 |
|------|----------|------------|
| 单例模式 | `creational/singleton/` | `sync.Once`、并发安全、懒加载 |
| 工厂方法 | `creational/factory/` | 接口抽象、解耦对象创建 |
| 建造者模式 | `creational/builder/` | 复杂对象构建、函数式选项模式 |
| 抽象工厂模式 | [`creational/abstract_factory/`](./creational/abstract_factory/) | 产品族、一致性创建、Go 中的适用边界 |
| 原型模式 | [`creational/prototype/`](./creational/prototype/) | 浅拷贝、深拷贝、切片与 Map 复制 |

### 第二阶段：结构型模式 (Structural Patterns) ✅
学习如何组合类和对象形成更大的结构。

| 模式 | 文件路径 | 核心知识点 |
|------|----------|------------|
| 适配器模式 | [`structural/adapter/`](./structural/adapter/) | 隐式接口、接口转换、兼容性处理 |
| 装饰器模式 | [`structural/decorator/`](./structural/decorator/) | 高阶函数、HTTP 中间件、运行时扩展 |
| 代理模式 | [`structural/proxy/`](./structural/proxy/) | 延迟加载、访问控制、缓存代理 |
| 组合模式 | [`structural/composite/`](./structural/composite/) | 树形结构、统一接口 |
| 外观模式 | [`structural/facade/`](./structural/facade/) | 简化接口、子系统封装 |
| 桥接模式 | [`structural/bridge/`](./structural/bridge/) | 抽象与实现解耦、接口组合 |
| 享元模式 | [`structural/flyweight/`](./structural/flyweight/) | 共享内部状态、sync.Pool 对象复用 |

### 第三阶段：行为型模式 (Behavioral Patterns) ✅
学习对象间的通信和责任分配。

| 模式 | 文件路径 | 核心知识点 |
|------|----------|------------|
| 策略模式 | [`behavioral/strategy/`](./behavioral/strategy/) | 接口策略、函数类型策略、闭包策略 |
| 观察者模式 | [`behavioral/observer/`](./behavioral/observer/) | 切片实现、Channel 异步讨论、事件总线 |
| 责任链模式 | [`behavioral/chain_of_responsibility/`](./behavioral/chain_of_responsibility/) | 接口链、函数链、中间件风格、构建器模式 |
| 状态模式 | [`behavioral/state/`](./behavioral/state/) | 状态对象、消除 switch-case、状态迁移 |
| 模板方法 | [`behavioral/template_method/`](./behavioral/template_method/) | 组合+接口、闭包骨架、无继承实现 |
| 命令模式 | [`behavioral/command/`](./behavioral/command/) | 请求封装、func 命令、任务调度 |
| 迭代器模式 | [`behavioral/iterator/`](./behavioral/iterator/) | iter.Seq、闭包迭代器、自定义遍历 |
| 备忘录模式 | [`behavioral/memento/`](./behavioral/memento/) | 快照保存、撤销恢复 |
| 中介者模式 | [`behavioral/mediator/`](./behavioral/mediator/) | 中心协调、解耦对象交互 |
| 访问者模式 | [`behavioral/visitor/`](./behavioral/visitor/) | 操作分离、对象结构遍历 |
| 解释器模式 | [`behavioral/interpreter/`](./behavioral/interpreter/) | 简单 DSL、表达式树、规则求值 |

### 第四阶段：Go 惯用模式 (Go Idioms) ✅
Go 语言社区演化出的独特模式，非传统 GoF 设计模式。

| 模式 | 文件路径 | 核心知识点 |
|------|----------|------------|
| 函数式选项 | [`go_idioms/functional_options/`](./go_idioms/functional_options/) | 类型安全配置、向后兼容 API、默认值管理 |
| 扇入扇出 | [`go_idioms/fan_in_fan_out/`](./go_idioms/fan_in_fan_out/) | Goroutine 池、Channel 流水线、并行 Map/Filter/Reduce |
| 工作池 | `go_idioms/worker_pool/` (TODO) | 并发控制、资源复用 |
| 管道模式 | `go_idioms/pipeline/` (TODO) | 数据流处理、stage 组合 |
| Context 模式 | `go_idioms/context/` (TODO) | 取消信号、超时控制 |

---

## ✅ GoF 23 模式完成度总表

| 类型 | 模式 | 目录 | 关键实现文件 |
|------|------|------|--------------|
| 创建型 | 单例模式 | [`creational/singleton/`](./creational/singleton/) | [`singleton.go`](./creational/singleton/singleton.go) |
| 创建型 | 工厂方法 | [`creational/factory/`](./creational/factory/) | [`factory.go`](./creational/factory/factory.go) |
| 创建型 | 建造者模式 | [`creational/builder/`](./creational/builder/) | [`builder.go`](./creational/builder/builder.go) |
| 创建型 | 抽象工厂模式 | [`creational/abstract_factory/`](./creational/abstract_factory/) | [`abstract_factory.go`](./creational/abstract_factory/abstract_factory.go) |
| 创建型 | 原型模式 | [`creational/prototype/`](./creational/prototype/) | [`prototype.go`](./creational/prototype/prototype.go) |
| 结构型 | 适配器模式 | [`structural/adapter/`](./structural/adapter/) | [`adapter.go`](./structural/adapter/adapter.go) |
| 结构型 | 装饰器模式 | [`structural/decorator/`](./structural/decorator/) | [`decorator.go`](./structural/decorator/decorator.go) |
| 结构型 | 代理模式 | [`structural/proxy/`](./structural/proxy/) | [`proxy.go`](./structural/proxy/proxy.go) |
| 结构型 | 组合模式 | [`structural/composite/`](./structural/composite/) | [`composite.go`](./structural/composite/composite.go) |
| 结构型 | 外观模式 | [`structural/facade/`](./structural/facade/) | [`facade.go`](./structural/facade/facade.go) |
| 结构型 | 桥接模式 | [`structural/bridge/`](./structural/bridge/) | [`bridge.go`](./structural/bridge/bridge.go) |
| 结构型 | 享元模式 | [`structural/flyweight/`](./structural/flyweight/) | [`flyweight.go`](./structural/flyweight/flyweight.go) |
| 行为型 | 策略模式 | [`behavioral/strategy/`](./behavioral/strategy/) | [`strategy.go`](./behavioral/strategy/strategy.go) |
| 行为型 | 观察者模式 | [`behavioral/observer/`](./behavioral/observer/) | [`observer.go`](./behavioral/observer/observer.go) |
| 行为型 | 责任链模式 | [`behavioral/chain_of_responsibility/`](./behavioral/chain_of_responsibility/) | [`chain_of_responsibility.go`](./behavioral/chain_of_responsibility/chain_of_responsibility.go) |
| 行为型 | 状态模式 | [`behavioral/state/`](./behavioral/state/) | [`state.go`](./behavioral/state/state.go) |
| 行为型 | 模板方法 | [`behavioral/template_method/`](./behavioral/template_method/) | [`template_method.go`](./behavioral/template_method/template_method.go) |
| 行为型 | 命令模式 | [`behavioral/command/`](./behavioral/command/) | [`command.go`](./behavioral/command/command.go) |
| 行为型 | 迭代器模式 | [`behavioral/iterator/`](./behavioral/iterator/) | [`iterator.go`](./behavioral/iterator/iterator.go) |
| 行为型 | 备忘录模式 | [`behavioral/memento/`](./behavioral/memento/) | [`memento.go`](./behavioral/memento/memento.go) |
| 行为型 | 中介者模式 | [`behavioral/mediator/`](./behavioral/mediator/) | [`mediator.go`](./behavioral/mediator/mediator.go) |
| 行为型 | 访问者模式 | [`behavioral/visitor/`](./behavioral/visitor/) | [`visitor.go`](./behavioral/visitor/visitor.go) |
| 行为型 | 解释器模式 | [`behavioral/interpreter/`](./behavioral/interpreter/) | [`interpreter.go`](./behavioral/interpreter/interpreter.go) |

## 📁 项目目录结构

```
go-design-patterns/
├── README.md                 # 本文件：项目总览
├── go.mod                    # Go 模块定义
├── creational/               # 创建型模式
│   ├── singleton/            # 单例模式
│   ├── factory/              # 工厂方法模式
│   ├── builder/              # 建造者模式
│   ├── abstract_factory/     # 抽象工厂模式 ✅
│   └── prototype/            # 原型模式 ✅
├── structural/               # 结构型模式
│   ├── adapter/              # 适配器模式 ✅
│   ├── decorator/            # 装饰器模式 ✅
│   ├── proxy/                # 代理模式 ✅
│   ├── composite/            # 组合模式 ✅
│   ├── facade/               # 外观模式 ✅
│   ├── bridge/               # 桥接模式 ✅
│   └── flyweight/            # 享元模式 ✅
├── behavioral/               # 行为型模式 ✅
│   ├── strategy/             # 策略模式 ✅
│   ├── observer/             # 观察者模式 ✅
│   ├── chain_of_responsibility/ # 责任链模式 ✅
│   ├── state/                # 状态模式 ✅
│   ├── template_method/      # 模板方法模式 ✅
│   ├── command/              # 命令模式 ✅
│   ├── iterator/             # 迭代器模式 ✅
│   ├── mediator/             # 中介者模式 ✅
│   ├── memento/              # 备忘录模式 ✅
│   ├── visitor/              # 访问者模式 ✅
│   └── interpreter/          # 解释器模式 ✅
└── go_idioms/                # Go 惯用模式 ✅
    ├── functional_options/   # 函数式选项模式 ✅
    ├── fan_in_fan_out/       # 扇入扇出模式 ✅
    ├── worker_pool/          # (TODO)
    ├── pipeline/             # (TODO)
    └── context/              # (TODO)
```

---

## 🚀 快速开始

### 环境要求
- Go 1.21 或更高版本

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定模式的测试
go test ./creational/singleton/...

# 运行并查看覆盖率
go test -cover ./creational/singleton/...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```

### 学习单个模式

每个模式目录包含三个文件：

1. **`xxx.go`** - 模式的核心实现
2. **`xxx_test.go`** - 使用示例和单元测试
3. **`README.md`** - 模式的详细说明

建议学习顺序：
1. 先阅读 `README.md` 理解模式概念
2. 查看 `xxx.go` 的实现细节
3. 运行 `xxx_test.go` 观察实际效果

---

## 🎓 Go 设计模式核心原则

### 1. 组合优于继承 (Composition over Inheritance)

Go 没有继承，使用 **embedding**（嵌入）实现代码复用：

```go
// 不推荐：试图模拟继承
type MyServer struct {
    http.Server  // 嵌入，不是继承
}
```

### 2. 隐式接口 (Implicit Interfaces)

Go 的接口是隐式实现的，这带来了极大的灵活性：

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

// 任何实现了 Read 方法的类型都自动满足 Reader 接口
// 不需要显式声明 `implements`
```

### 3. 零值有用 (Zero Value is Useful)

Go 的类型有默认零值，合理利用可以简化代码：

```go
var mu sync.Mutex  // 零值就是可用的锁
mu.Lock()
```

### 4. 并发原语优先 (Concurrency Primitives)

优先使用 `sync` 包和 channel，而不是手动管理锁：

```go
// 推荐：使用 sync.Once 实现单例
var once sync.Once
once.Do(func() {
    instance = &Singleton{}
})
```

### 5. 函数是一等公民 (First-Class Functions)

充分利用函数类型、闭包和高阶函数：

```go
type Option func(*Config)

func WithTimeout(d time.Duration) Option {
    return func(c *Config) {
        c.timeout = d
    }
}
```

---

## 📖 推荐资源

### 官方资源
- [Effective Go](https://go.dev/doc/effective_go) - Go 语言官方编写指南
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) - 代码审查最佳实践

### 书籍
- 《Go 语言设计模式》
- 《Go 语言高级编程》

### 社区资源
- [Go Patterns](https://github.com/tmrts/go-patterns) - 社区维护的模式集合

---

## 🤝 贡献

欢迎通过 Issue 或 PR 提出改进建议！

---

## 📄 License

MIT License

---

> 💡 **提示**：每个模式的 README.md 都包含详细的适用场景分析和代码示例，建议深入阅读。
