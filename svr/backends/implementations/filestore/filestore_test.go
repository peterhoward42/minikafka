package filestore

import (
	"fmt"
	"log"
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
	log.Printf("hello")
}
