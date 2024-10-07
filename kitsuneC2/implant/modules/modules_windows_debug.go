//go:build windows && debug

package modules

import (
	"errors"
	"log"
	"os/exec"
)

// TODO
func ShellcodeExec(shellcode []byte) error {
	log.Printf("[START WINDOWS SHELLCODE EXEC] starting new thread with following shellcode: % X\n", shellcode)
	log.Printf("[ERROR WINDOWS SHELLCODE EXEC] not yet implemented\n")
	return errors.New("Not yet implemented for this platform")
}

// Executes a command in Powershell and returns stdout
func Exec(cmd string) ([]byte, error) {
	log.Printf("[START WINDOWS EXEC] Powershell command: %s\n", cmd)
	command := exec.Command("C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe", cmd)
	byteOut, err := command.CombinedOutput()
	if err != nil {
		log.Printf("[ERROR WINDOWS EXEC] error: %s\n", err.Error())
		return nil, err
	}
	log.Printf("[SUCCESS WINDOWS EXEC] result: %s\n", string(byteOut))
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
