package sshClient

import (
	"bytes"
	"gatssh/models"
)

const (
	SshNetworkError        = 1
	SshAuthenticationError = 2
	SshUnknowError         = 3
	SshCommandError        = 4
	NoMatchPassInDB        = 5
	SaveHostAndPassErr     = 6
)


type Host struct {
	Addr string `json:"addr"`
	Port int    `json:"port"`
}

type Task struct {
	TaskId          string
	Host            Host
	Auth            []Auth
	GatUser         string
	Cmd             string
	Standard        Standard
	SshError        SshError
	SavePassword    bool `json:"savePassword"`
	UsePasswordInDB bool `json:"usePasswordInDB"`
}

type Auth struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type Standard struct {
	StdOut bytes.Buffer
	StdErr bytes.Buffer
}

type SshError struct {
	Code    int
	Content error
}

type Result struct {
	StdOut string
	StdErr string
	SshErr interface{}
	Ip     string
}

type CreateTask struct {
	HostList        []Host `json:"hostList"`
	AuthList        []Auth `json:"authList"`
	Command         string `json:"command"`
	TaskId          string
	GatUser         string
	PoolSize        int
	SavePassword    bool   `json:"savePassword"`
	UsePasswordInDB bool   `json:"usePasswordInDB"`
	TaskChan        chan *Task
	ResultChan      chan *models.TaskDetail
}

