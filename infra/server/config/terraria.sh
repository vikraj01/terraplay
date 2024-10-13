#!/bin/bash

sudo yum update -y

sudo amazon-linux-extras install docker -y

sudo systemctl enable docker
sudo systemctl start docker

sudo usermod -aG docker ec2-user

sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

sudo mkdir -p /opt/terraria/data
sudo chown -R ec2-user:ec2-user /opt/terraria

cat <<EOF >/opt/terraria/docker-compose.yml
version: '3.8'

services:
  terraria:
    image: ryshe/terraria:latest
    container_name: terraria
    environment:
      - TZ=Etc/UTC
      - world_name=MyWorld
      - maxplayers=16
      - difficulty=1
      - world_size=2
    volumes:
      - /opt/terraria/data:/world
    ports:
      - 7777:7777/tcp
    restart: unless-stopped
EOF

sudo chown -R ec2-user:ec2-user /opt/terraria

cd /opt/terraria || exit
sudo /usr/local/bin/docker-compose up -d

echo "Terraria server setup is complete!"
