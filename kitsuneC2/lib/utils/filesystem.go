package utils

import (
	"os"
	"path/filepath"
)

// Writes file to specified location.
func WriteFile(file []byte, location string) error {
	path, err := filepath.Abs(location)
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
