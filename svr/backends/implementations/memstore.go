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
func (m *MemStore) RemoveOldMessages(maxAge time.Time) (
	removed map[string][]int, err error) {
	mutex.Lock()
	defer mutex.Unlock()
	removed = map[string][]int{}
	for topic := range m.messagesPerTopic {
		err = m.removeOldMessagesFromTopic(topic, maxAge, removed)
		if err != nil {
			return
		}
	}
	return
}

// Poll is defined by, and documented in the backends/contract/BackingStore
// interface.
func (m *MemStore) Poll(topic string, readFrom int) (
	foundMessages []toykafka.Message, newReadFrom int, err error) {

	mutex.Lock()
	defer mutex.Unlock()

	storedMessages := m.messagesPerTopic[topic]
	serveFromIndex := sort.Search(len(storedMessages), func(i int) bool {
		return storedMessages[i].messageNumber >= readFrom
	})

	foundMessages = []toykafka.Message{}
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
func (m *MemStore) removeOldMessagesFromTopic(
	topic string, maxAge time.Time, removed map[string][]int) error {

	// Find the boundary between the messages to keep and those to remove.
	messages := m.messagesPerTopic[topic]
	removed[topic] = []int{}
	keepFromIndex := sort.Search(len(messages), func(i int) bool {
		return messages[i].creationTime.After(maxAge)
	})

	messagesToKeep := messages[keepFromIndex:] // Safe when keeping none.
	numberRemoving := keepFromIndex

	if numberRemoving == 0 {
		return nil
	}

	log.Printf("Removing %v old messages from topic: %v", numberRemoving, topic)

	// Harvest the message number of those that are being removed.
	for _, msg := range messages[:keepFromIndex] {
		removed[topic] = append(removed[topic], msg.messageNumber)
	}

	// Replace the incumbent queue slice with a newly minted one so that the
	// underlying array gets freed for garbage collection. Otherwise it would
	// grown inexorably.
	freshSlice := make([]storedMessage, len(messagesToKeep))
	copy(freshSlice, messagesToKeep)
	m.messagesPerTopic[topic] = freshSlice

	return nil
}

// ------------------------------------------------------------------------
// AUXILLIARY TYPES AND THEIR METHODS.
// ------------------------------------------------------------------------

// storedMessage is a private type for the MemStore implementation backing store
// implementation which encapsulates a message itself, along with its creation
// time, and message number.
type storedMessage struct {
	message       toykafka.Message
	creationTime  time.Time
	messageNumber int
}
