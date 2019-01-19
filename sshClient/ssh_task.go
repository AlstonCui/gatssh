package sshClient

import (
	"fmt"
	"gatssh/models"
	"golang.org/x/crypto/ssh"
	"time"
	"gatssh/utils"
)

func (ct *CreateTask) StartNewTask(taskChan chan *Task, resultChan chan *models.TaskDetail) (err error) {

	StartTaskWorkPool(ct, taskChan,resultChan)

	tr := models.TaskRecord{
		TaskId:           ct.TaskId,
		GatUser:          ct.GatUser,
		OperationContent: ct.Command,
		OperationTime:    time.Now(),
		HostCount:        len(ct.HostList),
	}
	err = tr.SaveTask()

	return
}

func StartTaskWorkPool(ct *CreateTask, taskChan chan *Task,resultChan chan *models.TaskDetail) {

	for _, h := range ct.HostList {

		t := &Task{
			TaskId:          ct.TaskId,
			Host:            h,
			Auth:            ct.AuthList,
			Cmd:             ct.Command,
			GatUser:         ct.GatUser,
			SavePassword:    false,
			UsePasswordInDB: false,
		}
		taskChan <- t
	}

	for i := 0; i < ct.PoolSize; i++ {
		go taskWorker(taskChan, resultChan)
	}
}

func taskWorker(taskChan chan *Task,resultChan chan *models.TaskDetail) {

	for t := range taskChan {
		var client *ssh.Client

		client, t.SshError = newGatSshClient(t)

		if t.SshError.Content != nil {
			t.createTaskDetail(resultChan)
			continue
		}

		t.Standard, t.SshError = sshExecution(client, t.Cmd)
		if t.SshError.Code != 0 {
			t.createTaskDetail(resultChan)
			continue
		}

		t.createTaskDetail(resultChan)
	}

}

func (t *Task) createTaskDetail(resultChan chan *models.TaskDetail) {

	td := &models.TaskDetail{
		TaskId:           t.TaskId,
		Ip:               t.Host.Addr,
		GatUser:          t.GatUser,
		OperationContent: t.Cmd,
		OperationTime:    time.Now(),
		ResultCode:       t.SshError.Code,
		ResultErr:        fmt.Sprintln(t.SshError.Content),
		ResultContent:    t.Standard.StdOut.String() + t.Standard.StdErr.String(),
	}

	err := td.SaveTaskDetail()
	if err != nil {
		utils.GatLog.Warning("Save task details:",err)
		return
	}

	resultChan <- td

	return
}
