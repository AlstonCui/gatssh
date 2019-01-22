package sshClient

import (
	"golang.org/x/crypto/ssh"
	"time"
	"net"
	"fmt"
	"gatssh/models"
)

func sshClient(h Host,user string,password string) (client *ssh.Client, sshErr SshError) {

	//signer, err := ssh.ParsePrivateKey(utils.Key)
	//SSH auth config
	config := &ssh.ClientConfig{

		User:    user,
		Timeout: 5 * time.Second,
		Auth: []ssh.AuthMethod{
	//		ssh.PublicKeys(signer),
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
	//Establishing TCP connection，if err，it's a network problem.
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", h.Addr, h.Port), config.Timeout)
	if err != nil {
		sshErr.Code = SshNetworkError
		sshErr.Content = err
		return
	}
	//Establishing SSH connection，if err，it's authentication problem.
	c, chans, reqs, err := ssh.NewClientConn(conn, fmt.Sprintf("%s:%d", h.Addr, h.Port), config)
	if err != nil {
		sshErr.Code = SshAuthenticationError
		sshErr.Content = err
		return
	}
	//Instantiate the ssh client
	client = ssh.NewClient(c, chans, reqs)

	return
}

func newGatSshClient(t *Task) (client *ssh.Client, sshErr SshError) {
	//Whether to use the password in the database
	if t.UsePasswordInDB {
		var pass string
		user,pass, err := models.QueryHost(t.Host.Addr, t.Host.Port, t.GatUser)

		if err != nil {
			sshErr.Code = NoMatchPassInDB
			sshErr.Content = err
			return
		}

		client, sshErr = sshClient(t.Host,user, pass)
		if sshErr.Code != 0 {
			return
		}
		return
	}
	//Whether to save the password into the database
	if t.SavePassword == true {
		for _, au := range t.Auth {

			client, sshErr = sshClient(t.Host,au.User,au.Password)

			if sshErr.Code == SshAuthenticationError {
				continue
			}
			if sshErr.Content != nil {
				return
			}

			h := &models.Host{
				Ip:       t.Host.Addr,
				Port:     t.Host.Port,
				User:     au.User,
				Owner:    t.GatUser,
				Password: au.Password,
			}

			err := h.SaveHost()
			if err != nil {
				sshErr.Code = SaveHostAndPassErr
				sshErr.Content = err
				return
			}
			return
		}
		return
	}
	//If more than one password is inputted, It's can try them one at a time
	for _, au := range t.Auth{
		client, sshErr = sshClient(t.Host,au.User,au.Password)
		if sshErr.Code == SshAuthenticationError {
			continue
		}
		if sshErr.Content != nil {
			return
		}
		return
	}
	return
}
//Command execution
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
