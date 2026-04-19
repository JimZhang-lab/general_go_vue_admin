/*
 * @Author: JimZhang
 * @Date: 2025-07-29 17:00:00
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-29 17:00:00
 * @FilePath: /server/pkg/breaker/circuit_breaker.go
 * @Description: 熔断器实现
 */

package breaker

import (
	"context"
	"fmt"
	"log"
	"server/pkg/redis"
	"sync"
	"time"

	redisClient "github.com/redis/go-redis/v9"
)

// State 熔断器状态
type State int

const (
	StateClosed State = iota
	StateHalfOpen
	StateOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateHalfOpen:
		return "HALF_OPEN"
	case StateOpen:
		return "OPEN"
	default:
		return "UNKNOWN"
	}
}

// CircuitBreaker 熔断器接口
type CircuitBreaker interface {
	Execute(fn func() error) error
	Call(fn func() (interface{}, error)) (interface{}, error)
	GetState() State
	GetStats() *BreakerStats
	Reset() error
}

// BreakerStats 熔断器统计信息
type BreakerStats struct {
	Name              string        `json:"name"`
	State             State         `json:"state"`
	TotalRequests     int64         `json:"total_requests"`
	SuccessfulRequests int64        `json:"successful_requests"`
	FailedRequests    int64         `json:"failed_requests"`
	ConsecutiveFailures int64       `json:"consecutive_failures"`
	LastFailureTime   time.Time     `json:"last_failure_time"`
	LastSuccessTime   time.Time     `json:"last_success_time"`
	NextRetryTime     time.Time     `json:"next_retry_time"`
	FailureRate       float64       `json:"failure_rate"`
	Timeout           time.Duration `json:"timeout"`
}

// Config 熔断器配置
type Config struct {
	Name                string        `json:"name"`
	MaxRequests         uint32        `json:"max_requests"`          // 半开状态下的最大请求数
	Interval            time.Duration `json:"interval"`              // 统计时间窗口
	Timeout             time.Duration `json:"timeout"`               // 熔断超时时间
	ReadyToTrip         func(counts Counts) bool `json:"-"`         // 触发熔断的条件
	OnStateChange       func(name string, from State, to State) `json:"-"` // 状态变化回调
}

// Counts 计数器
type Counts struct {
	Requests             uint32
	TotalSuccesses       uint32
	TotalFailures        uint32
	ConsecutiveSuccesses uint32
	ConsecutiveFailures  uint32
}

// DefaultCircuitBreaker 默认熔断器实现
type DefaultCircuitBreaker struct {
	name          string
	maxRequests   uint32
	interval      time.Duration
	timeout       time.Duration
	readyToTrip   func(counts Counts) bool
	onStateChange func(name string, from State, to State)

	mutex      sync.Mutex
	state      State
	generation uint64
	counts     Counts
	expiry     time.Time

	client *redisClient.Client
	ctx    context.Context
}

// BreakerManager 熔断器管理器
type BreakerManager struct {
	breakers map[string]CircuitBreaker
	mutex    sync.RWMutex
	client   *redisClient.Client
	ctx      context.Context
}

var (
	breakerManager *BreakerManager
	breakerOnce    sync.Once
)

// InitBreakerManager 初始化熔断器管理器
func InitBreakerManager() (*BreakerManager, error) {
	var err error
	breakerOnce.Do(func() {
		breakerManager, err = newBreakerManager()
	})
	return breakerManager, err
}

// GetBreakerManager 获取熔断器管理器实例
func GetBreakerManager() *BreakerManager {
	return breakerManager
}

// newBreakerManager 创建新的熔断器管理器
func newBreakerManager() (*BreakerManager, error) {
	if redis.RedisDb == nil {
		return nil, fmt.Errorf("Redis客户端未初始化")
	}

	manager := &BreakerManager{
		breakers: make(map[string]CircuitBreaker),
		client:   redis.RedisDb,
		ctx:      context.Background(),
	}

	log.Println("熔断器管理器初始化成功")
	return manager, nil
}

// NewCircuitBreaker 创建新的熔断器
func (bm *BreakerManager) NewCircuitBreaker(config Config) CircuitBreaker {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	if config.ReadyToTrip == nil {
		config.ReadyToTrip = defaultReadyToTrip
	}

	if config.OnStateChange == nil {
		config.OnStateChange = defaultOnStateChange
	}

	breaker := &DefaultCircuitBreaker{
		name:          config.Name,
		maxRequests:   config.MaxRequests,
		interval:      config.Interval,
		timeout:       config.Timeout,
		readyToTrip:   config.ReadyToTrip,
		onStateChange: config.OnStateChange,
		state:         StateClosed,
		client:        bm.client,
		ctx:           bm.ctx,
	}

	bm.breakers[config.Name] = breaker
	return breaker
}

// GetCircuitBreaker 获取熔断器
func (bm *BreakerManager) GetCircuitBreaker(name string) (CircuitBreaker, bool) {
	bm.mutex.RLock()
	defer bm.mutex.RUnlock()

	breaker, exists := bm.breakers[name]
	return breaker, exists
}

// GetAllBreakers 获取所有熔断器
func (bm *BreakerManager) GetAllBreakers() map[string]CircuitBreaker {
	bm.mutex.RLock()
	defer bm.mutex.RUnlock()

	result := make(map[string]CircuitBreaker)
	for name, breaker := range bm.breakers {
		result[name] = breaker
	}
	return result
}

// DefaultCircuitBreaker 实现

// Execute 执行函数（无返回值）
func (cb *DefaultCircuitBreaker) Execute(fn func() error) error {
	_, err := cb.Call(func() (interface{}, error) {
		return nil, fn()
	})
	return err
}

// Call 执行函数（有返回值）
func (cb *DefaultCircuitBreaker) Call(fn func() (interface{}, error)) (interface{}, error) {
	generation, err := cb.beforeRequest()
	if err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			cb.afterRequest(generation, false)
			panic(r)
		}
	}()

	result, err := fn()
	cb.afterRequest(generation, err == nil)
	return result, err
}

// beforeRequest 请求前检查
func (cb *DefaultCircuitBreaker) beforeRequest() (uint64, error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()
	state, generation := cb.currentState(now)

	if state == StateOpen {
		return generation, fmt.Errorf("熔断器 %s 处于开启状态", cb.name)
	} else if state == StateHalfOpen && cb.counts.Requests >= cb.maxRequests {
		return generation, fmt.Errorf("熔断器 %s 半开状态请求数已达上限", cb.name)
	}

	cb.counts.onRequest()
	return generation, nil
}

// afterRequest 请求后处理
func (cb *DefaultCircuitBreaker) afterRequest(before uint64, success bool) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()
	state, generation := cb.currentState(now)
	if generation != before {
		return
	}

	if success {
		cb.onSuccess(state, now)
	} else {
		cb.onFailure(state, now)
	}

	// 更新Redis中的统计信息
	go cb.updateRedisStats()
}

// onSuccess 成功处理
func (cb *DefaultCircuitBreaker) onSuccess(state State, now time.Time) {
	switch state {
	case StateClosed:
		cb.counts.onSuccess()
	case StateHalfOpen:
		cb.counts.onSuccess()
		if cb.counts.ConsecutiveSuccesses >= cb.maxRequests {
			cb.setState(StateClosed, now)
		}
	}
}

// onFailure 失败处理
func (cb *DefaultCircuitBreaker) onFailure(state State, now time.Time) {
	switch state {
	case StateClosed:
		cb.counts.onFailure()
		if cb.readyToTrip(cb.counts) {
			cb.setState(StateOpen, now)
		}
	case StateHalfOpen:
		cb.setState(StateOpen, now)
	}
}

// currentState 获取当前状态
func (cb *DefaultCircuitBreaker) currentState(now time.Time) (State, uint64) {
	switch cb.state {
	case StateClosed:
		if !cb.expiry.IsZero() && cb.expiry.Before(now) {
			cb.toNewGeneration(now)
		}
	case StateOpen:
		if cb.expiry.Before(now) {
			cb.setState(StateHalfOpen, now)
		}
	}
	return cb.state, cb.generation
}

// setState 设置状态
func (cb *DefaultCircuitBreaker) setState(state State, now time.Time) {
	if cb.state == state {
		return
	}

	prev := cb.state
	cb.state = state

	cb.toNewGeneration(now)

	if state == StateOpen {
		cb.expiry = now.Add(cb.timeout)
	} else {
		cb.expiry = time.Time{}
	}

	cb.onStateChange(cb.name, prev, state)
}

// toNewGeneration 进入新的周期
func (cb *DefaultCircuitBreaker) toNewGeneration(now time.Time) {
	cb.generation++
	cb.counts.clear()

	var zero time.Time
	switch cb.state {
	case StateClosed:
		if cb.interval == 0 {
			cb.expiry = zero
		} else {
			cb.expiry = now.Add(cb.interval)
		}
	case StateOpen:
		cb.expiry = now.Add(cb.timeout)
	default: // StateHalfOpen
		cb.expiry = zero
	}
}

// GetState 获取状态
func (cb *DefaultCircuitBreaker) GetState() State {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()
	state, _ := cb.currentState(now)
	return state
}

// GetStats 获取统计信息
func (cb *DefaultCircuitBreaker) GetStats() *BreakerStats {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()
	state, _ := cb.currentState(now)

	var failureRate float64
	if cb.counts.Requests > 0 {
		failureRate = float64(cb.counts.TotalFailures) / float64(cb.counts.Requests)
	}

	var nextRetryTime time.Time
	if state == StateOpen {
		nextRetryTime = cb.expiry
	}

	return &BreakerStats{
		Name:                cb.name,
		State:               state,
		TotalRequests:       int64(cb.counts.Requests),
		SuccessfulRequests:  int64(cb.counts.TotalSuccesses),
		FailedRequests:      int64(cb.counts.TotalFailures),
		ConsecutiveFailures: int64(cb.counts.ConsecutiveFailures),
		NextRetryTime:       nextRetryTime,
		FailureRate:         failureRate,
		Timeout:             cb.timeout,
	}
}

// Reset 重置熔断器
func (cb *DefaultCircuitBreaker) Reset() error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.toNewGeneration(time.Now())
	cb.setState(StateClosed, time.Now())

	// 清除Redis中的统计信息
	key := fmt.Sprintf("breaker:stats:%s", cb.name)
	return cb.client.Del(cb.ctx, key).Err()
}

// updateRedisStats 更新Redis中的统计信息
func (cb *DefaultCircuitBreaker) updateRedisStats() {
	key := fmt.Sprintf("breaker:stats:%s", cb.name)
	stats := cb.GetStats()

	pipe := cb.client.Pipeline()
	pipe.HSet(cb.ctx, key, map[string]interface{}{
		"state":                stats.State.String(),
		"total_requests":       stats.TotalRequests,
		"successful_requests":  stats.SuccessfulRequests,
		"failed_requests":      stats.FailedRequests,
		"consecutive_failures": stats.ConsecutiveFailures,
		"failure_rate":         stats.FailureRate,
		"last_update":          time.Now().Unix(),
	})
	pipe.Expire(cb.ctx, key, time.Hour)
	pipe.Exec(cb.ctx)
}

// Counts 方法

func (c *Counts) onRequest() {
	c.Requests++
}

func (c *Counts) onSuccess() {
	c.TotalSuccesses++
	c.ConsecutiveSuccesses++
	c.ConsecutiveFailures = 0
}

func (c *Counts) onFailure() {
	c.TotalFailures++
	c.ConsecutiveFailures++
	c.ConsecutiveSuccesses = 0
}

func (c *Counts) clear() {
	c.Requests = 0
	c.TotalSuccesses = 0
	c.TotalFailures = 0
	c.ConsecutiveSuccesses = 0
	c.ConsecutiveFailures = 0
}

// 默认配置函数

// defaultReadyToTrip 默认熔断触发条件
func defaultReadyToTrip(counts Counts) bool {
	return counts.Requests >= 20 && counts.TotalFailures >= 10
}

// defaultOnStateChange 默认状态变化回调
func defaultOnStateChange(name string, from State, to State) {
	log.Printf("熔断器 %s 状态变化: %s -> %s", name, from.String(), to.String())
}
