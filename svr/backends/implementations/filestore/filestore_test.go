package filestore

import (
	"os"
	"testing"

	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/ioutils"

	"github.com/peterhoward42/minikafka/svr/backends/contract"
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
	rootDir := ioutils.TmpRootDir(t)
	defer os.RemoveAll(rootDir)
	filestore := FileStore{rootDir}
	// Delegate to a test suite that takes a contract.BackingStore
	// (interface) argument.
	contract.RunBackingStoreTests(t, filestore)
}
