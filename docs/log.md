o  Persistence using filesystem
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
            *  does index components package need test?
            o  filestore.go is still too darned big
                *  move unexported store implementation into a new actions 
                   package, anticipating actions.poll etc.
                *  get it to build
                *  consider isolated tests for the store module
                o  isolate encode/decode factory for storedmessage with tests.
                o  revert to getting file store to build
        o  incrementally un-comment the tests in the BackingStore test suite,
           and get filestore to satisfy them.
            o  for each suite failure case, first add filestore specific 
               unit tests and get them to pass. Then suite tests should work.
        o  imports horribly deeply nested ?
        o  msgmeta duplication inside stored message?
        o  consider ditching mem store and / or making test suite more lucid
    o  get the whole system to build and decide where to leave things in head
       for scrutiny, maybe add dev log to readme
    o  this test suite mo - not much good cos doesn't tell you where the
               problem is.
    o  augment test for with tests that close and repopen the store



o  could / should hierarchy be flatter?
o  switch to doc.go
o  don't call it toy anything, maybe micro?
o  say host should be host but can be just port
o  test docco in godoc
o  Tidy and make consistent all log messages
o  should docco show install and run cli?
o  should docco explain about other language clients?
o  should docco show code examples
o  Update doc strings in API to sufficient quality for go doc being useful.
o  Add TLS / and or JWT auth?
o  Containerisation
    o  Maybe seperate repo?
o  Orchestration / Hosted
o  Video / pitch?

----------------------------------------------------------------
store action - what to test?
----------------------------------------------------------------
