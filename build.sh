#!/bin/bash
version="1.3.9"
sourceVer=$(git log --date=iso --pretty=format:"%h @%cd" -1)
buildTime="$(date '+%Y-%m-%d %H:%M:%S') by $(go version|sed 's/go version //')"
cat <<EOF | gofmt >common/version.go
package common
var (
    Version = "$version"
    SourceVer = "$sourceVer"
    BuildTime = "$buildTime"
)
EOF
# tags 支持linux、windows、macos
echo build for Linux-amd64...
GOOS=linux GOARCH=amd64 go build -tags "linux" -ldflags "-w -s" -o m4s-converter-linux_amd64 main.go
echo build for Linux-arm64...
GOOS=linux GOARCH=arm64 GOARM=7 go build -tags "linux" -ldflags "-w -s" -o m4s-converter-linux_arm64 main.go
# 设置目标操作系统为darwin（MacOS）
echo build for darwin-amd64...
GOOS=darwin GOARCH=amd64 go build -tags "darwin" -o m4s-converter-darwin_amd64 main.go
echo build for darwin-arm64...
GOOS=darwin GOARCH=arm64 go build -tags "darwin" -o m4s-converter-darwin_arm64 main.go
# 压缩
upx --lzma m4s-converter-*