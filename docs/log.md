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
    o  Bug - returned next readfrom returned by implementation Poll should be 
       based on final message returned, not counting!
       o  Need a failing test!
            o  Work up test suite
                o  Double check all old enough or none old enough
    o  Reduce the number of uint32 casts by using them natively consider
    o  defaultconfig is a smell

    o  seperate cli for server concedptually from that for clients
        o  should there be a server command in svr tree, not cli?
        o  careful design of what things injected, and from cmd line or from
           env
    o  inject all timings
        o  message age
        o  poll interval
    o  consider different style of docco - active / do this
    o  Use new logger with prefix
    o  switch to doc.go
    o  test docco in godoc
    o  Tidy and make consistent all log messages
    o  Several functions used named arguments but routinely re-create the
        objects?
    o  Test
        o  For resilience, and exemplars
    o  Make sure only client objects exposed to protocol types
o  Consider which error handling can do better than fatal.
o  Update doc strings in API to sufficient quality for go doc being useful.
o  Add TLS / and or JWT auth
o  Package level doc.go

----------------------------------------------------------------
----------------------------------------------------------------
