package filestore

import (
	"sync"
	"time"

	toykafka "github.com/peterhoward42/toy-kafka"
)

var mutex = &sync.Mutex{} // Guards concurrent access of the FileStore.

// FileStore implements the svr/backends/contract/BackingStore interface using
// files on disk.
type FileStore struct {
	rootDir string
}

// NewFileStore instantiates, initializes and returns a FileStore.
func NewFileStore(rootDir string) *FileStore {
	return &FileStore{
		rootDir: rootDir,
	}
}

// ------------------------------------------------------------------------
// METHODS TO SATISFY THE BackingStore INTERFACE.
// ------------------------------------------------------------------------

// DeleteContents removes all contents from the store.
func (m FileStore) DeleteContents() {
}

// Store is defined by, and documented in the backends/contract/BackingStore
// interface.
func (m FileStore) Store(topic string, message toykafka.Message) (
	messageNumber int, err error) {

	mutex.Lock()
	defer mutex.Unlock()

	/*
	   Psuedo code

	   Identify directory for topic, making one if necessary.

	   Read the index file
	   Serialise the message ready for storing (requires msg# from index)
	   Measure size of newest file. (or observe there are none)
	   Decide if going to use existing, open new, or open inaugral

	   capture file to add to based on existence of any and if latest is
	   big enough

	   if need to create one
	       do so
	       update index as to its existence

	   append message to chosen file
	   update index
	*/
	return -1, nil
}

// RemoveOldMessages is defined by, and documented in the
// backends/contract/BackingStore interface.
func (m FileStore) RemoveOldMessages(maxAge time.Time) (
	nRemoved int, err error) {
	return -1, nil
}

// Poll is defined by, and documented in the backends/contract/BackingStore
// interface.
func (m FileStore) Poll(topic string, readFrom int) (
	foundMessages []toykafka.Message, newReadFrom int, err error) {

	foundMessages = []toykafka.Message{}
	return foundMessages, 11, nil
}

// ------------------------------------------------------------------------
// Helper functions.
// ------------------------------------------------------------------------
