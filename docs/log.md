    *  Agree vernacular:
        *  message, consumer, producer (apis), readFrom / current
        *  drop payload
        *  embrace stored message
        *  Embrace stored message but only in backing impl and private
        o  Poll should say readFrom at every level
            o  protobuf
            o  contract
            o  backing store impl
            o  should find no remaining refs to msgNum or messageNum
        o  Centralise type for byte slice called message
            o  Should no longer need contract.Message
        o  Use readFrom in protobuf
    o  Bug - returned next readfrom returned by implementation Poll should be 
       based on final message returned, not counting!
       o  Need a failing test!

    o  seperate cli for server concedptually from that for clients
        o  should there be a server command in svr tree, not cli?
        o  careful design of what things injected, and from cmd line or from
           env
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
