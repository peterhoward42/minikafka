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
    *  Switch cli to using filestore and test it
            *  design
            o  hypothesise and write
            o  test
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
README overhaul
----------------------------------------------------------------
o say its go!!!!
o file rationale
o  cuter links

o  First section is good but add seconds/minutes/hours/weeks
o  Add storeage is either file system or in mem
o  Uses gRPC cleint server comms under hood
o  But go client libs provided to hide away behind proxy objects.

o  Soften disclaimer

o  Run server like this:

o  Not config options

o  Real world client apps would embed client libs in their own code. But 
   simple demo command line clients provided. Run like this.

o  This

o  If interested in file sys rationale see here
o  For inst in building client apps around client libs see here.