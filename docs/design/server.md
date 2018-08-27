# Server

## File System Schema

root/topic/unique_name.messages

    Provides fundamental message storage. The messages concatenated. Each file <=
    MAX_LOG_FILE_SIZE. A new file is created when to append the next message to
    the current file would make it too big.

root/topic/current_file_ptr

    The filename of the .msg file that should be used for the next write.
    Memory mapped.

root/topic/index
    
    Supports the fetch operation. Contains a mapping between every message number
    to the name of the .msg file it lives in, plus its seek index therein and
    size.

    Memory mapped.
    
root/topic/aging.txt

    A file to support automatic purging of old data. Contains records that record
    the latest time at which a message has been added to each .msg file.

    Memory mapped.

## Socket Server

Use https://golang.org/pkg/net/#example_Listener

## Concurrency

Concurrent tcp connections handled by socket server example code.

Singleton sync.Mutex to guard file operations.

## File Housekeeping

Launch a time.NewTicker() at boot time - sending event processing to a
housekeeping function in a seperate goroutine.
