package ioutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// DeleteDirectoryContents removes everything from the given directory,
// retaining the directory itself.
func DeleteDirectoryContents(dir string) error {
	dirInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("ioutil.ReadDir(): %v", err)
	}
	for _, entry := range dirInfo {
		fullpath := path.Join(dir, entry.Name())
		err = os.RemoveAll(fullpath)
		if err != nil {
			return fmt.Errorf("os.RemoveAll(): %v", err)
		}
	}
	return nil
}

// CreateDirIfDoesntExist
func CreateDirIfDoesntExist(path string) error {
	err := os.Mkdir(path, 0777)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		return nil
	}
	return fmt.Errorf("os.Mkdir(): %v", err)
}

// AppendToFile appends some bytes to the specified file, and re-closes it.
func AppendToFile(filepath string, someData []byte) error {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return fmt.Errorf("os.OpenFile(): %v", err)
	}
	defer file.Close()
	_, err = file.Write(someData)
	if err != nil {
		return fmt.Errorf("file.Write(): %v", err)
	}
	return nil
}
