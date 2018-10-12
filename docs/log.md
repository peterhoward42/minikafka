    *  Timeout - how to observer errors in server about housekeeping
        *  Coded strategy to receive errors from culling or grpc
        *  Test?
    *  Centralise parsing of command line for client and server.
    *  Shouldn't server serve on all NICs?
    o  Server should not die if consumer subscribes to unknown topic
        o  Backout putting error field in protobuf responses
            o  How go back on git?
                o  commit all local changes
                o  identify commit you want to be
                o  check that out
                o  commit and push with good comment
            o  Work out which commit we want
    o  Go through all of server replacing sensible fatalf calls to
        returning errors and handling with these in cli
    o  Use new logger with prefix
        *  Necessary? = no
        o  Review remaining loggng calls
    o  Consider which error handling can do better than fatal.
        o  ATM the server stops if a consumer polls an unknown topic
    o  should docco show install and run cli?
    o  should docco show code examples
o  Persistent using filesystem
    o  switch to doc.go
    o  test docco in godoc
    o  Tidy and make consistent all log messages
    o  Test
        o  Consider which code usefully and pracical to test
        o  For resilience, and exemplars
o  Update doc strings in API to sufficient quality for go doc being useful.
o  Add TLS / and or JWT auth
o  Package level doc.go
o  Containerisation
    o  Maybe seperate repo?
o  Orchestration

----------------------------------------------------------------
----------------------------------------------------------------
