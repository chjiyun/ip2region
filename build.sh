#!/bin/bash

appName="ip2region"
logDir="./info.log"
outputFile="./${appName}"
dbFileName="ip2region.db"
dbSource="https://fastly.jsdelivr.net/gh/bqf9979/ip2region@master/data/ip2region.db"

if [ -f "$dbFileName" ]; then
  cp "$dbFileName" "${dbFileName}.old"
fi

#download newest db file, otherwise do nothing
wget -N $dbSource

buildResult=$(go build -o "${outputFile}" main.go)

# 编译成功才能杀旧进程
if [ $? -eq 0 ]; then
  chmod 773 "${outputFile}"
  echo "build success, filename: ${outputFile}"
  pid=$(ps -ef |grep "${appName}" | grep -v grep|awk '{print $2}')
  echo "current pid is $pid"
  if [ -n "$pid" ]; then
      echo "Prepare to kill the process: ${pid}"
      kill -9 "$pid"
      sleep 1
  fi
else
  echo "build error $buildResult"
  exit
fi

nohup "${outputFile}" 1>"${logDir}" 2>&1 &
echo "complete..."
