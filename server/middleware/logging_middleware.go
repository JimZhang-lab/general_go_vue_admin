/*
 * @Author: JimZhang
 * @Date: 2025-07-29 17:30:00
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-29 17:30:00
 * @FilePath: /server/middleware/logging_middleware.go
 * @Description: 异步日志中间件
 */

package middleware

import (
	"bytes"
	"fmt"
	"io"
	"server/pkg/logger"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingConfig 日志配置
type LoggingConfig struct {
	// 是否启用请求日志
	EnableRequestLog bool
	// 是否启用响应日志
	EnableResponseLog bool
	// 是否记录请求体
	LogRequestBody bool
	// 是否记录响应体
	LogResponseBody bool
	// 跳过日志的路径
	SkipPaths []string
	// 慢请求阈值
	SlowThreshold time.Duration
	// 最大请求体大小（字节）
	MaxBodySize int64
}

// DefaultLoggingConfig 默认日志配置
var DefaultLoggingConfig = LoggingConfig{
	EnableRequestLog:  true,
	EnableResponseLog: true,
	LogRequestBody:    false, // 默认不记录请求体（可能包含敏感信息）
	LogResponseBody:   false, // 默认不记录响应体（可能很大）
	SkipPaths:         []string{"/health", "/metrics", "/favicon.ico"},
	SlowThreshold:     time.Second * 2,
	MaxBodySize:       1024 * 1024, // 1MB
}

// responseWriter 响应写入器包装
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// AsyncLoggingMiddleware 异步日志中间件
func AsyncLoggingMiddleware(config ...LoggingConfig) gin.HandlerFunc {
	cfg := DefaultLoggingConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	// 获取异步日志处理器
	asyncLogger := logger.GetAsyncLogger()
	if asyncLogger == nil {
		panic("异步日志处理器未初始化")
	}

	return func(c *gin.Context) {
		// 检查是否跳过日志
		if shouldSkipLogging(c.Request.URL.Path, cfg.SkipPaths) {
			c.Next()
			return
		}

		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// 简化User-Agent解析
		browser := "Unknown"
		os := "Unknown"
		if strings.Contains(userAgent, "Chrome") {
			browser = "Chrome"
		} else if strings.Contains(userAgent, "Firefox") {
			browser = "Firefox"
		} else if strings.Contains(userAgent, "Safari") {
			browser = "Safari"
		}

		if strings.Contains(userAgent, "Windows") {
			os = "Windows"
		} else if strings.Contains(userAgent, "Mac") {
			os = "macOS"
		} else if strings.Contains(userAgent, "Linux") {
			os = "Linux"
		}

		// 读取请求体（如果需要）
		var requestBody string
		if cfg.LogRequestBody && c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil && len(bodyBytes) <= int(cfg.MaxBodySize) {
				requestBody = string(bodyBytes)
				// 重新设置请求体
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// 包装响应写入器（如果需要记录响应体）
		var responseBody string
		if cfg.LogResponseBody {
			writer := &responseWriter{
				ResponseWriter: c.Writer,
				body:           bytes.NewBuffer(nil),
			}
			c.Writer = writer
			defer func() {
				if writer.body.Len() <= int(cfg.MaxBodySize) {
					responseBody = writer.body.String()
				}
			}()
		}

		// 处理请求
		c.Next()

		// 计算处理时间
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// 构建日志字段
		fields := map[string]interface{}{
			"method":      method,
			"path":        path,
			"status_code": statusCode,
			"latency_ms":  latency.Milliseconds(),
			"client_ip":   clientIP,
			"user_agent":  userAgent,
			"browser":     browser,
			"os":          os,
		}

		// 添加请求体（如果启用）
		if cfg.LogRequestBody && requestBody != "" {
			fields["request_body"] = requestBody
		}

		// 添加响应体（如果启用）
		if cfg.LogResponseBody && responseBody != "" {
			fields["response_body"] = responseBody
		}

		// 添加错误信息（如果有）
		if len(c.Errors) > 0 {
			fields["errors"] = c.Errors.String()
		}

		// 记录请求日志
		if cfg.EnableRequestLog {
			logLevel := getLogLevel(statusCode, latency, cfg.SlowThreshold)
			message := fmt.Sprintf("%s %s %d %v", method, path, statusCode, latency)

			asyncLogger.Log(logLevel, message, fields)
		}

		// 记录操作日志（如果是登录用户）
		if userID, exists := c.Get("user_id"); exists {
			if username, ok := c.Get("username"); ok {
				asyncLogger.LogOperation(
					userID.(int),
					username.(string),
					method,
					path,
					clientIP,
					userAgent,
					map[string]interface{}{
						"status_code": statusCode,
						"latency_ms":  latency.Milliseconds(),
					},
				)
			}
		}
	}
}

// LoginLoggingMiddleware 登录日志中间件
func LoginLoggingMiddleware() gin.HandlerFunc {
	asyncLogger := logger.GetAsyncLogger()
	if asyncLogger == nil {
		panic("异步日志处理器未初始化")
	}

	return func(c *gin.Context) {
		// 只处理登录相关的路径
		if !isLoginPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// 简化User-Agent解析
		browser := "Unknown"
		os := "Unknown"
		if strings.Contains(userAgent, "Chrome") {
			browser = "Chrome"
		} else if strings.Contains(userAgent, "Firefox") {
			browser = "Firefox"
		} else if strings.Contains(userAgent, "Safari") {
			browser = "Safari"
		}

		if strings.Contains(userAgent, "Windows") {
			os = "Windows"
		} else if strings.Contains(userAgent, "Mac") {
			os = "macOS"
		} else if strings.Contains(userAgent, "Linux") {
			os = "Linux"
		}

		// 处理请求
		c.Next()

		// 获取登录结果
		username := getUsername(c)
		status := getLoginStatus(c)
		message := getLoginMessage(c)
		location := getLocation(clientIP) // 可以集成IP地理位置服务

		// 记录登录日志
		asyncLogger.LogLogin(username, clientIP, location, browser, os, status, message)
	}
}

// OperationLoggingMiddleware 操作日志中间件
func OperationLoggingMiddleware() gin.HandlerFunc {
	asyncLogger := logger.GetAsyncLogger()
	if asyncLogger == nil {
		panic("异步日志处理器未初始化")
	}

	return func(c *gin.Context) {
		// 只记录需要记录的操作
		if !shouldLogOperation(c.Request.Method, c.Request.URL.Path) {
			c.Next()
			return
		}

		// 检查用户是否已登录
		userID, userExists := c.Get("user_id")
		username, nameExists := c.Get("username")

		if !userExists || !nameExists {
			c.Next()
			return
		}

		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// 处理请求
		c.Next()

		// 记录操作日志
		asyncLogger.LogOperation(
			userID.(int),
			username.(string),
			c.Request.Method,
			c.Request.URL.Path,
			clientIP,
			userAgent,
			map[string]interface{}{
				"status_code": c.Writer.Status(),
			},
		)
	}
}

// ErrorLoggingMiddleware 错误日志中间件
func ErrorLoggingMiddleware() gin.HandlerFunc {
	asyncLogger := logger.GetAsyncLogger()
	if asyncLogger == nil {
		panic("异步日志处理器未初始化")
	}

	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				fields := map[string]interface{}{
					"method":     c.Request.Method,
					"path":       c.Request.URL.Path,
					"client_ip":  c.ClientIP(),
					"user_agent": c.Request.UserAgent(),
					"error_type": err.Type,
				}

				// 添加用户信息（如果有）
				if userID, exists := c.Get("user_id"); exists {
					fields["user_id"] = userID
				}
				if username, exists := c.Get("username"); exists {
					fields["username"] = username
				}

				asyncLogger.Error(err.Error(), fields)
			}
		}

		// 检查HTTP状态码
		statusCode := c.Writer.Status()
		if statusCode >= 500 {
			fields := map[string]interface{}{
				"method":      c.Request.Method,
				"path":        c.Request.URL.Path,
				"status_code": statusCode,
				"client_ip":   c.ClientIP(),
				"user_agent":  c.Request.UserAgent(),
			}

			asyncLogger.Error(fmt.Sprintf("HTTP %d Error", statusCode), fields)
		}
	}
}

// 辅助函数

// shouldSkipLogging 检查是否应该跳过日志
func shouldSkipLogging(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// getLogLevel 根据状态码和延迟获取日志级别
func getLogLevel(statusCode int, latency time.Duration, slowThreshold time.Duration) logger.LogLevel {
	if statusCode >= 500 {
		return logger.ErrorLevel
	}
	if statusCode >= 400 {
		return logger.WarnLevel
	}
	if latency > slowThreshold {
		return logger.WarnLevel
	}
	return logger.InfoLevel
}

// isLoginPath 检查是否是登录路径
func isLoginPath(path string) bool {
	loginPaths := []string{"/login", "/logout", "/auth/login", "/api/login"}
	for _, loginPath := range loginPaths {
		if strings.Contains(path, loginPath) {
			return true
		}
	}
	return false
}

// shouldLogOperation 检查是否应该记录操作日志
func shouldLogOperation(method, path string) bool {
	// 只记录修改操作
	if method == "GET" || method == "HEAD" || method == "OPTIONS" {
		return false
	}

	// 跳过某些路径
	skipPaths := []string{"/health", "/metrics", "/ping"}
	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return false
		}
	}

	return true
}

// getUsername 从上下文获取用户名
func getUsername(c *gin.Context) string {
	if username, exists := c.Get("username"); exists {
		return username.(string)
	}

	// 尝试从请求中获取
	if username := c.PostForm("username"); username != "" {
		return username
	}

	return "unknown"
}

// getLoginStatus 获取登录状态
func getLoginStatus(c *gin.Context) int {
	statusCode := c.Writer.Status()
	if statusCode == 200 {
		return 1 // 成功
	}
	return 2 // 失败
}

// getLoginMessage 获取登录消息
func getLoginMessage(c *gin.Context) string {
	if len(c.Errors) > 0 {
		return c.Errors.String()
	}

	statusCode := c.Writer.Status()
	if statusCode == 200 {
		return "登录成功"
	}

	return "登录失败"
}

// getLocation 获取IP地理位置（简单实现）
func getLocation(ip string) string {
	// 这里可以集成第三方IP地理位置服务
	// 目前返回简单的本地/外网判断
	if strings.HasPrefix(ip, "127.") || strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "10.") {
		return "内网"
	}
	return "外网"
}
