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
)

var mutex = &sync.Mutex{} // Guards concurrent access of the FileStore.

// FileStore implements the svr/backends/contract/BackingStore interface using
// files on disk.
type FileStore struct {
	rootDir string
}

// NewFileStore instantiates, initializes and returns a FileStore.
func NewFileStore(rootDir string) (*FileStore, error) {
	// Make the root directory if it does not exist already.
	// That is all that is needed to be a viable empty store.
	err := os.MkdirAll(rootDir, 0777)
	// Tolerate already-exists error, but non others.
	if err != nil && os.IsExist(err) == false {
		return nil, fmt.Errorf("os.MkdirAll(): %v", err)
	}
	return &FileStore{rootDir: rootDir}, nil
}

// ------------------------------------------------------------------------
// METHODS TO SATISFY THE BackingStore INTERFACE.
// ------------------------------------------------------------------------

// DeleteContents removes all contents from the store.
func (s FileStore) DeleteContents() error {
	mutex.Lock()
	defer mutex.Unlock()
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

// Store is defined by, and documented in the backends/contract/BackingStore
// interface.
func (s FileStore) Store(topic string, message toykafka.Message) (
	messageNumber int, err error) {

	mutex.Lock()
	defer mutex.Unlock()

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
	var fileToUse string
	if s.needNewFileForTopic(topic, msgSize, index) {
		fileToUse, err = s.setupNewFileForTopic(topic, index)
		if err != nil {
			return fmt.Errorf("setupNewFileForTopic: %v", err)
		}
	}
	err = s.saveAndRegisterMessage(fileToUse, msgToStore, msgNumber, index)
	if err != nil {
		return fmt.Errorf("saveAndRegisterMessage: %v", err)
	}
	err = index.Save()
	return msgNumber, nil
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

func (s FileStore) retrieveIndexFromDisk() (*StoreIndex, error) {
	const indexFileName = "index"
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

type storedMessage struct {
	message       toykafka.Message
	creationTime  time.Time
	messageNumber uint32
}
