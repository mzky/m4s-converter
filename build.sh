#!/bin/bash
echo build for Linux-amd64 ...
GOOS=linux GOARCH=amd64 go build -tags "linux" -ldflags "-w -s" -o m4s-converter-amd64 main.go

echo build for Linux-arm64 ...
GOOS=linux GOARCH=arm64 GOARM=8 go build -tags "linux" -ldflags "-w -s" -o m4s-converter-arm64 main.go



