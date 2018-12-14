package ssh

import (
"golang.org/x/crypto/ssh"
"time"
"net"
"fmt"
)

func sshClient(h Host) (client *ssh.Client, sshErr SshError) {

	config := &ssh.ClientConfig{

		User:    h.User,
		Timeout: 5 * time.Second,
		Auth: []ssh.AuthMethod{
			ssh.Password(h.Passwd),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", h.Addr, h.Port), config.Timeout)
	if err != nil {
		sshErr.Code = 1
		sshErr.Content = err
		//checkErr(sshErr)
		return
	}
	//defer conn.Close()

	c, chans, reqs, err := ssh.NewClientConn(conn, fmt.Sprintf("%s:%d", h.Addr, h.Port), config)
	if err != nil {
		sshErr.Code = 2
		sshErr.Content = err
		//checkErr(sshErr)
		return
	}
	//defer c.Close()

	client = ssh.NewClient(c, chans, reqs)

	return
}
