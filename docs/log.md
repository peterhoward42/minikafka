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
            o  confusion in filestore between index the object and index the
               package - wha tto do?
        o  filestore.go is too darned big, as are some methods.
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
FileStore design
----------------------------------------------------------------
