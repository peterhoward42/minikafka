    o  Now consider appropriate (12 factor) source ofr each parameterisation 
       in turn
        *  Serving IP address for server
            *  run_server cmd should require this from env
        *  Host and Port to clients
        *  Topic provision to clients
        *  Max age / retention time
        *  Polling interval in housekeeping go-routine = 1/10 of retetion time
        *  Polling timing params to cli client (hard coded)
        o  Client message response Timeouts
            o  Should pass into lib methods and recommend env in cli(s)
        o  update readme? with these facts

    o  soften readme sayin gnot intended to be used
    o  Consider which error handling can do better than fatal.
    o  should docco show install and run cli?
    o  should docco show code examples
    o  Use new logger with prefix
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
