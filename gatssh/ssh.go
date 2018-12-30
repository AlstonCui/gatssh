package gatssh

import "bytes"

const (
	SshNetworkError        = 1
	SshAuthenticationError = 2
	SshUnknowError         = 3
	SshCommandError        = 4
	NoMatchPassInDB        = 5
	SaveHostAndPassErr     = 6
)

type Host struct {
	User string `json:"user"`
	Addr string `json:"address"`
	Port int    `json:"port"`
}

type Task struct {
	TaskId   string
	Host     Host
	Auth     Auth
	GatUser  string
	Cmd      string
	PoolSize int
	Standard Standard
	SshError SshError
}

type Auth struct {
	Passwords       []string `json:"password"`
	SavePassword    bool     `json:"savePassword"`    //1 = yes, -1 = no
	UsePasswordInDB bool     `json:"usePasswordInDB"` //1 = yes, -1 = no
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
	Hosts    []Host `json:"hosts"`
	Auth     Auth   `json:"auth"`
	Cmd      string `json:"cmd"`
	PoolSize int    `json:"poolSize"`
	TaskId   string
	GatUser	string
}
