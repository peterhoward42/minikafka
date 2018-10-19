package index

import (
	"time"
)

// The types' fields are exported so they can be automatically gob-encoded
// without bothering with structure tags.

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
