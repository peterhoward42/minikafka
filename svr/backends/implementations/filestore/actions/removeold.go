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

// RemoveOldMessages contains an optimisation as allowed by the interface,
// whereby it does not neccesarily remove all of the messages it is invited to.
// The optimisation is to only remove whole message files that are eligible 
// rather than crack any of them open. 
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
