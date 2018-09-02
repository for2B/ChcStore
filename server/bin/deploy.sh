#!/bin/bash

#指令执行错误的操作
exitOnErr(){
  exitCode=$?
  echo "[LINE:$1] Error: Command  exited with status $?"
  exit $exitCode
}

#捕获错误信号，并执行相应操作
trap 'exitOnErr $LINENO' ERR

PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin:/tmp/workspace/go/bin
export PATH 
export chcs_PATH=/var/lib/jenkins/workspace/ChcStore/server
export GOPATH=$GOPATH:$chcs_PATH
export SRC_PATH=$chcs_PATH/src
export BIN_PATH=$chcs_PATH/bin

cd $SRC_PATH/ChenHC/app

echo "go build...."
go build -o chcs  main.go
if [  $? -eq 0 ];then
      echo " go build succeess!"
fi

APP_PID=$(lsof -i:6618 -t)
if [ ! -z "$APP_PID" ];then
     kill -9 $APP_PID
fi

if [ $? -eq 0 ];then
      echo " kill chcs succeess!"
else  
      echo " no running chcs"
fi

if [ -f "$BIN_PATH/chcs" ];then
	rm $BIN_PATH/chcs
    echo "delete chcs success!"
else 
	echo "no chcs"
fi

mv ./chcs $BIN_PATH
if [ $? -eq 0 ];then
    echo "move succeess!"
fi

nohup $BIN_PATH/chcs &
if [ $? -ne 0 ]; then
     exit 1
else
    echo "chcs run succeess !!!"
fi


