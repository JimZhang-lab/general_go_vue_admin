/*
 * @Author: JimZhang
 * @Date: 2025-07-29 17:20:00
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-29 17:20:00
 * @FilePath: /server/pkg/logger/async_logger.go
 * @Description: 异步日志处理器
 */

package logger

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"server/pkg/rabbitmq"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// LogLevel 日志级别
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// LogEntry 日志条目
type LogEntry struct {
	ID        string                 `json:"id"`
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields"`
	Timestamp time.Time              `json:"timestamp"`
	Caller    string                 `json:"caller"`
	Stack     string                 `json:"stack,omitempty"`
}

// OperationLogEntry 操作日志条目
type OperationLogEntry struct {
	AdminID   int                    `json:"admin_id"`
	Username  string                 `json:"username"`
	Method    string                 `json:"method"`
	URL       string                 `json:"url"`
	IP        string                 `json:"ip"`
	UserAgent string                 `json:"user_agent"`
	Extra     map[string]interface{} `json:"extra"`
	Timestamp time.Time              `json:"timestamp"`
}

// LoginLogEntry 登录日志条目
type LoginLogEntry struct {
	Username  string    `json:"username"`
	IP        string    `json:"ip"`
	Location  string    `json:"location"`
	Browser   string    `json:"browser"`
	OS        string    `json:"os"`
	Status    int       `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// AsyncLogger 异步日志处理器
type AsyncLogger struct {
	// 日志通道
	logChan       chan *LogEntry
	operationChan chan *OperationLogEntry
	loginChan     chan *LoginLogEntry

	// 配置
	bufferSize    int
	flushInterval time.Duration
	maxRetries    int

	// 控制
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// 统计
	stats *LoggerStats
	mutex sync.RWMutex

	// 外部依赖
	queueManager *rabbitmq.QueueManager
	logger       *logrus.Logger
}

// LoggerStats 日志统计信息
type LoggerStats struct {
	TotalLogs      int64     `json:"total_logs"`
	OperationLogs  int64     `json:"operation_logs"`
	LoginLogs      int64     `json:"login_logs"`
	ErrorLogs      int64     `json:"error_logs"`
	DroppedLogs    int64     `json:"dropped_logs"`
	ProcessedLogs  int64     `json:"processed_logs"`
	LastLogTime    time.Time `json:"last_log_time"`
	BufferUsage    int       `json:"buffer_usage"`
	MaxBufferUsage int       `json:"max_buffer_usage"`
}

var (
	asyncLogger *AsyncLogger
	loggerOnce  sync.Once
)

// InitAsyncLogger 初始化异步日志处理器
func InitAsyncLogger() (*AsyncLogger, error) {
	var err error
	loggerOnce.Do(func() {
		asyncLogger, err = newAsyncLogger()
	})
	return asyncLogger, err
}

// GetAsyncLogger 获取异步日志处理器实例
func GetAsyncLogger() *AsyncLogger {
	return asyncLogger
}

// newAsyncLogger 创建新的异步日志处理器
func newAsyncLogger() (*AsyncLogger, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// 创建logrus实例
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	al := &AsyncLogger{
		logChan:       make(chan *LogEntry, 10000),         // 10k缓冲
		operationChan: make(chan *OperationLogEntry, 5000), // 5k缓冲
		loginChan:     make(chan *LoginLogEntry, 1000),     // 1k缓冲
		bufferSize:    10000,
		flushInterval: time.Second * 5,
		maxRetries:    3,
		ctx:           ctx,
		cancel:        cancel,
		stats:         &LoggerStats{},
		queueManager:  rabbitmq.GetQueueManager(),
		logger:        logger,
	}

	// 启动处理协程
	al.startWorkers()

	log.Println("异步日志处理器初始化成功")
	return al, nil
}

// startWorkers 启动工作协程
func (al *AsyncLogger) startWorkers() {
	// 启动通用日志处理器
	al.wg.Add(1)
	go al.logWorker()

	// 启动操作日志处理器
	al.wg.Add(1)
	go al.operationLogWorker()

	// 启动登录日志处理器
	al.wg.Add(1)
	go al.loginLogWorker()

	// 启动统计信息更新器
	al.wg.Add(1)
	go al.statsWorker()
}

// logWorker 通用日志处理工作协程
func (al *AsyncLogger) logWorker() {
	defer al.wg.Done()

	batch := make([]*LogEntry, 0, 100)
	ticker := time.NewTicker(al.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case entry := <-al.logChan:
			batch = append(batch, entry)
			al.updateStats(1, 0, 0, 0)

			// 批量处理
			if len(batch) >= 100 {
				al.processBatch(batch)
				batch = batch[:0]
			}

		case <-ticker.C:
			// 定时刷新
			if len(batch) > 0 {
				al.processBatch(batch)
				batch = batch[:0]
			}

		case <-al.ctx.Done():
			// 处理剩余日志
			if len(batch) > 0 {
				al.processBatch(batch)
			}
			return
		}
	}
}

// operationLogWorker 操作日志处理工作协程
func (al *AsyncLogger) operationLogWorker() {
	defer al.wg.Done()

	for {
		select {
		case entry := <-al.operationChan:
			al.processOperationLog(entry)
			al.updateStats(0, 1, 0, 0)

		case <-al.ctx.Done():
			return
		}
	}
}

// loginLogWorker 登录日志处理工作协程
func (al *AsyncLogger) loginLogWorker() {
	defer al.wg.Done()

	for {
		select {
		case entry := <-al.loginChan:
			al.processLoginLog(entry)
			al.updateStats(0, 0, 1, 0)

		case <-al.ctx.Done():
			return
		}
	}
}

// statsWorker 统计信息更新工作协程
func (al *AsyncLogger) statsWorker() {
	defer al.wg.Done()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			al.updateBufferStats()

		case <-al.ctx.Done():
			return
		}
	}
}

// Log 记录通用日志
func (al *AsyncLogger) Log(level LogLevel, message string, fields map[string]interface{}) {
	entry := &LogEntry{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Level:     level,
		Message:   message,
		Fields:    fields,
		Timestamp: time.Now(),
		Caller:    getCaller(2),
	}

	// 如果是错误级别，添加堆栈信息
	if level >= ErrorLevel {
		entry.Stack = getStack(2)
	}

	select {
	case al.logChan <- entry:
		// 成功发送
	default:
		// 通道满了，丢弃日志
		al.updateStats(0, 0, 0, 1)
	}
}

// LogOperation 记录操作日志
func (al *AsyncLogger) LogOperation(adminID int, username, method, url, ip, userAgent string, extra map[string]interface{}) {
	entry := &OperationLogEntry{
		AdminID:   adminID,
		Username:  username,
		Method:    method,
		URL:       url,
		IP:        ip,
		UserAgent: userAgent,
		Extra:     extra,
		Timestamp: time.Now(),
	}

	select {
	case al.operationChan <- entry:
		// 成功发送
	default:
		// 通道满了，丢弃日志
		al.updateStats(0, 0, 0, 1)
	}
}

// LogLogin 记录登录日志
func (al *AsyncLogger) LogLogin(username, ip, location, browser, os string, status int, message string) {
	entry := &LoginLogEntry{
		Username:  username,
		IP:        ip,
		Location:  location,
		Browser:   browser,
		OS:        os,
		Status:    status,
		Message:   message,
		Timestamp: time.Now(),
	}

	select {
	case al.loginChan <- entry:
		// 成功发送
	default:
		// 通道满了，丢弃日志
		al.updateStats(0, 0, 0, 1)
	}
}

// 便捷方法
func (al *AsyncLogger) Debug(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	al.Log(DebugLevel, message, f)
}

func (al *AsyncLogger) Info(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	al.Log(InfoLevel, message, f)
}

func (al *AsyncLogger) Warn(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	al.Log(WarnLevel, message, f)
}

func (al *AsyncLogger) Error(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	al.Log(ErrorLevel, message, f)
}

func (al *AsyncLogger) Fatal(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	al.Log(FatalLevel, message, f)
}

// processBatch 批量处理日志
func (al *AsyncLogger) processBatch(batch []*LogEntry) {
	for _, entry := range batch {
		// 写入到logrus
		logEntry := al.logger.WithFields(logrus.Fields(entry.Fields))
		logEntry = logEntry.WithField("caller", entry.Caller)
		logEntry = logEntry.WithField("log_id", entry.ID)

		switch entry.Level {
		case DebugLevel:
			logEntry.Debug(entry.Message)
		case InfoLevel:
			logEntry.Info(entry.Message)
		case WarnLevel:
			logEntry.Warn(entry.Message)
		case ErrorLevel:
			logEntry.Error(entry.Message)
		case FatalLevel:
			logEntry.Fatal(entry.Message)
		}

		// 如果有RabbitMQ，发送到消息队列
		if al.queueManager != nil {
			al.sendToQueue("general_logs", entry)
		}
	}

	al.updateStats(0, 0, 0, 0) // 更新处理计数
}

// processOperationLog 处理操作日志
func (al *AsyncLogger) processOperationLog(entry *OperationLogEntry) {
	// TODO: 这里应该保存到数据库，暂时只记录日志
	log.Printf("处理操作日志: 用户=%s, 方法=%s, URL=%s, IP=%s",
		entry.Username, entry.Method, entry.URL, entry.IP)

	// 发送到RabbitMQ
	if al.queueManager != nil {
		al.sendToQueue("operation_logs", entry)
	}
}

// processLoginLog 处理登录日志
func (al *AsyncLogger) processLoginLog(entry *LoginLogEntry) {
	// TODO: 这里应该保存到数据库，暂时只记录日志
	log.Printf("处理登录日志: 用户=%s, 状态=%d, IP=%s, 位置=%s, 浏览器=%s, 系统=%s",
		entry.Username, entry.Status, entry.IP, entry.Location, entry.Browser, entry.OS)

	// 发送到RabbitMQ
	if al.queueManager != nil {
		al.sendToQueue("login_logs", entry)
	}
}

// sendToQueue 发送到消息队列
func (al *AsyncLogger) sendToQueue(queueName string, data interface{}) {
	message := &rabbitmq.Message{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Type:      queueName,
		Data:      map[string]interface{}{"log_data": data},
		Timestamp: time.Now(),
		Retry:     0,
		MaxRetry:  al.maxRetries,
	}

	err := al.queueManager.PublishToQueue(queueName, message)
	if err != nil {
		log.Printf("发送日志到队列失败: %v", err)
	}
}

// updateStats 更新统计信息
func (al *AsyncLogger) updateStats(total, operation, login, dropped int64) {
	al.mutex.Lock()
	defer al.mutex.Unlock()

	al.stats.TotalLogs += total
	al.stats.OperationLogs += operation
	al.stats.LoginLogs += login
	al.stats.DroppedLogs += dropped
	al.stats.ProcessedLogs += total + operation + login
	al.stats.LastLogTime = time.Now()
}

// updateBufferStats 更新缓冲区统计信息
func (al *AsyncLogger) updateBufferStats() {
	al.mutex.Lock()
	defer al.mutex.Unlock()

	al.stats.BufferUsage = len(al.logChan) + len(al.operationChan) + len(al.loginChan)
	if al.stats.BufferUsage > al.stats.MaxBufferUsage {
		al.stats.MaxBufferUsage = al.stats.BufferUsage
	}
}

// GetStats 获取统计信息
func (al *AsyncLogger) GetStats() *LoggerStats {
	al.mutex.RLock()
	defer al.mutex.RUnlock()

	// 创建副本
	stats := *al.stats
	return &stats
}

// Flush 刷新所有待处理的日志
func (al *AsyncLogger) Flush() error {
	// 等待所有通道清空
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("刷新超时")
		case <-ticker.C:
			if len(al.logChan) == 0 && len(al.operationChan) == 0 && len(al.loginChan) == 0 {
				return nil
			}
		}
	}
}

// Close 关闭异步日志处理器
func (al *AsyncLogger) Close() error {
	// 刷新待处理的日志
	al.Flush()

	// 取消上下文
	al.cancel()

	// 等待所有工作协程结束
	al.wg.Wait()

	log.Println("异步日志处理器已关闭")
	return nil
}

// 辅助函数

// getCaller 获取调用者信息
func getCaller(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	return fmt.Sprintf("%s:%d", file, line)
}

// getStack 获取堆栈信息
func getStack(skip int) string {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}
