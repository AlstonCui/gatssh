package gatssh

import (
	"fmt"
	"gatlin/models"
	"log"
	"golang.org/x/crypto/ssh"
	"time"
)

func (ct *CreateTask) StartNewTask() (err error) {

	StartTaskWorkPool(ct)

	tr := models.TaskRecord{
		TaskId:           ct.TaskId,
		GatUser:          ct.GatUser,
		OperationContent: ct.Cmd,
		OperationTime:    time.Now(),
		HostCount:        len(ct.Hosts),
	}
	err = tr.TaskInsert()

	return
}

func StartTaskWorkPool(ct *CreateTask) {

	taskChan := make(chan *Task, 10000)

	for _, h := range ct.Hosts {

		t := &Task{
			TaskId:  ct.TaskId,
			Host:    h,
			Auth:    ct.Auth,
			Cmd:     ct.Cmd,
			GatUser: ct.GatUser,
		}
		taskChan <- t
	}

	for i := 0; i < ct.PoolSize; i++ {
		go taskWorker(taskChan)
	}

}

func taskWorker(taskChan chan *Task) {

	for t := range taskChan {

		var client *ssh.Client

		client, t.SshError = newGatSshClient(t)

		if t.SshError.Content != nil {

			err := t.createTaskDetail()
			if err != nil {
				log.Fatal(err)
			}
			continue
		}

		t.Standard, t.SshError = sshExecution(client, t.Cmd)
		if t.SshError.Code != 0 {
			err := t.createTaskDetail()
			if err != nil {
				log.Fatal(err)
				return
			}
			continue
		}
		err := t.createTaskDetail()
		if err != nil {
			log.Fatal(err)
		}
	}

}

func (t *Task) createTaskDetail() (err error) {

	td := &models.TaskDetail{
		TaskId:           t.TaskId,
		Ip:               t.Host.Addr,
		User:             t.GatUser,
		OperationContent: t.Cmd,
		OperationTime:    time.Now(),
		ResultCode:       t.SshError.Code,
		ResultErr:        fmt.Sprintln(t.SshError.Content),
		ResultContent:    "\n" + t.Standard.StdOut.String() + t.Standard.StdErr.String(),
	}
	fmt.Println(td)
	err = td.TaskDetailInsert()
	if err != nil {
		log.Fatal(err)
	}
	return
}

