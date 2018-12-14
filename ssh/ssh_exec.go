package ssh

import (
	"sync"
	"fmt"
	"strings"
)

func Execution(h Host, pws []string, cmd string) (std Std, sshErr SshError) {

	client, sshErr := tryPassRetClient(h, pws)
	if sshErr.Content != nil{
		return
	}

	session, err := client.NewSession()
	if err != nil {
		sshErr.Code = 3
		sshErr.Content = err
		return
	}


	defer session.Close()
	defer client.Close()

	session.Stdout = &std.StdOut
	session.Stderr = &std.StdErr

	err = session.Run(cmd)
	if err != nil {
		sshErr.Code = 4
		sshErr.Content = err
		return
	}


	return
}

func ConcurrentExecutionProcess(cs *CmdSession) (cmdRes []*CmdResult) {

	var wg sync.WaitGroup

	for _, h := range cs.Hosts {

		wg.Add(1)

		go func(h Host) {

			std, sshErr := Execution(h, cs.Passwords, cs.Cmd)

			cr := NewCmdResult(
				strings.Replace(std.StdOut.String(), "\n", "", -1),
				strings.Replace(std.StdErr.String(), "\n", "", -1),
				sshErr, h.Addr)

			cmdRes = append(cmdRes, cr)
			fmt.Println(cr)
			wg.Done()
		}(h)

	}

	wg.Wait()
	return

}

