# Brief for Golang Experience Deepener Exercise

A mini service inspired by Kafka.

# Overview
A Kafka-like server.
Cmd line tools to produce and consume.
Can specify topics.
Only granularity is topic.
No partitions and thus segments.
No clustering.
Start with in-memory (volatile) backing store.
Switch later to mounted file system based store.
Simple discard message expiry time.
Multiple concurrent clients.
gRPC comms protocol.
