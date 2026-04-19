package service

import (
	"server/api/dao"
	"server/api/entity"
	"server/common/result"
	"server/common/utils"
	"server/pkg/jwt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ISysNoticeService interface {
	GetSysNoticeList(c *gin.Context, dto entity.SysNoticeQueryDto)
	GetCurrentSysNoticeList(c *gin.Context, dto entity.CurrentSysNoticeQueryDto)
	CreateSysNotice(c *gin.Context, dto entity.SaveSysNoticeDto)
	UpdateSysNotice(c *gin.Context, dto entity.SaveSysNoticeDto)
	UpdateSysNoticeStatus(c *gin.Context, dto entity.UpdateSysNoticeStatusDto)
	DeleteSysNoticeById(c *gin.Context, dto entity.SysNoticeIdDto)
	BatchDeleteSysNotice(c *gin.Context, dto entity.BatchSysNoticeDto)
	MarkSysNoticeRead(c *gin.Context, dto entity.SysNoticeIdDto)
	MarkAllSysNoticeRead(c *gin.Context)
}

type SysNoticeServiceImpl struct{}

func normalizeNoticeStatus(notice *entity.SysNotice, status int) {
	now := utils.HTime{Time: time.Now()}
	notice.Status = status
	notice.UpdateTime = now
	if status == 2 {
		notice.PublishTime = now
	}
}

func (s SysNoticeServiceImpl) GetSysNoticeList(c *gin.Context, dto entity.SysNoticeQueryDto) {
	if dto.PageSize < 1 {
		dto.PageSize = 10
	}
	if dto.PageNum < 1 {
		dto.PageNum = 1
	}
	list, count := dao.GetSysNoticeList(dto.Title, dto.NoticeType, dto.Status, dto.PageSize, dto.PageNum)
	result.Success(c, map[string]interface{}{
		"total":    count,
		"pageSize": dto.PageSize,
		"pageNum":  dto.PageNum,
		"list":     list,
	})
}

func (s SysNoticeServiceImpl) GetCurrentSysNoticeList(c *gin.Context, dto entity.CurrentSysNoticeQueryDto) {
	if dto.Limit <= 0 {
		dto.Limit = 10
	}
	if !dao.GetSysSettingBoolValue("notification.enable_site_notice", true) {
		result.Success(c, make([]entity.CurrentSysNoticeVo, 0))
		return
	}
	adminID, err := jwt.GetAdminId(c)
	if err != nil {
		result.Failed(c, int(result.ApiCode.NOAUTH), "当前用户未登录")
		return
	}
	result.Success(c, dao.GetCurrentSysNoticeList(adminID, dto.Limit, dto.UnreadOnly))
}

func (s SysNoticeServiceImpl) CreateSysNotice(c *gin.Context, dto entity.SaveSysNoticeDto) {
	if err := validator.New().Struct(dto); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "通知参数不完整")
		return
	}
	admin, err := jwt.GetAdmin(c)
	if err != nil {
		result.Failed(c, int(result.ApiCode.NOAUTH), "当前用户未登录")
		return
	}
	now := utils.HTime{Time: time.Now()}
	notice := entity.SysNotice{
		Title:         dto.Title,
		Content:       dto.Content,
		NoticeType:    dto.NoticeType,
		NoticeLevel:   dto.NoticeLevel,
		TargetType:    dto.TargetType,
		CreatedBy:     admin.ID,
		CreatedByName: admin.Username,
		CreateTime:    now,
		UpdateTime:    now,
	}
	normalizeNoticeStatus(&notice, dto.Status)
	if err := dao.SaveSysNotice(&notice); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "新增通知失败")
		return
	}
	result.Success(c, notice)
}

func (s SysNoticeServiceImpl) UpdateSysNotice(c *gin.Context, dto entity.SaveSysNoticeDto) {
	if dto.Id == 0 {
		result.Failed(c, int(result.ApiCode.FAILED), "通知ID不能为空")
		return
	}
	if err := validator.New().Struct(dto); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "通知参数不完整")
		return
	}
	notice := dao.GetSysNoticeById(dto.Id)
	if notice.ID == 0 {
		result.Failed(c, int(result.ApiCode.FAILED), "通知不存在")
		return
	}
	notice.Title = dto.Title
	notice.Content = dto.Content
	notice.NoticeType = dto.NoticeType
	notice.NoticeLevel = dto.NoticeLevel
	notice.TargetType = dto.TargetType
	normalizeNoticeStatus(&notice, dto.Status)
	if err := dao.SaveSysNotice(&notice); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "更新通知失败")
		return
	}
	result.Success(c, notice)
}

func (s SysNoticeServiceImpl) UpdateSysNoticeStatus(c *gin.Context, dto entity.UpdateSysNoticeStatusDto) {
	if err := validator.New().Struct(dto); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "通知状态参数不完整")
		return
	}
	notice := dao.GetSysNoticeById(dto.Id)
	if notice.ID == 0 {
		result.Failed(c, int(result.ApiCode.FAILED), "通知不存在")
		return
	}
	normalizeNoticeStatus(&notice, dto.Status)
	if err := dao.SaveSysNotice(&notice); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "更新通知状态失败")
		return
	}
	result.Success(c, true)
}

func (s SysNoticeServiceImpl) DeleteSysNoticeById(c *gin.Context, dto entity.SysNoticeIdDto) {
	if err := dao.DeleteSysNoticeById(dto.Id); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "删除通知失败")
		return
	}
	result.Success(c, true)
}

func (s SysNoticeServiceImpl) BatchDeleteSysNotice(c *gin.Context, dto entity.BatchSysNoticeDto) {
	if len(dto.Ids) == 0 {
		result.Failed(c, int(result.ApiCode.FAILED), "缺少待删除的通知")
		return
	}
	if err := dao.BatchDeleteSysNotice(dto.Ids); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "批量删除通知失败")
		return
	}
	result.Success(c, true)
}

func (s SysNoticeServiceImpl) MarkSysNoticeRead(c *gin.Context, dto entity.SysNoticeIdDto) {
	adminID, err := jwt.GetAdminId(c)
	if err != nil {
		result.Failed(c, int(result.ApiCode.NOAUTH), "当前用户未登录")
		return
	}
	if dto.Id == 0 {
		result.Failed(c, int(result.ApiCode.FAILED), "通知ID不能为空")
		return
	}
	if err := dao.MarkSysNoticeRead(adminID, dto.Id); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "通知已读更新失败")
		return
	}
	result.Success(c, true)
}

func (s SysNoticeServiceImpl) MarkAllSysNoticeRead(c *gin.Context) {
	adminID, err := jwt.GetAdminId(c)
	if err != nil {
		result.Failed(c, int(result.ApiCode.NOAUTH), "当前用户未登录")
		return
	}
	if err := dao.MarkAllSysNoticeRead(adminID); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "全部已读操作失败")
		return
	}
	result.Success(c, true)
}

var sysNoticeService = SysNoticeServiceImpl{}

func SysNoticeService() ISysNoticeService {
	return &sysNoticeService
}
