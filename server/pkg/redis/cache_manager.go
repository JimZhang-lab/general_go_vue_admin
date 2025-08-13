/*
 * @Author: JimZhang
 * @Date: 2025-07-29 16:00:00
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-29 16:00:00
 * @FilePath: /server/pkg/redis/cache_manager.go
 * @Description: Redis缓存管理器
 */

package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheManager Redis缓存管理器
type CacheManager struct {
	client *redis.Client
	ctx    context.Context
}

// CacheConfig 缓存配置
type CacheConfig struct {
	DefaultTTL time.Duration
	MaxRetries int
	Prefix     string
}

var (
	cacheManager *CacheManager
	defaultConfig = &CacheConfig{
		DefaultTTL: time.Hour,
		MaxRetries: 3,
		Prefix:     "app:",
	}
)

// InitCacheManager 初始化缓存管理器
func InitCacheManager() *CacheManager {
	if RedisDb == nil {
		log.Fatal("Redis客户端未初始化")
	}

	cacheManager = &CacheManager{
		client: RedisDb,
		ctx:    context.Background(),
	}

	log.Println("Redis缓存管理器初始化成功")
	return cacheManager
}

// GetCacheManager 获取缓存管理器实例
func GetCacheManager() *CacheManager {
	return cacheManager
}

// Set 设置缓存
func (cm *CacheManager) Set(key string, value interface{}, ttl time.Duration) error {
	if ttl == 0 {
		ttl = defaultConfig.DefaultTTL
	}

	// 序列化值
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化缓存值失败: %v", err)
	}

	// 添加前缀
	fullKey := defaultConfig.Prefix + key

	// 设置缓存
	err = cm.client.Set(cm.ctx, fullKey, data, ttl).Err()
	if err != nil {
		return fmt.Errorf("设置缓存失败: %v", err)
	}

	return nil
}

// Get 获取缓存
func (cm *CacheManager) Get(key string, dest interface{}) error {
	fullKey := defaultConfig.Prefix + key

	// 获取缓存
	data, err := cm.client.Get(cm.ctx, fullKey).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("缓存不存在")
		}
		return fmt.Errorf("获取缓存失败: %v", err)
	}

	// 反序列化
	err = json.Unmarshal([]byte(data), dest)
	if err != nil {
		return fmt.Errorf("反序列化缓存值失败: %v", err)
	}

	return nil
}

// GetOrSet 获取缓存，如果不存在则设置
func (cm *CacheManager) GetOrSet(key string, dest interface{}, setter func() (interface{}, error), ttl time.Duration) error {
	// 尝试获取缓存
	err := cm.Get(key, dest)
	if err == nil {
		return nil // 缓存存在，直接返回
	}

	// 缓存不存在，调用setter获取数据
	value, err := setter()
	if err != nil {
		return fmt.Errorf("setter函数执行失败: %v", err)
	}

	// 设置缓存
	err = cm.Set(key, value, ttl)
	if err != nil {
		return fmt.Errorf("设置缓存失败: %v", err)
	}

	// 将值复制到dest
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化值失败: %v", err)
	}

	err = json.Unmarshal(data, dest)
	if err != nil {
		return fmt.Errorf("反序列化值失败: %v", err)
	}

	return nil
}

// Delete 删除缓存
func (cm *CacheManager) Delete(key string) error {
	fullKey := defaultConfig.Prefix + key
	return cm.client.Del(cm.ctx, fullKey).Err()
}

// DeletePattern 按模式删除缓存
func (cm *CacheManager) DeletePattern(pattern string) error {
	fullPattern := defaultConfig.Prefix + pattern
	
	// 获取匹配的键
	keys, err := cm.client.Keys(cm.ctx, fullPattern).Result()
	if err != nil {
		return fmt.Errorf("获取匹配键失败: %v", err)
	}

	if len(keys) == 0 {
		return nil
	}

	// 批量删除
	return cm.client.Del(cm.ctx, keys...).Err()
}

// Exists 检查缓存是否存在
func (cm *CacheManager) Exists(key string) (bool, error) {
	fullKey := defaultConfig.Prefix + key
	count, err := cm.client.Exists(cm.ctx, fullKey).Result()
	return count > 0, err
}

// TTL 获取缓存剩余时间
func (cm *CacheManager) TTL(key string) (time.Duration, error) {
	fullKey := defaultConfig.Prefix + key
	return cm.client.TTL(cm.ctx, fullKey).Result()
}

// Expire 设置缓存过期时间
func (cm *CacheManager) Expire(key string, ttl time.Duration) error {
	fullKey := defaultConfig.Prefix + key
	return cm.client.Expire(cm.ctx, fullKey, ttl).Err()
}

// Increment 递增计数器
func (cm *CacheManager) Increment(key string, delta int64) (int64, error) {
	fullKey := defaultConfig.Prefix + key
	return cm.client.IncrBy(cm.ctx, fullKey, delta).Result()
}

// Decrement 递减计数器
func (cm *CacheManager) Decrement(key string, delta int64) (int64, error) {
	fullKey := defaultConfig.Prefix + key
	return cm.client.DecrBy(cm.ctx, fullKey, delta).Result()
}

// SetNX 仅当键不存在时设置
func (cm *CacheManager) SetNX(key string, value interface{}, ttl time.Duration) (bool, error) {
	if ttl == 0 {
		ttl = defaultConfig.DefaultTTL
	}

	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("序列化值失败: %v", err)
	}

	fullKey := defaultConfig.Prefix + key
	return cm.client.SetNX(cm.ctx, fullKey, data, ttl).Result()
}

// HSet 设置哈希字段
func (cm *CacheManager) HSet(key, field string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化值失败: %v", err)
	}

	fullKey := defaultConfig.Prefix + key
	return cm.client.HSet(cm.ctx, fullKey, field, data).Err()
}

// HGet 获取哈希字段
func (cm *CacheManager) HGet(key, field string, dest interface{}) error {
	fullKey := defaultConfig.Prefix + key
	
	data, err := cm.client.HGet(cm.ctx, fullKey, field).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("哈希字段不存在")
		}
		return fmt.Errorf("获取哈希字段失败: %v", err)
	}

	return json.Unmarshal([]byte(data), dest)
}

// HGetAll 获取所有哈希字段
func (cm *CacheManager) HGetAll(key string) (map[string]string, error) {
	fullKey := defaultConfig.Prefix + key
	return cm.client.HGetAll(cm.ctx, fullKey).Result()
}

// HDel 删除哈希字段
func (cm *CacheManager) HDel(key string, fields ...string) error {
	fullKey := defaultConfig.Prefix + key
	return cm.client.HDel(cm.ctx, fullKey, fields...).Err()
}

// LPush 从左侧推入列表
func (cm *CacheManager) LPush(key string, values ...interface{}) error {
	fullKey := defaultConfig.Prefix + key
	
	// 序列化所有值
	serializedValues := make([]interface{}, len(values))
	for i, value := range values {
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("序列化值失败: %v", err)
		}
		serializedValues[i] = data
	}

	return cm.client.LPush(cm.ctx, fullKey, serializedValues...).Err()
}

// RPop 从右侧弹出列表元素
func (cm *CacheManager) RPop(key string, dest interface{}) error {
	fullKey := defaultConfig.Prefix + key
	
	data, err := cm.client.RPop(cm.ctx, fullKey).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("列表为空")
		}
		return fmt.Errorf("弹出列表元素失败: %v", err)
	}

	return json.Unmarshal([]byte(data), dest)
}

// LLen 获取列表长度
func (cm *CacheManager) LLen(key string) (int64, error) {
	fullKey := defaultConfig.Prefix + key
	return cm.client.LLen(cm.ctx, fullKey).Result()
}

// SAdd 添加集合成员
func (cm *CacheManager) SAdd(key string, members ...interface{}) error {
	fullKey := defaultConfig.Prefix + key
	
	// 序列化所有成员
	serializedMembers := make([]interface{}, len(members))
	for i, member := range members {
		data, err := json.Marshal(member)
		if err != nil {
			return fmt.Errorf("序列化成员失败: %v", err)
		}
		serializedMembers[i] = data
	}

	return cm.client.SAdd(cm.ctx, fullKey, serializedMembers...).Err()
}

// SMembers 获取集合所有成员
func (cm *CacheManager) SMembers(key string) ([]string, error) {
	fullKey := defaultConfig.Prefix + key
	return cm.client.SMembers(cm.ctx, fullKey).Result()
}

// SIsMember 检查是否为集合成员
func (cm *CacheManager) SIsMember(key string, member interface{}) (bool, error) {
	data, err := json.Marshal(member)
	if err != nil {
		return false, fmt.Errorf("序列化成员失败: %v", err)
	}

	fullKey := defaultConfig.Prefix + key
	return cm.client.SIsMember(cm.ctx, fullKey, data).Result()
}

// GetStats 获取缓存统计信息
func (cm *CacheManager) GetStats() (map[string]interface{}, error) {
	info, err := cm.client.Info(cm.ctx, "memory", "stats").Result()
	if err != nil {
		return nil, fmt.Errorf("获取Redis信息失败: %v", err)
	}

	// 获取数据库大小
	dbSize, err := cm.client.DBSize(cm.ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("获取数据库大小失败: %v", err)
	}

	return map[string]interface{}{
		"db_size": dbSize,
		"info":    info,
	}, nil
}

// FlushDB 清空当前数据库
func (cm *CacheManager) FlushDB() error {
	return cm.client.FlushDB(cm.ctx).Err()
}
