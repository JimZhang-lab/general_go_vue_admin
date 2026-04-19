package service

import (
	"server/api/dao"
	"server/api/entity"
	"server/common/result"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ISysSettingService interface {
	GetSysSettingList(c *gin.Context, groupKey string)
	BatchUpdateSysSettings(c *gin.Context, dto entity.BatchUpdateSysSettingDto)
}

type SysSettingServiceImpl struct{}

func (s SysSettingServiceImpl) GetSysSettingList(c *gin.Context, groupKey string) {
	result.Success(c, dao.GetSysSettingList(groupKey))
}

func (s SysSettingServiceImpl) BatchUpdateSysSettings(c *gin.Context, dto entity.BatchUpdateSysSettingDto) {
	if err := validator.New().Struct(dto); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "系统设置参数不完整")
		return
	}
	if err := dao.BatchUpsertSysSettings(dto.Items); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "系统设置保存失败")
		return
	}
	result.Success(c, true)
}

var sysSettingService = SysSettingServiceImpl{}

func SysSettingService() ISysSettingService {
	return &sysSettingService
}
