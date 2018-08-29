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
    *  The defer close cannot before control passes to the go routine with 
       interpreter
    *  See https://appliedgo.net/networking/ and anticipate gob payload
    o  Start scratch cli client to test command dispatch
        o  how does server respond?
            o  succeeds in decoding entire payload, but cant save in a 
                command key instance
                *  send only command key
                o  is dec.Decode() blocking or spinning? = blocking except when
                   EOF is stuck in pipe, then spins ?
                o  wait for ack, then continue to send message

------------------------------------------------------------------------------

github.com/peterhoward42/toy-kafka/...

    server
    client
        produce
        consume

------------------------------------------------------------------------------
Conn interface
------------------------------------------------------------------------------

command comes first as is short int
