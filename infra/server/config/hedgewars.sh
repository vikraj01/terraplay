#!/bin/bash

sudo yum update -y

sudo amazon-linux-extras install docker -y

sudo systemctl enable docker
sudo systemctl start docker

sudo usermod -aG docker ec2-user

sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

sudo mkdir -p /opt/hedgewars/data
sudo chown -R ec2-user:ec2-user /opt/hedgewars

cat <<EOF >/opt/hedgewars/docker-compose.yml
version: '3.8'

services:
  hedgewars:
    image: hedgewars/server
    container_name: hedgewars
    volumes:
      - /opt/hedgewars/data:/config
    ports:
      - 46631:46631/tcp
    restart: unless-stopped
EOF

sudo chown -R ec2-user:ec2-user /opt/hedgewars

cd /opt/hedgewars || exit
sudo /usr/local/bin/docker-compose up -d

echo "Hedgewars server setup is complete!"
