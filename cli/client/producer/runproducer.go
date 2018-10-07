package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/peterhoward42/toy-kafka/cli"
	"github.com/peterhoward42/toy-kafka/client"
)

// This is a command line program that encapsulates a Toy-Kafka
// *produce* client. You specify a topic using the -topic flag, and then are
// invited to type in (string) messages, followed by ENTER. These are each sent
// to the server using the *produce* API.
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

	producer, err := client.NewProducer(topic, host)
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
