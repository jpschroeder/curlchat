#!/bin/bash

echo "Start ssh agent"
eval `ssh-agent`
ssh-add

## Building

echo "Building"
go build -ldflags="-s -w"

## Deployment

echo "Create destination folder"
ssh root@chat.pipeto.me 'mkdir /root/data/curlchat'

echo "Copy executable file"
scp ./curlchat root@chat.pipeto.me:/root/data/curlchat/curlchat

echo "Grant execute permission"
ssh root@chat.pipeto.me 'chmod +x /root/data/curlchat/curlchat'

### Service

echo "Copy service definition"
scp ./scripts/curlchat.service root@chat.pipeto.me:/lib/systemd/system/curlchat.service

echo "Start service"
ssh root@chat.pipeto.me 'systemctl start curlchat'

echo "Enable service on boot"
ssh root@chat.pipeto.me 'systemctl enable curlchat'

echo "Check service status"
ssh root@chat.pipeto.me 'systemctl status curlchat'

### Nginx

echo "Make sure that nginx is installed"
ssh root@chat.pipeto.me 'apt-get install -y nginx'

echo "Copy nginx config"
scp ./scripts/curlchat.nginx.conf root@chat.pipeto.me:/etc/nginx/sites-available/curlchat.nginx.conf

echo "Enable the nginx proxy"
ssh root@chat.pipeto.me 'ln -s /etc/nginx/sites-available/curlchat.nginx.conf /etc/nginx/sites-enabled/curlchat.nginx.conf'

echo "Restart nginx to pick up the changes"
ssh root@chat.pipeto.me 'systemctl restart nginx'

### Nginx Https

echo "Install letsencrypt client"
ssh root@chat.pipeto.me 'add-apt-repository ppa:certbot/certbot'
ssh root@chat.pipeto.me 'apt-get install python-certbot-nginx'

echo "Generate and install certificate"
ssh root@chat.pipeto.me 'certbot --nginx -d chat.pipeto.me'


echo "Stop ssh agent"
killall ssh-agent
