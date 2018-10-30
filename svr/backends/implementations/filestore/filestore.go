// Package filestore provides a message storage system based on a mounted file
// system. It implements the backingstore.contract.BackingStore interface.
package filestore

import (
	"fmt"
	"sync"
	"time"

	minikafka "github.com/peterhoward42/minikafka"
	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/actions"
	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/filenamer"
	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/indexing"
	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/ioutils"
)

var mutex = &sync.Mutex{} // Guards concurrent access of the FileStore.

// FileStore encapsulates the store.
type FileStore struct {
	RootDir string
}

// NewFileStore provides an intialised FileStore object based on the root
// directory provided. It either consumes the file store that is already
// persisted there, or sets up a new one if there isn't one there.
func NewFileStore(rootDir string) (*FileStore, error) {
	// Create the root directory if it does not exist.
	err := ioutils.CreateDirIfDoesntExist(rootDir)
	if err != nil {
		return nil, fmt.Errorf("ioutils.CreateDirIfDoesntExist(): %v", err)
	}
	// Create and persist a blank index file if doesn't exist.
	indexFilePath := filenamer.IndexFile(rootDir)
	if ioutils.Exists(indexFilePath) == false {
		index := indexing.NewIndex()
		err := index.Save(indexFilePath)
		if err != nil {
			return nil, fmt.Errorf("index.Save(): %v", err)
		}
	}
	return &FileStore{rootDir}, nil
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
func (s FileStore) Store(topic string, message minikafka.Message) (
	messageNumber int, err error) {

	mutex.Lock()
	defer mutex.Unlock()

	// Establish the index, - either virgin, or deserialised from disk.
	index := indexing.NewIndex()
	indexPath := filenamer.IndexFile(s.RootDir)
	if ioutils.Exists(indexPath) {
		err = index.PopulateFromDisk(filenamer.IndexFile(s.RootDir))
		if err != nil {
			return -1, fmt.Errorf("index.PopulateFromDisk(): %v", err)
		}
	}

	// Delegate to a StoreAction instance.
	storeAction := actions.StoreAction{
		Topic: topic, Message: message, Index: index, RootDir: s.RootDir}
	messageNumber, _, err = storeAction.Store()

	// Finish up by mandating the index to re-save itself to disk, ready
	// for the next API operation to pick up.
	err = index.Save(filenamer.IndexFile(s.RootDir))
	if err != nil {
		return -1, fmt.Errorf("SaveIndex(): %v", err)
	}

	return messageNumber, nil
}

// RemoveOldMessages is defined by, and documented in the
// backends/contract/BackingStore interface.
func (s FileStore) RemoveOldMessages(maxAge time.Time) error {

	mutex.Lock()
	defer mutex.Unlock()

	// Establish the index, - either virgin, or deserialised from disk.
	index := indexing.NewIndex()
	indexPath := filenamer.IndexFile(s.RootDir)
	if ioutils.Exists(indexPath) {
		err := index.PopulateFromDisk(filenamer.IndexFile(s.RootDir))
		if err != nil {
			return fmt.Errorf("index.PopulateFromDisk(): %v", err)
		}
	}

	// Delegate to a RemoveOldMessagesAction instance.
	rmOldAction := actions.RemoveOldMessagesAction{
		MaxAge: maxAge, Index: index, RootDir: s.RootDir}
	_, _, err := rmOldAction.RemoveOldMessages()

	// Finish up by mandating the index to re-save itself to disk, ready
	// for the next API operation to pick up.
	err = index.Save(filenamer.IndexFile(s.RootDir))
	if err != nil {
		return fmt.Errorf("SaveIndex(): %v", err)
	}

	return nil
}

// Poll is defined by, and documented in the backends/contract/BackingStore
// interface.
func (s FileStore) Poll(topic string, readFrom int) (
	foundMessages []minikafka.Message, newReadFrom int, err error) {

	mutex.Lock()
	defer mutex.Unlock()

	// Establish the index, - either virgin, or deserialised from disk.
	index := indexing.NewIndex()
	indexPath := filenamer.IndexFile(s.RootDir)
	if ioutils.Exists(indexPath) {
		err = index.PopulateFromDisk(filenamer.IndexFile(s.RootDir))
		if err != nil {
			return nil, -1, fmt.Errorf("index.PopulateFromDisk(): %v", err)
		}
	}

	// Delegate to a PollAction instance.
	pollAction := actions.PollAction{
		Topic:    topic,
		ReadFrom: readFrom,
		Index:    index,
		RootDir:  s.RootDir}
	foundMessages, newReadFrom, err = pollAction.Poll()
	if err != nil {
		return nil, -1, fmt.Errorf("possAction.Poll(): %v", err)
	}

	// Finish up by mandating the index to re-save itself to disk, ready
	// for the next API operation to pick up.
	err = index.Save(filenamer.IndexFile(s.RootDir))
	if err != nil {
		return nil, -1, fmt.Errorf("SaveIndex(): %v", err)
	}

	return foundMessages, newReadFrom, nil
}

// ------------------------------------------------------------------------
// Miscellaneous Implementation functions.
// ------------------------------------------------------------------------

func (s FileStore) deleteContents() error {
	err := ioutils.DeleteDirectoryContents(s.RootDir)
	if err != nil {
		return fmt.Errorf("ioutils.DeleteDirectoryContents(): %v", err)
	}
	return nil
}
