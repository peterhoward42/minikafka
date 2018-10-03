package implementations

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCanStoreToVirginStore(t *testing.T) {
	ms := NewMemStore()
	msgNum, err := ms.Store("topicA", []byte("hello"))
	assert.Nil(t, err)
	assert.Equal(t, 1, msgNum)
}

func TestCanStoreToExistingTopic(t *testing.T) {
	ms := NewMemStore()
	msgNum, err := ms.Store("topicA", []byte("hello"))
	msgNum, err = ms.Store("topicA", []byte("goodbye"))
	assert.Nil(t, err)
	assert.Equal(t, 2, msgNum)
}

func TestMessageNumberAllocatedPerTopic(t *testing.T) {
	ms := NewMemStore()
	msgNum, _ := ms.Store("topicA", []byte("foo"))
	msgNum, _ = ms.Store("topicA", []byte("bar"))
	msgNum, _ = ms.Store("topicB", []byte("baz"))
	assert.Equal(t, 1, msgNum)
}

func TestRemoveMsgOperatesAcrossTopics(t *testing.T) {
	ms := NewMemStore()
	ms.Store("topicA", []byte("foo"))
	ms.Store("topicB", []byte("bar"))

	maxAge := time.Now()
	removed, _ := ms.RemoveOldMessages(maxAge)

	assert.Equal(t, 1, removed["topicA"][0])
	assert.Equal(t, 1, removed["topicB"][0])
}

func TestRemoveOnEmptyStore(t *testing.T) {
	ms := NewMemStore()

	maxAge := time.Now()
	removed, err := ms.RemoveOldMessages(maxAge)

	assert.Nil(t, err)
	assert.Equal(t, 0, len(removed))
}

func TestRemoveWhenNoneOldEnough(t *testing.T) {
	ms := NewMemStore()
	ms.Store("topicA", []byte("foo"))

	// Remove messages older than one hour ago.
	maxAge := time.Now().Add(time.Duration(-1 * time.Hour))
	removed, err := ms.RemoveOldMessages(maxAge)

	assert.Nil(t, err)
	assert.Equal(t, 0, len(removed["topicA"]))
}

func TestRemoveWhenAllOldEnough(t *testing.T) {
	ms := NewMemStore()
	ms.Store("topicA", []byte("foo"))

	// Remove messages older than one hour's hence.
	maxAge := time.Now().Add(time.Duration(1 * time.Hour))
	removed, err := ms.RemoveOldMessages(maxAge)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(removed["topicA"]))
}

func TestRemoveWhenOnlySomeOldEnough(t *testing.T) {
	ms := NewMemStore()
	// Store two messages immediately.
	ms.Store("topicA", []byte("abc"))
	ms.Store("topicA", []byte("def"))
	// Store two more, after a 500ms delay.
	time.Sleep(time.Millisecond * 500)
	ms.Store("topicA", []byte("ghi"))
	ms.Store("topicA", []byte("klm"))
	// Remove those older than 250ms.
	maxAge := time.Now().Add(time.Duration(-250 * time.Microsecond))
	removed, err := ms.RemoveOldMessages(maxAge)
	// Should be two removed.
	assert.Nil(t, err)
	assert.Equal(t, 2, len(removed["topicA"]))
}

func TestPollErrorHandlingWhenNoSuchTopic(t *testing.T) {
	ms := NewMemStore()

	_, _, err := ms.Poll("XXX", 1)
	assert.EqualError(t, err, "Unknown topic: XXX")
}

func TestPollWhenTopicIsEmpty(t *testing.T) {
	ms := NewMemStore()
	// Bring topic into being.
	ms.Store("topicA", []byte("foo"))
	// Remove all messages.
	maxAge := time.Now().Add(time.Duration(1 * time.Hour))
	removed, err := ms.RemoveOldMessages(maxAge)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(removed["topicA"]))

	messages, newReadFrom, err := ms.Poll("topicA", 1)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(messages))
	assert.Equal(t, 1, newReadFrom)
}

func TestNewReadFromAdvancement(t *testing.T) {
	ms := NewMemStore()
	// Add 3 messages.
	ms.Store("topicA", []byte("foo"))
	ms.Store("topicA", []byte("bar"))
	ms.Store("topicA", []byte("baz"))
	// Check returned values from a Poll that will empty the topic.
	messages, newReadFrom, err := ms.Poll("topicA", 1)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(messages))
	assert.Equal(t, 4, newReadFrom)
	// Check returned values when Polling for newever values when there
	// are none.
	messages, newReadFrom, err = ms.Poll("topicA", newReadFrom)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(messages))
	assert.Equal(t, 4, newReadFrom)

	// Check returned values when Polling for newever values when there
	// are some new ones.
	ms.Store("topicA", []byte("baz"))
	messages, newReadFrom, err = ms.Poll("topicA", newReadFrom)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(messages))
	assert.Equal(t, 5, newReadFrom)
}
