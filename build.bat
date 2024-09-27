@echo off
:: tags 支持linux、windows、macos
set "targetFile=common\version.go"
(
    echo package common
    echo.
    echo var Version = "1.3.7"
) > "%targetFile%"
::
echo build for windows ...
set GOARCH=386
go build -tags "windows" -ldflags "-w -s" -o m4s-converter-386.exe main.go
:: x86的exe文件兼容64位系统
set GOARCH=amd64
go build -tags "windows" -ldflags "-w -s" -o m4s-converter-amd64.exe main.go
:: 压缩
upx --lzma m4s-converter-*.exe
