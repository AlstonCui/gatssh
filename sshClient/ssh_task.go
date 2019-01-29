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

//start
func (ct *CreateTask) StartNewTask() (err error) {
	//Save task record
	tr := models.TaskRecord{
		TaskId:           ct.TaskId,
		GatUser:          ct.GatUser,
		OperationContent: ct.Command,
		OperationTime:    time.Now(),
		HostCount:        len(ct.HostList),
	}
	err = tr.SaveTask()
	if err != nil{
		return
	}

	go StartTaskWorkPool(ct)
	return
}

//Create a work pool
func StartTaskWorkPool(ct *CreateTask) {
	//Put the result queue into the global cache
	//ResultCatch.Store(ct.TaskId, ct.ResultChan)
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
		go taskWorker(ct.TaskChan, &wg)
	}

	wg.Wait()

	close(ct.TaskChan)
}

func taskWorker(taskChan chan *Task, wg *sync.WaitGroup) {
	//Blocking and listening to the task queue
	for t := range taskChan {

		wg.Add(1)

		var client *ssh.Client

		client, t.SshError = newGatSshClient(t)

		if t.SshError.Content != nil {
			t.createTaskDetail()
			wg.Done()
			continue
		}

		t.Standard, t.SshError = sshExecution(client, t.Cmd)
		if t.SshError.Code != 0 {
			t.createTaskDetail()
			wg.Done()
			continue
		}

		t.createTaskDetail()
		wg.Done()
	}
}

func (t *Task) createTaskDetail() {

	td := &models.TaskDetail{
		TaskId:           t.TaskId,
		Ip:               t.Host.Addr,
		Port:             t.Host.Port,
		GatUser:          t.GatUser,
		OperationContent: t.Cmd,
		OperationTime:    time.Now(),
		ResultCode:       t.SshError.Code,
		ResultErr:        fmt.Sprintln(t.SshError.Content),
		ResultContent:    t.Standard.StdOut.String() + t.Standard.StdErr.String(),
	}
	//Limit the rate of writing to DB,Because Sqlite does not support high concurrent writes.
	mutex.Lock()
	err := td.SaveTaskDetail()
	mutex.Unlock()
	if err != nil {
		utils.GatLog.Warning("Save task details:", err)
	}

	return
}
