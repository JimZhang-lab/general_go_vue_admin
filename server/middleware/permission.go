/*
 * @Author: JimZhang
 * @Date: 2026-04-19 14:26:57
 * @LastEditors: JimZhang
 * @LastEditTime: 2026-04-19 14:26:57
 * @FilePath: /server/middleware/permission.go
 * @Description: 
 * 
 */
package middleware

import (
	"server/api/dao"
	"server/common/result"
	"server/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func PermissionMiddleware(requiredPermissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(requiredPermissions) == 0 {
			c.Next()
			return
		}

		adminID, err := jwt.GetAdminId(c)
		if err != nil {
			result.Failed(c, int(result.ApiCode.NOAUTH), "当前用户未登录")
			c.Abort()
			return
		}

		permissionList := dao.QueryPermissionList(adminID)
		permissionSet := make(map[string]struct{}, len(permissionList))
		for _, item := range permissionList {
			if item.Value == "" {
				continue
			}
			permissionSet[item.Value] = struct{}{}
		}

		for _, permission := range requiredPermissions {
			if _, ok := permissionSet[permission]; ok {
				c.Next()
				return
			}
		}

		result.Failed(c, int(result.ApiCode.NOAUTH), "暂无访问权限")
		c.Abort()
	}
}
