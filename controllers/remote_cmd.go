package controllers

import (
	"gatlin/ssh"
	"encoding/json"
)

type RemoteCmd struct {
	baseController
}

func (this RemoteCmd)Post()  {
	var cmdSession ssh.CmdSession

	//json.Unmarshal(this.Ctx.Input.RequestBody,&cmdSession)


	json.NewDecoder(this.Ctx.Request.Body).Decode(&cmdSession)

	res:= ssh.ConcurrentExecutionProcess(cmdSession)




	this.ServeJSON(20000,res)
	return


}