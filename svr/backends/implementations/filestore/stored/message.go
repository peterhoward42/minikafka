package stored

import (
	"time"

	minikafka "github.com/peterhoward42/minikafka"
)

// Message is the structure that is serialized, and appended to the
// message log files.
type Message struct {
	Message       minikafka.Message
	CreationTime  time.Time
	MessageNumber int32
}
