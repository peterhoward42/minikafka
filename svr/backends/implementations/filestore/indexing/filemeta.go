package indexing

import (
	"time"
)

// The types' fields are exported so they can be automatically gob-encoded
// without bothering with structure tags.

//-----------------------------------------------------------------------

// FileMeta holds information about the oldest and newest message in
// one message file, its current size, and the file-seek-offsets at which each
// message starts.
type FileMeta struct {
	Oldest                     MsgMeta
	Newest                     MsgMeta
	Size                       int64
	SeekOffsetForMessageNumber map[int32]int64
}

// NewFileMeta provides an initialised FileMeta, ready to use.
func NewFileMeta() *FileMeta {
	return &FileMeta{SeekOffsetForMessageNumber: map[int32]int64{}}
}

// RegisterNewMessage updates the FileMeta object according to this new
// message arriving in the store.
func (fm *FileMeta) RegisterNewMessage(
	msgNumber int, creationTime time.Time, messageSize int64) {

	// Capture the seek offset for the incoming message before we mutate
	// the data structure.
	fm.SeekOffsetForMessageNumber[int32(msgNumber)] = fm.Size
	fm.Size += messageSize

	// Special case, when this is the first message to arrive for the file.
	if fm.Oldest.MsgNum == int32(0) {
		fm.Oldest = MsgMeta{int32(msgNumber), creationTime}
	}
	fm.Newest = MsgMeta{int32(msgNumber), creationTime}
}
