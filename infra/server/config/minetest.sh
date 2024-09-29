#!/bin/bash

sudo yum update -y

sudo amazon-linux-extras install docker -y

sudo systemctl enable docker
sudo systemctl start docker

sudo usermod -aG docker ec2-user

sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

sudo mkdir -p /opt/minetest/data
sudo chown -R ec2-user:ec2-user /opt/minetest

cat <<EOF > /opt/minetest/docker-compose.yml
version: '3.8'

services:
  minetest:
    image: lscr.io/linuxserver/minetest:latest
    container_name: minetest
    environment:
      - PUID=1000
      - PGID=1000
      - TZ=Etc/UTC
      - "CLI_ARGS=--gameid devtest"
    volumes:
      - /opt/minetest/data:/config/.minetest
    ports:
      - 30000:30000/udp
    restart: unless-stopped
EOF

sudo chown -R ec2-user:ec2-user /opt/minetest

cd /opt/minetest
sudo /usr/local/bin/docker-compose up -d

echo "Minetest server setup is complete!"
