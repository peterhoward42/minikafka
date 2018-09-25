package implementations

import (
	"log"
	"sort"
	"sync"
	"time"

	"github.com/peterhoward42/toy-kafka/svr/backends/contract"
)

// MemStore implements the svr/backends/contract/BackingStore interface using
// a volatile, in-process memory store.
type MemStore struct {
	// Fundamental storage is separated by topic, and comprises simply
	// time-ordered slices of messages held in *messageStorage* objects.
	// (These hold their own message payload, plus their creation time, and
	// message-number.
	messagesPerTopic    map[string][]messageStorage // Keyed on topic.
	newestMessageNumber map[string]int              // Keyed on topic.
}

// NewMemStore constructs and initializes an empty MemStore instance.
func NewMemStore() *MemStore {
	return &MemStore{
		messagesPerTopic:    map[string][]messageStorage{},
		newestMessageNumber: map[string]int{},
	}
}

var mutex = &sync.Mutex{} // Protect mutation of the MemStore.

// ------------------------------------------------------------------------
// METHODS TO SATISFY THE BackingStore INTERFACE.
// ------------------------------------------------------------------------

// Store is defined by, and documented in the backends/contract/BackingStore
// interface.
func (m *MemStore) Store(topic string, message contract.Message) (
	messageNumber int, err error) {

	mutex.Lock()

	// Special case, if this is a new topic.
	if _, ok := m.messagesPerTopic[topic]; ok == false {
		m.messagesPerTopic[topic] = []messageStorage{}
		m.newestMessageNumber[topic] = 0
	}

	// Drop into the general case.

	// Allocate the next available message number.
	m.newestMessageNumber[topic]++

	// Make and add the new message.
	msgToAdd := messageStorage{message, time.Now(), m.newestMessageNumber[topic]}
	m.messagesPerTopic[topic] = append(m.messagesPerTopic[topic], msgToAdd)

	mutex.Unlock()
	return m.newestMessageNumber[topic], nil
}

// RemoveOldMessages is defined by, and documented in the
// backends/contract/BackingStore interface.
func (m *MemStore) RemoveOldMessages(maxAge time.Time) {
	mutex.Lock()
	for topic := range m.messagesPerTopic {
		m.removeOldMessagesFromTopic(topic, maxAge)
	}
	mutex.Unlock()
}

// RemoveOldMessagesFromTopic is a topic-specific helper function for the
// whole-store RemoveOldMessages method.
func (m *MemStore) removeOldMessagesFromTopic(topic string, maxAge time.Time) {
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
		freshSlice := make([]messageStorage, len(messagesToKeep))
		copy(freshSlice, messagesToKeep)
		m.messagesPerTopic[topic] = freshSlice
	}
}

// ------------------------------------------------------------------------
// AUXILLIARY TYPES AND THEIR METHODS.
// ------------------------------------------------------------------------

// messageStorage encapsulates a message payload (bytes), plus its creation time,
// and message number.
type messageStorage struct {
	payload       contract.Message
	creationTime  time.Time
	messageNumber int
}
