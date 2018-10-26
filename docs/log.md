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
                        *  MsgMeta and fix tests
                        *  FileMeta and fix tests
                        *  MessageFileList and fix tests
                        *  ditch components module and fix tests
                        *  check comments
                *  Attempt to write the poll action to drive out what meta
                   data is needed to make it efficient.
                    o  List the generated new function demands.
                    o  Write and test them.
                    o  Rever to testing the poll action.
        o  Remove duplicated code to get index at start of Action 
           methods
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
