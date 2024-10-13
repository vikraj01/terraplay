#!/bin/bash

sudo yum update -y

sudo amazon-linux-extras install docker -y

sudo systemctl enable docker
sudo systemctl start docker

sudo usermod -aG docker ec2-user

sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

sudo mkdir -p /opt/wesnoth/data
sudo chown -R ec2-user:ec2-user /opt/wesnoth

cat <<EOF >/opt/wesnoth/docker-compose.yml
version: '3.8'

services:
  wesnoth:
    image: wesnoth/wesnoth
    container_name: wesnoth
    volumes:
      - /opt/wesnoth/data:/config
    ports:
      - 15000:15000/tcp
    restart: unless-stopped
EOF

sudo chown -R ec2-user:ec2-user /opt/wesnoth

cd /opt/wesnoth || exit
sudo /usr/local/bin/docker-compose up -d

echo "Wesnoth server setup is complete!"
