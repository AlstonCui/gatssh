package controllers

import (
	"gatssh/models"
)

type QueryGatSshTaskResults struct {
	baseController
}

type taskResults struct {
	TaskDetails  []models.TaskDetail `json:"taskDetails"`
	CurrentCount int                 `json:"currentCount"`
	LastID       int                 `json:"lastId"`
}

func (this *QueryGatSshTaskResults) Post() {
	if this.IsLogin != true {
		this.ServeJSON(40000, "Please login...")
		return
	}

	taskId := this.GetString("taskId")

	startID, err := this.GetInt("startId")
	if err != nil {
		startID = 0
	}

	taskDetails, currentCount, lastID, err := models.QueryTaskDetails(taskId, startID)
	if err != nil {
		this.ServeJSON(40000, err)
		return
	}

	tr := &taskResults{
		TaskDetails:  taskDetails,
		CurrentCount: currentCount,
		LastID:       lastID,
	}

	this.ServeJSON(20000, tr)
	return
}
