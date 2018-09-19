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

	topic := flag.String("topic", "", "Please specify a topic")
	flag.Parse()
	if *topic == "" {
		log.Fatal("You must specify a topic using the '-topic flag'")
	}

	host := "localhost"
	port := 8086
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
		messageNum, err := producer.SendMessage(messageTxt)

		if err != nil {
			log.Printf("Error SendMessage: %v", err)
			continue
		}
		log.Printf("Server acknowledged message number: %v", messageNum)
	}

	fmt.Printf("Finished\n")
}
