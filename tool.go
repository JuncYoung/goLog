package goLog

import (
	"io/ioutil"
	"os"
)

func countDirFileNum(path string) (int, error) {
	dirList, err := ioutil.ReadDir(path)
	if err != nil {
		return -1, err
	}
	return len(dirList), nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
