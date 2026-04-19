/*
 * @Author: JimZhang
 * @Date: 2025-07-25 20:40:00
 * @LastEditors: JimZhang
 * @LastEditTime: 2025-07-25 20:45:00
 * @FilePath: /server/api/service/base_service.go
 * @Description: 基础服务类，封装通用的服务层逻辑
 *
 */
package service

import (
	"server/common/result"
	"server/pkg/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CheckPageParams 检查分页参数
func CheckPageParams(pageSize, pageNum int) (int, int) {
	if pageSize < 1 {
		pageSize = 10
	}
	if pageNum < 1 {
		pageNum = 1
	}
	return pageSize, pageNum
}

// CheckEntityExists 检查实体是否存在
func CheckEntityExists(c *gin.Context, model interface{}, id int, errorMsg string) bool {
	findResult := db.Db.First(model, id)
	if findResult.Error != nil {
		if findResult.Error == gorm.ErrRecordNotFound {
			result.Failed(c, 400, errorMsg)
		} else {
			result.Failed(c, 500, "查询出错")
		}
		return false
	}
	return true
}

// UpdateEntityStatus 更新实体状态
func UpdateEntityStatus(c *gin.Context, model interface{}, id int, statusField string, statusValue interface{}, errorMsg string) {
	// 检查实体是否存在
	if !CheckEntityExists(c, model, id, errorMsg) {
		return
	}

	// 更新状态
	db.Db.Model(model).Update(statusField, statusValue)
	result.Success(c, true)
}

// BuildPageResult 构建分页结果
func BuildPageResult(c *gin.Context, data interface{}, count int64, pageSize, pageNum int) {
	result.Success(c, map[string]interface{}{
		"total":    count,
		"pageSize": pageSize,
		"pageNum":  pageNum,
		"list":     data,
	})
}

// CheckUpdateResult 检查更新结果
func CheckUpdateResult(c *gin.Context, saveResult *gorm.DB) bool {
	if saveResult.Error != nil || saveResult.RowsAffected == 0 {
		result.Failed(c, 500, "更新失败")
		return false
	}
	return true
}

// GetEntityById 根据ID获取实体
func GetEntityById(c *gin.Context, model interface{}, id int, notFoundMsg string) {
	findResult := db.Db.First(model, id)
	if findResult.Error != nil {
		if findResult.Error == gorm.ErrRecordNotFound {
			result.Failed(c, 400, notFoundMsg)
		} else {
			result.Failed(c, 500, "查询出错")
		}
		return
	}
	result.Success(c, model)
}

// DeleteEntityById 根据ID删除实体
func DeleteEntityById(c *gin.Context, model interface{}, id int, notFoundMsg string) {
	// 检查实体是否存在
	if !CheckEntityExists(c, model, id, notFoundMsg) {
		return
	}

	// 删除实体
	deleteResult := db.Db.Delete(model)
	if deleteResult.Error != nil || deleteResult.RowsAffected == 0 {
		result.Failed(c, 500, "删除失败")
		return
	}
	result.Success(c, true)
}

// CreateEntity 创建实体
func CreateEntity(c *gin.Context, model interface{}, checkFunc func() (bool, string)) {
	// 执行自定义检查
	if checkFunc != nil {
		if ok, msg := checkFunc(); !ok {
			result.Failed(c, 400, msg)
			return
		}
	}

	// 创建实体
	createResult := db.Db.Create(model)
	if createResult.Error != nil || createResult.RowsAffected == 0 {
		result.Failed(c, 500, "创建失败")
		return
	}
	result.Success(c, true)
}

// UpdateEntity 更新实体
func UpdateEntity(c *gin.Context, model interface{}, id int, notFoundMsg string) {
	// 检查实体是否存在
	if !CheckEntityExists(c, model, id, notFoundMsg) {
		return
	}

	// 更新实体
	saveResult := db.Db.Save(model)
	if !CheckUpdateResult(c, saveResult) {
		return
	}

	result.Success(c, true)
}
