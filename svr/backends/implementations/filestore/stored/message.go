package stored

import (
    "time"

	toykafka "github.com/peterhoward42/toy-kafka"
)

// Message is the structure that is serialized, and appended to the
// message log files.
type Message struct {
	Message       toykafka.Message
	CreationTime  time.Time
	MessageNumber int32
}
