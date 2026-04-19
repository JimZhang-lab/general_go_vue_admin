/*
 * @Author: JimZhang
 * @Date: 2025-07-29 16:10:00
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-29 16:10:00
 * @FilePath: /server/pkg/redis/distributed_lock.go
 * @Description: Redis分布式锁
 */

package redis

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// DistributedLock 分布式锁
type DistributedLock struct {
	client    *redis.Client
	ctx       context.Context
	key       string
	value     string
	ttl       time.Duration
	acquired  bool
	renewStop chan struct{}
}

// LockManager 锁管理器
type LockManager struct {
	client *redis.Client
	ctx    context.Context
	prefix string
}

var lockManager *LockManager

// InitLockManager 初始化锁管理器
func InitLockManager() *LockManager {
	if RedisDb == nil {
		log.Fatal("Redis客户端未初始化")
	}

	lockManager = &LockManager{
		client: RedisDb,
		ctx:    context.Background(),
		prefix: "lock:",
	}

	log.Println("Redis分布式锁管理器初始化成功")
	return lockManager
}

// GetLockManager 获取锁管理器实例
func GetLockManager() *LockManager {
	return lockManager
}

// NewLock 创建新的分布式锁
func (lm *LockManager) NewLock(key string, ttl time.Duration) *DistributedLock {
	if ttl == 0 {
		ttl = 30 * time.Second // 默认30秒
	}

	return &DistributedLock{
		client:    lm.client,
		ctx:       lm.ctx,
		key:       lm.prefix + key,
		value:     generateLockValue(),
		ttl:       ttl,
		acquired:  false,
		renewStop: make(chan struct{}),
	}
}

// TryLock 尝试获取锁（非阻塞）
func (dl *DistributedLock) TryLock() (bool, error) {
	// 使用SET命令的NX和EX选项实现原子操作
	result, err := dl.client.SetNX(dl.ctx, dl.key, dl.value, dl.ttl).Result()
	if err != nil {
		return false, fmt.Errorf("尝试获取锁失败: %v", err)
	}

	if result {
		dl.acquired = true
		// 启动自动续期
		go dl.autoRenew()
		log.Printf("成功获取分布式锁: %s", dl.key)
	}

	return result, nil
}

// Lock 获取锁（阻塞，带超时）
func (dl *DistributedLock) Lock(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	
	for time.Now().Before(deadline) {
		acquired, err := dl.TryLock()
		if err != nil {
			return err
		}
		
		if acquired {
			return nil
		}
		
		// 等待一小段时间后重试
		time.Sleep(100 * time.Millisecond)
	}
	
	return fmt.Errorf("获取锁超时: %s", dl.key)
}

// Unlock 释放锁
func (dl *DistributedLock) Unlock() error {
	if !dl.acquired {
		return fmt.Errorf("锁未被获取")
	}

	// 停止自动续期
	close(dl.renewStop)

	// 使用Lua脚本确保只有锁的持有者才能释放锁
	luaScript := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`

	result, err := dl.client.Eval(dl.ctx, luaScript, []string{dl.key}, dl.value).Result()
	if err != nil {
		return fmt.Errorf("释放锁失败: %v", err)
	}

	if result.(int64) == 1 {
		dl.acquired = false
		log.Printf("成功释放分布式锁: %s", dl.key)
		return nil
	}

	return fmt.Errorf("锁已被其他进程释放或过期")
}

// IsLocked 检查锁是否被持有
func (dl *DistributedLock) IsLocked() (bool, error) {
	value, err := dl.client.Get(dl.ctx, dl.key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, fmt.Errorf("检查锁状态失败: %v", err)
	}

	// 检查是否是当前实例持有的锁
	return value == dl.value, nil
}

// Extend 延长锁的过期时间
func (dl *DistributedLock) Extend(ttl time.Duration) error {
	if !dl.acquired {
		return fmt.Errorf("锁未被获取")
	}

	// 使用Lua脚本确保只有锁的持有者才能延长锁
	luaScript := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("expire", KEYS[1], ARGV[2])
		else
			return 0
		end
	`

	result, err := dl.client.Eval(dl.ctx, luaScript, []string{dl.key}, dl.value, int(ttl.Seconds())).Result()
	if err != nil {
		return fmt.Errorf("延长锁失败: %v", err)
	}

	if result.(int64) == 1 {
		dl.ttl = ttl
		return nil
	}

	return fmt.Errorf("锁已被其他进程释放或过期")
}

// autoRenew 自动续期
func (dl *DistributedLock) autoRenew() {
	ticker := time.NewTicker(dl.ttl / 3) // 每1/3的TTL时间续期一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := dl.Extend(dl.ttl)
			if err != nil {
				log.Printf("自动续期失败: %v", err)
				return
			}
		case <-dl.renewStop:
			return
		}
	}
}

// WithLock 使用锁执行函数
func (lm *LockManager) WithLock(key string, ttl time.Duration, timeout time.Duration, fn func() error) error {
	lock := lm.NewLock(key, ttl)
	
	// 获取锁
	err := lock.Lock(timeout)
	if err != nil {
		return fmt.Errorf("获取锁失败: %v", err)
	}
	
	// 确保释放锁
	defer func() {
		if unlockErr := lock.Unlock(); unlockErr != nil {
			log.Printf("释放锁失败: %v", unlockErr)
		}
	}()
	
	// 执行函数
	return fn()
}

// TryWithLock 尝试使用锁执行函数（非阻塞）
func (lm *LockManager) TryWithLock(key string, ttl time.Duration, fn func() error) error {
	lock := lm.NewLock(key, ttl)
	
	// 尝试获取锁
	acquired, err := lock.TryLock()
	if err != nil {
		return fmt.Errorf("尝试获取锁失败: %v", err)
	}
	
	if !acquired {
		return fmt.Errorf("锁被其他进程持有")
	}
	
	// 确保释放锁
	defer func() {
		if unlockErr := lock.Unlock(); unlockErr != nil {
			log.Printf("释放锁失败: %v", unlockErr)
		}
	}()
	
	// 执行函数
	return fn()
}

// GetActiveLocks 获取活跃的锁
func (lm *LockManager) GetActiveLocks() ([]string, error) {
	pattern := lm.prefix + "*"
	keys, err := lm.client.Keys(lm.ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("获取活跃锁失败: %v", err)
	}

	// 移除前缀
	locks := make([]string, len(keys))
	for i, key := range keys {
		locks[i] = key[len(lm.prefix):]
	}

	return locks, nil
}

// ForceUnlock 强制释放锁（谨慎使用）
func (lm *LockManager) ForceUnlock(key string) error {
	fullKey := lm.prefix + key
	err := lm.client.Del(lm.ctx, fullKey).Err()
	if err != nil {
		return fmt.Errorf("强制释放锁失败: %v", err)
	}

	log.Printf("强制释放锁: %s", fullKey)
	return nil
}

// generateLockValue 生成锁的唯一值
func generateLockValue() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GetLockStats 获取锁统计信息
func (lm *LockManager) GetLockStats() (map[string]interface{}, error) {
	activeLocks, err := lm.GetActiveLocks()
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"active_locks_count": len(activeLocks),
		"active_locks":       activeLocks,
		"prefix":            lm.prefix,
	}

	return stats, nil
}
