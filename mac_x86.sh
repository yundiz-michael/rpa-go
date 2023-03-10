#!/usr/bin/env bash
# 编译指令，会加载当前.mod文件对应的包进行编译，所以必须到main所在的目录
rm -rf /workspace/release/server
cd src/main/
go env -w GOARCH=amd64
go env -w CGO_ENABLED=1
go env -w GOOS=darwin
go env -w GOPROXY=https://goproxy.cn,direct
go build -o /workspace/release/mac_x86/merkaba /workspace/xpa/go/merkaba/src/main/main.go
yes | cp /workspace/xpa/go/merkaba/libs/libvncserver.dylib.x86 /workspace/release/mac_x86/libvncserver.dylib
fileDate=`date "+%m_%d_%H_%M"`
fileName="/workspace/release/mac_x86_"+$fileDate
cd /workspace/release/mac_x86
zip -r $fileName *
echo $fileName
