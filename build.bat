@echo off
set GOARCH=386
go build -ldflags "-w -s"
upx --brute m4s-converter.exe