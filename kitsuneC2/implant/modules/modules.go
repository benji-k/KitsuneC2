//This package contains all functions that an implant is able to perform. These functions usually get called through handlers.
//Only pure-go, platform independent modules should go in this file.

package modules

import (
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

// Given a filepath, returns info about a file such as name, size etc.
func FileInfo(path string) (fs.FileInfo, error) {
	results, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// Executes a command and returns stdout
func Exec(cmd string, args []string) ([]byte, error) {
	command := exec.Command(cmd, args...)
	byteOut, err := command.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return byteOut, nil
}

// Writes file to specified location with name=filename.
func WriteFile(file []byte, location string, filename string) error {
	path, err := filepath.Abs(path.Join(location, filename))
	if err != nil {
		return err
	}
	fp, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = fp.Write(file)
	if err != nil {
		return err
	}
	return nil
}

// Reas
func ReadFile(path string) ([]byte, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return content, nil
}
