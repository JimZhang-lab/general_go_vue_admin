/*
 * @Author: JimZhang
 * @Date: 2025-07-29 17:50:00
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-29 17:50:00
 * @FilePath: /server/middleware/metrics_middleware.go
 * @Description: 性能监控中间件
 */

package middleware

import (
	"fmt"
	"net/http"
	"server/common/result"
	"server/pkg/metrics"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// MetricsConfig 监控配置
type MetricsConfig struct {
	// 是否启用监控
	Enabled bool
	// 跳过监控的路径
	SkipPaths []string
	// 是否记录详细路径统计
	DetailedPathStats bool
	// 路径分组规则（将动态路径归类）
	PathGroupRules map[string]string
}

// DefaultMetricsConfig 默认监控配置
var DefaultMetricsConfig = MetricsConfig{
	Enabled:           true,
	SkipPaths:         []string{"/health", "/metrics", "/favicon.ico"},
	DetailedPathStats: true,
	PathGroupRules: map[string]string{
		"/api/admin/*/edit":   "/api/admin/:id/edit",
		"/api/admin/*/delete": "/api/admin/:id/delete",
		"/api/user/*/profile": "/api/user/:id/profile",
	},
}

// MetricsMiddleware 性能监控中间件
func MetricsMiddleware(config ...MetricsConfig) gin.HandlerFunc {
	cfg := DefaultMetricsConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	if !cfg.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// 获取指标收集器
	metricsCollector := metrics.GetMetricsCollector()
	if metricsCollector == nil {
		panic("性能监控指标收集器未初始化")
	}

	return func(c *gin.Context) {
		// 检查是否跳过监控
		if shouldSkipMetrics(c.Request.URL.Path, cfg.SkipPaths) {
			c.Next()
			return
		}

		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// 应用路径分组规则
		if cfg.DetailedPathStats {
			path = applyPathGroupRules(path, cfg.PathGroupRules)
		}

		// 处理请求
		c.Next()

		// 计算响应时间
		responseTime := time.Since(start)
		statusCode := c.Writer.Status()

		// 记录HTTP指标
		metricsCollector.RecordHTTPRequest(method, path, statusCode, responseTime)

		// 添加性能指标到响应头
		c.Header("X-Response-Time", fmt.Sprintf("%.2fms", float64(responseTime.Nanoseconds())/1e6))
		c.Header("X-Request-ID", getRequestID(c))
	}
}

// HealthCheckHandler 健康检查处理器
func HealthCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		metricsCollector := metrics.GetMetricsCollector()
		if metricsCollector == nil {
			result.Failed(c, http.StatusInternalServerError, "监控系统未初始化")
			return
		}

		systemMetrics := metricsCollector.GetSystemMetrics()

		health := map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now(),
			"uptime":    time.Since(systemMetrics.LastUpdate),
			"system": map[string]interface{}{
				"memory_usage_mb": systemMetrics.MemoryUsage / 1024 / 1024,
				"memory_percent":  systemMetrics.MemoryPercent,
				"goroutines":      systemMetrics.GoroutineCount,
				"gc_count":        systemMetrics.GCCount,
			},
		}

		// 检查系统健康状态
		if systemMetrics.MemoryPercent > 90 {
			health["status"] = "warning"
			health["warnings"] = []string{"内存使用率过高"}
		}

		if systemMetrics.GoroutineCount > 10000 {
			health["status"] = "warning"
			if warnings, exists := health["warnings"]; exists {
				health["warnings"] = append(warnings.([]string), "协程数量过多")
			} else {
				health["warnings"] = []string{"协程数量过多"}
			}
		}

		result.Success(c, health)
	}
}

// MetricsHandler 指标查看处理器
func MetricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		metricsCollector := metrics.GetMetricsCollector()
		if metricsCollector == nil {
			result.Failed(c, http.StatusInternalServerError, "监控系统未初始化")
			return
		}

		// 获取查询参数
		metricType := c.Query("type")
		format := c.DefaultQuery("format", "json")

		var data interface{}

		switch metricType {
		case "http":
			data = metricsCollector.GetHTTPMetrics()
		case "system":
			data = metricsCollector.GetSystemMetrics()
		case "summary":
			data = getSummaryMetrics(metricsCollector)
		default:
			data = metricsCollector.GetAllMetrics()
		}

		if format == "prometheus" {
			// 返回Prometheus格式的指标
			c.Header("Content-Type", "text/plain")
			c.String(http.StatusOK, convertToPrometheusFormat(data))
			return
		}

		result.Success(c, data)
	}
}

// QPSHandler QPS查看处理器
func QPSHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		metricsCollector := metrics.GetMetricsCollector()
		if metricsCollector == nil {
			result.Failed(c, http.StatusInternalServerError, "监控系统未初始化")
			return
		}

		httpMetrics := metricsCollector.GetHTTPMetrics()

		// 计算实时QPS
		duration := c.DefaultQuery("duration", "60") // 默认60秒
		durationInt, err := strconv.Atoi(duration)
		if err != nil {
			durationInt = 60
		}

		qpsData := map[string]interface{}{
			"current_qps":       httpMetrics.RequestsPerSec,
			"total_requests":    httpMetrics.TotalRequests,
			"avg_response_time": httpMetrics.AvgResponseTime.Milliseconds(),
			"max_response_time": httpMetrics.MaxResponseTime.Milliseconds(),
			"min_response_time": httpMetrics.MinResponseTime.Milliseconds(),
			"error_rate":        httpMetrics.ErrorRate,
			"status_codes":      httpMetrics.StatusCodes,
			"duration_seconds":  durationInt,
			"timestamp":         time.Now(),
		}

		// 添加热门路径统计
		topPaths := getTopPaths(httpMetrics.PathStats, 10)
		qpsData["top_paths"] = topPaths

		result.Success(c, qpsData)
	}
}

// PerformanceReportHandler 性能报告处理器
func PerformanceReportHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		metricsCollector := metrics.GetMetricsCollector()
		if metricsCollector == nil {
			result.Failed(c, http.StatusInternalServerError, "监控系统未初始化")
			return
		}

		allMetrics := metricsCollector.GetAllMetrics()
		httpMetrics := metricsCollector.GetHTTPMetrics()
		systemMetrics := metricsCollector.GetSystemMetrics()

		report := map[string]interface{}{
			"report_time": time.Now(),
			"summary": map[string]interface{}{
				"total_requests":    httpMetrics.TotalRequests,
				"current_qps":       httpMetrics.RequestsPerSec,
				"avg_response_time": httpMetrics.AvgResponseTime.Milliseconds(),
				"error_rate":        httpMetrics.ErrorRate,
				"memory_usage_mb":   systemMetrics.MemoryUsage / 1024 / 1024,
				"memory_percent":    systemMetrics.MemoryPercent,
				"goroutine_count":   systemMetrics.GoroutineCount,
			},
			"performance_analysis": analyzePerformance(httpMetrics, systemMetrics),
			"recommendations":      generateRecommendations(httpMetrics, systemMetrics),
			"detailed_metrics":     allMetrics,
		}

		result.Success(c, report)
	}
}

// 辅助函数

// shouldSkipMetrics 检查是否应该跳过监控
func shouldSkipMetrics(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// applyPathGroupRules 应用路径分组规则
func applyPathGroupRules(path string, rules map[string]string) string {
	for pattern, replacement := range rules {
		// 简单的通配符匹配
		if strings.Contains(pattern, "*") {
			prefix := strings.Split(pattern, "*")[0]
			suffix := strings.Split(pattern, "*")[1]
			if strings.HasPrefix(path, prefix) && strings.HasSuffix(path, suffix) {
				return replacement
			}
		}
	}
	return path
}

// getRequestID 获取请求ID
func getRequestID(c *gin.Context) string {
	if requestID := c.GetHeader("X-Request-ID"); requestID != "" {
		return requestID
	}

	// 生成简单的请求ID
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// getSummaryMetrics 获取摘要指标
func getSummaryMetrics(collector *metrics.MetricsCollector) map[string]interface{} {
	httpMetrics := collector.GetHTTPMetrics()
	systemMetrics := collector.GetSystemMetrics()

	return map[string]interface{}{
		"qps":               httpMetrics.RequestsPerSec,
		"total_requests":    httpMetrics.TotalRequests,
		"avg_response_time": httpMetrics.AvgResponseTime.Milliseconds(),
		"error_rate":        httpMetrics.ErrorRate,
		"memory_usage_mb":   systemMetrics.MemoryUsage / 1024 / 1024,
		"memory_percent":    systemMetrics.MemoryPercent,
		"goroutine_count":   systemMetrics.GoroutineCount,
		"gc_count":          systemMetrics.GCCount,
		"timestamp":         time.Now(),
	}
}

// getTopPaths 获取热门路径
func getTopPaths(pathStats map[string]*metrics.PathStat, limit int) []map[string]interface{} {
	type pathInfo struct {
		Path  string
		Count int64
		Stat  *metrics.PathStat
	}

	var paths []pathInfo
	for path, stat := range pathStats {
		paths = append(paths, pathInfo{
			Path:  path,
			Count: stat.Count,
			Stat:  stat,
		})
	}

	// 简单排序（按请求数量）
	for i := 0; i < len(paths)-1; i++ {
		for j := i + 1; j < len(paths); j++ {
			if paths[i].Count < paths[j].Count {
				paths[i], paths[j] = paths[j], paths[i]
			}
		}
	}

	// 取前N个
	if len(paths) > limit {
		paths = paths[:limit]
	}

	result := make([]map[string]interface{}, len(paths))
	for i, path := range paths {
		result[i] = map[string]interface{}{
			"path":              path.Path,
			"count":             path.Count,
			"avg_response_time": path.Stat.AvgResponseTime.Milliseconds(),
			"max_response_time": path.Stat.MaxResponseTime.Milliseconds(),
			"error_count":       path.Stat.ErrorCount,
			"last_access":       path.Stat.LastAccess,
		}
	}

	return result
}

// analyzePerformance 分析性能
func analyzePerformance(httpMetrics *metrics.HTTPMetrics, systemMetrics *metrics.SystemMetrics) map[string]interface{} {
	analysis := map[string]interface{}{
		"status": "good",
		"issues": []string{},
	}

	issues := []string{}

	// 检查响应时间
	if httpMetrics.AvgResponseTime > time.Second {
		issues = append(issues, "平均响应时间过长")
		analysis["status"] = "warning"
	}

	// 检查错误率
	if httpMetrics.ErrorRate > 5.0 {
		issues = append(issues, "错误率过高")
		analysis["status"] = "critical"
	}

	// 检查内存使用
	if systemMetrics.MemoryPercent > 80 {
		issues = append(issues, "内存使用率过高")
		analysis["status"] = "warning"
	}

	// 检查协程数量
	if systemMetrics.GoroutineCount > 5000 {
		issues = append(issues, "协程数量过多")
		analysis["status"] = "warning"
	}

	analysis["issues"] = issues
	return analysis
}

// generateRecommendations 生成建议
func generateRecommendations(httpMetrics *metrics.HTTPMetrics, systemMetrics *metrics.SystemMetrics) []string {
	recommendations := []string{}

	if httpMetrics.AvgResponseTime > time.Second {
		recommendations = append(recommendations, "考虑优化数据库查询或添加缓存")
	}

	if httpMetrics.ErrorRate > 5.0 {
		recommendations = append(recommendations, "检查错误日志，修复导致错误的问题")
	}

	if systemMetrics.MemoryPercent > 80 {
		recommendations = append(recommendations, "考虑增加内存或优化内存使用")
	}

	if systemMetrics.GoroutineCount > 5000 {
		recommendations = append(recommendations, "检查是否有协程泄漏")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "系统运行良好，继续保持")
	}

	return recommendations
}

// convertToPrometheusFormat 转换为Prometheus格式
func convertToPrometheusFormat(data interface{}) string {
	// 简单的Prometheus格式转换
	// 实际实现应该更完整
	return fmt.Sprintf("# HELP http_requests_total Total HTTP requests\n# TYPE http_requests_total counter\nhttp_requests_total %v\n", data)
}
