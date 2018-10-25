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
	testCanStoreToExistingTopic(t, implementation)
	testMessageNumberAllocatedPerTopic(t, implementation)
	testRemoveMsgOperatesAcrossTopics(t, implementation)
	testRemoveOnEmptyStore(t, implementation)
	testRemoveWhenNoneOldEnough(t, implementation)
	testRemoveWhenAllOldEnough(t, implementation)
	testRemoveWhenOnlySomeOldEnough(t, implementation)
	//testPollErrorHandlingWhenNoSuchTopic(t, implementation)
	//testPollWhenTopicIsEmpty(t, implementation)
	//testNewReadFromAdvancement(t, implementation)
}

//----------------------------------------------------------------------------
// Unexported tests.
//----------------------------------------------------------------------------

func testCanStoreToVirginStore(t *testing.T, store BackingStore) {
	err := store.DeleteContents()
	assert.Nil(t, err)
	msgNum, err := store.Store("topicA", []byte("hello"))
	assert.Nil(t, err)
	assert.Equal(t, 1, msgNum)
}

func testCanStoreToExistingTopic(t *testing.T, store BackingStore) {
	err := store.DeleteContents()
	assert.Nil(t, err)
	msgNum, err := store.Store("topicA", []byte("hello"))
	assert.Nil(t, err)
	msgNum, err = store.Store("topicA", []byte("goodbye"))
	assert.Nil(t, err)
	assert.Equal(t, 2, msgNum)
}

func testMessageNumberAllocatedPerTopic(t *testing.T, store BackingStore) {
	err := store.DeleteContents()
	assert.Nil(t, err)
	msgNum, err := store.Store("topicA", []byte("foo"))
	assert.Nil(t, err)
	msgNum, err = store.Store("topicA", []byte("bar"))
	assert.Nil(t, err)
	msgNum, err = store.Store("topicB", []byte("baz"))
	assert.Nil(t, err)
	assert.Equal(t, 1, msgNum)
}

func testRemoveMsgOperatesAcrossTopics(t *testing.T, store BackingStore) {
	err := store.DeleteContents()
	assert.Nil(t, err)
	_, err = store.Store("topicA", []byte("foo"))
	assert.Nil(t, err)
	_, err = store.Store("topicB", []byte("bar"))
	assert.Nil(t, err)

	maxAge := time.Now()
	nRemoved, err := store.RemoveOldMessages(maxAge)
	assert.Nil(t, err)
	assert.Equal(t, 2, nRemoved)
}

func testRemoveOnEmptyStore(t *testing.T, store BackingStore) {
	err := store.DeleteContents()
	assert.Nil(t, err)

	maxAge := time.Now()
	nRemoved, err := store.RemoveOldMessages(maxAge)

	assert.Nil(t, err)
	assert.Equal(t, 0, nRemoved)
}

func testRemoveWhenNoneOldEnough(t *testing.T, store BackingStore) {
	err := store.DeleteContents()
	assert.Nil(t, err)
	_, err = store.Store("topicA", []byte("foo"))
	assert.Nil(t, err)

	// Remove messages older than one hour ago.
	maxAge := time.Now().Add(time.Duration(-1 * time.Hour))
	nRemoved, err := store.RemoveOldMessages(maxAge)
	assert.Nil(t, err)

	assert.Equal(t, 0, nRemoved)
}

func testRemoveWhenAllOldEnough(t *testing.T, store BackingStore) {
	err := store.DeleteContents()
	assert.Nil(t, err)
	_, err = store.Store("topicA", []byte("foo"))
	assert.Nil(t, err)

	// Remove messages older than one hour's hence.
	maxAge := time.Now().Add(time.Duration(1 * time.Hour))
	nRemoved, err := store.RemoveOldMessages(maxAge)
	assert.Nil(t, err)

	assert.Equal(t, 1, nRemoved)
}

func testRemoveWhenOnlySomeOldEnough(t *testing.T, store BackingStore) {
	err := store.DeleteContents()
	assert.Nil(t, err)
	// Store two messages immediately.
	_, err = store.Store("topicA", []byte("abc"))
	assert.Nil(t, err)
	_, err = store.Store("topicA", []byte("def"))
	assert.Nil(t, err)
	// Store two more, after a 500ms delay.
	time.Sleep(time.Millisecond * 500)
	_, err = store.Store("topicA", []byte("ghi"))
	assert.Nil(t, err)
	_, err = store.Store("topicA", []byte("klm"))
	assert.Nil(t, err)
	// Remove those older than 250ms.
	maxAge := time.Now().Add(time.Duration(-250 * time.Microsecond))
	nRemoved, err := store.RemoveOldMessages(maxAge)
	assert.Nil(t, err)
	// Should be two removed.
	assert.Equal(t, 2, nRemoved)
}

func testPollErrorHandlingWhenNoSuchTopic(t *testing.T, store BackingStore) {
	err := store.DeleteContents()
	assert.Nil(t, err)
	_, _, err = store.Poll("XXX", 1)
	assert.EqualError(t, err, "No such topic: XXX")
}

func testPollWhenTopicIsEmpty(t *testing.T, store BackingStore) {
	err := store.DeleteContents()
	assert.Nil(t, err)
	// Bring topic into being.
	_, err = store.Store("topicA", []byte("foo"))
	assert.Nil(t, err)
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
	err := store.DeleteContents()
	assert.Nil(t, err)
	// Add 3 messages.
	_, err = store.Store("topicA", []byte("foo"))
	assert.Nil(t, err)
	_, err = store.Store("topicA", []byte("bar"))
	assert.Nil(t, err)
	_, err = store.Store("topicA", []byte("baz"))
	assert.Nil(t, err)
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
	_, err = store.Store("topicA", []byte("baz"))
	assert.Nil(t, err)
	messages, newReadFrom, err = store.Poll("topicA", newReadFrom)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(messages))
	assert.Equal(t, 5, newReadFrom)
}
