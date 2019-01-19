package models

import (
	"time"
	"errors"
)

type TaskRecord struct {
	ID               int    `gorm:"primary_key"`
	TaskId           string `gorm:"type:varchar(100);not null;"`
	GatUser          string `gorm:"type:varchar(100);not null;"`
	OperationContent string `gorm:"type:TEXT;not null;"`
	OperationTime    time.Time
	HostCount        int
}

func (this *TaskRecord) SaveTask() (err error) {
	if err := db.Create(this).Error; err != nil {
		return err
	}
	return
}


func (this *TaskRecord) QueryTask(taskId string) ( taskRecord TaskRecord,err error) {

	rows := db.Where("task_id = ? ", taskId).Find(&taskRecord)

	if rows.RowsAffected == 0 {
		err = errors.New("No match task record in DB,Please make sure the task_id is right...")
		return
	}
	return
}