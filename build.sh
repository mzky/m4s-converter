#!/bin/bash
# 交叉编译arm版本,需使用指定版本的gcc 下载地址：https://releases.linaro.org/components/toolchain/binaries/
gccPath="/usr/gcc-aarch64/bin/aarch64-linux-gnu-gcc" # 从外部定义gcc路径
CGO_ENABLED=1 CC=${gccPath} GOOS=linux GOARCH=arm64 GOARM=7 go build -ldflags "-w -s" -o m4s-converter-arm64 main.go
echo build for ARM64 success