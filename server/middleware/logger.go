/*
 * @Author: JimZhang
 * @Date: 2025-07-24 12:31:33
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-24 23:02:58
 * @FilePath: /server/middleware/logger.go
 * @Description: 日志中间件
 *
 */
package middleware

import (
	"io"
	"server/pkg/log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logger() gin.HandlerFunc {
	logger := log.Log()
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime) / time.Millisecond
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		header := c.Request.Header
		proto := c.Request.Proto
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		err := c.Err()
		body, _ := io.ReadAll(c.Request.Body)
		logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUri,
			"header":       header,
			"proto":        proto,
			"err":          err,
			"body":         string(body),
		}).Info("Request processed")
	}
}
