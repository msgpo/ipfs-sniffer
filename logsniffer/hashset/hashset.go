package hashset

import (
	"context"
	"github.com/ipfs-search/ipfs-sniffer/logsniffer"
	"github.com/wangjia184/sortedset"
	"log"
)

// HashSet wraps a sorted set with Hash -> HashProvider objects with Date as the score.
type HashSet struct {
	set *sortedset.SortedSet
}

// New creates a new HashSet.
func New() *HashSet {
	set := sortedset.New()

	return &HashSet{
		set: set,
	}
}

// Add HashProvider to HashSet, returns true when added, false when updated.
func (hs *HashSet) Add(hp *logsniffer.HashProvider) bool {
	return hs.set.AddOrUpdate(hp.Hash, sortedset.SCORE(hp.Date.Unix()), hp)
}

// FromChannel takes HashProvider objects from a channel and adds them to HashSet.
func (hs *HashSet) FromChannel(ctx context.Context, c <-chan logsniffer.HashProvider, errc chan<- error) {
	for {
		select {
		case <-ctx.Done():
			// Context closed, return
			errc <- ctx.Err()
			return
		case hp := <-c:
			added := hs.Add(&hp)

			if added {
				log.Printf("Added new hash: %v", hp)
			} else {
				log.Printf("Updating existing hash: %v", hp)
			}
		}
	}
}
