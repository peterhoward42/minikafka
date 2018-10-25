package indexing

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

//--------------------------------------------------------------------------
// API
//--------------------------------------------------------------------------

func TestNextMsgNumForTopic(t *testing.T) {
	index := MakeReferenceIndex()

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
	index := MakeReferenceIndex()
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
	index := MakeReferenceIndex()
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
	// Make sure that the identification of spent (i.e. expired) files
	// is just one file in a given topic, when we set the maximum age
	// to catch file2, but allow file1 to escape.
	index := MakeReferenceIndex()
	maxAge := nowMinusNMinutes(7)
	spentFiles := index.MessageFileLists["topicA"].SpentFiles(maxAge)
	sort.Strings(spentFiles)
	expected := []string{"file2"}
	assert.Equal(t, expected, spentFiles)
}

func TestForgetFiles(t *testing.T) {
	// Make sure the index does forget files when told to, and make sure
	// the special cases logic that occur in the implementaion.

	// Case when a name is first in the list of names held.
	index := MakeReferenceIndex()
	forgetThese := []string{"file1"}
	lst := index.MessageFileLists["topicA"]
	lst.ForgetFiles(forgetThese)

	assert.Contains(t, lst.Names, "file2")
	assert.NotContains(t, lst.Names, "file1")

	assert.Contains(t, lst.Meta, "file2")
	assert.NotContains(t, lst.Meta, "file1")

	// Case when a name is last in the list of names held.
	index = MakeReferenceIndex()
	forgetThese = []string{"file2"}
	lst = index.MessageFileLists["topicA"]
	lst.ForgetFiles(forgetThese)

	assert.Contains(t, lst.Names, "file1")
	assert.NotContains(t, lst.Names, "file2")

	assert.Contains(t, lst.Meta, "file1")
	assert.NotContains(t, lst.Meta, "file2")

	// Check when a name is not known to the list it copes silently.
	index = MakeReferenceIndex()
	forgetThese = []string{"neverheardof"}
	lst.ForgetFiles(forgetThese)
}

func TestNumMessagesInFile(t *testing.T) {
	// General case.
	index := MakeReferenceIndex()
	lst := index.MessageFileLists["topicA"]
	n := lst.NumMessagesInFile("file2")
	expected := 5
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
