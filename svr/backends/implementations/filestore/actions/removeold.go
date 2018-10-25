// Package actions is where the private implementation code lives for each
// of the main BackingStore actions. I.e. store/removeold/poll etc.
package actions

import (
	"fmt"
	"os"
	"time"

	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/filenamer"
	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/indexing"
)

// RemoveOldMessagesAction encapsulates a single execution of the
// remove-old-messages command.
type RemoveOldMessagesAction struct {
	MaxAge  time.Time
	Index   *indexing.Index
	RootDir string
}

// RemoveOldMessages works out which message files can be viewed as spent
// because they hold only messages  older than the threshold specified,
// physically removes these files and updates the index accordingly.
func (action RemoveOldMessagesAction) RemoveOldMessages() (
	filesRemoved []string, nMessagesRemoved int, err error) {
	filesRemoved = []string{}
	nMessagesRemoved = 0
	// Handle the action on a per-topic basis.
	for topic, msgFileList := range action.Index.MessageFileLists {
		// Capture the files to delete and how many messages they had
		// in them.
		oldFiles := msgFileList.SpentFiles(action.MaxAge)
		for _, fileName := range oldFiles {
			nMessages := msgFileList.NumMessagesInFile(fileName)
			nMessagesRemoved += nMessages
		}
		filesRemoved = append(filesRemoved, oldFiles...)
		// Mandate the index to forget about these files.
		msgFileList.ForgetFiles(oldFiles)
		// Physically remove the files.
		for _, fileName := range oldFiles {
			filePath := filenamer.MessageFilePath(
				fileName, topic, action.RootDir)
			err = os.Remove(filePath)
			if err != nil {
				return nil, -1, fmt.Errorf("os.Remove(): %v", err)
			}
		}
	}
	return filesRemoved, nMessagesRemoved, nil
}
