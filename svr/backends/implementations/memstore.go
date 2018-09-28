package implementations

import (
	"log"
	"sort"
	"sync"
	"time"

	toykafka "github.com/peterhoward42/toy-kafka"
)

// MemStore implements the svr/backends/contract/BackingStore interface using
// a volatile, in-process memory store.
type MemStore struct {
	// Fundamental storage is separated by topic, and comprises simply
	// time-ordered slices of messages held in *storedMessage* objects.
	// (These hold the message itself, plus its creation time, and
	// message-number.)
	messagesPerTopic    map[string][]storedMessage // Keyed on topic.
	newestMessageNumber map[string]int             // Keyed on topic.
}

// NewMemStore constructs and initializes an empty MemStore instance.
func NewMemStore() *MemStore {
	return &MemStore{
		messagesPerTopic:    map[string][]storedMessage{},
		newestMessageNumber: map[string]int{},
	}
}

var mutex = &sync.Mutex{} // Protects concurrent access of the MemStore.

// ------------------------------------------------------------------------
// METHODS TO SATISFY THE BackingStore INTERFACE.
// ------------------------------------------------------------------------

// Store is defined by, and documented in the backends/contract/BackingStore
// interface.
func (m *MemStore) Store(topic string, message toykafka.Message) (
	messageNumber int, err error) {

	mutex.Lock()
	defer mutex.Unlock()

	// Special case, if this is a new topic.
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
func (m *MemStore) RemoveOldMessages(maxAge time.Time) error {
	mutex.Lock()
	defer mutex.Unlock()
	for topic := range m.messagesPerTopic {
		err := m.removeOldMessagesFromTopic(topic, maxAge)
		if err != nil {
			return err
		}
	}
	return nil
}

// Poll is defined by, and documented in the backends/contract/BackingStore
// interface.
func (m *MemStore) Poll(topic string, readFrom int) (
	foundMessages []toykafka.Message, newReadFrom int, err error) {

	mutex.Lock()
	defer mutex.Unlock()

	storedMessages := m.messagesPerTopic[topic]
	serveFrom := sort.Search(len(storedMessages), func(i int) bool {
		return storedMessages[i].messageNumber >= readFrom
	})
	foundMessages = []toykafka.Message{}
	for _, msg := range storedMessages[serveFrom:] {
		foundMessages = append(foundMessages, msg.payload)
	}
	nServed := len(foundMessages)
	newReadFrom = readFrom + nServed
	return
}

// ------------------------------------------------------------------------
// Helper functions.
// ------------------------------------------------------------------------

// RemoveOldMessagesFromTopic is a topic-specific helper function for the
// whole-store RemoveOldMessages method.
func (m *MemStore) removeOldMessagesFromTopic(
	topic string, maxAge time.Time) error {
	messages := m.messagesPerTopic[topic]
	keepFrom := sort.Search(len(messages), func(i int) bool {
		return messages[i].creationTime.After(maxAge)
	})
	messagesToKeep := messages[keepFrom:] // Safe also when keeping none.

	// If we're keeping from index 3, then we're going to remove 0,1,2.
	nRemoving := keepFrom
	if nRemoving > 0 {
		log.Printf("Removing %v old messages from topic: %v", nRemoving, topic)
		// We make and use a copy of the messagesToKeep slice, to free up the
		// old backing array for garbage collection. Otherwise it would grow
		// inexorably.
		freshSlice := make([]storedMessage, len(messagesToKeep))
		copy(freshSlice, messagesToKeep)
		m.messagesPerTopic[topic] = freshSlice
	}
	return nil
}

// ------------------------------------------------------------------------
// AUXILLIARY TYPES AND THEIR METHODS.
// ------------------------------------------------------------------------

// storedMessage is a private type for the MemStore implementation backing store
// implementation which encapsulates a message itself, along with its creation
// time, and message number.
type storedMessage struct {
	payload       toykafka.Message
	creationTime  time.Time
	messageNumber int
}
