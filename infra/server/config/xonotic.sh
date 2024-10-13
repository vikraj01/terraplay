#!/bin/bash

sudo yum update -y

sudo amazon-linux-extras install docker -y

sudo systemctl enable docker
sudo systemctl start docker

sudo usermod -aG docker ec2-user

sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

sudo mkdir -p /opt/xonotic/data
sudo chown -R ec2-user:ec2-user /opt/xonotic

cat <<EOF >/opt/xonotic/docker-compose.yml
version: '3.8'

services:
  xonotic:
    image: xonotic/server
    container_name: xonotic
    volumes:
      - /opt/xonotic/data:/config
    ports:
      - 26000:26000/udp
    restart: unless-stopped
EOF

sudo chown -R ec2-user:ec2-user /opt/xonotic

cd /opt/xonotic || exit
sudo /usr/local/bin/docker-compose up -d

echo "Xonotic server setup is complete!"
