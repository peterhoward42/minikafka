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
func (fm *FileMeta) RegisterNewMessage(msgNumber int32, messageSize int64) {

	fm.SeekOffsetForMessageNumber[msgNumber] = fm.Size
	fm.Size += messageSize

	creationTime := time.Now()

	// Special case, when this is the first message to arrive for the file.
	if fm.Oldest.MsgNum == int32(0) {
		fm.Oldest.MsgNum = msgNumber
		fm.Oldest.Created = creationTime
	}
	fm.Newest = MsgMeta{msgNumber, creationTime}
}
