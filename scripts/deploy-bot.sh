#!/bin/bash

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

echo "Starting SSH connection to update and install Docker on EC2"
SSH_CMD="ssh -o StrictHostKeyChecking=no -i ${EC2_SSH_KEY} ${EC2_USER}@${EC2_HOST}"

$SSH_CMD << EOF
    echo "Updating EC2 instance"
    sudo yum update -y

    echo "Installing Docker"
    sudo amazon-linux-extras install docker -y

    echo "Starting Docker service"
    sudo systemctl start docker

    echo "Enabling Docker service to run on boot"
    sudo systemctl enable docker

    echo "Adding user to Docker group"
    sudo usermod -aG docker ec2-user

    echo "Restarting Docker service"
    sudo systemctl restart docker
EOF

echo "SSH connection complete. Docker setup finished."

echo "Deploying the application on EC2"
$SSH_CMD << EOF
    echo "Logging in to AWS ECR"
    aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com

    echo "Pulling Docker image from ECR"
    docker pull ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/nimbus-bot:${IMAGE_TAG}

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
      ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/nimbus-bot:${IMAGE_TAG}
EOF

echo "Deployment complete!"
