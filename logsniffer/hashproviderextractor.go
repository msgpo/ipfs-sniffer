package logsniffer

import (
	"context"
	"fmt"
	"time"
)

// HashProvider represents a discovered hash and its provider.
type HashProvider struct {
	Date     time.Time
	Hash     string
	Provider string
}

// extract returns nil when no HashProvider is found and an error in unexpected situations.
func extract(msg Message) (*HashProvider, error) {
	// Somehow, real life messages are divided into events and operations.
	// This is not properly documented anywhere.
	operationType, _ := msg["Operation"]
	if operationType == "handleAddProvider" {
		rawDate, ok := msg["Start"]
		if !ok {
			return nil, fmt.Errorf("'Start' not found in message: %v", msg)
		}

		date, err := time.Parse("2006-01-02T15:04:05.999999999-07:00", rawDate.(string))
		if err != nil {
			return nil, fmt.Errorf("Error converting 'Start' into time: %w", err)
		}

		rawTags, ok := msg["Tags"]
		if !ok {
			return nil, fmt.Errorf("'Tags' not found in message: %#v", msg)
		}

		tags, ok := rawTags.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Could not convert 'Tags' for message: %#v", msg)
		}

		key, ok := tags["key"].(string)
		if !ok {
			return nil, fmt.Errorf("Could not read 'key' in tags of message: %#v", msg)
		}

		peer, ok := tags["peer"].(string)
		if !ok {
			return nil, fmt.Errorf("Could not read 'peer' in tags of message: %#v", msg)
		}

		return &HashProvider{
			Date:     date,
			Hash:     key,
			Provider: peer,
		}, nil
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

			if hash != nil {
				hashes <- *hash
			}
		}
	}
}
