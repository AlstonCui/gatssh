package controllers

import (
	"gatlin/ssh"
	"encoding/json"
)

type RemoteCmd struct {
	baseController
}

func (this RemoteCmd) Post() {

	var cs *ssh.CmdSession

	err := json.NewDecoder(this.Ctx.Request.Body).Decode(&cs)
	if err != nil {
		this.ServeJSON(40000, err)
		return
	}

	res := ssh.ConcurrentExecutionProcess(cs)

	this.ServeJSON(20000, res)
	return

}
