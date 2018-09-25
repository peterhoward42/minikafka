*  Switch to grpc
    *  Follow go tutorial here: https://grpc.io/docs/tutorials/basic/go.html
        almost verbatim so have skel to modify

*  Clean up comments a bit
o  Now do time based message culling
    *  Let test be > N seconds old. Server setting. 
    *  Swithcing to feature branch
    *  Have server launch a goroutine for this - constructed with a age 
       limit param
    *  Server launch param for now.
    *  Decision making run by svr
    *  Passes message to store (extend interface)
    *  Debug the cycling call and depletion
        *  Compiles
        *  Fix inexorable growth problem of backing array
        *  Runs client/server?
            *  Is returning message number as one?
        *  Add logging of message numbers culled
        *  Runs
        *  Removing the right ones
    *  Add and housekeeping?
o  Consider error handling from the startCulling go routine?
o  Now do getter client
o  Configure host from envars
o  Add tests
o  Add TLS / and or JWT auth

----------------------------------------------------------------
----------------------------------------------------------------
Policy any message been in store longer than N seconds.

Threshold setter on storage interface and constructor, to be stored in the 
store.

