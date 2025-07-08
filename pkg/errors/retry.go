package errors

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// RetryConfig 重试配置
type RetryConfig struct {
	MaxAttempts     int           // 最大重试次数
	InitialDelay    time.Duration // 初始延迟
	MaxDelay        time.Duration // 最大延迟
	BackoffFactor   float64       // 退避因子
	Jitter          bool          // 是否添加抖动
	RetryableErrors []ErrorCode   // 可重试的错误码
}

// DefaultRetryConfig 默认重试配置
var DefaultRetryConfig = RetryConfig{
	MaxAttempts:   3,
	InitialDelay:  1 * time.Second,
	MaxDelay:      30 * time.Second,
	BackoffFactor: 2.0,
	Jitter:        true,
	RetryableErrors: []ErrorCode{
		ErrCodeDatabaseConnection,
		ErrCodeDatabaseTimeout,
		ErrCodeCacheConnection,
		ErrCodeCacheTimeout,
		ErrCodeWebhookTimeout,
		ErrCodeTimeout,
		ErrCodeServiceUnavailable,
	},
}

// RetryableFunc 可重试的函数类型
type RetryableFunc func() error

// RetryableFuncWithContext 带上下文的可重试函数类型
type RetryableFuncWithContext func(ctx context.Context) error

// Retry 执行重试逻辑
func Retry(fn RetryableFunc, config RetryConfig) error {
	return RetryWithContext(context.Background(), func(ctx context.Context) error {
		return fn()
	}, config)
}

// RetryWithContext 带上下文的执行重试逻辑
func RetryWithContext(ctx context.Context, fn RetryableFuncWithContext, config RetryConfig) error {
	var lastErr error

	for attempt := 0; attempt <= config.MaxAttempts; attempt++ {
		// 检查上下文是否已取消
		select {
		case <-ctx.Done():
			return NewTimeoutError("Retry cancelled by context", ctx.Err())
		default:
		}

		// 执行函数
		err := fn(ctx)
		if err == nil {
			return nil // 成功，返回
		}

		lastErr = err

		// 检查是否为可重试错误
		if !isRetryableError(err, config.RetryableErrors) {
			return err // 不可重试，直接返回
		}

		// 如果是最后一次尝试，返回错误
		if attempt == config.MaxAttempts {
			break
		}

		// 计算延迟时间
		delay := calculateDelay(attempt, config)

		// 等待延迟时间
		select {
		case <-time.After(delay):
			continue
		case <-ctx.Done():
			return NewTimeoutError("Retry cancelled by context", ctx.Err())
		}
	}

	// 所有重试都失败了
	return WrapError(lastErr, fmt.Sprintf("Retry failed after %d attempts", config.MaxAttempts+1))
}

// RetryWithExponentialBackoff 使用指数退避的重试
func RetryWithExponentialBackoff(fn RetryableFunc, maxAttempts int) error {
	config := DefaultRetryConfig
	config.MaxAttempts = maxAttempts
	return Retry(fn, config)
}

// RetryWithFixedDelay 使用固定延迟的重试
func RetryWithFixedDelay(fn RetryableFunc, maxAttempts int, delay time.Duration) error {
	config := RetryConfig{
		MaxAttempts:     maxAttempts,
		InitialDelay:    delay,
		MaxDelay:        delay,
		BackoffFactor:   1.0,
		Jitter:          false,
		RetryableErrors: DefaultRetryConfig.RetryableErrors,
	}
	return Retry(fn, config)
}

// isRetryableError 检查错误是否可重试
func isRetryableError(err error, retryableCodes []ErrorCode) bool {
	// 检查自定义错误
	if customErr, ok := err.(*CustomError); ok {
		for _, code := range retryableCodes {
			if customErr.Code == code {
				return true
			}
		}
		return false
	}

	// 检查通用错误类型
	return IsRetryableError(err)
}

// calculateDelay 计算延迟时间
func calculateDelay(attempt int, config RetryConfig) time.Duration {
	// 计算基础延迟
	delay := float64(config.InitialDelay) * math.Pow(config.BackoffFactor, float64(attempt))

	// 限制最大延迟
	if delay > float64(config.MaxDelay) {
		delay = float64(config.MaxDelay)
	}

	// 添加抖动
	if config.Jitter {
		jitter := delay * 0.1 * rand.Float64() // 10% 的抖动
		delay += jitter
	}

	return time.Duration(delay)
}

// RetryableOperation 可重试操作的结构体
type RetryableOperation struct {
	config RetryConfig
	ctx    context.Context
}

// NewRetryableOperation 创建新的可重试操作
func NewRetryableOperation(ctx context.Context, config RetryConfig) *RetryableOperation {
	return &RetryableOperation{
		config: config,
		ctx:    ctx,
	}
}

// Execute 执行可重试操作
func (r *RetryableOperation) Execute(fn RetryableFuncWithContext) error {
	return RetryWithContext(r.ctx, fn, r.config)
}

// WithMaxAttempts 设置最大重试次数
func (r *RetryableOperation) WithMaxAttempts(maxAttempts int) *RetryableOperation {
	r.config.MaxAttempts = maxAttempts
	return r
}

// WithInitialDelay 设置初始延迟
func (r *RetryableOperation) WithInitialDelay(delay time.Duration) *RetryableOperation {
	r.config.InitialDelay = delay
	return r
}

// WithMaxDelay 设置最大延迟
func (r *RetryableOperation) WithMaxDelay(delay time.Duration) *RetryableOperation {
	r.config.MaxDelay = delay
	return r
}

// WithBackoffFactor 设置退避因子
func (r *RetryableOperation) WithBackoffFactor(factor float64) *RetryableOperation {
	r.config.BackoffFactor = factor
	return r
}

// WithJitter 设置是否添加抖动
func (r *RetryableOperation) WithJitter(jitter bool) *RetryableOperation {
	r.config.Jitter = jitter
	return r
}

// WithRetryableErrors 设置可重试的错误码
func (r *RetryableOperation) WithRetryableErrors(codes []ErrorCode) *RetryableOperation {
	r.config.RetryableErrors = codes
	return r
}

// RetryableResult 重试结果
type RetryableResult[T any] struct {
	Result T
	Error  error
}

// RetryableFuncWithResult 带结果的可重试函数类型
type RetryableFuncWithResult[T any] func() (T, error)

// RetryableFuncWithContextAndResult 带上下文和结果的可重试函数类型
type RetryableFuncWithContextAndResult[T any] func(ctx context.Context) (T, error)

// RetryWithResult 执行带结果的重试逻辑
func RetryWithResult[T any](fn RetryableFuncWithResult[T], config RetryConfig) RetryableResult[T] {
	result, err := RetryWithContextAndResult(context.Background(), func(ctx context.Context) (T, error) {
		return fn()
	}, config)
	return RetryableResult[T]{Result: result, Error: err}
}

// RetryWithContextAndResult 带上下文的执行带结果的重试逻辑
func RetryWithContextAndResult[T any](ctx context.Context, fn RetryableFuncWithContextAndResult[T], config RetryConfig) (T, error) {
	var lastErr error
	var zero T

	for attempt := 0; attempt <= config.MaxAttempts; attempt++ {
		// 检查上下文是否已取消
		select {
		case <-ctx.Done():
			return zero, NewTimeoutError("Retry cancelled by context", ctx.Err())
		default:
		}

		// 执行函数
		result, err := fn(ctx)
		if err == nil {
			return result, nil // 成功，返回
		}

		lastErr = err

		// 检查是否为可重试错误
		if !isRetryableError(err, config.RetryableErrors) {
			return zero, err // 不可重试，直接返回
		}

		// 如果是最后一次尝试，返回错误
		if attempt == config.MaxAttempts {
			break
		}

		// 计算延迟时间
		delay := calculateDelay(attempt, config)

		// 等待延迟时间
		select {
		case <-time.After(delay):
			continue
		case <-ctx.Done():
			return zero, NewTimeoutError("Retry cancelled by context", ctx.Err())
		}
	}

	// 所有重试都失败了
	return zero, WrapError(lastErr, fmt.Sprintf("Retry failed after %d attempts", config.MaxAttempts+1))
}

// 初始化随机数种子
func init() {
	// Go 1.20+ 不再需要手动设置随机数种子
	// rand.Seed(time.Now().UnixNano())
}
