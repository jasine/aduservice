#!/bin/bash

if [ $# -ne 1 ]; then
	echo "./update.sh new_version[eg: 0.9.2]"
	exit
fi

TAG=$1

#go build dockermanager.go
./build4server.sh
#./build4armv7.sh

cp ../basic/auth .

docker build -t 192.168.5.46:5000/aduservice-t2-test:$TAG .

docker push 192.168.5.46:5000/aduservice-t2-test:$TAG

echo aduservice-t2-test:$TAG >> versions