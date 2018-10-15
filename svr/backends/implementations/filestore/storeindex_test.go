package filestore

import (
	"encoding/json"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestSaveAndRestoreIndex creates an index programmatically and then
// mandates the index to save itself to disk. It then attempts to make a new
// index by de-serializing from that disk file, and checks the restored index
// has the expected contents.
func TestSaveAndRestoreIndex(t *testing.T) {
	index := makeIndexProgrammatically()
	fileName := path.Join(os.TempDir(), "saveandrestoreindex")
	err := index.Save(fileName)
	if err != nil {
		t.Fatalf("index.Save: %v", err)
	}
	restoredIndex, err := LoadStoreIndex(fileName)
	if err != nil {
		t.Fatalf("NewStoreIndex: %v", err)
	}
	if restoredIndex == nil {
		t.Fatalf("Restored index is nil.")
	}
	// The restored and original index objects can be compared for
	// equality via their conversion to json.
	origJSON, err := json.Marshal(index)
	if err != nil {
		t.Fatalf("json.Marshal(): %v", err)
	}
	restoredJSON, err := json.Marshal(restoredIndex)
	if err != nil {
		t.Fatalf("json.Marshal(): %v", err)
	}
	if string(origJSON) != string(restoredJSON) {
		t.Fatalf("Restored index differs from the one saved.")
	}
}

func TestNextMsgNumForTopic(t *testing.T) {
	index := makeIndexProgrammatically()

	// Should be 1 for virgin topic.
	nextNum := index.nextMessageNumberFor("nosuchtopic")
	assert.Equal(t, 1, nextNum)

	// Should be 21 in a prepared case.
	nextNum = index.nextMessageNumberFor("foo_topic")
	assert.Equal(t, 21, nextNum)
}

//--------------------------------------------------------------------------------
// Auxilliary code.
//--------------------------------------------------------------------------------

func makeIndexProgrammatically() *StoreIndex {

	idx := StoreIndex{}

	idx["foo_topic"] = []FileMeta{}
	foo1Meta := FileMeta{
		"foo1",
		MsgMeta{1, time.Now().Add(-9 * 24 * time.Hour)},
		MsgMeta{10, time.Now().Add(-8 * 24 * time.Hour)},
	}
	foo2Meta := FileMeta{
		"foo2",
		MsgMeta{11, time.Now().Add(-7 * 24 * time.Hour)},
		MsgMeta{20, time.Now().Add(-6 * 24 * time.Hour)},
	}
	idx["foo_topic"] = append(idx["foo_topic"], foo1Meta)
	idx["foo_topic"] = append(idx["foo_topic"], foo2Meta)

	idx["bar_topic"] = []FileMeta{}
	bar1Meta := FileMeta{
		"bar1",
		MsgMeta{1, time.Now().Add(-5 * 24 * time.Hour)},
		MsgMeta{10, time.Now().Add(-4 * 24 * time.Hour)},
	}
	idx["bar_topic"] = append(idx["bar_topic"], bar1Meta)

	return &idx
}
