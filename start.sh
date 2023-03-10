#!/usr/bin/env bash
# 如何创建一个编号为1的2k屏幕（4k=3840x2160x32）
# sudo apt install xvfb
hasx1=`ls /tmp/.X11-unix|grep -w X1|wc -l`
if [ 0 == $hasx1 ];then
   sudo Xvfb :1 -ac -screen 0 2560x1440x24 &
   echo "创建屏幕X1"
else
   echo "屏幕X1已经存在"
fi
# 查看系统有多少个屏幕及编号
# ls /tmp/.X11-unix
cd /workspace/merkaba/
rm -f nohup.out
rm -f log/*
ps aux | grep server | awk '{print $2}' | xargs kill
ps aux | grep chrome | awk '{print $2}' | xargs kill
count=`ps -ef |grep server|grep -v "grep"|wc -l`
if [ 0 == $count ];then
   echo "merkaba 已经停止"
fi
DISPLAY=:1 nohup ./merkaba  >&1 &
sleep 0.5
tail -30f nohup.out