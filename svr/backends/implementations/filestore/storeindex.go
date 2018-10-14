package filestore

import (
    "time"
)

//-----------------------------------------------------------------------
// The types from which the index is (hierarchically) composed.
// Described bottom-up.
//-----------------------------------------------------------------------

// Type MsgMeta holds the message number and creation time for stored message.
type MsgMeta struct {
    msgNum int32
    created time.Time
}

// Type FileMeta holds information about the oldest and newest message in
// a file.
type FileMeta struct {
    fileName string
    oldestMsg MsgMeta
    newestMsg  MsgMeta
}

// Type TopicMessages holds information about the set of files that 
// contain a topic's messages.
type TopicMessages []FileMeta // Ordered by file creation sequence.

// Type Index provides an index over the entire store by providing a map
// of topic names to their corresponding TopicMessages object.
type StoreIndex map[string] TopicMessages

//-----------------------------------------------------------------------
// 
//-----------------------------------------------------------------------
