package logsniffer

import (
	"time"
)

// Actual log message for handleAddProvider
var handleAddProviderMessage = Message{
	"Logs":         []string{},
	"Operation":    "handleAddProvider",
	"ParentSpanID": 0,
	"SpanID":       6.999711555735423e+18,
	"Start":        time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	"Tags": map[string]string{
		"key":    "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
		"peer":   "QmeTtFXm42Jb2todcKR538j6qHYxXt6suUzpF3rtT9FPSd",
		"system": "dht",
	},
	"TraceID": 4.483443946463055e+18,
}
