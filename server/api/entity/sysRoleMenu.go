/*
 * @Author: JimZhang
 * @Date: 2023-07-22 15:51:48
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-25 20:25:51
 * @FilePath: /server/api/entity/sysRoleMenu.go
 * @Description:
 *
 */

package entity

// SysRoleMenu 角色与菜单关系模型
type SysRoleMenu struct {
	RoleId uint `gorm:"column:role_id;comment:'角色id';NOT NULL" json:"roleId"` // 角色id
	MenuId uint `gorm:"column:menu_id;comment:'菜单id';NOT NULL" json:"menuId"` // 菜单id
}

func (SysRoleMenu) TableName() string {
	return "sys_role_menu"
}
