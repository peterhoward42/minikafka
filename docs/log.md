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
        o  start with plain producer client with hard code emitssion of produce
           message
            o  design payload for whole produce message, starting with 
               command code
            o  code it
                *  cmd line app taking topic from args
                *  invite enter messages
                *  send each one as entered
                *  echo ack
                o  read how to initiate connn and sned gob, suspect make 
                   persistent connection and hold it
        o  how does server respond?

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