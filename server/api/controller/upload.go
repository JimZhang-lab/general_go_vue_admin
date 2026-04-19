/*
 * @Author: JimZhang
 * @Date: 2023-07-23 12:42:46
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-25 20:31:24
 * @FilePath: /server/api/controller/upload.go
 * @Description:
 *
 */

package controller

import (
	"server/api/service"

	"github.com/gin-gonic/gin"
)

// 单图片上传
// @Summary 单图片上传接口
// @Description 单图片上传接口
// @Produce json
// @Accept multipart/form-data
// @Param file formData file true "file"
// @Success 200 {object} result.Result
// @Router /api/upload [post]
// @Security ApiKeyAuth
func Upload(c *gin.Context) {
	service.UploadService().Upload(c)
}
