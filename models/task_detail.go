package models

import (
	"time"
	"errors"
)

type TaskDetail struct {
	Id               int    `gorm:"primary_key"`
	IdInTask         int
	TaskId           string `gorm:"type:varchar(100);not null;"`
	Ip               string `gorm:"type:varchar(100);not null;"`
	Port             int
	GatUser          string `gorm:"type:varchar(100);not null;"`
	OperationContent string `gorm:"type:TEXT;not null;"`
	OperationTime    time.Time
	ResultCode       int    `gorm:"type:int;not null;"`
	ResultContent    string `gorm:"type:TEXT;not null;"`
	ResultErr        string
}

func (this *TaskDetail) SaveTaskDetail() error {

	if err := db.Create(this).Error; err != nil {
		return err
	}
	return nil
}

func (this *TaskDetail) QueryTaskDetails(taskId string) (taskDetails []TaskDetail, err error) {

	rows := db.Where("task_id = ? ", taskId).Find(&taskDetails)

	if rows.RowsAffected == 0 {
		err = errors.New("No match task details in DB,Please make sure the task_id is right...")
		return
	}
	return
}
