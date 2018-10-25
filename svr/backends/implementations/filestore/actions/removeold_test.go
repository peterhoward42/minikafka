package actions

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/ioutils"

	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/filenamer"

	"github.com/stretchr/testify/assert"

	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/indexing"
)

// TestRemoveOld makes sure that the construction and operation of
// RemoveOldMessagesAction removes the files it should, and returns the correct
// information about which message files have been deleted.
func TestRemoveOld(t *testing.T) {

	// Prepare a root directory that we can delete after the test.
	rootDir, err := ioutil.TempDir("", "filestore")
	if err != nil {
		msg := fmt.Sprintf("ioutil.TempDir(): %v", err)
		assert.Fail(t, msg)
	}
	defer os.RemoveAll(rootDir)

	// Create an empty index.
	index := indexing.NewIndex()

	// We use a store-action we can use multiple times so as to spawn
	// several message files.
	message := make([]byte, 100000) // Plenty will fit in each file.
	const topic string = "neverheardof"
	storeAction := StoreAction{
		Topic:   topic,
		Message: message,
		Index:   index,
		RootDir: rootDir,
	}

	// Store messages at slight time intervals until the fifth file
	// has been spawned.
	delay := time.Duration(50 * time.Millisecond)
	filesUsed := map[string]bool{}
	var newestInFile2 time.Time
	for len(filesUsed) < 5 {
		_, fileUsed, err := storeAction.Store()
		if err != nil {
			msg := fmt.Sprintf("storeAction.Store(): %v", err)
			assert.Fail(t, msg)
		}
		filesUsed[fileUsed] = true
		if len(filesUsed) == 2 {
			newestInFile2 = time.Now()
		}
		time.Sleep(delay)
	}
	// Set maxAge to target the first two files for deletion.
	maxAge := newestInFile2.Add(time.Duration(10 * time.Millisecond))
	removeAction := RemoveOldMessagesAction{maxAge, index, rootDir}
	filesRemoved, _, err := removeAction.RemoveOldMessages()
	if err != nil {
		msg := fmt.Sprintf("removeAction.RemoveOldMessages(): %v", err)
		assert.Fail(t, msg)
	}
	// Were the correct number of files reported as being removed?
	expected := 2
	assert.Equal(t, expected, len(filesRemoved))

	// Are there exactly 3 files remaining on disk?
	dir := filenamer.DirectoryForTopic(topic, rootDir)
	nFilesRemaining, err := ioutils.CountEntitiesInDir(dir)
	if err != nil {
		msg := fmt.Sprintf("ioutils.CountFilesInDir(): %v", err)
		assert.Fail(t, msg)
	}
	expected = 3
	assert.Equal(t, expected, nFilesRemaining)
}
