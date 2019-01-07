package controllers

import (
	"gatlin/gatssh"
	"encoding/json"
	"github.com/satori/go.uuid"
)

type GatSshMultiShoot struct {
	baseController
}

func (this *GatSshMultiShoot) Post() {
	if this.IsLogin != true {
		this.ServeJSON(40000,"Please login...")
		return
	}
	var ct *gatssh.CreateTask

	err := json.NewDecoder(this.Ctx.Request.Body).Decode(&ct)
	if err != nil {
		this.ServeJSON(40000, err)
		return
	}
	if ct.Cmd ==""{
		this.ServeJSON(40000, "No command received...")
		return
	}
	ct.TaskId = uuid.NewV4().String()

	ct.GatUser = this.User

	err = ct.StartNewTask()
	if err != nil{
		this.ServeJSON(40000,err)
	}

	this.ServeJSON(20000,ct.TaskId)
	return


}
