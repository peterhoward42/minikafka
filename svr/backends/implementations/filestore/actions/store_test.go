package actions

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	toykafka "github.com/peterhoward42/minikafka"
	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/indexing"
)

//--------------------------------------------------------------------------
// API
//--------------------------------------------------------------------------

// Operate the StoreAction in a context where it is obliged to make a new
// topic directory, and by implication also the inaugural message file for that
// topic, and make sure it doesn't crash, or report errors.
func TestWhenHasToMakeDirectory(t *testing.T) {

	// Prepare a root directory that we can delete after the test.
	rootDir, err := ioutil.TempDir("", "filestore")
	if err != nil {
		msg := fmt.Sprintf("ioutil.TempDir(): %v", err)
		assert.Fail(t, msg)
	}
	defer os.RemoveAll(rootDir)

	// Use the reference index which has just two well-known topics specified.
	index := indexing.MakeReferenceIndex()

	// Create a store-action that cites a topic that is unknown to the index.
	msg := toykafka.Message("some message")
	storeAction := StoreAction{
		Topic:   "neverheardof",
		Message: msg,
		Index:   index,
		RootDir: rootDir,
	}

	// Make sure that executing the store action doesn't fail.
	_, _, err = storeAction.Store()
	if err != nil {
		msg := fmt.Sprintf("storeAction.Store(): %v", err)
		assert.Fail(t, msg)
	}
}

// This test exercises the store action on a virgin store, and thus tests
// the logic used to create topic directories and a virgin message storage file.
func TestVirginState(t *testing.T) {

	// Prepare a root directory that we can delete after the test.
	rootDir, err := ioutil.TempDir("", "filestore")
	if err != nil {
		msg := fmt.Sprintf("ioutil.TempDir(): %v", err)
		assert.Fail(t, msg)
	}
	defer os.RemoveAll(rootDir)

	// Use the reference index which has just two well-known topics specified.
	index := indexing.MakeReferenceIndex()

	// Create a store-action with a small payload that we can use twice.
	msg := toykafka.Message("some message")
	storeAction := StoreAction{
		Topic:   "neverheardof",
		Message: msg,
		Index:   index,
		RootDir: rootDir,
	}

	// Storage works without reporting errors and a plausible message file
	// got used.
	_, msgFileUsed, err := storeAction.Store()
	if err != nil {
		msg := fmt.Sprintf("storeAction.Store(): %v", err)
		assert.Fail(t, msg)
	}
	// Plausible message file used?
	assert.Equal(t, 8, len(msgFileUsed))
}

// Test messages get stored in the same message file, while there is
// is plenty of room.
func TestStorageFileReuse(t *testing.T) {

	// Prepare a root directory that we can delete after the test.
	rootDir, err := ioutil.TempDir("", "filestore")
	if err != nil {
		msg := fmt.Sprintf("ioutil.TempDir(): %v", err)
		assert.Fail(t, msg)
	}
	defer os.RemoveAll(rootDir)

	// Use the reference index which has just two well-known topics specified.
	index := indexing.MakeReferenceIndex()

	// Create a store-action with a small payload that we can use twice.
	msg := toykafka.Message("some message")
	storeAction := StoreAction{
		Topic:   "neverheardof",
		Message: msg,
		Index:   index,
		RootDir: rootDir,
	}

	// Call the store action twice
	msgFilesUsed := make([]string, 2)
	for i := 0; i < 2; i++ {
		_, msgFilesUsed[i], err = storeAction.Store()
		if err != nil {
			msg := fmt.Sprintf("storeAction.Store(): %v", err)
			assert.Fail(t, msg)
		}
	}
	// Make sure the same single file got used.
	assert.Equal(t, msgFilesUsed[0], msgFilesUsed[1])
}

// Operate the StoreAction with two very large messages and make sure that the
// second one causes a new storage file to be used opened.
func TestTwoLargeMessagesGetPutInDifferentFiles(t *testing.T) {

	// Prepare a root directory that we can delete after the test.
	rootDir, err := ioutil.TempDir("", "filestore")
	if err != nil {
		msg := fmt.Sprintf("ioutil.TempDir(): %v", err)
		assert.Fail(t, msg)
	}
	defer os.RemoveAll(rootDir)

	// Use the reference index which has just two well-known topics specified.
	index := indexing.MakeReferenceIndex()

	// Create a store-action with a large payload that we can use twice.
	largeMsg := make([]byte, 0.75*maximumFileSize)
	storeAction := StoreAction{
		Topic:   "neverheardof",
		Message: largeMsg,
		Index:   index,
		RootDir: rootDir,
	}

	// Call the store action twice
	msgFilesUsed := make([]string, 2)
	for i := 0; i < 2; i++ {
		_, msgFilesUsed[i], err = storeAction.Store()
		if err != nil {
			msg := fmt.Sprintf("storeAction.Store(): %v", err)
			assert.Fail(t, msg)
		}
	}
	// Make sure different files got used.
	assert.NotEqual(t, msgFilesUsed[0], msgFilesUsed[1])
}

// Test that the index is left in a properly updated state after some
// messages are stored.
func TestIndexIsUpdated(t *testing.T) {

	// Prepare a root directory that we can delete after the test.
	rootDir, err := ioutil.TempDir("", "filestore")
	if err != nil {
		msg := fmt.Sprintf("ioutil.TempDir(): %v", err)
		assert.Fail(t, msg)
	}
	defer os.RemoveAll(rootDir)

	// Use the reference index which has just two well-known topics specified.
	index := indexing.MakeReferenceIndex()

	// Create a store-action with a small payload that we can use twice.
	topic := "justforthistest"
	msg := toykafka.Message("some message")
	storeAction := StoreAction{
		Topic:   topic,
		Message: msg,
		Index:   index,
		RootDir: rootDir,
	}

	// Call the store action twice
	var msgFileUsed string
	for i := 0; i < 2; i++ {
		_, msgFileUsed, err = storeAction.Store()
		if err != nil {
			msg := fmt.Sprintf("storeAction.Store(): %v", err)
			assert.Fail(t, msg)
		}
	}
	// The tests preceding this one rely on most of the index updating being
	// correct, so we check here only other things.

	msgFileList := index.MessageFileLists[topic]

	// Seems these two messages when gob serialised occupy 248 bytes.
	assert.Equal(t, 248, msgFileList.Meta[msgFileUsed].Size)

	assert.Equal(t, int32(1), msgFileList.Meta[msgFileUsed].Oldest.MsgNum)
	assert.Equal(t, int32(2), msgFileList.Meta[msgFileUsed].Newest.MsgNum)

	expectedT := time.Now() // approx

	oldestT := msgFileList.Meta[msgFileUsed].Oldest.Created
	newestT := msgFileList.Meta[msgFileUsed].Newest.Created
	tolerance := time.Duration(1 * time.Second)

	assert.WithinDuration(t, expectedT, oldestT, tolerance)
	assert.WithinDuration(t, expectedT, newestT, tolerance)
}
