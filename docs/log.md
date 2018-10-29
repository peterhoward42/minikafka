o  Persistence using filesystem
    *  stuff I did before...collapsed
    *  substantial redesign experiment whereby only the raw messages go
       into the message files...
    o  upgrade store action tests to use variable sized messages.
    o  extend poll action tests to use variable sized messages.
    o  extend poll action tests to check all code pathways 
    o  revert to opening out all backingstore suite tests.
    o  Remove duplicated code to get index at start of Action methods
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
