package stored

import (
	"fmt"
	"testing"
    "time"
    "bytes"

	"github.com/stretchr/testify/assert"

	"github.com/peterhoward42/minikafka"
)

// Check the Encode/Decode methods work as they should using a round-trip,
// equality test.
func TestEncodeDecode(t *testing.T) {
    // Make a message and serialize it.
    msgAlone := minikafka.Message("hello")
    now := time.Now()
    msgToStore := Message{msgAlone, now, 42}

    var buf bytes.Buffer
    err := msgToStore.Encode(&buf)
    if err != nil {
        msg := fmt.Sprintf("msgToStore.Encode(): %v", err)
        assert.Fail(t, msg)
    }

    // Deserialize the byte stream just produces into a fresh message.
    retrievedMsg := Message{}
    err = retrievedMsg.Decode(&buf)
    if err != nil {
        msg := fmt.Sprintf("retrievedMsg.Decode(): %v", err)
        assert.Fail(t, msg)
    }

    // Check its field are what they should be.
    assert.Equal(t, msgAlone, retrievedMsg.Message)

    timeDelta := now.Sub(retrievedMsg.CreationTime)
    assert.Zero(t, timeDelta)

    assert.Equal(t, int32(42), retrievedMsg.MessageNumber)
}
