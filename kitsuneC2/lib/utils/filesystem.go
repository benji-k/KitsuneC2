package utils

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

// Writes file to specified location.
func WriteFile(file []byte, location string) error {
	path, err := filepath.Abs(location)
	if err != nil {
		return err
	}

	//we don't want this function to override existing files for 2 reasons. 1st being that one might accidentally corrupt files on the remote host
	//where the implant is running. If the operator want's to replace a file, they should delete it first. 2nd reason: We don't want the server
	//to overwrite existing files on it's OS. A malicious operator could, in theory, write implant binaries to arbitrary locations on the system.
	//we don't want a malicous operator to potentially corrupt the system.
	_, err = os.Stat(path)
	if err == nil {
		return errors.New("cannot write file because it already exists")
	}

	err = os.MkdirAll(filepath.Dir(path), 0700)
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

// Reads file from "path" and returns content as bytes
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

// Recursively copies contents of "source folder" into "destination"
func CopyFolder(source string, destination string) error {
	log.Printf("[INFO] Attempting to copy folder %s into %s...", source, destination)
	dirEntries, err := os.ReadDir(source)
	if err != nil {
		return err
	}
	for _, entry := range dirEntries {
		if entry.IsDir() {
			src := filepath.Join(source, entry.Name())
			dst := filepath.Join(destination, entry.Name())
			entryInfo, _ := entry.Info()
			err = os.Mkdir(dst, entryInfo.Mode())
			if err != nil {
				log.Printf("[ERROR] Could not copy implant source files to temp folder. Reason %s", err)
				return err
			}
			err = CopyFolder(src, dst)
			if err != nil {
				log.Printf("[ERROR] Could not copy implant source files to temp folder. Reason %s", err)
				return err
			}
		} else {
			src := filepath.Join(source, entry.Name())
			dst := filepath.Join(destination, entry.Name())
			err := os.Link(src, dst)
			if err != nil {
				log.Printf("[ERROR ]Could not copy implant source files to temp folder. Reason %s", err)
				return err
			}
		}
	}
	return nil
}
