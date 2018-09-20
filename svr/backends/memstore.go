package backends

import (
	"sync"
)

type messageStream []Message

// These maps are keyed on topic.
type messageStreams map[string]messageStream
type newestMessageNumbers map[string]uint32

// MemStore implements the svr.BackingStore interface using a volatile,
// in-process memory store.
type MemStore struct {
	messageStreams   messageStreams
	newestMsgNumbers newestMessageNumbers
}

var mutex = &sync.Mutex{} // Protect mutation of the MemStore.

// NewMemStore constructs and initializes an empty MemStore instance.
func NewMemStore() *MemStore {
	return &MemStore{messageStreams{}, newestMessageNumbers{}}
}

// Store offers to add a Message to the sequence already held
// for the Topic, and to return the message number thus asigned
// to it.
func (m *MemStore) Store(topic string, message Message) (
	messageNumber uint32, err error) {

	// Protect the memory from concurrent access.
	mutex.Lock()

	// Starting a new topic?
	msgStream, ok := m.messageStreams[topic]
	if ok == false {
		m.messageStreams[topic] = messageStream{}
		m.newestMsgNumbers[topic] = 0
	}

	// Append the message to the stream - coping with the slice being
	// re-allocated.
	msgStream = append(msgStream, message)
	m.messageStreams[topic] = msgStream

	// Increment the newest message number for this topic to become the number
	// to allocate to this message.
	m.newestMsgNumbers[topic]++

	// Take a snapshot of the number, in case a concurrent operation changes
	// it again after we release the mutex, but before we return.
	savedNumber := m.newestMsgNumbers[topic]

	// Release the mutex
	mutex.Unlock()

	return savedNumber, nil
}
