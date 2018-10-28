// Package actions is where the private implementation code lives for each
// of the main BackingStore actions. I.e. store/poll etc.
package actions

import (
	"bytes"
	"fmt"
	"os"
	"time"

	minikafka "github.com/peterhoward42/minikafka"
	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/filenamer"
	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/indexing"
	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/ioutils"
	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/stored"
)

const maximumFileSize = 1048576 // 1 MiB

// StoreAction encapsulates a single execution of the store (message) command.
type StoreAction struct {
	Topic   string
	Message minikafka.Message
	Index   *indexing.Index
	RootDir string
}

// Store is the internal entry point function to store a new message in the
// filestore. Its responsibility to perform the storage operation and to update
// the in-memory index. It is not responsible for mutex protection, nor re-saving
// the index afterwards. These are the responsibility of the caller.
func (action StoreAction) Store() (
	messageNumber int, msgFileUsed string, err error) {

	// Special case when the store has never encountered this topic before.
	err = action.createTopicDirIfNotExists()
	if err != nil {
		return -1, "", fmt.Errorf("createTopicDirIfNotExists(): %v", err)
	}
	// Prepare the payload that will get stored, now, so we can measure how big
	// it is, to inform the decisions about whether we should spill over to a new
	// storage file.
	msgNumber := action.Index.NextMessageNumberFor(action.Topic)
	storedMessage := stored.Message{
		Message:       action.Message,
		CreationTime:  time.Now(),
		MessageNumber: int32(msgNumber),
	}
	var buf bytes.Buffer
	err = storedMessage.Encode(&buf)
	if err != nil {
		return -1, "", fmt.Errorf("storedMessage.Encode(): %v", err)
	}
	bytesToStore := buf.Bytes()
	msgSize := int64(len(bytesToStore))

	// Establish which storage file to use - including the case for needing to
	// start a new one.
	var msgFileName string
	msgFileName = action.Index.CurrentMsgFileNameFor(action.Topic)
	var needNewFile = false
	if msgFileName == "" {
		needNewFile = true
	} else {
		needNewFile = action.fileHasInsufficentRoom(msgFileName, msgSize)
	}
	if needNewFile {
		msgFileName, err = action.setupNewFileForTopic()
		if err != nil {
			return -1, "", fmt.Errorf("setupNewFileForTopic(): %v", err)
		}
	}
	// Append the message object to the storage file, and mandate the
	// index to update itself with this new info.
	err = action.saveAndRegisterMessage(msgNumber, msgFileName, bytesToStore)
	if err != nil {
		return -1, "", fmt.Errorf("saveAndRegisterMessage(): %v", err)
	}
	return int(msgNumber), msgFileName, nil
}

// createTopicDirIfNotExists looks to see if a directory already exists
// for the given topic, and when not so, it creates one. It seeks the help of
// the filenamer module about file-naming rules.
func (action *StoreAction) createTopicDirIfNotExists() error {
	dirPath := filenamer.DirectoryForTopic(action.Topic, action.RootDir)
	err := ioutils.CreateDirIfDoesntExist(dirPath)
	if err != nil {
		return fmt.Errorf("os.Mkdir(): %v", err)
	}
	return nil
}

func (action *StoreAction) fileHasInsufficentRoom(
	msgFileName string, msgSize int64) bool {
	msgFileList := action.Index.MessageFileLists[action.Topic]
	return msgFileList.Meta[msgFileName].Size+msgSize > maximumFileSize
}

// setupNewFileForTopic works out what the new file should be called, creates it,
// and then registers this new information with the index.
func (action *StoreAction) setupNewFileForTopic() (msgFileName string, err error) {
	fileName := filenamer.NewMsgFilenameFor(action.Topic, action.Index)
	filePath := filenamer.MessageFilePath(
		fileName, action.Topic, action.RootDir)
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("os.Create(): %v", err)
	}
	err = file.Close()
	if err != nil {
		return "", fmt.Errorf("file.Close(): %v", err)
	}
	msgFileList := action.Index.GetMessageFileListFor(action.Topic)
	msgFileList.RegisterNewFile(fileName)
	return fileName, nil
}

// saveAndRegisteMessage appends the message to the specified file and updates
// the index with this new info. Note this is the point at which the
// message creation time is evaluated and associated with the message.
func (action *StoreAction) saveAndRegisterMessage(
	msgNumber int, msgFileName string, msgToStore []byte) error {
	filepath := filenamer.MessageFilePath(
		msgFileName, action.Topic, action.RootDir)
	err := ioutils.AppendToFile(filepath, msgToStore)
	if err != nil {
		return fmt.Errorf("ioutils.AppendToFile(): %v", err)
	}
	creationTime := time.Now()
	msgFileList := action.Index.GetMessageFileListFor(action.Topic)
	fileMeta := msgFileList.Meta[msgFileName]
	fileMeta.RegisterNewMessage(msgNumber, creationTime, int64(len(msgToStore)))
	return nil
}
