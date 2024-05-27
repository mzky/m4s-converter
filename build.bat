@echo off
set "targetFile=common\version.go"
echo package common> %targetFile%
echo var Version = "1.3.6">> %targetFile%
set GOARCH=386
go build -ldflags "-w -s"
upx --lzma m4s-converter.exe