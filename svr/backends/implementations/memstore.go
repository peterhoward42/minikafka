package implementations

import (
	"sort"
	"sync"
	"time"

	"github.com/peterhoward42/toy-kafka/svr/backends/contract"
)

// MemStore implements the svr/backends/contract/BackingStore interface using
// a volatile, in-process memory store.
type MemStore struct {
	allTopics map[string]*oneTopic // Keyed on topic.
}

// NewMemStore constructs and initializes an empty MemStore instance.
func NewMemStore() *MemStore {
	return &MemStore{map[string]*oneTopic{}}
}

var mutex = &sync.Mutex{} // Protect mutation of the MemStore.

// METHODS TO SATISFY THE BackingStore INTERFACE.

// Store is defined by, and documented in the backends/contract/BackingStore
// interface.
func (m *MemStore) Store(topic string, message contract.Message) (
	messageNumber uint32, err error) {

	mutex.Lock()

	// New topic?
	tpc, ok := m.allTopics[topic]
	if ok == false {
		tpc := newOneTopic()
		m.allTopics[topic] = tpc
	}

	tpc.messages = append(tpc.messages, message)
	tpc.newestMessageNumber++
	tpc.created[tpc.newestMessageNumber] = time.Now()

	n := tpc.newestMessageNumber // Copy avoids mutation risk before return.
	mutex.Unlock()
	return n, nil
}

// RemoveOldMessages is defined by, and documented in the
// backends/contract/BackingStore interface.
func (m *MemStore) RemoveOldMessages(maxAge time.Time) {
	// Protect the memory from concurrent access.
	mutex.Lock()

	earliest := sort.Search(fibble)

	// What datastructures to releate messages to creation time?
	// a map[topic][msgnumber] = creation Time obj.
	// When set? time of addition.
	// How to do reduction in one go? = use sort.Search to find slice indices.
	// Coping with none such. await
	// What side effects pruning needed? = remove spent keys in created

	// Release the mutex
	mutex.Unlock()
}

// ------------------------------------------------------------------------
// AUXILLIARY TYPES AND THEIR METHODS.
// ------------------------------------------------------------------------

type messageStream []contract.Message

type oneTopic struct {
	messages            messageStream
	newestMessageNumber uint32
	created             map[uint32]time.Time
}

func newOneTopic() *oneTopic {
	return &oneTopic{
		messages: []contract.Message{},
	}
}
