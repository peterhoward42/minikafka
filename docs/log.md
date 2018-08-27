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
    o  Read how receive packets from listener
    o  Make separate object to listen to protocol messages and post debug from it
    o  How test conn?

------------------------------------------------------------------------------

github.com/peterhoward42/toy-kafka/...

    server
    client
        produce
        consume
