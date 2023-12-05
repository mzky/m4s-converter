@echo off
set GOARCH=386
go build -ldflags "-w -s"