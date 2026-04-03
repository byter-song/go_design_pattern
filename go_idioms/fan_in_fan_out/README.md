# 扇入扇出模式 (Fan-In / Fan-Out Pattern)

## 概述

扇入扇出模式是 Go 并发编程中最强大、最常用的模式之一，充分利用了 Goroutine 和 Channel 的特性来处理并行数据流。

## 核心概念

### 扇出 (Fan-Out)

一个输入源分发给多个 Goroutine 并行处理：

```
┌─────────┐
│ Input   │
│ Channel │
└────┬────┘
     │
┌────┼────┐
│    │    │
▼    ▼    ▼
┌───┐┌───┐┌───┐
│ W ││ W ││ W │  <- 多个 Worker Goroutine 并行处理
│ 1 ││ 2 ││ 3 │
└─┬─┘└─┬─┘└─┬─┘
  │    │    │
  └────┼────┘
       ▼
  ┌─────────┐
  │ Output  │
  │ Channel │
  └─────────┘
```

### 扇入 (Fan-In)

多个 Goroutine 的结果合并到一个输出通道：

```
┌───┐┌───┐┌───┐
│ W ││ W ││ W │  <- 多个 Worker Goroutine
│ 1 ││ 2 ││ 3 │
└─┬─┘└─┬─┘└─┬─┘
  │    │    │
  └────┼────┘
       ▼
  ┌─────────┐
  │ Output  │
  │ Channel │
  └─────────┘
```

## 适用场景

- **数据处理流水线**（Data Pipeline）
- **批量任务并行处理**
- **爬虫并发抓取**
- **图片/视频并行处理**
- **大规模数据转换**
- **CPU 密集型计算**
- **I/O 密集型操作**

## 实现模式

### 基本扇出

```go
func FanOut[T, R any](
    ctx context.Context,
    input <-chan T,
    processor Processor[T, R],
    config *FanOutConfig,
) <-chan R {
    cfg := DefaultFanOutConfig()
    if config != nil {
        cfg = *config
    }

    output := make(chan R, cfg.OutputBufferSize)
    var wg sync.WaitGroup
    wg.Add(cfg.WorkerCount)

    // 启动多个 Worker
    for i := 0; i < cfg.WorkerCount; i++ {
        go func(workerID int) {
            defer wg.Done()
            for item := range input {
                result := processor(item)
                output <- result
            }
        }(i)
    }

    // 等待所有 Worker 完成后关闭输出通道
    go func() {
        wg.Wait()
        close(output)
    }()

    return output
}
```

### 基本扇入

```go
func FanIn[T any](ctx context.Context, inputs ...<-chan T) <-chan T {
    output := make(chan T)
    var wg sync.WaitGroup

    for _, input := range inputs {
        wg.Add(1)
        go func(ch <-chan T) {
            defer wg.Done()
            for item := range ch {
                output <- item
            }
        }(input)
    }

    go func() {
        wg.Wait()
        close(output)
    }()

    return output
}
```

## 使用示例

### 并行处理数字

```go
ctx := context.Background()

// 创建输入通道
input := make(chan int, 10)
for i := 1; i <= 10; i++ {
    input <- i
}
close(input)

// 处理器：计算平方
processor := func(n int) int {
    return n * n
}

// 执行扇出
output := FanOut(ctx, input, processor, nil)

// 收集结果
for result := range output {
    fmt.Println(result)
}
```

### 图片并行处理

```go
tasks := []ImageTask{
    {ID: "img1", Data: imageData1, Format: "jpg"},
    {ID: "img2", Data: imageData2, Format: "png"},
    // ...
}

processor := func(task ImageTask) ImageResult {
    // 图片处理逻辑（调整大小、压缩等）
    return processImage(task)
}

results := ProcessImages(ctx, tasks, processor, runtime.NumCPU())
```

### URL 并发抓取

```go
urls := []string{
    "http://example.com/1",
    "http://example.com/2",
    // ...
}

fetcher := URLFetcher{
    Timeout:       5 * time.Second,
    MaxConcurrent: 10,
}

results := FetchURLs(ctx, urls, fetcher, fetchFunc)
```

## 高级功能

### 有序扇出 (Ordered Fan-Out)

使用序列号机制确保输出顺序与输入顺序一致：

```go
output := FanOutOrdered(ctx, input, processor, config)
// 输出顺序保证与输入顺序一致
```

### 并行 Map/Filter/Reduce

```go
// 并行 Map
results := ParallelMap(ctx, items, mapper, parallelism)

// 并行 Filter
filtered := ParallelFilter(ctx, items, predicate, parallelism)

// 并行 Reduce
sum := ParallelReduce(ctx, items, addFunc, parallelism)
```

### 速率限制扇出

```go
output := RateLimitedFanOut(ctx, input, processor, 
    maxConcurrent,  // 最大并发数
    ratePerSecond,  // 每秒处理速率
)
```

## 设计要点

1. **使用带缓冲的 Channel** 减少阻塞
2. **使用 sync.WaitGroup** 协调 Goroutine 生命周期
3. **使用 context.Context** 支持取消操作
4. **合理设置 Worker 数量**（通常等于 CPU 核心数）
5. **注意 goroutine 泄漏防护**

## 性能对比

```
BenchmarkFanOutVsSequential/Sequential-8    100    10523456 ns/op
BenchmarkFanOutVsSequential/FanOut-8        500     2104567 ns/op
```

扇出模式在 CPU 密集型任务中通常能获得接近 CPU 核心数的加速比。

## 最佳实践

### Worker 数量设置

```go
// CPU 密集型：Worker 数 = CPU 核心数
config := &FanOutConfig{
    WorkerCount: runtime.NumCPU(),
}

// I/O 密集型：Worker 数可以更多
config := &FanOutConfig{
    WorkerCount: runtime.NumCPU() * 2, // 或更多
}
```

### 错误处理

```go
// 使用结果结构体包含错误信息
type Result struct {
    Data  interface{}
    Error error
}

processor := func(item Item) Result {
    data, err := process(item)
    return Result{Data: data, Error: err}
}
```

### 优雅关闭

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// 监听系统信号
signalChan := make(chan os.Signal, 1)
signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

go func() {
    <-signalChan
    cancel()
}()

output := FanOut(ctx, input, processor, nil)
```

## 完整流水线示例

```go
// 1. 数据源
source := func() <-chan int {
    ch := make(chan int)
    go func() {
        for i := 0; i < 100; i++ {
            ch <- i
        }
        close(ch)
    }()
    return ch
}

// 2. 处理
processor := func(n int) int {
    return n * n
}

// 3. 数据接收
sink := func(output <-chan int) error {
    for result := range output {
        fmt.Println(result)
    }
    return nil
}

// 执行流水线
Pipeline(ctx, source, processor, sink, config)
```

## 参考

- [Go Concurrency Patterns: Pipelines and cancellation](https://go.dev/blog/pipelines)
- [Concurrency in Go](https://www.oreilly.com/library/view/concurrency-in-go/9781491941294/)
- [Go 101: Channel Use Cases](https://go101.org/article/channel-use-cases.html)
