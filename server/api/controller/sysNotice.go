package controller

import (
	"server/api/entity"
	"server/api/service"

	"github.com/gin-gonic/gin"
)

func GetSysNoticeList(c *gin.Context) {
	var dto entity.SysNoticeQueryDto
	_ = c.ShouldBindQuery(&dto)
	service.SysNoticeService().GetSysNoticeList(c, dto)
}

func GetCurrentSysNoticeList(c *gin.Context) {
	var dto entity.CurrentSysNoticeQueryDto
	_ = c.ShouldBindQuery(&dto)
	service.SysNoticeService().GetCurrentSysNoticeList(c, dto)
}

func CreateSysNotice(c *gin.Context) {
	var dto entity.SaveSysNoticeDto
	_ = c.BindJSON(&dto)
	service.SysNoticeService().CreateSysNotice(c, dto)
}

func UpdateSysNotice(c *gin.Context) {
	var dto entity.SaveSysNoticeDto
	_ = c.BindJSON(&dto)
	service.SysNoticeService().UpdateSysNotice(c, dto)
}

func UpdateSysNoticeStatus(c *gin.Context) {
	var dto entity.UpdateSysNoticeStatusDto
	_ = c.BindJSON(&dto)
	service.SysNoticeService().UpdateSysNoticeStatus(c, dto)
}

func DeleteSysNoticeById(c *gin.Context) {
	var dto entity.SysNoticeIdDto
	_ = c.BindJSON(&dto)
	service.SysNoticeService().DeleteSysNoticeById(c, dto)
}

func BatchDeleteSysNotice(c *gin.Context) {
	var dto entity.BatchSysNoticeDto
	_ = c.BindJSON(&dto)
	service.SysNoticeService().BatchDeleteSysNotice(c, dto)
}

func MarkSysNoticeRead(c *gin.Context) {
	var dto entity.SysNoticeIdDto
	_ = c.BindJSON(&dto)
	service.SysNoticeService().MarkSysNoticeRead(c, dto)
}

func MarkAllSysNoticeRead(c *gin.Context) {
	service.SysNoticeService().MarkAllSysNoticeRead(c)
}
