/*
 * @Author: JimZhang
 * @Date: 2025-07-29 18:20:00
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-29 18:20:00
 * @FilePath: /server/pkg/bootstrap/bootstrap.go
 * @Description: 系统启动引导程序
 */

package bootstrap

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"server/pkg/breaker"
	"server/pkg/database"
	"server/pkg/limiter"
	"server/pkg/logger"
	"server/pkg/metrics"
	"server/pkg/rabbitmq"
	"server/pkg/redis"
	"syscall"
	"time"
)

// Application 应用程序结构
type Application struct {
	// 中间件组件
	metricsCollector *metrics.MetricsCollector
	asyncLogger      *logger.AsyncLogger
	limiterManager   *limiter.LimiterManager
	breakerManager   *breaker.BreakerManager
	enhancedDB       *database.EnhancedDB
	
	// 消息队列
	rabbitMQ *rabbitmq.RabbitMQ
	
	// 控制
	ctx    context.Context
	cancel context.CancelFunc
}

// Config 启动配置
type Config struct {
	// 是否启用RabbitMQ
	EnableRabbitMQ bool
	// 是否启用Redis缓存管理器
	EnableRedisManagers bool
	// 是否启用增强数据库
	EnableEnhancedDB bool
	// 是否启用性能监控
	EnableMetrics bool
	// 是否启用异步日志
	EnableAsyncLogger bool
	// 是否启用限流器
	EnableRateLimiter bool
	// 是否启用熔断器
	EnableCircuitBreaker bool
}

// DefaultConfig 默认配置
var DefaultConfig = Config{
	EnableRabbitMQ:       true,
	EnableRedisManagers:  true,
	EnableEnhancedDB:     true,
	EnableMetrics:        true,
	EnableAsyncLogger:    true,
	EnableRateLimiter:    true,
	EnableCircuitBreaker: true,
}

// NewApplication 创建应用程序实例
func NewApplication(config ...Config) (*Application, error) {
	cfg := DefaultConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	ctx, cancel := context.WithCancel(context.Background())

	app := &Application{
		ctx:    ctx,
		cancel: cancel,
	}

	// 初始化各个组件
	if err := app.initializeComponents(cfg); err != nil {
		cancel()
		return nil, fmt.Errorf("初始化组件失败: %v", err)
	}

	log.Println("🚀 应用程序初始化完成")
	return app, nil
}

// initializeComponents 初始化组件
func (app *Application) initializeComponents(cfg Config) error {
	log.Println("开始初始化系统组件...")

	// 1. 初始化RabbitMQ（如果启用）
	if cfg.EnableRabbitMQ {
		log.Println("初始化RabbitMQ...")
		if err := rabbitmq.Init(); err != nil {
			log.Printf("⚠️  RabbitMQ初始化失败，跳过: %v", err)
		} else {
			app.rabbitMQ = rabbitmq.GetRabbitMQ()
			log.Println("✅ RabbitMQ初始化成功")
		}
	}

	// 2. 初始化Redis管理器（如果启用）
	if cfg.EnableRedisManagers {
		log.Println("初始化Redis管理器...")
		
		// 初始化缓存管理器
		redis.InitCacheManager()
		log.Println("✅ Redis缓存管理器初始化成功")
		
		// 初始化分布式锁管理器
		redis.InitLockManager()
		log.Println("✅ Redis分布式锁管理器初始化成功")
		
		// 初始化会话管理器
		redis.InitSessionManager()
		log.Println("✅ Redis会话管理器初始化成功")
	}

	// 3. 初始化增强数据库（如果启用）
	if cfg.EnableEnhancedDB {
		log.Println("初始化增强数据库...")
		enhancedDB, err := database.InitEnhancedDB()
		if err != nil {
			log.Printf("⚠️  增强数据库初始化失败，使用默认数据库: %v", err)
		} else {
			app.enhancedDB = enhancedDB
			log.Println("✅ 增强数据库初始化成功")
		}
	}

	// 4. 初始化性能监控（如果启用）
	if cfg.EnableMetrics {
		log.Println("初始化性能监控...")
		metricsCollector, err := metrics.InitMetricsCollector()
		if err != nil {
			return fmt.Errorf("性能监控初始化失败: %v", err)
		}
		app.metricsCollector = metricsCollector
		log.Println("✅ 性能监控初始化成功")
	}

	// 5. 初始化异步日志（如果启用）
	if cfg.EnableAsyncLogger {
		log.Println("初始化异步日志...")
		asyncLogger, err := logger.InitAsyncLogger()
		if err != nil {
			return fmt.Errorf("异步日志初始化失败: %v", err)
		}
		app.asyncLogger = asyncLogger
		log.Println("✅ 异步日志初始化成功")
	}

	// 6. 初始化限流器（如果启用）
	if cfg.EnableRateLimiter {
		log.Println("初始化限流器...")
		limiterManager, err := limiter.InitLimiterManager()
		if err != nil {
			return fmt.Errorf("限流器初始化失败: %v", err)
		}
		app.limiterManager = limiterManager
		log.Println("✅ 限流器初始化成功")
	}

	// 7. 初始化熔断器（如果启用）
	if cfg.EnableCircuitBreaker {
		log.Println("初始化熔断器...")
		breakerManager, err := breaker.InitBreakerManager()
		if err != nil {
			return fmt.Errorf("熔断器初始化失败: %v", err)
		}
		app.breakerManager = breakerManager
		log.Println("✅ 熔断器初始化成功")
	}

	return nil
}

// Start 启动应用程序
func (app *Application) Start() error {
	log.Println("🎯 应用程序启动中...")

	// 启动健康检查
	go app.startHealthCheck()

	// 启动性能监控报告
	if app.metricsCollector != nil {
		go app.startPerformanceReporting()
	}

	log.Println("✅ 应用程序启动完成")
	return nil
}

// Stop 停止应用程序
func (app *Application) Stop() error {
	log.Println("🛑 应用程序停止中...")

	// 取消上下文
	app.cancel()

	// 关闭各个组件
	if app.asyncLogger != nil {
		app.asyncLogger.Close()
	}

	if app.metricsCollector != nil {
		app.metricsCollector.Close()
	}

	if app.enhancedDB != nil {
		app.enhancedDB.Close()
	}

	if app.rabbitMQ != nil {
		app.rabbitMQ.Close()
	}

	log.Println("✅ 应用程序已停止")
	return nil
}

// WaitForShutdown 等待关闭信号
func (app *Application) WaitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("收到信号: %v", sig)

	// 优雅关闭
	app.Stop()
}

// startHealthCheck 启动健康检查
func (app *Application) startHealthCheck() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			app.performHealthCheck()
		case <-app.ctx.Done():
			return
		}
	}
}

// performHealthCheck 执行健康检查
func (app *Application) performHealthCheck() {
	if app.asyncLogger != nil {
		stats := app.asyncLogger.GetStats()
		if stats.DroppedLogs > 0 {
			log.Printf("⚠️  异步日志丢弃数量: %d", stats.DroppedLogs)
		}
	}

	if app.metricsCollector != nil {
		systemMetrics := app.metricsCollector.GetSystemMetrics()
		if systemMetrics.MemoryPercent > 90 {
			log.Printf("⚠️  内存使用率过高: %.2f%%", systemMetrics.MemoryPercent)
		}
		
		if systemMetrics.GoroutineCount > 10000 {
			log.Printf("⚠️  协程数量过多: %d", systemMetrics.GoroutineCount)
		}
	}
}

// startPerformanceReporting 启动性能报告
func (app *Application) startPerformanceReporting() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			app.generatePerformanceReport()
		case <-app.ctx.Done():
			return
		}
	}
}

// generatePerformanceReport 生成性能报告
func (app *Application) generatePerformanceReport() {
	if app.metricsCollector == nil {
		return
	}

	httpMetrics := app.metricsCollector.GetHTTPMetrics()
	systemMetrics := app.metricsCollector.GetSystemMetrics()

	log.Printf("📊 性能报告 - QPS: %.2f, 平均响应时间: %v, 内存使用: %.2f%%, 协程数: %d",
		httpMetrics.RequestsPerSec,
		httpMetrics.AvgResponseTime,
		systemMetrics.MemoryPercent,
		systemMetrics.GoroutineCount,
	)

	// 记录到异步日志
	if app.asyncLogger != nil {
		app.asyncLogger.Info("性能报告", map[string]interface{}{
			"qps":              httpMetrics.RequestsPerSec,
			"avg_response_time": httpMetrics.AvgResponseTime.Milliseconds(),
			"memory_percent":   systemMetrics.MemoryPercent,
			"goroutine_count":  systemMetrics.GoroutineCount,
			"total_requests":   httpMetrics.TotalRequests,
			"error_rate":       httpMetrics.ErrorRate,
		})
	}
}

// GetMetricsCollector 获取指标收集器
func (app *Application) GetMetricsCollector() *metrics.MetricsCollector {
	return app.metricsCollector
}

// GetAsyncLogger 获取异步日志器
func (app *Application) GetAsyncLogger() *logger.AsyncLogger {
	return app.asyncLogger
}

// GetLimiterManager 获取限流器管理器
func (app *Application) GetLimiterManager() *limiter.LimiterManager {
	return app.limiterManager
}

// GetBreakerManager 获取熔断器管理器
func (app *Application) GetBreakerManager() *breaker.BreakerManager {
	return app.breakerManager
}

// GetEnhancedDB 获取增强数据库
func (app *Application) GetEnhancedDB() *database.EnhancedDB {
	return app.enhancedDB
}

// GetRabbitMQ 获取RabbitMQ
func (app *Application) GetRabbitMQ() *rabbitmq.RabbitMQ {
	return app.rabbitMQ
}

// PrintSystemInfo 打印系统信息
func (app *Application) PrintSystemInfo() {
	fmt.Println("========================================")
	fmt.Println("         系统组件状态")
	fmt.Println("========================================")
	
	if app.rabbitMQ != nil {
		fmt.Println("✅ RabbitMQ: 已启用")
	} else {
		fmt.Println("❌ RabbitMQ: 未启用")
	}
	
	if app.metricsCollector != nil {
		fmt.Println("✅ 性能监控: 已启用")
	} else {
		fmt.Println("❌ 性能监控: 未启用")
	}
	
	if app.asyncLogger != nil {
		fmt.Println("✅ 异步日志: 已启用")
	} else {
		fmt.Println("❌ 异步日志: 未启用")
	}
	
	if app.limiterManager != nil {
		fmt.Println("✅ 限流器: 已启用")
	} else {
		fmt.Println("❌ 限流器: 未启用")
	}
	
	if app.breakerManager != nil {
		fmt.Println("✅ 熔断器: 已启用")
	} else {
		fmt.Println("❌ 熔断器: 未启用")
	}
	
	if app.enhancedDB != nil {
		fmt.Println("✅ 增强数据库: 已启用")
	} else {
		fmt.Println("❌ 增强数据库: 未启用")
	}
	
	fmt.Println("========================================")
}
