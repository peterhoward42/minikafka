package filestore

import (
	"testing"
    "io/ioutil"
    "log"
    "os"

	"github.com/peterhoward42/toy-kafka/svr/backends/contract"
)

// First we test the error handling of the concrete implementation, before
// delegating the remainder of the testing to a suite that checks that the
// implementation's behaviour does what is required from the BackingStore interface.



// TestIndexInitialisation makes sure that when you call any of the BackingStore 
// API methods, the common preliminary step of building the index by 
// deserializing it from disk reports errors as it should.
func TestIndexInitialisation(t *testing.T) {
    rootDir := makeTempDirOrFatal()
    defer clearTempDirOrFatal(rootDir)
    store := NewFileStore(rootDir)

    _, err := store.Store("fibble_topic", []byte("a message"))
    if err == nil {
        t.Fatalf("store.Store() should have produced error")
    }
}

//--------------------------------------------------------------------------------
// Now delegate to the BackingStore interface test suite.
//--------------------------------------------------------------------------------

// TestBackingStoreConformance ensures that this passes all the tests
// defined for the BackingStore interface it claims to satisfy. 
func _TestBackingStoreConformance(t *testing.T) {
	filestore := NewFileStore("/tmp/store")
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
