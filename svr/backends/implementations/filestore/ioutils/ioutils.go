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

// CreateDirIfDoesntExist creates a directory with the given path,
// if one is not there already.
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

// Exists evaluates whether there is an entity in the file system at the
// given path. Note it does not guarantee that this is a file.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// CountEntitiesInDir provides the number of entities in the given directory.
func CountEntitiesInDir(dir string) (int, error) {
	entities, err := ioutil.ReadDir(dir)
	if err != nil {
		return -1, fmt.Errorf("ioutil.ReadDir(): %v", err)
	}
	return len(entities), nil
}
