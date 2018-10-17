o  Persistent using filesystem
    *  Design
    *  Write up file sys rationale
    *  Refactor memstore tests into a reusable (for file store) suite also.
    *  Deprecate Create in favour of DeleteAll
    o  TDD filesystem solution using now-reusable test suite.
        *  stores must exist in differing packages
        o  start satisfying methods method
            *  write psuedo code in-situ in terms of analysis document
            *  find a way to isolate tests
            o  have a go at the Store virgin test pass
                o  digress to unit test index independently
                    *  start with creation programmatically and save to given
                       file
                    *  upgrade test to being save and restore
                    *  rename topicmessages
                o   resurrect suite test to use index and drive use cases
                    o  just Store API to start with
                        *  implement delete contents
                        o  start working up Store
                            *  NewFileStore should create root dir?
                            *  digress to impl and write test for
                               index.nextMsgNumberForTopic
                            o  made life needlessly difficult working in basenames
                               and index not knowing about directories
                                *  change index
                                o  change filestore
                                    o  focus on store
                                        *  bring up a name resolver
                                        *  impl and test index.CurrentMsgFileNameFor
                                        o  newmsgfilename
                                            o  need test
                                                *  expose reference index
                                                *  implement and test HasBeenUsedForTopic
                                    o  bring in a ioutils package
                                    o  bring in index persister
                    o  should mutex and code per dir?
    o  this test suite mo - not much good cos doesn't tell you where the
               problem is.
    o  augment test for with tests that close and repopen the store



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
