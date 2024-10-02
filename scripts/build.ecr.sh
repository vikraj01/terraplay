#!/bin/bash

if [ -z "$AWS_REGION" ] || [ -z "$AWS_ACCOUNT_ID" ]; then
  echo "Error: Missing required environment variables."
  echo "AWS_REGION and AWS_ACCOUNT_ID must be set."
  exit 1
fi

ECR_URL="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com"
IMAGE_NAME="global_terraplay_ecr"
IMAGE_TAG=$(date +"%Y%m%d-%H%M")

echo "Logging in to Amazon ECR at $ECR_URL..."
aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${ECR_URL}
if [ $? -ne 0 ]; then
  echo "Error: Failed to log in to Amazon ECR"
  exit 1
fi

echo "Building Docker image ${IMAGE_NAME}:${IMAGE_TAG}..."
docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .
if [ $? -ne 0 ]; then
  echo "Error: Failed to build Docker image"
  exit 1
fi

echo "Tagging Docker image as ${ECR_URL}/${IMAGE_NAME}:${IMAGE_TAG}..."
docker tag ${IMAGE_NAME}:${IMAGE_TAG} ${ECR_URL}/${IMAGE_NAME}:${IMAGE_TAG}
if [ $? -ne 0 ]; then
  echo "Error: Failed to tag Docker image"
  exit 1
fi

echo "Pushing Docker image to ${ECR_URL}..."
docker push ${ECR_URL}/${IMAGE_NAME}:${IMAGE_TAG}
if [ $? -ne 0 ]; then
  echo "Error: Failed to push Docker image"
  exit 1
fi

echo "Docker image ${IMAGE_NAME}:${IMAGE_TAG} successfully pushed to ${ECR_URL}!"
