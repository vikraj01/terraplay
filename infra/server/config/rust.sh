#!/bin/bash

sudo yum update -y

sudo amazon-linux-extras install docker -y

sudo systemctl enable docker
sudo systemctl start docker

sudo usermod -aG docker ec2-user

sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

sudo mkdir -p /opt/rust/data
sudo chown -R ec2-user:ec2-user /opt/rust

cat <<EOF >/opt/rust/docker-compose.yml
version: '3.8'

services:
  rust:
    image: didstopia/rust-server:latest
    container_name: rust
    environment:
      - RUST_SERVER_STARTUP_ARGUMENTS=-batchmode
      - RUST_SERVER_NAME="My Rust Server"
    volumes:
      - /opt/rust/data:/steamcmd/rust
    ports:
      - 28015:28015/tcp
      - 28015:28015/udp
    restart: unless-stopped
EOF

sudo chown -R ec2-user:ec2-user /opt/rust

cd /opt/rust || exit
sudo /usr/local/bin/docker-compose up -d

echo "Rust server setup is complete!"
