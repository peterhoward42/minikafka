package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/peterhoward42/minikafka/svr/backends/implementations/filestore"
	"github.com/peterhoward42/minikafka/svr/backends/implementations/memstore"

	"github.com/peterhoward42/minikafka/svr/backends/contract"

	"github.com/peterhoward42/minikafka/svr"
)

// This commmand-line program instantiates a MiniKafka server and
// mandates it to start serving.
func main() {

	host, retentionTime, rootDir := readEnvironmentVariables()

	// Create an in-memory, or file-based backing store according
	// to the environment variables.
	var err error
	var backingStore contract.BackingStore
	var storeMessage string
	if rootDir == "" {
		backingStore = memstore.NewMemStore()
		storeMessage = "In-memory (volatile) store"
	} else {
		backingStore, err = filestore.NewFileStore(rootDir)
		if err != nil {
			log.Fatalf("filestore.NewFileStore(): %v", err)
		}
		storeMessage = fmt.Sprintf("File-system store rooted at: %s", rootDir)
	}

	svr := svr.NewServer(backingStore)

	log.Printf("Launching server on host: %v", host)
	log.Printf("Using backing store: %s", storeMessage)
	log.Printf("Messages retained for: %v", retentionTime)

	// Server forever, or until an error condition.
	err = svr.Serve(host, retentionTime)
	if err != nil {
		log.Fatalf("svr.Serve: %s", err)
	}

	log.Print("Server Finished")
}

// readEnvironmentVariables fetches the configuration parameters parameterise
// the operation of the server from environment variables.
func readEnvironmentVariables() (
	host string, retentionTime time.Duration, rootDir string) {

	const hostEnvVar string = "MINIKAFKA_HOST"
	const retentionEnvVar string = "MINIKAFKA_RETENTIONTIME"
	const rootDirEnvVar string = "MINIKAFKA_ROOT_DIR"

	host = os.Getenv(hostEnvVar)
	rt := os.Getenv(retentionEnvVar)
	rootDir = os.Getenv(rootDirEnvVar)

	// Host and retention time (unlike root directory) are obligatory.
	if host == "" {
		log.Fatalf("Please set the %s environment variable\n"+
			"E.g. :9999", hostEnvVar)
	}
	if rt == "" {
		log.Fatalf("Please set the %s environment variable\n"+
			"E.g. 3s or 10m", retentionEnvVar)
	}

	retentionTime, err := time.ParseDuration(rt)
	if err != nil {
		log.Fatalf("Error parsing this retention time (%s) from \n"+
			"the %s environment variable: %s", rt, retentionEnvVar, err)
	}
	return host, retentionTime, rootDir
}
