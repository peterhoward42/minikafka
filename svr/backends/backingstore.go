package backends

// Message contains is an (unopinionated) message to store.
type Message []byte

// BackingStore offers to add a Message to the sequence already held in the
// container for the Topic, and to return the message number thus asigned
// to it.
type BackingStore interface {
	Store(topic string, message Message) (messageNumber uint32, err error)
}
