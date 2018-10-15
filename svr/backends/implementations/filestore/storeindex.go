package filestore

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"
)

//-----------------------------------------------------------------------
// The types from which the index is (hierarchically) composed.
// Described bottom-up. These export their fields to enable automatic
// serialization by gob.
//-----------------------------------------------------------------------

// MsgMeta holds the message number and creation time for stored message.
type MsgMeta struct {
	MsgNum  int32
	Created time.Time
}

// FileMeta holds information about the oldest and newest message in
// a file.
type FileMeta struct {
	FileName  string
	OldestMsg MsgMeta
	NewestMsg MsgMeta
}

// TopicMessages holds information about the set of files that
// contain a topic's messages.
type TopicMessages []FileMeta // Ordered by file creation sequence.

// StoreIndex provides an index over the entire store by providing a map
// of topic names to their corresponding TopicMessages object.
type StoreIndex map[string]TopicMessages

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
