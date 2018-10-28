o  Persistence using filesystem
    *  stuff I did before...collapsed
    o  resume TDD using a FileStore instance in the BackingStore interface 
       test suite.
        *  start by disabling all the tests in the suite except for
           doing some storage in a virgin store.
        o  now incrementally re-enable and satisfy the other suite methods
            o  uncomment these tests and fix up one at a time...
                *  Currently failing on testRemoveWhenOnlySomeOldEnough
                o  Move on to poll functionality
                    *  dispand indexcomponents in favour of a module per
                       component
                *  Attempt to write the poll action to drive out what meta
                   data is needed to make it efficient.
                    *  Code (mostly) action.AddMessagesFromFile
                    *  Rename SerializeToBytes encode and make both byte slice
                       and file desc versions of both, and unit tests
                    *  Create and test the generated new function demands.
                        *  msgFileList.FilesContainingThisMessageAndNewer(
                        *  fileMeta.SeekOffsetForMessageNumber
                    o  Revert to testing the poll action.
                        *  The message encode/decode are useless layers?
                        o  Simplest possible test.
                            o  Get to compile and run
                                o  says extra data in buffer
                                    *  someome suggested new decoder for each
                                    o  second iteration is erroring on EOF
                                    o  think gonna have to seek after each decode
                                       to start of next message
                                        o  but doesn't that undermine architecture?
                                            o  stop and think.
                                            o  read entire file into memory and iteratively slice according to seek offsets,
                                               then don't need to seek in filesystem at all!
                                               o  start new branch for this
                                    o  smell having to make a new decoder for
                                       each iteraton.
                            o  Add more checks
        o  Remove duplicated code to get index at start of Action 
           methods
        o  Are there some other gob.encode/decode wrappers?
    o  What tests are appropriate for the filestore that are not covered by
       the interface conformity tests?


TODO
o  find better way of writing backing store test suite
    o  failures should say which method you're in
    o  avoid repetition of set up code
o  add example of how to run svr and cli locally
    o  postponed because Google have broken the way gRPC is structured, such
       that it can't be fetched as a *go get* dependency.
o  move ioutils up out of filestore
o  imports horribly deeply nested ?
o  msgmeta duplication inside stored message?
o  consider ditching mem store - is there any point in keeping it:
    o  to have the test suite in common with filestore, makes the tests
       (being in a suite this way) - less user friendly.
o  switch to doc.go for packages
o  test docco fitness for purpose in godoc
o  Tidy and make consistent all log messages
o  Consider style of top level readme?
    o  consider merit of making ready made docker containers for cli(s) and svr?
    o  should docco have instructions to run svr and cli or just be a code ref?
    o  should docco explain about how gRPC allows other language clients?
    o  should docco show code examples
o  Invest in authentication for this?
o  Consider putting up the service using K8s in GCP

----------------------------------------------------------------
----------------------------------------------------------------