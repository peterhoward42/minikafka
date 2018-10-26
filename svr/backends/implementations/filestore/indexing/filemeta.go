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
