#!/bin/bash

sudo yum update -y

sudo amazon-linux-extras install docker -y

sudo systemctl enable docker
sudo systemctl start docker

sudo usermod -aG docker ec2-user

sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

sudo mkdir -p /opt/valheim/data
sudo chown -R ec2-user:ec2-user /opt/valheim

cat <<EOF >/opt/valheim/docker-compose.yml
version: '3.8'

services:
  valheim:
    image: lloesche/valheim-server:latest
    container_name: valheim
    environment:
      - SERVER_NAME=MyValheimServer
      - WORLD_NAME=MyValheimWorld
      - SERVER_PASSWORD=secret
      - PUID=1000
      - PGID=1000
    volumes:
      - /opt/valheim/data:/config
    ports:
      - 2456-2458:2456-2458/udp
    restart: unless-stopped
EOF

sudo chown -R ec2-user:ec2-user /opt/valheim

cd /opt/valheim || exit
sudo /usr/local/bin/docker-compose up -d

echo "Valheim server setup is complete!"
