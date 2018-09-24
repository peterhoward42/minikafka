// Package backends defines the interface that storage backend variants must
// provide.
package backends

// Message is a data structure that encapsulates an unopinionated /
// opaque message to store.
type Message []byte

// BackingStore is an interface that offers to add a Message to the
// sequence of Messages already held in the container for a Topic, and to
// return the message number thus asigned to it.
type BackingStore interface {
	Store(topic string, message Message) (messageNumber uint32, err error)
}
