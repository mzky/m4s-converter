@echo off
:: tags 支持linux、windows、macos
echo build for windows-x86...
set GOARCH=386
go build -tags "windows" -ldflags "-w -s" -o m4s-converter-386.exe
:: x86的exe文件兼容64位系统，实际都使用32位程序即可
echo build for windows-amd64...
set GOARCH=amd64
go build -tags "windows" -ldflags "-w -s" -o m4s-converter-amd64.exe
:: 压缩
upx --lzma m4s-converter-386.exe
upx --lzma m4s-converter-amd64.exe
::
echo build for windows-arm64...
set GOARCH=arm64
set GOARM=7
go build -tags "windows" -ldflags "-w -s" -o m4s-converter-arm64.exe main.go

