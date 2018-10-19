package ioutils

import (
)

// DeleteDirectoryContents removes everything from the given directory,
// retaining the directory itself.
func DeleteDirectoryContents(dir string) error {
	dir, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("ioutil.ReadDir(): %v", err)
	}
	for _, entry := range dir {
		fullpath := path.Join(dir, entry.Name())
		err = os.RemoveAll(fullpath)
		if err != nil {
			return fmt.Errorf("os.RemoveAll(): %v", err)
		}
	}
	return nil
}

// CreateDirIfDoesntExist 
func CreateDirIfDoesntExist(pathname) error {
	err := os.Mkdir(pathname, 0777)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		return nil
	}
	return fmt.Errorf("os.Mkdir(): %v", err)
}

func FileSize(pathname string) (int64, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return -1, fmt.Errorf("os.Open(): %v", err)
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return -1, fmt.Errorf("file.Stat(): %v", err)
	}
	return fileInfo.Size(), nil
}

func AppendToFile(filepath string, someData []bytes) error {
	file, err := os.OpenFile(filePath, os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("os.OpenFile(): %v", err)
	}
	defer file.Close()
    what is something?
	_, err := file.Write(someData)
	if err != nil {
		return fmt.Errorf("file.Write(): %v", err)
	}
    return nil
}
