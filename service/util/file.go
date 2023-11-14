package util

import (
	"os"
)

func FileExists(fileName string) bool {
	if fileName == "" {
		return false
	}

	info, err := os.Stat(fileName)

	return err == nil && !info.IsDir()
}
