# Project Overview

# Intro / Preamble

# Service Overview

The server offers to store *messages*, when network clients post them to it
from a *produce* client. A messages is just an opaque, arbitrary-length, 
sequence of bytes.

Messages belong to a *topic* (a string), which is specified in calls to the
*produce* API. This can be an existing topic, or it can be used to introduce
a new one.

The messages are stored in strict received-order per topic, in a *stream*,
and messages are said to have a message number - which is their ordinal
position in this stream. Kafka calls this the message offset.

The other type of client is a *consume* client.  These fetch messages from the
server using the *consume* API. Each consuming client *subsribes* to a single
topic, and then periodically calls the *poll* method to fetch messages.

Each consumer client is said to have a current read-from position (the message
number to read next). Then, when it calls the API *poll* method, it will receive
all messages in the store with message numbers greater than or equal to its
current read-from position. The poll operation automatically then advances the
consumer client's read-from position accordingly.

A message fetched with *poll* is only *consumed* from the perspective of that
client. It is consumed only in the sense that the client's current read-from
position is advanced. The message remains in the store to be *consumed* by other
clients. The current read-from position of each client is held by the client
itself; the server has no such concept. Thus clients govern their own 
consumption of the stream autonomously.

Messages do not remain in the store indefinitely; they have an expiry date, after
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
