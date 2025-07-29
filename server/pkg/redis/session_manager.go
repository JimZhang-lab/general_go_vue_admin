/*
 * @Author: JimZhang
 * @Date: 2025-07-29 16:15:00
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-29 16:15:00
 * @FilePath: /server/pkg/redis/session_manager.go
 * @Description: Redis会话管理器
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

// SessionManager 会话管理器
type SessionManager struct {
	client *redis.Client
	ctx    context.Context
	prefix string
	ttl    time.Duration
}

// SessionData 会话数据
type SessionData struct {
	UserID    int                    `json:"user_id"`
	Username  string                 `json:"username"`
	LoginTime time.Time              `json:"login_time"`
	LastSeen  time.Time              `json:"last_seen"`
	IP        string                 `json:"ip"`
	UserAgent string                 `json:"user_agent"`
	Extra     map[string]interface{} `json:"extra"`
}

var sessionManager *SessionManager

// InitSessionManager 初始化会话管理器
func InitSessionManager() *SessionManager {
	if RedisDb == nil {
		log.Fatal("Redis客户端未初始化")
	}

	sessionManager = &SessionManager{
		client: RedisDb,
		ctx:    context.Background(),
		prefix: "session:",
		ttl:    24 * time.Hour, // 默认24小时
	}

	log.Println("Redis会话管理器初始化成功")
	return sessionManager
}

// GetSessionManager 获取会话管理器实例
func GetSessionManager() *SessionManager {
	return sessionManager
}

// CreateSession 创建会话
func (sm *SessionManager) CreateSession(sessionID string, data *SessionData) error {
	data.LoginTime = time.Now()
	data.LastSeen = time.Now()

	// 序列化会话数据
	sessionData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("序列化会话数据失败: %v", err)
	}

	// 存储会话
	fullKey := sm.prefix + sessionID
	err = sm.client.Set(sm.ctx, fullKey, sessionData, sm.ttl).Err()
	if err != nil {
		return fmt.Errorf("创建会话失败: %v", err)
	}

	// 添加到用户会话集合（支持多设备登录）
	userSessionKey := fmt.Sprintf("user_sessions:%d", data.UserID)
	err = sm.client.SAdd(sm.ctx, userSessionKey, sessionID).Err()
	if err != nil {
		log.Printf("添加到用户会话集合失败: %v", err)
	}

	// 设置用户会话集合过期时间
	sm.client.Expire(sm.ctx, userSessionKey, sm.ttl)

	log.Printf("创建会话成功: 用户=%s, 会话ID=%s", data.Username, sessionID)
	return nil
}

// GetSession 获取会话
func (sm *SessionManager) GetSession(sessionID string) (*SessionData, error) {
	fullKey := sm.prefix + sessionID

	// 获取会话数据
	sessionData, err := sm.client.Get(sm.ctx, fullKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("会话不存在或已过期")
		}
		return nil, fmt.Errorf("获取会话失败: %v", err)
	}

	// 反序列化会话数据
	var data SessionData
	err = json.Unmarshal([]byte(sessionData), &data)
	if err != nil {
		return nil, fmt.Errorf("反序列化会话数据失败: %v", err)
	}

	return &data, nil
}

// UpdateSession 更新会话
func (sm *SessionManager) UpdateSession(sessionID string, data *SessionData) error {
	data.LastSeen = time.Now()

	// 序列化会话数据
	sessionData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("序列化会话数据失败: %v", err)
	}

	// 更新会话
	fullKey := sm.prefix + sessionID
	err = sm.client.Set(sm.ctx, fullKey, sessionData, sm.ttl).Err()
	if err != nil {
		return fmt.Errorf("更新会话失败: %v", err)
	}

	return nil
}

// RefreshSession 刷新会话过期时间
func (sm *SessionManager) RefreshSession(sessionID string) error {
	fullKey := sm.prefix + sessionID

	// 检查会话是否存在
	exists, err := sm.client.Exists(sm.ctx, fullKey).Result()
	if err != nil {
		return fmt.Errorf("检查会话存在性失败: %v", err)
	}

	if exists == 0 {
		return fmt.Errorf("会话不存在")
	}

	// 刷新过期时间
	err = sm.client.Expire(sm.ctx, fullKey, sm.ttl).Err()
	if err != nil {
		return fmt.Errorf("刷新会话过期时间失败: %v", err)
	}

	// 更新最后访问时间
	session, err := sm.GetSession(sessionID)
	if err != nil {
		return err
	}

	session.LastSeen = time.Now()
	return sm.UpdateSession(sessionID, session)
}

// DeleteSession 删除会话
func (sm *SessionManager) DeleteSession(sessionID string) error {
	// 获取会话数据以便从用户会话集合中移除
	session, err := sm.GetSession(sessionID)
	if err == nil {
		userSessionKey := fmt.Sprintf("user_sessions:%d", session.UserID)
		sm.client.SRem(sm.ctx, userSessionKey, sessionID)
	}

	// 删除会话
	fullKey := sm.prefix + sessionID
	err = sm.client.Del(sm.ctx, fullKey).Err()
	if err != nil {
		return fmt.Errorf("删除会话失败: %v", err)
	}

	log.Printf("删除会话成功: 会话ID=%s", sessionID)
	return nil
}

// GetUserSessions 获取用户的所有会话
func (sm *SessionManager) GetUserSessions(userID int) ([]string, error) {
	userSessionKey := fmt.Sprintf("user_sessions:%d", userID)
	
	sessions, err := sm.client.SMembers(sm.ctx, userSessionKey).Result()
	if err != nil {
		return nil, fmt.Errorf("获取用户会话失败: %v", err)
	}

	// 过滤掉已过期的会话
	validSessions := make([]string, 0)
	for _, sessionID := range sessions {
		fullKey := sm.prefix + sessionID
		exists, err := sm.client.Exists(sm.ctx, fullKey).Result()
		if err != nil {
			continue
		}
		
		if exists > 0 {
			validSessions = append(validSessions, sessionID)
		} else {
			// 从集合中移除过期的会话ID
			sm.client.SRem(sm.ctx, userSessionKey, sessionID)
		}
	}

	return validSessions, nil
}

// DeleteUserSessions 删除用户的所有会话
func (sm *SessionManager) DeleteUserSessions(userID int) error {
	sessions, err := sm.GetUserSessions(userID)
	if err != nil {
		return fmt.Errorf("获取用户会话失败: %v", err)
	}

	// 删除所有会话
	for _, sessionID := range sessions {
		err := sm.DeleteSession(sessionID)
		if err != nil {
			log.Printf("删除会话失败: %v", err)
		}
	}

	// 删除用户会话集合
	userSessionKey := fmt.Sprintf("user_sessions:%d", userID)
	sm.client.Del(sm.ctx, userSessionKey)

	log.Printf("删除用户所有会话成功: 用户ID=%d, 会话数=%d", userID, len(sessions))
	return nil
}

// IsSessionValid 检查会话是否有效
func (sm *SessionManager) IsSessionValid(sessionID string) (bool, error) {
	fullKey := sm.prefix + sessionID
	
	exists, err := sm.client.Exists(sm.ctx, fullKey).Result()
	if err != nil {
		return false, fmt.Errorf("检查会话有效性失败: %v", err)
	}

	return exists > 0, nil
}

// GetActiveSessions 获取所有活跃会话
func (sm *SessionManager) GetActiveSessions() ([]string, error) {
	pattern := sm.prefix + "*"
	keys, err := sm.client.Keys(sm.ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("获取活跃会话失败: %v", err)
	}

	// 移除前缀
	sessions := make([]string, len(keys))
	for i, key := range keys {
		sessions[i] = key[len(sm.prefix):]
	}

	return sessions, nil
}

// GetSessionStats 获取会话统计信息
func (sm *SessionManager) GetSessionStats() (map[string]interface{}, error) {
	activeSessions, err := sm.GetActiveSessions()
	if err != nil {
		return nil, err
	}

	// 统计在线用户数
	userSet := make(map[int]bool)
	for _, sessionID := range activeSessions {
		session, err := sm.GetSession(sessionID)
		if err != nil {
			continue
		}
		userSet[session.UserID] = true
	}

	stats := map[string]interface{}{
		"active_sessions_count": len(activeSessions),
		"online_users_count":    len(userSet),
		"session_ttl_seconds":   int(sm.ttl.Seconds()),
		"prefix":               sm.prefix,
	}

	return stats, nil
}

// CleanExpiredSessions 清理过期会话（定期任务）
func (sm *SessionManager) CleanExpiredSessions() error {
	// 获取所有用户会话集合
	pattern := "user_sessions:*"
	keys, err := sm.client.Keys(sm.ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("获取用户会话集合失败: %v", err)
	}

	cleanedCount := 0
	for _, userSessionKey := range keys {
		sessions, err := sm.client.SMembers(sm.ctx, userSessionKey).Result()
		if err != nil {
			continue
		}

		for _, sessionID := range sessions {
			fullKey := sm.prefix + sessionID
			exists, err := sm.client.Exists(sm.ctx, fullKey).Result()
			if err != nil {
				continue
			}

			if exists == 0 {
				// 会话已过期，从集合中移除
				sm.client.SRem(sm.ctx, userSessionKey, sessionID)
				cleanedCount++
			}
		}

		// 如果集合为空，删除集合
		count, err := sm.client.SCard(sm.ctx, userSessionKey).Result()
		if err == nil && count == 0 {
			sm.client.Del(sm.ctx, userSessionKey)
		}
	}

	if cleanedCount > 0 {
		log.Printf("清理过期会话完成: 清理数量=%d", cleanedCount)
	}

	return nil
}

// SetTTL 设置会话TTL
func (sm *SessionManager) SetTTL(ttl time.Duration) {
	sm.ttl = ttl
	log.Printf("会话TTL已更新: %v", ttl)
}
