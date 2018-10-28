package actions

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/peterhoward42/minikafka"
	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore/indexing"
)

func TestSimplestPossibleCase(t *testing.T) {
	// Store a handful of tiny messages in a virgin store, and make sure
	// that a Poll from msg number 1, gives us back all of them, and returns
	// the correct next read-from message number.

	// Prepare a root directory that we can delete after the test.
	rootDir, err := ioutil.TempDir("", "filestore")
	if err != nil {
		msg := fmt.Sprintf("ioutil.TempDir(): %v", err)
		assert.Fail(t, msg)
	}
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
		_, _, err = storeAction.Store()
		if err != nil {
			msg := fmt.Sprintf("storeAction.Store(): %v", err)
			assert.Fail(t, msg)
		}
	}
	readFrom := 1
	action := PollAction{topic, readFrom, index, rootDir}
	messages, _, err := action.Poll()
	if err != nil {
		msg := fmt.Sprintf("action.Poll(): %v", err)
		assert.Fail(t, msg)
	}
	assert.Equal(t, 5, len(messages))
}
