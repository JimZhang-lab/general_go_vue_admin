/*
 * @Author: JimZhang
 * @Date: 2026-04-19 14:26:57
 * @LastEditors: JimZhang
 * @LastEditTime: 2026-04-19 14:26:57
 * @FilePath: /server/api/controller/sys_dict_controller.go
 * @Description:
 *
 */
package controller

import (
	"server/api/dao"
	"server/api/entity"
	"server/common/result"

	"github.com/gin-gonic/gin"
)

var sysDictTypeDao = &dao.SysDictTypeDao{}
var sysDictDataDao = &dao.SysDictDataDao{}

// GetDictTypeList 获取字典类型列表
func GetDictTypeList(c *gin.Context) {
	var params struct {
		PageNum  int    `form:"pageNum"`
		PageSize int    `form:"pageSize"`
		DictName string `json:"dictName" form:"dictName"`
		DictType string `json:"dictType" form:"dictType"`
		Status   string `json:"status" form:"status"`
	}
	if err := c.ShouldBind(&params); err != nil {
		result.Failed(c, 400, "参数错误")
		return
	}

	if params.PageNum <= 0 {
		params.PageNum = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	list, total, err := sysDictTypeDao.GetDictTypeList(params.DictName, params.DictType, params.Status, params.PageNum, params.PageSize)
	if err != nil {
		result.Failed(c, 500, "获取数据失败")
		return
	}
	result.SuccessWithPage(c, list, total, params.PageNum, params.PageSize)
}

func AddDictType(c *gin.Context) {
	var dt entity.SysDictType
	if err := c.ShouldBindJSON(&dt); err != nil {
		result.Failed(c, 400, "参数错误")
		return
	}
	if err := sysDictTypeDao.AddDictType(&dt); err != nil {
		result.Failed(c, 500, "操作失败")
		return
	}
	result.Success(c, nil)
}

func UpdateDictType(c *gin.Context) {
	var dt entity.SysDictType
	if err := c.ShouldBindJSON(&dt); err != nil {
		result.Failed(c, 400, "参数错误")
		return
	}
	if err := sysDictTypeDao.UpdateDictType(&dt); err != nil {
		result.Failed(c, 500, "操作失败")
		return
	}
	result.Success(c, nil)
}

func DeleteDictType(c *gin.Context) {
	var params struct{ Ids []int `json:"ids"` }
	if err := c.ShouldBindJSON(&params); err != nil || len(params.Ids) == 0 {
		result.Failed(c, 400, "参数错误")
		return
	}
	if err := sysDictTypeDao.DeleteDictType(params.Ids); err != nil {
		result.Failed(c, 500, "删除失败")
		return
	}
	result.Success(c, nil)
}

// GetDictDataList 获取字典数据列表
func GetDictDataList(c *gin.Context) {
	var params struct {
		PageNum   int    `form:"pageNum"`
		PageSize  int    `form:"pageSize"`
		DictType  string `json:"dictType" form:"dictType"`
		DictLabel string `json:"dictLabel" form:"dictLabel"`
		Status    string `json:"status" form:"status"`
	}
	if err := c.ShouldBind(&params); err != nil {
		result.Failed(c, 400, "参数错误")
		return
	}

	if params.PageNum <= 0 {
		params.PageNum = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	list, total, err := sysDictDataDao.GetDictDataList(params.DictType, params.DictLabel, params.Status, params.PageNum, params.PageSize)
	if err != nil {
		result.Failed(c, 500, "获取数据失败")
		return
	}
	result.SuccessWithPage(c, list, total, params.PageNum, params.PageSize)
}

func AddDictData(c *gin.Context) {
	var dd entity.SysDictData
	if err := c.ShouldBindJSON(&dd); err != nil {
		result.Failed(c, 400, "参数错误")
		return
	}
	if err := sysDictDataDao.AddDictData(&dd); err != nil {
		result.Failed(c, 500, "操作失败")
		return
	}
	result.Success(c, nil)
}

func UpdateDictData(c *gin.Context) {
	var dd entity.SysDictData
	if err := c.ShouldBindJSON(&dd); err != nil {
		result.Failed(c, 400, "参数错误")
		return
	}
	if err := sysDictDataDao.UpdateDictData(&dd); err != nil {
		result.Failed(c, 500, "操作失败")
		return
	}
	result.Success(c, nil)
}

func DeleteDictData(c *gin.Context) {
	var params struct{ Ids []int `json:"ids"` }
	if err := c.ShouldBindJSON(&params); err != nil || len(params.Ids) == 0 {
		result.Failed(c, 400, "参数错误")
		return
	}
	if err := sysDictDataDao.DeleteDictData(params.Ids); err != nil {
		result.Failed(c, 500, "删除失败")
		return
	}
	result.Success(c, nil)
}
