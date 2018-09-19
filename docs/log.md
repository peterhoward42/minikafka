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
        o  Make client
            *  Cut down example to min viable
            *  Integrate with cli command
            *  Get to compile
                *  Two main()s !
            o  Try against server
                o  Need to configure security
                    o  Server
                    o  Client
        o  Must we have timeout on request
    o  How get go generate to work
    o  Make it do the real thing for a produce against a in-memory backend.
    o  Add in get command


o  Add tests
o  Add TLS
o  Add flags or env config


----------------------------------------------------------------
Correct boundary for server
----------------------------------------------------------------
Cmd should instantiate  and configure server
    o  inject host/port
Then mandate it to serve