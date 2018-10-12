package client

import (
	"flag"
	"log"
)

// ParseCommandLine retrieves the host and topic command line flags.
// It treats their absence as a fatal error.
func ParseCommandLine() (topic, host string) {

	flag.StringVar(&topic, "topic", "", "Specify a topic.")
	flag.StringVar(&host, "host", "", "Specify a host.")
	flag.Parse()

	if topic == "" {
		log.Fatal("You must specify a topic with the -topic flag.")
	}
	if host == "" {
		log.Fatal("You must specify a host with the -host flag.\n" +
			"E.g. localhost:9999")
	}
	return topic, host
}
