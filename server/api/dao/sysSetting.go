/*
 * @Author: JimZhang
 * @Date: 2026-04-19 14:26:57
 * @LastEditors: JimZhang
 * @LastEditTime: 2026-04-19 14:26:57
 * @FilePath: /server/api/dao/sysSetting.go
 * @Description: 
 * 
 */
package dao

import (
	"server/api/entity"
	"server/common/utils"
	"server/pkg/db"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

func GetSysSettingList(groupKey string) (settings []entity.SysSetting) {
	settings = make([]entity.SysSetting, 0)
	curDb := db.Db.Model(&entity.SysSetting{})
	if groupKey != "" {
		curDb = curDb.Where("group_key = ?", groupKey)
	}
	curDb.Order("sort ASC, id ASC").Find(&settings)
	return settings
}

func GetSysSettingByKey(settingKey string) (setting entity.SysSetting, err error) {
	err = db.Db.Where("setting_key = ?", settingKey).First(&setting).Error
	return setting, err
}

func GetSysSettingBoolValue(settingKey string, fallback bool) bool {
	setting, err := GetSysSettingByKey(settingKey)
	if err != nil {
		return fallback
	}
	value := strings.TrimSpace(strings.ToLower(setting.SettingValue))
	switch value {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}

func GetSysSettingIntValue(settingKey string, fallback int) int {
	setting, err := GetSysSettingByKey(settingKey)
	if err != nil {
		return fallback
	}
	value, convErr := strconv.Atoi(strings.TrimSpace(setting.SettingValue))
	if convErr != nil {
		return fallback
	}
	return value
}

func BatchUpsertSysSettings(items []entity.SysSettingItemDto) error {
	return db.Db.Transaction(func(tx *gorm.DB) error {
		now := utils.HTime{Time: time.Now()}
		for _, item := range items {
			var setting entity.SysSetting
			err := tx.Where("setting_key = ?", item.SettingKey).First(&setting).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}

			if err == gorm.ErrRecordNotFound {
				setting = entity.SysSetting{
					GroupKey:     item.GroupKey,
					SettingKey:   item.SettingKey,
					SettingName:  item.SettingName,
					SettingValue: item.SettingValue,
					ValueType:    item.ValueType,
					OptionsJSON:  item.OptionsJSON,
					IsEncrypted:  item.IsEncrypted,
					Sort:         item.Sort,
					Remark:       item.Remark,
					CreateTime:   now,
					UpdateTime:   now,
				}
				if createErr := tx.Create(&setting).Error; createErr != nil {
					return createErr
				}
				continue
			}

			setting.GroupKey = item.GroupKey
			setting.SettingName = item.SettingName
			setting.SettingValue = item.SettingValue
			setting.ValueType = item.ValueType
			setting.OptionsJSON = item.OptionsJSON
			setting.IsEncrypted = item.IsEncrypted
			setting.Sort = item.Sort
			setting.Remark = item.Remark
			setting.UpdateTime = now

			if saveErr := tx.Save(&setting).Error; saveErr != nil {
				return saveErr
			}
		}
		return nil
	})
}
