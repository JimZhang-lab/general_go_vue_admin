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
	"time"

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

	gormLogLevel := logger.LogLevel(dbConfig.LogLevel)
	if gormLogLevel < logger.Silent || gormLogLevel > logger.Info {
		gormLogLevel = logger.Info
	}

	Db, err = gorm.Open(mysql.Open(url), &gorm.Config{
		Logger:                                   logger.Default.LogMode(gormLogLevel),
		DisableForeignKeyConstraintWhenMigrating: true,
		PrepareStmt:                              dbConfig.PrepareStmt,
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}
	if Db.Error != nil {
		return fmt.Errorf("数据库实例初始化失败: %w", Db.Error)
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
		&entity.SysSetting{},
		&entity.SysNotice{},
		&entity.SysNoticeRead{},
		&entity.SysDictType{},
		&entity.SysDictData{},
	)
	if err != nil {
		return err
	}

	// 首次运行初始化数据（幂等）
	if err := seed.InitDataIfNeeded(Db); err != nil {
		// 不阻断启动，但打印初始化失败信息，避免静默失败
		fmt.Printf("seed init failed: %v\n", err)
	}

	sqlDB, err := Db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	if dbConfig.SetConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.SetConnMaxLifetime) * time.Second)
	}
	if dbConfig.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(time.Duration(dbConfig.ConnMaxIdleTime) * time.Second)
	}
	return nil
}
