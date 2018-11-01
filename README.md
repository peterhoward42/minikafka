# MiniKafka: Project Overview

This code implements a client-server message-queue system like Kafka, but simpler.

- The server stores *Messages* in topics.
- Messages are just byte sequences. Of arbitrary-length.
- Clients can post messages to a topic using the *Produce* client library.
- Other Clients can subscribe to messages that arrive in topics using the
  *Consume* client library.
- Messages are not removed from the server when they are *consumed*, and
  consuming clients can individually consume the stream at their own rate.
- But Messages **are** removed from the server when they reach a configurable
  *age*; specified in anything from milliseconds to days.

Messages can either be stored to a file-system, or (volatile) memory.

The underlying client/server communciations uses gRPC (unlike Kafka's custom tcp
protocol). But the client libraries hide that away, and expose a few simple API 
methods you can call on proxy objects.

# Getting and Running the Server

    go get google.golang.org/grpc
    go get github.com/peterhoward42/minikafka

    export MINIKAFKA_HOST=":9999"
    export MINIKAFKA_RETENTIONTIME="10s"
    export MINIKAFKA_ROOT_DIR="/tmp/minikafka"

    mkfk-server

# Running a Producer Client

You can try out a simple command line wrapper to the client library:

    mkfk-producer -host localhost:999 -topic topic_foo

This sends each line of text you enter to the server as a *Produce* message.
See below for the more realistic option of embedding the producer (and consumer) 
client in your own code.

# Running a Consumer Client

    mkfk-consumer -host localhost:999 -topic topic_foo

This polls the server every 3 seconds and tells you what it got back.
Remember though, that the messages only live on the server with these settings
for 10 seconds.

# Using the Producer Client Library in Your Own Code

    import (
        "time"
        "github.com/peterhoward42/minikafka/client"
    )

	timeout := time.Duration(500 * time.Millisecond)
	producer, err := client.NewProducer("some_topic", timeout, ":9999")
    message := make([]byte, 300)
    msgNumber, err := producer.SendMessage(message)

# Using the Consumer Client Library in Your Own Code

    import (
        "time"
        "github.com/peterhoward42/minikafka/client"
    )

    // Your app will probably persist the message number to read from (***).
	readFrom := 1 

	timeout := time.Duration(500 * time.Millisecond)
	consumer, err := clientlib.NewConsumer(
            "some_topic", readFrom, timeout, ":9999")

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		messages, newReadFrom, err := consumer.Poll()
        // (***) Persist newReadFrom for the next re-launch.
	}

# Launching the Server From Your Own Code

You can launch the server programmatically (instead of using the CLI wrapper) - and 
inject its configuration like this:


    import (
        "time"

        "github.com/peterhoward42/minikafka/svr/backends/implementations/filestore"
        "github.com/peterhoward42/minikafka/svr"
    )

    rootDir := "/tmp/mkfk-store"
    backingStore, err = filestore.NewFileStore(rootDir) // See also NewMemStore()

	svr := svr.NewServer(backingStore)
    retentionTime := time.Duration(10 * time.Seconds)
	err = svr.Serve(":9999", retentionTime)

# Making Clients in Other Languages

The beauty of gRPC is that you can auto-generate client code in most languages
using the protobuf file in ./protocol/minikafka.proto, and the tools described here:
https://grpc.io/ .








# Components and Architecture

It is a client-server system, with the client libraries offering Go, Kafka-like
(proxy) APIs. The underlying communication protocol is gRPC - unlike Kafka's
custom TCP protocol.

The server is implemented as a Go type, with a *Serve* method. A command line
program is also provided to instantiate and run it - aspiring to 12-factor design
principles whereby runtime configuration choices (such as which 
port to serve on, and the maximum retention time for messages), are read from
environment variables.

The server delegates the storage CRUD operations to a backend component - coupled
only by a contractual Go interface. There are two backend implementations in 
the code. The first is complete, and is a volatile in-memory store created to 
test the client/server and cli parts. The second (real) storage backend uses a
mounted file system - and is a work in progress.

There are two distinct client libraries; one a *produce* client (for sending
messages), and the other a *consume* client for receiving them. These too expose
themselves as Go types with appropriate high-level proxy methods that hide away
the underlying gRPC communications. Users of the client libraries will normally
embed a producer client or a consumer client object in their own client apps, and
configure that object programmatically.

# Service Definition (Conceptual)

This section explains the sub-set of the Kafka services that are available.

The server offers to store *messages*, when a network produce-client posts one 
to it.

Messages belong to a *topic* (a string), which is specified in calls to the
*produce* API. This can be an existing topic, or it can be used to introduce
a new topic.

The messages are stored in strict received-order per topic, in a *stream*,
and messages are said to have a message number - which is their ordinal
position in this stream (1,2,3...). Kafka calls this the message offset.

The other type of client is a *consume* client. These fetch messages from the
server using the *consume* API. A consuming client *subsribes* to a topic,
and then periodically calls the *poll* method to fetch new messages from that
topic.

Each consumer client holds internally a current read-from position in the
stream; (the message number to read next). Then, when it calls the API *poll*
method, it will receive all messages from the store that are more recent than 
this. The poll operation automatically then advances the consumer client's 
read-from position accordingly.

The idea of *consumption* is a concept that exists only on the client-side. I.e.
any one client fetches new messages from a stream and "consumes" them as it
thinks fit.

It follows that the server should not remove messages from its store after any
one client fetches them, However messages do not remain in the store
indefinitely; they have an expiry date, after which the server will remove them
from the store. This does not affect the message numbering scheme. Message
numbers in a stream continue to increment regardless of messages having been
removed.

The expiry date is based on a *maximum message age*, which is set 
programatically at server boot time.
