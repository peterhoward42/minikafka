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

	maxAge := time.Now().Add(time.Duration(-1 * time.Hour))
	removed, err := ms.RemoveOldMessages(maxAge)

	assert.Nil(t, err)
	assert.Equal(t, 0, len(removed["topicA"]))
}

func TestRemoveWhenAllOldEnough(t *testing.T) {
	ms := NewMemStore()
	ms.Store("topicA", []byte("foo"))

	maxAge := time.Now().Add(time.Duration(1 * time.Hour))
	removed, err := ms.RemoveOldMessages(maxAge)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(removed["topicA"]))
}

/*
Poll deals with no such topic as it should.
Poll behaves when there no messages in the store.
Poll behaves when no messages are newer than time given.
Poll behaves when all messages are newer than time given.
Poll behaves when some messages are newer than time given.

*/
