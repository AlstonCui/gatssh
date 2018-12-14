package ssh

import (
	"sync"
	"fmt"
	"strings"
)

func Execution(h Host, cmd string) (std Std, sshErr SshError) {

	client, sshErr := sshClient(h)
	if sshErr.Content != nil {
		//checkErr(sshErr)
		return
	}

	session, err := client.NewSession()
	if err != nil {
		sshErr.Code = 3
		sshErr.Content = err
		//checkErr(sshErr)
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
		//checkErr(sshErr)
		return
	}

	return
}

func ConcurrentExecutionProcess(cmdse CmdSession) (cmdRes []*CmdResult) {

	var wg sync.WaitGroup
	exec := func(h Host) {
		fmt.Println("in", h)
		std, sshErr := Execution(h, cmdse.Cmd)


		cr := NewCmdResult(h.Addr,
			strings.Replace(std.StdOut.String(), "\n", "", -1),
			strings.Replace(std.StdErr.String(), "\n", "", -1),sshErr)

		//fmt.Println(cr)

		cmdRes = append(cmdRes, cr)
		fmt.Println(cr)
		wg.Done()
	}

	for _, host := range cmdse.Hosts {
		wg.Add(1)
		go exec(host)

	}

	wg.Wait()
	return

}
