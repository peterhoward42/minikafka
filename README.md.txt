# Project Overview

# Intro / Preamble

# Service Overview

The server offers to store *messages*, when a network produce-client posts one 
to it. A messages is just an opaque, arbitrary-length, sequence of bytes.

Messages belong to a *topic* (a string), which is specified in calls to the
*produce* API. This can be an existing topic, or it can be used to introduce
a new one.

The messages are stored in strict received-order per topic, in a *stream*,
and messages are said to have a message number - which is their ordinal
position in this stream (1,2,3...). Kafka calls this the message offset.

The other type of client is a *consume* client. These fetch messages from the
server using the *consume* API. A consuming client *subsribes* to a topic,
and then periodically calls the *poll* method to fetch new messages from that
topic.

Each consumer client is said to have a current read-from position inn the
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



## Architecture
Mention grpc
mention indirection to backing store
plus dockerize

## Tools
clis
