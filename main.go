package main

import (
	"context"
	"fmt"
	"github.com/ipfs-search/ipfs-sniffer/logsniffer"
	"github.com/ipfs-search/ipfs-sniffer/logsniffer/hashset"
	shell "github.com/ipfs/go-ipfs-api"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var ipfsURL = "localhost:5001"

// onSigTerm calls f() when SIGTERM (control-C) is received
func onSigTerm(f func()) {
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	var fail = func() {
		<-sigChan
		os.Exit(1)
	}

	var quit = func() {
		<-sigChan

		go fail()

		fmt.Println("Received SIGTERM, quitting... One more SIGTERM and we'll abort!")
		f()
	}

	go quit()
}

func main() {
	// Create context
	ctx, cancel := context.WithCancel(context.Background())

	// Allow SIGTERM / Control-C quit through context
	onSigTerm(cancel)

	// Open shell
	sh := shell.NewShell(ipfsURL)

	// Initialize reader
	reader := logsniffer.Reader{}

	log.Printf("Opening log reader")
	err := reader.Open(ctx, sh)
	if err != nil {
		log.Fatal(err)
	}

	// Close when we're done
	defer func() {
		log.Printf("Closing log reader")

		err := reader.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Create channels for messages/errors/hashes
	msgs := make(chan logsniffer.Message)
	hashes := make(chan logsniffer.HashProvider)
	errc := make(chan error, 1)

	hs := hashset.New()

	// Read messages, asynchroneously
	go reader.Read(msgs, errc)

	// Extract hashes, asynchroneously
	go logsniffer.HashProviderExtractor(ctx, msgs, hashes, errc)

	// Update hashset, asynchroneously
	go hs.FromChannel(ctx, hashes, errc)

	// Process messages
	for {
		select {
		case err := <-errc:
			log.Fatalf("Error processing log messages: %s", err)
			cancel()
		}
	}
}
