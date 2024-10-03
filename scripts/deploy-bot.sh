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
IMAGE_NAME="global_terraplay_ecr"
CONTAINER_NAME="nimbus-bot"

# Write the SSH key to a file and set correct permissions
echo "${EC2_SSH_KEY}" > ec2_key.pem
chmod 600 ec2_key.pem

# Define the SSH command
SSH_CMD="ssh -o StrictHostKeyChecking=no -i ec2_key.pem ${EC2_USER}@${EC2_HOST}"

echo "Starting SSH connection to update and install Docker and AWS CLI on EC2"
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

    # Check if AWS CLI is installed and update it, otherwise install it
    if aws --version 2>&1 >/dev/null; then
        echo "AWS CLI is already installed, updating it"
        sudo ./aws/install --update || { echo "Failed to update AWS CLI"; exit 1; }
    else
        echo "Installing AWS CLI"
        sudo yum install -y unzip || { echo "Failed to install unzip"; exit 1; }
        curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" || { echo "Failed to download AWS CLI"; exit 1; }
        unzip awscliv2.zip || { echo "Failed to unzip AWS CLI"; exit 1; }
        sudo ./aws/install || { echo "Failed to install AWS CLI"; exit 1; }
        echo "Cleaning up AWS CLI installation files"
        rm -rf awscliv2.zip aws
    fi
EOF
echo "SSH connection complete. Docker and AWS CLI setup finished."

echo "Deploying the application on EC2"
$SSH_CMD << EOF
    echo "Configuring AWS CLI with passed credentials"
    aws configure set aws_access_key_id ${AWS_ACCESS_KEY_ID}
    aws configure set aws_secret_access_key ${AWS_SECRET_ACCESS_KEY}
    aws configure set region ${AWS_REGION}

    # Check if the container exists
    if [ \$(docker ps -a -q --filter name=${CONTAINER_NAME}) ]; then
        echo "Stopping and removing the existing container: ${CONTAINER_NAME}"
        docker stop ${CONTAINER_NAME} || true
        docker rm ${CONTAINER_NAME} || true
    else
        echo "No existing container named ${CONTAINER_NAME} found"
    fi

    echo "Logging in to AWS ECR"
    aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com || { echo "Failed to login to AWS ECR"; exit 1; }

    echo "Pulling Docker image from ECR"
    docker pull ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${IMAGE_NAME}:${IMAGE_TAG} || { echo "Failed to pull Docker image"; exit 1; }

    echo "Running the new Docker container"
    docker run -d --name ${CONTAINER_NAME} -p 8080:8080 \
      -e DISCORD_BOT_TOKEN=${DISCORD_BOT_TOKEN} \
      -e DISCORD_CHANNEL_ID=${DISCORD_CHANNEL_ID} \
      -e REPO_TOKEN=${REPO_TOKEN} \
      -e ACTIONS_WEBHOOK_SECRET=${ACTIONS_WEBHOOK_SECRET} \
      -e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
      -e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
      -e AWS_REGION=${AWS_REGION} \
      -e DYNAMO_TABLE=${DYNAMO_TABLE} \
      -e APP_ENV=${APP_ENV} \
      ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${IMAGE_NAME}:${IMAGE_TAG} || { echo "Failed to run Docker container"; exit 1; }
EOF

# Clean up the SSH key file after deployment
rm ec2_key.pem

echo "Deployment complete!"
