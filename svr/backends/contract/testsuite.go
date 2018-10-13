package contract

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// RunBackingStoreTests is a test entry point function that checks all the
// functionality that implementations should provide - by delegating to a set
// of individual test functions. It takes an object that claims to offer the
// interface as an argument. Most of the tests create new instances of the
// implementation object - which they can do readily because the BackingStore
// interface includes a *Create()* method.
func RunBackingStoreTests(t *testing.T, impl BackingStore) {
	testCanStoreToVirginStore(t, impl)
	testCanStoreToVirginStore(t, impl)
	testCanStoreToExistingTopic(t, impl)
	testMessageNumberAllocatedPerTopic(t, impl)
	testRemoveMsgOperatesAcrossTopics(t, impl)
	testRemoveOnEmptyStore(t, impl)
	testRemoveWhenNoneOldEnough(t, impl)
	testRemoveWhenAllOldEnough(t, impl)
	testRemoveWhenOnlySomeOldEnough(t, impl)
	testPollErrorHandlingWhenNoSuchTopic(t, impl)
	testPollWhenTopicIsEmpty(t, impl)
	testNewReadFromAdvancement(t, impl)
}

//----------------------------------------------------------------------------
// Unexported tests.
//----------------------------------------------------------------------------

func testCanStoreToVirginStore(t *testing.T, impl BackingStore) {
	ms := impl.Create()
	msgNum, err := ms.Store("topicA", []byte("hello"))
	assert.Nil(t, err)
	assert.Equal(t, 1, msgNum)
}

func testCanStoreToExistingTopic(t *testing.T, impl BackingStore) {
	ms := impl.Create()
	msgNum, err := ms.Store("topicA", []byte("hello"))
	msgNum, err = ms.Store("topicA", []byte("goodbye"))
	assert.Nil(t, err)
	assert.Equal(t, 2, msgNum)
}

func testMessageNumberAllocatedPerTopic(t *testing.T, impl BackingStore) {
	ms := impl.Create()
	msgNum, _ := ms.Store("topicA", []byte("foo"))
	msgNum, _ = ms.Store("topicA", []byte("bar"))
	msgNum, _ = ms.Store("topicB", []byte("baz"))
	assert.Equal(t, 1, msgNum)
}

func testRemoveMsgOperatesAcrossTopics(t *testing.T, impl BackingStore) {
	ms := impl.Create()
	ms.Store("topicA", []byte("foo"))
	ms.Store("topicB", []byte("bar"))

	maxAge := time.Now()
	nRemoved, _ := ms.RemoveOldMessages(maxAge)
	assert.Equal(t, 2, nRemoved)
}

func testRemoveOnEmptyStore(t *testing.T, impl BackingStore) {
	ms := impl.Create()

	maxAge := time.Now()
	nRemoved, err := ms.RemoveOldMessages(maxAge)

	assert.Nil(t, err)
	assert.Equal(t, 0, nRemoved)
}

func testRemoveWhenNoneOldEnough(t *testing.T, impl BackingStore) {
	ms := impl.Create()
	ms.Store("topicA", []byte("foo"))

	// Remove messages older than one hour ago.
	maxAge := time.Now().Add(time.Duration(-1 * time.Hour))
	nRemoved, err := ms.RemoveOldMessages(maxAge)

	assert.Nil(t, err)
	assert.Equal(t, 0, nRemoved)
}

func testRemoveWhenAllOldEnough(t *testing.T, impl BackingStore) {
	ms := impl.Create()
	ms.Store("topicA", []byte("foo"))

	// Remove messages older than one hour's hence.
	maxAge := time.Now().Add(time.Duration(1 * time.Hour))
	nRemoved, err := ms.RemoveOldMessages(maxAge)

	assert.Nil(t, err)
	assert.Equal(t, 1, nRemoved)
}

func testRemoveWhenOnlySomeOldEnough(t *testing.T, impl BackingStore) {
	ms := impl.Create()
	// Store two messages immediately.
	ms.Store("topicA", []byte("abc"))
	ms.Store("topicA", []byte("def"))
	// Store two more, after a 500ms delay.
	time.Sleep(time.Millisecond * 500)
	ms.Store("topicA", []byte("ghi"))
	ms.Store("topicA", []byte("klm"))
	// Remove those older than 250ms.
	maxAge := time.Now().Add(time.Duration(-250 * time.Microsecond))
	nRemoved, err := ms.RemoveOldMessages(maxAge)
	// Should be two removed.
	assert.Nil(t, err)
	assert.Equal(t, 2, nRemoved)
}

func testPollErrorHandlingWhenNoSuchTopic(t *testing.T, impl BackingStore) {
	ms := impl.Create()

	_, _, err := ms.Poll("XXX", 1)
	assert.EqualError(t, err, "No such topic: XXX")
}

func testPollWhenTopicIsEmpty(t *testing.T, impl BackingStore) {
	ms := impl.Create()
	// Bring topic into being.
	ms.Store("topicA", []byte("foo"))
	// Remove all messages.
	maxAge := time.Now().Add(time.Duration(1 * time.Hour))
	nRemoved, err := ms.RemoveOldMessages(maxAge)
	assert.Nil(t, err)
	assert.Equal(t, 1, nRemoved)

	messages, newReadFrom, err := ms.Poll("topicA", 1)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(messages))
	assert.Equal(t, 1, newReadFrom)
}

func testNewReadFromAdvancement(t *testing.T, impl BackingStore) {
	ms := impl.Create()
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
