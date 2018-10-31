package filestore

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/ioutils"
	"github.com/stretchr/testify/assert"

	"github.com/peterhoward42/minikafka/svr/backends/contract"
)

//---------------------------------------------------------------------------
// We have unit tests for most of the internals of this package, so
// we can test the bulk of the top level behaviour by using the test suite
// that checks FileStore behaviour against the required behaviour of a
// BackingStore.
//---------------------------------------------------------------------------

// TestBackingStoreConformance ensures that FileStore passes all the tests
// defined for the BackingStore interface it claims to satisfy.
func TestBackingStoreConformance(t *testing.T) {
	// Prepare a root directory that we can delete after the test.
	rootDir := ioutils.TmpRootDir(t)
	defer os.RemoveAll(rootDir)
	filestore, err := NewFileStore(rootDir)
	if err != nil {
		msg := fmt.Sprintf("NewFileStore(): %v", err)
		assert.Fail(t, msg)
	}
	// Delegate to a test suite that takes a contract.BackingStore
	// (interface) argument.
	contract.RunBackingStoreTests(t, filestore)
}

//---------------------------------------------------------------------------
// Some additional tests as the FileStore API level - testing behaviour
// that is not covered by the BackingStore suite test suite above.
//---------------------------------------------------------------------------

func TestConstructionOnADirectoryThatDoesntExist(t *testing.T) {
	// This makes sure that the root directory is created when it doesn't
	// already exist.

	// Use our utility to make a root directory.
	rootDir := ioutils.TmpRootDir(t)
	defer os.RemoveAll(rootDir)
	// But construct the file store on a non-existent directory inside it.
	nonExistent := path.Join(rootDir, "doesnotexist")
	filestore, err := NewFileStore(nonExistent)
	if err != nil {
		msg := fmt.Sprintf("NewFileStore(): %v", err)
		assert.Fail(t, msg)
	}
	// Make sure we can store something in it without error.
	_, err = filestore.Store("some_topic", []byte("a message"))
	if err != nil {
		msg := fmt.Sprintf("filestore.Store(): %v", err)
		assert.Fail(t, msg)
	}
}

func TestPersistence(t *testing.T) {
	// This test makes sure that if we store some messages in one
	// FileStore instance, then when we create a new instance based on
	// the same root directory - it picks up the stored message and index
	// left behind by the first instance as it should.

	rootDir := ioutils.TmpRootDir(t)
	defer os.RemoveAll(rootDir)

	// Create the first file store instance and store something in it.
	filestore, err := NewFileStore(rootDir)
	if err != nil {
		msg := fmt.Sprintf("NewFileStore(): %v", err)
		assert.Fail(t, msg)
	}
	topic := "some topic"
	msgNumber, err := filestore.Store(topic, []byte("a message"))
	if err != nil {
		msg := fmt.Sprintf("filestore.Store(): %v", err)
		assert.Fail(t, msg)
	}
	assert.Equal(t, 1, msgNumber)

	// Create a second file store over the same root directory, store something
	// in it, and make sure a Poll returns both messages.
	newFileStore, err := NewFileStore(rootDir)
	if err != nil {
		msg := fmt.Sprintf("NewFileStore(): %v", err)
		assert.Fail(t, msg)
	}
	msgNumber, err = newFileStore.Store(topic, []byte("a message"))
	if err != nil {
		msg := fmt.Sprintf("filestore.Store(): %v", err)
		assert.Fail(t, msg)
	}
	assert.Equal(t, 2, msgNumber)

	readFrom := 1
	messages, newReadFrom, err := newFileStore.Poll(topic, readFrom)
	if err != nil {
		msg := fmt.Sprintf("newFileStore.Poll(): %v", err)
		assert.Fail(t, msg)
	}
	assert.Equal(t, 2, len(messages))
	assert.Equal(t, 3, newReadFrom)
}
