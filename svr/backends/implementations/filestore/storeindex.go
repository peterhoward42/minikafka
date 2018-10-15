package filestore

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"
)

//-----------------------------------------------------------------------
// The types from which the store index is (hierarchically) composed are
// defined here bottom-up. These export their fields to enable automatic
// serialization by gob.
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

// TopicFiles holds information about the set of files that
// contain a topic's messages.
type TopicFiles []FileMeta // Ordered by file creation sequence.

// StoreIndex provides an index over the entire store by providing a map
// of topic names to their corresponding TopicFiles object.
type StoreIndex map[string]TopicFiles

//-----------------------------------------------------------------------
// API methods
//-----------------------------------------------------------------------

// Save writes a representation of itself to disk.
func (index *StoreIndex) Save(path string) error {
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

// LoadStoreIndex deserializes a StoreIndex from disk and returns it.
func LoadStoreIndex(fileName string) (*StoreIndex, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("os.Open(): %v", err)
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	var index StoreIndex
	err = decoder.Decode(&index)
	if err != nil {
		return nil, fmt.Errorf("decoder.Decode: %v", err)
	}
	return &index, nil
}

func (index StoreIndex) nextMessageNumberFor(topic string) uint32 {
	topicFiles, ok := index[topic]
	if ok == false { // Haven't heard of this topic before.
		return 1
	}
	nTopicFiles := len(topicFiles) // No files ever created for this topic.
	if nTopicFiles == 0 {
		return 1
	}
	// Consult the meta data for the newest file.
	newestFileMeta := topicFiles[nTopicFiles-1]
	return newestFileMeta.NewestMsg.MsgNum + 1
}
