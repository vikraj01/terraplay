#!/bin/bash

# Update the system
sudo yum update -y

# Install Docker using Amazon Linux Extras
sudo amazon-linux-extras install docker -y

# Enable and start Docker service
sudo systemctl enable docker
sudo systemctl start docker

# Add ec2-user to the Docker group to avoid using sudo with Docker
sudo usermod -aG docker ec2-user

# Download Docker Compose and set executable permissions
sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Create directories for Minetest data
sudo mkdir -p /opt/minetest/data
sudo chown -R ec2-user:ec2-user /opt/minetest

# Create Docker Compose file for Minetest
cat <<EOF > /opt/minetest/docker-compose.yml
version: '3.8'

services:
  minetest:
    image: lscr.io/linuxserver/minetest:latest
    container_name: minetest
    environment:
      - PUID=1000
      - PGID=1000
      - TZ=Etc/UTC
      - "CLI_ARGS=--gameid devtest"
    volumes:
      - /opt/minetest/data:/config/.minetest
    ports:
      - 30000:30000/udp
    restart: unless-stopped
EOF

# Ensure the permissions for the Minetest directory are correct
sudo chown -R ec2-user:ec2-user /opt/minetest

# Run Docker Compose using sudo (since the ec2-user group change might not apply yet)
cd /opt/minetest
sudo /usr/local/bin/docker-compose up -d

# Print completion message
echo "Minetest server setup is complete!"
