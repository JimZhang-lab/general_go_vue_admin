/*
 * @Author: JimZhang
 * @Date: 2025-07-24 22:16:32
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-24 22:19:35
 * @FilePath: /server/api/controller/captcha.go
 * @Description: 验证码控制层
 *
 */
package controller

import (
	"server/api/service"
	"server/common/result"

	"github.com/gin-gonic/gin"
)

// @Summary 验证码 接口
// @Produce  application/json
// @Description 验证码接口
// @Success 200 {object} result.Result
// @router /api/captcha [get]
func Captcha(c *gin.Context) {
	id, base64Image := service.CaptMake()
	result.Success(c, map[string]interface{}{"idKey": id, "image": base64Image})
}
