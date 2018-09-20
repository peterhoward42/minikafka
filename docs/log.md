o  Switch to grpc
    *  Follow go tutorial here: https://grpc.io/docs/tutorials/basic/go.html
        almost verbatim so have skel to modify
    o  Ditch selected existing content apart of skel
        *  Commands to run server and client can remain as is.
        *  Producer cut down to almost nothing
        *  send payload ditch
        *  ditch svr interpreter
        *  ditch srv protocol
        *  Cut down server to almost nothing awaiting tutorial shape
        *  Commit and push explaining what is
    o  Literally start adding code as per tutorial
        *  Initiate a protocol dir with toykafka.proto
        *  Define and wrk only with a unary produce message to start with
            *  Make proto fil
        *  Compile it inti toykafka.pb.go and grok tpes and interfaces
        *  Maker server 
            *  Cut down example to min viable
            *  Get it to compile
            *  Call it from cmd
                *  Division of responsibilities wrong between command and
                   server class
            *  Try running it
        *  Make client
            *  Cut down example to min viable
            *  Integrate with cli command
            *  Get to compile
                *  Two main()s !
            *  Try against server
                *  Need to configure security on client
        *  Make backlog / shift stuff into it
        *  No good doing defer conn.Close() in NewProducer, but when?
    *  How get go generate to work
    o  Shared default host/port
    o  Make it do the *real* thing for a produce against a in-memory backend.
        o  Add to store (with mutex)
        o  Increment that topic's next message number (to use)
    o  Now do getter client
    o  Now do time based message evition housekeeping
    o  Configure host from envars
o  Add tests
o  Add TLS / and or JWT auth

----------------------------------------------------------------
----------------------------------------------------------------