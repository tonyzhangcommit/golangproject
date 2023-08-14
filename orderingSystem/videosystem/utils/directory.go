package utils

import (
	"os"
	"io/fs"
	
)
type FileMode = fs.FileMode

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateDir(path string, perm FileMode) (err error) {
	if ok, _ := PathExists(path); !ok {
		err = os.Mkdir(path, perm)
	}
	return
}
