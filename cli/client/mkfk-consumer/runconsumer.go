package main

import (
	"log"
	"time"

	clientcli "github.com/peterhoward42/minikafka/cli/client"
	clientlib "github.com/peterhoward42/minikafka/client"
)

// This command line program contains a MiniKafka consumer client.
// You provide a topic for it to subscribe to with the -topic flag. It will then
// consume (and report on), both the existing messages in this topic, and newly
// arriving ones, by polling the server every 3 seconds.
func main() {

	topic, host := clientcli.ParseCommandLine()

	// Unlike this examplar consumer command line app, most real-world consumer
	// apps will not start consuming from message 1 at every boot time. But will
	// instead persist their current readFrom message number, so as to only
	// consume messages they have not previously seen.

	readFrom := 1 // Start consuming at message 1.

	// You specify the response timeout for each consumer.Poll() at consumer
	// construction time.
	timeout := time.Duration(500 * time.Millisecond)
	consumer, err := clientlib.NewConsumer(topic, readFrom, timeout, host)
	if err != nil {
		log.Fatalf("client.NewConsumer: %v", err)
	}

	// Poll at regular intervals, reporting one what is thus received.

	// Using a 3 second polling interval is intended to be convenient
	// for a human keeping an eye on what the CLI is reporting.
	ticker := time.NewTicker(3 * time.Second)

	defer ticker.Stop()
	for range ticker.C {
		log.Printf("Polling")
		messages, newReadFrom, err := consumer.Poll()
		if err != nil {
			log.Fatalf("consumer.Poll: %v", err)
		}
		log.Printf("Received %d messages.", len(messages))
		log.Printf("Next message to read advanced to: %d", newReadFrom)
		for _, msg := range messages {
			log.Printf("  %s", string(msg))
		}
		log.Printf("")
	}
}
