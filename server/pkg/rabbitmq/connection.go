/*
 * @Author: JimZhang
 * @Date: 2025-07-29 15:30:00
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-29 15:30:00
 * @FilePath: /server/pkg/rabbitmq/connection.go
 * @Description: RabbitMQ连接池管理
 */

package rabbitmq

import (
	"fmt"
	"log"
	"server/common/config"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ConnectionPool RabbitMQ连接池
type ConnectionPool struct {
	connections []*amqp.Connection
	channels    chan *amqp.Channel
	mu          sync.RWMutex
	closed      bool
}

var (
	pool *ConnectionPool
	once sync.Once
)

// InitConnectionPool 初始化RabbitMQ连接池
func InitConnectionPool() (*ConnectionPool, error) {
	var err error
	once.Do(func() {
		pool, err = newConnectionPool()
	})
	return pool, err
}

// GetConnectionPool 获取连接池实例
func GetConnectionPool() *ConnectionPool {
	return pool
}

// newConnectionPool 创建新的连接池
func newConnectionPool() (*ConnectionPool, error) {
	cfg := config.Config.RabbitMQ

	// 构建连接URL
	url := fmt.Sprintf("amqp://%s:%s@%s:%d%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Vhost,
	)

	pool := &ConnectionPool{
		connections: make([]*amqp.Connection, 0, cfg.MaxConnections),
		channels:    make(chan *amqp.Channel, cfg.MaxChannels),
		closed:      false,
	}

	// 创建连接
	for i := 0; i < cfg.MaxConnections; i++ {
		conn, err := amqp.Dial(url)
		if err != nil {
			// 清理已创建的连接
			pool.Close()
			return nil, fmt.Errorf("创建RabbitMQ连接失败: %v", err)
		}
		pool.connections = append(pool.connections, conn)
	}

	// 预创建通道
	for i := 0; i < cfg.MaxChannels; i++ {
		ch, err := pool.createChannel()
		if err != nil {
			log.Printf("预创建通道失败: %v", err)
			continue
		}
		pool.channels <- ch
	}

	// 启动连接监控
	go pool.monitorConnections()

	log.Printf("RabbitMQ连接池初始化成功: %d个连接, %d个通道",
		len(pool.connections), len(pool.channels))

	return pool, nil
}

// GetChannel 获取通道
func (p *ConnectionPool) GetChannel() (*amqp.Channel, error) {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return nil, fmt.Errorf("连接池已关闭")
	}
	p.mu.RUnlock()

	select {
	case ch := <-p.channels:
		// 检查通道是否可用
		if ch.IsClosed() {
			// 重新创建通道
			newCh, err := p.createChannel()
			if err != nil {
				return nil, err
			}
			return newCh, nil
		}
		return ch, nil
	default:
		// 通道池为空，创建新通道
		return p.createChannel()
	}
}

// ReturnChannel 归还通道
func (p *ConnectionPool) ReturnChannel(ch *amqp.Channel) {
	if ch == nil || ch.IsClosed() {
		return
	}

	select {
	case p.channels <- ch:
		// 成功归还
	default:
		// 通道池已满，关闭通道
		ch.Close()
	}
}

// createChannel 创建新通道
func (p *ConnectionPool) createChannel() (*amqp.Channel, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if len(p.connections) == 0 {
		return nil, fmt.Errorf("没有可用的连接")
	}

	// 轮询选择连接
	for _, conn := range p.connections {
		if !conn.IsClosed() {
			ch, err := conn.Channel()
			if err == nil {
				return ch, nil
			}
		}
	}

	return nil, fmt.Errorf("无法创建通道，所有连接都不可用")
}

// monitorConnections 监控连接状态
func (p *ConnectionPool) monitorConnections() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.checkConnections()
		}
	}
}

// checkConnections 检查连接状态
func (p *ConnectionPool) checkConnections() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return
	}

	cfg := config.Config.RabbitMQ
	url := fmt.Sprintf("amqp://%s:%s@%s:%d%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Vhost,
	)

	// 检查并重连断开的连接
	for i, conn := range p.connections {
		if conn.IsClosed() {
			log.Printf("检测到连接 %d 已断开，尝试重连...", i)

			newConn, err := amqp.Dial(url)
			if err != nil {
				log.Printf("重连失败: %v", err)
				continue
			}

			p.connections[i] = newConn
			log.Printf("连接 %d 重连成功", i)
		}
	}
}

// GetStats 获取连接池统计信息
func (p *ConnectionPool) GetStats() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	activeConnections := 0
	for _, conn := range p.connections {
		if !conn.IsClosed() {
			activeConnections++
		}
	}

	return map[string]interface{}{
		"total_connections":  len(p.connections),
		"active_connections": activeConnections,
		"available_channels": len(p.channels),
		"max_channels":       cap(p.channels),
		"closed":             p.closed,
	}
}

// Close 关闭连接池
func (p *ConnectionPool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil
	}

	p.closed = true

	// 关闭所有通道
	close(p.channels)
	for ch := range p.channels {
		if ch != nil && !ch.IsClosed() {
			ch.Close()
		}
	}

	// 关闭所有连接
	for _, conn := range p.connections {
		if conn != nil && !conn.IsClosed() {
			conn.Close()
		}
	}

	log.Println("RabbitMQ连接池已关闭")
	return nil
}
