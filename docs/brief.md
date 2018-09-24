# Brief for Golang Refresh Exercise

A toy service inspired by Kafka.

# Overview
A Kafka-like server.
Cmd line tools to produce and consume.
Can specify topics.
Only granularity is topic.
No partitions and thus segments.
No clustering.
Start with in-memory (volatile) backing store.
Switch later to Redit backing store.
Simple discard time.
Multiple concurrent clients.
Protocol simplified to minimum viable to make *something*.
