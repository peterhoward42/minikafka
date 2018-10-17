// Package index keeps track of which message filenames have been used for each
// topic, and for each, which range of message numbers and creation times they
// hold. It deals only in file basenames, and has no awareness of actual storage
// locations, nor involvement in IO. However it does to offer to serialize
// itself to and from a Reader/Writer, so that clients can persist it.
package index

import (
	"encoding/gob"
	"fmt"
	"io"
	"time"
)

// The types' fields are exported only so that they can easily be gob-encoded.

//-----------------------------------------------------------------------

// MsgMeta holds the message number and creation time for one stored message.
type MsgMeta struct {
	MsgNum  int32
	Created time.Time
}

//-----------------------------------------------------------------------

// FileMeta holds information about the oldest and newest message in
// one message file.
type FileMeta struct {
	Oldest MsgMeta
	Newest MsgMeta
}

//-----------------------------------------------------------------------

// MessageFileList holds information about the set of files that
// contain one topic's messages.
type MessageFileList struct {
	Names []string
	Meta  map[string]*FileMeta
}

// NewMessageFileList creates and initializes a MessageFileList.
func NewMessageFileList() *MessageFileList {
	return &MessageFileList{
		[]string{},
		map[string]*FileMeta{},
	}
}

// RegisterNewFile .
func (lst *MessageFileList) RegisterNewFile(filename string) {
	lst.Names = append(lst.Names, filename)
	lst.Meta[filename] = &FileMeta{}
}

//-----------------------------------------------------------------------

// Index is the top level index object.
type Index struct {
	MessageFileLists map[string]*MessageFileList
}

// NewIndex creates and initialized an Index.
func NewIndex() *Index {
	return &Index{map[string]*MessageFileList{}}
}

// GetMessageFileListFor provides access to the MesageFileList for the
// given topic. Copes with the topic being hithertoo unknown.
func (index *Index) GetMessageFileListFor(topic string) *MessageFileList {
	_, ok := index.MessageFileLists[topic]
	if ok != true {
		index.MessageFileLists[topic] = NewMessageFileList()
	}
	return index.MessageFileLists[topic]
}

// Encode is a serializer. It encodes the index into a byte stream and writes
// them to the output writer provided. See also the Decode sister method.
func (index *Index) Encode(writer io.Writer) error {
	encoder := gob.NewEncoder(writer)
	err := encoder.Encode(index)
	if err != nil {
		return fmt.Errorf("encoder.Encode(): %v", err)
	}
	return nil
}

// Decode is a de-serializer. It populates the index by decoding the bytes
// read from the input reader provided. See also the Encode sister method.
func (index *Index) Decode(reader io.Reader) error {
	decoder := gob.NewDecoder(reader)
	err := decoder.Decode(index)
	if err != nil {
		return fmt.Errorf("decoder.Decode: %v", err)
	}
	return nil
}

// NextMessageNumberFor provides the next-available message number for a topic.
// It copes with the special cases of the index having no record of that topic,
// or it having never yet contained any messages.
func (index Index) NextMessageNumberFor(topic string) int32 {
	messageFileList, ok := index.MessageFileLists[topic]
	if ok == false {
		return 1
	}
	nTopicFiles := len(messageFileList.Names)
	if nTopicFiles == 0 {
		return 1
	}
	// Consult the meta data for the newest file.
	newestName := messageFileList.Names[nTopicFiles-1]
	newestFileMeta := messageFileList.Meta[newestName]
	return newestFileMeta.Newest.MsgNum + 1
}

// CurrentMsgFileNameFor .
func (index Index) CurrentMsgFileNameFor(topic string) string {
    msgFileList, ok :=  index.MessageFileLists[topic]
    if ok == false {
        return ""
    }
    if len(msgFileList.Names) == 0 {
        return ""
    }
    n := len(msgFileList.Names)
    return msgFileList.Names[n-1]
}


//-----------------------------------------------------------------------
