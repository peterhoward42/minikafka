o  Persistent using filesystem
    o  Design
    o  Write up file sys rationale
o  switch to doc.go
o  test docco in godoc
o  Tidy and make consistent all log messages
o  Test
    o  Consider which code usefully and pracical to test
    o  For resilience, and exemplars
o  should docco show install and run cli?
o  should docco explain about other language clients?
o  should docco show code examples
o  Update doc strings in API to sufficient quality for go doc being useful.
o  Add TLS / and or JWT auth
o  Package level doc.go
o  Containerisation
    o  Maybe seperate repo?
o  Orchestration / Hosted
o  Video / pitch?

----------------------------------------------------------------
FileStore design
----------------------------------------------------------------

o  scheme similar to Kafka's
o  file sys 
    o  reason
        o  because performant because os does great job of optimisation
           and memory caching. typ cache as much as can according to avail
           memory. (all as argued by Kafka here)
        o  Way simpler than inventing caching, and os system likely bug-free.
o  dir per topic - no brainer
o  load messages into sequence of files in topic dir as they arrive, using file 
   size limit as split factor. (follows that must limit message size), 
    o  Reason
        o  Most ops need only load few (often one) files into memory.
        o  File mutations need only be whole file deletions for purging 
           old messages.
        o  Only one file mutation needed for produce, and that smallish
o  seek to avoid extra meta files or even in memory indices to eable fast
   operations. Holds given caching arguments above. Seek then only to avoid
   unecessary big-o complexity costs to serve ops.
o  analysis of navigation and read/write needs
    o  publish
        o  which file to use or start new one
        o  what message number to allocate to it
    o  poll
        o  which files to look in?
        o  which messages are higher than msgN?
        o  how to demarcate messages
        o  highest message number that got served,
    o  housekeep
        o  which files can be deleted? - ie which contain onl messages older than
           the expiry date (or no messages)
o consider self-describing message files
    o  ie what can be put in message files themselves to serve all these needs
    o  we can use gob serialization to store non trivial datastructures in the
       files.
    o  so go for fundamental of concatentated sequence of stored message 
       objects. ie the payload byte slice, packed into a structure alongside 
       the message creation time and its message number.
    o  is this logically sufficient (regardless of efficiency)?
        o  yes, provided the filenames imply creation time ordering
o  efficiency analysis
    o  publish
        o  dir contents file names query
        o  sorting file names by name (and thus by creation time) (a short list)
        o  query size of newest file
        o  if using current file
            o  complete sequential read/update/write of whole file contents 
               (XXX BAD)
        o  if using new file
            o  complete sequential read only of current file whole contents 
               to get message number (XXX BAD)
    o  poll
    o  housekeep


o  what about accumulation of topic dirs?
