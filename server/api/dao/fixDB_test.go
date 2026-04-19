package dao

import (
	"fmt"
	"server/api/entity"
	"server/pkg/db"
	"testing"
)

func TestFixDB(t *testing.T) {
	fmt.Println("Starting FixDB test")
	if err := db.SetupDBLink(); err != nil {
		t.Fatalf("failed to setup db: %v", err)
	}

	// Update Post
	db.Db.Model(&entity.SysMenu{}).Where("menu_name = ?", "岗位管理").Update("value", "system:post:list")
	db.Db.Model(&entity.SysMenu{}).Where("menu_name = ?", "新增岗位").Update("value", "system:post:add")
	db.Db.Model(&entity.SysMenu{}).Where("menu_name = ?", "修改岗位").Update("value", "system:post:update")
	db.Db.Model(&entity.SysMenu{}).Where("menu_name = ?", "删除岗位").Update("value", "system:post:delete")

	// Update Dept
	db.Db.Model(&entity.SysMenu{}).Where("menu_name = ?", "部门管理").Update("value", "system:dept:list")
	db.Db.Model(&entity.SysMenu{}).Where("menu_name = ?", "新增部门").Update("value", "system:dept:add")
	db.Db.Model(&entity.SysMenu{}).Where("menu_name = ?", "修改部门").Update("value", "system:dept:update")
	db.Db.Model(&entity.SysMenu{}).Where("menu_name = ?", "删除部门").Update("value", "system:dept:delete")

	// Add missing ones if any
	db.Db.Model(&entity.SysMenu{}).Where("menu_name = ?", "数据字典").Update("value", "system:dict:list")
	db.Db.Model(&entity.SysMenu{}).Where("menu_name = ?", "系统监控").Update("value", "system:monitor:list")

	fmt.Println("Fixed DB successfully")
}
