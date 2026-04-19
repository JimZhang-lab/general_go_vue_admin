/*
 * @Author: JimZhang
 * @Date: 2026-04-19 14:26:57
 * @LastEditors: JimZhang
 * @LastEditTime: 2026-04-19 14:26:57
 * @FilePath: /server/api/entity/sys_dict.go
 * @Description: 
 * 
 */
package entity

import "time"

// SysDictType 字典类型表
type SysDictType struct {
	Id         int       `gorm:"primaryKey;autoIncrement" json:"id"` // 字典主键
	DictName   string    `gorm:"size:100;default:''" json:"dictName"` // 字典名称
	DictType   string    `gorm:"size:100;uniqueIndex;default:''" json:"dictType"` // 字典类型
	Status     string    `gorm:"size:1;default:'1'" json:"status"` // 状态（1正常 2停用）
	Remark     string    `gorm:"size:500;default:''" json:"remark"` // 备注
	CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"` // 创建时间
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"updateTime"` // 更新时间
}

// TableName 表名
func (SysDictType) TableName() string {
	return "sys_dict_type"
}

// SysDictData 字典数据表
type SysDictData struct {
	Id         int       `gorm:"primaryKey;autoIncrement" json:"id"` // 字典编码
	DictSort   int       `gorm:"default:0" json:"dictSort"` // 字典排序
	DictLabel  string    `gorm:"size:100;default:''" json:"dictLabel"` // 字典标签
	DictValue  string    `gorm:"size:100;default:''" json:"dictValue"` // 字典键值
	DictType   string    `gorm:"size:100;default:''" json:"dictType"` // 字典类型（关联SysDictType.DictType）
	Status     string    `gorm:"size:1;default:'1'" json:"status"` // 状态（1正常 2停用）
	Remark     string    `gorm:"size:500;default:''" json:"remark"` // 备注
	CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"` // 创建时间
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"updateTime"` // 更新时间
}

// TableName 表名
func (SysDictData) TableName() string {
	return "sys_dict_data"
}
