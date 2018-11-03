# Using the File System for Storage - Rationale and Design

# Why use the file system for storage?

- Because Kafka does. Reasoning that modern operating systems do a great 
  job of caching recently accessed files in memory, which means that most of 
  the blocking IO is to/from memory and therefore fast.

- And these caching systems are more reliable and battle-hardened than trying to
  write one yourself.

  See more [here](https://medium.freecodecamp.org/what-makes-apache-kafka-so-fast-a8d4f94ab145)

# The file/directory schema

- One directory per topic.
- Messages are stored as they arrive, concatenated in files.
- Once a file has grown to a certain size, a new file is started.
- The files are given arbitrarily unique names.

- The parent directory contains an index file that enumerates, for each topic,
  the filename sequence, and for each: it's lowest and highest message 
  number, the oldest and newest message age, and the seek offset for each
  message number.

# What's in a message storage file?

- Message storage files are simply the byte sequences comprising the messages -
  concatenated. A message file, in of itself, has no way of knowing where
  one message stops, and the next starts.

# Rationale

- The availability of the index almost completely avoids any (slow) seeking 
  operations inside files.
- The seek-like behaviour to delimit and fetch messages for the Poll operation
  happens on memory slices, after the necessary message store 
  files have been read, in their entirety into memory.
- Makes it possible to determine which message files are relavent to each of the
  operations without looking inside any of them.
- Moderates the size of message files, so that when one must be read into memory 
  the cost is constrained.
- Reduces the message data-writing cost of the produce operation to only one 
  append operation to one file.
- Makes it possible to do the old-message eviction operation without mutating
  files - it need only delete whole files.
- The random-looking file names for message storage files avoids any risk of
  people thinking the names have semantic significance and then mistakenly 
  relying on this.

# Flip-Side of the Rationale Benefits
- It does not scale horizontally.
- The index file must be read and re-written for each of the 3 
  (produce, consume, evict operations. Although it should remain a 
  relative small file in comparison with the message storage files. And the
  serialize/deserialize steps are relatively fast - using Gob encoding.
- Access to the the index file is required to be protected with a mutex, thus 
  serializing access to the entire store.  (Possible enhancement: Topics could 
  be made completely independent, and each have an index of their own.
