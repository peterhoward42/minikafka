package memstore

import (
	"testing"

	"github.com/peterhoward42/minikafka/svr/backends/contract"
)

// TestMemStore ensures that implementations.MemStore passes all the tests
// defined for the BackingStore interface it claims to satisfy. It delegates
// the real work to an external BackingStore interface test suite.
func TestMemStore(t *testing.T) {
	memstore := NewMemStore()
	// Delegate to a test suite that takes a contract.BackingStore
	// (interface) argument.
	contract.RunBackingStoreTests(t, *memstore)
}
