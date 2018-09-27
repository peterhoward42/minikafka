package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/peterhoward42/toy-kafka/client"
	"github.com/peterhoward42/toy-kafka/protocol"
)

func main() {

	// Extract topic from command line args.

	topic := flag.String("topic", "", "Please specify a topic")
	flag.Parse()
	if *topic == "" {
		log.Fatal("You must specify a topic using the '-topic flag'")
	}

	// Todo - required host, and override port from environment variables.
	host := "localhost"
	port := protocol.DefaultPort
	// We will start reading at the beginning of the stream when the ClI
	// boots.
	readFrom := 1
	consumer, err := client.NewConsumer(*topic, readFrom, host, port)
	if err != nil {
		log.Fatalf("Failed to create Consumer, with error: %v", err)
	}

	// Poll at intervals reporting on these events and what comes back.
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		fmt.Printf("About to poll...")
		messages, err := consumer.Poll()
		if err != nil {
			log.Fatalf("Error generated in call to Poll(): %v", err)
		}
		fmt.Printf("Received these messages...")
		for _, msg := range messages {
			fmt.Printf("%v", msg)
		}
	}
}
