#!/bin/bash

# Exit immediately if any command exits with a non-zero status
set -e

# Function to handle errors and print custom messages
error_handler() {
    echo "Error occurred in script at line: ${1}"
    exit 1
}

# Trap errors and call the error_handler function with the failing line number
trap 'error_handler $LINENO' ERR

# Set variables
EC2_HOST=${EC2_HOST}
EC2_USER=${EC2_USER}
EC2_SSH_KEY=${EC2_SSH_KEY}
AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}
AWS_REGION=${AWS_REGION}
AWS_ACCOUNT_ID=${AWS_ACCOUNT_ID}
IMAGE_TAG=${IMAGE_TAG}
DISCORD_BOT_TOKEN=${DISCORD_BOT_TOKEN}
DISCORD_CHANNEL_ID=${DISCORD_CHANNEL_ID}
REPO_TOKEN=${REPO_TOKEN}
ACTIONS_WEBHOOK_SECRET=${ACTIONS_WEBHOOK_SECRET}
DYNAMO_TABLE=${DYNAMO_TABLE}
APP_ENV=${APP_ENV}

# Define the SSH command
SSH_CMD="ssh -o StrictHostKeyChecking=no -i ${EC2_SSH_KEY} ${EC2_USER}@${EC2_HOST}"

echo "Starting SSH connection to update and install Docker on EC2"
$SSH_CMD << EOF
    echo "Updating EC2 instance"
    sudo yum update -y || { echo "Failed to update EC2 instance"; exit 1; }

    echo "Installing Docker"
    sudo amazon-linux-extras install docker -y || { echo "Failed to install Docker"; exit 1; }

    echo "Starting Docker service"
    sudo systemctl start docker || { echo "Failed to start Docker"; exit 1; }

    echo "Enabling Docker service to run on boot"
    sudo systemctl enable docker || { echo "Failed to enable Docker service"; exit 1; }

    echo "Adding user to Docker group"
    sudo usermod -aG docker ec2-user || { echo "Failed to add user to Docker group"; exit 1; }

    echo "Restarting Docker service"
    sudo systemctl restart docker || { echo "Failed to restart Docker service"; exit 1; }
EOF
echo "SSH connection complete. Docker setup finished."

echo "Deploying the application on EC2"
$SSH_CMD << EOF
    echo "Logging in to AWS ECR"
    aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com || { echo "Failed to login to AWS ECR"; exit 1; }

    echo "Pulling Docker image from ECR"
    docker pull ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/nimbus-bot:${IMAGE_TAG} || { echo "Failed to pull Docker image"; exit 1; }

    echo "Stopping existing container if running"
    docker stop nimbus-bot || true

    echo "Removing existing container if present"
    docker rm nimbus-bot || true

    echo "Running the new Docker container"
    docker run -d --name nimbus-bot -p 80:8000 \
      -e DISCORD_BOT_TOKEN=${DISCORD_BOT_TOKEN} \
      -e DISCORD_CHANNEL_ID=${DISCORD_CHANNEL_ID} \
      -e REPO_TOKEN=${REPO_TOKEN} \
      -e ACTIONS_WEBHOOK_SECRET=${ACTIONS_WEBHOOK_SECRET} \
      -e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
      -e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
      -e AWS_REGION=${AWS_REGION} \
      -e DYNAMO_TABLE=${DYNAMO_TABLE} \
      -e APP_ENV=${APP_ENV} \
      ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/nimbus-bot:${IMAGE_TAG} || { echo "Failed to run Docker container"; exit 1; }
EOF
echo "Deployment complete!"
