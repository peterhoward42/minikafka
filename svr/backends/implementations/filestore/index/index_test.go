package index

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
	index := makeReferenceIndex()
	fileName := path.Join(os.TempDir(), "saveandrestoreindex")
	err := index.Save(fileName)
	if err != nil {
		t.Fatalf("index.Save: %v", err)
	}
	restoredIndex := NewIndex()
	err = restoredIndex.PopulateFromDisk(fileName)
	if err != nil {
		t.Fatalf("PopulateFromDisk: %v", err)
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

	// Make some entries for the "foo_topic".
	msgFileList := []FileMeta{}
	fileMeta := FileMeta{
		"foo1",
		MsgMeta{1, time.Now().Add(-9 * 24 * time.Hour)},
		MsgMeta{10, time.Now().Add(-8 * 24 * time.Hour)},
	}
	msgFileList = append(msgFileList, fileMeta)

	fileMeta = FileMeta{
		"foo2",
		MsgMeta{11, time.Now().Add(-7 * 24 * time.Hour)},
		MsgMeta{20, time.Now().Add(-6 * 24 * time.Hour)},
	}
	msgFileList = append(msgFileList, fileMeta)
	idx.MessageFileLists["foo_topic"] = msgFileList

	// Make some entries for the "bar_topic".
	msgFileList = []FileMeta{}
	fileMeta = FileMeta{
		"bar1",
		MsgMeta{21, time.Now().Add(-5 * 24 * time.Hour)},
		MsgMeta{30, time.Now().Add(-4 * 24 * time.Hour)},
	}
	msgFileList = append(msgFileList, fileMeta)
	idx.MessageFileLists["bar_topic"] = msgFileList

	// Introduce "baz_topic", but record no message files for it.
	msgFileList = []FileMeta{}
	idx.MessageFileLists["baz_topic"] = msgFileList

	return idx
}
