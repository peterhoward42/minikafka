package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/peterhoward42/toy-kafka/client"
	"github.com/peterhoward42/toy-kafka/protocol"
)

// This is a command line program that encapsulates a Toy-Kafka
// *produce* client. You specify a topic using the -topic flag, and then are
// invited to type in (string) messages, followed by ENTER. These are each sent
// to the server using the *produce* API.
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
	producer, err := client.NewProducer(*topic, host, port)
	if err != nil {
		log.Fatalf("Failed to create Producer, with error: %v", err)
	}

	// Invite user to enter lines.
	fmt.Printf("Enter messages, one per line, and press ENTER.\n")

	// Ingest lines and send each one as a *produce* command.
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		messageTxt := scanner.Text()
		_, err := producer.SendMessage([]byte(messageTxt))

		if err != nil {
			log.Fatalf("Error SendMessage: %v", err)
		}
	}
}
