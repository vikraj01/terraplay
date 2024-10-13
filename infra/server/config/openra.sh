#!/bin/bash

sudo yum update -y

sudo amazon-linux-extras install docker -y

sudo systemctl enable docker
sudo systemctl start docker

sudo usermod -aG docker ec2-user

sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

sudo mkdir -p /opt/openra/data
sudo chown -R ec2-user:ec2-user /opt/openra

cat <<EOF >/opt/openra/docker-compose.yml
version: '3.8'

services:
  openra:
    image: openra/server
    container_name: openra
    volumes:
      - /opt/openra/data:/config
    ports:
      - 1234:1234/tcp
    restart: unless-stopped
EOF

sudo chown -R ec2-user:ec2-user /opt/openra

cd /opt/openra || exit
sudo /usr/local/bin/docker-compose up -d

echo "OpenRA server setup is complete!"
