# Contracts

# Components

- The server
- Consume CLI
- Produce CLI
- Consume API client
- Produce API client

# Server Contract

CreateTopic(name)

ProduceRequest(topic, key, value)

    - Advances offset
    - Retained for N minutes
    - Errors: NoSuchTopicErr, FileWriteErr, FileCreateErr, TimeoutErr

FetchRequest(topic, from_idx, to_idx)

    - Returns key-value sequence (somehow)
    - Errors: NoSuchTopicErr, NoLongerAvailErr, NotAvailYetErr.

# Housekeeping Contract

    - Message storage deleted after RETENTION (hard coded) minutes

# Produce CLI Contract

    - Invoke CLI with topic cmd line argument.
    - Enter text and press return - sends message synchronously.
    - Auto advances offset index implicitly server-side (competing with
      concurrently sent requests)
    - Echos TCP response status, new offset, server error code

# Consume CLI Contract

    - Invoke CLI with topic
    - Echos successful connection or aborts
    - Polls server, and displays messages as they occur
    - Exits if polling fails




    

- 
