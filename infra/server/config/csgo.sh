#!/bin/bash

sudo yum update -y

sudo amazon-linux-extras install docker -y

sudo systemctl enable docker
sudo systemctl start docker

sudo usermod -aG docker ec2-user

sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

sudo mkdir -p /opt/csgo/data
sudo chown -R ec2-user:ec2-user /opt/csgo

cat <<EOF >/opt/csgo/docker-compose.yml
version: '3.8'

services:
  csgo:
    image: cm2network/csgo:latest
    container_name: csgo
    environment:
      - SRCDS_TOKEN=your_token_here
      - SRCDS_PORT=27015
      - SRCDS_TICKRATE=128
    volumes:
      - /opt/csgo/data:/home/steam/csgo
    ports:
      - 27015:27015/udp
      - 27015:27015/tcp
    restart: unless-stopped
EOF

sudo chown -R ec2-user:ec2-user /opt/csgo

cd /opt/csgo || exit
sudo /usr/local/bin/docker-compose up -d

echo "CS:GO server setup is complete!"
