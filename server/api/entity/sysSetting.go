/*
 * @Author: JimZhang
 * @Date: 2026-04-19 14:26:57
 * @LastEditors: JimZhang
 * @LastEditTime: 2026-04-19 14:26:57
 * @FilePath: /server/api/entity/sysSetting.go
 * @Description: 
 * 
 */
package entity

import "server/common/utils"

// SysSetting 系统设置
type SysSetting struct {
	ID           uint        `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	GroupKey     string      `gorm:"column:group_key;type:varchar(64);index:idx_sys_setting_group_key;comment:'分组标识'" json:"groupKey"`
	SettingKey   string      `gorm:"column:setting_key;type:varchar(128);uniqueIndex:idx_sys_setting_setting_key;comment:'设置键'" json:"settingKey"`
	SettingName  string      `gorm:"column:setting_name;type:varchar(128);comment:'设置名称'" json:"settingName"`
	SettingValue string      `gorm:"column:setting_value;type:text;comment:'设置值'" json:"settingValue"`
	ValueType    string      `gorm:"column:value_type;type:varchar(32);comment:'值类型'" json:"valueType"`
	OptionsJSON  string      `gorm:"column:options_json;type:text;comment:'扩展选项'" json:"optionsJson"`
	IsEncrypted  bool        `gorm:"column:is_encrypted;default:false;comment:'是否加密'" json:"isEncrypted"`
	Sort         int         `gorm:"column:sort;default:0;comment:'排序'" json:"sort"`
	Remark       string      `gorm:"column:remark;type:varchar(255);comment:'备注'" json:"remark"`
	CreateTime   utils.HTime `gorm:"column:create_time;comment:'创建时间'" json:"createTime"`
	UpdateTime   utils.HTime `gorm:"column:update_time;comment:'更新时间'" json:"updateTime"`
}

func (SysSetting) TableName() string {
	return "sys_setting"
}

type SysSettingQueryDto struct {
	GroupKey string `form:"groupKey" json:"groupKey"`
}

type SysSettingItemDto struct {
	GroupKey     string `json:"groupKey" validate:"required"`
	SettingKey   string `json:"settingKey" validate:"required"`
	SettingName  string `json:"settingName" validate:"required"`
	SettingValue string `json:"settingValue"`
	ValueType    string `json:"valueType"`
	OptionsJSON  string `json:"optionsJson"`
	IsEncrypted  bool   `json:"isEncrypted"`
	Sort         int    `json:"sort"`
	Remark       string `json:"remark"`
}

type BatchUpdateSysSettingDto struct {
	Items []SysSettingItemDto `json:"items" validate:"required,min=1,dive"`
}
