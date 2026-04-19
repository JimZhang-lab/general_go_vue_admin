/*
 * @Author: JimZhang
 * @Date: 2025-07-29 17:10:00
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-29 17:10:00
 * @FilePath: /server/middleware/rate_limit_middleware.go
 * @Description: 限流中间件
 */

package middleware

import (
	"fmt"
	"net/http"
	"server/common/result"
	"server/pkg/breaker"
	"server/pkg/limiter"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	// 限流器类型: "token", "sliding", "fixed"
	LimiterType string
	// 是否启用IP限流
	EnableIPLimit bool
	// 是否启用用户限流
	EnableUserLimit bool
	// 是否启用API限流
	EnableAPILimit bool
	// 自定义键生成函数
	KeyGenerator func(*gin.Context) string
	// 限流触发时的响应
	OnLimitReached func(*gin.Context)
	// 跳过限流的路径
	SkipPaths []string
}

// CircuitBreakerConfig 熔断器配置
type CircuitBreakerConfig struct {
	// 是否启用熔断器
	Enabled bool
	// 熔断器名称前缀
	NamePrefix string
	// 最大请求数（半开状态）
	MaxRequests uint32
	// 统计时间窗口
	Interval time.Duration
	// 熔断超时时间
	Timeout time.Duration
	// 自定义熔断条件
	ReadyToTrip func(breaker.Counts) bool
	// 跳过熔断的路径
	SkipPaths []string
}

// DefaultRateLimitConfig 默认限流配置
var DefaultRateLimitConfig = RateLimitConfig{
	LimiterType:     "sliding",
	EnableIPLimit:   true,
	EnableUserLimit: true,
	EnableAPILimit:  true,
	KeyGenerator:    defaultKeyGenerator,
	OnLimitReached:  defaultOnLimitReached,
	SkipPaths:       []string{"/health", "/metrics"},
}

// DefaultCircuitBreakerConfig 默认熔断器配置
var DefaultCircuitBreakerConfig = CircuitBreakerConfig{
	Enabled:     true,
	NamePrefix:  "api",
	MaxRequests: 10,
	Interval:    time.Minute,
	Timeout:     time.Minute * 2,
	ReadyToTrip: defaultReadyToTrip,
	SkipPaths:   []string{"/health", "/metrics"},
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(config ...RateLimitConfig) gin.HandlerFunc {
	cfg := DefaultRateLimitConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	// 获取限流器管理器
	limiterManager := limiter.GetLimiterManager()
	if limiterManager == nil {
		panic("限流器管理器未初始化")
	}

	return func(c *gin.Context) {
		// 检查是否跳过限流
		if shouldSkip(c.Request.URL.Path, cfg.SkipPaths) {
			c.Next()
			return
		}

		// 生成限流键
		keys := generateLimitKeys(c, cfg)

		// 检查每个键的限流状态
		for keyType, key := range keys {
			allowed, err := limiterManager.AllowWithType(key, cfg.LimiterType)
			if err != nil {
				// 限流器错误，记录日志但不阻止请求
				c.Header("X-RateLimit-Error", err.Error())
				continue
			}

			if !allowed {
				// 触发限流
				c.Header("X-RateLimit-Limit-Type", keyType)
				c.Header("X-RateLimit-Key", key)

				if cfg.OnLimitReached != nil {
					cfg.OnLimitReached(c)
				} else {
					defaultOnLimitReached(c)
				}
				return
			}
		}

		// 添加限流信息到响应头
		for keyType, key := range keys {
			// TODO: 实现GetStats方法或移除此功能
			c.Header(fmt.Sprintf("X-RateLimit-%s-Key", keyType), key)
		}

		c.Next()
	}
}

// CircuitBreakerMiddleware 熔断器中间件
func CircuitBreakerMiddleware(config ...CircuitBreakerConfig) gin.HandlerFunc {
	cfg := DefaultCircuitBreakerConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	if !cfg.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// 获取熔断器管理器
	breakerManager := breaker.GetBreakerManager()
	if breakerManager == nil {
		panic("熔断器管理器未初始化")
	}

	return func(c *gin.Context) {
		// 检查是否跳过熔断
		if shouldSkip(c.Request.URL.Path, cfg.SkipPaths) {
			c.Next()
			return
		}

		// 生成熔断器名称
		breakerName := generateBreakerName(c, cfg)

		// 获取或创建熔断器
		circuitBreaker, exists := breakerManager.GetCircuitBreaker(breakerName)
		if !exists {
			circuitBreaker = breakerManager.NewCircuitBreaker(breaker.Config{
				Name:        breakerName,
				MaxRequests: cfg.MaxRequests,
				Interval:    cfg.Interval,
				Timeout:     cfg.Timeout,
				ReadyToTrip: cfg.ReadyToTrip,
			})
		}

		// 使用熔断器执行请求
		err := circuitBreaker.Execute(func() error {
			c.Next()

			// 检查响应状态码，5xx认为是失败
			if c.Writer.Status() >= 500 {
				return fmt.Errorf("服务器内部错误: %d", c.Writer.Status())
			}

			return nil
		})

		if err != nil {
			// 熔断器阻止了请求或请求失败
			stats := circuitBreaker.GetStats()
			c.Header("X-Circuit-Breaker-State", stats.State.String())
			c.Header("X-Circuit-Breaker-Name", breakerName)

			if stats.State.String() == "OPEN" {
				// 熔断器开启，返回503
				result.Failed(c, http.StatusServiceUnavailable, "服务暂时不可用，请稍后重试")
				c.Abort()
				return
			}
		}

		// 添加熔断器状态到响应头
		stats := circuitBreaker.GetStats()
		c.Header("X-Circuit-Breaker-State", stats.State.String())
		c.Header("X-Circuit-Breaker-Requests", fmt.Sprintf("%d", stats.TotalRequests))
		c.Header("X-Circuit-Breaker-Failures", fmt.Sprintf("%d", stats.FailedRequests))
	}
}

// CombinedMiddleware 组合中间件（限流+熔断）
func CombinedMiddleware(rateLimitConfig RateLimitConfig, breakerConfig CircuitBreakerConfig) gin.HandlerFunc {
	rateLimitMW := RateLimitMiddleware(rateLimitConfig)
	breakerMW := CircuitBreakerMiddleware(breakerConfig)

	return func(c *gin.Context) {
		// 先执行限流检查
		rateLimitMW(c)
		if c.IsAborted() {
			return
		}

		// 再执行熔断检查
		breakerMW(c)
	}
}

// 辅助函数

// shouldSkip 检查是否应该跳过中间件
func shouldSkip(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// generateLimitKeys 生成限流键
func generateLimitKeys(c *gin.Context, cfg RateLimitConfig) map[string]string {
	keys := make(map[string]string)

	if cfg.EnableIPLimit {
		keys["IP"] = fmt.Sprintf("ip:%s", c.ClientIP())
	}

	if cfg.EnableUserLimit {
		if userID, exists := c.Get("user_id"); exists {
			keys["User"] = fmt.Sprintf("user:%v", userID)
		}
	}

	if cfg.EnableAPILimit {
		keys["API"] = fmt.Sprintf("api:%s:%s", c.Request.Method, c.Request.URL.Path)
	}

	// 使用自定义键生成器
	if cfg.KeyGenerator != nil {
		customKey := cfg.KeyGenerator(c)
		if customKey != "" {
			keys["Custom"] = customKey
		}
	}

	return keys
}

// generateBreakerName 生成熔断器名称
func generateBreakerName(c *gin.Context, cfg CircuitBreakerConfig) string {
	return fmt.Sprintf("%s:%s:%s", cfg.NamePrefix, c.Request.Method, c.Request.URL.Path)
}

// defaultKeyGenerator 默认键生成器
func defaultKeyGenerator(c *gin.Context) string {
	return fmt.Sprintf("global:%s", c.ClientIP())
}

// defaultOnLimitReached 默认限流触发处理
func defaultOnLimitReached(c *gin.Context) {
	result.Failed(c, http.StatusTooManyRequests, "请求过于频繁，请稍后重试")
	c.Abort()
}

// defaultReadyToTrip 默认熔断触发条件
func defaultReadyToTrip(counts breaker.Counts) bool {
	failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
	return counts.Requests >= 10 && failureRatio >= 0.5
}

// GetLimitStats 获取限流统计信息
func GetLimitStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		limiterManager := limiter.GetLimiterManager()
		if limiterManager == nil {
			result.Failed(c, http.StatusInternalServerError, "限流器管理器未初始化")
			return
		}

		// 这里可以实现获取所有限流统计信息的逻辑
		// 由于当前实现中没有全局统计，返回基本信息
		stats := map[string]interface{}{
			"status": "active",
			"type":   "rate_limiter",
		}

		result.Success(c, stats)
	}
}

// GetBreakerStats 获取熔断器统计信息
func GetBreakerStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		breakerManager := breaker.GetBreakerManager()
		if breakerManager == nil {
			result.Failed(c, http.StatusInternalServerError, "熔断器管理器未初始化")
			return
		}

		allBreakers := breakerManager.GetAllBreakers()
		stats := make(map[string]*breaker.BreakerStats)

		for name, cb := range allBreakers {
			stats[name] = cb.GetStats()
		}

		result.Success(c, stats)
	}
}

// ResetBreaker 重置熔断器
func ResetBreaker() gin.HandlerFunc {
	return func(c *gin.Context) {
		breakerName := c.Param("name")
		if breakerName == "" {
			result.Failed(c, http.StatusBadRequest, "熔断器名称不能为空")
			return
		}

		breakerManager := breaker.GetBreakerManager()
		if breakerManager == nil {
			result.Failed(c, http.StatusInternalServerError, "熔断器管理器未初始化")
			return
		}

		circuitBreaker, exists := breakerManager.GetCircuitBreaker(breakerName)
		if !exists {
			result.Failed(c, http.StatusNotFound, "熔断器不存在")
			return
		}

		err := circuitBreaker.Reset()
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, fmt.Sprintf("重置熔断器失败: %v", err))
			return
		}

		result.Success(c, gin.H{"message": "熔断器重置成功"})
	}
}
