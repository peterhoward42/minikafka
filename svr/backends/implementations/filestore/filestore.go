package filestore

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"

	toykafka "github.com/peterhoward42/toy-kafka"
	"github.com/peterhoward42/toy-kafka/svr/backends/implementations/filestore/index"
)

var mutex = &sync.Mutex{} // Guards concurrent access of the FileStore.

// FileStore implements the svr/backends/contract/BackingStore interface using
// files on disk.
type FileStore struct {
	rootDir string
}

// NewFileStore creates and initializes a FileStore.
func NewFileStore(rootDir string) *FileStore {
	return &FileStore{rootDir}
}

// ------------------------------------------------------------------------
// METHODS TO SATISFY THE BackingStore INTERFACE.
// ------------------------------------------------------------------------

// DeleteContents removes all contents from the store.
func (s FileStore) DeleteContents() error {
	mutex.Lock()
	return s.deleteContents()
	defer mutex.Unlock()
}

// Store is defined by, and documented in the backends/contract/BackingStore
// interface.
func (s FileStore) Store(topic string, message toykafka.Message) (
	messageNumber int, err error) {

	mutex.Lock()
	return s.store(topic, message)
	defer mutex.Unlock()
}

// RemoveOldMessages is defined by, and documented in the
// backends/contract/BackingStore interface.
func (s FileStore) RemoveOldMessages(maxAge time.Time) (
	nRemoved int, err error) {
	return -1, nil
}

// Poll is defined by, and documented in the backends/contract/BackingStore
// interface.
func (s FileStore) Poll(topic string, readFrom int) (
	foundMessages []toykafka.Message, newReadFrom int, err error) {

	foundMessages = []toykafka.Message{}
	return foundMessages, 11, nil
}

// ------------------------------------------------------------------------
// Helper functions.
// ------------------------------------------------------------------------

func (s FileStore) deleteContents() error {
	dir, err := ioutil.ReadDir(s.rootDir)
	if err != nil {
		return fmt.Errorf("ioutil.ReadDir(): %v", err)
	}
	for _, entry := range dir {
		fullpath := path.Join(s.rootDir, entry.Name())
		err = os.RemoveAll(fullpath)
		if err != nil {
			return fmt.Errorf("os.RemoveAll(): %v", err)
		}
	}
	return nil
}

func (s FileStore) store(topic string, message toykafka.Message) (
	messageNumber int, err error) {

	index, err := s.retrieveIndexFromDisk()
	if err != nil {
		return -1, fmt.Errorf("RetrieveIndexFromDisk(): %v", err)
	}

	err = s.createTopicDirIfNotExists(topic)
	if err != nil {
		return -1, fmt.Errorf("createTopicDirIfNotExists: %v", err)
	}

	msgNumber := index.nextMessageNumberFor(topic)
	msgToStore := s.makeMsgToStore(message, msgNumber)
	msgSize := len(msgToStore)
	var fileBaseNameToUse string

	fileBaseNameToUse, needNewOne := s.selectFileBasenameForTopic(topic, msgSize, index)
	if needNewOne {
		fileBaseNameToUse, err = s.setupNewFileForTopic(topic, index)
		if err != nil {
			return -1, fmt.Errorf("setupNewFileForTopic: %v", err)
		}
	}
	err = s.saveAndRegisterMessage(fileBaseNameToUse, msgToStore, 
		msgNumber, index)
	if err != nil {
		return fmt.Errorf("saveAndRegisterMessage: %v", err)
	}
	err = index.Save()
	return msgNumber, nil
}

func (s FileStore) retrieveIndexFromDisk() (*StoreIndex, error) {
	indexPath := path.Join(s.rootDir, indexFileName)
	index, err := LoadStoreIndex(indexPath)
	if err != nil {
		return nil, fmt.Errorf("LoadStoreIndex(): %v", err)
	}
	return index, nil

}

func (s FileStore) createTopicDirIfNotExists(topic string) error {
	dirPath := s.directoryForTopic(topic)
	err := os.Mkdir(dirPath, 0777) // Todo what should permissions be?
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		return nil
	}
	return fmt.Errorf("os.Mkdir(): %v", err)
}

func (s FileStore) setupNewFileForTopic(topic string, index index.Index) (
	fileBaseNameToUse string, err error) {
	fileBaseName := fmt.Sprintf("%d", time.Now().UnixNano()/1000)
	filePath := path.Join(s.directoryForTopic(topic), fileBaseName)
	err = ioutil.WriteFile(filePath, []byte{}, "0777")
	if err != nil {
		return nil, fmt.Errorf("ioutil.WriteFile(): %v", err)
	}
	index.RegisterNewFile(topic, fileBaseName)
	return fileBaseName, nil
}

func (s FileStore) directoryForTopic(topic string) string {
	return path.Join(s.rootDir, topic)
}

func (s FileStore) makeMsgToStore(
	message toykafka.Message, msgNumber uint32) []byte {
	msg := storedMessage{message, time.Now(), msgNumber}
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	encoder.Encode(msg)
	return buf.Bytes()
}
func (s FileStore) selectFileBasenameForTopic(
	topic string, msgSize int, index StoreIndex) (
	fileBasenameToUse string, needNewOne bool) {

	// Does the index know of any files used for this topic?
	mostRecentFileBaseName := index.mostRecentFileFor(topic)
	if mostRecentFile == "" {
		return nil, true
	}
	// Does that file still exist? (Could have been cleared out).
	// And does it have enough room?
	filePath := s.fullPathToStorageFile(topic, mostRecentFile)
	if s.fileDoesNotExistOrHasInsufficientRoom(filePath) {
		return nil, true
	}
	return "", mostRecentFile
}

func (*s FileStore) saveAndRegisterMessage(fileToUse string, 
	topic string, msgToStore []byte, msgNumber int32, index index.Index) error { 
	filePath := s.fullPathToStorageFile(topic, fileToUse)
		// evaluate file full path
		// append write to it
		// mandate index to ack done incl updating next msg for topic
	}

type storedMessage struct {
	message       toykafka.Message
	creationTime  time.Time
	messageNumber uint32
}
