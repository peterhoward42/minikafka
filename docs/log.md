o  Persistence using filesystem
    *  stuff I did before...collapsed
    o  resume TDD using a FileStore instance in the BackingStore interface 
       test suite.
        *  start by disabling all the tests in the suite except for
           doing some storage in a virgin store.
        o  freeze in respectable state
            *  sarch for todo etc words
            *  do all the repo tests pass now?
            *  scan the repo - is it respectable from 10,000 feet?
            *  check all docs
            o  ditch the use of the word toy from everything
            o  complete go lint sweep
            o  complete go fmt sweep
            o  make sure clis work today
            o  review update the README
                o  include tree?
            o  commit to main
        o  add example of how to run svr and cli locally
        o  now incrementally re-enable and satisfy the other suite methods
    o  What tests are appropriate for the filestore that are not covered by
       the interface conformity tests?


TODO
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
