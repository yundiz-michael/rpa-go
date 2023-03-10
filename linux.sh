#!/usr/bin/env bash
# 编译指令，会加载当前.mod文件对应的包进行编译，所以必须到main所在的目录
#rm -f /workspace/release/linux/server
cd src/main/
go env -w GOARCH=amd64
go env -w CGO_ENABLED=1
go env -w GOOS=linux
go env -w GOPROXY=https://goproxy.cn,direct
go build -o  /workspace/release/linux/merkaba /workspace/xpa/go/merkaba/src/main/main.go
yes | cp /workspace/xpa/go/merkaba/libs/libvncserver.so* /workspace/release/linux/.
fileDate=`date "+%m_%d_%H_%M"`
fileName="/workspace/release/linux_"+$fileDate
cd /workspace/release/linux
zip -r $fileName *
echo $fileName

