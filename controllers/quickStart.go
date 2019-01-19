package controllers

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"github.com/gorilla/websocket"
	"gatssh/models"
	"gatssh/gatssh"
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

	var ct *gatssh.CreateTask
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
	ct.PoolSize = 10
	ct.GatUser = this.User
	ct.UsePasswordInDB = false
	ct.SavePassword = false
	ct.TaskChan = make(chan *gatssh.Task, len(ct.HostList))
	ct.ResultChan = make(chan *models.TaskDetail, len(ct.HostList))

	TaskCatch.Store(ct.TaskId,ct.TaskChan)
	ResultCatch.Store(ct.TaskId, ct.ResultChan)


	err = ct.StartNewTask(ct.TaskChan,ct.ResultChan)
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
		this.ServeJSON(40000,err)
		return
	}

	defer ws.Close()

	resultChan,ok := ResultCatch.Load(TaskId)
	if ok {

		for i := 1; i <= cap(resultChan.(chan *models.TaskDetail)); i++ {

			result := <-resultChan.(chan *models.TaskDetail)

			result.Id = i

			err := ws.WriteJSON(result)
			if err != nil {
				utils.GatLog.Alert("WebSocket Write JSON error: %v", err)
				break
			}
		}
	}else {
		err = errors.New("This task id is not correct:")
		utils.GatLog.Alert("%v %v", err,TaskId)
		this.ServeJSON(40000,err)
		return
	}


	taskChan,_:= TaskCatch.Load(TaskId)

	close(taskChan.(chan *gatssh.Task))
	close(resultChan.(chan *models.TaskDetail))

	TaskCatch.Delete(TaskId)
	ResultCatch.Delete(TaskId)

	this.ServeJSON(20000, "Task results are all transmitted")

	return
}
