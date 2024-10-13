#!/bin/bash

sudo yum update -y

sudo amazon-linux-extras install docker -y

sudo systemctl enable docker
sudo systemctl start docker

sudo usermod -aG docker ec2-user

sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

sudo mkdir -p /opt/openttd/data
sudo chown -R ec2-user:ec2-user /opt/openttd

cat <<EOF >/opt/openttd/docker-compose.yml
version: '3.8'

services:
  openttd:
    image: bateau/openttd
    container_name: openttd
    volumes:
      - /opt/openttd/data:/config
    ports:
      - 3979:3979/tcp
    restart: unless-stopped
EOF

sudo chown -R ec2-user:ec2-user /opt/openttd

cd /opt/openttd || exit
sudo /usr/local/bin/docker-compose up -d

echo "OpenTTD server setup is complete!"
