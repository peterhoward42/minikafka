package main

import (
	"flag"
	"log"
	"time"

	"github.com/peterhoward42/toy-kafka/cli"
	"github.com/peterhoward42/toy-kafka/client"
)

// This command line program contains a Toy-Kafka consumer client.
// You provide a topic for it to subscribe to with the -topic flag. It will then
// consume (and report on), both the existing messages in this topic, and newly
// arriving ones, by polling the server every 3 seconds.
func main() {

	// Extract command line arguments.

	var topic, host string
	flag.StringVar(&topic, "topic", "", "Specify a topic.")
	flag.StringVar(&host, "host", cli.DefaultHost, "Specify a host.")
	flag.Parse()

	if topic == "" {
		log.Fatal("You must specify a topic with the -topic flag.")
	}
	if host == cli.DefaultHost {
		log.Printf(
			"Warning, using default host: %s.\nBetter to specify one with -host flag.",
			cli.DefaultHost)
	}

	// Unlike this examplar consumer command line app, most real-world consumer
	// apps will not start consuming from message 1 at every boot time. But will
	// instead persist their current readFrom message number, so as to only
	// consume messages they have not previously already seen.

	readFrom := 1 // Start consuming at message 1.
	consumer, err := client.NewConsumer(topic, readFrom, host)
	if err != nil {
		log.Fatalf("Failed to create Consumer, with error: %v", err)
	}

	// Poll at regular intervals, reporting one what is thus received.
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		log.Printf("CLI Poll")
		messages, err := consumer.Poll()
		if err != nil {
			log.Fatalf("Error generated in call to Poll(): %v", err)
		}
		for n, msg := range messages {
			log.Printf("%d: %s", n, string(msg))
		}
	}
}
