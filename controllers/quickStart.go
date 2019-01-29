package controllers

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"gatssh/sshClient"
)

type GatSshQuickStart struct {
	baseController
}

func (this *GatSshQuickStart) QuickStart() {

	if this.IsLogin != true {
		this.ServeJSON(40000, "Please login...")
		return
	}

	var ct *sshClient.CreateTask
	err := json.NewDecoder(this.Ctx.Request.Body).Decode(&ct)
	if err != nil {
		this.ServeJSON(40000, err)
		return
	}

	if ct.Command == "" {
		this.ServeJSON(40000, "No command received...")
		return
	}
	//Create task
	ct.TaskId = uuid.NewV4().String()
	ct.PoolSize = 1000
	ct.GatUser = this.User
	//ct.UsePasswordInDB = false
	//ct.SavePassword = false
	ct.TaskChan = make(chan *sshClient.Task, len(ct.HostList))

	err = ct.StartNewTask()
	if err != nil {
		this.ServeJSON(40000, err)
	}

	this.ServeJSON(20000, ct.TaskId)

	return
}

