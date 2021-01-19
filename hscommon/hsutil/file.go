// Package hsutil holds utility and helper functions
package hsutil

import (
	"log"
	"os"
	"path/filepath"
)

const (
	defaultFilePermissions os.FileMode = 0755
)

func CreateFileAtPath(pathToFile string, data []byte) bool {
	folder := filepath.Dir(pathToFile)

	err := os.MkdirAll(folder, defaultFilePermissions)
	if err != nil && !os.IsExist(err) {
		// if this error isn't just "the folder already exists", something went wrong
		log.Printf("failed to mkdir %s: %s", folder, err)
		return false
	}

	newFile, err := os.Create(pathToFile)
	if err != nil {
		log.Printf("failed to create new file %s: %s", pathToFile, err)
		return false
	}
	defer newFile.Close()

	bytesWritten, err := newFile.Write(data)
	if err != nil || bytesWritten != len(data) {
		log.Printf("failed to write data to new file %s, wrote %d bytes out of %d: %s",
			pathToFile, bytesWritten, len(data), err)
		return false
	}

	return true
}
