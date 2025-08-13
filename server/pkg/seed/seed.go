package seed

import (
	"fmt"
	"time"

	"server/api/entity"
	"server/common/config"
	"server/common/utils"

	"gorm.io/gorm"
)

// InitDataIfNeeded 在首启时初始化基础数据（幂等）
func InitDataIfNeeded(db *gorm.DB) error {
	// 1) 角色
	var roleCount int64
	if err := db.Model(&entity.SysRole{}).Count(&roleCount).Error; err != nil {
		return fmt.Errorf("count roles failed: %w", err)
	}
	var adminRole entity.SysRole
	if roleCount == 0 {
		adminRole = entity.SysRole{
			RoleName:    "管理员",
			RoleKey:     "admin",
			Status:      1,
			Description: "系统管理员角色",
			CreateTime:  utils.HTime{Time: time.Now()},
		}
		if err := db.Create(&adminRole).Error; err != nil {
			return fmt.Errorf("create admin role failed: %w", err)
		}
	} else {
		// 尝试获取已有管理员角色
		db.Where("role_key = ?", "admin").First(&adminRole)
	}

	// 2) 部门
	var deptCount int64
	if err := db.Model(&entity.SysDept{}).Count(&deptCount).Error; err != nil {
		return fmt.Errorf("count depts failed: %w", err)
	}
	var rootDept entity.SysDept
	if deptCount == 0 {
		rootDept = entity.SysDept{
			ParentId:   0,
			DeptType:   1,
			DeptName:   "总部",
			DeptStatus: 1,
			CreateTime: utils.HTime{Time: time.Now()},
		}
		if err := db.Create(&rootDept).Error; err != nil {
			return fmt.Errorf("create root dept failed: %w", err)
		}
	} else {
		// 取一个作为默认部门
		db.First(&rootDept)
	}

	// 3) 岗位
	var postCount int64
	if err := db.Model(&entity.SysPost{}).Count(&postCount).Error; err != nil {
		return fmt.Errorf("count posts failed: %w", err)
	}
	var adminPost entity.SysPost
	if postCount == 0 {
		adminPost = entity.SysPost{
			PostCode:   "ADMIN",
			PostName:   "管理员",
			PostStatus: 1,
			CreateTime: utils.HTime{Time: time.Now()},
			Remark:     "系统默认岗位",
		}
		if err := db.Create(&adminPost).Error; err != nil {
			return fmt.Errorf("create admin post failed: %w", err)
		}
	} else {
		db.First(&adminPost)
	}

	// 4) 管理员账号
	var adminCount int64
	if err := db.Model(&entity.SysAdmin{}).Count(&adminCount).Error; err != nil {
		return fmt.Errorf("count admins failed: %w", err)
	}
	if adminCount == 0 && config.Config.Seed.Enable {
		// 允许在配置中自定义默认管理员信息
		u := config.Config.Seed.Admin
		if u.Username == "" {
			u.Username = "admin"
		}
		if u.Password == "" {
			u.Password = "admin123"
		}
		if u.Nickname == "" {
			u.Nickname = "系统管理员"
		}
		if u.Email == "" {
			u.Email = "admin@example.com"
		}
		if u.Phone == "" {
			u.Phone = "13800138000"
		}

		admin := entity.SysAdmin{
			PostId:     int(adminPost.ID),
			DeptId:     int(rootDept.ID),
			Username:   u.Username,
			Password:   utils.EncryptionMd5(u.Password),
			Nickname:   u.Nickname,
			Status:     1,
			Email:      u.Email,
			Phone:      u.Phone,
			Note:       "系统管理员账号",
			CreateTime: utils.HTime{Time: time.Now()},
		}
		if err := db.Create(&admin).Error; err != nil {
			return fmt.Errorf("create admin failed: %w", err)
		}
		// 关联角色
		if adminRole.ID > 0 {
			link := entity.SysAdminRole{AdminId: admin.ID, RoleId: adminRole.ID}
			_ = db.Create(&link).Error
		}
	}

	return nil
}
