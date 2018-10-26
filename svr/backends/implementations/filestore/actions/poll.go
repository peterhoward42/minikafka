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

type PollAction struct {
    Topic string
    ReadFrom int
    Index indexing.Index
    RootDir string
}

// Poll is the internal entry point function to poll for messages beyond a given
// message number.  Its responsibility is to perform the poll operation. It is
// not responsible for mutex protection, 
func (action PollAction) Poll() (
    foundMessages []minikafka.Message, newReadFrom int, err error) {

    // Access the topic-specific indexing information.
    msgFileList, ok := action.Index.MessageFilesLists[action.Topic]
    if ok == false {
        return nil, -1, fmt.Errorf("Unknown topic: %v", topic)
    }

    // Which message storage files must we look in?
    fileNames := msgFileList.ReadFromFiles(action.ReadFrom)

    // If there are none return benign data.
    if len(fileNames) == 0 {
        return []minikafka.Message{}, action.ReadFrom, nil)
    }

    // For the first file, whereabouts in the file should we start reading?
    firstFile := fileNames[0]
    fileMeta := msgFileList.Meta[firstFile]
    firstFileSeekIndex := fileMeta.SeekOffsetForMessageNumber[action.ReadFrom]

    // Harvest the messages from this list of files.
    messages := []minikafka.Message{}
    for _, fileName := range(fileNames) {
        filepath := filenamer.MessageFilePath(
            fileName, action.topic, action.RootDir)
        var seekIndex = 0
        if fileName == firstFile {
            seekIndex = firstFileSeekIndex
        }
        messages, err := action.AddMessagesFromFile(
            messages, filepath, seekIndex)
        if err != nil {
            return nil, -1, fmt.Errorf("action.AddMessagesFromFile(): %v", err)
        }
    }

    newReadFrom = action.Index.NextMessageNumberFor(action.topic)

    return messages, newReadFrom, nil
}
