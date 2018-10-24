// Package actions is where the private implementation code lives for each
// of the main BackingStore actions. I.e. store/removeold/poll etc.
package actions

import (
	"fmt"
	"os"
	"time"

	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/filenamer"
	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/indexing"
	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/ioutils"
)

// RemoveOldMessagesAction encapsulates a single execution of the
// remove-old-messages command.
type RemoveOldMessagesAction struct {
	MaxAge  time.Time
	Index   indexing.Index
	RootDir string
}

// RemoveOldMessages works out which message files can be viewed as spent
// because they hold only messages  older than the threshold specified,
// physically removes these files and updates the index accordingly.
func (action RemoveOldMessagesAction) RemoveOldMessages() (
	nRemoved int, err error) {
	nRemoved = 0
	// Handle the action on a per-topic basis.
	for topic, msgFileList := range action.Index.MessageFileLists {
		// Capture the files to delete and how many messages they had
		// in them.
		oldFiles := msgFileList.SpentFiles(action.MaxAge)
		for file := range oldFiles {
			count := msgFileList.NumMessagesIn(file)
			nRemoved += count
		}
		// Mandate the index to forget about these files.
		msgFileList.ForgetFiles(oldFiles)
		// Physically remove the files.
		for _, fileName := range oldFiles {
			filePath := filenamer.MessageFilePath(
				fileName, topic, action.RootDir)
			err = ioutils.RemoveFile(filePath)
			if err != nil {
				return -1, fmt.Errorf("ioutils.RemoveFile(): %v", err)
			}
		}
	}
	return nRemoved, nil
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
	msgFileName string, msgSize int) bool {
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
// message creation time is evauated and associated with the message.
func (action *StoreAction) saveAndRegisterMessage(
	msgFileName string, msgToStore []byte, msgNumber int32) error {
	filepath := filenamer.MessageFilePath(
		msgFileName, action.Topic, action.RootDir)
	err := ioutils.AppendToFile(filepath, msgToStore)
	if err != nil {
		return fmt.Errorf("ioutils.AppendToFile(): %v", err)
	}
	creationTime := time.Now()
	msgFileList := action.Index.GetMessageFileListFor(action.Topic)
	fileMeta := msgFileList.Meta[msgFileName]
	fileMeta.RegisterNewMessage(msgNumber, creationTime, len(msgToStore))
	return nil
}
