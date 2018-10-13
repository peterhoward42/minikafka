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
- The files are given names that encode both their creation time, and their 
  first constituent message number - in such a way that sorting the names,
  orders them by creation time.
- No files are used other than these message storage files. Specifically there
  are no supplementary index files.

# What's in a message storage file?

- Message storage files are Go GOB serialized binary files and thus can contain
  non-trivial data structures including variable length elements like slices.
- Most of the file is taken up with stored message data objects simply
  concatenated. 
- A stored message object contains not only the payload byte slice, but also
  its message number and creation time.
- When a file has been closed because its maximum size has been reached, it gains
  a header that states the highest message number stored in the file. Thus only
  the first few bytes need be read to find this out.

# Rationale

- Avoids the need for seeking inside files for any of the produce, consume or
  old-message eviction operations.
- Makes it possible to determine which files are relavent to each of the
  operations without looking inside more than one of them; and then only the
  first few bytes.
- Reduces the data-writing cost of the produce operation to only an append
  operation to one file. (Except when a new file needs to be created, in which
  case it costs a re-write of the old file to put in the header).
- Makes it possible to do the old-message eviction operation without mutating
  files - it need only delete whole files.

This rationale is predicated on the assumption that reading a directory's file
listing on a modern operating system's file system (that caches) is fast.
