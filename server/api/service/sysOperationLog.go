/*
 * @Author: JimZhang
 * @Date: 2026-04-19 14:26:57
 * @LastEditors: JimZhang
 * @LastEditTime: 2026-04-19 14:26:57
 * @FilePath: /server/api/service/sysOperationLog.go
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

type ISysOperationLogService interface {
	GetSysOperationLogList(c *gin.Context, Username, BeginTime, EndTime string, PageSize, PageNum int)
	DeleteSysOperationLogById(c *gin.Context, dto entity.SysOperationLogIdDto)
	BatchDeleteSysOperationLog(c *gin.Context, dto entity.BatchDeleteSysOperationLogDto)
	CleanSysOperationLog(c *gin.Context)
	ExportSysOperationLogList(c *gin.Context, Username, BeginTime, EndTime string)
}

type SysOperationLogServiceImpl struct{}

// 清空操作日志
func (s SysOperationLogServiceImpl) CleanSysOperationLog(c *gin.Context) {
	dao.CleanSysOperationLog()
	result.Success(c, true)
}

// 批量删除操作日志
func (s SysOperationLogServiceImpl) BatchDeleteSysOperationLog(c *gin.Context, dto entity.BatchDeleteSysOperationLogDto) {
	dao.BatchDeleteSysOperationLog(dto)
	result.Success(c, true)
}

// 根据id删除操作日志
func (s SysOperationLogServiceImpl) DeleteSysOperationLogById(c *gin.Context, dto entity.SysOperationLogIdDto) {
	dao.DeleteSysOperationLogById(dto)
	result.Success(c, true)
}

// 分页查询操作日志列表
func (s SysOperationLogServiceImpl) GetSysOperationLogList(c *gin.Context, Username, BeginTime, EndTime string, PageSize, PageNum int) {
	if PageSize < 1 {
		PageSize = 10
	}
	if PageNum < 1 {
		PageNum = 1
	}
	sysOperationLog, count := dao.GetSysOperationLogList(Username, BeginTime, EndTime, PageSize, PageNum)
	result.Success(c, map[string]interface{}{"total": count, "pageSize": PageSize, "pageNum": PageNum, "list": sysOperationLog})
}

// 导出操作日志
func (s SysOperationLogServiceImpl) ExportSysOperationLogList(c *gin.Context, Username, BeginTime, EndTime string) {
	// 一次性获取最大支持导出数量，比如10000条
	sysOperationLog, _ := dao.GetSysOperationLogList(Username, BeginTime, EndTime, 10000, 1)

	b := &bytes.Buffer{}
	// 写入BOM头，防止Excel乱码
	b.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(b)

	w.Write([]string{"ID", "用户名", "请求方法", "请求URL", "IP地址", "操作时间"})
	for _, row := range sysOperationLog {
		w.Write([]string{
			fmt.Sprint(row.ID),
			row.Username,
			row.Method,
			row.Url,
			row.Ip,
			row.CreateTime.Time.Format("2006-01-02 15:04:05"),
		})
	}
	w.Flush()

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment;filename=operation_logs.csv")
	c.Data(200, "text/csv;charset=utf-8", b.Bytes())
}

var sysOperationLogService = SysOperationLogServiceImpl{}

func SysOperationLogService() ISysOperationLogService {
	return &sysOperationLogService
}
