/*
 * @Author: JimZhang
 * @Date: 2025-07-24 12:31:33
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-08-11 00:00:00
 * @FilePath: /server/middleware/logger.go
 * @Description: 日志中间件（Zap）
 */
package middleware

import (
	"io"
	"server/pkg/log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	logger := log.Log()
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime) / time.Millisecond

		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		header := c.Request.Header
		proto := c.Request.Proto
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		err := c.Err()
		body, _ := io.ReadAll(c.Request.Body)

		logger.Infow("Request processed",
			"status_code", statusCode,
			"latency_time_ms", latencyTime,
			"client_ip", clientIP,
			"req_method", reqMethod,
			"req_uri", reqUri,
			"header", header,
			"proto", proto,
			"err", err,
			"body", string(body),
		)
	}
}
