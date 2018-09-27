package main

import (
	"flag"
	"log"
	"time"

	"github.com/peterhoward42/toy-kafka/client"
	"github.com/peterhoward42/toy-kafka/protocol"
)

// This command launches a command-line interface to a Toy-Kafka consumer client.
// You provide a topic for it to subscribe to with the -topic flag. It will then
// consume (and report on), both the existing messages in this topic, and newly 
// arriving ones, by polling the server every 3 seconds.
func main() {

	// Extract topic from command line args.

	topic := flag.String("topic", "", "Please specify a topic")
	flag.Parse()
	if *topic == "" {
		log.Fatal("You must specify a topic using the '-topic flag'")
	}

	// Todo - fetch host, and override port from environment variables.
	host := "localhost"
	port := protocol.DefaultPort
	readFrom := 1 // Start consuming at message 1.
	consumer, err := client.NewConsumer(*topic, readFrom, host, port)
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
