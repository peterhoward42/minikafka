package actions

import (
	"testing"
    "io/ioutil"
    "fmt"
    "os"

	"github.com/stretchr/testify/assert"

	toykafka "github.com/peterhoward42/toy-kafka"
	"github.com/peterhoward42/toy-kafka/svr/backends/implementations/filestore/indexing"
)

//--------------------------------------------------------------------------
// API
//--------------------------------------------------------------------------

// Operate the StoreAction in a context where it is obliged to make a new
// topic directory, and make sure it doesn't crash, or report errors.
func TestWhenHasToMakeDirectory(t *testing.T) {
    rootDir, err := ioutil.TempDir("", "filestore")

    // Prepare a root directory that we can delete after the test.
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
        topic: "neverheardof",
        message: msg,
        index: index,
        rootDir: rootDir,
    }

    // Make sure that executing the store action doesn't fail.
    _, err = storeAction.Store()
    if err != nil {
        msg := fmt.Sprintf("storeAction.Store(): %v", err)
        assert.Fail(t, msg)
    }
}

