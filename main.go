package main

import (
	"context"
	"log"

	shell "github.com/ipfs/go-ipfs-api"
)

var ipfsURL = "localhost:5001"

func printlogs(ctx context.Context, sh *shell.Shell) error {
	log.Printf("Opening logger")
	logger, err := sh.GetLogs(ctx)
	if err != nil {
		return err
	}
	defer func() {
		log.Printf("Closing logger")

		err := logger.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Printf("Printing log messages")
	for {
		msg, err := logger.Next()
		if err != nil {
			return err
		}

		eventType, isEvent := msg["event"]
		if isEvent {
			log.Printf("Found event: %s\n", eventType)
		} else {
			operationType, isOperation := msg["Operation"]
			if isOperation {
				log.Printf("Found operation: %s\n", operationType)

			} else {
				log.Printf("Unknown log message: %v\n", msg)

			}
		}
	}

	return nil
}

func main() {
	// Open shell
	sh := shell.NewShell(ipfsURL)

	// Create context
	ctx := context.Background()

	err := printlogs(ctx, sh)
	if err != nil {
		log.Fatal(err)
	}
}
