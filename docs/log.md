    *  Agree vernacular:
        *  message, consumer, producer (apis), readFrom / current
        *  drop payload
        *  embrace stored message
        *  Embrace stored message but only in backing impl and private
        *  Poll should say readFrom at every level
        *  Centralise type for byte slice called Message
            *  Delete contract.Message
            *  Work through build errors arising to fix up
        *  Now scan all source for correct nomenclature
    *  Does it still work?
    *  Bug - returned next readfrom returned by implementation Poll should be 
       based on final message returned, not counting!
       *  Need a failing test!
            *  Work up test suite
                *  Double check all old enough or none old enough
    *  Reduce the number of uint32 casts by using them natively consider
        *  Don't bother.

    *  Refactor to make server a first class cli thing whereas clients focussed
       on libs with illustrative cli
        *  Start new branch
        *  Start by documenting in readme the intent, including DI
        *  Now make it so
            *  Directory hierarchy
            *  Name of sserver launch program
    o  Timeout seems build is broken?
    o  Now consider appropriate (12 factor) source ofr each parameterisation 
       in turn
        o  Serving IP address for server
        o  Host and Port to client
        o  Topic provision to clients
        o  Max age
        o  Polling timing params to client
        o  Timeouts

    o  soften readme sayin gnot intended to be used
    o  should docco show install and run cli?
    o  should docco show code examples
    o  Use new logger with prefix
    o  switch to doc.go
    o  test docco in godoc
    o  Tidy and make consistent all log messages
    o  Test
        o  For resilience, and exemplars
o  Consider which error handling can do better than fatal.
o  Update doc strings in API to sufficient quality for go doc being useful.
o  Add TLS / and or JWT auth
o  Package level doc.go
o  Redit backend
o  Containerisation
    o  Maybe seperate repo?
o  Orchestration

----------------------------------------------------------------
----------------------------------------------------------------
