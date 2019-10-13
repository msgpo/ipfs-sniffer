package logsniffer

import (
	"context"
	shell "github.com/ipfs/go-ipfs-api"
	"log"
)

// Reader provides a reader with a channel of IPFS log messages.
type Reader struct {
	Errors   chan error
	Messages chan Message
	logger   shell.Logger
	ctx      context.Context
}

// Open a logger for the shell.
func (l *Reader) Open(ctx context.Context, sh *shell.Shell) error {
	logger, err := sh.GetLogs(ctx)

	l.logger = logger
	l.ctx = ctx

	return err
}

// Close the logger.
func (l *Reader) Close() error {
	return l.logger.Close()
}

// Read messages from log channel until context closes, writing errors to the error channel.
func (l *Reader) Read() {
	for {
		msg, err := l.logger.Next()
		if err != nil {
			l.Errors <- err
		}

		select {
		case <-l.ctx.Done():
			// Context closed
			log.Printf("Stop reading log messages: %s", l.ctx.Err())
			l.Errors <- l.ctx.Err()
		case l.Messages <- msg:
		}
	}
}
