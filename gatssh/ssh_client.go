package gatssh

import (
	"golang.org/x/crypto/ssh"
	"time"
	"net"
	"fmt"
	"gatlin/models"
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
		Config: ssh.Config{
			Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr",
				"aes128-gcm@openssh.com", "arcfour256", "arcfour128",
				"aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
		},
	}

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", h.Addr, h.Port), config.Timeout)
	if err != nil {
		sshErr.Code = SshNetworkError
		sshErr.Content = err
		return
	}

	c, chans, reqs, err := ssh.NewClientConn(conn, fmt.Sprintf("%s:%d", h.Addr, h.Port), config)
	if err != nil {
		sshErr.Code = SshAuthenticationError
		sshErr.Content = err
		return
	}

	client = ssh.NewClient(c, chans, reqs)

	return
}

func newGatSshClient(t *Task) (client *ssh.Client, sshErr SshError) {

	if t.Auth.UsePasswordInDB {
		var pass string
		pass, err := models.HostSearch(t.Host.Addr, t.Host.Port, t.GatUser)

		if err != nil {
			sshErr.Code = NoMatchPassInDB
			sshErr.Content = err
			return
		}

		client, sshErr = sshClient(t.Host, pass)
		if sshErr.Code != 0 {
			return
		}
		return
	}

	if t.Auth.SavePassword {
		for _, pw := range t.Auth.Passwords {
			client, sshErr = sshClient(t.Host, pw)

			if sshErr.Code == SshAuthenticationError {
				continue
			}
			if sshErr.Content != nil {
				return
			}

			h := &models.Host{
				Ip:       t.Host.Addr,
				Port:     t.Host.Port,
				User:     t.Host.User,
				Owner:    t.GatUser,
				Password: pw,
			}

			err := h.HostSave()
			if err != nil {
				sshErr.Code = SaveHostAndPassErr
				sshErr.Content = err
				return
			}
			return
		}
		return
	}

	for _, pw := range t.Auth.Passwords {
		client, sshErr = sshClient(t.Host, pw)
		if sshErr.Code == SshAuthenticationError {
			continue
		}
		if sshErr.Content != nil {
			return
		}
	}
	return
}

func sshExecution(client *ssh.Client, cmd string) (std Standard, sshErr SshError) {

	session, err := client.NewSession()
	if err != nil {
		sshErr.Code = SshUnknowError
		sshErr.Content = err
		return
	}

	defer session.Close()
	defer client.Close()

	session.Stdout = &std.StdOut
	session.Stderr = &std.StdErr

	err = session.Run(cmd)
	if err != nil {
		sshErr.Code = SshCommandError
		sshErr.Content = err
		return
	}

	return
}
