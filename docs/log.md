*  Persistence using filesystem
    *  stuff I did before...collapsed
    *  substantial redesign experiment whereby only the raw messages go
       into the message files...
    *  upgrade store action tests to use variable sized messages.
    *  upgrade all relevant tests to use ioutils.TmpRootDir
    *  extend poll action tests to use variable sized messages.
    *  extend poll action tests to check all code pathways 
    *  fix up file store rationale
    *  complete system regression tests
    *  revert to opening out all backingstore suite tests.
    *  Remove duplicated code to get index at start of Action methods
    o  Switch cli to using filestore and test it
        *  How / where
            *  cli should accept an env var that specifies the root dir
               for the store, or empty implies memstore, and thus bring
               up accordingly.
            *  harvest env vars with erroro handling
            *  create instance of memstore or backing store accordingly
            *  server.NewServer should take backing store instance as an
               argument and use it
            *  provide the new argument from cli
        *  log should say which sort of store
        *  test config as mem
        *  in prep for below have NewFileServer create root dir and null index file
           when not present and test these.
            *  code
            *  test
                *  when rootdir does not exist
                *  round trip virgin to not virgin next message number
        *  Full regression test run
        *  Have the server log the culling time
        o  Restart testing cli for both storeage options
            *  In mem store
            *  filestore where no such dir
            o  filestore on prev used dir
                o  when server from cli using existing rootDir, it says next 
                   stored message is being saved as message number 1?
                    *  get to bottom of
                        *  message culling is losing knowledge of current file
                           and thencenext msg nuber
                            *  make sure capture with failing suite test
                            *  design solution
                            *  code and test the solution incrementally
                                *  code
                                *  test up from index upwards
                            o  rename method in index to say next
                o  after production of 4 messages, consumer clie saying
                   poll only got 1 (which may be ok), but that nextmsgNum
                   advanced to 2 - which is wrong.
                    o  put temp logging in to debug
                    o  get to bottom of
                    o  make sure capture fix in a test
                o  after production of 4 messages, consumer clie saying
        o  can the installed artefacts have less clashing names?
        o  update readme for rootDir env var and expiry
o  Update main git branch with this
    o  Check in and push flat-msg-storage-files
    o  Check out and pull main
    o  Merge flat-xxx
    o  Push to main
o  update readem with build / run script
o  should there be some tests for the higher level packages?
o  some protection for filling disks





TODO
o  Are there some other gob.encode/decode wrappers?
o  What tests are appropriate for the filestore that are not covered by
       the interface conformity tests?
o  switch server cli to use file store
o  find better way of writing backing store test suite
    o  failures should say which method you're in
    o  avoid repetition of set up code
o  add example of how to run svr and cli locally
    o  postponed because Google have broken the way gRPC is structured, such
       that it can't be fetched as a *go get* dependency.
o  move ioutils up out of filestore
o  imports horribly deeply nested ?
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
o  Feasible to performance test?
o  Consider vid / article

----------------------------------------------------------------
----------------------------------------------------------------