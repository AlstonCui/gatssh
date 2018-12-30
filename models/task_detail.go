package models

import (
	"time"
)

type TaskDetail struct {
	ID             int    `gorm:"primary_key"`
	TaskId         string `gorm:"type:varchar(100);not null;"`
	Ip             string `gorm:"type:varchar(100);not null;"`
	User           string `gorm:"type:varchar(100);not null;"`
	OperationContent string `gorm:"type:TEXT;not null;"`
	OperationTime    time.Time
	ResultCode     int    `gorm:"type:int;not null;"`
	ResultContent  string `gorm:"type:TEXT;not null;"`
	ResultErr      string
}

func (this *TaskDetail) TaskDetailInsert() error {

	if err := db.Create(this).Error; err != nil {
		return err
	}
	return nil
}

/*func NewRecord(job *gatssh.Job) (record *Record) {
	return &Record{
		Ip:             job.Host.Addr,
		User:           job.GatUser,
		OperateContent: job.Cmd,
		OperateTime:    time.Now(),
		ResultCode:     job.SshError.Code,
		ResultErr:      fmt.Sprintln(job.SshError.Content),
		ResultContent:  "\n" + job.Standard.StdOut.String() + job.Standard.StdErr.String(),
	}
}*/
