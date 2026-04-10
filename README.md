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

## 🧠 设计模式速查总览

为了便于后续复习，可以先看这一节，再进入各个子目录深入学习。

**Emoji 说明：**

- `⭐` = 核心设计模式，建议优先掌握
- `🐹` = 在 Go 中更常用、更贴近日常工程实践
- `⭐🐹` = 同时属于核心模式且在 Go 中常见

### 创建型模式

| 模式 | 主要用途 | 核心点 |
|------|----------|--------|
| 单例模式 ⭐🐹 | 保证全局只有一个实例，用于共享资源或全局服务 | `sync.Once`、懒加载、并发安全、生命周期控制 |
| 工厂方法 ⭐🐹 | 隐藏对象创建细节，调用方只依赖抽象接口 | 创建与使用解耦、返回接口、便于扩展实现 |
| 建造者模式 ⭐ | 分步骤构建复杂对象或配置对象 | 构建过程拆分、参数组织清晰、常可类比 Go 的函数式选项 |
| 抽象工厂模式 | 创建一组相关对象，保证产品族的一致性 | 面向产品族设计、隔离具体实现、适合多套配套方案 |
| 原型模式 | 通过复制已有对象快速创建新对象 | 浅拷贝/深拷贝、引用字段复制、避免重复初始化成本 |

### 结构型模式

| 模式 | 主要用途 | 核心点 |
|------|----------|--------|
| 适配器模式 ⭐🐹 | 将不兼容接口转换为可协作接口 | 兼容旧接口、包装转换、利用 Go 隐式接口做适配 |
| 装饰器模式 ⭐🐹 | 在不修改原对象的前提下动态增强能力 | 包装而非继承、职责叠加、Go 中常见于 HTTP middleware |
| 代理模式 🐹 | 为目标对象增加访问控制、缓存、延迟加载等能力 | 控制访问入口、延迟初始化、附加横切逻辑 |
| 组合模式 ⭐ | 用统一方式处理树形结构中的单个对象和组合对象 | 递归结构、统一接口、部分-整体一致对待 |
| 外观模式 ⭐🐹 | 为复杂子系统提供更简单的统一入口 | 降低使用复杂度、封装子系统、收敛调用流程 |
| 桥接模式 ⭐ | 将抽象与实现解耦，使两者可独立变化 | 维度拆分、组合代替继承、避免类爆炸 |
| 享元模式 | 共享可复用状态，降低大量小对象的内存开销 | 内外部状态分离、对象复用、Go 中常联想到 `sync.Pool` |

### 行为型模式

| 模式 | 主要用途 | 核心点 |
|------|----------|--------|
| 策略模式 ⭐🐹 | 在多种算法或行为之间灵活切换 | 面向接口/函数注入策略、消除大量 `if-else` |
| 观察者模式 ⭐🐹 | 一处状态变化通知多个依赖方 | 发布-订阅、事件分发、同步/异步通知权衡 |
| 责任链模式 ⭐🐹 | 将请求按顺序交给多个处理者逐步处理 | 请求链传递、解耦发送者与处理者、Go 中常见于中间件 |
| 状态模式 ⭐ | 将状态相关行为拆分到不同状态对象中 | 状态迁移、消除复杂分支、行为随状态变化 |
| 模板方法 | 固定流程骨架，把部分步骤留给具体实现 | 流程复用、变与不变分离、Go 中多用组合替代继承 |
| 命令模式 🐹 | 将请求封装为对象或函数，便于排队、记录、撤销 | 请求封装、解耦调用者与执行者、Go 中可直接用 `func()` |
| 迭代器模式 | 统一遍历聚合对象而不暴露内部结构 | 遍历与容器分离、顺序访问、Go 1.22+ 可结合 `iter.Seq` |
| 备忘录模式 | 保存对象快照，用于撤销、恢复或历史回滚 | 状态快照、封装内部细节、恢复历史状态 |
| 中介者模式 | 通过中心对象协调多个对象交互 | 降低对象间网状依赖、集中通信规则 |
| 访问者模式 ⭐ | 在不修改对象结构的前提下，为其增加新操作 | 数据结构与操作分离、双分派思想、适合稳定结构+多操作场景 |
| 解释器模式 | 为简单语法或规则系统构建解释执行机制 | 表达式树、DSL、递归求值 |

### Go 惯用模式

| 模式 | 主要用途 | 核心点 |
|------|----------|--------|
| 函数式选项 ⭐🐹 | 解决复杂配置构造与 API 向后兼容问题 | 可选参数、默认值管理、配置扩展性强 |
| 扇入扇出 ⭐🐹 | 聚合多个输入或并发分发任务，提高吞吐量 | goroutine 协作、channel 汇聚/分发、并行处理 |
| 工作池 🐹 | 控制并发数量并复用工作 goroutine | 限流、任务调度、资源复用 |
| 管道模式 🐹 | 将处理流程拆成多个 stage 串联处理数据流 | stage 解耦、流式处理、背压与关闭控制 |
| Context 模式 ⭐🐹 | 在调用链中传递取消信号、超时和请求级元数据 | `context.Context`、超时控制、协程协作取消 |

### 快速记忆建议

- 先掌握带 `⭐🐹` 的模式：它们最能体现设计模式思想，也最容易在 Go 项目里碰到
- 遇到接口不兼容先想适配器，遇到流程增强先想装饰器，遇到流程编排先想责任链
- 处理配置构建优先联想到函数式选项，处理并发协作优先联想到扇入扇出、工作池、Context
- Go 中很多模式会弱化“类层次”，更强调接口、组合、函数和并发原语

---

## 🗺️ 学习路径

建议按照以下顺序学习：

### 第一阶段：创建型模式 (Creational Patterns)
学习如何优雅地创建对象，这是理解其他模式的基础。

| 模式 | 文件路径 | 核心知识点 |
|------|----------|------------|
| 单例模式 ⭐🐹 | `creational/singleton/` | `sync.Once`、并发安全、懒加载 |
| 工厂方法 ⭐🐹 | `creational/factory/` | 接口抽象、解耦对象创建 |
| 建造者模式 ⭐ | `creational/builder/` | 复杂对象构建、函数式选项模式 |
| 抽象工厂模式 | [`creational/abstract_factory/`](./creational/abstract_factory/) | 产品族、一致性创建、Go 中的适用边界 |
| 原型模式 | [`creational/prototype/`](./creational/prototype/) | 浅拷贝、深拷贝、切片与 Map 复制 |

### 第二阶段：结构型模式 (Structural Patterns) ✅
学习如何组合类和对象形成更大的结构。

| 模式 | 文件路径 | 核心知识点 |
|------|----------|------------|
| 适配器模式 ⭐🐹 | [`structural/adapter/`](./structural/adapter/) | 隐式接口、接口转换、兼容性处理 |
| 装饰器模式 ⭐🐹 | [`structural/decorator/`](./structural/decorator/) | 高阶函数、HTTP 中间件、运行时扩展 |
| 代理模式 🐹 | [`structural/proxy/`](./structural/proxy/) | 延迟加载、访问控制、缓存代理 |
| 组合模式 ⭐ | [`structural/composite/`](./structural/composite/) | 树形结构、统一接口 |
| 外观模式 ⭐🐹 | [`structural/facade/`](./structural/facade/) | 简化接口、子系统封装 |
| 桥接模式 ⭐ | [`structural/bridge/`](./structural/bridge/) | 抽象与实现解耦、接口组合 |
| 享元模式 | [`structural/flyweight/`](./structural/flyweight/) | 共享内部状态、sync.Pool 对象复用 |

### 第三阶段：行为型模式 (Behavioral Patterns) ✅
学习对象间的通信和责任分配。

| 模式 | 文件路径 | 核心知识点 |
|------|----------|------------|
| 策略模式 ⭐🐹 | [`behavioral/strategy/`](./behavioral/strategy/) | 接口策略、函数类型策略、闭包策略 |
| 观察者模式 ⭐🐹 | [`behavioral/observer/`](./behavioral/observer/) | 切片实现、Channel 异步讨论、事件总线 |
| 责任链模式 ⭐🐹 | [`behavioral/chain_of_responsibility/`](./behavioral/chain_of_responsibility/) | 接口链、函数链、中间件风格、构建器模式 |
| 状态模式 ⭐ | [`behavioral/state/`](./behavioral/state/) | 状态对象、消除 switch-case、状态迁移 |
| 模板方法 | [`behavioral/template_method/`](./behavioral/template_method/) | 组合+接口、闭包骨架、无继承实现 |
| 命令模式 🐹 | [`behavioral/command/`](./behavioral/command/) | 请求封装、func 命令、任务调度 |
| 迭代器模式 | [`behavioral/iterator/`](./behavioral/iterator/) | iter.Seq、闭包迭代器、自定义遍历 |
| 备忘录模式 | [`behavioral/memento/`](./behavioral/memento/) | 快照保存、撤销恢复 |
| 中介者模式 | [`behavioral/mediator/`](./behavioral/mediator/) | 中心协调、解耦对象交互 |
| 访问者模式 ⭐ | [`behavioral/visitor/`](./behavioral/visitor/) | 操作分离、对象结构遍历 |
| 解释器模式 | [`behavioral/interpreter/`](./behavioral/interpreter/) | 简单 DSL、表达式树、规则求值 |

### 第四阶段：Go 惯用模式 (Go Idioms) ✅
Go 语言社区演化出的独特模式，非传统 GoF 设计模式。

| 模式 | 文件路径 | 核心知识点 |
|------|----------|------------|
| 函数式选项 ⭐🐹 | [`go_idioms/functional_options/`](./go_idioms/functional_options/) | 类型安全配置、向后兼容 API、默认值管理 |
| 扇入扇出 ⭐🐹 | [`go_idioms/fan_in_fan_out/`](./go_idioms/fan_in_fan_out/) | Goroutine 池、Channel 流水线、并行 Map/Filter/Reduce |
| 工作池 🐹 | `go_idioms/worker_pool/` (TODO) | 并发控制、资源复用 |
| 管道模式 🐹 | `go_idioms/pipeline/` (TODO) | 数据流处理、stage 组合 |
| Context 模式 ⭐🐹 | `go_idioms/context/` (TODO) | 取消信号、超时控制 |

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
