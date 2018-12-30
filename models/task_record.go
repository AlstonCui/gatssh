package models

import (
	"time"
)

type TaskRecord struct {
	ID               int    `gorm:"primary_key"`
	TaskId           string `gorm:"type:varchar(100);not null;"`
	GatUser          string `gorm:"type:varchar(100);not null;"`
	OperationContent string `gorm:"type:TEXT;not null;"`
	OperationTime    time.Time
	HostCount        int
}

func (this *TaskRecord) TaskInsert() (err error) {
	if err := db.Create(this).Error; err != nil {
		return err
	}
	return
}
