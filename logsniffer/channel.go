package logsniffer

import (
	"context"
	shell "github.com/ipfs/go-ipfs-api"
)

// Reader provides a reader with a channel of IPFS log messages.
type Reader struct {
	Errors   chan error
	Messages chan Message
	logger   shell.Logger
	ctx      context.Context
}

// Open a logger for the shell.
func (r *Reader) Open(ctx context.Context, sh *shell.Shell) error {
	logger, err := sh.GetLogs(ctx)

	r.logger = logger
	r.ctx = ctx

	return err
}

// Close the logger.
func (r *Reader) Close() error {
	return r.logger.Close()
}

// Read messages from log channel until context closes, writing errors to the error channel.
func (r *Reader) Read() {
	for {
		msg, err := r.logger.Next()
		if err != nil {
			r.Errors <- err
		}

		select {
		case <-r.ctx.Done():
			// Context closed
			r.Errors <- r.ctx.Err()
		case r.Messages <- msg:
		}
	}
}
