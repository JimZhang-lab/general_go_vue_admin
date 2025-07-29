/*
 * @Author: JimZhang
 * @Date: 2025-07-25 11:54:33
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-25 12:04:22
 * @FilePath: /server/router/sysAdminRouter.go
 * @Description:
 *
 */
package router

import (
	"server/api/controller"

	"github.com/gin-gonic/gin"
)

func InitAdminRouter(rgAdmin *gin.RouterGroup, rgAuth *gin.RouterGroup) {
	RegisterRouter(func(rgAdmin *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		rgAdminUser := rgAdmin.Group("sysAdmin").Use(func() gin.HandlerFunc {
			return func(ctx *gin.Context) {
				// ctx.AbortWithStatusJSON(200, gin.H{
				// 	"msg": "login Middleware",
				// })
			}
		}())
		{
			rgAdminUser.POST("/login", controller.Login)
		}
	})
}
