package implementations

import (
	"sync"
	"time"

	"github.com/peterhoward42/toy-kafka/svr/backends/contract"
)

// MemStore implements the svr/backends/contract/BackingStore interface using
// a volatile, in-process memory store.
type MemStore struct {
	// Fundamental storage is separated by topic, and comprises simply an
	// time-ordered slice of *messageStorage* objects. (These hold
	// their own payload, creation time, and message-number.
	allTopics map[string][]messageStorage // Keyed on topic.
}

// NewMemStore constructs and initializes an empty MemStore instance.
func NewMemStore() *MemStore {
	return &MemStore{map[string][]messageStorage{}}
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

	// Special case, this is a new topic.
	_, ok := m.allTopics[topic]
	if ok == false {
		m.allTopics[topic] = []messageStorage{}
	}

	// Now the general case.

	// Allocate the next available message number.
	count := len(m.allTopics[topic])
	var messageNum int
	if count == 0 {
		messageNum = 1
	} else {
		messageNum = m.allTopics[topic][count-1].messageNumber + 1
	}

	// Make and add the new message.
	msgToAdd := messageStorage{message, time.Now(), messageNum}
	m.allTopics[topic] = append(m.allTopics[topic], msgToAdd)

	mutex.Unlock()
	return messageNum, nil
}

// RemoveOldMessages is defined by, and documented in the
// backends/contract/BackingStore interface.
func (m *MemStore) RemoveOldMessages(maxAge time.Time) {
	// Protect the memory from concurrent access.
	mutex.Lock()
	for _, tpc := range m.allTopics {
	}
	mutex.Unlock()
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
