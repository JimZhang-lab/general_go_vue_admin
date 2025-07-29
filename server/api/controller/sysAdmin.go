/*
 * @Author: JimZhang
 * @Date: 2025-07-25 11:12:49
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-26 00:36:14
 * @FilePath: /server/api/controller/sysAdmin.go
 * @Description: 用户控制层
 *
 */
package controller

import (
	"crypto/rand"
	"encoding/hex"
	"server/api/entity"
	"server/api/service"
	"server/common/controller"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SysAdminController 系统管理员控制器
type SysAdminController struct {
	controller.BaseController
}

// generateTraceID 生成追踪ID
func generateTraceID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// @Summary 用户登录接口
// @Produce json
// @Description 用户登录接口，支持JSON、表单、URI多种数据格式
// @Param data body entity.LoginDto true "data"
// @Success 200 {object} result.Result
// @router /api/login [post]
func Login(c *gin.Context) {
	ctrl := &SysAdminController{}

	var dto entity.LoginDto
	if err := ctrl.BindRequest(c, &dto); err != nil {
		ctrl.FailedWithError(c, err)
		return
	}

	// 设置追踪ID
	traceID := ctrl.GetTraceID(c)
	if traceID == "" {
		traceID = generateTraceID()
		ctrl.SetTraceID(c, traceID)
	}

	// 记录请求日志
	ctrl.LogRequest(c, "Login", dto)

	// 调用优化后的服务层
	service.SysAdminService().LoginLegacy(c, dto)
}

// @Summary 用户登出接口
// @Produce json
// @Description 用户登出接口
// @Success 200 {object} result.Result
// @router /api/logout [post]
// @Security ApiKeyAuth
func Logout(c *gin.Context) {
	service.SysAdminService().Logout(c)
}

// 新增用户
// @Summary 新增用户接口
// @Produce json
// @Description 新增用户接口，支持多种数据格式
// @Param data body entity.AddSysAdminDto true "data"
// @Success 200 {object} result.Result
// @router /api/admin/add [post]
// @Security ApiKeyAuth
func CreateSysAdmin(c *gin.Context) {
	ctrl := &SysAdminController{}

	var dto entity.AddSysAdminDto
	if err := ctrl.BindRequest(c, &dto); err != nil {
		ctrl.FailedWithError(c, err)
		return
	}

	// 设置追踪ID
	traceID := ctrl.GetTraceID(c)
	if traceID == "" {
		traceID = generateTraceID()
		ctrl.SetTraceID(c, traceID)
	}

	ctrl.LogRequest(c, "CreateSysAdmin", dto)
	service.SysAdminService().CreateSysAdmin(c, dto)
}

// 根据id查询用户
// @Summary 根据id查询用户接口
// @Produce json
// @Description 根据id查询用户接口
// @Param id query int true "Id"
// @Success 200 {object} result.Result
// @router /api/admin/info [get]
// @Security ApiKeyAuth
func GetSysAdminInfo(c *gin.Context) {
	Id, _ := strconv.Atoi(c.Query("id"))
	service.SysAdminService().GetSysAdminInfo(c, Id)
}

// 修改用户
// @Summary 修改用户接口
// @Produce json
// @Description 修改用户接口
// @Param data body entity.UpdateSysAdminDto true "data"
// @Success 200 {object} result.Result
// @router /api/admin/update [put]
// @Security ApiKeyAuth
func UpdateSysAdmin(c *gin.Context) {
	var dto entity.UpdateSysAdminDto
	_ = c.BindJSON(&dto)
	service.SysAdminService().UpdateSysAdmin(c, dto)
}

// 根据id删除用户
// @Summary 根据id删除接口
// @Produce json
// @Description 根据id删除接口
// @Param data body entity.SysAdminIdDto true "data"
// @Success 200 {object} result.Result
// @router /api/admin/delete [delete]
// @Security ApiKeyAuth
func DeleteSysAdminById(c *gin.Context) {
	var dto entity.SysAdminIdDto
	_ = c.BindJSON(&dto)
	service.SysAdminService().DeleteSysAdminById(c, dto)
}

//	用户状态启用/停用
//
// @Summary 用户状态启用/停用接口
// @Produce json
// @Description 用户状态启用/停用接口
// @Param data body entity.UpdateSysAdminStatusDto true "data"
// @Success 200 {object} result.Result
// @router /api/admin/updateStatus [put]
// @Security ApiKeyAuth
func UpdateSysAdminStatus(c *gin.Context) {
	var dto entity.UpdateSysAdminStatusDto
	_ = c.BindJSON(&dto)
	service.SysAdminService().UpdateSysAdminStatus(c, dto)
}

// 重置密码
// @Summary 重置密码接口
// @Produce json
// @Description 重置密码接口
// @Param data body entity.ResetSysAdminPasswordDto true "data"
// @Success 200 {object} result.Result
// @router /api/admin/updatePassword [put]
// @Security ApiKeyAuth
func ResetSysAdminPassword(c *gin.Context) {
	var dto entity.ResetSysAdminPasswordDto
	_ = c.BindJSON(&dto)
	service.SysAdminService().ResetSysAdminPassword(c, dto)
}

// 分页获取用户列表
// @Summary 分页获取用户列表接口
// @Produce json
// @Description 分页获取用户列表接口
// @Param pageNum query int false "分页数"
// @Param pageSize query int false "每页数"
// @Param username query string false "用户名"
// @Param status query string false "帐号启用状态：1->启用,2->禁用"
// @Param beginTime query string false "开始时间"
// @Param endTime query string false "结束时间"
// @Success 200 {object} result.Result
// @router /api/admin/list [get]
// @Security ApiKeyAuth
func GetSysAdminList(c *gin.Context) {
	PageNum, _ := strconv.Atoi(c.Query("pageNum"))
	PageSize, _ := strconv.Atoi(c.Query("pageSize"))
	Username := c.Query("username")
	Status := c.Query("status")
	BeginTime := c.Query("beginTime")
	EndTime := c.Query("endTime")
	service.SysAdminService().GetSysAdminList(c, PageSize, PageNum, Username, Status, BeginTime, EndTime)
}

// 修改个人信息
// @Summary 修改个人信息接口
// @Produce json
// @Description 修改个人信息接口
// @Param data body entity.UpdatePersonalDto true "data"
// @Success 200 {object} result.Result
// @router /api/admin/updatePersonal [put]
// @Security ApiKeyAuth
func UpdatePersonal(c *gin.Context) {
	var dto entity.UpdatePersonalDto
	_ = c.BindJSON(&dto)
	service.SysAdminService().UpdatePersonal(c, dto)
}

// 修改密码
// @Summary 修改密码接口
// @Produce json
// @Description 修改密码接口
// @Param data body entity.UpdatePersonalPasswordDto true "data"
// @Success 200 {object} result.Result
// @router /api/admin/updatePersonalPassword [put]
// @Security ApiKeyAuth
func UpdatePersonalPassword(c *gin.Context) {
	var dto entity.UpdatePersonalPasswordDto
	_ = c.BindJSON(&dto)
	service.SysAdminService().UpdatePersonalPassword(c, dto)
}
