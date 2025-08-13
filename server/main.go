/*
 * @Author: JimZhang
 * @Date: 2025-07-24 10:43:16
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-26 00:53:58
 * @FilePath: /go-vue-general-admin/server/main.go
 * @Description:
 *
 */
package main

import (
	"context"
	"net/http"
	"os/signal"
	"server/common/config"
	"server/common/result"
	_ "server/docs"
	"server/pkg/db"
	"server/pkg/log"
	"server/pkg/redis"
	"server/router"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// @title 通用后台管理系统
// @version 1.0
// @description 后台管理系统API接口文档
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 加载日志
	logger := log.Log()
	gin.SetMode(config.Config.Server.Model)
	router := router.InitRouter()
	setupAddress := config.Config.Server.Host + ":" + strconv.Itoa(config.Config.Server.Port)
	srv := &http.Server{
		Addr:    setupAddress,
		Handler: router,
	}

	// 启动服务
	go func() {
		logger.Infof("Listening and serving HTTP on %s", setupAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Infof("listen: %s", err.Error())
			return
		}
	}()

	<-ctx.Done()

	logger.Info("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Infof("Server Shutdown: %v", err.Error())
		return
	}
	logger.Info("Server exiting")
}

func init() {
	// 初始化数据库
	db.SetupDBLink()
	// 初始化redis
	redis.SetupRedisDb()
	// 初始化API状态码
	result.Init()
}
