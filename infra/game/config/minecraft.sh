#!/bin/bash

# Redirect all output to /var/log/user-data.log for troubleshooting
exec > >(tee /var/log/user-data.log | logger -t user-data) 2>&1

# Update and upgrade system packages
sudo apt update -y && sudo apt upgrade -y

# Install necessary packages for Docker and Nginx
sudo apt-get install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common

# Add Docker's official GPG key
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

# Add Docker APT repository
sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"

# Update package lists again to include Docker's repository
sudo apt update -y

# Install Docker and necessary components
sudo apt-get install -y docker-ce docker-ce-cli containerd.io

# Enable and start Docker
sudo systemctl enable docker
sudo systemctl start docker

# Verify Docker installation
docker --version

# Pull the Minecraft server Docker image
docker pull itzg/minecraft-server

# Create a directory to persist Minecraft data
sudo mkdir -p /data
sudo chown ubuntu:ubuntu /data

# Run the Minecraft server container with environment variables and volume
docker run -d \
  --name mc-server \
  -p 25565:25565 \
  -e EULA=TRUE \
  -e ONLINE_MODE=TRUE \  # Setting ONLINE_MODE to TRUE for secure, online mode
  -v /data:/data \
  itzg/minecraft-server

# Verify Minecraft server is running
docker ps

# Install Nginx
sudo apt-get install -y nginx

# Enable and start Nginx
sudo systemctl enable nginx
sudo systemctl start nginx

# Retrieve EC2 metadata
export META_INST_ID=$(curl -s http://169.254.169.254/latest/meta-data/instance-id)
export META_INST_TYPE=$(curl -s http://169.254.169.254/latest/meta-data/instance-type)
export META_INST_AZ=$(curl -s http://169.254.169.254/latest/meta-data/placement/availability-zone)

# Create a simple HTML status page for the EC2 instance
sudo tee /var/www/html/index.html > /dev/null <<EOF
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: 'Open Sans', sans-serif;
            background-color: #f4f4f4;
            color: #333;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
        }
        .instance-card {
            background-color: white;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0 0 15px rgba(0,0,0,0.1);
            text-align: center;
        }
        h1 {
            color: #6944ff;
        }
        p {
            margin: 10px 0;
            font-size: 18px;
        }
        .data {
            font-weight: bold;
        }
    </style>
    <title>EC2 Instance Status</title>
</head>
<body>
    <div class="instance-card">
        <h1>EC2 Instance Information</h1>
        <p>Instance ID: <span class="data">$META_INST_ID</span></p>
        <p>Instance Type: <span class="data">$META_INST_TYPE</span></p>
        <p>Availability Zone: <span class="data">$META_INST_AZ</span></p>
    </div>
</body>
</html>
EOF

# Ensure Nginx is running properly
sudo systemctl restart nginx

# Check Nginx status to ensure it's running
systemctl status nginx
