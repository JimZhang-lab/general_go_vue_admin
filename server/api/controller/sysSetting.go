/*
 * @Author: JimZhang
 * @Date: 2026-04-19 14:26:57
 * @LastEditors: JimZhang
 * @LastEditTime: 2026-04-19 14:26:57
 * @FilePath: /server/api/controller/sysSetting.go
 * @Description: 
 * 
 */
package controller

import (
	"server/api/entity"
	"server/api/service"

	"github.com/gin-gonic/gin"
)

// GetSysSettingList 获取系统设置列表
func GetSysSettingList(c *gin.Context) {
	var dto entity.SysSettingQueryDto
	_ = c.ShouldBindQuery(&dto)
	service.SysSettingService().GetSysSettingList(c, dto.GroupKey)
}

// BatchUpdateSysSettings 批量更新系统设置
func BatchUpdateSysSettings(c *gin.Context) {
	var dto entity.BatchUpdateSysSettingDto
	_ = c.BindJSON(&dto)
	service.SysSettingService().BatchUpdateSysSettings(c, dto)
}
