//go:build !linux && !windows

package modules

import (
	"errors"
)

func ShellcodeExec(shellcode []byte) error {
	return errors.New("Not yet implemented for this platform")
}
