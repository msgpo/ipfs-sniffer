#!/bin/sh

env GOOS=linux GOARCH=amd64 go build -o build/logsniffer.linux64 main.go && \
rsync -v --progress build/logsniffer.linux64 freeman.ipfs-search.com:~/bin/ && \
ssh freeman.ipfs-search.com ./bin/logsniffer.linux64
