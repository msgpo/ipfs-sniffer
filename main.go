package main

import (
	"context"
	"fmt"
	"github.com/ipfs-search/ipfs-sniffer/logsniffer"
	shell "github.com/ipfs/go-ipfs-api"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var ipfsURL = "localhost:5001"

func processMessage(msg logsniffer.Message) {
	_, isEvent := msg["event"]
	if isEvent {
		// log.Printf("Found event: %s\n", eventType)
		fmt.Printf(".")
	} else {
		operationType, isOperation := msg["Operation"]
		if isOperation {
			// log.Printf("Found operation: %s\n", operationType)
			fmt.Printf(".")

			if operationType == "handleAddProvider" {
				log.Printf("------------- Whooop!!!!")
				log.Printf("%v", msg)
			}
		} else {
			log.Printf("Unknown log message: %v\n", msg)

		}
	}

}

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

	// Create channels for messages/errors
	msgs := make(chan logsniffer.Message, 1)
	errs := make(chan error, 1)

	// Initialize reader
	reader := logsniffer.Reader{
		Errors:   errs,
		Messages: msgs,
	}

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

	// Read messages, asynchroneously
	go reader.Read()

	// Process messages
	for {
		select {
		case err := <-reader.Errors:
			log.Fatalf("Error reading log messages: %s", err)
		case msg := <-reader.Messages:
			processMessage(msg)
		}
	}
}
