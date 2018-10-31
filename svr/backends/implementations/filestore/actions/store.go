// Package actions is where the private implementation code lives for each
// of the main BackingStore actions. I.e. store/poll etc.
package actions

import (
	"fmt"
	"os"

	minikafka "github.com/peterhoward42/minikafka"
	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/filenamer"
	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/indexing"
	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/ioutils"
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

	// Special case when the store has never stored a message for this
	// this topic before.
	err = action.createTopicDirIfNotExists()
	if err != nil {
		return -1, "", fmt.Errorf("createTopicDirIfNotExists(): %v", err)
	}

	// Establish which storage file to use - including the case for needing to
	// start a new one.
	var msgFileName string
	msgFileName = action.Index.CurrentMsgFileNameFor(action.Topic)
	var needNewFile = false
	if msgFileName == "" {
		needNewFile = true
	} else {
		needNewFile = action.fileHasInsufficentRoom(msgFileName)
	}
	if needNewFile {
		msgFileName, err = action.setupNewFileForTopic()
		if err != nil {
			return -1, "", fmt.Errorf("setupNewFileForTopic(): %v", err)
		}
	}
	// Append the message bytes to the storage file, and mandate the
	// index to update itself with this new info.
	messageNumber, err = action.saveAndRegisterMessage(msgFileName)
	if err != nil {
		return -1, "", fmt.Errorf("saveAndRegisterMessage(): %v", err)
	}
	return messageNumber, msgFileName, nil
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

func (action *StoreAction) fileHasInsufficentRoom(msgFileName string) bool {
	msgFileList := action.Index.MessageFileLists[action.Topic]
	msgSize := int64(len(action.Message))
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
// the index with this new info.
func (action *StoreAction) saveAndRegisterMessage(
	msgFileName string) (msgNumber int, err error) {
	filepath := filenamer.MessageFilePath(
		msgFileName, action.Topic, action.RootDir)
	err = ioutils.AppendToFile(filepath, action.Message)
	if err != nil {
		return 0, fmt.Errorf("ioutils.AppendToFile(): %v", err)
	}
	msgNumber = int(action.Index.GetAndIncrementMessageNumberFor(action.Topic))
	msgFileList := action.Index.GetMessageFileListFor(action.Topic)
	fileMeta := msgFileList.Meta[msgFileName]
	fileMeta.RegisterNewMessage(int32(msgNumber), int64(len(action.Message)))
	return msgNumber, nil
}
