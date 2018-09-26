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
        o  Decide division of responsibilities between cli and consumer lib for:
            o  setting topic
            o  pollling interval
            o  deciding what message number a new consumer should start from
            o  holding next message number state
        o  Check against readme
    o  Test polling using the CLI
o  Update doc strings in API to sufficient quality for go doc being useful.
o  Configure host from envars
o  Add tests
o  Add TLS / and or JWT auth
o  Package level readme with refs

----------------------------------------------------------------
----------------------------------------------------------------