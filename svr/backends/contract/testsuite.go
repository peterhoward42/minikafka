package contract

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// RunBackingStoreTests is a test suite entry point function that checks all the
// functionality that implementations should provide - by delegating to a set
// of individual test functions. You pass in the implementation object you want
// to test.
func RunBackingStoreTests(t *testing.T, implementation BackingStore) {
	testCanStoreToVirginStore(t, implementation)
	testCanStoreToVirginStore(t, implementation)
	testCanStoreToExistingTopic(t, implementation)
	testMessageNumberAllocatedPerTopic(t, implementation)
	testRemoveMsgOperatesAcrossTopics(t, implementation)
	testRemoveOnEmptyStore(t, implementation)
	testRemoveWhenNoneOldEnough(t, implementation)
	testRemoveWhenAllOldEnough(t, implementation)
	testRemoveWhenOnlySomeOldEnough(t, implementation)
	testPollErrorHandlingWhenNoSuchTopic(t, implementation)
	testPollWhenTopicIsEmpty(t, implementation)
	testNewReadFromAdvancement(t, implementation)
}

//----------------------------------------------------------------------------
// Unexported tests.
//----------------------------------------------------------------------------

func testCanStoreToVirginStore(t *testing.T, store BackingStore) {
	store.DeleteContents()
	msgNum, err := store.Store("topicA", []byte("hello"))
	assert.Nil(t, err)
	assert.Equal(t, 1, msgNum)
}

func testCanStoreToExistingTopic(t *testing.T, store BackingStore) {
	store.DeleteContents()
	msgNum, err := store.Store("topicA", []byte("hello"))
	msgNum, err = store.Store("topicA", []byte("goodbye"))
	assert.Nil(t, err)
	assert.Equal(t, 2, msgNum)
}

func testMessageNumberAllocatedPerTopic(t *testing.T, store BackingStore) {
	store.DeleteContents()
	msgNum, _ := store.Store("topicA", []byte("foo"))
	msgNum, _ = store.Store("topicA", []byte("bar"))
	msgNum, _ = store.Store("topicB", []byte("baz"))
	assert.Equal(t, 1, msgNum)
}

func testRemoveMsgOperatesAcrossTopics(t *testing.T, store BackingStore) {
	store.DeleteContents()
	store.Store("topicA", []byte("foo"))
	store.Store("topicB", []byte("bar"))

	maxAge := time.Now()
	nRemoved, _ := store.RemoveOldMessages(maxAge)
	assert.Equal(t, 2, nRemoved)
}

func testRemoveOnEmptyStore(t *testing.T, store BackingStore) {
	store.DeleteContents()

	maxAge := time.Now()
	nRemoved, err := store.RemoveOldMessages(maxAge)

	assert.Nil(t, err)
	assert.Equal(t, 0, nRemoved)
}

func testRemoveWhenNoneOldEnough(t *testing.T, store BackingStore) {
	store.DeleteContents()
	store.Store("topicA", []byte("foo"))

	// Remove messages older than one hour ago.
	maxAge := time.Now().Add(time.Duration(-1 * time.Hour))
	nRemoved, err := store.RemoveOldMessages(maxAge)

	assert.Nil(t, err)
	assert.Equal(t, 0, nRemoved)
}

func testRemoveWhenAllOldEnough(t *testing.T, store BackingStore) {
	store.DeleteContents()
	store.Store("topicA", []byte("foo"))

	// Remove messages older than one hour's hence.
	maxAge := time.Now().Add(time.Duration(1 * time.Hour))
	nRemoved, err := store.RemoveOldMessages(maxAge)

	assert.Nil(t, err)
	assert.Equal(t, 1, nRemoved)
}

func testRemoveWhenOnlySomeOldEnough(t *testing.T, store BackingStore) {
	store.DeleteContents()
	// Store two messages immediately.
	store.Store("topicA", []byte("abc"))
	store.Store("topicA", []byte("def"))
	// Store two more, after a 500ms delay.
	time.Sleep(time.Millisecond * 500)
	store.Store("topicA", []byte("ghi"))
	store.Store("topicA", []byte("klm"))
	// Remove those older than 250ms.
	maxAge := time.Now().Add(time.Duration(-250 * time.Microsecond))
	nRemoved, err := store.RemoveOldMessages(maxAge)
	// Should be two removed.
	assert.Nil(t, err)
	assert.Equal(t, 2, nRemoved)
}

func testPollErrorHandlingWhenNoSuchTopic(t *testing.T, store BackingStore) {
	store.DeleteContents()

	_, _, err := store.Poll("XXX", 1)
	assert.EqualError(t, err, "No such topic: XXX")
}

func testPollWhenTopicIsEmpty(t *testing.T, store BackingStore) {
	store.DeleteContents()
	// Bring topic into being.
	store.Store("topicA", []byte("foo"))
	// Remove all messages.
	maxAge := time.Now().Add(time.Duration(1 * time.Hour))
	nRemoved, err := store.RemoveOldMessages(maxAge)
	assert.Nil(t, err)
	assert.Equal(t, 1, nRemoved)

	messages, newReadFrom, err := store.Poll("topicA", 1)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(messages))
	assert.Equal(t, 1, newReadFrom)
}

func testNewReadFromAdvancement(t *testing.T, store BackingStore) {
	store.DeleteContents()
	// Add 3 messages.
	store.Store("topicA", []byte("foo"))
	store.Store("topicA", []byte("bar"))
	store.Store("topicA", []byte("baz"))
	// Check returned values from a Poll that will empty the topic.
	messages, newReadFrom, err := store.Poll("topicA", 1)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(messages))
	assert.Equal(t, 4, newReadFrom)
	// Check returned values when Polling for newever values when there
	// are none.
	messages, newReadFrom, err = store.Poll("topicA", newReadFrom)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(messages))
	assert.Equal(t, 4, newReadFrom)

	// Check returned values when Polling for newever values when there
	// are some new ones.
	store.Store("topicA", []byte("baz"))
	messages, newReadFrom, err = store.Poll("topicA", newReadFrom)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(messages))
	assert.Equal(t, 5, newReadFrom)
}
