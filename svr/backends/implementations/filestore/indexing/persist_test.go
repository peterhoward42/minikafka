package indexing

import (
	"testing"
    "os"
    "io/ioutil"
    "fmt"

	"github.com/stretchr/testify/assert"
)

// Make sure the saving of an index to disk runs without crashing, and that
// restoring an index from that file produces an index which has been
// correctly derived from that file. (We need not test the serialize and
// deserialise logic because that is tested elsewhere. He we are concerned with
// the file IO.
func TestSaveAndRetrieve(t *testing.T) {
    file, err := ioutil.TempFile("", "index_")
    if err != nil {
        msg := fmt.Sprintf("ioutil.Tempfile(): %v", err)
        assert.FailNow(t, msg)
    }
    filepath := file.Name()
    file.Close()
    defer os.Remove(filepath)

	index := makeReferenceIndexForTesting()
    err = index.Save(filepath)
    if err != nil {
        msg := fmt.Sprintf("SaveIndex(): %v", err)
        assert.FailNow(t, msg)
    }
    newIndex := NewIndex()
    err = newIndex.PopulateFromDisk(filepath)
    if err != nil {
        msg := fmt.Sprintf("index.RetrieveIndexFromDisk(): %v", err)
        assert.FailNow(t, msg)
    }
    expected := int32(16)
    assert.Equal(t, expected, newIndex.NextMessageNumberFor("topicA"))
}
