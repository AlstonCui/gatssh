package ssh

import (
	"bytes"
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
	Hosts     []Host   `json:"hosts"`
	Cmd       string   `json:"cmd"`
	Passwords []string `json:"password"`
}

type SshError struct {
	/*  code:
		SshNetworkError         //1 = 网络错误
		SshAuthenticationError  //2 = 验证错误
		SshUnknowError          //3 = 未知错误
		SshCommandError         //4 = 命令错误*/
	Code    int
	Content error
}

type CmdResult struct {
	StdOut string
	StdErr string
	SshErr interface{}
	Ip     string
}

func NewCmdResult(StdOut string, StdErr string, SshErr interface{}, Ip string) *CmdResult {
	return &CmdResult{
		Ip:     Ip,
		StdOut: StdOut,
		StdErr: StdErr,
		SshErr: SshErr,
	}
}
