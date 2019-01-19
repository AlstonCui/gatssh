package controllers

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"github.com/gorilla/websocket"
	"gatssh/models"
	"gatssh/sshClient"
	"gatssh/utils"
	"errors"
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

	ct.TaskId = uuid.NewV4().String()
	ct.PoolSize = 1000
	ct.GatUser = this.User
	ct.UsePasswordInDB = false
	ct.SavePassword = false
	ct.TaskChan = make(chan *sshClient.Task, len(ct.HostList))
	ct.ResultChan = make(chan *models.TaskDetail, len(ct.HostList))

	err = ct.StartNewTask()
	if err != nil {
		this.ServeJSON(40000, err)
	}

	this.ServeJSON(20000, ct.TaskId)

	return
}

func (this *GatSshQuickStart) StartSendByWS() {

	if this.IsLogin != true {
		this.Delete()
		return
	}

	TaskId := this.GetString("taskId")

	var upGrader = websocket.Upgrader{}

	ws, err := upGrader.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil)
	if err != nil {
		utils.GatLog.Alert("WebSocket Upgrade: %v", err)
		this.ServeJSON(40000, err)
		return
	}

	defer ws.Close()

	resultChan, ok := sshClient.ResultCatch.Load(TaskId)
	if ok {

		for i := 1; i <= cap(resultChan.(chan *models.TaskDetail)); i++ {

			result := <-resultChan.(chan *models.TaskDetail)

			result.IdInTask = i

			err := ws.WriteJSON(result)
			if err != nil {
				utils.GatLog.Alert("WebSocket Write JSON error: %v", err)
				break
			}
		}
	} else {
		err = errors.New("This task id is not correct:")
		utils.GatLog.Alert("%v %v", err, TaskId)
		this.ServeJSON(40000, err)
		return
	}

	this.ServeJSON(20000, "Task results are all transmitted")

	return
}
