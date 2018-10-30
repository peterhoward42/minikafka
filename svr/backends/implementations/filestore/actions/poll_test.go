package actions

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/peterhoward42/minikafka"
	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/indexing"
	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/ioutils"
)

func TestSimplestCase(t *testing.T) {
	// Store a handful of tiny messages in a virgin store, and make sure
	// that a Poll from msg number 1, gives us back all of them, and returns
	// the correct next read-from message number.

	// This test also implicitly checks the special case logic for
	// determining the end of the last message in a file properly identifies
	// the correct final byte.

	rootDir := ioutils.TmpRootDir(t)
	defer os.RemoveAll(rootDir)

	index := indexing.NewIndex()

	msg := minikafka.Message("some message")
	topic := "sometopic"
	storeAction := StoreAction{
		Topic:   topic,
		Message: msg,
		Index:   index,
		RootDir: rootDir,
	}
	for i := 0; i < 5; i++ {
		_, _, err := storeAction.Store()
		if err != nil {
			msg := fmt.Sprintf("storeAction.Store(): %v", err)
			assert.Fail(t, msg)
		}
	}
	readFrom := 1
	action := PollAction{topic, readFrom, index, rootDir}
	messages, newReadFrom, err := action.Poll()
	if err != nil {
		msg := fmt.Sprintf("action.Poll(): %v", err)
		assert.Fail(t, msg)
	}
	assert.Equal(t, 5, len(messages))
	assert.Equal(t, 6, newReadFrom)

	assert.Equal(t, 6, newReadFrom)
}

func TestSeekOffsetsCalculatedRight(t *testing.T) {
	// Store some messages of differing size, and then make sure that a
	// Poll operation that retrieves all of them, gets back the right
	// reconstructed messages.

	// This will validate the storage and use of message seek offsets inside
	// their storage file.

	rootDir := ioutils.TmpRootDir(t)
	defer os.RemoveAll(rootDir)

	index := indexing.NewIndex()

	topic := "sometopic"
	storeAction := StoreAction{
		Topic:   topic,
		Message: []byte{}, // Overwritten before use.
		Index:   index,
		RootDir: rootDir,
	}
	for i := 0; i < 3; i++ {
		msgString := strings.Repeat("X", i+1)
		storeAction.Message = []byte(msgString)
		_, _, err := storeAction.Store()
		if err != nil {
			msg := fmt.Sprintf("storeAction.Store(): %v", err)
			assert.Fail(t, msg)
		}
	}
	readFrom := 1
	action := PollAction{topic, readFrom, index, rootDir}
	messages, newReadFrom, err := action.Poll()
	if err != nil {
		msg := fmt.Sprintf("action.Poll(): %v", err)
		assert.Fail(t, msg)
	}
	assert.Equal(t, 3, len(messages))
	assert.Equal(t, "X", string(messages[0]))
	assert.Equal(t, "XX", string(messages[1]))
	assert.Equal(t, "XXX", string(messages[2]))

	assert.Equal(t, 4, newReadFrom)
}

func TestOnStoreWithNoMessagesRegistered(t *testing.T) {
	// Make sure that the Poll operation on a store that has no messages
	// stored for the poll-topic behaves as it should.

	rootDir := ioutils.TmpRootDir(t)
	defer os.RemoveAll(rootDir)

	index := indexing.NewIndex()
	topic := "foo_topic"
	// This call initialises the index' data structures to know about the
	// topic, but without creating any message files yet.
	index.GetMessageFileListFor(topic)

	readFrom := 1
	action := PollAction{topic, readFrom, index, rootDir}
	messages, newReadFrom, err := action.Poll()
	if err != nil {
		msg := fmt.Sprintf("action.Poll(): %v", err)
		assert.Fail(t, msg)
	}
	assert.Equal(t, 0, len(messages))

	assert.Equal(t, 1, newReadFrom)
}
func TestWhenTopicIsUnknown(t *testing.T) {
	// Check an error is reported when we poll for an unknown topic.

	rootDir := ioutils.TmpRootDir(t)
	defer os.RemoveAll(rootDir)

	index := indexing.NewIndex()

	readFrom := 1
	action := PollAction{"nosuchtopic", readFrom, index, rootDir}
	_, _, err := action.Poll()
	assert.EqualError(t, err, "Unknown topic: nosuchtopic")
}

func TestWhenReadFromIsEarlierThanAllFiles(t *testing.T) {
	// Store a handful of tiny messages and make sure that the Poll action
	// that specifies a read-from message number lower than any stored,
	// provides all the topic's messages.

	rootDir := ioutils.TmpRootDir(t)
	defer os.RemoveAll(rootDir)

	index := indexing.NewIndex()

	msg := minikafka.Message("some message")
	topic := "sometopic"
	storeAction := StoreAction{
		Topic:   topic,
		Message: msg,
		Index:   index,
		RootDir: rootDir,
	}
	for i := 0; i < 5; i++ {
		_, _, err := storeAction.Store()
		if err != nil {
			msg := fmt.Sprintf("storeAction.Store(): %v", err)
			assert.Fail(t, msg)
		}
	}
	readFrom := -999
	action := PollAction{topic, readFrom, index, rootDir}
	messages, newReadFrom, err := action.Poll()
	if err != nil {
		msg := fmt.Sprintf("action.Poll(): %v", err)
		assert.Fail(t, msg)
	}
	assert.Equal(t, 5, len(messages))

	assert.Equal(t, 6, newReadFrom)
}
func TestWhenReadFromIsLaterThanAllFiles(t *testing.T) {
	// Store a handful of tiny messages and make sure that the Poll action
	// that specifies a read-from message number higher than any stored,
	// provides an empty list.

	rootDir := ioutils.TmpRootDir(t)
	defer os.RemoveAll(rootDir)

	index := indexing.NewIndex()

	msg := minikafka.Message("some message")
	topic := "sometopic"
	storeAction := StoreAction{
		Topic:   topic,
		Message: msg,
		Index:   index,
		RootDir: rootDir,
	}
	for i := 0; i < 5; i++ {
		_, _, err := storeAction.Store()
		if err != nil {
			msg := fmt.Sprintf("storeAction.Store(): %v", err)
			assert.Fail(t, msg)
		}
	}
	readFrom := 999
	action := PollAction{topic, readFrom, index, rootDir}
	messages, newReadFrom, err := action.Poll()
	if err != nil {
		msg := fmt.Sprintf("action.Poll(): %v", err)
		assert.Fail(t, msg)
	}
	assert.Equal(t, 0, len(messages))

	assert.Equal(t, 999, newReadFrom)
}
func TestWhenReadFromIsInMiddleOfFile(t *testing.T) {
	// Store a handful of tiny messages and make sure that the Poll action
	// that specifies a read-from message number somewhere in the middle of
	// a file provides those messages from the targeted on onwards.

	rootDir := ioutils.TmpRootDir(t)
	defer os.RemoveAll(rootDir)

	index := indexing.NewIndex()

	topic := "sometopic"
	storeAction := StoreAction{
		Topic:   topic,
		Message: []byte{}, // Overwritten before use.
		Index:   index,
		RootDir: rootDir,
	}
	for i := 0; i < 5; i++ {
		msgString := strings.Repeat("X", i+1)
		storeAction.Message = []byte(msgString)
		_, _, err := storeAction.Store()
		if err != nil {
			msg := fmt.Sprintf("storeAction.Store(): %v", err)
			assert.Fail(t, msg)
		}
	}
	readFrom := 3
	action := PollAction{topic, readFrom, index, rootDir}
	messages, newReadFrom, err := action.Poll()
	if err != nil {
		msg := fmt.Sprintf("action.Poll(): %v", err)
		assert.Fail(t, msg)
	}
	assert.Equal(t, 3, len(messages))
	assert.Equal(t, "XXX", string(messages[0]))
	assert.Equal(t, "XXXX", string(messages[1]))
	assert.Equal(t, "XXXXX", string(messages[2]))

	assert.Equal(t, 6, newReadFrom)
}
func TestResultsWhenFromMultipleFiles(t *testing.T) {
	// Store some big messages that force several files to be created, and
	// make sure the Poll results aggregates them as it should.

	rootDir := ioutils.TmpRootDir(t)
	defer os.RemoveAll(rootDir)

	index := indexing.NewIndex()

	topic := "sometopic"
	message := make([]byte, 200e3) // Big.
	storeAction := StoreAction{
		Topic:   topic,
		Message: message,
		Index:   index,
		RootDir: rootDir,
	}
	for i := 0; i < 20; i++ {
		_, _, err := storeAction.Store()
		if err != nil {
			msg := fmt.Sprintf("storeAction.Store(): %v", err)
			assert.Fail(t, msg)
		}
	}
	readFrom := 1
	action := PollAction{topic, readFrom, index, rootDir}
	messages, newReadFrom, err := action.Poll()
	if err != nil {
		msg := fmt.Sprintf("action.Poll(): %v", err)
		assert.Fail(t, msg)
	}
	assert.Equal(t, 20, len(messages))
	assert.Equal(t, 6, newReadFrom)
}
