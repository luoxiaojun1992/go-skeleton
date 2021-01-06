package models

import (
	"github.com/luoxiaojun1992/go-skeleton/services/db"
	"gorm.io/gorm/clause"
	"time"
)

type PhabricatorBugTasks struct {
	BaseModel
	JingUUID         string    `gorm:"primaryKey,column:jing_uuid"`
	TaskID           int64     `gorm:"column:task_id"`
	TaskPHID         string    `gorm:"column:task_phid"`
	TaskName         string    `gorm:"column:task_name"`
	TaskStatus       string    `gorm:"column:task_status"`
	TaskPriority     float64   `gorm:"column:task_priority"`
	TaskSubPriority  float64   `gorm:"column:task_sub_priority"`
	TaskAuthorPHID   string    `gorm:"column:task_author_phid"`
	TaskAuthorName   string    `gorm:"column:task_author_name"`
	TaskOwnerPHID    string    `gorm:"column:task_owner_phid"`
	TaskOwnerName    string    `gorm:"column:task_owner_name"`
	TaskCloserPHID   string    `gorm:"column:task_closer_phid"`
	TaskCloserName   string    `gorm:"column:task_closer_name"`
	TaskClosed       int64     `gorm:"column:task_closed"`
	TaskClosedDate   time.Time `gorm:"column:task_closed_datetime"`
	TaskCreated      int64     `gorm:"column:task_created"`
	TaskCreatedDate  time.Time `gorm:"column:task_created_datetime"`
	TaskModified     int64     `gorm:"column:task_modified"`
	TaskModifiedDate time.Time `gorm:"column:task_modified_datetime"`
	TaskTags         string    `gorm:"column:task_tags"`
}

// ORM框架自动通过此方法获取表名，如果不存在此方法默认根据model struct名解析表名
func (pbt *PhabricatorBugTasks) TableName() string {
	return "phabricator_bug_tasks"
}

func (pbt *PhabricatorBugTasks) Connection() string {
	return ""
}

func (pbt *PhabricatorBugTasks) BatchInsert(tasks []PhabricatorBugTasks, retry bool) error {
	doBatchInsert := func() error {
		return pbt.Query(pbt).DBClient.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "task_id"}, {Name: "task_status"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"task_phid", "task_name", "task_priority", "task_sub_priority", "task_author_phid",
				"task_author_name", "task_owner_phid", "task_owner_name", "task_closer_phid",
				"task_closer_name", "task_closed", "task_closed_datetime", "task_created",
				"task_created_datetime", "task_modified", "task_modified_datetime", "task_tags",
			}),
		}).Create(&tasks).Error
	}

	err := doBatchInsert()

	if retry {
		if db.CausedByLostConnection(err) {
			return doBatchInsert()
		}
	}

	return err
}
