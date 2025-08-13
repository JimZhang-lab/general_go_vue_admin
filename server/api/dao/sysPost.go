/*
 * @Author: JimZhang
 * @Date: 2025-07-25 11:30:31
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-25 11:34:33
 * @FilePath: /server/api/dao/sysPost.go
 * @Description: 系统岗位
 *
 */
package dao

import (
	"server/api/entity"
	"server/common/utils"
	"server/pkg/db"
	"time"
)

func GetSysPostByCode(postCode string) (sysPost entity.SysPost) {
	db.Db.Where("post_code = ?", postCode).First(&sysPost)
	return sysPost
}

// 根据名称查询
func GetSysPostByName(postName string) (sysPost entity.SysPost) {
	db.Db.Where("post_name = ?", postName).First(&sysPost)
	return sysPost
}

// 新增岗位
func CreateSysPost(sysPost entity.SysPost) bool {
	sysPostByCode := GetSysPostByCode(sysPost.PostCode)
	if sysPostByCode.ID > 0 {
		return false
	}
	sysPostByName := GetSysPostByName(sysPost.PostName)
	if sysPostByName.ID > 0 {
		return false
	}
	addSysPost := entity.SysPost{
		PostCode:   sysPost.PostCode,
		PostName:   sysPost.PostName,
		PostStatus: sysPost.PostStatus,
		CreateTime: utils.HTime{Time: time.Now()},
		Remark:     sysPost.Remark,
	}
	tx := db.Db.Save(&addSysPost)

	return tx.RowsAffected > 0
}

// 分页查询岗位列表
func GetSysPostList(PageNum, PageSize int, PostName, PostStatus, BeginTime, EndTime string) (sysPost []entity.SysPost, count int64) {
	curDb := db.Db.Table("sys_post")
	if PostName != "" {
		curDb = curDb.Where("post_name = ?", PostName)
	}
	if PostStatus != "" {
		curDb = curDb.Where("post_status = ?", PostStatus)
	}
	if BeginTime != "" && EndTime != "" {
		curDb = curDb.Where("`create_time` BETWEEN ? AND ?", BeginTime, EndTime)
	}
	curDb.Count(&count)
	curDb.Limit(PageSize).Offset((PageNum - 1) * PageSize).Order("create_time desc").Find(&sysPost)
	return sysPost, count
}

// 根据id查询岗位
func GetSysPostById(Id int) (sysPost entity.SysPost) {
	db.Db.First(&sysPost, Id)
	return sysPost
}

// 修改岗位
func UpdateSysPost(post entity.SysPost) (sysPost entity.SysPost) {
	db.Db.First(&sysPost, post.ID)
	sysPost.PostName = post.PostName
	sysPost.PostCode = post.PostCode
	sysPost.PostStatus = post.PostStatus
	if post.Remark != "" {
		sysPost.Remark = post.Remark
	}
	db.Db.Save(&sysPost)
	return sysPost
}

// 根据id删除岗位
func DeleteSysPostById(dto entity.SysPostIdDto) {
	db.Db.Delete(&entity.SysPost{}, dto.Id)
}

// 批量删除岗位
func BatchDeleteSysPost(dto entity.DelSysPostDto) {
	db.Db.Where("id in (?)", dto.Ids).Delete(&entity.SysPost{})
}

// 修改状态
func UpdateSysPostStatus(dto entity.UpdateSysPostStatusDto) {
	var sysPost entity.SysPost
	result := db.Db.First(&sysPost, dto.Id)
	if result.Error != nil {
		// 未找到记录或其他错误
		return
	}
	sysPost.PostStatus = dto.PostStatus
	saveResult := db.Db.Save(&sysPost)
	if saveResult.Error != nil {
		// 保存出错
		return
	}
}

// 岗位下拉列表
func QuerySysPostVoList() (sysPostVo []entity.SysPostVo) {
	db.Db.Table("sys_post").Select("id, post_name").Scan(&sysPostVo)
	return sysPostVo
}
