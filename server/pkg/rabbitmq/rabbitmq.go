/*
 * @Author: JimZhang
 * @Date: 2025-07-29 15:45:00
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-29 15:45:00
 * @FilePath: /server/pkg/rabbitmq/rabbitmq.go
 * @Description: RabbitMQ初始化和管理
 */

package rabbitmq

import (
	"fmt"
	"log"
)

// RabbitMQ 管理器
type RabbitMQ struct {
	connectionPool *ConnectionPool
	queueManager   *QueueManager
	asyncLogger    *AsyncLogger
}

var rabbitMQ *RabbitMQ

// Init 初始化RabbitMQ
func Init() error {
	var err error

	// 初始化连接池
	connectionPool, err := InitConnectionPool()
	if err != nil {
		return fmt.Errorf("初始化RabbitMQ连接池失败: %v", err)
	}

	// 初始化队列管理器
	queueManager, err := InitQueueManager()
	if err != nil {
		return fmt.Errorf("初始化队列管理器失败: %v", err)
	}

	// 初始化异步日志处理器
	asyncLogger, err := InitAsyncLogger()
	if err != nil {
		return fmt.Errorf("初始化异步日志处理器失败: %v", err)
	}

	rabbitMQ = &RabbitMQ{
		connectionPool: connectionPool,
		queueManager:   queueManager,
		asyncLogger:    asyncLogger,
	}

	log.Println("RabbitMQ初始化成功")
	return nil
}

// GetRabbitMQ 获取RabbitMQ实例
func GetRabbitMQ() *RabbitMQ {
	return rabbitMQ
}

// GetConnectionPool 获取连接池
func (r *RabbitMQ) GetConnectionPool() *ConnectionPool {
	return r.connectionPool
}

// GetQueueManager 获取队列管理器
func (r *RabbitMQ) GetQueueManager() *QueueManager {
	return r.queueManager
}

// GetAsyncLogger 获取异步日志处理器
func (r *RabbitMQ) GetAsyncLogger() *AsyncLogger {
	return r.asyncLogger
}

// GetStats 获取RabbitMQ统计信息
func (r *RabbitMQ) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})

	// 连接池统计
	if r.connectionPool != nil {
		stats["connection_pool"] = r.connectionPool.GetStats()
	}

	// 异步日志统计
	if r.asyncLogger != nil {
		logStats, err := r.asyncLogger.GetStats()
		if err == nil {
			stats["async_logger"] = logStats
		}
	}

	return stats
}

// Close 关闭RabbitMQ
func (r *RabbitMQ) Close() error {
	var lastErr error

	// 关闭异步日志处理器
	if r.asyncLogger != nil {
		if err := r.asyncLogger.Close(); err != nil {
			log.Printf("关闭异步日志处理器失败: %v", err)
			lastErr = err
		}
	}

	// 关闭连接池
	if r.connectionPool != nil {
		if err := r.connectionPool.Close(); err != nil {
			log.Printf("关闭连接池失败: %v", err)
			lastErr = err
		}
	}

	log.Println("RabbitMQ已关闭")
	return lastErr
}

// Close 全局关闭函数
func Close() error {
	if rabbitMQ != nil {
		return rabbitMQ.Close()
	}
	return nil
}
