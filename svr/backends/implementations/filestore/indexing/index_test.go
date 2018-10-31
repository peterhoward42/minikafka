package indexing

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetAndIncrementMessageNumberFor(t *testing.T) {
	index, _ := MakeReferenceIndex()

	// Check a prepared case.
	nextNum := index.GetAndIncrementMessageNumberFor("topicB")
	assert.Equal(t, int32(7), nextNum)
	// Check the auto-increment side effect.
	nextNum = index.GetAndIncrementMessageNumberFor("topicB")
	assert.Equal(t, int32(8), nextNum)
}

func TestCurrentMsgFileNameFor(t *testing.T) {
	index, _ := MakeReferenceIndex()
	// Check correct when topic is known and has files registered.
	currentName := index.CurrentMsgFileNameFor("topicA")
	expected := "file2"
	assert.Equal(t, expected, currentName)
	// Check correct when topic is unknown.
	currentName = index.CurrentMsgFileNameFor("nosuchtopic")
	expected = ""
	assert.Equal(t, expected, currentName)
}

func TestPreviouslyUsed(t *testing.T) {
	index, _ := MakeReferenceIndex()
	// When should say yes.
	used := index.PreviouslyUsed("file1", "topicA")
	expected := true
	assert.Equal(t, expected, used)
	// When should say no because names exist but not this one.
	used = index.PreviouslyUsed("unknownname", "topicA")
	expected = false
	assert.Equal(t, expected, used)
	// When should say no because no names exist.
	used = index.PreviouslyUsed("file1", "unknowntopic")
	expected = false
	assert.Equal(t, expected, used)
}

func TestSpentFiles(t *testing.T) {
	// Using the reference index, make sure that the identification of spent
	// (i.e. expired) files is just the oldest (file1), when we set the
	// maximum age to be a fraction older than file2.
	index, times := MakeReferenceIndex()
	file2Time := times[3]
	maxAge := file2Time.Add(-time.Duration(20 * time.Millisecond))
	spentFiles := index.MessageFileLists["topicA"].SpentFiles(maxAge)
	sort.Strings(spentFiles)
	expected := []string{"file1"}
	assert.Equal(t, expected, spentFiles)
}

func TestForgetFiles(t *testing.T) {
	// Make sure the index does forget files when told to, and make sure
	// the special cases logic that occur in the implementaion.

	// Case when a name is first in the list of names held.
	index, _ := MakeReferenceIndex()
	forgetThese := []string{"file1"}
	lst := index.MessageFileLists["topicA"]
	lst.ForgetFiles(forgetThese)

	assert.Contains(t, lst.Names, "file2")
	assert.NotContains(t, lst.Names, "file1")

	assert.Contains(t, lst.Meta, "file2")
	assert.NotContains(t, lst.Meta, "file1")

	// Case when a name is last in the list of names held.
	index, _ = MakeReferenceIndex()
	forgetThese = []string{"file2"}
	lst = index.MessageFileLists["topicA"]
	lst.ForgetFiles(forgetThese)

	assert.Contains(t, lst.Names, "file1")
	assert.NotContains(t, lst.Names, "file2")

	assert.Contains(t, lst.Meta, "file1")
	assert.NotContains(t, lst.Meta, "file2")

	// Check when a name is not known to the list it copes silently.
	index, _ = MakeReferenceIndex()
	forgetThese = []string{"neverheardof"}
	lst.ForgetFiles(forgetThese)
}

func TestNumMessagesInFile(t *testing.T) {
	// General case.
	index, _ := MakeReferenceIndex()
	lst := index.MessageFileLists["topicA"]
	n := lst.NumMessagesInFile("file2")
	expected := 3
	assert.Equal(t, expected, n)

	// Case when file is not known.
	lst = NewMessageFileList()
	n = lst.NumMessagesInFile("some file")
	expected = 0
	assert.Equal(t, expected, n)

	// Case when file known, but no messages registered.
	lst = NewMessageFileList()
	lst.RegisterNewFile("some file")
	n = lst.NumMessagesInFile("some file")
	expected = 0
	assert.Equal(t, expected, n)
}

func TestMessageFilesForMessagesFrom(t *testing.T) {
	index, _ := MakeReferenceIndex()
	lst := index.MessageFileLists["topicA"]

	// A message number less than any of those used, should provide
	// the full list.
	files := lst.MessageFilesForMessagesFrom(-99)
	expected := []string{"file1", "file2"}
	assert.Equal(t, expected, files)

	// Message #1 should provide file1 and file2.
	files = lst.MessageFilesForMessagesFrom(1)
	expected = []string{"file1", "file2"}
	assert.Equal(t, expected, files)

	// Similar to above, but use a message number in the middle of file1.
	files = lst.MessageFilesForMessagesFrom(1)
	expected = []string{"file1", "file2"}
	assert.Equal(t, expected, files)

	// Check a message number that should provide just file 2.
	files = lst.MessageFilesForMessagesFrom(5)
	expected = []string{"file2"}
	assert.Equal(t, expected, files)

	// Check a higher message number than any used. Should produce
	// no file names.
	files = lst.MessageFilesForMessagesFrom(9999)
	expected = []string{}
	assert.Equal(t, expected, files)
}

// Add other cases.
