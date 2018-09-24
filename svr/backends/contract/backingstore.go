// Package backends.contract defines the interface that storage backend
// variants must provide.

package contract

import "time"

// Message is a data structure that encapsulates an unopinionated /
// opaque message to store.
type Message []byte

// BackingStore is an interface that offers a core set of CRUD methods
// over a backing store for messages.
type BackingStore interface {

	// Store adds the given message to the sequence of Messages already
	// held in the store for a Topic, and return the message number thus
	// asigned to it.
	Store(topic string, message Message) (messageNumber uint32, err error)

	// RemoveOldMessages removes any messages in the store that were stored
	// before the time specified.
	RemoveOldMessages(maxAge time.Time)
}
