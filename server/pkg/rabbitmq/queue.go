/*
 * @Author: JimZhang
 * @Date: 2025-07-29 15:35:00
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-29 15:35:00
 * @FilePath: /server/pkg/rabbitmq/queue.go
 * @Description: RabbitMQ队列管理
 */

package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// QueueManager 队列管理器
type QueueManager struct {
	pool *ConnectionPool
}

// Message 消息结构
type Message struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
	Retry     int                    `json:"retry"`
	MaxRetry  int                    `json:"max_retry"`
}

// QueueConfig 队列配置
type QueueConfig struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

// ExchangeConfig 交换机配置
type ExchangeConfig struct {
	Name       string
	Type       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

var queueManager *QueueManager

// InitQueueManager 初始化队列管理器
func InitQueueManager() (*QueueManager, error) {
	pool := GetConnectionPool()
	if pool == nil {
		return nil, fmt.Errorf("连接池未初始化")
	}

	queueManager = &QueueManager{
		pool: pool,
	}

	// 初始化默认队列和交换机
	err := queueManager.setupDefaultQueues()
	if err != nil {
		return nil, fmt.Errorf("设置默认队列失败: %v", err)
	}

	log.Println("RabbitMQ队列管理器初始化成功")
	return queueManager, nil
}

// GetQueueManager 获取队列管理器实例
func GetQueueManager() *QueueManager {
	return queueManager
}

// setupDefaultQueues 设置默认队列
func (qm *QueueManager) setupDefaultQueues() error {
	// 默认交换机配置
	exchanges := []ExchangeConfig{
		{
			Name:    "logs.direct",
			Type:    "direct",
			Durable: true,
		},
		{
			Name:    "notifications.fanout",
			Type:    "fanout",
			Durable: true,
		},
		{
			Name:    "tasks.topic",
			Type:    "topic",
			Durable: true,
		},
	}

	// 默认队列配置
	queues := []QueueConfig{
		{
			Name:    "operation_logs",
			Durable: true,
		},
		{
			Name:    "login_logs",
			Durable: true,
		},
		{
			Name:    "email_notifications",
			Durable: true,
		},
		{
			Name:    "sms_notifications",
			Durable: true,
		},
		{
			Name:    "async_tasks",
			Durable: true,
		},
	}

	ch, err := qm.pool.GetChannel()
	if err != nil {
		return err
	}
	defer qm.pool.ReturnChannel(ch)

	// 创建交换机
	for _, exchange := range exchanges {
		err = ch.ExchangeDeclare(
			exchange.Name,
			exchange.Type,
			exchange.Durable,
			exchange.AutoDelete,
			exchange.Internal,
			exchange.NoWait,
			exchange.Args,
		)
		if err != nil {
			return fmt.Errorf("创建交换机 %s 失败: %v", exchange.Name, err)
		}
	}

	// 创建队列
	for _, queue := range queues {
		_, err = ch.QueueDeclare(
			queue.Name,
			queue.Durable,
			queue.AutoDelete,
			queue.Exclusive,
			queue.NoWait,
			queue.Args,
		)
		if err != nil {
			return fmt.Errorf("创建队列 %s 失败: %v", queue.Name, err)
		}
	}

	// 绑定队列到交换机
	bindings := map[string][]string{
		"logs.direct":          {"operation_logs", "login_logs"},
		"notifications.fanout": {"email_notifications", "sms_notifications"},
		"tasks.topic":          {"async_tasks"},
	}

	for exchange, queueNames := range bindings {
		for _, queueName := range queueNames {
			routingKey := queueName
			if exchange == "notifications.fanout" {
				routingKey = ""
			}

			err = ch.QueueBind(
				queueName,
				routingKey,
				exchange,
				false,
				nil,
			)
			if err != nil {
				return fmt.Errorf("绑定队列 %s 到交换机 %s 失败: %v", queueName, exchange, err)
			}
		}
	}

	log.Println("默认队列和交换机设置完成")
	return nil
}

// PublishMessage 发布消息
func (qm *QueueManager) PublishMessage(exchange, routingKey string, message *Message) error {
	ch, err := qm.pool.GetChannel()
	if err != nil {
		return fmt.Errorf("获取通道失败: %v", err)
	}
	defer qm.pool.ReturnChannel(ch)

	// 序列化消息
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %v", err)
	}

	// 发布消息
	err = ch.Publish(
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // 持久化消息
			Timestamp:    time.Now(),
			MessageId:    message.ID,
		},
	)

	if err != nil {
		return fmt.Errorf("发布消息失败: %v", err)
	}

	return nil
}

// PublishToQueue 直接发布消息到队列
func (qm *QueueManager) PublishToQueue(queueName string, message *Message) error {
	return qm.PublishMessage("", queueName, message)
}

// ConsumeMessages 消费消息
func (qm *QueueManager) ConsumeMessages(queueName string, handler func(*Message) error) error {
	ch, err := qm.pool.GetChannel()
	if err != nil {
		return fmt.Errorf("获取通道失败: %v", err)
	}

	// 设置QoS
	err = ch.Qos(
		10,    // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		qm.pool.ReturnChannel(ch)
		return fmt.Errorf("设置QoS失败: %v", err)
	}

	msgs, err := ch.Consume(
		queueName,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		qm.pool.ReturnChannel(ch)
		return fmt.Errorf("开始消费失败: %v", err)
	}

	log.Printf("开始消费队列: %s", queueName)

	go func() {
		defer qm.pool.ReturnChannel(ch)

		for d := range msgs {
			var message Message
			err := json.Unmarshal(d.Body, &message)
			if err != nil {
				log.Printf("反序列化消息失败: %v", err)
				d.Nack(false, false) // 拒绝消息，不重新入队
				continue
			}

			// 处理消息
			err = handler(&message)
			if err != nil {
				log.Printf("处理消息失败: %v", err)

				// 检查重试次数
				if message.Retry < message.MaxRetry {
					message.Retry++
					// 重新发布消息
					qm.PublishToQueue(queueName, &message)
				}

				d.Nack(false, false)
			} else {
				d.Ack(false)
			}
		}
	}()

	return nil
}

// CreateQueue 创建队列
func (qm *QueueManager) CreateQueue(config QueueConfig) error {
	ch, err := qm.pool.GetChannel()
	if err != nil {
		return err
	}
	defer qm.pool.ReturnChannel(ch)

	_, err = ch.QueueDeclare(
		config.Name,
		config.Durable,
		config.AutoDelete,
		config.Exclusive,
		config.NoWait,
		config.Args,
	)

	return err
}

// CreateExchange 创建交换机
func (qm *QueueManager) CreateExchange(config ExchangeConfig) error {
	ch, err := qm.pool.GetChannel()
	if err != nil {
		return err
	}
	defer qm.pool.ReturnChannel(ch)

	err = ch.ExchangeDeclare(
		config.Name,
		config.Type,
		config.Durable,
		config.AutoDelete,
		config.Internal,
		config.NoWait,
		config.Args,
	)

	return err
}

// GetQueueInfo 获取队列信息
func (qm *QueueManager) GetQueueInfo(queueName string) (amqp.Queue, error) {
	ch, err := qm.pool.GetChannel()
	if err != nil {
		return amqp.Queue{}, err
	}
	defer qm.pool.ReturnChannel(ch)

	return ch.QueueInspect(queueName)
}
