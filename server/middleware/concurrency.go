/*
 * @Author: JimZhang
 * @Date: 2025-07-27
 * @Description: 高并发处理中间件，包含限流、熔断、超时控制等
 */
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"server/common/errors"
	"server/common/result"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter 限流器接口
type RateLimiter interface {
	Allow() bool
	Wait(ctx context.Context) error
}

// TokenBucketLimiter 简化版令牌桶限流器
type TokenBucketLimiter struct {
	mu       sync.Mutex
	tokens   float64
	capacity float64
	rate     float64
	lastTime time.Time
}

// NewTokenBucketLimiter 创建令牌桶限流器
func NewTokenBucketLimiter(rps float64, burst int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		tokens:   float64(burst),
		capacity: float64(burst),
		rate:     rps,
		lastTime: time.Now(),
	}
}

func (t *TokenBucketLimiter) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(t.lastTime).Seconds()
	t.lastTime = now

	// 添加新的令牌
	t.tokens += elapsed * t.rate
	if t.tokens > t.capacity {
		t.tokens = t.capacity
	}

	// 检查是否有足够的令牌
	if t.tokens >= 1 {
		t.tokens--
		return true
	}

	return false
}

func (t *TokenBucketLimiter) Wait(ctx context.Context) error {
	for {
		if t.Allow() {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Millisecond * 10):
			// 继续等待
		}
	}
}

// SlidingWindowLimiter 滑动窗口限流器
type SlidingWindowLimiter struct {
	mu       sync.RWMutex
	requests []time.Time
	limit    int
	window   time.Duration
}

// NewSlidingWindowLimiter 创建滑动窗口限流器
func NewSlidingWindowLimiter(limit int, window time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		requests: make([]time.Time, 0),
		limit:    limit,
		window:   window,
	}
}

func (s *SlidingWindowLimiter) Allow() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-s.window)

	// 清理过期的请求记录
	validRequests := make([]time.Time, 0)
	for _, req := range s.requests {
		if req.After(cutoff) {
			validRequests = append(validRequests, req)
		}
	}
	s.requests = validRequests

	// 检查是否超过限制
	if len(s.requests) >= s.limit {
		return false
	}

	// 记录当前请求
	s.requests = append(s.requests, now)
	return true
}

func (s *SlidingWindowLimiter) Wait(ctx context.Context) error {
	for {
		if s.Allow() {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Millisecond * 100):
			// 继续等待
		}
	}
}

// CircuitBreaker 熔断器
type CircuitBreaker struct {
	mu                sync.RWMutex
	state             CircuitState
	failureCount      int
	successCount      int
	lastFailureTime   time.Time
	failureThreshold  int
	recoveryThreshold int
	timeout           time.Duration
}

// CircuitState 熔断器状态
type CircuitState int

const (
	StateClosed CircuitState = iota
	StateOpen
	StateHalfOpen
)

// NewCircuitBreaker 创建熔断器
func NewCircuitBreaker(failureThreshold, recoveryThreshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:             StateClosed,
		failureThreshold:  failureThreshold,
		recoveryThreshold: recoveryThreshold,
		timeout:           timeout,
	}
}

// Execute 执行操作
func (cb *CircuitBreaker) Execute(fn func() error) error {
	if !cb.canExecute() {
		return errors.New(errors.ErrCircuitBreaker, "熔断器开启，请求被拒绝")
	}

	err := fn()
	cb.recordResult(err)
	return err
}

// canExecute 检查是否可以执行
func (cb *CircuitBreaker) canExecute() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		return time.Since(cb.lastFailureTime) > cb.timeout
	case StateHalfOpen:
		return true
	default:
		return false
	}
}

// recordResult 记录执行结果
func (cb *CircuitBreaker) recordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failureCount++
		cb.lastFailureTime = time.Now()

		if cb.state == StateClosed && cb.failureCount >= cb.failureThreshold {
			cb.state = StateOpen
		} else if cb.state == StateHalfOpen {
			cb.state = StateOpen
		}
	} else {
		cb.successCount++

		if cb.state == StateHalfOpen && cb.successCount >= cb.recoveryThreshold {
			cb.state = StateClosed
			cb.failureCount = 0
			cb.successCount = 0
		}
	}

	// 从开启状态转换到半开状态
	if cb.state == StateOpen && time.Since(cb.lastFailureTime) > cb.timeout {
		cb.state = StateHalfOpen
		cb.successCount = 0
	}
}

// GetState 获取当前状态
func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// ConcurrencyConfig 并发控制配置
type ConcurrencyConfig struct {
	// 限流配置
	EnableRateLimit   bool
	RateLimitType     string        // "token_bucket" 或 "sliding_window"
	RequestsPerSecond float64       // 每秒请求数
	BurstSize         int           // 突发大小
	WindowSize        time.Duration // 滑动窗口大小

	// 熔断器配置
	EnableCircuitBreaker bool
	FailureThreshold     int           // 失败阈值
	RecoveryThreshold    int           // 恢复阈值
	CircuitTimeout       time.Duration // 熔断超时时间

	// 超时控制
	EnableTimeout  bool
	RequestTimeout time.Duration

	// 并发控制
	EnableConcurrencyLimit bool
	MaxConcurrentRequests  int

	// 优雅关闭
	EnableGracefulShutdown bool
	ShutdownTimeout        time.Duration
}

// DefaultConcurrencyConfig 默认并发配置
var DefaultConcurrencyConfig = &ConcurrencyConfig{
	EnableRateLimit:        true,
	RateLimitType:          "token_bucket",
	RequestsPerSecond:      100,
	BurstSize:              200,
	WindowSize:             time.Minute,
	EnableCircuitBreaker:   true,
	FailureThreshold:       10,
	RecoveryThreshold:      5,
	CircuitTimeout:         time.Minute * 2,
	EnableTimeout:          true,
	RequestTimeout:         time.Second * 30,
	EnableConcurrencyLimit: true,
	MaxConcurrentRequests:  1000,
	EnableGracefulShutdown: true,
	ShutdownTimeout:        time.Second * 30,
}

// ConcurrencyManager 并发管理器
type ConcurrencyManager struct {
	config         *ConcurrencyConfig
	rateLimiter    RateLimiter
	circuitBreaker *CircuitBreaker
	semaphore      chan struct{}
	activeRequests int64
	mu             sync.RWMutex
}

// NewConcurrencyManager 创建并发管理器
func NewConcurrencyManager(config *ConcurrencyConfig) *ConcurrencyManager {
	if config == nil {
		config = DefaultConcurrencyConfig
	}

	cm := &ConcurrencyManager{
		config: config,
	}

	// 初始化限流器
	if config.EnableRateLimit {
		switch config.RateLimitType {
		case "sliding_window":
			cm.rateLimiter = NewSlidingWindowLimiter(
				int(config.RequestsPerSecond*config.WindowSize.Seconds()),
				config.WindowSize,
			)
		default:
			cm.rateLimiter = NewTokenBucketLimiter(
				config.RequestsPerSecond,
				config.BurstSize,
			)
		}
	}

	// 初始化熔断器
	if config.EnableCircuitBreaker {
		cm.circuitBreaker = NewCircuitBreaker(
			config.FailureThreshold,
			config.RecoveryThreshold,
			config.CircuitTimeout,
		)
	}

	// 初始化并发控制信号量
	if config.EnableConcurrencyLimit {
		cm.semaphore = make(chan struct{}, config.MaxConcurrentRequests)
	}

	return cm
}

// RateLimitMiddleware 限流中间件
func (cm *ConcurrencyManager) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !cm.config.EnableRateLimit {
			c.Next()
			return
		}

		if !cm.rateLimiter.Allow() {
			result.Failed(c, int(errors.ErrRateLimit), "请求频率过高，请稍后重试")
			c.Abort()
			return
		}

		c.Next()
	}
}

// CircuitBreakerMiddleware 熔断器中间件
func (cm *ConcurrencyManager) CircuitBreakerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !cm.config.EnableCircuitBreaker {
			c.Next()
			return
		}

		err := cm.circuitBreaker.Execute(func() error {
			c.Next()

			// 检查响应状态码，判断是否为失败
			if c.Writer.Status() >= 500 {
				return fmt.Errorf("服务器内部错误")
			}
			return nil
		})

		if err != nil {
			if appErr, ok := err.(*errors.AppError); ok && appErr.Code == errors.ErrCircuitBreaker {
				result.Failed(c, int(errors.ErrCircuitBreaker), "服务暂时不可用，请稍后重试")
				c.Abort()
				return
			}
		}
	}
}

// TimeoutMiddleware 超时控制中间件
func (cm *ConcurrencyManager) TimeoutMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !cm.config.EnableTimeout {
			c.Next()
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), cm.config.RequestTimeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{})
		go func() {
			defer close(done)
			c.Next()
		}()

		select {
		case <-done:
			// 请求正常完成
		case <-ctx.Done():
			// 请求超时
			result.Failed(c, int(errors.ErrTimeout), "请求超时")
			c.Abort()
		}
	}
}

// ConcurrencyLimitMiddleware 并发限制中间件
func (cm *ConcurrencyManager) ConcurrencyLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !cm.config.EnableConcurrencyLimit {
			c.Next()
			return
		}

		select {
		case cm.semaphore <- struct{}{}:
			defer func() { <-cm.semaphore }()
			c.Next()
		default:
			result.Failed(c, http.StatusTooManyRequests, "服务器繁忙，请稍后重试")
			c.Abort()
		}
	}
}

// GetActiveRequests 获取活跃请求数
func (cm *ConcurrencyManager) GetActiveRequests() int64 {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.activeRequests
}

// GetCircuitBreakerState 获取熔断器状态
func (cm *ConcurrencyManager) GetCircuitBreakerState() CircuitState {
	if cm.circuitBreaker == nil {
		return StateClosed
	}
	return cm.circuitBreaker.GetState()
}
