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
	router.POST("/api/register", controller.Register)

	// 需要JWT认证的登出接口
	router.POST("/api/logout", middleware.AuthMiddleware(), controller.Logout)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 需要JWT认证的接口
	jwtGroup := router.Group("/api", middleware.AuthMiddleware(), middleware.LogMiddleware())
	{
		jwtGroup.POST("/post/add", middleware.PermissionMiddleware("system:admin:add"), controller.CreateSysPost)
		jwtGroup.GET("/post/list", middleware.PermissionMiddleware("system:admin:list"), controller.GetSysPostList)
		jwtGroup.GET("/post/info", middleware.PermissionMiddleware("system:admin:list"), controller.GetSysPostById)
		jwtGroup.PUT("/post/update", middleware.PermissionMiddleware("system:admin:update"), controller.UpdateSysPost)
		jwtGroup.DELETE("/post/delete", middleware.PermissionMiddleware("system:admin:delete"), controller.DeleteSysPostById)
		jwtGroup.DELETE("/post/batch/delete", middleware.PermissionMiddleware("system:admin:delete"), controller.BatchDeleteSysPost)
		jwtGroup.PUT("/post/updateStatus", middleware.PermissionMiddleware("system:admin:updateStatus"), controller.UpdateSysPostStatus)
		jwtGroup.GET("/post/vo/list", middleware.PermissionMiddleware("system:admin:list"), controller.QuerySysPostVoList)

		jwtGroup.GET("/dept/list", middleware.PermissionMiddleware("system:admin:list"), controller.GetSysDeptList)
		jwtGroup.POST("/dept/add", middleware.PermissionMiddleware("system:admin:add"), controller.CreateSysDept)
		jwtGroup.GET("/dept/info", middleware.PermissionMiddleware("system:admin:list"), controller.GetSysDeptById)
		jwtGroup.PUT("/dept/update", middleware.PermissionMiddleware("system:admin:update"), controller.UpdateSysDept)
		jwtGroup.DELETE("/dept/delete", middleware.PermissionMiddleware("system:admin:delete"), controller.DeleteSysDeptById)
		jwtGroup.GET("/dept/vo/list", middleware.PermissionMiddleware("system:admin:list"), controller.QuerySysDeptVoList)

		jwtGroup.POST("/menu/add", middleware.PermissionMiddleware("system:menu:add"), controller.CreateSysMenu)
		jwtGroup.GET("/menu/vo/list", middleware.PermissionMiddleware("system:menu:list"), controller.QuerySysMenuVoList)
		jwtGroup.GET("/menu/info", middleware.PermissionMiddleware("system:menu:list"), controller.GetSysMenu)
		jwtGroup.PUT("/menu/update", middleware.PermissionMiddleware("system:menu:update"), controller.UpdateSysMenu)
		jwtGroup.DELETE("/menu/delete", middleware.PermissionMiddleware("system:menu:delete"), controller.DeleteSysMenu)
		jwtGroup.GET("/menu/list", middleware.PermissionMiddleware("system:menu:list"), controller.GetSysMenuList)

		jwtGroup.POST("/role/add", middleware.PermissionMiddleware("system:role:add"), controller.CreateSysRole)
		jwtGroup.GET("/role/info", middleware.PermissionMiddleware("system:role:list"), controller.GetSysRoleById)
		jwtGroup.PUT("/role/update", middleware.PermissionMiddleware("system:role:update"), controller.UpdateSysRole)
		jwtGroup.DELETE("/role/delete", middleware.PermissionMiddleware("system:role:delete"), controller.DeleteSysRoleById)
		jwtGroup.PUT("/role/updateStatus", middleware.PermissionMiddleware("system:role:update"), controller.UpdateSysRoleStatus)
		jwtGroup.GET("/role/list", middleware.PermissionMiddleware("system:role:list"), controller.GetSysRoleList)
		jwtGroup.GET("/role/vo/list", middleware.PermissionMiddleware("system:role:list"), controller.QuerySysRoleVoList)
		jwtGroup.GET("/role/vo/idList", middleware.PermissionMiddleware("system:role:list"), controller.QueryRoleMenuIdList)
		jwtGroup.PUT("/role/assignPermissions", middleware.PermissionMiddleware("system:role:assign"), controller.AssignPermissions)

		jwtGroup.POST("/admin/add", middleware.PermissionMiddleware("system:admin:add"), controller.CreateSysAdmin)
		jwtGroup.GET("/admin/info", middleware.PermissionMiddleware("system:profile:view", "system:admin:list"), controller.GetSysAdminInfo)
		jwtGroup.PUT("/admin/update", middleware.PermissionMiddleware("system:admin:update"), controller.UpdateSysAdmin)
		jwtGroup.DELETE("/admin/delete", middleware.PermissionMiddleware("system:admin:delete"), controller.DeleteSysAdminById)
		jwtGroup.PUT("/admin/updateStatus", middleware.PermissionMiddleware("system:admin:updateStatus"), controller.UpdateSysAdminStatus)
		jwtGroup.PUT("/admin/updatePassword", middleware.PermissionMiddleware("system:admin:resetPassword"), controller.ResetSysAdminPassword)
		jwtGroup.GET("/admin/list", middleware.PermissionMiddleware("system:admin:list"), controller.GetSysAdminList)
		jwtGroup.POST("/upload", controller.Upload)
		jwtGroup.PUT("/admin/updatePersonal", controller.UpdatePersonal)
		jwtGroup.PUT("/admin/updatePersonalPassword", controller.UpdatePersonalPassword)

		jwtGroup.GET("/sysLoginInfo/list", middleware.PermissionMiddleware("system:log:list"), controller.GetSysLoginInfoList)
		jwtGroup.DELETE("/sysLoginInfo/batch/delete", middleware.PermissionMiddleware("system:log:delete"), controller.BatchDeleteSysLoginInfo)
		jwtGroup.DELETE("/sysLoginInfo/delete", middleware.PermissionMiddleware("system:log:delete"), controller.DeleteSysLoginInfoById)
		jwtGroup.DELETE("/sysLoginInfo/clean", middleware.PermissionMiddleware("system:log:clean"), controller.CleanSysLoginInfo)

		jwtGroup.GET("/sysOperationLog/list", middleware.PermissionMiddleware("system:log:list"), controller.GetSysOperationLogList)
		jwtGroup.DELETE("/sysOperationLog/delete", middleware.PermissionMiddleware("system:log:delete"), controller.DeleteSysOperationLogById)
		jwtGroup.DELETE("/sysOperationLog/batch/delete", middleware.PermissionMiddleware("system:log:delete"), controller.BatchDeleteSysOperationLog)
		jwtGroup.DELETE("/sysOperationLog/clean", middleware.PermissionMiddleware("system:log:clean"), controller.CleanSysOperationLog)

		jwtGroup.GET("/setting/list", middleware.PermissionMiddleware("system:setting:list"), controller.GetSysSettingList)
		jwtGroup.PUT("/setting/batch/update", middleware.PermissionMiddleware("system:setting:update"), controller.BatchUpdateSysSettings)

		jwtGroup.GET("/notice/list", middleware.PermissionMiddleware("system:notice:list"), controller.GetSysNoticeList)
		jwtGroup.GET("/notice/current", controller.GetCurrentSysNoticeList)
		jwtGroup.POST("/notice/add", middleware.PermissionMiddleware("system:notice:add"), controller.CreateSysNotice)
		jwtGroup.PUT("/notice/update", middleware.PermissionMiddleware("system:notice:update"), controller.UpdateSysNotice)
		jwtGroup.PUT("/notice/updateStatus", middleware.PermissionMiddleware("system:notice:publish"), controller.UpdateSysNoticeStatus)
		jwtGroup.DELETE("/notice/delete", middleware.PermissionMiddleware("system:notice:delete"), controller.DeleteSysNoticeById)
		jwtGroup.DELETE("/notice/batch/delete", middleware.PermissionMiddleware("system:notice:delete"), controller.BatchDeleteSysNotice)
		jwtGroup.PUT("/notice/read", controller.MarkSysNoticeRead)
		jwtGroup.PUT("/notice/readAll", controller.MarkAllSysNoticeRead)
	}
}
