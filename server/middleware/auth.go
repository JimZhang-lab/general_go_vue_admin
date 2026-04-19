/*
 * @Author: JimZhang
 * @Date: 2025-07-24 12:06:22
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-26 00:39:59
 * @FilePath: /server/middleware/auth.go
 * @Description: 鉴权中间件
 *
 */
package middleware

import (
	"server/common/constant"
	"server/common/result"
	"server/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 鉴权中间件，用于验证请求是否包含有效的Authorization头
// 返回一个gin.HandlerFunc处理函数，该函数会检查请求头中的Authorization字段
// 如果Authorization头为空，则返回未授权错误并中止请求处理
func AuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			result.Failed(c, int(result.ApiCode.NOAUTH), result.ApiCode.GetMessage(result.ApiCode.NOAUTH))
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			result.Failed(c, int(result.ApiCode.AUTHFORMATERROR), result.ApiCode.GetMessage(result.ApiCode.AUTHFORMATERROR))
			c.Abort()
			return
		}
		mc, err := jwt.ValidateToken(parts[1])
		if err != nil {
			result.Failed(c, int(result.ApiCode.INVALIDTOKEN), result.ApiCode.GetMessage(result.ApiCode.INVALIDTOKEN))
			c.Abort()
			return
		}
		c.Set(constant.ContextKeyUserObject, mc)
		c.Next()
	}
}
