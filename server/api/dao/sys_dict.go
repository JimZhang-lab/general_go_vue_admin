/*
 * @Author: JimZhang
 * @Date: 2026-04-19 14:26:57
 * @LastEditors: JimZhang
 * @LastEditTime: 2026-04-19 14:26:57
 * @FilePath: /server/api/dao/sys_dict.go
 * @Description: 
 * 
 */
package dao

import (
	"server/api/entity"
	"server/pkg/db"
)

type SysDictTypeDao struct{}

func (s *SysDictTypeDao) GetDictTypeList(dictName, dictType, status string, pageNum, pageSize int) ([]*entity.SysDictType, int64, error) {
	var list []*entity.SysDictType
	var total int64
	query := db.Db.Model(&entity.SysDictType{})
	if dictName != "" {
		query = query.Where("dict_name LIKE ?", "%"+dictName+"%")
	}
	if dictType != "" {
		query = query.Where("dict_type LIKE ?", "%"+dictType+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = query.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (s *SysDictTypeDao) AddDictType(dt *entity.SysDictType) error {
	return db.Db.Create(dt).Error
}

func (s *SysDictTypeDao) UpdateDictType(dt *entity.SysDictType) error {
	return db.Db.Updates(dt).Error
}

func (s *SysDictTypeDao) DeleteDictType(ids []int) error {
	return db.Db.Delete(&entity.SysDictType{}, "id IN ?", ids).Error
}

type SysDictDataDao struct{}

func (s *SysDictDataDao) GetDictDataList(dictType, dictLabel, status string, pageNum, pageSize int) ([]*entity.SysDictData, int64, error) {
	var list []*entity.SysDictData
	var total int64
	query := db.Db.Model(&entity.SysDictData{})
	if dictType != "" {
		query = query.Where("dict_type = ?", dictType)
	}
	if dictLabel != "" {
		query = query.Where("dict_label LIKE ?", "%"+dictLabel+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = query.Order("dict_sort ASC").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (s *SysDictDataDao) AddDictData(dd *entity.SysDictData) error {
	return db.Db.Create(dd).Error
}

func (s *SysDictDataDao) UpdateDictData(dd *entity.SysDictData) error {
	return db.Db.Updates(dd).Error
}

func (s *SysDictDataDao) DeleteDictData(ids []int) error {
	return db.Db.Delete(&entity.SysDictData{}, "id IN ?", ids).Error
}
