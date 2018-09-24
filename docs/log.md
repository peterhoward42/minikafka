*  Switch to grpc
    *  Follow go tutorial here: https://grpc.io/docs/tutorials/basic/go.html
        almost verbatim so have skel to modify

*  Clean up comments a bit
o  Now do time based message culling
    *  Let test be > N seconds old. Server setting. 
    o  Server launch param for now.
    o  Decision making run by svr
    o  Passes message to store (extend interface)
    o  Log the culls from the svr not the store
    o  memstore impl...
o  Now do getter client
o  Configure host from envars
o  Add tests
o  Add TLS / and or JWT auth

----------------------------------------------------------------
----------------------------------------------------------------
Policy any message been in store longer than N seconds.

Threshold setter on storage interface and constructor, to be stored in the 
store.

