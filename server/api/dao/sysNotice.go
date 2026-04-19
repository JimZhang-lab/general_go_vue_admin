package dao

import (
	"server/api/entity"
	"server/common/utils"
	"server/pkg/db"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetSysNoticeList(title, noticeType, status string, pageSize, pageNum int) (notices []entity.SysNotice, count int64) {
	notices = make([]entity.SysNotice, 0)
	curDb := db.Db.Model(&entity.SysNotice{})
	if title != "" {
		curDb = curDb.Where("title LIKE ?", "%"+title+"%")
	}
	if noticeType != "" {
		curDb = curDb.Where("notice_type = ?", noticeType)
	}
	if status != "" {
		curDb = curDb.Where("status = ?", status)
	}
	curDb.Count(&count)
	curDb.Order("status ASC, publish_time DESC, create_time DESC").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Find(&notices)
	return notices, count
}

func GetSysNoticeById(id uint) (notice entity.SysNotice) {
	db.Db.First(&notice, id)
	return notice
}

func SaveSysNotice(notice *entity.SysNotice) error {
	return db.Db.Save(notice).Error
}

func DeleteSysNoticeById(id uint) error {
	return db.Db.Delete(&entity.SysNotice{}, id).Error
}

func BatchDeleteSysNotice(ids []uint) error {
	return db.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("notice_id IN (?)", ids).Delete(&entity.SysNoticeRead{}).Error; err != nil {
			return err
		}
		return tx.Where("id IN (?)", ids).Delete(&entity.SysNotice{}).Error
	})
}

func GetCurrentSysNoticeList(adminID uint, limit int, unreadOnly bool) (list []entity.CurrentSysNoticeVo) {
	list = make([]entity.CurrentSysNoticeVo, 0)
	query := db.Db.Table("sys_notice AS n").
		Select("n.id, n.title, n.content, n.notice_type, n.notice_level, n.status, n.target_type, n.created_by, n.created_by_name, n.publish_time, n.create_time, CASE WHEN nr.id IS NULL THEN false ELSE true END AS is_read").
		Joins("LEFT JOIN sys_notice_read AS nr ON nr.notice_id = n.id AND nr.admin_id = ?", adminID).
		Where("n.status = ?", 2).
		Where("n.target_type IN ?", []string{"all", "admin"})

	if unreadOnly {
		query = query.Where("nr.id IS NULL")
	}

	query.Order("n.publish_time DESC, n.create_time DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	query.Find(&list)
	return list
}

func MarkSysNoticeRead(adminID uint, noticeID uint) error {
	var noticeRead entity.SysNoticeRead
	err := db.Db.Where("admin_id = ? AND notice_id = ?", adminID, noticeID).First(&noticeRead).Error
	if err == nil {
		noticeRead.ReadTime = utils.HTime{Time: time.Now()}
		return db.Db.Save(&noticeRead).Error
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	noticeRead = entity.SysNoticeRead{
		AdminId:  adminID,
		NoticeId: noticeID,
		ReadTime: utils.HTime{Time: time.Now()},
	}
	return db.Db.Create(&noticeRead).Error
}

func MarkAllSysNoticeRead(adminID uint) error {
	notices := GetCurrentSysNoticeList(adminID, 0, true)
	if len(notices) == 0 {
		return nil
	}

	now := utils.HTime{Time: time.Now()}
	reads := make([]entity.SysNoticeRead, 0, len(notices))
	for _, notice := range notices {
		reads = append(reads, entity.SysNoticeRead{
			AdminId:  adminID,
			NoticeId: notice.ID,
			ReadTime: now,
		})
	}

	return db.Db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "notice_id"}, {Name: "admin_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"read_time"}),
	}).Create(&reads).Error
}
