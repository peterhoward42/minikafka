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
    *  Update docco
        *  check docs dir
*  Update main git branch with this
*  test / fix readme using Windows testbed, having uninstalled gk and grpc






TODO
o  should there be some tests for the higher level packages?
o  benchmark and tune
o  some protection for filling disks
o  Add authentication option
o  Consider putting up the service using K8s in GCP
o  Consider vid / article

----------------------------------------------------------------
----------------------------------------------------------------
