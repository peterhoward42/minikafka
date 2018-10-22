// Package backends.contract defines the interface that storage backend
// variants must provide.

package contract

import (
	"time"

	toykafka "github.com/peterhoward42/minikafka"
)

// BackingStore is an interface that offers a core set of CRUD methods
// on a backing store for messages.
type BackingStore interface {

	// Store adds the given message to the sequence of Messages already
	// held in the store for a Topic, and returns the message number thus
	// asigned to it.
	Store(topic string, message toykafka.Message) (
		messageNumber int, err error)

	// RemoveOldMessages removes any messages in the store that were stored
	// before the time specified.
	RemoveOldMessages(maxAge time.Time) (nRemoved int, err error)

	// Provide a list of all the messages held for this topic, whose message
	// number is greater than or equal to the specified read-from message
	// number. Returns the messages, and also the advised new read-from message
	// number. (beyond those returned by this invocation).
	Poll(topic string, readFrom int) (messages []toykafka.Message,
		newReadFrom int, err error)

	// DeleteContents empties the store of all its contents.
	DeleteContents() error
}
