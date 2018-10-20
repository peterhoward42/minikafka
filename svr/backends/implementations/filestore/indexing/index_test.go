package indexing

import (
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
