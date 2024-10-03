#!/bin/bash

# This script assumes that the following environment variables are passed:
# - EC2_HOST
# - EC2_USER
# - EC2_SSH_KEY (the private key for SSH access)
# - AWS_ACCESS_KEY_ID
# - AWS_SECRET_ACCESS_KEY
# - AWS_REGION
# - AWS_ACCOUNT_ID (your AWS account ID for ECR)
# - IMAGE_TAG (the Docker image tag to pull)
# - DISCORD_BOT_TOKEN
# - DISCORD_CHANNEL_ID
# - REPO_TOKEN
# - ACTIONS_WEBHOOK_SECRET
# - DYNAMO_TABLE
# - APP_ENV

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

# Update the EC2 instance and install Docker
$SSH_CMD << EOF
    sudo yum update -y
    sudo amazon-linux-extras install docker -y
    sudo systemctl start docker
    sudo systemctl enable docker
    sudo usermod -aG docker ec2-user
    sudo systemctl restart docker
EOF

# Deploy the application
$SSH_CMD << EOF
    # Log in to AWS ECR securely
    aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com

    # Pull the Docker image with the correct tag
    docker pull ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/nimbus-bot:${IMAGE_TAG}

    # Stop and remove the existing container, if it's running
    docker stop nimbus-bot || true
    docker rm nimbus-bot || true

    # Run the new Docker container with environment variables passed in
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
      ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/nimbus-bot:${IMAGE_TAG}
EOF

echo "Deployment complete!"
