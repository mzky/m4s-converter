@echo off
:: tags 支持linux、windows、macos
set "targetFile=common\version.go"
echo package common> %targetFile%
echo var Version = "1.3.7">> %targetFile%
::
echo build for windows success
set GOARCH=386
go build -tags "windows" -ldflags "-w -s" -o m4s-converter-386.exe main.go
upx --lzma m4s-converter-386.exe
