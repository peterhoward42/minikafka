// Package actions is where the private implementation code lives for each
// of the main BackingStore actions. I.e. store/poll etc.
package actions

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"time"

	toykafka "github.com/peterhoward42/toy-kafka"
	"github.com/peterhoward42/toy-kafka/svr/backends/implementations/filestore/filenamer"
	"github.com/peterhoward42/toy-kafka/svr/backends/implementations/filestore/indexing"
	"github.com/peterhoward42/toy-kafka/svr/backends/implementations/filestore/ioutils"
	"github.com/peterhoward42/toy-kafka/svr/backends/implementations/filestore/stored"
)

const maximumFileSize = 1048576 // 1 MiB

// StoreAction encapsulates a single execution of the store (message) command.
type StoreAction struct {
	topic   string
	message toykafka.Message
	index   *indexing.Index
	rootDir string
}

// Store is the internal entry point function to store a new message in the
// filestore. Its responsibility to perform the storage operation and to update
// the in-memory index. It is not responsible for mutex protection, nor re-saving
// the index afterwards. These are the responsibility of the caller.
func (action StoreAction) Store() (messageNumber int, err error) {

	// Special case when the store has never encountered this topic before.
	err = action.createTopicDirIfNotExists()
	if err != nil {
		return -1, fmt.Errorf("createTopicDirIfNotExists: %v", err)
	}
	// Prepare the payload that will get stored, now, so we can measure how big
	// it is, to inform the decisions about whether we should spill over to a new
	// storage file.
	msgNumber := action.index.NextMessageNumberFor(action.topic)
	msgToStore := action.makeMsgToStore(action.message, msgNumber)
	msgSize := len(msgToStore)

	// Establish which storage file to use - including the case for needing to
	// start a new one.
	var msgFileName string
	msgFileName = action.index.CurrentMsgFileNameFor(action.topic)
	var needNewFile = false
	if msgFileName == "" {
		needNewFile = true
	} else {
		needNewFile, err = action.fileHasInsufficentRoom(msgFileName, msgSize)
		if err != nil {
			return -1, fmt.Errorf("fileHasInsufficietRoom(): %v", err)
		}
	}
	if needNewFile {
		msgFileName, err = action.setupNewFileForTopic()
		if err != nil {
			return -1, fmt.Errorf("setupNewFileForTopic(): %v", err)
		}
	}
	// Append the message object to the storage file, and mandate the
	// index to update itself with this new info.
	err = action.saveAndRegisterMessage(msgFileName, msgToStore, msgNumber)
	if err != nil {
		return -1, fmt.Errorf("saveAndRegisterMessage(): %v", err)
	}
	return int(msgNumber), nil
}

// makeMsgToStore builds a storedMessage structure to represent the
// incoming message and returns its byte-serialized form.
func (action *StoreAction) makeMsgToStore(
	message toykafka.Message, msgNumber int32) []byte {
	msg := stored.Message{message, time.Now(), msgNumber}
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	encoder.Encode(msg)
	return buf.Bytes()
}

// createTopicDirIfNotExists looks to see if a directory already exists
// for the given topic, and when not so, it creates one. It seeks the help of
// the filenamer module about file-naming rules.
func (action *StoreAction) createTopicDirIfNotExists() error {
	dirPath := filenamer.DirectoryForTopic(action.topic, action.rootDir)
	err := ioutils.CreateDirIfDoesntExist(dirPath)
	if err != nil {
		return fmt.Errorf("os.Mkdir(): %v", err)
	}
	return nil
}

// fileHasInsufficientRoom works out if the message storage file that is
// being used for a particular topic, has enough room to accomodate the
// incoming message without breaking the maximum file size limit.
func (action *StoreAction) fileHasInsufficentRoom(
	msgFileName string, msgSize int) (bool, error) {
	size, err := ioutils.FileSize(msgFileName)
	if err != nil {
		return false, fmt.Errorf("ioutils.FileSize(): %v", err)
	}
	insufficient := size+int64(msgSize) > maximumFileSize
	return insufficient, nil
}

// setupNewFileForTopic works out what the new file should be called, creates it,
// and then registers this new information with the index.
func (action *StoreAction) setupNewFileForTopic() (msgFileName string, err error) {
	fileName := filenamer.NewMsgFilenameFor(action.topic, action.index)
	filePath := filenamer.MessageFilePath(
		fileName, action.topic, action.rootDir)
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("os.Create(): %v", err)
	}
	err = file.Close()
	if err != nil {
		return "", fmt.Errorf("file.Close(): %v", err)
	}
	msgFileList := action.index.GetMessageFileListFor(action.topic)
	msgFileList.RegisterNewFile(fileName)
	return fileName, nil
}

// saveAndRegisteMessage appends the message to the specified file and updates
// the index with this new info. Note this is the point at which the
// message creation time is evauated and associated with the message.
func (action *StoreAction) saveAndRegisterMessage(
	msgFileName string, msgToStore []byte, msgNumber int32) error {
	filepath := filenamer.MessageFilePath(
		msgFileName, action.topic, action.rootDir)
	err := ioutils.AppendToFile(filepath, msgToStore)
	if err != nil {
		return fmt.Errorf("ioutils.AppendToFile(): %v", err)
	}
	creationTime := time.Now()
	msgFileList := action.index.GetMessageFileListFor(action.topic)
	fileMeta := msgFileList.Meta[msgFileName]
	fileMeta.RegisterNewMessage(msgNumber, creationTime)
	return nil
}
