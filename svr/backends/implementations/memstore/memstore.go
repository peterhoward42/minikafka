package memstore

import (
	"fmt"
	"sort"
	"sync"
	"time"

	minikafka "github.com/peterhoward42/minikafka"
)

var mutex = &sync.Mutex{} // Guards concurrent access of the MemStore.

// MemStore implements the svr/backends/contract/BackingStore interface using
// a volatile, in-process memory store. It exists principally to aid
// development and testing without being dependent on real storage.
type MemStore struct {
	// Fundamental storage is separated by topic, and comprises simply
	// time-ordered slices of messages held in *storedMessage* objects.
	// (These hold the message itself, plus its creation time, and
	// message-number.)
	messagesPerTopic    map[string][]storedMessage // Keyed on topic.
	newestMessageNumber map[string]int             // Keyed on topic.
}

// NewMemStore instantiates, initializes and returns a MemStore.
func NewMemStore() *MemStore {
	return &MemStore{
		messagesPerTopic:    map[string][]storedMessage{},
		newestMessageNumber: map[string]int{},
	}
}

// ------------------------------------------------------------------------
// METHODS TO SATISFY THE BackingStore INTERFACE.
// ------------------------------------------------------------------------

// DeleteContents is defined in the BackingStore interface.
func (m MemStore) DeleteContents() error {
	mutex.Lock()
	defer mutex.Unlock()
	for k := range m.messagesPerTopic {
		delete(m.messagesPerTopic, k)
	}
	for k := range m.newestMessageNumber {
		delete(m.newestMessageNumber, k)
	}
	return nil
}

// Store is defined by, and documented in the backends/contract/BackingStore
// interface.
func (m MemStore) Store(topic string, message minikafka.Message) (
	messageNumber int, err error) {

	mutex.Lock()
	defer mutex.Unlock()

	// Bit of extra work if this is a new topic.
	if _, ok := m.messagesPerTopic[topic]; ok == false {
		m.messagesPerTopic[topic] = []storedMessage{}
		m.newestMessageNumber[topic] = 0
	}

	// Drop into the general case.

	// Allocate the next available message number.
	m.newestMessageNumber[topic]++

	// Make and add the new message.
	msgToAdd := storedMessage{message, time.Now(),
		m.newestMessageNumber[topic]}
	m.messagesPerTopic[topic] = append(m.messagesPerTopic[topic], msgToAdd)

	return m.newestMessageNumber[topic], nil
}

// RemoveOldMessages is defined by, and documented in the
// backends/contract/BackingStore interface.
func (m MemStore) RemoveOldMessages(maxAge time.Time) (
	nRemoved int, err error) {
	mutex.Lock()
	defer mutex.Unlock()
	for topic := range m.messagesPerTopic {
		n, err := m.removeOldMessagesFromTopic(topic, maxAge)
		if err != nil {
			return 0, err
		}
		nRemoved += n
	}
	return nRemoved, nil
}

// Poll is defined by, and documented in the backends/contract/BackingStore
// interface.
func (m MemStore) Poll(topic string, readFrom int) (
	foundMessages []minikafka.Message, newReadFrom int, err error) {

	mutex.Lock()
	defer mutex.Unlock()

	storedMessages, ok := m.messagesPerTopic[topic]
	if !ok {
		return nil, -1, fmt.Errorf("No such topic: %s", topic)
	}
	serveFromIndex := sort.Search(len(storedMessages), func(i int) bool {
		return storedMessages[i].messageNumber >= readFrom
	})

	foundMessages = []minikafka.Message{}
	var highest int
	for _, msg := range storedMessages[serveFromIndex:] {
		foundMessages = append(foundMessages, msg.message)
		highest = msg.messageNumber
	}
	nFound := len(foundMessages)
	if nFound > 0 {
		newReadFrom = highest + 1
		return foundMessages, newReadFrom, nil
	}
	unchangedReadFrom := readFrom
	return foundMessages, unchangedReadFrom, nil
}

// ------------------------------------------------------------------------
// Helper functions.
// ------------------------------------------------------------------------

// RemoveOldMessagesFromTopic is a topic-specific helper function for the
// whole-store RemoveOldMessages method.
func (m MemStore) removeOldMessagesFromTopic(
	topic string, maxAge time.Time) (nRemoved int, err error) {

	// Find the boundary between the messages to keep and those to remove.
	messages := m.messagesPerTopic[topic]
	keepFromIndex := sort.Search(len(messages), func(i int) bool {
		return messages[i].creationTime.After(maxAge)
	})

	messagesToKeep := messages[keepFromIndex:] // Safe when keeping none.
	nRemoved = keepFromIndex

	if nRemoved == 0 {
		return 0, nil
	}

	// Replace the incumbent queue slice with a newly minted one so that the
	// underlying array gets freed for garbage collection. Otherwise it would
	// grown inexorably.
	freshSlice := make([]storedMessage, len(messagesToKeep))
	copy(freshSlice, messagesToKeep)
	m.messagesPerTopic[topic] = freshSlice

	return nRemoved, nil
}

// ------------------------------------------------------------------------
// AUXILLIARY CODE
// ------------------------------------------------------------------------

// storedMessage is a private type for the MemStore backing store
// implementation which encapsulates a message itself, along with its creation
// time, and message number.

type storedMessage struct {
	message       minikafka.Message
	creationTime  time.Time
	messageNumber int
}
