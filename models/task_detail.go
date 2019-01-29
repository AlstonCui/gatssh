package models

import (
	"time"
	"errors"
)

type TaskDetail struct {
	Id               int    `gorm:"primary_key"`
	//IdInTask         int
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


func QueryTaskDetails(taskId string, startID int) (taskDetails []TaskDetail, currentCount int, lastID int, err error) {

	if startID == 0 {
		rows := db.Where("task_id = ? ", taskId).Find(&taskDetails)

		if rows.RowsAffected == 0 {
			err = errors.New("The task is not finished yet, Or no match task details in DB,Please make sure the task_id is right...")
			return
		}
		currentCount = len(taskDetails)

		lastID = taskDetails[len(taskDetails)-1].Id

		return
	}

	rows := db.Where("task_id = ? AND id > ?", taskId,startID).Find(&taskDetails)
	if rows.RowsAffected == 0 {
		err = errors.New("The task is not finished yet, Or no match task details in DB,Please make sure the task_id is right...")
		return
	}

	currentCount = len(taskDetails)

	lastID = taskDetails[len(taskDetails)-1].Id

	return
}
