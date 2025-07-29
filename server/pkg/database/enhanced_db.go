/*
 * @Author: JimZhang
 * @Date: 2025-07-29 16:30:00
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-29 16:30:00
 * @FilePath: /server/pkg/database/enhanced_db.go
 * @Description: 增强的数据库管理器
 */

package database

import (
	"context"
	"fmt"
	"log"
	"server/common/config"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBStats 数据库统计信息
type DBStats struct {
	OpenConnections   int           `json:"open_connections"`
	InUseConnections  int           `json:"in_use_connections"`
	IdleConnections   int           `json:"idle_connections"`
	WaitCount         int64         `json:"wait_count"`
	WaitDuration      time.Duration `json:"wait_duration"`
	MaxIdleClosed     int64         `json:"max_idle_closed"`
	MaxIdleTimeClosed int64         `json:"max_idle_time_closed"`
	MaxLifetimeClosed int64         `json:"max_lifetime_closed"`
}

// EnhancedDB 增强的数据库管理器
type EnhancedDB struct {
	master      *gorm.DB
	slaves      []*gorm.DB
	stats       *DBStats
	mu          sync.RWMutex
	slaveIndex  int
	healthCheck *HealthChecker
}

// HealthChecker 健康检查器
type HealthChecker struct {
	interval time.Duration
	timeout  time.Duration
	stopCh   chan struct{}
}

// QueryMetrics 查询指标
type QueryMetrics struct {
	TotalQueries   int64         `json:"total_queries"`
	SlowQueries    int64         `json:"slow_queries"`
	FailedQueries  int64         `json:"failed_queries"`
	AverageLatency time.Duration `json:"average_latency"`
	MaxLatency     time.Duration `json:"max_latency"`
	LastQueryTime  time.Time     `json:"last_query_time"`
}

var (
	enhancedDB *EnhancedDB
	dbOnce     sync.Once
	metrics    *QueryMetrics
)

// InitEnhancedDB 初始化增强数据库管理器
func InitEnhancedDB() (*EnhancedDB, error) {
	var err error
	dbOnce.Do(func() {
		enhancedDB, err = newEnhancedDB()
	})
	return enhancedDB, err
}

// GetEnhancedDB 获取增强数据库管理器实例
func GetEnhancedDB() *EnhancedDB {
	return enhancedDB
}

// newEnhancedDB 创建新的增强数据库管理器
func newEnhancedDB() (*EnhancedDB, error) {
	cfg := config.Config.DB

	// 创建主数据库连接
	master, err := createDBConnection(true)
	if err != nil {
		return nil, fmt.Errorf("创建主数据库连接失败: %v", err)
	}

	// 创建从数据库连接（如果启用读写分离）
	var slaves []*gorm.DB
	if cfg.ReadReplicas.Enabled && len(cfg.ReadReplicas.Hosts) > 0 {
		for _, host := range cfg.ReadReplicas.Hosts {
			// TODO: 这里需要为每个从库创建单独的连接
			// 暂时跳过从库配置
			log.Printf("跳过从数据库配置: %s", host)
		}
		log.Printf("从数据库配置已跳过，使用主库")
	}

	// 初始化指标
	metrics = &QueryMetrics{
		LastQueryTime: time.Now(),
	}

	db := &EnhancedDB{
		master: master,
		slaves: slaves,
		stats:  &DBStats{},
		healthCheck: &HealthChecker{
			interval: 30 * time.Second,
			timeout:  5 * time.Second,
			stopCh:   make(chan struct{}),
		},
	}

	// 启动健康检查
	go db.startHealthCheck()

	// 启动统计信息收集
	go db.collectStats()

	log.Println("增强数据库管理器初始化成功")
	return db, nil
}

// createDBConnection 创建数据库连接
func createDBConnection(isMaster bool) (*gorm.DB, error) {
	dbConfig := config.Config.DB

	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=30s",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
		dbConfig.Charset,
	)

	// 配置日志级别
	var logLevel logger.LogLevel
	switch dbConfig.LogLevel {
	case 1:
		logLevel = logger.Silent
	case 2:
		logLevel = logger.Error
	case 3:
		logLevel = logger.Warn
	case 4:
		logLevel = logger.Info
	default:
		logLevel = logger.Info
	}

	// 创建自定义日志器
	customLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Duration(dbConfig.SlowThreshold) * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// GORM配置
	gormConfig := &gorm.Config{
		Logger:                                   customLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
		PrepareStmt:                              dbConfig.PrepareStmt,
		// 跳过默认事务（提高性能）
		SkipDefaultTransaction: true,
		// 命名策略
		NamingStrategy: nil,
	}

	// 创建数据库连接
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %v", err)
	}

	// 获取底层SQL DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层数据库连接失败: %v", err)
	}

	// 配置连接池
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.SetConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(dbConfig.ConnMaxIdleTime) * time.Second)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("数据库连接测试失败: %v", err)
	}

	connType := "主库"
	if !isMaster {
		connType = "从库"
	}
	log.Printf("%s连接成功: %s:%d", connType, dbConfig.Host, dbConfig.Port)

	return db, nil
}

// Master 获取主数据库连接
func (edb *EnhancedDB) Master() *gorm.DB {
	return edb.master
}

// Slave 获取从数据库连接（负载均衡）
func (edb *EnhancedDB) Slave() *gorm.DB {
	if len(edb.slaves) == 0 {
		return edb.master
	}

	edb.mu.Lock()
	defer edb.mu.Unlock()

	// 轮询选择从库
	db := edb.slaves[edb.slaveIndex]
	edb.slaveIndex = (edb.slaveIndex + 1) % len(edb.slaves)

	return db
}

// Read 读操作（使用从库）
func (edb *EnhancedDB) Read() *gorm.DB {
	return edb.Slave()
}

// Write 写操作（使用主库）
func (edb *EnhancedDB) Write() *gorm.DB {
	return edb.Master()
}

// Transaction 事务操作
func (edb *EnhancedDB) Transaction(fn func(*gorm.DB) error) error {
	return edb.Master().Transaction(fn)
}

// WithTimeout 带超时的数据库操作
func (edb *EnhancedDB) WithTimeout(timeout time.Duration, fn func(*gorm.DB) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- fn(edb.Master().WithContext(ctx))
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return fmt.Errorf("数据库操作超时")
	}
}

// BatchInsert 批量插入
func (edb *EnhancedDB) BatchInsert(data interface{}, batchSize int) error {
	start := time.Now()
	defer func() {
		edb.updateMetrics(time.Since(start), nil)
	}()

	err := edb.Master().CreateInBatches(data, batchSize).Error
	if err != nil {
		edb.updateMetrics(time.Since(start), err)
	}
	return err
}

// BulkUpdate 批量更新
func (edb *EnhancedDB) BulkUpdate(model interface{}, updates map[string]interface{}, where string, args ...interface{}) error {
	start := time.Now()
	defer func() {
		edb.updateMetrics(time.Since(start), nil)
	}()

	err := edb.Master().Model(model).Where(where, args...).Updates(updates).Error
	if err != nil {
		edb.updateMetrics(time.Since(start), err)
	}
	return err
}

// startHealthCheck 启动健康检查
func (edb *EnhancedDB) startHealthCheck() {
	ticker := time.NewTicker(edb.healthCheck.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			edb.performHealthCheck()
		case <-edb.healthCheck.stopCh:
			return
		}
	}
}

// performHealthCheck 执行健康检查
func (edb *EnhancedDB) performHealthCheck() {
	// 检查主库
	if err := edb.pingDB(edb.master, "主库"); err != nil {
		log.Printf("主库健康检查失败: %v", err)
	}

	// 检查从库
	for i, slave := range edb.slaves {
		if err := edb.pingDB(slave, fmt.Sprintf("从库%d", i+1)); err != nil {
			log.Printf("从库%d健康检查失败: %v", i+1, err)
		}
	}
}

// pingDB 测试数据库连接
func (edb *EnhancedDB) pingDB(db *gorm.DB, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), edb.healthCheck.timeout)
	defer cancel()

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.PingContext(ctx)
}

// collectStats 收集统计信息
func (edb *EnhancedDB) collectStats() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			edb.updateStats()
		case <-edb.healthCheck.stopCh:
			return
		}
	}
}

// updateStats 更新统计信息
func (edb *EnhancedDB) updateStats() {
	sqlDB, err := edb.master.DB()
	if err != nil {
		return
	}

	stats := sqlDB.Stats()
	edb.mu.Lock()
	edb.stats.OpenConnections = stats.OpenConnections
	edb.stats.InUseConnections = stats.InUse
	edb.stats.IdleConnections = stats.Idle
	edb.stats.WaitCount = stats.WaitCount
	edb.stats.WaitDuration = stats.WaitDuration
	edb.stats.MaxIdleClosed = stats.MaxIdleClosed
	edb.stats.MaxIdleTimeClosed = stats.MaxIdleTimeClosed
	edb.stats.MaxLifetimeClosed = stats.MaxLifetimeClosed
	edb.mu.Unlock()
}

// updateMetrics 更新查询指标
func (edb *EnhancedDB) updateMetrics(duration time.Duration, err error) {
	metrics.TotalQueries++
	metrics.LastQueryTime = time.Now()

	if err != nil {
		metrics.FailedQueries++
	}

	// 更新延迟统计
	if duration > metrics.MaxLatency {
		metrics.MaxLatency = duration
	}

	// 简单的平均延迟计算
	metrics.AverageLatency = (metrics.AverageLatency + duration) / 2

	// 慢查询统计
	slowThreshold := time.Duration(config.Config.DB.SlowThreshold) * time.Millisecond
	if duration > slowThreshold {
		metrics.SlowQueries++
	}
}

// GetStats 获取数据库统计信息
func (edb *EnhancedDB) GetStats() map[string]interface{} {
	edb.mu.RLock()
	defer edb.mu.RUnlock()

	return map[string]interface{}{
		"connection_stats":      edb.stats,
		"query_metrics":         metrics,
		"master_status":         "healthy",
		"slaves_count":          len(edb.slaves),
		"read_replicas_enabled": config.Config.DB.ReadReplicas.Enabled,
	}
}

// Close 关闭数据库连接
func (edb *EnhancedDB) Close() error {
	// 停止健康检查
	close(edb.healthCheck.stopCh)

	// 关闭主库连接
	if sqlDB, err := edb.master.DB(); err == nil {
		sqlDB.Close()
	}

	// 关闭从库连接
	for _, slave := range edb.slaves {
		if sqlDB, err := slave.DB(); err == nil {
			sqlDB.Close()
		}
	}

	log.Println("增强数据库管理器已关闭")
	return nil
}
