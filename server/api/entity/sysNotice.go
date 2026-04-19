/*
 * @Author: JimZhang
 * @Date: 2026-04-19 14:26:57
 * @LastEditors: JimZhang
 * @LastEditTime: 2026-04-19 14:26:57
 * @FilePath: /server/api/entity/sysNotice.go
 * @Description: 
 * 
 */
package entity

import "server/common/utils"

// SysNotice 系统通知公告
type SysNotice struct {
	ID            uint        `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	Title         string      `gorm:"column:title;type:varchar(120);comment:'标题'" json:"title"`
	Content       string      `gorm:"column:content;type:text;comment:'内容'" json:"content"`
	NoticeType    string      `gorm:"column:notice_type;type:varchar(32);comment:'通知类型'" json:"noticeType"`
	NoticeLevel   string      `gorm:"column:notice_level;type:varchar(32);comment:'通知等级'" json:"noticeLevel"`
	Status        int         `gorm:"column:status;default:1;comment:'状态：1草稿 2已发布 3已归档'" json:"status"`
	TargetType    string      `gorm:"column:target_type;type:varchar(32);comment:'目标范围'" json:"targetType"`
	CreatedBy     uint        `gorm:"column:created_by;comment:'创建人ID'" json:"createdBy"`
	CreatedByName string      `gorm:"column:created_by_name;type:varchar(64);comment:'创建人名称'" json:"createdByName"`
	PublishTime   utils.HTime `gorm:"column:publish_time;comment:'发布时间'" json:"publishTime"`
	CreateTime    utils.HTime `gorm:"column:create_time;comment:'创建时间'" json:"createTime"`
	UpdateTime    utils.HTime `gorm:"column:update_time;comment:'更新时间'" json:"updateTime"`
}

func (SysNotice) TableName() string {
	return "sys_notice"
}

// SysNoticeRead 通知已读状态
type SysNoticeRead struct {
	ID       uint        `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	NoticeId uint        `gorm:"column:notice_id;uniqueIndex:idx_notice_admin_read;comment:'通知ID'" json:"noticeId"`
	AdminId  uint        `gorm:"column:admin_id;uniqueIndex:idx_notice_admin_read;comment:'管理员ID'" json:"adminId"`
	ReadTime utils.HTime `gorm:"column:read_time;comment:'已读时间'" json:"readTime"`
}

func (SysNoticeRead) TableName() string {
	return "sys_notice_read"
}

type SysNoticeIdDto struct {
	Id uint `json:"id"`
}

type BatchSysNoticeDto struct {
	Ids []uint `json:"ids"`
}

type SaveSysNoticeDto struct {
	Id          uint   `json:"id"`
	Title       string `json:"title" validate:"required,min=2,max=120"`
	Content     string `json:"content" validate:"required,min=2,max=5000"`
	NoticeType  string `json:"noticeType" validate:"required"`
	NoticeLevel string `json:"noticeLevel" validate:"required"`
	Status      int    `json:"status" validate:"required"`
	TargetType  string `json:"targetType" validate:"required"`
}

type UpdateSysNoticeStatusDto struct {
	Id     uint `json:"id" validate:"required"`
	Status int  `json:"status" validate:"required"`
}

type SysNoticeQueryDto struct {
	PageNum    int    `form:"pageNum" json:"pageNum"`
	PageSize   int    `form:"pageSize" json:"pageSize"`
	Title      string `form:"title" json:"title"`
	NoticeType string `form:"noticeType" json:"noticeType"`
	Status     string `form:"status" json:"status"`
}

type CurrentSysNoticeQueryDto struct {
	Limit      int  `form:"limit" json:"limit"`
	UnreadOnly bool `form:"unreadOnly" json:"unreadOnly"`
}

type CurrentSysNoticeVo struct {
	ID          uint        `json:"id"`
	Title       string      `json:"title"`
	Content     string      `json:"content"`
	NoticeType  string      `json:"noticeType"`
	NoticeLevel string      `json:"noticeLevel"`
	Status      int         `json:"status"`
	TargetType  string      `json:"targetType"`
	CreatedBy   uint        `json:"createdBy"`
	CreatedByName string    `json:"createdByName"`
	PublishTime utils.HTime `json:"publishTime"`
	CreateTime  utils.HTime `json:"createTime"`
	Read        bool        `gorm:"column:is_read" json:"read"`
}
