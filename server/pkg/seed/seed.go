package seed

import (
	"fmt"
	"time"

	"server/api/entity"
	"server/common/config"
	"server/common/utils"

	"gorm.io/gorm"
)

type seedMenuItem struct {
	Key       string
	ParentKey string
	MenuName  string
	Icon      string
	Value     string
	MenuType  uint
	URL       string
	Status    uint
	Sort      uint
}

var defaultMenus = []seedMenuItem{
	{Key: "auth_root", MenuName: "权限管理", Icon: "bar-chart", MenuType: 1, URL: "/auth", Status: 2, Sort: 10},
	{Key: "auth_dashboard", ParentKey: "auth_root", MenuName: "权限总览", Icon: "bar-chart", Value: "system:auth:view", MenuType: 2, URL: "/auth/dashboard", Status: 2, Sort: 10},
	{Key: "auth_admin", ParentKey: "auth_root", MenuName: "管理员管理", Icon: "users", Value: "system:admin:list", MenuType: 2, URL: "/auth/admin", Status: 2, Sort: 20},
	{Key: "auth_admin_add", ParentKey: "auth_admin", MenuName: "新增管理员", Value: "system:admin:add", MenuType: 3, Status: 2, Sort: 10},
	{Key: "auth_admin_update", ParentKey: "auth_admin", MenuName: "编辑管理员", Value: "system:admin:update", MenuType: 3, Status: 2, Sort: 20},
	{Key: "auth_admin_delete", ParentKey: "auth_admin", MenuName: "删除管理员", Value: "system:admin:delete", MenuType: 3, Status: 2, Sort: 30},
	{Key: "auth_admin_reset_password", ParentKey: "auth_admin", MenuName: "重置管理员密码", Value: "system:admin:resetPassword", MenuType: 3, Status: 2, Sort: 40},
	{Key: "auth_admin_update_status", ParentKey: "auth_admin", MenuName: "修改管理员状态", Value: "system:admin:updateStatus", MenuType: 3, Status: 2, Sort: 50},
	{Key: "auth_role", ParentKey: "auth_root", MenuName: "角色管理", Icon: "shield", Value: "system:role:list", MenuType: 2, URL: "/auth/role", Status: 2, Sort: 30},
	{Key: "auth_role_add", ParentKey: "auth_role", MenuName: "新增角色", Value: "system:role:add", MenuType: 3, Status: 2, Sort: 10},
	{Key: "auth_role_update", ParentKey: "auth_role", MenuName: "编辑角色", Value: "system:role:update", MenuType: 3, Status: 2, Sort: 20},
	{Key: "auth_role_delete", ParentKey: "auth_role", MenuName: "删除角色", Value: "system:role:delete", MenuType: 3, Status: 2, Sort: 30},
	{Key: "auth_role_assign", ParentKey: "auth_role", MenuName: "分配角色权限", Value: "system:role:assign", MenuType: 3, Status: 2, Sort: 40},
	{Key: "auth_menu", ParentKey: "auth_root", MenuName: "权限管理", Icon: "key", Value: "system:menu:list", MenuType: 2, URL: "/auth/permission", Status: 2, Sort: 40},
	{Key: "auth_menu_add", ParentKey: "auth_menu", MenuName: "新增权限", Value: "system:menu:add", MenuType: 3, Status: 2, Sort: 10},
	{Key: "auth_menu_update", ParentKey: "auth_menu", MenuName: "编辑权限", Value: "system:menu:update", MenuType: 3, Status: 2, Sort: 20},
	{Key: "auth_menu_delete", ParentKey: "auth_menu", MenuName: "删除权限", Value: "system:menu:delete", MenuType: 3, Status: 2, Sort: 30},
	{Key: "auth_logs", ParentKey: "auth_root", MenuName: "系统日志", Icon: "logs", Value: "system:log:list", MenuType: 2, URL: "/auth/logs", Status: 2, Sort: 50},
	{Key: "auth_log_delete", ParentKey: "auth_logs", MenuName: "删除日志", Value: "system:log:delete", MenuType: 3, Status: 2, Sort: 10},
	{Key: "auth_log_clean", ParentKey: "auth_logs", MenuName: "清空日志", Value: "system:log:clean", MenuType: 3, Status: 2, Sort: 20},
	{Key: "personal_root", MenuName: "个人中心", Icon: "user-circle", MenuType: 1, URL: "/auth/profile", Status: 2, Sort: 20},
	{Key: "personal_profile", ParentKey: "personal_root", MenuName: "个人资料", Icon: "user-circle", Value: "system:profile:view", MenuType: 2, URL: "/auth/profile", Status: 2, Sort: 10},
	{Key: "notice_center", ParentKey: "personal_root", MenuName: "通知中心", Icon: "bell", Value: "system:notice:list", MenuType: 2, URL: "/auth/notifications", Status: 2, Sort: 20},
	{Key: "notice_add", ParentKey: "notice_center", MenuName: "新增通知", Value: "system:notice:add", MenuType: 3, Status: 2, Sort: 10},
	{Key: "notice_update", ParentKey: "notice_center", MenuName: "编辑通知", Value: "system:notice:update", MenuType: 3, Status: 2, Sort: 20},
	{Key: "notice_delete", ParentKey: "notice_center", MenuName: "删除通知", Value: "system:notice:delete", MenuType: 3, Status: 2, Sort: 30},
	{Key: "notice_publish", ParentKey: "notice_center", MenuName: "发布通知", Value: "system:notice:publish", MenuType: 3, Status: 2, Sort: 40},
	{Key: "settings_root", MenuName: "系统设置", Icon: "settings", MenuType: 1, URL: "/auth/settings/basic", Status: 2, Sort: 30},
	{Key: "settings_basic", ParentKey: "settings_root", MenuName: "基础设置", Icon: "settings", Value: "system:setting:list", MenuType: 2, URL: "/auth/settings/basic", Status: 2, Sort: 10},
	{Key: "settings_security", ParentKey: "settings_root", MenuName: "安全设置", Icon: "shield", Value: "system:setting:list", MenuType: 2, URL: "/auth/settings/security", Status: 2, Sort: 20},
	{Key: "settings_notification", ParentKey: "settings_root", MenuName: "通知设置", Icon: "bell", Value: "system:setting:list", MenuType: 2, URL: "/auth/settings/notification", Status: 2, Sort: 30},
	{Key: "settings_update", ParentKey: "settings_basic", MenuName: "更新系统设置", Value: "system:setting:update", MenuType: 3, Status: 2, Sort: 10},
}

func ensureDefaultMenus(db *gorm.DB) (map[string]entity.SysMenu, error) {
	now := utils.HTime{Time: time.Now()}
	keyToMenu := make(map[string]entity.SysMenu, len(defaultMenus))

	for _, item := range defaultMenus {
		parentID := uint(0)
		if item.ParentKey != "" {
			parentMenu, ok := keyToMenu[item.ParentKey]
			if !ok {
				return nil, fmt.Errorf("parent menu %s not found", item.ParentKey)
			}
			parentID = parentMenu.ID
		}

		var menu entity.SysMenu
		err := db.Where("parent_id = ? AND menu_name = ? AND menu_type = ?", parentID, item.MenuName, item.MenuType).First(&menu).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}

		menu.ParentId = parentID
		menu.MenuName = item.MenuName
		menu.Icon = item.Icon
		menu.Value = item.Value
		menu.MenuType = item.MenuType
		menu.Url = item.URL
		menu.MenuStatus = item.Status
		menu.Sort = item.Sort

		if err == gorm.ErrRecordNotFound {
			menu.CreateTime = now
			if createErr := db.Create(&menu).Error; createErr != nil {
				return nil, createErr
			}
		} else if saveErr := db.Save(&menu).Error; saveErr != nil {
			return nil, saveErr
		}

		keyToMenu[item.Key] = menu
	}

	return keyToMenu, nil
}

func ensureAdminRolePermissions(db *gorm.DB, roleID uint) error {
	var menuIDs []uint
	if err := db.Model(&entity.SysMenu{}).Pluck("id", &menuIDs).Error; err != nil {
		return err
	}

	for _, menuID := range menuIDs {
		var relation entity.SysRoleMenu
		err := db.Where("role_id = ? AND menu_id = ?", roleID, menuID).First(&relation).Error
		if err == nil {
			continue
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}

		relation = entity.SysRoleMenu{
			RoleId: roleID,
			MenuId: menuID,
		}
		if createErr := db.Create(&relation).Error; createErr != nil {
			return createErr
		}
	}

	return nil
}

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

	if _, err := ensureDefaultMenus(db); err != nil {
		return fmt.Errorf("create default menus failed: %w", err)
	}
	if adminRole.ID > 0 {
		if err := ensureAdminRolePermissions(db, adminRole.ID); err != nil {
			return fmt.Errorf("assign admin permissions failed: %w", err)
		}
	}

	// 5) 系统设置
	var settingCount int64
	if err := db.Model(&entity.SysSetting{}).Count(&settingCount).Error; err != nil {
		return fmt.Errorf("count settings failed: %w", err)
	}
	if settingCount == 0 {
		now := utils.HTime{Time: time.Now()}
		defaultSettings := []entity.SysSetting{
			{GroupKey: "basic", SettingKey: "basic.site_name", SettingName: "站点名称", SettingValue: "通用后台管理系统", ValueType: "text", Sort: 10, Remark: "系统标题", CreateTime: now, UpdateTime: now},
			{GroupKey: "basic", SettingKey: "basic.site_subtitle", SettingName: "站点副标题", SettingValue: "稳定、安全、可扩展的管理后台", ValueType: "text", Sort: 20, Remark: "系统副标题", CreateTime: now, UpdateTime: now},
			{GroupKey: "basic", SettingKey: "basic.allow_register", SettingName: "允许注册", SettingValue: "false", ValueType: "switch", Sort: 30, Remark: "是否允许公开注册", CreateTime: now, UpdateTime: now},
			{GroupKey: "security", SettingKey: "security.login_failed_limit", SettingName: "登录失败阈值", SettingValue: "5", ValueType: "number", Sort: 10, Remark: "触发锁定前的失败次数", CreateTime: now, UpdateTime: now},
			{GroupKey: "security", SettingKey: "security.lock_minutes", SettingName: "锁定分钟数", SettingValue: "15", ValueType: "number", Sort: 20, Remark: "登录锁定时长", CreateTime: now, UpdateTime: now},
			{GroupKey: "security", SettingKey: "security.captcha_enabled", SettingName: "启用验证码", SettingValue: "true", ValueType: "switch", Sort: 30, Remark: "登录是否启用验证码", CreateTime: now, UpdateTime: now},
			{GroupKey: "notification", SettingKey: "notification.enable_site_notice", SettingName: "站内通知", SettingValue: "true", ValueType: "switch", Sort: 10, Remark: "开启站内通知中心", CreateTime: now, UpdateTime: now},
			{GroupKey: "notification", SettingKey: "notification.enable_email_notice", SettingName: "邮件通知", SettingValue: "false", ValueType: "switch", Sort: 20, Remark: "是否启用邮件通知", CreateTime: now, UpdateTime: now},
			{GroupKey: "notification", SettingKey: "notification.notice_retention_days", SettingName: "通知保留天数", SettingValue: "30", ValueType: "number", Sort: 30, Remark: "通知保留时间", CreateTime: now, UpdateTime: now},
		}
		if err := db.Create(&defaultSettings).Error; err != nil {
			return fmt.Errorf("create default settings failed: %w", err)
		}
	}

	// 6) 默认通知
	var noticeCount int64
	if err := db.Model(&entity.SysNotice{}).Count(&noticeCount).Error; err != nil {
		return fmt.Errorf("count notices failed: %w", err)
	}
	if noticeCount == 0 {
		now := utils.HTime{Time: time.Now()}
		welcomeNotice := entity.SysNotice{
			Title:         "欢迎使用后台管理系统",
			Content:       "系统已完成初始化，您可以开始配置权限、系统参数与通知策略。",
			NoticeType:    "system",
			NoticeLevel:   "info",
			Status:        2,
			TargetType:    "all",
			CreatedBy:     1,
			CreatedByName: "system",
			PublishTime:   now,
			CreateTime:    now,
			UpdateTime:    now,
		}
		if err := db.Create(&welcomeNotice).Error; err != nil {
			return fmt.Errorf("create welcome notice failed: %w", err)
		}
	}

	return nil
}
