#!/bin/bash

echo "Start ssh agent"
eval `ssh-agent`
ssh-add


echo "Building"
go build -ldflags="-s -w"

echo "Stop service"
ssh root@chat.pipeto.me 'systemctl stop curlchat'

echo "Copy executable file"
scp ./curlchat root@chat.pipeto.me:/root/data/curlchat/curlchat

echo "Start service"
ssh root@chat.pipeto.me 'systemctl start curlchat'


echo "Stop ssh agent"
killall ssh-agent
