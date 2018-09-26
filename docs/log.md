*  Switch to grpc
    *  Follow go tutorial here: https://grpc.io/docs/tutorials/basic/go.html
        almost verbatim so have skel to modify

*  Clean up comments a bit
*  Now do time based message culling
o  Timeout to specify *my* client api contracts
    o  In client code or pan-client?
o  Now do getter client
    *  Read / digest kafka model
    o  Add in consume api to docs above
        o  top level readme with section on services and client apis
        o  upgrade produce api code to self doc
        o  check doc is good enough
o  Configure host from envars
o  Add tests
o  Add TLS / and or JWT auth
o  Package level readme with refs

----------------------------------------------------------------
Consumer API
----------------------------------------------------------------
Call it Poll.
Consumer maintains own offset.
Offset defined as message number of first message returned by next Poll.
Consider secondary committed offset maybe.

// might be good to inject starting message number when construct consumer 
// object, and thereafter have the consumer remember it.

consumer = Consumer(connection details)

consumer.Subscribe(topics...) // register interests

consumer.Poll()     // gets any records with offsets >= [topic] onwards.
                    // auto updates consumers offset[topic] accordingly

                    // nb if offset msg has been removed will start from
                    // next avail

                    // might be good for Poll's returned payload to include
                    // which message numbers each message is.

                    // is an error if no topics are subscribed

what happens if you ask for message numbers that don't exist yet?
    = silent return no messages