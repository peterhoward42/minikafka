// Package backends.contract defines the interface that storage backend
// variants must provide.

package contract

import (
	"time"

	minikafka "github.com/peterhoward42/minikafka"
)

// BackingStore is an interface that offers a core set of CRUD methods
// on a backing store for messages.
type BackingStore interface {

	// Store adds the given message to the sequence of Messages already
	// held in the store for a Topic, and returns the message number thus
	// asigned to it.
	Store(topic string, message minikafka.Message) (
		messageNumber int, err error)

	// RemoveOldMessages invites the store to remove any messages in the 
    // store that were stored before the time specified. The store is allowed to
    // deploy some internal optimisation to **not** remove these messages at
    // this time.
	RemoveOldMessages(maxAge time.Time) error

	// Provide a list of all the messages held for this topic, whose message
	// number is greater than or equal to the specified read-from message
	// number. Returns the messages, and also the advised new read-from message
	// number. (beyond those returned by this invocation).
	Poll(topic string, readFrom int) (messages []minikafka.Message,
		newReadFrom int, err error)

	// DeleteContents empties the store of all its contents.
	DeleteContents() error
}
