#!/bin/bash

sudo yum update -y

sudo amazon-linux-extras install docker -y

sudo systemctl enable docker
sudo systemctl start docker

sudo usermod -aG docker ec2-user

sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

sudo mkdir -p /opt/minecraft/data
sudo chown -R ec2-user:ec2-user /opt/minecraft

cat <<EOF >/opt/minecraft/docker-compose.yml
version: '3.8'

services:
  minecraft:
    image: itzg/minecraft-server:latest
    container_name: minecraft
    environment:
      - EULA=TRUE
      - VERSION=LATEST
      - PUID=1000
      - PGID=1000
      - MEMORY=2G
    volumes:
      - /opt/minecraft/data:/data
    ports:
      - 25565:25565/tcp
    restart: unless-stopped
EOF

sudo chown -R ec2-user:ec2-user /opt/minecraft

cd /opt/minecraft || exit
sudo /usr/local/bin/docker-compose up -d

echo "Minecraft server setup is complete!"
