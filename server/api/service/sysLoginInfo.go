/*
 * @Author: JimZhang
 * @Date: 2026-04-19 14:26:57
 * @LastEditors: JimZhang
 * @LastEditTime: 2026-04-19 14:26:57
 * @FilePath: /server/api/service/sysLoginInfo.go
 * @Description:
 *
 */
package service

import (
	"server/api/dao"
	"server/api/entity"
	"server/common/result"
	"bytes"
	"encoding/csv"
	"fmt"

	"github.com/gin-gonic/gin"
)

type ISysLoginInfoService interface {
	GetSysLoginInfoList(c *gin.Context, Username, LoginStatus, BeginTime, EndTime string, PageSize, pageNum int)
	BatchDeleteSysLoginInfo(c *gin.Context, dto entity.DelSysLoginInfoDto)
	DeleteSysLoginInfo(c *gin.Context, dto entity.SysLoginInfoIdDto)
	CleanSysLoginInfo(c *gin.Context)
	ExportSysLoginInfoList(c *gin.Context, Username, LoginStatus, BeginTime, EndTime string)
}

type SysLoginInfoServiceImpl struct{}

// 清空登录日志
func (s SysLoginInfoServiceImpl) CleanSysLoginInfo(c *gin.Context) {
	dao.CleanSysLoginInfo()
	result.Success(c, true)
}

// 批量删除登录日志
func (s SysLoginInfoServiceImpl) BatchDeleteSysLoginInfo(c *gin.Context, dto entity.DelSysLoginInfoDto) {
	dao.BatchDeleteSysLoginInfo(dto)
	result.Success(c, true)
}

// 根据id删除登录日志
func (s SysLoginInfoServiceImpl) DeleteSysLoginInfo(c *gin.Context, dto entity.SysLoginInfoIdDto) {
	dao.DeleteSysLoginInfoById(dto)
	result.Success(c, true)
}

// 分页获取登录日志列表
func (s SysLoginInfoServiceImpl) GetSysLoginInfoList(c *gin.Context, Username, LoginStatus, BeginTime, EndTime string, PageSize, pageNum int) {
	if PageSize < 1 {
		PageSize = 10
	}
	if pageNum < 1 {
		pageNum = 1
	}
	sysLoginInfo, count := dao.GetSysLoginInfoList(Username, LoginStatus, BeginTime, EndTime, PageSize, pageNum)
	result.Success(c, map[string]interface{}{"total": count, "pageSize": PageSize, "pageNum": pageNum, "list": sysLoginInfo})
}

// 导出登录日志
func (s SysLoginInfoServiceImpl) ExportSysLoginInfoList(c *gin.Context, Username, LoginStatus, BeginTime, EndTime string) {
	sysLoginInfo, _ := dao.GetSysLoginInfoList(Username, LoginStatus, BeginTime, EndTime, 10000, 1)

	b := &bytes.Buffer{}
	b.WriteString("\xEF\xBB\xBF") // BOM
	w := csv.NewWriter(b)

	w.Write([]string{"ID", "用户名", "IP地址", "登录地点", "浏览器", "操作系统", "状态", "提示消息", "登录时间"})
	for _, row := range sysLoginInfo {
		statusStr := "成功"
		if row.LoginStatus != 1 {
			statusStr = "失败"
		}
		w.Write([]string{
			fmt.Sprint(row.ID),
			row.Username,
			row.IpAddress,
			row.LoginLocation,
			row.Browser,
			row.Os,
			statusStr,
			row.Message,
			row.LoginTime.Time.Format("2006-01-02 15:04:05"),
		})
	}
	w.Flush()

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment;filename=login_logs.csv")
	c.Data(200, "text/csv;charset=utf-8", b.Bytes())
}

var sysLoginInfoService = SysLoginInfoServiceImpl{}

func SysLoginInfoService() ISysLoginInfoService {
	return &sysLoginInfoService
}
