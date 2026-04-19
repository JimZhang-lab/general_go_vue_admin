/*
 * @Author: JimZhang
 * @Date: 2025-07-27
 * @Description: 优化的服务基类，支持错误链式传递和高并发处理
 */
package service

import (
	"context"
	"fmt"
	"server/common/errors"
	"server/pkg/redis"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BaseService 基础服务类
type BaseService struct {
	validator *validator.Validate
	mu        sync.RWMutex
}

// NewBaseService 创建基础服务实例
func NewBaseService() *BaseService {
	return &BaseService{
		validator: validator.New(),
	}
}

// ValidateStruct 验证结构体
func (bs *BaseService) ValidateStruct(obj interface{}) *errors.AppError {
	if err := bs.validator.Struct(obj); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errorMessages []string
			for _, e := range validationErrors {
				errorMessages = append(errorMessages, bs.formatValidationError(e))
			}
			return errors.ValidationError(fmt.Sprintf("参数验证失败: %v", errorMessages))
		}
		return errors.Wrap(err, errors.ErrValidation, "参数验证失败")
	}
	return nil
}

// formatValidationError 格式化验证错误
func (bs *BaseService) formatValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("字段 %s 是必需的", err.Field())
	case "min":
		return fmt.Sprintf("字段 %s 的最小值为 %s", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("字段 %s 的最大值为 %s", err.Field(), err.Param())
	case "email":
		return fmt.Sprintf("字段 %s 必须是有效的邮箱地址", err.Field())
	case "len":
		return fmt.Sprintf("字段 %s 的长度必须为 %s", err.Field(), err.Param())
	default:
		return fmt.Sprintf("字段 %s 验证失败: %s", err.Field(), err.Tag())
	}
}

// GetTraceID 从上下文获取追踪ID
func (bs *BaseService) GetTraceID(c *gin.Context) string {
	if traceID := c.GetString("trace_id"); traceID != "" {
		return traceID
	}
	return c.GetHeader("X-Trace-ID")
}

// WithTimeout 为操作添加超时控制
func (bs *BaseService) WithTimeout(ctx context.Context, timeout time.Duration, fn func(ctx context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- fn(ctx)
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return errors.TimeoutError("操作超时")
	}
}

// CacheGet 从缓存获取数据
func (bs *BaseService) CacheGet(key string) (string, error) {
	if redis.RedisDb == nil {
		return "", errors.RedisError(fmt.Errorf("Redis客户端未初始化"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	result, err := redis.RedisDb.Get(ctx, key).Result()
	if err != nil {
		return "", errors.RedisError(err)
	}
	return result, nil
}

// CacheSet 设置缓存数据
func (bs *BaseService) CacheSet(key string, value interface{}, expiration time.Duration) error {
	if redis.RedisDb == nil {
		return errors.RedisError(fmt.Errorf("Redis客户端未初始化"))
	}

	ctx := context.Background()
	err := redis.RedisDb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return errors.RedisError(err)
	}
	return nil
}

// CacheDel 删除缓存数据
func (bs *BaseService) CacheDel(keys ...string) error {
	if redis.RedisDb == nil {
		return errors.RedisError(fmt.Errorf("Redis客户端未初始化"))
	}

	ctx := context.Background()
	err := redis.RedisDb.Del(ctx, keys...).Err()
	if err != nil {
		return errors.RedisError(err)
	}
	return nil
}

// CacheExists 检查缓存是否存在
func (bs *BaseService) CacheExists(key string) (bool, error) {
	if redis.RedisDb == nil {
		return false, errors.RedisError(fmt.Errorf("Redis客户端未初始化"))
	}

	ctx := context.Background()
	result, err := redis.RedisDb.Exists(ctx, key).Result()
	if err != nil {
		return false, errors.RedisError(err)
	}
	return result > 0, nil
}

// ParallelExecute 并行执行多个任务
func (bs *BaseService) ParallelExecute(tasks ...func() error) *errors.ErrorChain {
	errorChain := errors.NewErrorChain()
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, task := range tasks {
		wg.Add(1)
		go func(t func() error) {
			defer wg.Done()
			if err := t(); err != nil {
				mu.Lock()
				if appErr, ok := err.(*errors.AppError); ok {
					errorChain.Add(appErr)
				} else {
					errorChain.Add(errors.Wrap(err, errors.ErrInternal, "并行任务执行失败"))
				}
				mu.Unlock()
			}
		}(task)
	}

	wg.Wait()
	return errorChain
}

// RetryWithBackoff 带退避的重试机制
func (bs *BaseService) RetryWithBackoff(maxRetries int, initialDelay time.Duration, fn func() error) error {
	var lastErr error
	delay := initialDelay

	for i := 0; i < maxRetries; i++ {
		if err := fn(); err != nil {
			lastErr = err
			if i < maxRetries-1 {
				time.Sleep(delay)
				delay *= 2 // 指数退避
			}
		} else {
			return nil
		}
	}

	return errors.Wrapf(lastErr, errors.ErrInternal, "重试 %d 次后仍然失败", maxRetries)
}

// BatchProcess 批量处理数据
func (bs *BaseService) BatchProcess(items []interface{}, batchSize int, processor func(batch []interface{}) error) *errors.ErrorChain {
	errorChain := errors.NewErrorChain()

	for i := 0; i < len(items); i += batchSize {
		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}

		batch := items[i:end]
		if err := processor(batch); err != nil {
			if appErr, ok := err.(*errors.AppError); ok {
				errorChain.Add(appErr.WithContext("batch_index", i))
			} else {
				errorChain.Add(errors.Wrap(err, errors.ErrInternal, fmt.Sprintf("批次 %d 处理失败", i/batchSize+1)))
			}
		}
	}

	return errorChain
}

// ServiceResult 服务层统一返回结构
type ServiceResult struct {
	Data    interface{}      `json:"data"`
	Error   *errors.AppError `json:"error,omitempty"`
	TraceID string           `json:"trace_id,omitempty"`
}

// NewServiceResult 创建服务结果
func NewServiceResult(data interface{}, err *errors.AppError) *ServiceResult {
	return &ServiceResult{
		Data:  data,
		Error: err,
	}
}

// IsSuccess 检查是否成功
func (sr *ServiceResult) IsSuccess() bool {
	return sr.Error == nil
}

// WithTraceID 添加追踪ID
func (sr *ServiceResult) WithTraceID(traceID string) *ServiceResult {
	sr.TraceID = traceID
	if sr.Error != nil {
		sr.Error.WithTraceID(traceID)
	}
	return sr
}

// PageResult 分页结果
type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	Pages    int         `json:"pages"`
}

// NewPageResult 创建分页结果
func NewPageResult(list interface{}, total int64, page, pageSize int) *PageResult {
	pages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		pages++
	}

	return &PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Pages:    pages,
	}
}

// ConcurrentMap 并发安全的Map
type ConcurrentMap struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

// NewConcurrentMap 创建并发安全的Map
func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{
		data: make(map[string]interface{}),
	}
}

// Set 设置值
func (cm *ConcurrentMap) Set(key string, value interface{}) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.data[key] = value
}

// Get 获取值
func (cm *ConcurrentMap) Get(key string) (interface{}, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	value, exists := cm.data[key]
	return value, exists
}

// Delete 删除值
func (cm *ConcurrentMap) Delete(key string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.data, key)
}

// Keys 获取所有键
func (cm *ConcurrentMap) Keys() []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	keys := make([]string, 0, len(cm.data))
	for key := range cm.data {
		keys = append(keys, key)
	}
	return keys
}
