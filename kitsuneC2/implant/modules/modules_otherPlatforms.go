//go:build !linux && !windows && !debug

package modules

import (
	"errors"
)

func ShellcodeExec(shellcode []byte) error {
	return errors.New("Not yet implemented for this platform")
}

func Exec(cmd string) ([]byte, error) {
	return errors.New("Not yet implemented for this platform")
}
