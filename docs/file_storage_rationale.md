# Using the File System for Storage - Rationale and Design

# Why use the file system for storage?

- If you organise it to avoid seeking inside the files, it gains the performance
  arising from the implicity provided, very advanced and reliable in-memory 
  caching. See more [here]
  (https://medium.freecodecamp.org/what-makes-apache-kafka-so-fast-a8d4f94ab145)


# The file/directory schema

- One directory per topic.
- In everything that follows, assume it's per topic directory.
- Messages are stored as they arrive, concatenated in files.
- Once a file has grown to a certain size, a new file is started.
- The files are given arbitrarily unique names.
- The directory contains an index file that enumerates the filename sequence
  and for each: it's lowest and highest message number, and oldest and
  newest message age.

# What's in a message storage file?

- Message storage files are Go GOB serialized binary files and thus can contain
  non-trivial data structures including variable length elements like slices.
- The file contains simply stored message objects concatenated.
- A stored message object contains not only the payload byte slice, but also
  its message number and creation time.

# Rationale

- Avoids the need for seeking inside files for any of the produce, consume or
  old-message eviction operations.
- Makes it possible to determine which files are relavent to each of the
  operations without looking inside any of them.
- Reduces the message data-writing cost of the produce operation to only an append
  operation to one file.
- Makes it possible to do the old-message eviction operation without mutating
  files - it need only delete whole files.
