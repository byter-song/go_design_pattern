package decorator

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestLoggingMiddleware 测试日志装饰器
func TestLoggingMiddleware(t *testing.T) {
	handler := LoggingMiddleware(BusinessHandler)

	ctx := context.WithValue(context.Background(), "userID", "user123")
	req := &Request{Data: map[string]interface{}{"orderID": "order456"}}

	resp, err := handler(ctx, req)
	if err != nil {
		t.Fatalf("处理失败: %v", err)
	}

	if resp.Data["message"] != "order processed" {
		t.Errorf("响应消息不正确: %v", resp.Data["message"])
	}

	t.Log("✓ 日志装饰器测试通过")
}

// TestAuthMiddleware 测试认证装饰器
func TestAuthMiddleware(t *testing.T) {
	// 模拟 token 验证器
	validator := func(token string) (string, bool) {
		if token == "valid-token" {
			return "user123", true
		}
		return "", false
	}

	handler := AuthMiddleware(validator)(BusinessHandler)

	t.Run("有效 token", func(t *testing.T) {
		req := &Request{Data: map[string]interface{}{"token": "valid-token"}}
		resp, err := handler(context.Background(), req)
		if err != nil {
			t.Fatalf("处理失败: %v", err)
		}
		if resp.Data["userID"] != "user123" {
			t.Errorf("userID 不正确: %v", resp.Data["userID"])
		}
		t.Log("✓ 有效 token 测试通过")
	})

	t.Run("无效 token", func(t *testing.T) {
		req := &Request{Data: map[string]interface{}{"token": "invalid-token"}}
		_, err := handler(context.Background(), req)
		if err == nil {
			t.Error("无效 token 应该返回错误")
		}
		t.Log("✓ 无效 token 测试通过")
	})

	t.Run("缺少 token", func(t *testing.T) {
		req := &Request{Data: map[string]interface{}{}}
		_, err := handler(context.Background(), req)
		if err == nil {
			t.Error("缺少 token 应该返回错误")
		}
		t.Log("✓ 缺少 token 测试通过")
	})
}

// TestRecoveryMiddleware 测试恢复装饰器
func TestRecoveryMiddleware(t *testing.T) {
	// 创建一个会 panic 的 handler
	panicHandler := func(ctx context.Context, req *Request) (*Response, error) {
		panic("something went wrong")
	}

	handler := RecoveryMiddleware(panicHandler)

	req := &Request{Data: map[string]interface{}{}}
	resp, err := handler(context.Background(), req)

	if resp != nil {
		t.Error("panic 后应该返回 nil 响应")
	}
	if err == nil {
		t.Error("panic 应该被转换为错误")
	}
	if err.Error() != "internal server error: something went wrong" {
		t.Errorf("错误消息不正确: %v", err.Error())
	}

	t.Log("✓ 恢复装饰器测试通过")
}

// TestTimeoutMiddleware 测试超时装饰器
func TestTimeoutMiddleware(t *testing.T) {
	// 创建一个慢 handler
	slowHandler := func(ctx context.Context, req *Request) (*Response, error) {
		time.Sleep(200 * time.Millisecond)
		return &Response{Data: map[string]interface{}{"status": "ok"}}, nil
	}

	t.Run("超时", func(t *testing.T) {
		handler := TimeoutMiddleware(50 * time.Millisecond)(slowHandler)
		req := &Request{Data: map[string]interface{}{}}
		_, err := handler(context.Background(), req)
		if err == nil {
			t.Error("应该返回超时错误")
		}
		t.Log("✓ 超时测试通过")
	})

	t.Run("未超时", func(t *testing.T) {
		handler := TimeoutMiddleware(500 * time.Millisecond)(slowHandler)
		req := &Request{Data: map[string]interface{}{}}
		resp, err := handler(context.Background(), req)
		if err != nil {
			t.Fatalf("不应该返回错误: %v", err)
		}
		if resp.Data["status"] != "ok" {
			t.Error("响应不正确")
		}
		t.Log("✓ 未超时测试通过")
	})
}

// TestMetricsMiddleware 测试指标装饰器
func TestMetricsMiddleware(t *testing.T) {
	counter := &RequestCounter{}

	handler := MetricsMiddleware(counter)(BusinessHandler)

	ctx := context.WithValue(context.Background(), "userID", "user123")

	// 发送多个请求
	for i := 0; i < 5; i++ {
		req := &Request{Data: map[string]interface{}{"orderID": fmt.Sprintf("order%d", i)}}
		handler(ctx, req)
	}

	total, errors, avgTime := counter.Stats()
	if total != 5 {
		t.Errorf("期望总请求数为 5，实际为 %d", total)
	}
	if errors != 0 {
		t.Errorf("期望错误数为 0，实际为 %d", errors)
	}
	if avgTime <= 0 {
		t.Error("平均时间应该大于 0")
	}

	t.Logf("✓ 指标装饰器测试通过: total=%d, errors=%d, avgTime=%v", total, errors, avgTime)
}

// TestChain 测试装饰器链
func TestChain(t *testing.T) {
	// 创建一个会 panic 的 handler
	panicHandler := func(ctx context.Context, req *Request) (*Response, error) {
		panic("test panic")
	}

	// 组合多个装饰器：日志 -> 恢复 -> 认证
	validator := func(token string) (string, bool) {
		return "user123", true
	}

	handler := Chain(
		LoggingMiddleware,
		RecoveryMiddleware,
		AuthMiddleware(validator),
	)(panicHandler)

	req := &Request{Data: map[string]interface{}{"token": "valid"}}
	resp, err := handler(context.Background(), req)

	if resp != nil {
		t.Error("panic 后应该返回 nil")
	}
	if err == nil {
		t.Error("应该返回错误")
	}

	t.Log("✓ 装饰器链测试通过")
}

// TestChainExecutionOrder 测试装饰器链执行顺序
func TestChainExecutionOrder(t *testing.T) {
	var order []string

	middleware1 := func(next Handler) Handler {
		return func(ctx context.Context, req *Request) (*Response, error) {
			order = append(order, "before1")
			resp, err := next(ctx, req)
			order = append(order, "after1")
			return resp, err
		}
	}

	middleware2 := func(next Handler) Handler {
		return func(ctx context.Context, req *Request) (*Response, error) {
			order = append(order, "before2")
			resp, err := next(ctx, req)
			order = append(order, "after2")
			return resp, err
		}
	}

	handler := func(ctx context.Context, req *Request) (*Response, error) {
		order = append(order, "handler")
		return &Response{}, nil
	}

	chained := Chain(middleware1, middleware2)(handler)
	chained(context.Background(), &Request{})

	// 期望顺序: before1 -> before2 -> handler -> after2 -> after1
	expected := []string{"before1", "before2", "handler", "after2", "after1"}
	for i, v := range expected {
		if order[i] != v {
			t.Errorf("位置 %d: 期望 %s, 实际 %s", i, v, order[i])
		}
	}

	t.Log("✓ 执行顺序测试通过")
}

// TestHTTPMiddleware 测试 HTTP 风格中间件
func TestHTTPMiddleware(t *testing.T) {
	// 创建测试 handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 应用中间件
	chained := ChainHTTP(handler, LoggingHTTPMiddleware)

	// 创建测试请求
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	chained.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("期望状态码 200，实际 %d", rr.Code)
	}

	t.Log("✓ HTTP 中间件测试通过")
}

// TestAuthHTTPMiddleware 测试 HTTP 认证中间件
func TestAuthHTTPMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	chained := ChainHTTP(handler, AuthHTTPMiddleware)

	t.Run("无认证头", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()
		chained.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("期望状态码 401，实际 %d", rr.Code)
		}
		t.Log("✓ 无认证头测试通过")
	})

	t.Run("有认证头", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer token123")
		rr := httptest.NewRecorder()
		chained.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("期望状态码 200，实际 %d", rr.Code)
		}
		t.Log("✓ 有认证头测试通过")
	})
}

// TestGenericMiddleware 测试泛型中间件
func TestGenericMiddleware(t *testing.T) {
	handler := func(req int) (string, error) {
		return fmt.Sprintf("result: %d", req*2), nil
	}

	logged := GenericLogging(handler)

	result, err := logged(21)
	if err != nil {
		t.Fatalf("处理失败: %v", err)
	}
	if result != "result: 42" {
		t.Errorf("结果不正确: %s", result)
	}

	t.Log("✓ 泛型中间件测试通过")
}

// TestGenericChain 测试泛型中间件链
func TestGenericChain(t *testing.T) {
	handler := func(req int) (int, error) {
		return req * 2, nil
	}

	addMiddleware := func(n int) GenericMiddleware[int, int] {
		return func(next GenericHandler[int, int]) GenericHandler[int, int] {
			return func(req int) (int, error) {
				return next(req + n)
			}
		}
	}

	chained := GenericChain(addMiddleware(5), addMiddleware(3))(handler)

	// (5 + 3 + 10) * 2 = 36
	result, err := chained(10)
	if err != nil {
		t.Fatalf("处理失败: %v", err)
	}
	if result != 36 {
		t.Errorf("期望结果 36，实际 %d", result)
	}

	t.Log("✓ 泛型中间件链测试通过")
}

// BenchmarkMiddlewareChain 基准测试：中间件链性能
func BenchmarkMiddlewareChain(b *testing.B) {
	handler := BusinessHandler

	// 添加多个中间件
	chained := Chain(
		LoggingMiddleware,
		RecoveryMiddleware,
	)(handler)

	ctx := context.WithValue(context.Background(), "userID", "user123")
	req := &Request{Data: map[string]interface{}{"orderID": "order123"}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		chained(ctx, req)
	}
}

// BenchmarkHTTPMiddleware 基准测试：HTTP 中间件性能
func BenchmarkHTTPMiddleware(b *testing.B) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	chained := ChainHTTP(handler, LoggingHTTPMiddleware)

	req := httptest.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		chained.ServeHTTP(rr, req)
	}
}

// ExampleChain 示例：装饰器链用法
func ExampleChain() {
	// 创建一个简单的处理器
	handler := func(ctx context.Context, req *Request) (*Response, error) {
		return &Response{Data: map[string]interface{}{"result": "ok"}}, nil
	}

	// 使用 TimeoutMiddleware 包装
	wrapped := TimeoutMiddleware(100 * time.Millisecond)(handler)

	ctx := context.Background()
	req := &Request{Data: map[string]interface{}{}}

	resp, err := wrapped(ctx, req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Success: %v\n", resp.Data["result"])
	// Output: Success: ok
}

// ExampleTimeoutMiddleware 示例：超时装饰器用法
func ExampleTimeoutMiddleware() {
	slowHandler := func(ctx context.Context, req *Request) (*Response, error) {
		time.Sleep(200 * time.Millisecond)
		return &Response{Data: map[string]interface{}{"status": "ok"}}, nil
	}

	handler := TimeoutMiddleware(50 * time.Millisecond)(slowHandler)

	ctx := context.Background()
	req := &Request{Data: map[string]interface{}{}}

	_, err := handler(ctx, req)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
	}
	// Output: Request failed: request timeout after 50ms
}
