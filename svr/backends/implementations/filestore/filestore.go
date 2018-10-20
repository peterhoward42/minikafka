// Package filestore provides a message storage system based on a mounted file
// system. It implements the backingstore.contract.BackingStore interface.
package filestore

import (
	"fmt"
	"sync"
	"time"

	toykafka "github.com/peterhoward42/toy-kafka"
	"github.com/peterhoward42/toy-kafka/svr/backends/implementations/filestore/ioutils"
	"github.com/peterhoward42/toy-kafka/svr/backends/implementations/filestore/indexing"
	"github.com/peterhoward42/toy-kafka/svr/backends/implementations/filestore/filenamer"
	"github.com/peterhoward42/toy-kafka/svr/backends/implementations/filestore/actions"

)

var mutex = &sync.Mutex{} // Guards concurrent access of the FileStore.

// FileStore encapsulates the store.
type FileStore struct {
	rootDir string
}

// ------------------------------------------------------------------------
// METHODS TO SATISFY THE BackingStore INTERFACE.
//
// Most of these methods delegate to a helper function, but wrap it the
// call in a mutex().
// ------------------------------------------------------------------------

// DeleteContents removes all contents from the store.
func (s FileStore) DeleteContents() error {
	mutex.Lock()
	defer mutex.Unlock()
	return s.deleteContents()
}

// Store is defined by, and documented in the backends/contract/BackingStore
// interface.
func (s FileStore) Store(topic string, message toykafka.Message) (
	messageNumber int, err error) {

	mutex.Lock()
	defer mutex.Unlock()

    // Acquire the index from disk.
    index := indexing.NewIndex()
    err = index.PopulateFromDisk(filenamer.IndexFile(s.rootDir))
	if err != nil {
		return -1, fmt.Errorf("index.PopulateFromDisk(): %v", err)
	}

    // Delegate to a StoreAction instance.
    storeAction := actions.StoreAction{index, topic, message}
    messageNumber, err := storeAction.Store()

    // Finish up by mandating the index to re-save itself to disk, ready
    // for the next API operation to pick up.
	err = index.Save(filenamer.IndexFile(s.rootDir))
	if err != nil {
		return -1, fmt.Errorf("SaveIndex(): %v", err)
	}

	return messageNumber, nil
}

// RemoveOldMessages is defined by, and documented in the
// backends/contract/BackingStore interface.
func (s FileStore) RemoveOldMessages(maxAge time.Time) (
	nRemoved int, err error) {
	return -1, nil
}

// Poll is defined by, and documented in the backends/contract/BackingStore
// interface.
func (s FileStore) Poll(topic string, readFrom int) (
	foundMessages []toykafka.Message, newReadFrom int, err error) {

	foundMessages = []toykafka.Message{}
	return foundMessages, 11, nil
}

// ------------------------------------------------------------------------
// Miscellaneous Implementation functions.
// ------------------------------------------------------------------------

func (s FileStore) deleteContents() error {
	err := ioutils.DeleteDirectoryContents(s.rootDir)
	if err != nil {
		return fmt.Errorf("ioutils.DeleteDirectoryContents(): %v", err)
	}
	return nil
}


