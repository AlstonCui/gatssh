package ssh

import (
	"golang.org/x/crypto/ssh"
	"time"
	"net"
	"fmt"
)

func sshClient(h Host, password string) (client *ssh.Client, sshErr SshError) {

	config := &ssh.ClientConfig{

		User:    h.User,
		Timeout: 5 * time.Second,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", h.Addr, h.Port), config.Timeout)
	if err != nil {
		sshErr.Code = 1
		sshErr.Content = err
		return
	}

	c, chans, reqs, err := ssh.NewClientConn(conn, fmt.Sprintf("%s:%d", h.Addr, h.Port), config)
	if err != nil {
		sshErr.Code = 2
		sshErr.Content = err
		return
	}

	client = ssh.NewClient(c, chans, reqs)

	return
}

func tryPassRetClient(h Host, pws []string) (client *ssh.Client, sshErr SshError) {

	for _, pw := range pws {
		client, sshErr = sshClient(h, pw)
		if sshErr.Code == 2 {
			fmt.Println("unauthed")
			continue
		}
	}

	return
}
