# Toy Kafka: Project Overview

This code implements a client-server message queue system, along the same lines
conceptually as Kafka.

- Pub/sub model, for streams of *messages* in topics.
- Messages are binary, arbitrary-length, opaque, sequences of bytes.
- Multiple clients may consume the streams independently at their own rate.
- The server preserves messages until they reach a prescribed age.
- Clients post messages using a *produce* API.
- And receive messages using the *consume* API.

It's not intended to be used for real, and is very much simpler than Kafka. It's 
a project conceived to refinforce the author's Go design/coding experience, 
and offer that for scrutiny.

# Components and Architecture

It is a client-server system, with the client libraries offering a Go, Kafka-like
(proxy) API. The underlying communication protocol is gRPC - unlike Kafka's
custom TCP protocol.

The server is implemented as a Go type, with a *Serve* method. And a simple
command-line program is used to instantiate and run it. The server delegates the
storage CRUD operations to a backend component - coupled only by a contractual Go
interface.

There are two distinct client libraries; one a *produce* client (for sending
messages), and the other a *consume* client for receiving them. These too expose
themselves as Go types with appropriate methods. The expected way for a user of
the client code to use it, is to incorporate one or the other of the client 
packages into their own consuming or producing app, and drive it 
programmatically, to suit the needs of the app.

However to make it possible to play with it straight out of the box, the code
provides very simple command line programs to operate the clients. One is a
*produce* CLI that lets you type in strings and have them sent to the server as
messages. The other is *consumer* CLI that polls the server for newly arriving
messages and displays them as they arrive.

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
method, it will receive all messages in the store with message numbers
greater than or equal to its current read-from position. The poll operation
automatically then advances the consumer client's read-from position
accordingly.

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
