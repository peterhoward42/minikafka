package index

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//--------------------------------------------------------------------------
// API
//--------------------------------------------------------------------------

func TestNextMsgNumForTopic(t *testing.T) {
	index := makeReferenceIndexForTesting()

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
	index := makeReferenceIndexForTesting()
	// Check correct when topic is known and has files registered.
	currentName := index.CurrentMsgFileNameFor("topicA")
	expected := "file2"
	assert.Equal(t, expected, currentName)
	// Check correct when topic is unknown.
	currentName = index.CurrentMsgFileNameFor("nosuchtopic")
	expected = ""
	assert.Equal(t, expected, currentName)
}

func TestIsFilenameOk(t *testing.T) {
	index := makeReferenceIndexForTesting()
	// When should say yes.
	used := index.IsFilenameOk("file1", "topicA")
	expected := true
	assert.Equal(t, expected, used)
	// When should say no because names exist but not this one.
	used = index.IsFilenameOk("unknownname", "topicA")
	expected = false
	assert.Equal(t, expected, used)
	// When should say no because no names exist.
	used = index.IsFilenameOk("file1", "unknowntopic")
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

// makeReferenceIndexForTesting provides a useful, repeatable Index for
// testing purposes.
func makeReferenceIndexForTesting() *Index {

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
