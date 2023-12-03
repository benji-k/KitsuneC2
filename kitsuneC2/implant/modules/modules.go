//This package contains all functions that an implant is able to perform. These functions usually get called through handlers.
//Only pure-go, platform independent modules should go in this file.

package modules

import (
	"io/fs"
	"os"
)

// Given a filepath, returns info about a file such as name, size etc.
func FileInfo(path string) (fs.FileInfo, error) {
	results, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return results, nil
}
