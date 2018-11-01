// Package actions is where the private implementation code lives for each
// of the main BackingStore actions. I.e. store/removeold/poll etc.
package actions

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/peterhoward42/minikafka"
	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/filenamer"
	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/indexing"
)

// PollAction encapsulates a single execution of the Poll command.
type PollAction struct {
	Topic    string
	ReadFrom int
	Index    *indexing.Index
	RootDir  string
}

// Poll is the internal entry point function to poll for messages beyond a given
// message number.  Its responsibility is to perform the poll operation. It is
// not responsible for mutex protection,
func (action PollAction) Poll() (
	foundMessages []minikafka.Message, newReadFrom int, err error) {

	// Access the topic-specific indexing information.
	msgFileList, ok := action.Index.MessageFileLists[action.Topic]
	if ok == false {
		return nil, -1, fmt.Errorf("Unknown topic: %v", action.Topic)
	}

	// Which message storage files must we look in?
	messageNumberToReadFrom := action.ReadFrom
	fileNames := msgFileList.MessageFilesForMessagesFrom(
		messageNumberToReadFrom)

	// If there are none, return benign data.
	if len(fileNames) == 0 {
		return []minikafka.Message{}, int(action.ReadFrom), nil
	}

	// Harvest the messages from this list of files.
	messages := []minikafka.Message{}
	for _, fileName := range fileNames {
		messages, err = action.addMessagesFromFile(
			messages, fileName, int32(messageNumberToReadFrom))
		if err != nil {
			return nil, -1, fmt.Errorf("action.AddMessagesFromFile(): %v", err)
		}
	}

	newReadFrom = int(action.Index.NextMessageNumbers[action.Topic])

	return messages, newReadFrom, nil
}

// addMessagesFromFile appends all the messages in the file beyond (incl.)
// messageNumberToReadFrom, to the addTo slice, and returns it.
func (action PollAction) addMessagesFromFile(
	addTo []minikafka.Message, fileName string, messageNumberToReadFrom int32) (
	[]minikafka.Message, error) {

	// Read the file contents into memory.
	filePath := filenamer.MessageFilePath(fileName, action.Topic, action.RootDir)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("os.Open(): %v", err)
	}
	defer file.Close()
	fileContents, err := ioutil.ReadAll(file)

	// Which message numbers should we harvest?
	msgFileList, _ := action.Index.MessageFileLists[action.Topic]
	fileMeta := msgFileList.Meta[fileName]
	startMsgNum := messageNumberToReadFrom
	if fileMeta.Oldest.MsgNum > startMsgNum {
		startMsgNum = fileMeta.Oldest.MsgNum
	}
	endMsgNum := fileMeta.Newest.MsgNum

	// For each targeted message number, harvest the slice of bytes in the
	// file that represents it.
	for msgNum := startMsgNum; msgNum <= endMsgNum; msgNum++ {
		start := fileMeta.SeekOffsetForMessageNumber[msgNum]
		end, ok := fileMeta.SeekOffsetForMessageNumber[msgNum+1]
		if ok == false {
			end = int64(len(fileContents))
		}
		msgBytes := fileContents[start:end]
		addTo = append(addTo, msgBytes)
	}

	return addTo, nil
}
