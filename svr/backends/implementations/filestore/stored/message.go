package stored

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	toykafka "github.com/peterhoward42/minikafka"
)

// Message is the structure that is serialized, and appended to the
// message log files.
type Message struct {
	Message       toykafka.Message
	CreationTime  time.Time
	MessageNumber int32
}

// SerializeToBytes provides a gob-encoded serialization of the Message as a
// slice of bytes.
func (m *Message) SerializeToBytes() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(m)
	if err != nil {
		return nil, fmt.Errorf("encoder.Encode(): %v", err)
	}
	return buf.Bytes(), nil
}
