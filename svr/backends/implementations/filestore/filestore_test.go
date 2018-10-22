package filestore

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/peterhoward42/toy-kafka/svr/backends/contract"
	"github.com/stretchr/testify/assert"
)

//--------------------------------------------------------------------------------
// We have unit tests for most of the internals of this package, so
// we can test the bulk of the top level behaviour by using the test suite
// that checks FileStore behaviour against the required behaviour of a
// BackingStore.
//--------------------------------------------------------------------------------

// TestBackingStoreConformance ensures that FileStore passes all the tests
// defined for the BackingStore interface it claims to satisfy.
func TestBackingStoreConformance(t *testing.T) {
	// Prepare a root directory that we can delete after the test.
	rootDir, err := ioutil.TempDir("", "filestore")
	if err != nil {
		msg := fmt.Sprintf("ioutil.TempDir(): %v", err)
		assert.Fail(t, msg)
	}
	defer os.RemoveAll(rootDir)
	filestore := FileStore{rootDir}
	// Delegate to a test suite that takes a contract.BackingStore
	// (interface) argument.
	contract.RunBackingStoreTests(t, filestore)
}
