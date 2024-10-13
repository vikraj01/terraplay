#!/bin/bash

sudo yum update -y

sudo amazon-linux-extras install docker -y

sudo systemctl enable docker
sudo systemctl start docker

sudo usermod -aG docker ec2-user

sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

sudo mkdir -p /opt/unturned/data
sudo chown -R ec2-user:ec2-user /opt/unturned

cat <<EOF >/opt/unturned/docker-compose.yml
version: '3.8'

services:
  unturned:
    image: ghcr.io/linuxserver/unturned:latest
    container_name: unturned
    environment:
      - PUID=1000
      - PGID=1000
      - TZ=Etc/UTC
    volumes:
      - /opt/unturned/data:/config
    ports:
      - 27015:27015/udp
      - 27016:27016/tcp
    restart: unless-stopped
EOF

sudo chown -R ec2-user:ec2-user /opt/unturned

cd /opt/unturned || exit
sudo /usr/local/bin/docker-compose up -d

echo "Unturned server setup is complete!"
