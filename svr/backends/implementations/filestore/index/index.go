package index

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"
)

//-----------------------------------------------------------------------
// The types in here form collectively the in-memory index into the files that
// messages are stored in. Every FileStore API method, is expected to
// de-hydrate the index from disk in order to do its job, and then to update and
// re-save the index accordingly before returning.
//
// The types' fields are exported only so that they can be gob-encoded.
//-----------------------------------------------------------------------

// MsgMeta holds the message number and creation time for one stored message.
type MsgMeta struct {
	MsgNum  int32
	Created time.Time
}

// FileMeta holds information about the oldest and newest message in
// one file.
type FileMeta struct {
	FileName  string
	OldestMsg MsgMeta
	NewestMsg MsgMeta
}

// MessageFileList holds information about the set of files that
// contain one topic's messages.
type MessageFileList []FileMeta // Ordered by file creation sequence.

// Index is the top level index object for the file store.
type Index struct {
	MessageFileLists map[string]MessageFileList
}

//-----------------------------------------------------------------------
// API methods
//-----------------------------------------------------------------------

// NewIndex provides a zero-value Index instance.
func NewIndex() *Index {
	return &Index{map[string]MessageFileList{}}
}

// Save writes a representation of the Index to disk.
func (index *Index) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("os.Create(): %v", err)
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(index)
	if err != nil {
		return fmt.Errorf("encoder.Encode(): %v", err)
	}
	return nil
}

// PopulateFromDisk de-serializes the contents of a Index from disk,
// and populates the instance accordingly.
func (index *Index) PopulateFromDisk(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("os.Open(): %v", err)
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(index)
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
	nTopicFiles := len(messageFileList)
	if nTopicFiles == 0 {
		return 1
	}
	// Consult the meta data for the newest file.
	newestFileMeta := messageFileList[nTopicFiles-1]
	return newestFileMeta.NewestMsg.MsgNum + 1
}
