# Using the File System for Storage - Rationale and Design

# Why use the file system for storage?

- Becuase, provided you organise it to avoid seeking inside the message store 
  files, you can take advantage of the operating system's file caching in memory.
  Thus getting similar performance to that of an in-memory store.

  And the operating system's caching is hardened and reliable.

  See more [here]
  (https://medium.freecodecamp.org/what-makes-apache-kafka-so-fast-a8d4f94ab145)


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

- Message storage files are Go GOB serialized binary files and thus can contain
  non-trivial data structures including variable length elements like slices.
- The file simply contains stored message objects concatenated.
- A stored message object contains not only the payload byte slice, but also
  its message number and creation time.

# Rationale

- Avoids the need for seeking inside any of the  message files for the 
  produce, or old-message eviction operations. Reduces the number of seek
  operations needed for a poll operation to one.
- Makes it possible to determine which message files are relavent to each of the
  operations without looking inside any of them.
- Moderates the size of message files, so that when one must be re-written, 
  the cost is constrained.
- Reduces the message data-writing cost of the produce operation to only one 
  append operation to one file.
- Makes it possible to do the old-message eviction operation without mutating
  files - it need only delete whole files.

# Flip-Side of the Rationale Benefits
- The index file must be read and re-written for each of the 3 
  (produce, consume, evict operations. Although it should remain a 
  smallish file.
- Access to the the index file is required to be protected with a mutex, thus 
  serializing access to the entire store.  (Possible enhancement: Topics could 
  be made completely independent, and each have an index of their own.
