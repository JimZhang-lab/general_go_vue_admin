/*
 * @Author: JimZhang
 * @Date: 2025-07-29 16:45:00
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-29 16:45:00
 * @FilePath: /server/pkg/limiter/rate_limiter.go
 * @Description: API限流器
 */

package limiter

import (
	"context"
	"fmt"
	"log"
	"server/pkg/redis"
	"strconv"
	"sync"
	"time"

	redisClient "github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

// RateLimiter 限流器接口
type RateLimiter interface {
	Allow(key string) (bool, error)
	AllowN(key string, n int) (bool, error)
	Reset(key string) error
	GetStats(key string) (*LimitStats, error)
}

// LimitStats 限流统计
type LimitStats struct {
	Key           string    `json:"key"`
	Requests      int64     `json:"requests"`
	Allowed       int64     `json:"allowed"`
	Rejected      int64     `json:"rejected"`
	LastRequest   time.Time `json:"last_request"`
	ResetTime     time.Time `json:"reset_time"`
	RemainingQuota int64    `json:"remaining_quota"`
}

// TokenBucketLimiter 令牌桶限流器
type TokenBucketLimiter struct {
	client    *redisClient.Client
	ctx       context.Context
	buckets   map[string]*rate.Limiter
	mu        sync.RWMutex
	rate      rate.Limit
	burst     int
	keyPrefix string
}

// SlidingWindowLimiter 滑动窗口限流器
type SlidingWindowLimiter struct {
	client    *redisClient.Client
	ctx       context.Context
	window    time.Duration
	limit     int64
	keyPrefix string
}

// FixedWindowLimiter 固定窗口限流器
type FixedWindowLimiter struct {
	client    *redisClient.Client
	ctx       context.Context
	window    time.Duration
	limit     int64
	keyPrefix string
}

// LimiterManager 限流器管理器
type LimiterManager struct {
	tokenBucket   *TokenBucketLimiter
	slidingWindow *SlidingWindowLimiter
	fixedWindow   *FixedWindowLimiter
	defaultType   string
}

var (
	limiterManager *LimiterManager
	limiterOnce    sync.Once
)

// InitLimiterManager 初始化限流器管理器
func InitLimiterManager() (*LimiterManager, error) {
	var err error
	limiterOnce.Do(func() {
		limiterManager, err = newLimiterManager()
	})
	return limiterManager, err
}

// GetLimiterManager 获取限流器管理器实例
func GetLimiterManager() *LimiterManager {
	return limiterManager
}

// newLimiterManager 创建新的限流器管理器
func newLimiterManager() (*LimiterManager, error) {
	if redis.RedisDb == nil {
		return nil, fmt.Errorf("Redis客户端未初始化")
	}

	ctx := context.Background()

	// 创建令牌桶限流器
	tokenBucket := &TokenBucketLimiter{
		client:    redis.RedisDb,
		ctx:       ctx,
		buckets:   make(map[string]*rate.Limiter),
		rate:      rate.Limit(100), // 每秒100个请求
		burst:     200,             // 突发200个请求
		keyPrefix: "rate_limit:token:",
	}

	// 创建滑动窗口限流器
	slidingWindow := &SlidingWindowLimiter{
		client:    redis.RedisDb,
		ctx:       ctx,
		window:    time.Minute,     // 1分钟窗口
		limit:     1000,            // 1000个请求
		keyPrefix: "rate_limit:sliding:",
	}

	// 创建固定窗口限流器
	fixedWindow := &FixedWindowLimiter{
		client:    redis.RedisDb,
		ctx:       ctx,
		window:    time.Minute,     // 1分钟窗口
		limit:     1000,            // 1000个请求
		keyPrefix: "rate_limit:fixed:",
	}

	manager := &LimiterManager{
		tokenBucket:   tokenBucket,
		slidingWindow: slidingWindow,
		fixedWindow:   fixedWindow,
		defaultType:   "sliding", // 默认使用滑动窗口
	}

	log.Println("限流器管理器初始化成功")
	return manager, nil
}

// Allow 检查是否允许请求
func (lm *LimiterManager) Allow(key string) (bool, error) {
	switch lm.defaultType {
	case "token":
		return lm.tokenBucket.Allow(key)
	case "sliding":
		return lm.slidingWindow.Allow(key)
	case "fixed":
		return lm.fixedWindow.Allow(key)
	default:
		return lm.slidingWindow.Allow(key)
	}
}

// AllowWithType 使用指定类型的限流器检查请求
func (lm *LimiterManager) AllowWithType(key, limiterType string) (bool, error) {
	switch limiterType {
	case "token":
		return lm.tokenBucket.Allow(key)
	case "sliding":
		return lm.slidingWindow.Allow(key)
	case "fixed":
		return lm.fixedWindow.Allow(key)
	default:
		return false, fmt.Errorf("未知的限流器类型: %s", limiterType)
	}
}

// TokenBucketLimiter 实现

// Allow 检查令牌桶是否允许请求
func (tbl *TokenBucketLimiter) Allow(key string) (bool, error) {
	return tbl.AllowN(key, 1)
}

// AllowN 检查令牌桶是否允许N个请求
func (tbl *TokenBucketLimiter) AllowN(key string, n int) (bool, error) {
	tbl.mu.Lock()
	defer tbl.mu.Unlock()

	limiter, exists := tbl.buckets[key]
	if !exists {
		limiter = rate.NewLimiter(tbl.rate, tbl.burst)
		tbl.buckets[key] = limiter
	}

	// 检查是否允许
	allowed := limiter.AllowN(time.Now(), n)

	// 更新Redis中的统计信息
	go tbl.updateStats(key, int64(n), allowed)

	return allowed, nil
}

// Reset 重置令牌桶
func (tbl *TokenBucketLimiter) Reset(key string) error {
	tbl.mu.Lock()
	defer tbl.mu.Unlock()

	delete(tbl.buckets, key)
	
	// 清除Redis中的统计信息
	fullKey := tbl.keyPrefix + key
	return tbl.client.Del(tbl.ctx, fullKey).Err()
}

// GetStats 获取令牌桶统计信息
func (tbl *TokenBucketLimiter) GetStats(key string) (*LimitStats, error) {
	fullKey := tbl.keyPrefix + key
	
	data, err := tbl.client.HGetAll(tbl.ctx, fullKey).Result()
	if err != nil {
		return nil, err
	}

	stats := &LimitStats{Key: key}
	if val, ok := data["requests"]; ok {
		stats.Requests, _ = strconv.ParseInt(val, 10, 64)
	}
	if val, ok := data["allowed"]; ok {
		stats.Allowed, _ = strconv.ParseInt(val, 10, 64)
	}
	if val, ok := data["rejected"]; ok {
		stats.Rejected, _ = strconv.ParseInt(val, 10, 64)
	}

	return stats, nil
}

// updateStats 更新统计信息
func (tbl *TokenBucketLimiter) updateStats(key string, requests int64, allowed bool) {
	fullKey := tbl.keyPrefix + key
	
	pipe := tbl.client.Pipeline()
	pipe.HIncrBy(tbl.ctx, fullKey, "requests", requests)
	
	if allowed {
		pipe.HIncrBy(tbl.ctx, fullKey, "allowed", requests)
	} else {
		pipe.HIncrBy(tbl.ctx, fullKey, "rejected", requests)
	}
	
	pipe.HSet(tbl.ctx, fullKey, "last_request", time.Now().Unix())
	pipe.Expire(tbl.ctx, fullKey, time.Hour)
	
	pipe.Exec(tbl.ctx)
}

// SlidingWindowLimiter 实现

// Allow 检查滑动窗口是否允许请求
func (swl *SlidingWindowLimiter) Allow(key string) (bool, error) {
	return swl.AllowN(key, 1)
}

// AllowN 检查滑动窗口是否允许N个请求
func (swl *SlidingWindowLimiter) AllowN(key string, n int) (bool, error) {
	fullKey := swl.keyPrefix + key
	now := time.Now()
	windowStart := now.Add(-swl.window)

	// 使用Lua脚本实现原子操作
	luaScript := `
		local key = KEYS[1]
		local window_start = ARGV[1]
		local now = ARGV[2]
		local limit = tonumber(ARGV[3])
		local increment = tonumber(ARGV[4])

		-- 清除过期的记录
		redis.call('ZREMRANGEBYSCORE', key, 0, window_start)
		
		-- 获取当前窗口内的请求数
		local current = redis.call('ZCARD', key)
		
		if current + increment <= limit then
			-- 添加新的请求记录
			for i = 1, increment do
				redis.call('ZADD', key, now, now .. ':' .. i)
			end
			redis.call('EXPIRE', key, 3600)
			return {1, current + increment}
		else
			return {0, current}
		end
	`

	result, err := swl.client.Eval(swl.ctx, luaScript, []string{fullKey},
		windowStart.UnixNano(),
		now.UnixNano(),
		swl.limit,
		n,
	).Result()

	if err != nil {
		return false, err
	}

	resultSlice := result.([]interface{})
	allowed := resultSlice[0].(int64) == 1

	// 更新统计信息
	go swl.updateStats(key, int64(n), allowed)

	return allowed, nil
}

// Reset 重置滑动窗口
func (swl *SlidingWindowLimiter) Reset(key string) error {
	fullKey := swl.keyPrefix + key
	return swl.client.Del(swl.ctx, fullKey).Err()
}

// GetStats 获取滑动窗口统计信息
func (swl *SlidingWindowLimiter) GetStats(key string) (*LimitStats, error) {
	fullKey := swl.keyPrefix + key
	
	// 获取当前窗口内的请求数
	now := time.Now()
	windowStart := now.Add(-swl.window)
	
	count, err := swl.client.ZCount(swl.ctx, fullKey, 
		strconv.FormatInt(windowStart.UnixNano(), 10),
		strconv.FormatInt(now.UnixNano(), 10),
	).Result()
	
	if err != nil {
		return nil, err
	}

	stats := &LimitStats{
		Key:            key,
		Requests:       count,
		RemainingQuota: swl.limit - count,
		ResetTime:      now.Add(swl.window),
	}

	return stats, nil
}

// updateStats 更新统计信息
func (swl *SlidingWindowLimiter) updateStats(key string, requests int64, allowed bool) {
	statsKey := swl.keyPrefix + "stats:" + key
	
	pipe := swl.client.Pipeline()
	pipe.HIncrBy(swl.ctx, statsKey, "total_requests", requests)
	
	if allowed {
		pipe.HIncrBy(swl.ctx, statsKey, "allowed", requests)
	} else {
		pipe.HIncrBy(swl.ctx, statsKey, "rejected", requests)
	}
	
	pipe.HSet(swl.ctx, statsKey, "last_request", time.Now().Unix())
	pipe.Expire(swl.ctx, statsKey, time.Hour)
	
	pipe.Exec(swl.ctx)
}

// FixedWindowLimiter 实现

// Allow 检查固定窗口是否允许请求
func (fwl *FixedWindowLimiter) Allow(key string) (bool, error) {
	return fwl.AllowN(key, 1)
}

// AllowN 检查固定窗口是否允许N个请求
func (fwl *FixedWindowLimiter) AllowN(key string, n int) (bool, error) {
	now := time.Now()
	window := now.Truncate(fwl.window)
	fullKey := fmt.Sprintf("%s%s:%d", fwl.keyPrefix, key, window.Unix())

	// 使用Lua脚本实现原子操作
	luaScript := `
		local key = KEYS[1]
		local limit = tonumber(ARGV[1])
		local increment = tonumber(ARGV[2])
		local ttl = tonumber(ARGV[3])

		local current = redis.call('GET', key)
		if current == false then
			current = 0
		else
			current = tonumber(current)
		end

		if current + increment <= limit then
			local new_val = redis.call('INCRBY', key, increment)
			redis.call('EXPIRE', key, ttl)
			return {1, new_val}
		else
			return {0, current}
		end
	`

	result, err := fwl.client.Eval(fwl.ctx, luaScript, []string{fullKey},
		fwl.limit,
		n,
		int(fwl.window.Seconds()),
	).Result()

	if err != nil {
		return false, err
	}

	resultSlice := result.([]interface{})
	allowed := resultSlice[0].(int64) == 1

	// 更新统计信息
	go fwl.updateStats(key, int64(n), allowed)

	return allowed, nil
}

// Reset 重置固定窗口
func (fwl *FixedWindowLimiter) Reset(key string) error {
	pattern := fwl.keyPrefix + key + ":*"
	keys, err := fwl.client.Keys(fwl.ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return fwl.client.Del(fwl.ctx, keys...).Err()
	}

	return nil
}

// GetStats 获取固定窗口统计信息
func (fwl *FixedWindowLimiter) GetStats(key string) (*LimitStats, error) {
	now := time.Now()
	window := now.Truncate(fwl.window)
	fullKey := fmt.Sprintf("%s%s:%d", fwl.keyPrefix, key, window.Unix())

	current, err := fwl.client.Get(fwl.ctx, fullKey).Int64()
	if err != nil && err != redisClient.Nil {
		return nil, err
	}

	stats := &LimitStats{
		Key:            key,
		Requests:       current,
		RemainingQuota: fwl.limit - current,
		ResetTime:      window.Add(fwl.window),
	}

	return stats, nil
}

// updateStats 更新统计信息
func (fwl *FixedWindowLimiter) updateStats(key string, requests int64, allowed bool) {
	statsKey := fwl.keyPrefix + "stats:" + key
	
	pipe := fwl.client.Pipeline()
	pipe.HIncrBy(fwl.ctx, statsKey, "total_requests", requests)
	
	if allowed {
		pipe.HIncrBy(fwl.ctx, statsKey, "allowed", requests)
	} else {
		pipe.HIncrBy(fwl.ctx, statsKey, "rejected", requests)
	}
	
	pipe.HSet(fwl.ctx, statsKey, "last_request", time.Now().Unix())
	pipe.Expire(fwl.ctx, statsKey, time.Hour)
	
	pipe.Exec(fwl.ctx)
}
