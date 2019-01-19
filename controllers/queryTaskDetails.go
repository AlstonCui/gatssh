package controllers

import (
	"gatssh/models"
)

type QueryGatSshTaskDetails struct {
	baseController
}

func (this *QueryGatSshTaskDetails) Post() {
	if this.IsLogin != true {
		this.ServeJSON(40000, "Please login...")
		return
	}

	taskId:= this.GetString("taskId")

	var td *models.TaskDetail

	tds ,err := td.QueryTaskDetails(taskId)
	if err != nil {
		this.ServeJSON(40000,err)
		return
	}

	this.ServeJSON(20000,tds)
	return
}




