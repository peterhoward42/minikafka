*  Refresh Kafka knowledge
*  Decide components and rough contracts
*  Main design decisions
*  Package names / hierarchy
*  Repo github.com/peterhoward42/toy-kafka/...
o  Getting started
    *  Install go
    *  Create and checkout repo with readme
    *  Add go .gitignore
    *  Sync
o  Start with server
    *  Boiler plate from example in notes
    *  Read how receive packets from listener
    *  Make separate object to listen to protocol messages and post debug 
       from it
    *  How test conn? (netcat)
    o  Bundle wait for command word then hand off
        o  See https://appliedgo.net/networking/ and anticipate gob payload
        o  Do it

------------------------------------------------------------------------------

github.com/peterhoward42/toy-kafka/...

    server
    client
        produce
        consume

------------------------------------------------------------------------------
Conn interface
------------------------------------------------------------------------------

Read(b []byte) (n int, err error)
