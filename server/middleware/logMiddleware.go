/*
 * @Author: JimZhang
 * @Date: 2023-07-23 15:49:28
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-26 00:37:06
 * @FilePath: /server/middleware/logMiddleware.go
 * @Description:
 *
 */

package middleware

import (
	"server/api/dao"
	"server/api/entity"
	"server/common/utils"
	"server/pkg/jwt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := strings.ToLower(c.Request.Method)
		sysAdmin, _ := jwt.GetAdmin(c)
		if method != "get" {
			log := entity.SysOperationLog{
				AdminId:    sysAdmin.ID,
				Username:   sysAdmin.Username,
				Method:     method,
				Ip:         c.ClientIP(),
				Url:        c.Request.URL.Path,
				CreateTime: utils.HTime{Time: time.Now()},
			}
			dao.CreateSysOperationLog(log)
		}
	}
}
