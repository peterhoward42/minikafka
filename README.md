# MiniKafka in Go: Project Overview

This Go code implements a client-server message-queue system 
like [Kafka](https://kafka.apache.org/) but simpler.

- The server stores *Messages* in topics.
- Messages are just byte sequences; of arbitrary-length.
- Clients can post messages to a topic using the *Produce* client library.
- Other Clients can subscribe to messages as they arrive using the
  *Consumer* client library.
- Messages are not removed from the server when they are consumed.
- Thus, individual consuming clients can consume the stream at their own rate.
- But Messages **are** removed from the server when they reach a configurable
  *age*; specified in anything from milliseconds to days.

Messages can either be stored to a file-system, or (volatile) memory.

Read the [rationale](docs/file_storage_rationale.md) for using the file-system 
for storage.

The underlying client/server communciations uses gRPC (unlike Kafka's custom tcp
protocol). But the client libraries hide that away, and expose a few simple API 
methods you can call on proxy objects.

# Status

- It works
- It has decent unit test coverage
- The design, code and documentation is respectable, but has not been peer
  reviewed yet.
- It is still a work in progress with a few minor todos.
- It hasn't been benchmarked, tuned or stress-tested yet.

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

You can launch the server programmatically (instead of using the command line 
wrapper) - and inject its configuration like this:


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
using the protobuf file [here](protocol/minikafka.proto), and 
the [gRPC tools](https://grpc.io/).
