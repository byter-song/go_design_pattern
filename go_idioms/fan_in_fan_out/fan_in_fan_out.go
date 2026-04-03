// Package fan_in_fan_out 展示了 Go 语言中扇入（Fan-In）和扇出（Fan-Out）并发模式
//
// 这是 Go 并发编程中最强大、最常用的模式之一，充分利用了 Goroutine 和 Channel
// 的特性来处理并行数据流。
//
// 核心概念：
//
// 扇出（Fan-Out）：一个输入源分发给多个 Goroutine 并行处理
//   ┌─────────┐
//   │ Input   │
//   │ Channel │
//   └────┬────┘
//        │
//   ┌────┼────┐
//   │    │    │
//   ▼    ▼    ▼
// ┌───┐┌───┐┌───┐
// │ W ││ W ││ W │  <- 多个 Worker Goroutine 并行处理
// │ 1 ││ 2 ││ 3 │
// └─┬─┘└─┬─┘└─┬─┘
//   │    │    │
//   └────┼────┘
//        ▼
//   ┌─────────┐
//   │ Output  │
//   │ Channel │
//   └─────────┘
//
// 扇入（Fan-In）：多个 Goroutine 的结果合并到一个输出通道
//   ┌───┐┌───┐┌───┐
//   │ W ││ W ││ W │  <- 多个 Worker Goroutine
//   │ 1 ││ 2 ││ 3 │
//   └─┬─┘└─┬─┘└─┬─┘
//     │    │    │
//     └────┼────┘
//          ▼
//    ┌─────────┐
//    │ Output  │
//    │ Channel │
//    └─────────┘
//
// 适用场景：
// - 数据处理流水线（Data Pipeline）
// - 批量任务并行处理
// - 爬虫并发抓取
// - 图片/视频并行处理
// - 大规模数据转换
//
// 设计要点：
// 1. 使用带缓冲的 Channel 减少阻塞
// 2. 使用 sync.WaitGroup 协调 Goroutine 生命周期
// 3. 使用 context.Context 支持取消操作
// 4. 合理设置 Worker 数量（通常等于 CPU 核心数）
// 5. 注意 goroutine 泄漏防护
package fan_in_fan_out

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Processor 定义处理函数类型
type Processor[T, R any] func(T) R

// FanOutConfig 扇出配置
type FanOutConfig struct {
	// Worker 数量，默认为 CPU 核心数
	WorkerCount int
	// 输入通道缓冲大小
	InputBufferSize int
	// 输出通道缓冲大小
	OutputBufferSize int
	// 是否保持输入顺序
	PreserveOrder bool
}

// DefaultFanOutConfig 返回默认配置
func DefaultFanOutConfig() FanOutConfig {
	return FanOutConfig{
		WorkerCount:      runtime.NumCPU(),
		InputBufferSize:  100,
		OutputBufferSize: 100,
		PreserveOrder:    false,
	}
}

// FanOut 执行扇出操作：将输入分发给多个 Worker 并行处理
//
// 类型参数：
//   - T: 输入数据类型
//   - R: 输出数据类型
//
// 参数：
//   - ctx: 上下文，用于取消操作
//   - input: 输入通道
//   - processor: 处理函数
//   - config: 配置（可选，使用默认配置传 nil）
//
// 返回：输出通道
//
// 使用示例：
//
//	input := make(chan int)
//	output := FanOut(context.Background(), input, func(n int) int {
//	    return n * n
//	}, nil)
//
//	for result := range output {
//	    fmt.Println(result)
//	}
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

	// 启动多个 Worker Goroutine
	for i := 0; i < cfg.WorkerCount; i++ {
		go func(workerID int) {
			defer wg.Done()
			worker(ctx, input, output, processor, workerID)
		}(i)
	}

	// 等待所有 Worker 完成后关闭输出通道
	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

// worker 单个 Worker 的实现
func worker[T, R any](
	ctx context.Context,
	input <-chan T,
	output chan<- R,
	processor Processor[T, R],
	workerID int,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case item, ok := <-input:
			if !ok {
				return
			}
			result := processor(item)
			select {
			case <-ctx.Done():
				return
			case output <- result:
			}
		}
	}
}

// FanIn 执行扇入操作：将多个输入通道合并到一个输出通道
//
// 参数：
//   - ctx: 上下文
//   - inputs: 输入通道切片
//
// 返回：合并后的输出通道
func FanIn[T any](ctx context.Context, inputs ...<-chan T) <-chan T {
	output := make(chan T)
	var wg sync.WaitGroup

	// 为每个输入通道启动一个转发 Goroutine
	for _, input := range inputs {
		wg.Add(1)
		go func(ch <-chan T) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case item, ok := <-ch:
					if !ok {
						return
					}
					select {
					case <-ctx.Done():
						return
					case output <- item:
					}
				}
			}
		}(input)
	}

	// 等待所有输入通道关闭后关闭输出通道
	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

// FanOutOrdered 保持顺序的扇出（Ordered Fan-Out）
// 使用序列号机制确保输出顺序与输入顺序一致
//
// 适用场景：需要保持数据处理顺序的场景
// 性能代价：需要额外的排序开销
func FanOutOrdered[T, R any](
	ctx context.Context,
	input <-chan T,
	processor Processor[T, R],
	config *FanOutConfig,
) <-chan R {
	cfg := DefaultFanOutConfig()
	if config != nil {
		cfg = *config
	}
	cfg.PreserveOrder = true

	type orderedItem struct {
		seq    uint64
		result R
	}

	output := make(chan R, cfg.OutputBufferSize)
	orderedOutput := make(chan orderedItem, cfg.OutputBufferSize)

	var wg sync.WaitGroup
	wg.Add(cfg.WorkerCount)

	// 序列号生成器
	var seqCounter uint64
	var seqMu sync.Mutex

	// 启动 Worker
	for i := 0; i < cfg.WorkerCount; i++ {
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case item, ok := <-input:
					if !ok {
						return
					}

					seqMu.Lock()
					seq := seqCounter
					seqCounter++
					seqMu.Unlock()

					result := processor(item)
					select {
					case <-ctx.Done():
						return
					case orderedOutput <- orderedItem{seq: seq, result: result}:
					}
				}
			}
		}(i)
	}

	// 排序协程
	go func() {
		wg.Wait()
		close(orderedOutput)
	}()

	// 排序并输出
	go func() {
		buffer := make(map[uint64]R)
		var nextSeq uint64

		for item := range orderedOutput {
			buffer[item.seq] = item.result

			// 按顺序输出
			for {
				if result, ok := buffer[nextSeq]; ok {
					select {
					case <-ctx.Done():
						return
					case output <- result:
					}
					delete(buffer, nextSeq)
					nextSeq++
				} else {
					break
				}
			}
		}

		// 输出剩余项
		for len(buffer) > 0 {
			if result, ok := buffer[nextSeq]; ok {
				select {
				case <-ctx.Done():
					return
				case output <- result:
				}
				delete(buffer, nextSeq)
				nextSeq++
			}
		}

		close(output)
	}()

	return output
}

// Pipeline 创建完整的处理流水线：Source -> FanOut -> FanIn -> Sink
//
// 参数：
//   - ctx: 上下文
//   - source: 数据源函数，返回输入通道
//   - processor: 处理器函数
//   - sink: 数据接收函数，消费输出通道
//   - config: 配置
//
// 返回：处理的项目数量和可能的错误
func Pipeline[T, R any](
	ctx context.Context,
	source func() <-chan T,
	processor Processor[T, R],
	sink func(<-chan R) error,
	config *FanOutConfig,
) (int, error) {
	input := source()
	output := FanOut(ctx, input, processor, config)

	if err := sink(output); err != nil {
		return 0, err
	}

	return 0, nil
}

// ==================== 实际应用示例 ====================

// ImageTask 图片处理任务
type ImageTask struct {
	ID       string
	Data     []byte
	Format   string
}

// ImageResult 图片处理结果
type ImageResult struct {
	TaskID   string
	Data     []byte
	Format   string
	Duration time.Duration
	Error    error
}

// ProcessImages 并行处理图片
// 展示如何使用 Fan-Out 模式处理 CPU 密集型任务
func ProcessImages(
	ctx context.Context,
	tasks []ImageTask,
	processor func(ImageTask) ImageResult,
	workerCount int,
) []ImageResult {
	// 创建输入通道
	taskChan := make(chan ImageTask, len(tasks))
	go func() {
		for _, task := range tasks {
			select {
			case <-ctx.Done():
				close(taskChan)
				return
			case taskChan <- task:
			}
		}
		close(taskChan)
	}()

	// 配置
	config := &FanOutConfig{
		WorkerCount:      workerCount,
		InputBufferSize:  len(tasks),
		OutputBufferSize: len(tasks),
	}

	// 执行扇出处理
	resultChan := FanOut(ctx, taskChan, processor, config)

	// 收集结果
	var results []ImageResult
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

// URLFetcher URL 获取器配置
type URLFetcher struct {
	// HTTP 客户端超时
	Timeout time.Duration
	// 最大并发数
	MaxConcurrent int
	// 重试次数
	RetryCount int
}

// FetchResult 获取结果
type FetchResult struct {
	URL      string
	Status   int
	Body     []byte
	Duration time.Duration
	Error    error
}

// FetchURLs 并发获取多个 URL
// 展示如何使用 Fan-Out 模式处理 I/O 密集型任务
func FetchURLs(
	ctx context.Context,
	urls []string,
	fetcher URLFetcher,
	fetchFunc func(string) FetchResult,
) []FetchResult {
	// 创建 URL 通道
	urlChan := make(chan string, len(urls))
	go func() {
		for _, url := range urls {
			select {
			case <-ctx.Done():
				close(urlChan)
				return
			case urlChan <- url:
			}
		}
		close(urlChan)
	}()

	// 配置
	config := &FanOutConfig{
		WorkerCount:      fetcher.MaxConcurrent,
		InputBufferSize:  len(urls),
		OutputBufferSize: len(urls),
	}

	// 执行扇出获取
	resultChan := FanOut(ctx, urlChan, fetchFunc, config)

	// 收集结果
	var results []FetchResult
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

// DataTransformer 数据转换器
type DataTransformer[T, R any] struct {
	Transform func(T) (R, error)
	Workers   int
}

// TransformBatch 批量转换数据
// 展示 Fan-In 的使用场景
func TransformBatch[T, R any](
	ctx context.Context,
	inputs []T,
	transformer DataTransformer[T, R],
) ([]R, []error) {
	// 分割数据为多个批次
	batchSize := len(inputs) / transformer.Workers
	if batchSize == 0 {
		batchSize = len(inputs)
	}

	// 为每个批次创建一个通道
	var channels []<-chan R
	var errChannels []<-chan error

	for i := 0; i < len(inputs); i += batchSize {
		end := i + batchSize
		if end > len(inputs) {
			end = len(inputs)
		}
		batch := inputs[i:end]

		resultChan := make(chan R, len(batch))
		errChan := make(chan error, len(batch))

		go func(data []T) {
			defer close(resultChan)
			defer close(errChan)

			for _, item := range data {
				select {
				case <-ctx.Done():
					return
				default:
					result, err := transformer.Transform(item)
					if err != nil {
						select {
						case <-ctx.Done():
							return
						case errChan <- err:
						}
					} else {
						select {
						case <-ctx.Done():
							return
						case resultChan <- result:
						}
					}
				}
			}
		}(batch)

		channels = append(channels, resultChan)
		errChannels = append(errChannels, errChan)
	}

	// 使用 Fan-In 合并结果
	mergedResults := FanIn(ctx, channels...)
	mergedErrors := FanIn(ctx, errChannels...)

	// 收集结果
	var results []R
	var errors []error

	// 使用 WaitGroup 等待两个收集完成
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for r := range mergedResults {
			results = append(results, r)
		}
	}()

	go func() {
		defer wg.Done()
		for e := range mergedErrors {
			errors = append(errors, e)
		}
	}()

	wg.Wait()

	return results, errors
}

// RateLimitedFanOut 带速率限制的扇出
// 控制并发速度，防止系统过载
func RateLimitedFanOut[T, R any](
	ctx context.Context,
	input <-chan T,
	processor Processor[T, R],
	maxConcurrent int,
	ratePerSecond int,
) <-chan R {
	output := make(chan R)

	// 速率限制器
	ticker := time.NewTicker(time.Second / time.Duration(ratePerSecond))
	defer ticker.Stop()

	semaphore := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup

	go func() {
		defer close(output)

		for item := range input {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// 等待速率限制
			}

			semaphore <- struct{}{} // 获取信号量
			wg.Add(1)

			go func(t T) {
				defer wg.Done()
				defer func() { <-semaphore }() // 释放信号量

				result := processor(t)
				select {
				case <-ctx.Done():
					return
				case output <- result:
				}
			}(item)
		}

		wg.Wait()
	}()

	return output
}

// ParallelMap 并行 Map 操作
// 类似于函数式编程中的 map，但是并行执行
func ParallelMap[T, R any](
	ctx context.Context,
	items []T,
	mapper func(T) R,
	parallelism int,
) []R {
	if len(items) == 0 {
		return nil
	}

	// 创建输入通道
	inputChan := make(chan struct {
		index int
		item  T
	}, len(items))

	go func() {
		for i, item := range items {
			select {
			case <-ctx.Done():
				close(inputChan)
				return
			case inputChan <- struct {
				index int
				item  T
			}{index: i, item: item}:
			}
		}
		close(inputChan)
	}()

	// 结果切片
	results := make([]R, len(items))
	var mu sync.Mutex

	// 处理函数
	processor := func(input struct {
		index int
		item  T
	}) R {
		result := mapper(input.item)
		mu.Lock()
		results[input.index] = result
		mu.Unlock()
		return result
	}

	config := &FanOutConfig{
		WorkerCount: parallelism,
	}

	// 执行并行处理
	outputChan := FanOut(ctx, inputChan, processor, config)

	// 等待所有处理完成（通过消费输出通道）
	for range outputChan {
		// 消费输出以等待完成
	}

	return results
}

// ParallelFilter 并行 Filter 操作
// 并行过滤切片中的元素
func ParallelFilter[T any](
	ctx context.Context,
	items []T,
	predicate func(T) bool,
	parallelism int,
) []T {
	if len(items) == 0 {
		return nil
	}

	type itemWithIndex struct {
		index int
		item  T
		keep  bool
	}

	inputChan := make(chan itemWithIndex, len(items))
	go func() {
		for i, item := range items {
			select {
			case <-ctx.Done():
				close(inputChan)
				return
			case inputChan <- itemWithIndex{index: i, item: item}:
			}
		}
		close(inputChan)
	}()

	processor := func(input itemWithIndex) itemWithIndex {
		input.keep = predicate(input.item)
		return input
	}

	config := &FanOutConfig{
		WorkerCount: parallelism,
	}

	resultChan := FanOut(ctx, inputChan, processor, config)

	// 收集结果并保持原始顺序
	marked := make([]itemWithIndex, len(items))
	for result := range resultChan {
		marked[result.index] = result
	}

	// 过滤并返回
	var filtered []T
	for _, m := range marked {
		if m.keep {
			filtered = append(filtered, m.item)
		}
	}

	return filtered
}

// ParallelReduce 并行 Reduce 操作
// 使用分治法实现并行归约
func ParallelReduce[T any](
	ctx context.Context,
	items []T,
	reducer func(T, T) T,
	parallelism int,
) T {
	if len(items) == 0 {
		var zero T
		return zero
	}
	if len(items) == 1 {
		return items[0]
	}

	// 如果数据量小，直接串行处理
	if len(items) <= parallelism*2 {
		result := items[0]
		for i := 1; i < len(items); i++ {
			result = reducer(result, items[i])
		}
		return result
	}

	// 分割数据
	chunkSize := len(items) / parallelism
	if chunkSize == 0 {
		chunkSize = 1
	}

	// 部分结果通道
	partialResults := make(chan T, parallelism)

	var wg sync.WaitGroup
	for i := 0; i < len(items); i += chunkSize {
		end := i + chunkSize
		if end > len(items) {
			end = len(items)
		}
		chunk := items[i:end]

		wg.Add(1)
		go func(data []T) {
			defer wg.Done()

			result := data[0]
			for j := 1; j < len(data); j++ {
				select {
				case <-ctx.Done():
					return
				default:
					result = reducer(result, data[j])
				}
			}

			select {
			case <-ctx.Done():
				return
			case partialResults <- result:
			}
		}(chunk)
	}

	// 等待所有部分计算完成
	go func() {
		wg.Wait()
		close(partialResults)
	}()

	// 合并部分结果
	var partials []T
	for r := range partialResults {
		partials = append(partials, r)
	}

	// 递归合并
	return ParallelReduce(ctx, partials, reducer, parallelism)
}

// ==================== 监控和统计 ====================

// PipelineStats 流水线统计信息
type PipelineStats struct {
	ProcessedCount int64
	ErrorCount     int64
	StartTime      time.Time
	EndTime        time.Time
	Duration       time.Duration
}

// String 返回统计信息的字符串表示
func (s PipelineStats) String() string {
	return fmt.Sprintf(
		"Pipeline Stats:\n"+
			"  Processed: %d\n"+
			"  Errors: %d\n"+
			"  Duration: %v\n"+
			"  Throughput: %.2f items/sec",
		s.ProcessedCount,
		s.ErrorCount,
		s.Duration,
		float64(s.ProcessedCount)/s.Duration.Seconds(),
	)
}
