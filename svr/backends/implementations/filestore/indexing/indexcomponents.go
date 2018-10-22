package indexing

import (
	"time"
)

// The types' fields are exported so they can be automatically gob-encoded
// without bothering with structure tags.

//-----------------------------------------------------------------------

// MsgMeta holds the message number and creation time for one stored message.
type MsgMeta struct {
	MsgNum  int32 // Zero value of 0 used to signify uninitialised.
	Created time.Time
}

//-----------------------------------------------------------------------

// FileMeta holds information about the oldest and newest message in
// one message file, and its current size.
type FileMeta struct {
	Oldest MsgMeta
	Newest MsgMeta
	Size   int
}

// RegisterNewMessage updates the FileMeta object according to this new
// message arriving in the store.
func (fm *FileMeta) RegisterNewMessage(
	msgNumber int32, creationTime time.Time, messageSize int) {
	// Special case; set the Oldest field if this is the first
	// message to arrive.
	if fm.Oldest.MsgNum == int32(0) {
		fm.Oldest = MsgMeta{msgNumber, creationTime}
	}
	fm.Newest = MsgMeta{msgNumber, creationTime}
	fm.Size += messageSize
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
