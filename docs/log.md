    *  Timeout - how to observer errors in server about housekeeping
        *  Coded strategy to receive errors from culling or grpc
        *  Test?
    *  Centralise parsing of command line for client and server.
    *  Shouldn't server serve on all NICs?
    o  Server should not die if consumer subscribes to unknown topic
        o  Put error payload into responses, and read from these at 
           the coalfaces.
            *  change protobuf
            *  generate the go code
            *  work through compile errors
            o  change servers handlers to populate the err field for 
               application level errors
            o  change producer to read err in response
            o  change consumer to read err in response
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
