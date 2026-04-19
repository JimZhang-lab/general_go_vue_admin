package dao

import (
	"server/pkg/db"
	"time"
)

type DashboardStatsDto struct {
	AdminCount         int64   `json:"adminCount"`
	RoleCount          int64   `json:"roleCount"`
	PermissionCount    int64   `json:"permissionCount"`
	OnlineCount        int64   `json:"onlineCount"`
	LoginSuccessRate   float64 `json:"loginSuccessRate"`
	EnabledRate        float64 `json:"enabledRate"`
	OnlineRate         float64 `json:"onlineRate"`
	RecentActivities   []interface{} `json:"recentActivities"`
}

func GetDashboardStats() (DashboardStatsDto, error) {
	var stats DashboardStatsDto

	// Count total admins
	db.Db.Table("sys_admin").Count(&stats.AdminCount)

	var enabledAdminCount int64
	db.Db.Table("sys_admin").Where("status = ?", 1).Count(&enabledAdminCount)
	if stats.AdminCount > 0 {
		stats.EnabledRate = float64(enabledAdminCount) / float64(stats.AdminCount) * 100
	}

	// Count roles
	db.Db.Table("sys_role").Count(&stats.RoleCount)

	// Count permissions
	db.Db.Table("sys_menu").Count(&stats.PermissionCount)

	// Count online users (logged in within 30 minutes)
	time30MinsAgo := time.Now().Add(-30 * time.Minute)

	// Use Distinct to count unique users who logged in within 30 min
	db.Db.Table("sys_login_info").Where("login_status = ? AND login_time >= ?", 1, time30MinsAgo).Select("count(distinct username)").Row().Scan(&stats.OnlineCount)

	if stats.AdminCount > 0 {
		stats.OnlineRate = float64(stats.OnlineCount) / float64(stats.AdminCount) * 100
	}

	// Calculate overall login success rate
	var totalLogins int64
	var successLogins int64
	db.Db.Table("sys_login_info").Count(&totalLogins)
	db.Db.Table("sys_login_info").Where("login_status = ?", 1).Count(&successLogins)

	if totalLogins > 0 {
		stats.LoginSuccessRate = float64(successLogins) / float64(totalLogins) * 100
	}

	// Get 8 most recent activities from sys_operation_log
	var recentLogs []map[string]interface{}
	err := db.Db.Table("sys_operation_log").Select("id, username as user, CONCAT(method, ' ', url) as action, create_time as time").Order("create_time desc").Limit(8).Find(&recentLogs).Error

	if len(recentLogs) == 0 {
		// Fallback to login logs if no op logs
		db.Db.Table("sys_login_info").Select("id, username as user, IF(login_status=1, '用户登录成功', '用户登录失败') as action, login_time as time").Order("login_time desc").Limit(8).Find(&recentLogs)
	}

	stats.RecentActivities = make([]interface{}, len(recentLogs))
	for i, v := range recentLogs {
		stats.RecentActivities[i] = v
	}

	return stats, err
}
