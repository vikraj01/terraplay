#!/bin/bash

sudo yum update -y

sudo amazon-linux-extras install docker -y

sudo systemctl enable docker
sudo systemctl start docker

sudo usermod -aG docker ec2-user

sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

sudo mkdir -p /opt/veloren/data
sudo chown -R ec2-user:ec2-user /opt/veloren

cat <<EOF >/opt/veloren/docker-compose.yml
version: '3.8'

services:
  veloren:
    image: veloren/veloren-server
    container_name: veloren
    volumes:
      - /opt/veloren/data:/config
    ports:
      - 14004:14004/tcp
    restart: unless-stopped
EOF

sudo chown -R ec2-user:ec2-user /opt/veloren

cd /opt/veloren || exit
sudo /usr/local/bin/docker-compose up -d

echo "Veloren server setup is complete!"
