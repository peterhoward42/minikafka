package index

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//--------------------------------------------------------------------------
// API
//--------------------------------------------------------------------------

// MakeReferenceIndexForTesting provides a useful, repeatable Index for
// testing purposes.
func MakeReferenceIndexForTesting() *Index {

	idx := NewIndex()

	msgNum := int32(0)
	minutes := 1
	for _, topic := range []string{"topicA", "topicB"} {
		msgFileList := idx.GetMessageFileListFor(topic)
		for _, fileName := range []string{"file1", "file2"} {
			msgFileList.RegisterNewFile(fileName)

			fileMeta := msgFileList.Meta[fileName]

			fileMeta.Oldest.MsgNum = msgNum + 1
			fileMeta.Oldest.Created = nowMinusNMinutes(minutes)

			fileMeta.Newest.MsgNum = msgNum + 5
			fileMeta.Newest.Created = nowMinusNMinutes(minutes + 5)

			msgNum += 10
			minutes += 15
		}
	}
	return idx
}

//--------------------------------------------------------------------------
// API
//--------------------------------------------------------------------------

// TestSerialization tests the serialization methods for the index.
func TestSerialization(t *testing.T) {
	// Create an index programmatically.
	index := MakeReferenceIndexForTesting()

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

func TestNextMsgNumForTopic(t *testing.T) {
	index := MakeReferenceIndexForTesting()

	// Should be 1 for unknown topic.
	nextNum := index.NextMessageNumberFor("nosuchtopic")
	expected := int32(1)
	assert.Equal(t, expected, nextNum)

	// Should be 36 in a prepared case.
	nextNum = index.NextMessageNumberFor("topicB")
	expected = int32(36)
	assert.Equal(t, expected, nextNum)
}

func TestCurrentMsgFileNameFor(t *testing.T) {
	index := MakeReferenceIndexForTesting()
    // Check correct when topic is known and has files registered.
    currentName := index.CurrentMsgFileNameFor("topicA")
    expected := "file2"
    assert.Equal(t, expected, currentName)
    // Check correct when topic is unknown.
    currentName = index.CurrentMsgFileNameFor("nosuchtopic")
    expected = ""
    assert.Equal(t, expected, currentName)
}

func TestHasNameBeenUsedForTopic(t *testing.T) {
	index := MakeReferenceIndexForTesting()
    // When should say yes.
    used := index.HasNameBeenUsedForTopic("file1", "topicA")
    expected := true
    assert.Equal(t, expected, used)
    // When should say no because names exist but not this one. 
    used = index.HasNameBeenUsedForTopic("unknownname", "topicA")
    expected = false
    assert.Equal(t, expected, used)
    // When should say no because no names exist.
    used = index.HasNameBeenUsedForTopic("file1", "unknowntopic")
    expected = false
    assert.Equal(t, expected, used)
}

//--------------------------------------------------------------------------------
// Auxilliary code.
//--------------------------------------------------------------------------------


func nowMinusNMinutes(minutes int) time.Time {
	now := time.Now()
	duration := time.Duration(minutes) * time.Minute
	return now.Add(-duration)
}
