package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "os"

    "github.com/peterhoward42/toy-kafka/client"
)

func main() {

    // Extract topic from command line args.

    topic_ptr := flag.String("topic", "", "Please specify a topic")
    flag.Parse()
    if *topic_ptr == "" {
        log.Fatal("You must specify a topic using the '-topic flag'")
    }

    producer, err := client.NewProducer(*topic_ptr)
    if err != nil {
        log.Fatalf("Failed to create Producer, with error: %v", err)
    }

    // Invite user to enter lines.
    fmt.Printf("Enter messages, one per line, and press ENTER.\n")

    // Ingest lines and send each one as a *produce* command.
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        msg_text := scanner.Text()
        message_num, err := producer.SendMessage(msg_text)

        if err != nil {
            log.Printf("Error from producer: %v", err)
            continue
        }
        log.Printf("Server acknowledged message number: %v", message_num)
    }

	fmt.Printf("Finished\n")
}
