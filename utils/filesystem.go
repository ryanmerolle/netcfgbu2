package utils

import (
	"os"
)

// EnsureDir ensures that a directory exists, and creates it if it doesn't.
func EnsureDir(dirName string) error {
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
