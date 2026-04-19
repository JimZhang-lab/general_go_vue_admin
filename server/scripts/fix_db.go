package main

import (
	"fmt"
	"server/api/entity"
	"server/pkg/db"
)

func main() {
	db.SetupDBLink()

	// Print all menu items
	var menus []entity.SysMenu
	db.Db.Find(&menus)
	for _, m := range menus {
		fmt.Printf("Menu: %s, Value: %s\n", m.MenuName, m.Value)
	}

	// Update Post
	db.Db.Model(&entity.SysMenu{}).Where("menu_name LIKE ?", "%岗位%").Update("value", "system:post:list")
	db.Db.Model(&entity.SysMenu{}).Where("menu_name = ?", "新增岗位").Update("value", "system:post:add")
	db.Db.Model(&entity.SysMenu{}).Where("menu_name = ?", "修改岗位").Update("value", "system:post:update")
	db.Db.Model(&entity.SysMenu{}).Where("menu_name = ?", "删除岗位").Update("value", "system:post:delete")

	// Update Dept
	db.Db.Model(&entity.SysMenu{}).Where("menu_name LIKE ?", "%部门%").Update("value", "system:dept:list")
	db.Db.Model(&entity.SysMenu{}).Where("menu_name = ?", "新增部门").Update("value", "system:dept:add")
	db.Db.Model(&entity.SysMenu{}).Where("menu_name = ?", "修改部门").Update("value", "system:dept:update")
	db.Db.Model(&entity.SysMenu{}).Where("menu_name = ?", "删除部门").Update("value", "system:dept:delete")

	fmt.Println("Fixed db")
}
