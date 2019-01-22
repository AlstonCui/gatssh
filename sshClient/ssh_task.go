package sshClient

import (
	"fmt"
	"gatssh/models"
	"golang.org/x/crypto/ssh"
	"time"
	"gatssh/utils"
	"sync"
)


//Global cache for task execution results, key = taskId, value = ResultChan
var ResultCatch = sync.Map{}

var mutex sync.Mutex

//start
func (ct *CreateTask) StartNewTask() (err error) {

	go StartTaskWorkPool(ct)
	//Save task record
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

//Create a work pool
func StartTaskWorkPool(ct *CreateTask) {
	//Put the result queue into the global cache
	ResultCatch.Store(ct.TaskId, ct.ResultChan)
	//Break task into subtasks，and put it on the task queue.
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
	//Create worker in work pool
	for i := 0; i < ct.PoolSize; i++ {
		go taskWorker(ct.TaskChan, ct.ResultChan, &wg)
	}

	wg.Wait()

	close(ct.TaskChan)
	close(ct.ResultChan)

	//ResultCatch.Delete(ct.TaskId)
}

func taskWorker(taskChan chan *Task, resultChan chan *models.TaskDetail, wg *sync.WaitGroup) {
	//Blocking and listening to the task queue
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
	//Put the result into the result queue
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
