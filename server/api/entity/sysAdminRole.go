/*
 * @Author: JimZhang
 * @Date: 2023-07-22 15:51:48
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-25 20:25:12
 * @FilePath: /server/api/entity/sysAdminRole.go
 * @Description:
 *
 */

package entity

type SysAdminRole struct {
	RoleId  uint `gorm:"column:role_id;comment:'角色id';NOT NULL" json:"roleId"`  // 角色id
	AdminId uint `gorm:"column:admin_id;comment:'用户id';NOT NULL" json:"menuId"` // 用户id
}

func (SysAdminRole) TableName() string {
	return "sys_admin_role"
}
