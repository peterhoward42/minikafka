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

    cd $GOPATH/src/github.com/peterhoward42/minikafka
    go install ./...

    export MINIKAFKA_HOST=":9999"
    export MINIKAFKA_RETENTIONTIME="10s"
    export MINIKAFKA_ROOT_DIR="/tmp/minikafka"

    mkfk-server

If you want to use the in-memory store instead of a file-system store:

    export MINIKAFKA_ROOT_DIR=""

# Running a Producer Client

You can try out a simple command line wrapper to the client library:

    mkfk-producer -host localhost:9999 -topic topic_foo

This sends each line of text you enter to the server as a *Produce* message.
See below for the more realistic option of embedding the producer (and consumer) 
client in your own code.

# Running a Consumer Client

    mkfk-consumer -host localhost:9999 -topic topic_foo

This polls the server every 3 seconds and tells you what it got back.
Remember though, that the messages only live on the server with these settings
for 10 seconds.

# Using the Client Libraries in Your Own Code

The more realistic use-case is to incorporate a producer or consumer client
library in your own app - as illustrated by the command line 
[consumer wrapper code](cli/client/mkfk-consumer/runconsumer.go)., or the 
[producer wrapper code](cli/client/mkfk-producer/runproducer.go).


# Launching the Server From Your Own Code

You can similarly wrap the server library in your own code, perhaps to obtain the
configuration from something other than environment variables. See the [server
wrapper code](cli/mkfk-server/runsvr.go).


# Making Clients in Other Languages

The beauty of gRPC is that you can auto-generate client code in most languages
using the protobuf file [here](protocol/minikafka.proto), and 
the [gRPC tools](https://grpc.io/).
