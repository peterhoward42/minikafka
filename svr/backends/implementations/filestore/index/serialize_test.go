package index

import (
	"bytes"
    "testing"
	"encoding/json"
)

// TestSerialization tests the serialization methods for the index.
func TestSerialization(t *testing.T) {
	// Create an index programmatically.
	index := makeReferenceIndexForTesting()

	// Serialize it into a buffer.
	var buf bytes.Buffer
	err := index.Encode(&buf)
	if err != nil {
		t.Fatalf("index.Encode: %v", err)
	}

	// Make a new index by deserializing one from the buffer.
	restored := NewIndex()
	err = restored.Decode(&buf)
	if err != nil {
		t.Fatalf("index.Decode: %v", err)
	}

	// The restored and original index objects can be compared for
	// equality via their conversion to json.

	// This could break in the future due to how gob serialization will
	// serialize an empty slice as simply a nil value, or because iteration
	// over maps does not guarantee repeatable ordering over keys. TODO.
	origJSON, err := json.Marshal(index)
	if err != nil {
		t.Fatalf("json.Marshal(): %v", err)
	}
	restoredJSON, err := json.Marshal(restored)
	if err != nil {
		t.Fatalf("json.Marshal(): %v", err)
	}
	sOrig := string(origJSON)
	sRestored := string(restoredJSON)
	if sOrig != sRestored {
		t.Fatalf("Restored index differs from the one saved.")
	}
}
