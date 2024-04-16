@echo off
set GOARCH=386
go build -ldflags "-w -s"
upx --lzma m4s-converter.exe