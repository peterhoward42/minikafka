// Package actions is where the private implementation code lives for each
// of the main BackingStore actions. I.e. store/removeold/poll etc.
package actions

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"

	"github.com/peterhoward42/minikafka"
	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/filenamer"
	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/indexing"
	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/stored"
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
	fileNames := msgFileList.FilesContainingThisMessageAndNewer(
		messageNumberToReadFrom)

	// If there are none, return benign data.
	if len(fileNames) == 0 {
		return []minikafka.Message{}, int(action.ReadFrom), nil
	}

	// For the first file, whereabouts in the file should we start reading?
	firstFile := fileNames[0]
	fileMeta := msgFileList.Meta[firstFile]
	firstFileSeekOffset := fileMeta.SeekOffsetForMessageNumber[int32(action.ReadFrom)]

	// Harvest the messages from this list of files.
	messages := []minikafka.Message{}
	var seekOffset int64
	for _, fileName := range fileNames {
		filepath := filenamer.MessageFilePath(
			fileName, action.Topic, action.RootDir)
		seekOffset = 0             // The general case.
		if fileName == firstFile { // Special case for first file.
			seekOffset = firstFileSeekOffset
		}
		messages, err = action.addMessagesFromFile(
			messages, filepath, seekOffset)
		if err != nil {
			return nil, -1, fmt.Errorf("action.AddMessagesFromFile(): %v", err)
		}
	}

	newReadFrom = action.Index.NextMessageNumberFor(action.Topic)

	return messages, newReadFrom, nil
}

// addMessagesFromFile appends all the messages that can be read from the given
// file to the slice  given and returns it. The caller can pass in a seek offset
// to cause the function to start reading messages from that offset.
func (action PollAction) addMessagesFromFile(
	addTo []minikafka.Message, filepath string, seekOffset int64) (
	[]minikafka.Message, error) {

	// Open the file specified and seek to the position requested.
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("os.Open(): %v", err)
	}
	defer file.Close()
	_, err = file.Seek(seekOffset, 0)
	if err != nil {
		return nil, fmt.Errorf("file.Seek(): %v", err)
	}

	// Repeatedly do a gob.Decode from the file, which will generate a sequence
	// of stored.Message. Keep going until EOF. In each iteration, harvest the
	// minikafka.Message inside the stored.Message, accumulating them in the
	// caller's slice container provided.
	for {
		decoder := gob.NewDecoder(file)
		var storedMessage stored.Message
		err = decoder.Decode(&storedMessage)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("gob.Decoder.Decode(): %v", err)
		}
		addTo = append(addTo, storedMessage.Message)
	}
	return addTo, nil
}
