o  Persistent using filesystem
    *  bring in filenamer
    *  bring in a ioutils package
    *  index saver / loader
    *  make index serializer be separate?
    *  some quality issues with index.go incl duff logic in ismsgnumberok
        *  change interface to being used
        *  back fix namer and index tests
    o  resume testing at store method in suite test
        o  work thru towards getting call to suite to work, with only store
           enabled in the suite
            *  call the package indexing
            *  move most indexer functions into being methods
            *  get to build
            *  does index components need test?
            o  filestore.go is still too darned big
                o  move unexported store into own module
                o  get it to build
                o  consider isolated tests for the store module
                    o  adopt actions model with a StoreAction which has access
                       to a store and an index. Each action updates the index,
                       but the caller saves it.
                        *  get storeaction to build
                        o  write storeaction tests
                            *  consider what to test (bottom)
                            *  make reference store public and package of own
                            o  design and write/pass tests
                        o  move encode/decode of stored.message into stored module
                           and do round trip test plus read of creation time
                        o  get them to pass
                o  revert to getting file store to build
        o  vet package names to be short single words
        o  imports horribly deeply nested
        o  msgmeta reused inside stored message?
        o  consider ditching mem store and / or making test suite more lucid
    o  this test suite mo - not much good cos doesn't tell you where the
               problem is.
    o  augment test for with tests that close and repopen the store



o  could hierarchy be flatter?
o  switch to doc.go
o  don't call it toy anything
o  say host should be host but can be just port
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
store action - what to test?
----------------------------------------------------------------

o  works when has to create dir
o  splills to new file when wont fit
o  works when has to create msg file for first time
o  works when has to reuse existing msg file
o  does append the payload to the file (gets bigger)
o  index is properly updated
o  stored message gets now creation time
