# Toy Kafka: Project Overview

This code implements a client-server message queue system, along the same lines
conceptually as Kafka.

- Pub/sub model, for streams of *messages* in topics.
- Messages are binary, arbitrary-length, opaque, sequences of bytes.
- Multiple clients may to consume the streams at their own rate.
- The server preserves messages until they reach a prescribed age.
- Clients post messages using a *produce* API.
- And receive messages using the *consume* API.

It's not intended to be used for real and is very much simpler than Kafka. It's 
a project conceived to refinforce the author's Go design/coding experience, 
and offer that for scrutiny.

# Components and Architecture

- Client-Server.
- gRPC communication protocol (unlike Kafka's custom TCP protocol).

Logical components:

- server // Package svr.
- single definition of the gRPC schema // package protocol
- produce client(s) // package client
- consume client(s) // package client
- CLI // package cli

The server and clients are exposed by the code as types which can be instantiated
and then operated through their methods. There's also a CLI for each to show
simple exemplar usage of the components and system.

The storage back-end is pluggable, governed by the *BackingStore* interface.

There's only one back-end implementation so far - which is an in-memory, volatile
store.

# Service Definition (Conceptual)

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
method, it will receive all messages in the store with message numbers
greater than or equal to its current read-from position. The poll operation
automatically then advances the consumer client's read-from position
accordingly.

The idea of *consumption* is a concept on the client-side only. I.e. any one 
client fetches new messages from a stream and "consumes" them as it thinks fit.

Whilst the server does not remove messages in response to a client fetching them,
messages do not remain in the store indefinitely; they have an expiry date, after
which the server will remove them from the store. This does not affect the
message numbering scheme. Message numbers in a stream continue to increment
regardless of messages having been removed.

At the time of writing - the expiry date is based on a *maximum message age*,
which is set programatically at server boot time.
