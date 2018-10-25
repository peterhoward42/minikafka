package indexing

import (
	"sort"
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
	// In addition to the map, we must track the order in which names
	// are introduced to the list, so that we can identify the most recent.
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

// SpentFiles provides a list of filenames from the list, for which the
// constituent messages are all older than the time specified.
func (lst *MessageFileList) SpentFiles(maxAge time.Time) []string {
	spent := []string{}
	for name, fileMeta := range lst.Meta {
		if fileMeta.Newest.Created.Before(maxAge) {
			spent = append(spent, name)
		}
	}
	return spent
}

// ForgetFiles mandates the MessageFileList to forget about the given
// set of file names.
func (lst *MessageFileList) ForgetFiles(names []string) {
	for _, name := range names {
		// Get rid of this name from the map of file names to FileMeta.
		delete(lst.Meta, name)
		// Take the name out of the ordered list of filenames also.
		n := len(lst.Names)
		index := sort.SearchStrings(lst.Names, name)
		if index == n { // Not present.
			continue
		}
		if index == n-1 { // Special case when it's the last element.
			lst.Names = lst.Names[0 : n-1]
			continue
		}
		// General case.
		newList := []string{}
		newList = append(newList, lst.Names[0:index]...)
		newList = append(newList, lst.Names[index+1:]...)
		lst.Names = newList
	}
}

// NumMessagesInFile provides a count of how many messages are held
// in the given file.
func (lst *MessageFileList) NumMessagesInFile(name string) int {
	fileMeta, ok := lst.Meta[name]
	// Unknown name?
	if ok == false {
		return 0
	}
	// Name known, but no messages recorded for it yet.
	if fileMeta.Oldest.MsgNum == 0 {
		return 0
	}
	// General case.
	return int(fileMeta.Newest.MsgNum) - int(fileMeta.Oldest.MsgNum) + 1
}
