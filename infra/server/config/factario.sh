#!/bin/bash

sudo yum update -y

sudo amazon-linux-extras install docker -y

sudo systemctl enable docker
sudo systemctl start docker

sudo usermod -aG docker ec2-user

sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

sudo mkdir -p /opt/factorio/data
sudo chown -R ec2-user:ec2-user /opt/factorio

cat <<EOF >/opt/factorio/docker-compose.yml
version: '3.8'

services:
  factorio:
    image: factoriotools/factorio:latest
    container_name: factorio
    volumes:
      - /opt/factorio/data:/factorio
    ports:
      - 34197:34197/udp
    restart: unless-stopped
EOF

sudo chown -R ec2-user:ec2-user /opt/factorio

cd /opt/factorio || exit
sudo /usr/local/bin/docker-compose up -d

echo "Factorio server setup is complete!"
