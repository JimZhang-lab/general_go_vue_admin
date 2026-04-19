/*
 * @Author: JimZhang
 * @Date: 2025-07-29 15:40:00
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-29 15:40:00
 * @FilePath: /server/pkg/rabbitmq/async_logger.go
 * @Description: 异步日志处理器
 */

package rabbitmq

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// AsyncLogger 异步日志处理器
type AsyncLogger struct {
	queueManager *QueueManager
}

// LogData 日志数据结构
type LogData struct {
	Type      string                 `json:"type"`
	AdminID   int                    `json:"admin_id"`
	Username  string                 `json:"username"`
	Method    string                 `json:"method"`
	IP        string                 `json:"ip"`
	URL       string                 `json:"url"`
	Message   string                 `json:"message"`
	Status    int                    `json:"status"`
	Browser   string                 `json:"browser"`
	OS        string                 `json:"os"`
	Location  string                 `json:"location"`
	Extra     map[string]interface{} `json:"extra"`
	Timestamp time.Time              `json:"timestamp"`
}

var asyncLogger *AsyncLogger

// InitAsyncLogger 初始化异步日志处理器
func InitAsyncLogger() (*AsyncLogger, error) {
	qm := GetQueueManager()
	if qm == nil {
		return nil, fmt.Errorf("队列管理器未初始化")
	}

	asyncLogger = &AsyncLogger{
		queueManager: qm,
	}

	// 启动日志消费者
	err := asyncLogger.startConsumers()
	if err != nil {
		return nil, fmt.Errorf("启动日志消费者失败: %v", err)
	}

	log.Println("异步日志处理器初始化成功")
	return asyncLogger, nil
}

// GetAsyncLogger 获取异步日志处理器实例
func GetAsyncLogger() *AsyncLogger {
	return asyncLogger
}

// startConsumers 启动消费者
func (al *AsyncLogger) startConsumers() error {
	// 启动操作日志消费者
	err := al.queueManager.ConsumeMessages("operation_logs", al.handleOperationLog)
	if err != nil {
		return fmt.Errorf("启动操作日志消费者失败: %v", err)
	}

	// 启动登录日志消费者
	err = al.queueManager.ConsumeMessages("login_logs", al.handleLoginLog)
	if err != nil {
		return fmt.Errorf("启动登录日志消费者失败: %v", err)
	}

	return nil
}

// LogOperation 记录操作日志
func (al *AsyncLogger) LogOperation(data LogData) error {
	message := &Message{
		ID:   uuid.New().String(),
		Type: "operation_log",
		Data: map[string]interface{}{
			"admin_id":  data.AdminID,
			"username":  data.Username,
			"method":    data.Method,
			"ip":        data.IP,
			"url":       data.URL,
			"extra":     data.Extra,
			"timestamp": data.Timestamp,
		},
		Timestamp: time.Now(),
		Retry:     0,
		MaxRetry:  3,
	}

	return al.queueManager.PublishMessage("logs.direct", "operation_logs", message)
}

// LogLogin 记录登录日志
func (al *AsyncLogger) LogLogin(data LogData) error {
	message := &Message{
		ID:   uuid.New().String(),
		Type: "login_log",
		Data: map[string]interface{}{
			"username":  data.Username,
			"ip":        data.IP,
			"location":  data.Location,
			"browser":   data.Browser,
			"os":        data.OS,
			"status":    data.Status,
			"message":   data.Message,
			"timestamp": data.Timestamp,
		},
		Timestamp: time.Now(),
		Retry:     0,
		MaxRetry:  3,
	}

	return al.queueManager.PublishMessage("logs.direct", "login_logs", message)
}

// handleOperationLog 处理操作日志
func (al *AsyncLogger) handleOperationLog(message *Message) error {
	data := message.Data

	// 提取数据
	username, _ := data["username"].(string)
	method, _ := data["method"].(string)
	ip, _ := data["ip"].(string)
	url, _ := data["url"].(string)

	// 这里应该保存到数据库，暂时只记录日志
	log.Printf("异步处理操作日志: 用户=%s, 方法=%s, URL=%s, IP=%s", username, method, url, ip)

	// TODO: 集成实际的数据库操作
	// 可以在这里调用数据库服务保存日志

	return nil
}

// handleLoginLog 处理登录日志
func (al *AsyncLogger) handleLoginLog(message *Message) error {
	data := message.Data

	// 提取数据
	username, _ := data["username"].(string)
	ip, _ := data["ip"].(string)
	location, _ := data["location"].(string)
	browser, _ := data["browser"].(string)
	os, _ := data["os"].(string)
	status, _ := data["status"].(float64)
	logMessage, _ := data["message"].(string)

	// 这里应该保存到数据库，暂时只记录日志
	log.Printf("异步处理登录日志: 用户=%s, 状态=%d, IP=%s, 位置=%s, 浏览器=%s, 系统=%s, 消息=%s",
		username, int(status), ip, location, browser, os, logMessage)

	// TODO: 集成实际的数据库操作
	// 可以在这里调用数据库服务保存登录日志

	return nil
}

// GetStats 获取异步日志统计信息
func (al *AsyncLogger) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 获取操作日志队列信息
	operationQueue, err := al.queueManager.GetQueueInfo("operation_logs")
	if err != nil {
		log.Printf("获取操作日志队列信息失败: %v", err)
	} else {
		stats["operation_logs"] = map[string]interface{}{
			"messages":  operationQueue.Messages,
			"consumers": operationQueue.Consumers,
		}
	}

	// 获取登录日志队列信息
	loginQueue, err := al.queueManager.GetQueueInfo("login_logs")
	if err != nil {
		log.Printf("获取登录日志队列信息失败: %v", err)
	} else {
		stats["login_logs"] = map[string]interface{}{
			"messages":  loginQueue.Messages,
			"consumers": loginQueue.Consumers,
		}
	}

	return stats, nil
}

// Flush 刷新所有待处理的日志
func (al *AsyncLogger) Flush() error {
	// 这里可以实现强制处理所有队列中的消息
	// 目前只是返回成功，实际实现可能需要等待队列清空
	log.Println("异步日志刷新完成")
	return nil
}

// Close 关闭异步日志处理器
func (al *AsyncLogger) Close() error {
	// 刷新待处理的日志
	err := al.Flush()
	if err != nil {
		log.Printf("刷新异步日志失败: %v", err)
	}

	log.Println("异步日志处理器已关闭")
	return nil
}
