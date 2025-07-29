package test

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"server/common/config"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SysAdmin struct {
	ID         uint      `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	PostId     int       `gorm:"column:post_id;comment:'岗位id'" json:"postId"`
	DeptId     int       `gorm:"column:dept_id;comment:'部门id'" json:"deptId"`
	Username   string    `gorm:"column:username;varchar(64);comment:'用户账号';NOT NULL" json:"username"`
	Password   string    `gorm:"column:password;varchar(64);comment:'密码';NOT NULL" json:"password"`
	Nickname   string    `gorm:"column:nickname;varchar(64);comment:'昵称'" json:"nickname"`
	Status     int       `gorm:"column:status;default:1;comment:'帐号启用状态：1->启用,2->禁用';NOT NULL" json:"status"`
	Icon       string    `gorm:"column:icon;varchar(500);comment:'头像'" json:"icon"`
	Email      string    `gorm:"column:email;varchar(64);comment:'邮箱'" json:"email"`
	Phone      string    `gorm:"column:phone;varchar(64);comment:'电话'" json:"phone"`
	Note       string    `gorm:"column:note;varchar(500);comment:'备注'" json:"note"`
	CreateTime time.Time `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
}

func (SysAdmin) TableName() string {
	return "sys_admin"
}

type SysAdminRole struct {
	ID      uint `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	AdminId uint `gorm:"column:admin_id;comment:'用户id'" json:"adminId"`
	RoleId  uint `gorm:"column:role_id;comment:'角色id'" json:"roleId"`
}

func (SysAdminRole) TableName() string {
	return "sys_admin_role"
}

func encryptionMd5(s string) string {
	ctx := md5.New()
	ctx.Write([]byte(s))
	return hex.EncodeToString(ctx.Sum(nil))
}

func CreateUser() {
	// 连接数据库
	var dbConfig = config.Config.DB
	url := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local&timeout=10s&collation=utf8mb4_general_ci",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
		dbConfig.Charset)

	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 检查是否已存在admin用户
	var existingUser SysAdmin
	result := db.Where("username = ?", "admin").First(&existingUser)
	if result.Error == nil {
		fmt.Println("用户 admin 已存在")
		fmt.Printf("用户名: %s, 昵称: %s, 状态: %d\n", existingUser.Username, existingUser.Nickname, existingUser.Status)
		return
	}

	// 创建测试用户
	testUser := SysAdmin{
		PostId:     1,
		DeptId:     1,
		Username:   "admin",
		Password:   encryptionMd5("admin123"), // 密码: admin123
		Nickname:   "系统管理员",
		Status:     1, // 启用状态
		Email:      "admin@example.com",
		Phone:      "13800138000",
		Note:       "系统管理员账号",
		CreateTime: time.Now(),
	}

	// 插入用户
	if err := db.Create(&testUser).Error; err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
		return
	}

	// 创建用户角色关联（假设角色ID为1）
	userRole := SysAdminRole{
		AdminId: testUser.ID,
		RoleId:  1,
	}

	if err := db.Create(&userRole).Error; err != nil {
		fmt.Printf("创建用户角色关联失败: %v\n", err)
		return
	}

	fmt.Println("测试用户创建成功!")
	fmt.Printf("用户名: %s\n", testUser.Username)
	fmt.Printf("密码: admin123\n")
	fmt.Printf("昵称: %s\n", testUser.Nickname)
	fmt.Printf("用户ID: %d\n", testUser.ID)
}

func TestCreateCatalog(t *testing.T) {
	fmt.Println("TestCreateCatalog")
	CreateUser()
}
