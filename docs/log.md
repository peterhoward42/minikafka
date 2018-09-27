*  Switch to grpc
    *  Follow go tutorial here: https://grpc.io/docs/tutorials/basic/go.html
        almost verbatim so have skel to modify

*  Clean up comments a bit
*  Now do time based message culling
*  Timeout to specify *my* client api contracts
    * Do in project README at conceptual level.
o  Now do getter client
    *  Read / digest kafka model
    *  Add in consume api to docs above
    *  Augment protocol and compile it
    *  Code server handler
    *  Code back-end implementation
    *  Get this part to compile and run passively (run server)
    o  Create CLI poll client, similar to produce client
        *  becide division of responsibilities between cli and consumer lib for:
            *  setting topic client = object construction
            *  polling interval = outside world
            *  deciding what message number a new consumer should start from
                    = client object construction
            *  holding next message number state - the client object
        *  Check against readme
        *  Code API for client object
        *  Build CLI round it
            *  Get what have to compile and debug timed loop
            *  Put real reads inside the loop
                *  Timeout - produce is using string as message - clear up
                    *  What types do we have outside of protocol for the 
                       message type?
        *  Double check conume protected by a mutex

    o  Nail vernacular - starting with readme and write here.
        *  message
        *  consumer producer api
        *  read-from position // current / nex
        o  update manually written code
            o  ditch MessagePayload in favour of Message
            o  ditch all use of word payload
            o  ditch messagestorage somehow
        o  updae protobuf to say read from instead of message number
    o  Use new logger with prefix
    o  switch to doc.go
    o  test docco in godoc

    o  Improve names also in protocol land
    o  Use throughout
    o  Tidy and make consistent all log messages
    o  Several functions used named arguments but routinely re-create the
        objects?
    o  Test
            o  Maybe a demo of all doing stuff from differeing go routines?
    o  Make sure not needlessly duplicated types like message
    o  Make sure only client objects exposed to protocol types
o  Update doc strings in API to sufficient quality for go doc being useful.
o  Configure host from envars
o  Add tests
o  Add TLS / and or JWT auth
o  Package level readme with refs

----------------------------------------------------------------
----------------------------------------------------------------