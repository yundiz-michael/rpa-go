#!/usr/bin/env bash
# 链接动态库到 /usr/local/lib
ln -s /workspace/xpa/go/merkaba/libs/libvncserver.dylib  /usr/local/lib/libvncserver.1.dylib
echo "链接动态库到/usr/local/lib,成功"