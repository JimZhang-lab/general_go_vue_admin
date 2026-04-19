package controller

import (
	"fmt"
	"runtime"
	"server/common/result"
	"time"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

// 获取服务器运行状态监控信息
// @Summary 获取服务器监控信息
// @Produce json
// @router /api/monitor/server [get]
// @Security ApiKeyAuth
func GetServerMonitorInfo(c *gin.Context) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	osInfo := map[string]interface{}{
		"go_version": runtime.Version(),
		"os":         runtime.GOOS,
		"arch":       runtime.GOARCH,
		"num_cpu":    runtime.NumCPU(),
		"num_goroutine": runtime.NumGoroutine(),
		"uptime":     time.Since(startTime).String(),
	}

	memInfo := map[string]interface{}{
		"alloc":      fmt.Sprintf("%.2f MB", float64(memStats.Alloc)/1024/1024),
		"total_alloc": fmt.Sprintf("%.2f MB", float64(memStats.TotalAlloc)/1024/1024),
		"sys":        fmt.Sprintf("%.2f MB", float64(memStats.Sys)/1024/1024),
		"num_gc":     memStats.NumGC,
	}

	result.Success(c, map[string]interface{}{
		"os":  osInfo,
		"mem": memInfo,
	})
}
