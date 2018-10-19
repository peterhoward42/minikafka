// Package filestore provides a message storage system based on mounted file
// system. It implements the backingstore.contract.BackingStore interface.
package filestore

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"sync"
	"time"

	toykafka "github.com/peterhoward42/toy-kafka"
	"github.com/peterhoward42/toy-kafka/svr/backends/implementations/filestore/filenamer"
	"github.com/peterhoward42/toy-kafka/svr/backends/implementations/filestore/index"
	"github.com/peterhoward42/toy-kafka/svr/backends/implementations/filestore/ioutils"
)

const maximumFileSize = 1048576 // 1 MiB

var mutex = &sync.Mutex{} // Guards concurrent access of the FileStore.

// FileStore encapsulates the store.
type FileStore struct {
	rootDir string
}

// ------------------------------------------------------------------------
// METHODS TO SATISFY THE BackingStore INTERFACE.
// ------------------------------------------------------------------------

// DeleteContents removes all contents from the store.
func (s FileStore) DeleteContents() error {
	mutex.Lock()
	defer mutex.Unlock()
	return s.deleteContents()
}

// Store is defined by, and documented in the backends/contract/BackingStore
// interface.
func (s FileStore) Store(topic string, message toykafka.Message) (
	messageNumber int, err error) {

	mutex.Lock()
	defer mutex.Unlock()
	return s.store(topic, message)
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
	err := ioutils.DeleteDirectoryContents(s.rootDir)
	if err != nil {
		return fmt.Errorf("ioutils.DeleteDirectoryContents(): %v", err)
	}
	return nil
}

func (s FileStore) store(topic string, message toykafka.Message) (
	messageNumber int, err error) {

	currentIndex, err := index.RetrieveIndexFromDisk(
		filenamer.IndexFile(s.rootDir))
	if err != nil {
		return -1, fmt.Errorf("RetrieveIndexFromDisk(): %v", err)
	}
	err = s.createTopicDirIfNotExists(topic)
	if err != nil {
		return -1, fmt.Errorf("createTopicDirIfNotExists: %v", err)
	}
	msgNumber := storeIndex.NextMessageNumberFor(topic)
	msgToStore := s.makeMsgToStore(message, msgNumber)
	msgSize := len(msgToStore)

	var msgFileName string
	msgFileName = storeIndex.CurrentMsgFileNameFor(topic)
	var needNewFile = false
	if msgFileName == "" {
		needNewFile = true
	} else {
		needNewFile, err = s.fileHasInsufficentRoom(msgFileName, topic, msgSize)
		if err != nil {
			return -1, fmt.Errorf("fileHasInsufficietRoom(): %v", err)
		}
	}
	if needNewFile {
		msgFileName, err = s.setupNewFileForTopic(topic, index)
		if err != nil {
			return -1, fmt.Errorf("setupNewFileForTopic(): %v", err)
		}
	}
	err = s.saveAndRegisterMessage(
		msgFileName, topic, msgToStore, msgNumber, index)
	if err != nil {
		return -1, fmt.Errorf("saveAndRegisterMessage(): %v", err)
	}
	err = index.SaveIndex(index, filenamer.IndexFile(s.rootDir))
	if err != nil {
		return -1, fmt.Errorf("SaveIndex(): %v", err)
	}
	return int(msgNumber), nil
}

func (s FileStore) makeMsgToStore(
	message toykafka.Message, msgNumber int32) []byte {
	msg := storedMessage{message, time.Now(), msgNumber}
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	encoder.Encode(msg)
	return buf.Bytes()
}

func (s FileStore) createTopicDirIfNotExists(topic string) error {
	dirPath := filenamer.DirectoryForTopic(topic, s.rootDir)
	err := ioutils.CreateDirIfDoesntExist(dirPath)
	if err != nil {
		return fmt.Errorf("os.Mkdir(): %v", err)
	}
}

func (s FileStore) fileHasInsufficentRoom(
	msgFileName string, topic string, msgSize int) (bool, error) {
	size, err := ioutils.FileSize(msgFileName)
	if err != nil {
		return false, fmt.Errorf("ioutils.FileSize(): %v", err)
	}
	insufficient := size+int64(msgSize) > maximumFileSize
	return insufficient, nil
}

func (s FileStore) setupNewFileForTopic(
	topic string, index *index.Index) (msgFileName string, err error) {
	fileName := filenamer.NewMsgFilenameFor(topic, index)
	filepath := filenamer.MessageFilePath(fileName, topic, s.rootDir)

	file, err := os.Create(filepath)
	if err != nil {
		return false, fmt.Errorf("os.Create(): %v", err)
	}
	defer file.Close()

	msgFileList := index.GetMessageFileListFor(topic)
	msgFileList.RegisterNewFile(fileName)
	return fileName, nil
}

func (s FileStore) saveAndRegisterMessage(
	msgFileName string, topic string, msgToStore []byte,
	msgNumber int32, index *index.Index) err {
	filepath := filenamer.MessageFilePath(msgFileName, topic, s.rootDir)
	err = ioutils.AppendToFile(filepath, msgToStore)
	if err != nil {
		return fmt.Errorf("ioutils.AppendToFile(): %v", err)
	}
	creationTime := time.Now()
	msgFileList := index.GetMessageFileListFor(topic)
	fileMeta := msgFileList.Meta[msgFileName]
	fileMeta.RegisterNewMessage(msgNumber, creationTime)
	return nil
}

// ------------------------------------------------------------------------
// Auxilliary types.
// ------------------------------------------------------------------------

type storedMessage struct {
	message       toykafka.Message
	creationTime  time.Time
	messageNumber int32
}
