@echo off
:: tags 支持linux、windows、macos
echo build for windows-amd64...
set GOARCH=amd64
go build -tags "windows" -ldflags "-w -s" -o m4s-converter-amd64.exe
:: 压缩
upx --lzma m4s-converter-amd64.exe
