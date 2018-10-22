# MiniKafka: Project Overview

This code implements a client-server message queue system, along the same lines
conceptually as Kafka.

- Pub/sub model, for streams of *messages* in topics.
- Messages are binary, arbitrary-length, opaque, sequences of bytes.
- Multiple clients may post messages into topics.
- Multiple clients may consume the streams of messages independently at their 
  own rate.
- A client's consumption of a message does not remove it from the server's store.
- The server does however, evict messages once they reach a prescribed age.
- Clients post messages using a *produce* API.
- Clients recevie messages using the *consume* API.

The capabilities are very much simpler than those of Kafka. The project was
undertaken to extend and reinforce the author's Go design/coding experience.

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

# Demonstration using the command line tools

To make it possible to play with with the server and clients straight out of the
box, (and to support development), the code also provides a pair of very simple 
command line client apps. One wraps the produce client and lets you type in
strings that it then sends using the produce API. The other wraps the consume
client and shows you messages as it receives them.

## Build and Installation

```
go get github.com/peterhoward42/minikafka
cd $GOPATH/github.com/peterhoward42/minikafka
go install ./...
```

## To run the server from the command line

export MINIKAFKA_HOST=":9999"
export MINIKAFKA_RETENTIONTIME="3s"
server

## To run the produce client from the command line

In a different terminal...

```
produce -host localhost:999 -topic "demo topic"
```

## To run the consume client from the command line

In a different terminal...

```
consume -host localhost:999 -topic "demo topic"
```







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
