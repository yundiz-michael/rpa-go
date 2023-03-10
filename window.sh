#!/usr/bin/env bash
# 编译指令，会加载当前.mod文件对应的包进行编译，所以必须到main所在的目录
rm -rf /mnt/hgfs/release/merkaba
cd src/main/
go env -w GOARCH=amd64
go env -w CGO_ENABLED=0
go env -w GOOS=window
go env -w GOPROXY=https://goproxy.cn,direct
go build -o  /mnt/hgfs/release/merkaba /mnt/hgfs/merkaba/src/main/main.go
