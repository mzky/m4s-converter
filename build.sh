#!/bin/bash
echo build for Linux-amd64 ...
GOOS=linux GOARCH=amd64 go build -tags "linux" -ldflags "-w -s" -o m4s-converter-amd64 main.go
#
echo build for Linux-arm64 ...
GOOS=linux GOARCH=arm64 GOARM=8 go build -tags "linux" -ldflags "-w -s" -o m4s-converter-arm64 main.go
# 压缩
#upx --lzma  m4s-converter-*

# 设置目标操作系统为darwin（MacOS）
echo build for darwin-amd64 ...
GOOS=darwin GOARCH=amd64 go build -tags "darwin" -o m4s-converter-darwin_amd64 main.go
echo build for darwin-arm64 ...
GOOS=darwin GOARCH=arm64 go build -tags "darwin" -o m4s-converter-darwin_arm64 main.go
