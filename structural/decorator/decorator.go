// Package decorator 展示了 Go 语言中实现装饰器模式的惯用方法。
//
// 装饰器模式允许在不修改原有对象的情况下，动态地给对象添加额外的职责。
// 在 Go 中，我们利用高阶函数（Higher-order functions）来实现这一模式，
// 这也是 Go HTTP 中间件（Middleware）的核心实现方式。
//
// 关键概念：
//   - 高阶函数：接收函数作为参数，返回函数的函数
//   - 函数组合：将多个装饰器链式组合
//   - 洋葱模型：请求像穿过洋葱一样逐层经过装饰器
//
// Go 语言特色：
//   - 函数是一等公民，可以像值一样传递
//   - 闭包捕获上下文
//   - 与 interface 结合实现类型安全
package decorator

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// ============================================================================
// 核心接口定义
// ============================================================================

// Handler 是处理函数的核心接口，类似于 http.HandlerFunc。
//
// 设计决策：
//   1. 使用 context.Context 传递上下文（Go 惯用法）
//   2. 返回 error 便于错误处理
//   3. 使用 Request/Response 结构体封装数据
type Handler func(ctx context.Context, req *Request) (*Response, error)

// Request 表示请求数据
type Request struct {
	Data map[string]interface{}
}

// Response 表示响应数据
type Response struct {
	Data    map[string]interface{}
	Headers map[string]string
}

// Middleware 是装饰器类型，接收一个 Handler 并返回包装后的 Handler。
//
// 这是装饰器模式在 Go 中的核心表达：
//   - 输入：被装饰的 Handler
//   - 输出：装饰后的 Handler（增加了新功能）
type Middleware func(Handler) Handler

// ============================================================================
// 具体业务 Handler 实现
// ============================================================================

// BusinessHandler 是核心业务逻辑处理器。
// 模拟一个处理用户订单的 Handler。
func BusinessHandler(ctx context.Context, req *Request) (*Response, error) {
	// 模拟业务处理
	userID := ctx.Value("userID")
	if userID == nil {
		return nil, fmt.Errorf("userID not found in context")
	}

	return &Response{
		Data: map[string]interface{}{
			"message": "order processed",
			"userID":  userID,
		},
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

// ============================================================================
// 基础装饰器实现
// ============================================================================

// LoggingMiddleware 日志记录装饰器。
//
// 功能：
//   - 记录请求开始和结束时间
//   - 记录处理耗时
//   - 记录请求和响应数据
func LoggingMiddleware(next Handler) Handler {
	return func(ctx context.Context, req *Request) (*Response, error) {
		start := time.Now()
		fmt.Printf("[LOG] Request started at %s, data: %v\n", start.Format(time.RFC3339), req.Data)

		resp, err := next(ctx, req)

		duration := time.Since(start)
		if err != nil {
			fmt.Printf("[LOG] Request failed after %v, error: %v\n", duration, err)
		} else {
			fmt.Printf("[LOG] Request completed in %v, response: %v\n", duration, resp.Data)
		}

		return resp, err
	}
}

// AuthMiddleware 认证装饰器。
//
// 功能：
//   - 验证用户身份
//   - 将用户信息注入 context
//   - 未认证则返回错误
func AuthMiddleware(tokenValidator func(string) (string, bool)) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req *Request) (*Response, error) {
			// 从请求头或数据中获取 token
			token, ok := req.Data["token"].(string)
			if !ok || token == "" {
				return nil, fmt.Errorf("authentication required: no token provided")
			}

			// 验证 token
			userID, valid := tokenValidator(token)
			if !valid {
				return nil, fmt.Errorf("authentication failed: invalid token")
			}

			// 将用户信息注入 context
			ctx = context.WithValue(ctx, "userID", userID)

			return next(ctx, req)
		}
	}
}

// RecoveryMiddleware 恐慌恢复装饰器。
//
// 功能：
//   - 捕获 panic 防止程序崩溃
//   - 将 panic 转换为 error
//   - 记录堆栈信息（实际项目中）
func RecoveryMiddleware(next Handler) Handler {
	return func(ctx context.Context, req *Request) (resp *Response, err error) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("[RECOVERY] Panic recovered: %v\n", r)
				err = fmt.Errorf("internal server error: %v", r)
			}
		}()

		return next(ctx, req)
	}
}

// TimeoutMiddleware 超时控制装饰器。
//
// 功能：
//   - 设置请求处理超时时间
//   - 超时后取消 context
//   - 返回超时错误
func TimeoutMiddleware(timeout time.Duration) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req *Request) (*Response, error) {
			// 创建带超时的 context
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			// 使用 channel 处理并发
			type result struct {
				resp *Response
				err  error
			}
			done := make(chan result, 1)

			go func() {
				resp, err := next(ctx, req)
				done <- result{resp, err}
			}()

			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("request timeout after %v", timeout)
			case r := <-done:
				return r.resp, r.err
			}
		}
	}
}

// MetricsMiddleware 指标收集装饰器。
//
// 功能：
//   - 记录请求计数
//   - 记录响应时间分布
//   - 记录错误率
func MetricsMiddleware(counter *RequestCounter) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req *Request) (*Response, error) {
			start := time.Now()
			counter.IncrementTotal()

			resp, err := next(ctx, req)

			duration := time.Since(start)
			counter.RecordDuration(duration)

			if err != nil {
				counter.IncrementError()
			}

			return resp, err
		}
	}
}

// RequestCounter 请求计数器
type RequestCounter struct {
	total  int
	errors int
	times  []time.Duration
}

// IncrementTotal 增加总请求数
func (c *RequestCounter) IncrementTotal() {
	c.total++
}

// IncrementError 增加错误数
func (c *RequestCounter) IncrementError() {
	c.errors++
}

// RecordDuration 记录处理时间
func (c *RequestCounter) RecordDuration(d time.Duration) {
	c.times = append(c.times, d)
}

// Stats 返回统计信息
func (c *RequestCounter) Stats() (total, errors int, avgTime time.Duration) {
	if len(c.times) == 0 {
		return c.total, c.errors, 0
	}
	var sum time.Duration
	for _, t := range c.times {
		sum += t
	}
	return c.total, c.errors, sum / time.Duration(len(c.times))
}

// ============================================================================
// 装饰器组合工具
// ============================================================================

// Chain 将多个中间件组合成一个中间件。
//
// 执行顺序：从左到右（先应用的装饰器先执行）
//
// 使用示例：
//
//	handler := Chain(
//	    LoggingMiddleware,
//	    RecoveryMiddleware,
//	    AuthMiddleware(validator),
//	)(BusinessHandler)
func Chain(middlewares ...Middleware) Middleware {
	return func(final Handler) Handler {
		// 从后向前包装，确保执行顺序正确
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}

// ChainReverse 将多个中间件反向组合。
//
// 执行顺序：从右到左（后应用的装饰器先执行）
func ChainReverse(middlewares ...Middleware) Middleware {
	return func(final Handler) Handler {
		for _, m := range middlewares {
			final = m(final)
		}
		return final
	}
}

// ============================================================================
// HTTP 风格实现（与标准库兼容）
// ============================================================================

// HTTPMiddleware 是标准库 http.Handler 风格的中间件
type HTTPMiddleware func(http.Handler) http.Handler

// LoggingHTTPMiddleware HTTP 日志中间件
func LoggingHTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		fmt.Printf("[HTTP] %s %s started\n", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		fmt.Printf("[HTTP] %s %s completed in %v\n", r.Method, r.URL.Path, time.Since(start))
	})
}

// AuthHTTPMiddleware HTTP 认证中间件
func AuthHTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// 实际项目中验证 token...
		next.ServeHTTP(w, r)
	})
}

// ChainHTTP 组合 HTTP 中间件
func ChainHTTP(handler http.Handler, middlewares ...HTTPMiddleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// ============================================================================
// 泛型装饰器（Go 1.18+）
// ============================================================================

// GenericHandler 是泛型 Handler 接口
type GenericHandler[T, R any] func(T) (R, error)

// GenericMiddleware 是泛型中间件
type GenericMiddleware[T, R any] func(GenericHandler[T, R]) GenericHandler[T, R]

// GenericLogging 泛型日志装饰器
func GenericLogging[T, R any](next GenericHandler[T, R]) GenericHandler[T, R] {
	return func(req T) (R, error) {
		fmt.Printf("[GENERIC] Request: %v\n", req)
		resp, err := next(req)
		if err != nil {
			fmt.Printf("[GENERIC] Error: %v\n", err)
		} else {
			fmt.Printf("[GENERIC] Response: %v\n", resp)
		}
		return resp, err
	}
}

// GenericChain 泛型中间件链
func GenericChain[T, R any](middlewares ...GenericMiddleware[T, R]) GenericMiddleware[T, R] {
	return func(final GenericHandler[T, R]) GenericHandler[T, R] {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}
