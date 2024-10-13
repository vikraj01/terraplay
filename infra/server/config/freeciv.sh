#!/bin/bash

sudo yum update -y

sudo amazon-linux-extras install docker -y

sudo systemctl enable docker
sudo systemctl start docker

sudo usermod -aG docker ec2-user

sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

sudo mkdir -p /opt/freeciv/data
sudo chown -R ec2-user:ec2-user /opt/freeciv

cat <<EOF >/opt/freeciv/docker-compose.yml
version: '3.8'

services:
  freeciv:
    image: chedy007/freeciv:1.0
    container_name: freeciv
    volumes:
      - /opt/freeciv/data:/config
    ports:
      - 5556:5556/tcp
    restart: unless-stopped
EOF

sudo chown -R ec2-user:ec2-user /opt/freeciv

cd /opt/freeciv || exit
sudo /usr/local/bin/docker-compose up -d

echo "Freeciv server setup is complete!"
