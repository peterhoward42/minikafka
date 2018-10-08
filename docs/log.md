    *  there is a conceptual bug - that the consumer cli cannot tell which message
       number it has got up to, so how could a real client take responsibility for
       persisting this so that after a restart it could mandate the client object to
       start at the right piont in the stream?
        *  Check is IS available in the underlying protocol
        *  So is it just that the API offered by the client type is deficient?
        *  Upgrade Poll method contract on proxy client
        *  Satisfy compile needs cascade
        *  Upgrade consume cli to do something with the returned nextreadfrom
    o  Use new logger with prefix
    o  Consider which error handling can do better than fatal.
        o  ATM the server stops if a consumer polls an unknown topic
    o  should docco show install and run cli?
    o  should docco show code examples
o  Persistent using filesystem
    o  switch to doc.go
    o  test docco in godoc
    o  Tidy and make consistent all log messages
    o  Test
        o  Consider which code usefully and pracical to test
        o  For resilience, and exemplars
o  Update doc strings in API to sufficient quality for go doc being useful.
o  Add TLS / and or JWT auth
o  Package level doc.go
o  Containerisation
    o  Maybe seperate repo?
o  Orchestration

----------------------------------------------------------------
----------------------------------------------------------------
