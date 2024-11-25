#!/bin/bash
version="1.5.2"
sourceVer=$(git log --date=iso --pretty=format:"%h @%cd" -1)
buildTime="$(date '+%Y-%m-%d %H:%M:%S') by $(go version|sed 's/go version //')"
cat <<EOF | gofmt >common/version.go
package common
var (
    version = "$version"
    sourceVer = "$sourceVer"
    buildTime = "$buildTime"
)
EOF
# tags 支持linux、windows、macos
echo build for Linux-amd64...
GOOS=linux GOARCH=amd64 go build -tags "linux" -ldflags "-w -s" -o m4s-converter-linux_amd64 main.go
#
#echo build for Linux-arm64...
#GOOS=linux GOARCH=arm64 GOARM=7 go build -tags "linux" -ldflags "-w -s" -o m4s-converter-linux_arm64 main.go
# 设置目标操作系统为darwin（MacOS）
echo build for darwin-amd64...
GOOS=darwin GOARCH=amd64 go build -tags "darwin" -o m4s-converter-darwin_amd64 main.go
#
echo build for darwin-arm64...
GOOS=darwin GOARCH=arm64 go build -tags "darwin" -o m4s-converter-darwin_arm64 main.go
#
echo build for windows-amd64...
GOOS=windows GOARCH=amd64 go build -tags "windows" -ldflags "-w -s" -o m4s-converter-amd64.exe
# 压缩
upx --lzma m4s-converter-*