package filestore

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/peterhoward42/toy-kafka/svr/backends/contract"
)

//--------------------------------------------------------------------------------
// Delegate to the BackingStore interface test suite.
//--------------------------------------------------------------------------------

// TestBackingStoreConformance ensures that FileStore passes all the tests
// defined for the BackingStore interface it claims to satisfy.
func TestBackingStoreConformance(t *testing.T) {
	filestore, err := NewFileStore("/tmp/store")
	if err != nil {
		t.Fatalf("NewFileStore(): %v", err)
	}
	// Delegate to a test suite that takes a contract.BackingStore
	// (interface) argument.
	contract.RunBackingStoreTests(t, *filestore)
}

//--------------------------------------------------------------------------------
// Auxilliary code.
//--------------------------------------------------------------------------------

// makeTempDirOrExit makes a temporary directory and returns its name.
// It responds to errors with log.Fatalf().
func makeTempDirOrFatal() string {
	dir, err := ioutil.TempDir("", "file_store")
	if err != nil {
		log.Fatalf("ioutil.TempDir(): %v", err)
	}
	return dir
}

func clearTempDirOrFatal(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		log.Fatalf("os.RemoveAll(): %v", err)
	}
}
