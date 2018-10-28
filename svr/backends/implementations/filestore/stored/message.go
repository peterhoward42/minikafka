package stored

import (
	"encoding/gob"
	"fmt"
	"time"
    "io"

	minikafka "github.com/peterhoward42/minikafka"
)

// Message is the structure that is serialized, and appended to the
// message log files.
type Message struct {
	Message       minikafka.Message
	CreationTime  time.Time
	MessageNumber int32
}

// Encode gob-encodes the Message into the writer provided.
func (m *Message) Encode(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(m)
	if err != nil {
		return fmt.Errorf("encoder.Encode(): %v", err)
	}
	return nil
}

// Decode gob-decodes the Message from the reader provided.
func (m *Message) Decode(r io.Reader) error {
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(m)
	if err != nil {
		return fmt.Errorf("decoder.Decode: %v", err)
	}
	return nil
}
