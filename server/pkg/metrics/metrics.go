/*
 * @Author: JimZhang
 * @Date: 2025-07-29 17:40:00
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-29 17:40:00
 * @FilePath: /server/pkg/metrics/metrics.go
 * @Description: 性能监控与指标收集
 */

package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"server/pkg/redis"
	"sync"
	"sync/atomic"
	"time"

	redisClient "github.com/redis/go-redis/v9"
)

// MetricsCollector 指标收集器
type MetricsCollector struct {
	// HTTP指标
	httpMetrics *HTTPMetrics

	// 系统指标
	systemMetrics *SystemMetrics

	// 数据库指标
	dbMetrics *DatabaseMetrics

	// Redis指标
	redisMetrics *RedisMetrics

	// 自定义指标
	customMetrics map[string]*CustomMetric

	// 控制
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	mutex  sync.RWMutex

	// Redis客户端
	client *redisClient.Client
}

// HTTPMetrics HTTP指标
type HTTPMetrics struct {
	TotalRequests   int64                `json:"total_requests"`
	RequestsPerSec  float64              `json:"requests_per_sec"`
	AvgResponseTime time.Duration        `json:"avg_response_time"`
	MaxResponseTime time.Duration        `json:"max_response_time"`
	MinResponseTime time.Duration        `json:"min_response_time"`
	StatusCodes     map[int]int64        `json:"status_codes"`
	PathStats       map[string]*PathStat `json:"path_stats"`
	ErrorRate       float64              `json:"error_rate"`
	LastRequestTime time.Time            `json:"last_request_time"`

	// 内部计数器
	totalResponseTime int64
	requestCount      int64
	errorCount        int64
	mutex             sync.RWMutex
}

// PathStat 路径统计
type PathStat struct {
	Count           int64         `json:"count"`
	AvgResponseTime time.Duration `json:"avg_response_time"`
	MaxResponseTime time.Duration `json:"max_response_time"`
	ErrorCount      int64         `json:"error_count"`
	LastAccess      time.Time     `json:"last_access"`

	totalResponseTime int64
	mutex             sync.RWMutex
}

// SystemMetrics 系统指标
type SystemMetrics struct {
	CPUUsage       float64       `json:"cpu_usage"`
	MemoryUsage    int64         `json:"memory_usage"`
	MemoryTotal    int64         `json:"memory_total"`
	MemoryPercent  float64       `json:"memory_percent"`
	GoroutineCount int           `json:"goroutine_count"`
	GCCount        uint32        `json:"gc_count"`
	GCPauseTime    time.Duration `json:"gc_pause_time"`
	HeapSize       int64         `json:"heap_size"`
	HeapInUse      int64         `json:"heap_in_use"`
	StackInUse     int64         `json:"stack_in_use"`
	LastUpdate     time.Time     `json:"last_update"`
}

// DatabaseMetrics 数据库指标
type DatabaseMetrics struct {
	OpenConnections   int           `json:"open_connections"`
	InUseConnections  int           `json:"in_use_connections"`
	IdleConnections   int           `json:"idle_connections"`
	WaitCount         int64         `json:"wait_count"`
	WaitDuration      time.Duration `json:"wait_duration"`
	MaxIdleClosed     int64         `json:"max_idle_closed"`
	MaxLifetimeClosed int64         `json:"max_lifetime_closed"`
	QueryCount        int64         `json:"query_count"`
	SlowQueryCount    int64         `json:"slow_query_count"`
	ErrorCount        int64         `json:"error_count"`
	LastUpdate        time.Time     `json:"last_update"`
}

// RedisMetrics Redis指标
type RedisMetrics struct {
	ConnectedClients  int64     `json:"connected_clients"`
	UsedMemory        int64     `json:"used_memory"`
	UsedMemoryPeak    int64     `json:"used_memory_peak"`
	KeyspaceHits      int64     `json:"keyspace_hits"`
	KeyspaceMisses    int64     `json:"keyspace_misses"`
	HitRate           float64   `json:"hit_rate"`
	CommandsProcessed int64     `json:"commands_processed"`
	CommandsPerSec    float64   `json:"commands_per_sec"`
	LastUpdate        time.Time `json:"last_update"`
}

// CustomMetric 自定义指标
type CustomMetric struct {
	Name        string      `json:"name"`
	Value       interface{} `json:"value"`
	Type        string      `json:"type"` // counter, gauge, histogram
	Description string      `json:"description"`
	LastUpdate  time.Time   `json:"last_update"`
	mutex       sync.RWMutex
}

var (
	metricsCollector *MetricsCollector
	metricsOnce      sync.Once
)

// InitMetricsCollector 初始化指标收集器
func InitMetricsCollector() (*MetricsCollector, error) {
	var err error
	metricsOnce.Do(func() {
		metricsCollector, err = newMetricsCollector()
	})
	return metricsCollector, err
}

// GetMetricsCollector 获取指标收集器实例
func GetMetricsCollector() *MetricsCollector {
	return metricsCollector
}

// newMetricsCollector 创建新的指标收集器
func newMetricsCollector() (*MetricsCollector, error) {
	ctx, cancel := context.WithCancel(context.Background())

	mc := &MetricsCollector{
		httpMetrics: &HTTPMetrics{
			StatusCodes:     make(map[int]int64),
			PathStats:       make(map[string]*PathStat),
			MinResponseTime: time.Hour, // 初始化为很大的值
		},
		systemMetrics: &SystemMetrics{},
		dbMetrics:     &DatabaseMetrics{},
		redisMetrics:  &RedisMetrics{},
		customMetrics: make(map[string]*CustomMetric),
		ctx:           ctx,
		cancel:        cancel,
		client:        redis.RedisDb,
	}

	// 启动指标收集协程
	mc.startCollectors()

	log.Println("性能监控指标收集器初始化成功")
	return mc, nil
}

// startCollectors 启动收集器
func (mc *MetricsCollector) startCollectors() {
	// 启动系统指标收集
	mc.wg.Add(1)
	go mc.collectSystemMetrics()

	// 启动数据库指标收集
	mc.wg.Add(1)
	go mc.collectDatabaseMetrics()

	// 启动Redis指标收集
	mc.wg.Add(1)
	go mc.collectRedisMetrics()

	// 启动指标持久化
	mc.wg.Add(1)
	go mc.persistMetrics()
}

// RecordHTTPRequest 记录HTTP请求
func (mc *MetricsCollector) RecordHTTPRequest(method, path string, statusCode int, responseTime time.Duration) {
	mc.httpMetrics.mutex.Lock()
	defer mc.httpMetrics.mutex.Unlock()

	// 更新总请求数
	atomic.AddInt64(&mc.httpMetrics.TotalRequests, 1)
	atomic.AddInt64(&mc.httpMetrics.requestCount, 1)

	// 更新响应时间统计
	atomic.AddInt64(&mc.httpMetrics.totalResponseTime, int64(responseTime))

	if responseTime > mc.httpMetrics.MaxResponseTime {
		mc.httpMetrics.MaxResponseTime = responseTime
	}

	if responseTime < mc.httpMetrics.MinResponseTime {
		mc.httpMetrics.MinResponseTime = responseTime
	}

	// 更新状态码统计
	mc.httpMetrics.StatusCodes[statusCode]++

	// 更新错误计数
	if statusCode >= 400 {
		atomic.AddInt64(&mc.httpMetrics.errorCount, 1)
	}

	// 更新路径统计
	pathKey := fmt.Sprintf("%s %s", method, path)
	pathStat, exists := mc.httpMetrics.PathStats[pathKey]
	if !exists {
		pathStat = &PathStat{}
		mc.httpMetrics.PathStats[pathKey] = pathStat
	}

	pathStat.mutex.Lock()
	pathStat.Count++
	pathStat.totalResponseTime += int64(responseTime)
	pathStat.AvgResponseTime = time.Duration(pathStat.totalResponseTime / pathStat.Count)

	if responseTime > pathStat.MaxResponseTime {
		pathStat.MaxResponseTime = responseTime
	}

	if statusCode >= 400 {
		pathStat.ErrorCount++
	}

	pathStat.LastAccess = time.Now()
	pathStat.mutex.Unlock()

	// 更新最后请求时间
	mc.httpMetrics.LastRequestTime = time.Now()

	// 计算平均响应时间和错误率
	if mc.httpMetrics.requestCount > 0 {
		mc.httpMetrics.AvgResponseTime = time.Duration(mc.httpMetrics.totalResponseTime / mc.httpMetrics.requestCount)
		mc.httpMetrics.ErrorRate = float64(mc.httpMetrics.errorCount) / float64(mc.httpMetrics.requestCount) * 100
	}
}

// collectSystemMetrics 收集系统指标
func (mc *MetricsCollector) collectSystemMetrics() {
	defer mc.wg.Done()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mc.updateSystemMetrics()
		case <-mc.ctx.Done():
			return
		}
	}
}

// updateSystemMetrics 更新系统指标
func (mc *MetricsCollector) updateSystemMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	mc.systemMetrics.MemoryUsage = int64(m.Alloc)
	mc.systemMetrics.MemoryTotal = int64(m.Sys)
	mc.systemMetrics.MemoryPercent = float64(m.Alloc) / float64(m.Sys) * 100
	mc.systemMetrics.GoroutineCount = runtime.NumGoroutine()
	mc.systemMetrics.GCCount = m.NumGC
	mc.systemMetrics.HeapSize = int64(m.HeapSys)
	mc.systemMetrics.HeapInUse = int64(m.HeapInuse)
	mc.systemMetrics.StackInUse = int64(m.StackInuse)
	mc.systemMetrics.LastUpdate = time.Now()

	// 计算GC暂停时间
	if m.NumGC > 0 {
		mc.systemMetrics.GCPauseTime = time.Duration(m.PauseNs[(m.NumGC+255)%256])
	}
}

// collectDatabaseMetrics 收集数据库指标
func (mc *MetricsCollector) collectDatabaseMetrics() {
	defer mc.wg.Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mc.updateDatabaseMetrics()
		case <-mc.ctx.Done():
			return
		}
	}
}

// updateDatabaseMetrics 更新数据库指标
func (mc *MetricsCollector) updateDatabaseMetrics() {
	// 这里需要根据实际的数据库连接池实现来获取指标
	// 示例实现
	mc.dbMetrics.LastUpdate = time.Now()
}

// collectRedisMetrics 收集Redis指标
func (mc *MetricsCollector) collectRedisMetrics() {
	defer mc.wg.Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mc.updateRedisMetrics()
		case <-mc.ctx.Done():
			return
		}
	}
}

// updateRedisMetrics 更新Redis指标
func (mc *MetricsCollector) updateRedisMetrics() {
	if mc.client == nil {
		return
	}

	_, err := mc.client.Info(mc.ctx, "memory", "stats", "clients").Result()
	if err != nil {
		log.Printf("获取Redis信息失败: %v", err)
		return
	}

	// 解析Redis INFO命令的输出
	// 这里简化处理，实际应该解析具体的字段
	mc.redisMetrics.LastUpdate = time.Now()
}

// persistMetrics 持久化指标
func (mc *MetricsCollector) persistMetrics() {
	defer mc.wg.Done()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mc.saveMetricsToRedis()
		case <-mc.ctx.Done():
			return
		}
	}
}

// saveMetricsToRedis 保存指标到Redis
func (mc *MetricsCollector) saveMetricsToRedis() {
	if mc.client == nil {
		return
	}

	metrics := mc.GetAllMetrics()
	data, err := json.Marshal(metrics)
	if err != nil {
		log.Printf("序列化指标数据失败: %v", err)
		return
	}

	key := fmt.Sprintf("metrics:%d", time.Now().Unix())
	err = mc.client.Set(mc.ctx, key, data, time.Hour*24).Err()
	if err != nil {
		log.Printf("保存指标到Redis失败: %v", err)
		return
	}

	// 保留最近24小时的指标数据
	pattern := "metrics:*"
	keys, err := mc.client.Keys(mc.ctx, pattern).Result()
	if err == nil && len(keys) > 1440 { // 24小时 * 60分钟
		// 删除最旧的指标
		oldestKeys := keys[:len(keys)-1440]
		mc.client.Del(mc.ctx, oldestKeys...)
	}
}

// SetCustomMetric 设置自定义指标
func (mc *MetricsCollector) SetCustomMetric(name string, value interface{}, metricType, description string) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	metric, exists := mc.customMetrics[name]
	if !exists {
		metric = &CustomMetric{
			Name:        name,
			Type:        metricType,
			Description: description,
		}
		mc.customMetrics[name] = metric
	}

	metric.mutex.Lock()
	metric.Value = value
	metric.LastUpdate = time.Now()
	metric.mutex.Unlock()
}

// IncrementCounter 递增计数器
func (mc *MetricsCollector) IncrementCounter(name string, delta int64) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	metric, exists := mc.customMetrics[name]
	if !exists {
		metric = &CustomMetric{
			Name:  name,
			Type:  "counter",
			Value: int64(0),
		}
		mc.customMetrics[name] = metric
	}

	metric.mutex.Lock()
	if currentValue, ok := metric.Value.(int64); ok {
		metric.Value = currentValue + delta
	} else {
		metric.Value = delta
	}
	metric.LastUpdate = time.Now()
	metric.mutex.Unlock()
}

// GetAllMetrics 获取所有指标
func (mc *MetricsCollector) GetAllMetrics() map[string]interface{} {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	// 计算QPS
	if time.Since(mc.httpMetrics.LastRequestTime) < time.Minute {
		mc.httpMetrics.RequestsPerSec = float64(mc.httpMetrics.requestCount) / time.Since(mc.httpMetrics.LastRequestTime).Seconds()
	}

	return map[string]interface{}{
		"http":      mc.httpMetrics,
		"system":    mc.systemMetrics,
		"db":        mc.dbMetrics,
		"redis":     mc.redisMetrics,
		"custom":    mc.customMetrics,
		"timestamp": time.Now(),
	}
}

// GetHTTPMetrics 获取HTTP指标
func (mc *MetricsCollector) GetHTTPMetrics() *HTTPMetrics {
	return mc.httpMetrics
}

// GetSystemMetrics 获取系统指标
func (mc *MetricsCollector) GetSystemMetrics() *SystemMetrics {
	return mc.systemMetrics
}

// Close 关闭指标收集器
func (mc *MetricsCollector) Close() error {
	mc.cancel()
	mc.wg.Wait()
	log.Println("性能监控指标收集器已关闭")
	return nil
}
