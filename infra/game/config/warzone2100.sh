#!/bin/bash

sudo yum update -y

sudo amazon-linux-extras install epel -y

sudo yum install -y warzone2100

sudo firewall-cmd --zone=public --add-port=9990/tcp --permanent
sudo firewall-cmd --reload

nohup warzone2100 --server > /var/log/warzone.log 2>&1 &

sleep 5  # Give it a few seconds to start
if pgrep -x "warzone2100" > /dev/null; then
    echo "Warzone 2100 server started successfully."
else
    echo "Failed to start Warzone 2100 server. Check logs at /var/log/warzone.log"
fi
