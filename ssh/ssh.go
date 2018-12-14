package ssh

import (
	"bytes"
	"log"
)

type Host struct {
	User string `json:"user"`
	Addr string `json:"address"`
	Port int    `json:"port"`
}
type Std struct {
	StdOut bytes.Buffer
	StdErr bytes.Buffer
}

type CmdSession struct {
	Hosts    []Host   `json:"hosts"`
	Cmd      string   `json:"cmd"`
	Password []string `json:"password"`
}

type SshError struct {
	// code:
	//1 = 网络错误
	//2 = 验证错误
	//3 = 未知错误
	//4 = 命令错误
	Code    int
	Content error
}

type CmdResult struct {
	Ip     string
	StdOut string
	StdErr string
	SshErr interface{}
}

func NewCmdResult(Ip string, StdOut string, StdErr string, SshErr interface{}) *CmdResult {
	return &CmdResult{
		Ip:     Ip,
		StdOut: StdOut,
		StdErr: StdErr,
		SshErr: SshErr,
	}
}

const (
	SshSucceed             = iota
	SshNetworkError         //1 = 网络错误
	SshAuthenticationError  //2 = 验证错误
	SshUnknowError          //3 = 未知错误
	SshCommandError         //4 = 命令错误
)

func checkErr(sshErr SshError) {
	if sshErr.Content != nil {
		log.Printf("code: %v, error: %v\n", sshErr.Code, sshErr.Content)
		return
	}
}
