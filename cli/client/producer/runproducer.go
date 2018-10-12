package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	clientcli "github.com/peterhoward42/toy-kafka/cli/client"
	clientlib "github.com/peterhoward42/toy-kafka/client"
)

// This is a command line program that encapsulates a Toy-Kafka
// *produce* client. You specify a topic using the -topic flag, and then are
// invited to type in (string) messages, followed by ENTER. These are each sent
// to the server using the *produce* API.
func main() {
	topic, host := clientcli.ParseCommandLine()

	timeout := time.Duration(500 * time.Millisecond)
	producer, err := clientlib.NewProducer(topic, timeout, host)
	if err != nil {
		log.Fatalf("client.NewProducer: %v", err)
	}

	// Invite user to enter lines.
	fmt.Printf("Enter messages, one per line, and press ENTER.\n")

	// Ingest lines and send each one as a *produce* command.
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		messageTxt := scanner.Text()
		_, err := producer.SendMessage([]byte(messageTxt))

		if err != nil {
			log.Fatalf("producer.sendMessage: %v", err)
		}
	}
}
