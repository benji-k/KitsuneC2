//go:build !linux && !windows && debug

package modules

import (
	"errors"
	"log"
)

func ShellcodeExec(shellcode []byte) error {
	log.Printf("[START OTHER-PLATFORM SHELLCODE EXEC] starting new thread with following shellcode: % X\n", shellcode)
	log.Printf("[ERROR OTHER-PLATFORM SHELLCODE EXEC] not yet implemented\n")
}

func Exec(cmd string) ([]byte, error) {
	log.Printf("[START OTHER-PLATFORM EXEC] cmd: %s\n", cmd)
	log.Printf("[ERROR OTHER-PLATFORM EXEC] not yet implemented\n")

	return errors.New("Not yet implemented for this platform")
}
