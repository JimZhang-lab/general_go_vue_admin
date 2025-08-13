/*
 * @Author: JimZhang
 * @Date: 2025-07-24 11:31:20
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-25 20:59:50
 * @FilePath: /server/pkg/db/db.go
 * @Description: 初始化数据库连接
 *
 */
package db

import (
	"fmt"
	"server/api/entity"
	"server/common/config"
	"server/pkg/seed"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func SetupDBLink() error {

	var err error
	var dbConfig = config.Config.DB
	url := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local&timeout=10s&collation=utf8mb4_general_ci",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
		dbConfig.Charset)

	Db, err = gorm.Open(mysql.Open(url), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}
	if Db.Error != nil {
		panic(Db.Error)
	}

	// 自动建表
	err = Db.AutoMigrate(
		&entity.SysAdmin{},
		&entity.SysAdminRole{},
		&entity.SysPost{},
		&entity.SysDept{},
		&entity.SysRole{},
		&entity.SysRoleMenu{},
		&entity.SysMenu{},
		&entity.SysLoginInfo{},
		&entity.SysOperationLog{},
	)
	if err != nil {
		return err
	}

	// 首次运行初始化数据（幂等）
	if err := seed.InitDataIfNeeded(Db); err != nil {
		// 不阻断启动，仅打印错误
		// 可根据需要改为返回错误
	}

	sqlDB, err := Db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	return nil
}
