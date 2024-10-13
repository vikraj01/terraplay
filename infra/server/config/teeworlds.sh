#!/bin/bash

sudo yum update -y

sudo amazon-linux-extras install docker -y

sudo systemctl enable docker
sudo systemctl start docker

sudo usermod -aG docker ec2-user

sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

sudo mkdir -p /opt/teeworlds/data
sudo chown -R ec2-user:ec2-user /opt/teeworlds

cat <<EOF >/opt/teeworlds/docker-compose.yml
version: '3.8'

services:
  teeworlds:
    image: teeworlds/teeworlds
    container_name: teeworlds
    volumes:
      - /opt/teeworlds/data:/config
    ports:
      - 8303:8303/udp
    restart: unless-stopped
EOF

sudo chown -R ec2-user:ec2-user /opt/teeworlds

cd /opt/teeworlds || exit
sudo /usr/local/bin/docker-compose up -d

echo "Teeworlds server setup is complete!"
