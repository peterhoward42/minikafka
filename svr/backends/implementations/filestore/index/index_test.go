package index

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestSerialization tests the serialization methods for the index.
func TestSerialization(t *testing.T) {
	// Create an index programmatically.
	index := makeReferenceIndex()

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
	index := makeReferenceIndex()

	// Should be 1 for unknown topic.
	nextNum := index.NextMessageNumberFor("nosuchtopic")
	expected := int32(1)
	assert.Equal(t, expected, nextNum)

	// Should be 1 for known topic with no message files.
	nextNum = index.NextMessageNumberFor("baz_topic")
	expected = int32(1)
	assert.Equal(t, expected, nextNum)

	// Should be 21 in a prepared case.
	nextNum = index.NextMessageNumberFor("foo_topic")
	expected = int32(21)
	assert.Equal(t, expected, nextNum)
}

//--------------------------------------------------------------------------------
// Auxilliary code.
//--------------------------------------------------------------------------------

func makeReferenceIndex() *Index {

	idx := NewIndex()

	msgNum := int32(1)
	minutesAgo := 1
	for _, topic := range []string{"topicA", "topicB"} {
		msgFileList := idx.RegisterTopic(topic)
		for _, fileName := range []string{"file1", "file2"} {
			msgFileList.RegisterFile(fileName)
			fileMeta := msgFileList.Meta[fileName]

			oldestMeta := fileMeta.Oldest
			oldestMeta.Set(msgNum, nowMinusNMinutes(minutesAgo))
			newestMeta := fileMeta.Newest
			newestMeta.Set(msgNum+5, nowMinusNMinutes(minutesAgo+2))

			msgNum += 10
			minutesAgo += 15
		}
	}
	return idx
}

func nowMinusNMinutes(minutes int) time.Time {
	now := time.Now()
	duration := time.Duration(minutes) * time.Minute
	return now.Add(-duration)
}
