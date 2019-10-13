import (
	"error"
	"log"
	"time"
)

// HashProvider represents a discovered hash and its provider.
type HashProvider struct {
	Date     time.Date
	Hash     string
	Provider string
}

func extract(msg Message) (*HashProvider, error) {
	// Somehow, real life messages are divided into events and operations.
	// This is not properly documented anywhere.
	operationType, _ := msg["Operation"]
	if operationType == "handleAddProvider" {
		date, exists := msg["Start"]
		if !exists {
			return nil, error.Errf("No date in message: %v", msg)
		}

		tags, exists := msg["Tags"]
		if !exists {
			return nil, error.Errf("No tags in message: %v", msg)
		}

		key, exists := tags["key"]
		if !exists {
			return nil, error.Errf("No key in tags of message: %v", msg)
		}
		peer, exists := tags["peer"]
		if !exists {
			return nil, error.Errf("No peer in tags of message: %v", msg)
		}

		return &HashProvider{
			Date:     date,
			Hash:     key,
			Provider: peer,
		}
	}

	return nil, nil
}

// HashProviderExtractor filters out handleAddProvider messages and writes the relevant data to a channel.
func HashProviderExtractor(ctx context.Context, msgs <-chan Message, hashes chan<- HashProvider, errc chan<- error) {
	for {
		select {
		case <-ctx.Done():
			// Context closed, return
			errc <- ctx.Err()
			return
		case msg := <-msgs:
			hash, err := extract(msg)
			if err != nil {
				errc <- err
			}

			hashes <- hash
		}
	}
}
