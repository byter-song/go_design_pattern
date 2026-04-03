package fan_in_fan_out

import (
	"context"
	"fmt"
	"runtime"
	"sort"
	"sync/atomic"
	"testing"
	"time"
)

// TestFanOutBasic 测试基本扇出功能
func TestFanOutBasic(t *testing.T) {
	ctx := context.Background()

	// 创建输入通道
	input := make(chan int, 10)
	for i := 1; i <= 10; i++ {
		input <- i
	}
	close(input)

	// 处理器：平方
	processor := func(n int) int {
		return n * n
	}

	// 执行扇出
	output := FanOut(ctx, input, processor, nil)

	// 收集结果
	var results []int
	for result := range output {
		results = append(results, result)
	}

	// 验证结果数量
	if len(results) != 10 {
		t.Errorf("Expected 10 results, got %d", len(results))
	}

	// 验证结果内容（顺序可能不同）
	sort.Ints(results)
	expected := []int{1, 4, 9, 16, 25, 36, 49, 64, 81, 100}
	for i, exp := range expected {
		if results[i] != exp {
			t.Errorf("Expected result %d at index %d, got %d", exp, i, results[i])
		}
	}
}

// TestFanOutWithConfig 测试带配置的扇出
func TestFanOutWithConfig(t *testing.T) {
	ctx := context.Background()

	input := make(chan int, 100)
	for i := 0; i < 100; i++ {
		input <- i
	}
	close(input)

	var workerCount int32
	processor := func(n int) int {
		atomic.AddInt32(&workerCount, 1)
		time.Sleep(1 * time.Millisecond) // 模拟工作
		return n * 2
	}

	config := &FanOutConfig{
		WorkerCount:      4,
		InputBufferSize:  100,
		OutputBufferSize: 100,
	}

	output := FanOut(ctx, input, processor, config)

	var results int
	for range output {
		results++
	}

	if results != 100 {
		t.Errorf("Expected 100 results, got %d", results)
	}
}

// TestFanOutCancellation 测试扇出取消
func TestFanOutCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	input := make(chan int, 100)
	go func() {
		for i := 0; i < 100; i++ {
			input <- i
			time.Sleep(1 * time.Millisecond)
		}
		close(input)
	}()

	processor := func(n int) int {
		time.Sleep(10 * time.Millisecond) // 模拟耗时操作
		return n
	}

	output := FanOut(ctx, input, processor, nil)

	// 收集几个结果后取消
	var count int
	for range output {
		count++
		if count >= 5 {
			cancel()
			break
		}
	}

	// 给一点时间让 goroutine 退出
	time.Sleep(50 * time.Millisecond)

	// 验证没有 goroutine 泄漏（简化检查）
	if count < 5 {
		t.Errorf("Expected at least 5 results before cancellation, got %d", count)
	}
}

// TestFanIn 测试扇入功能
func TestFanIn(t *testing.T) {
	ctx := context.Background()

	// 创建多个输入通道
	ch1 := make(chan int, 10)
	ch2 := make(chan int, 10)
	ch3 := make(chan int, 10)

	// 发送数据
	go func() {
		for i := 0; i < 10; i++ {
			ch1 <- i
		}
		close(ch1)
	}()

	go func() {
		for i := 10; i < 20; i++ {
			ch2 <- i
		}
		close(ch2)
	}()

	go func() {
		for i := 20; i < 30; i++ {
			ch3 <- i
		}
		close(ch3)
	}()

	// 扇入
	output := FanIn(ctx, ch1, ch2, ch3)

	// 收集结果
	var results []int
	for result := range output {
		results = append(results, result)
	}

	if len(results) != 30 {
		t.Errorf("Expected 30 results, got %d", len(results))
	}

	// 验证所有数据都在
	seen := make(map[int]bool)
	for _, r := range results {
		seen[r] = true
	}
	for i := 0; i < 30; i++ {
		if !seen[i] {
			t.Errorf("Missing result: %d", i)
		}
	}
}

// TestFanOutOrdered 测试有序扇出
func TestFanOutOrdered(t *testing.T) {
	ctx := context.Background()

	input := make(chan int, 10)
	for i := 0; i < 10; i++ {
		input <- i
	}
	close(input)

	// 模拟不同处理时间
	processor := func(n int) int {
		// 数字越大处理时间越长
		time.Sleep(time.Duration(10-n) * time.Millisecond)
		return n * n
	}

	config := &FanOutConfig{
		WorkerCount: 4,
	}

	output := FanOutOrdered(ctx, input, processor, config)

	// 收集结果
	var results []int
	for result := range output {
		results = append(results, result)
	}

	// 验证顺序
	expected := []int{0, 1, 4, 9, 16, 25, 36, 49, 64, 81}
	for i, exp := range expected {
		if results[i] != exp {
			t.Errorf("At index %d: expected %d, got %d", i, exp, results[i])
		}
	}
}

// TestParallelMap 测试并行 Map
func TestParallelMap(t *testing.T) {
	ctx := context.Background()

	items := []int{1, 2, 3, 4, 5}
	mapper := func(n int) int {
		return n * n
	}

	results := ParallelMap(ctx, items, mapper, 2)

	if len(results) != len(items) {
		t.Errorf("Expected %d results, got %d", len(items), len(results))
	}

	// 验证顺序
	expected := []int{1, 4, 9, 16, 25}
	for i, exp := range expected {
		if results[i] != exp {
			t.Errorf("At index %d: expected %d, got %d", i, exp, results[i])
		}
	}
}

// TestParallelMapEmpty 测试空切片
func TestParallelMapEmpty(t *testing.T) {
	ctx := context.Background()

	var items []int
	mapper := func(n int) int {
		return n * n
	}

	results := ParallelMap(ctx, items, mapper, 2)

	if results != nil {
		t.Error("Expected nil for empty input")
	}
}

// TestParallelFilter 测试并行 Filter
func TestParallelFilter(t *testing.T) {
	ctx := context.Background()

	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	predicate := func(n int) bool {
		return n%2 == 0 // 偶数
	}

	results := ParallelFilter(ctx, items, predicate, 2)

	expected := []int{2, 4, 6, 8, 10}
	if len(results) != len(expected) {
		t.Errorf("Expected %d results, got %d", len(expected), len(results))
	}

	for i, exp := range expected {
		if results[i] != exp {
			t.Errorf("At index %d: expected %d, got %d", i, exp, results[i])
		}
	}
}

// TestParallelReduce 测试并行 Reduce
func TestParallelReduce(t *testing.T) {
	ctx := context.Background()

	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	reducer := func(a, b int) int {
		return a + b
	}

	result := ParallelReduce(ctx, items, reducer, 4)

	expected := 55 // 1+2+...+10
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

// TestParallelReduceEmpty 测试空切片 Reduce
func TestParallelReduceEmpty(t *testing.T) {
	ctx := context.Background()

	var items []int
	reducer := func(a, b int) int {
		return a + b
	}

	result := ParallelReduce(ctx, items, reducer, 4)

	if result != 0 {
		t.Errorf("Expected 0 for empty slice, got %d", result)
	}
}

// TestParallelReduceSingle 测试单元素 Reduce
func TestParallelReduceSingle(t *testing.T) {
	ctx := context.Background()

	items := []int{42}
	reducer := func(a, b int) int {
		return a + b
	}

	result := ParallelReduce(ctx, items, reducer, 4)

	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

// TestProcessImages 测试图片处理场景
func TestProcessImages(t *testing.T) {
	ctx := context.Background()

	tasks := []ImageTask{
		{ID: "img1", Data: []byte("data1"), Format: "jpg"},
		{ID: "img2", Data: []byte("data2"), Format: "png"},
		{ID: "img3", Data: []byte("data3"), Format: "gif"},
	}

	processor := func(task ImageTask) ImageResult {
		// 模拟处理
		time.Sleep(5 * time.Millisecond)
		return ImageResult{
			TaskID: task.ID,
			Data:   append([]byte("processed_"), task.Data...),
			Format: task.Format,
		}
	}

	results := ProcessImages(ctx, tasks, processor, 2)

	if len(results) != len(tasks) {
		t.Errorf("Expected %d results, got %d", len(tasks), len(results))
	}

	// 验证所有任务都被处理
	processedIDs := make(map[string]bool)
	for _, r := range results {
		processedIDs[r.TaskID] = true
	}
	for _, task := range tasks {
		if !processedIDs[task.ID] {
			t.Errorf("Task %s was not processed", task.ID)
		}
	}
}

// TestFetchURLs 测试 URL 获取场景
func TestFetchURLs(t *testing.T) {
	ctx := context.Background()

	urls := []string{
		"http://example.com/1",
		"http://example.com/2",
		"http://example.com/3",
	}

	fetcher := URLFetcher{
		Timeout:       5 * time.Second,
		MaxConcurrent: 2,
	}

	fetchFunc := func(url string) FetchResult {
		// 模拟获取
		time.Sleep(5 * time.Millisecond)
		return FetchResult{
			URL:    url,
			Status: 200,
			Body:   []byte("content of " + url),
		}
	}

	results := FetchURLs(ctx, urls, fetcher, fetchFunc)

	if len(results) != len(urls) {
		t.Errorf("Expected %d results, got %d", len(urls), len(results))
	}

	// 验证所有 URL 都被获取
	fetchedURLs := make(map[string]bool)
	for _, r := range results {
		fetchedURLs[r.URL] = true
	}
	for _, url := range urls {
		if !fetchedURLs[url] {
			t.Errorf("URL %s was not fetched", url)
		}
	}
}

// TestTransformBatch 测试批量转换
func TestTransformBatch(t *testing.T) {
	ctx := context.Background()

	inputs := []string{"1", "2", "3", "4", "5"}
	transformer := DataTransformer[string, int]{
		Transform: func(s string) (int, error) {
			var n int
			_, err := fmt.Sscanf(s, "%d", &n)
			return n * 10, err
		},
		Workers: 2,
	}

	results, errors := TransformBatch(ctx, inputs, transformer)

	if len(results) != len(inputs) {
		t.Errorf("Expected %d results, got %d", len(inputs), len(results))
	}

	if len(errors) != 0 {
		t.Errorf("Expected 0 errors, got %d", len(errors))
	}

	// 验证结果
	resultMap := make(map[int]bool)
	for _, r := range results {
		resultMap[r] = true
	}
	for i := 1; i <= 5; i++ {
		if !resultMap[i*10] {
			t.Errorf("Expected result %d not found", i*10)
		}
	}
}

// TestDefaultFanOutConfig 测试默认配置
func TestDefaultFanOutConfig(t *testing.T) {
	cfg := DefaultFanOutConfig()

	if cfg.WorkerCount != runtime.NumCPU() {
		t.Errorf("Expected WorkerCount to be %d, got %d", runtime.NumCPU(), cfg.WorkerCount)
	}

	if cfg.InputBufferSize != 100 {
		t.Errorf("Expected InputBufferSize to be 100, got %d", cfg.InputBufferSize)
	}

	if cfg.OutputBufferSize != 100 {
		t.Errorf("Expected OutputBufferSize to be 100, got %d", cfg.OutputBufferSize)
	}

	if cfg.PreserveOrder != false {
		t.Error("Expected PreserveOrder to be false by default")
	}
}

// TestPipelineStats 测试统计信息
func TestPipelineStats(t *testing.T) {
	stats := PipelineStats{
		ProcessedCount: 1000,
		ErrorCount:     10,
		StartTime:      time.Now().Add(-10 * time.Second),
		EndTime:        time.Now(),
		Duration:       10 * time.Second,
	}

	str := stats.String()

	if str == "" {
		t.Error("Expected non-empty stats string")
	}

	// 验证包含关键信息
	if !contains(str, "Processed: 1000") {
		t.Error("Stats should contain processed count")
	}
	if !contains(str, "Errors: 10") {
		t.Error("Stats should contain error count")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// BenchmarkFanOut 基准测试：扇出性能
func BenchmarkFanOut(b *testing.B) {
	ctx := context.Background()

	processor := func(n int) int {
		// 模拟 CPU 密集型操作
		sum := 0
		for i := 0; i < 1000; i++ {
			sum += i
		}
		return sum
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input := make(chan int, 100)
		for j := 0; j < 100; j++ {
			input <- j
		}
		close(input)

		output := FanOut(ctx, input, processor, nil)
		for range output {
		}
	}
}

// BenchmarkFanOutVsSequential 对比基准测试
func BenchmarkFanOutVsSequential(b *testing.B) {
	ctx := context.Background()

	processor := func(n int) int {
		sum := 0
		for i := 0; i < 10000; i++ {
			sum += i
		}
		return sum
	}

	b.Run("Sequential", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for j := 0; j < 100; j++ {
				processor(j)
			}
		}
	})

	b.Run("FanOut", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			input := make(chan int, 100)
			for j := 0; j < 100; j++ {
				input <- j
			}
			close(input)

			output := FanOut(ctx, input, processor, nil)
			for range output {
			}
		}
	})
}

// BenchmarkParallelMap 基准测试：并行 Map
func BenchmarkParallelMap(b *testing.B) {
	ctx := context.Background()

	items := make([]int, 1000)
	for i := range items {
		items[i] = i
	}

	mapper := func(n int) int {
		sum := 0
		for i := 0; i < 1000; i++ {
			sum += i
		}
		return sum
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParallelMap(ctx, items, mapper, runtime.NumCPU())
	}
}

// ExampleFanOut 可运行的示例
func ExampleFanOut() {
	ctx := context.Background()

	// 创建输入通道
	input := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		input <- i
	}
	close(input)

	// 处理器：计算平方
	processor := func(n int) int {
		return n * n
	}

	// 执行扇出
	output := FanOut(ctx, input, processor, nil)

	// 收集并打印结果
	var results []int
	for result := range output {
		results = append(results, result)
	}

	sort.Ints(results)
	fmt.Println(results)

	// Output: [1 4 9 16 25]
}

// ExampleFanIn 可运行的示例
func ExampleFanIn() {
	ctx := context.Background()

	ch1 := make(chan int, 3)
	ch2 := make(chan int, 3)

	ch1 <- 1
	ch1 <- 2
	ch1 <- 3
	close(ch1)

	ch2 <- 4
	ch2 <- 5
	ch2 <- 6
	close(ch2)

	output := FanIn(ctx, ch1, ch2)

	var sum int
	for result := range output {
		sum += result
	}

	fmt.Println(sum)

	// Output: 21
}
