package controllers

import (
	"gatlin/gatssh"
	"encoding/json"
	"github.com/satori/go.uuid"
)

type GatSshOneShoot struct {
	baseController
}

func (this GatSshOneShoot) Post() {

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

	ct.GatUser = "zhangsan"

	err = ct.StartNewTask()
	if err != nil{
		this.ServeJSON(40000,err)
	}

	this.ServeJSON(20000,ct.TaskId)
	return

}
