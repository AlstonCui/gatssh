package sshClient

import (
	"fmt"
	"gatssh/models"
	"golang.org/x/crypto/ssh"
	"time"
	"gatssh/utils"
	"sync"
)

var mutex sync.Mutex

func (ct *CreateTask) StartNewTask() (err error) {

	go StartTaskWorkPool(ct)

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

func StartTaskWorkPool(ct *CreateTask) {

	ResultCatch.Store(ct.TaskId, ct.ResultChan)

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
		ct.TaskChan <- t
	}

	var wg sync.WaitGroup

	for i := 0; i < ct.PoolSize; i++ {
		go taskWorker(ct.TaskChan, ct.ResultChan, &wg)
	}

	wg.Wait()

	close(ct.TaskChan)
	close(ct.ResultChan)

	//ResultCatch.Delete(ct.TaskId)
}

func taskWorker(taskChan chan *Task, resultChan chan *models.TaskDetail, wg *sync.WaitGroup) {

	for t := range taskChan {

		wg.Add(1)

		var client *ssh.Client

		client, t.SshError = newGatSshClient(t)

		if t.SshError.Content != nil {
			t.createTaskDetail(resultChan, wg)
			continue
		}

		t.Standard, t.SshError = sshExecution(client, t.Cmd)
		if t.SshError.Code != 0 {
			t.createTaskDetail(resultChan, wg)
			continue
		}

		t.createTaskDetail(resultChan, wg)
	}
}

func (t *Task) createTaskDetail(resultChan chan *models.TaskDetail, wg *sync.WaitGroup) {

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

	resultChan <- td

	mutex.Lock()
	err := td.SaveTaskDetail()
	mutex.Unlock()
	if err != nil {
		utils.GatLog.Warning("Save task details:", err)
	}

	wg.Done()

	return
}
