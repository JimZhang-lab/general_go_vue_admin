/*
 * @Author: JimZhang
 * @Date: 2025-07-24 12:40:50
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-26 00:35:11
 * @FilePath: /server/router/routers.go
 * @Description: 访问接口
 *
 */
package router

import (
	"net/http"
	"server/api/controller"
	"server/common/config"
	"server/middleware"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type IFnRegisterRouter = func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup)

var (
	gfnRoutes []IFnRegisterRouter
)

func RegisterRouter(fn IFnRegisterRouter) {
	if fn == nil {
		return
	}
	gfnRoutes = append(gfnRoutes, fn)
}

func InitRouter() *gin.Engine {
	router := gin.New()

	// 创建并发管理器
	concurrencyManager := middleware.NewConcurrencyManager(middleware.DefaultConcurrencyConfig)

	// 基础中间件
	router.Use(gin.Recovery())
	router.Use(middleware.Cors())

	// 高并发处理中间件
	router.Use(concurrencyManager.RateLimitMiddleware())
	router.Use(concurrencyManager.ConcurrencyLimitMiddleware())
	router.Use(concurrencyManager.TimeoutMiddleware())
	router.Use(concurrencyManager.CircuitBreakerMiddleware())

	// 智能绑定中间件
	router.Use(middleware.SmartBindMiddleware(middleware.DefaultBindingConfig))

	// 日志中间件
	router.Use(middleware.Logger())

	// 静态文件服务
	router.StaticFS(config.Config.ImageSettings.UploadDir, http.Dir(config.Config.ImageSettings.UploadDir))

	regsiterAdminPlatformRouter(router)
	return router
}

func regsiterAdminPlatformRouter(router *gin.Engine) {
	// 不需要认证的接口
	router.GET("/api/captcha", controller.Captcha)
	router.POST("/api/login", controller.Login)

	// 需要JWT认证的登出接口
	router.POST("/api/logout", middleware.AuthMiddleware(), controller.Logout)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// router.GET("/swagger", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 需要JWT认证的接口
	jwt := router.Group("/api", middleware.AuthMiddleware(), middleware.LogMiddleware())
	{
		jwt.POST("/post/add", controller.CreateSysPost)
		jwt.GET("/post/list", controller.GetSysPostList)
		jwt.GET("/post/info", controller.GetSysPostById)
		jwt.PUT("/post/update", controller.UpdateSysPost)
		jwt.DELETE("/post/delete", controller.DeleteSysPostById)
		jwt.DELETE("/post/batch/delete", controller.BatchDeleteSysPost)
		jwt.PUT("/post/updateStatus", controller.UpdateSysPostStatus)
		jwt.GET("/post/vo/list", controller.QuerySysPostVoList)
		jwt.GET("/dept/list", controller.GetSysDeptList)
		jwt.POST("/dept/add", controller.CreateSysDept)
		jwt.GET("/dept/info", controller.GetSysDeptById)
		jwt.PUT("/dept/update", controller.UpdateSysDept)
		jwt.DELETE("/dept/delete", controller.DeleteSysDeptById)
		jwt.GET("/dept/vo/list", controller.QuerySysDeptVoList)
		jwt.POST("/menu/add", controller.CreateSysMenu)
		jwt.GET("/menu/vo/list", controller.QuerySysMenuVoList)
		jwt.GET("/menu/info", controller.GetSysMenu)
		jwt.PUT("/menu/update", controller.UpdateSysMenu)
		jwt.DELETE("/menu/delete", controller.DeleteSysMenu)
		jwt.GET("/menu/list", controller.GetSysMenuList)
		jwt.POST("/role/add", controller.CreateSysRole)
		jwt.GET("/role/info", controller.GetSysRoleById)
		jwt.PUT("/role/update", controller.UpdateSysRole)
		jwt.DELETE("/role/delete", controller.DeleteSysRoleById)
		jwt.PUT("/role/updateStatus", controller.UpdateSysRoleStatus)
		jwt.GET("/role/list", controller.GetSysRoleList)
		jwt.GET("/role/vo/list", controller.QuerySysRoleVoList)
		jwt.GET("/role/vo/idList", controller.QueryRoleMenuIdList)
		jwt.PUT("/role/assignPermissions", controller.AssignPermissions)
		jwt.POST("/admin/add", controller.CreateSysAdmin)
		jwt.GET("/admin/info", controller.GetSysAdminInfo)
		jwt.PUT("/admin/update", controller.UpdateSysAdmin)
		jwt.DELETE("/admin/delete", controller.DeleteSysAdminById)
		jwt.PUT("/admin/updateStatus", controller.UpdateSysAdminStatus)
		jwt.PUT("/admin/updatePassword", controller.ResetSysAdminPassword)
		jwt.GET("/admin/list", controller.GetSysAdminList)
		jwt.POST("/upload", controller.Upload)
		jwt.PUT("/admin/updatePersonal", controller.UpdatePersonal)
		jwt.PUT("/admin/updatePersonalPassword", controller.UpdatePersonalPassword)
		jwt.GET("/sysLoginInfo/list", controller.GetSysLoginInfoList)
		jwt.DELETE("/sysLoginInfo/batch/delete", controller.BatchDeleteSysLoginInfo)
		jwt.DELETE("/sysLoginInfo/delete", controller.DeleteSysLoginInfoById)
		jwt.DELETE("/sysLoginInfo/clean", controller.CleanSysLoginInfo)
		jwt.GET("/sysOperationLog/list", controller.GetSysOperationLogList)
		jwt.DELETE("/sysOperationLog/delete", controller.DeleteSysOperationLogById)
		jwt.DELETE("/sysOperationLog/batch/delete", controller.BatchDeleteSysOperationLog)
		jwt.DELETE("/sysOperationLog/clean", controller.CleanSysOperationLog)
	}
}
