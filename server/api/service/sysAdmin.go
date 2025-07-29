/*
 * @Author: JimZhang
 * @Date: 2025-07-25 00:44:50
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-28 11:21:13
 * @FilePath: /go-vue-general-admin/server/api/service/sysAdmin.go
 * @Description:
 *
 */
package service

import (
	"context"
	"fmt"
	"server/api/dao"
	"server/api/entity"
	"server/common/errors"
	"server/common/result"
	"server/common/service"
	"server/common/utils"
	"server/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 定义接口 - 优化后支持错误链式传递
type ISysAdminService interface {
	Login(ctx context.Context, dto entity.LoginDto) *service.ServiceResult
	CreateSysAdmin(ctx context.Context, dto entity.AddSysAdminDto) *service.ServiceResult
	GetSysAdminInfo(ctx context.Context, id int) *service.ServiceResult
	UpdateSysAdmin(ctx context.Context, dto entity.UpdateSysAdminDto) *service.ServiceResult
	DeleteSysAdminById(ctx context.Context, dto entity.SysAdminIdDto) *service.ServiceResult
	UpdateSysAdminStatus(ctx context.Context, dto entity.UpdateSysAdminStatusDto) *service.ServiceResult
	ResetSysAdminPassword(ctx context.Context, dto entity.ResetSysAdminPasswordDto) *service.ServiceResult
	GetSysAdminList(ctx context.Context, pageSize, pageNum int, username, status, beginTime, endTime string) *service.ServiceResult
	UpdatePersonal(ctx context.Context, dto entity.UpdatePersonalDto) *service.ServiceResult
	UpdatePersonalPassword(ctx context.Context, dto entity.UpdatePersonalPasswordDto) *service.ServiceResult
}

// SysAdminServiceImpl 系统管理员服务实现
type SysAdminServiceImpl struct {
	*service.BaseService
}

// NewSysAdminService 创建系统管理员服务实例（暂时注释，保持向后兼容）
// func NewSysAdminService() ISysAdminService {
// 	return &SysAdminServiceImpl{
// 		BaseService: service.NewBaseService(),
// 	}
// }

// 修改个人密码
func (s SysAdminServiceImpl) UpdatePersonalPassword(c *gin.Context, dto entity.UpdatePersonalPasswordDto) {
	err := validator.New().Struct(dto)
	if err != nil {
		result.Failed(c, int(result.ApiCode.MissingChangePasswordParameter), result.ApiCode.GetMessage(result.ApiCode.MissingChangePasswordParameter))
		return
	}
	sysAdmin, _ := jwt.GetAdmin(c)
	dto.Id = sysAdmin.ID
	sysAdminExist := dao.GetSysAdminByUsername(sysAdmin.Username)
	if sysAdminExist.Password != utils.EncryptionMd5(dto.Password) {
		result.Failed(c, int(result.ApiCode.PASSWORDNOTTRUE), result.ApiCode.GetMessage(result.ApiCode.PASSWORDNOTTRUE))
		return
	}
	if dto.NewPassword != dto.ResetPassword {
		result.Failed(c, int(result.ApiCode.RESETPASSWORD), result.ApiCode.GetMessage(result.ApiCode.RESETPASSWORD))
		return
	}
	dto.NewPassword = utils.EncryptionMd5(dto.NewPassword)
	sysAdminUpdatePwd := dao.UpdatePersonalPassword(dto)
	tokenString, _ := jwt.GenerateTokenByAdmin(sysAdminUpdatePwd)
	result.Success(c, map[string]interface{}{"token": tokenString, "sysAdmin": sysAdminUpdatePwd})
	return
}

// 修改个人信息
func (s SysAdminServiceImpl) UpdatePersonal(c *gin.Context, dto entity.UpdatePersonalDto) {
	err := validator.New().Struct(dto)
	if err != nil {
		result.Failed(c, int(result.ApiCode.MissingModificationOfPersonalParameters), result.ApiCode.GetMessage(result.ApiCode.MissingModificationOfPersonalParameters))
		return
	}
	id, _ := jwt.GetAdminId(c)
	dto.Id = id
	result.Success(c, dao.UpdatePersonal(dto))
}

// 分页查询用户列表
func (s SysAdminServiceImpl) GetSysAdminList(c *gin.Context, PageSize, PageNum int, Username, Status, BeginTime, EndTime string) {
	if PageSize < 1 {
		PageSize = 10
	}
	if PageNum < 1 {
		PageNum = 1
	}
	sysAdmin, count := dao.GetSysAdminList(PageSize, PageNum, Username, Status, BeginTime, EndTime)
	result.Success(c, map[string]interface{}{"total": count, "pageSize": PageSize, "pageNum": PageNum, "list": sysAdmin})
	return
}

// 重置密码
func (s SysAdminServiceImpl) ResetSysAdminPassword(c *gin.Context, dto entity.ResetSysAdminPasswordDto) {
	dao.ResetSysAdminPassword(dto)
	result.Success(c, true)
}

// 修改用户状态
func (s SysAdminServiceImpl) UpdateSysAdminStatus(c *gin.Context, dto entity.UpdateSysAdminStatusDto) {
	dao.UpdateSysAdminStatus(dto)
	result.Success(c, true)
}

// 根据id删除用户
func (s SysAdminServiceImpl) DeleteSysAdminById(c *gin.Context, dto entity.SysAdminIdDto) {
	dao.DeleteSysAdminById(dto)
	result.Success(c, true)
}

// 修改用户
func (s SysAdminServiceImpl) UpdateSysAdmin(c *gin.Context, dto entity.UpdateSysAdminDto) {
	result.Success(c, dao.UpdateSysAdmin(dto))
}

// 根据id查询用户信息
func (s SysAdminServiceImpl) GetSysAdminInfo(c *gin.Context, Id int) {
	result.Success(c, dao.GetSysAdminInfo(Id))
}

// 新增用户
func (s SysAdminServiceImpl) CreateSysAdmin(c *gin.Context, dto entity.AddSysAdminDto) {
	err := validator.New().Struct(dto)
	if err != nil {
		result.Failed(c, int(result.ApiCode.MissingNewAdminParameter), result.ApiCode.GetMessage(result.ApiCode.MissingNewAdminParameter))
		return
	}
	bool := dao.CreateSysAdmin(dto)
	if !bool {
		result.Failed(c, int(result.ApiCode.USERNAMEALREADYEXISTS), result.ApiCode.GetMessage(result.ApiCode.USERNAMEALREADYEXISTS))
		return
	}
	result.Success(c, bool)
	return
}

// Login 用户登录 - 优化版本支持错误链式传递
func (s *SysAdminServiceImpl) Login(ctx context.Context, dto entity.LoginDto) *service.ServiceResult {
	// 参数验证
	if err := s.ValidateStruct(dto); err != nil {
		return service.NewServiceResult(nil, err)
	}

	// 从上下文获取客户端IP（如果是从gin.Context转换的）
	var ip string
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ip = ginCtx.ClientIP()
	} else {
		ip = "unknown"
	}

	// 并行执行验证码检查和用户信息查询
	var code string
	var sysAdmin entity.SysAdmin

	err := s.ParallelExecute(
		// 验证码检查任务
		func() error {
			code = utils.RedisStore{}.Get(dto.IdKey, true)
			if len(code) == 0 {
				dao.CreateSysLoginInfo(dto.Username, ip, utils.GetRealAddressByIP(ip), "", "", "验证码已过期", 2)
				return errors.AuthenticationError("验证码已过期")
			}

			// 添加调试日志
			fmt.Printf("🔍 验证码调试: ID=%s, 存储值=%s, 输入值=%s\n", dto.IdKey, code, dto.Image)

			// 校验验证码
			if !CaptVerify(dto.IdKey, dto.Image) {
				dao.CreateSysLoginInfo(dto.Username, ip, utils.GetRealAddressByIP(ip), "", "", "验证码不正确", 2)
				return errors.AuthenticationError("验证码不正确")
			}
			return nil
		},
		// 用户信息查询任务
		func() error {
			sysAdmin = dao.SysAdminDetail(dto)
			if sysAdmin.ID == 0 {
				return errors.NotFoundError("用户")
			}
			return nil
		},
	)

	if err.HasErrors() {
		return service.NewServiceResult(nil, err.First())
	}

	// 密码验证
	if sysAdmin.Password != utils.EncryptionMd5(dto.Password) {
		dao.CreateSysLoginInfo(dto.Username, ip, utils.GetRealAddressByIP(ip), "", "", "密码不正确", 2)
		return service.NewServiceResult(nil, errors.AuthenticationError("密码不正确"))
	}

	// 账号状态检查
	const disabledStatus int = 2
	if sysAdmin.Status == disabledStatus {
		dao.CreateSysLoginInfo(dto.Username, ip, utils.GetRealAddressByIP(ip), "", "", "账号已停用", 2)
		return service.NewServiceResult(nil, errors.AuthenticationError("账号已停用"))
	}

	// 生成token和构建响应数据
	var tokenString string
	var leftMenuVo []entity.LeftMenuVo
	var permissionList []entity.ValueVo

	// 并行执行token生成和菜单权限查询
	err = s.ParallelExecute(
		// Token生成任务
		func() error {
			var genErr error
			tokenString, genErr = jwt.GenerateTokenByAdmin(sysAdmin)
			if genErr != nil {
				return errors.Wrap(genErr, errors.ErrInternal, "Token生成失败")
			}
			return nil
		},
		// 菜单查询任务
		func() error {
			leftMenuList := dao.QueryLeftMenuList(sysAdmin.ID)
			for _, value := range leftMenuList {
				menuSvoList := dao.QueryMenuVoList(sysAdmin.ID, value.Id)
				item := entity.LeftMenuVo{
					Id:          value.Id,
					MenuName:    value.MenuName,
					Icon:        value.Icon,
					Url:         value.Url,
					MenuSvoList: menuSvoList,
				}
				leftMenuVo = append(leftMenuVo, item)
			}
			return nil
		},
		// 权限查询任务
		func() error {
			permissionList = dao.QueryPermissionList(sysAdmin.ID)
			return nil
		},
	)

	if err.HasErrors() {
		return service.NewServiceResult(nil, err.First())
	}

	// 记录登录成功日志
	dao.CreateSysLoginInfo(dto.Username, ip, utils.GetRealAddressByIP(ip), "", "", "登录成功", 1)

	// 构建权限字符串列表
	var stringList = make([]string, 0, len(permissionList))
	for _, value := range permissionList {
		stringList = append(stringList, value.Value)
	}

	// 构建响应数据
	loginResult := gin.H{
		"token":          tokenString,
		"sysAdmin":       sysAdmin,
		"leftMenuList":   leftMenuVo,
		"permissionList": stringList,
	}

	return service.NewServiceResult(loginResult, nil)
}

// LoginLegacy 保持向后兼容的登录方法
func (s *SysAdminServiceImpl) LoginLegacy(c *gin.Context, dto entity.LoginDto) {
	serviceResult := s.Login(c, dto)
	if serviceResult.IsSuccess() {
		result.Success(c, serviceResult.Data)
	} else {
		result.FailedWithError(c, serviceResult.Error)
	}
}

// Logout 用户登出
func (s *SysAdminServiceImpl) Logout(c *gin.Context) {
	// 从请求头中获取token
	token := c.GetHeader("Authorization")
	if token != "" {
		// 移除Bearer前缀
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		// 这里可以将token加入黑名单或者记录登出日志
		// 由于JWT是无状态的，我们主要是记录登出操作

		// 获取用户信息
		if claims, err := jwt.ValidateToken(token); err == nil {
			// 记录登出日志
			dao.CreateSysLoginInfo(
				claims.Username,
				c.ClientIP(),
				utils.GetRealAddressByIP(c.ClientIP()),
				utils.GetBrowser(c),
				utils.GetOs(c),
				"登出成功",
				1, // 成功
			)
		}
	}

	result.Success(c, "登出成功")
}

// 临时保持向后兼容的服务实例
var sysAdminService = &SysAdminServiceImpl{
	BaseService: service.NewBaseService(),
}

// SysAdminService 返回服务实例（保持向后兼容）
func SysAdminService() *SysAdminServiceImpl {
	return sysAdminService
}
