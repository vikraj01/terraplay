#!/bin/bash

IMAGE_NAME="nimbus-bot"
IMAGE_TAG="latest"
ECR_URL="009160067122.dkr.ecr.ap-south-1.amazonaws.com/global_terraplay_ecr"
ENV_FILE=".env"
AWS_REGION="ap-south-1"

# Sensitive

# Sensitive


LOCAL_CONTAINER_NAME="nimbus-bot-local"
PROD_CONTAINER_NAME="nimbus-bot-prod"

stop_and_remove_container() {
  local container_name=$1
  if [ "$(docker ps -aq -f name=$container_name)" ]; then
    echo "Stopping and removing existing container: $container_name"
    docker stop $container_name >/dev/null 2>&1 || true
    docker rm $container_name >/dev/null 2>&1 || true
  fi
}

build() {
  docker build -t $IMAGE_NAME:$IMAGE_TAG .
}

run_local() {
  stop_and_remove_container "$LOCAL_CONTAINER_NAME"
  docker run -d \
    --name "$LOCAL_CONTAINER_NAME" \
    -p 8080:8080 \
    --env-file "$ENV_FILE" \
    "$IMAGE_NAME:$IMAGE_TAG"
}

build_run_local() {
  build
  run_local
}

pull_from_ecr() {
  aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $ECR_URL
  docker pull "$ECR_URL:$IMAGE_TAG"
}

build_and_push_to_ecr() {
  aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $ECR_URL
  docker build -t "$ECR_URL:$IMAGE_TAG" .
  docker push "$ECR_URL:$IMAGE_TAG"
}

run_prod() {
  stop_and_remove_container "$PROD_CONTAINER_NAME"
  docker run -d \
    --name "$PROD_CONTAINER_NAME" \
    -p 80:8000 \
    -e DISCORD_BOT_TOKEN="$DISCORD_BOT_TOKEN" \
    -e DISCORD_CHANNEL_ID="$DISCORD_CHANNEL_ID" \
    -e REPO_TOKEN="$GITHUB_TOKEN" \
    -e ACTIONS_WEBHOOK_SECRET="$GITHUB_WEBHOOK_SECRET" \
    -e AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID" \
    -e AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY" \
    -e AWS_REGION="$AWS_REGION" \
    -e DYNAMO_TABLE="$DYNAMO_TABLE" \
    -e APP_ENV="production" \
    -e DEBUG="false" \
    "$ECR_URL:$IMAGE_TAG"
}

build_run_prod() {
  build_and_push_to_ecr
  run_prod
}

help_menu() {
  echo "Usage: ./manage.sh [build | run-local | build-run-local | pull-from-ecr | build-push-ecr | run-prod | build-run-prod]"
  echo "Options:"
  echo "  build            : Build the Docker image"
  echo "  run-local        : Run the Docker container locally with .env file"
  echo "  build-run-local  : Build and run the Docker container locally"
  echo "  pull-from-ecr    : Pull Docker image from ECR"
  echo "  build-push-ecr   : Build Docker image and push to ECR"
  echo "  run-prod         : Run the Docker container for production with environment variables"
  echo "  build-run-prod   : Build and run the Docker container for production"
}

case "$1" in
  build)
    build
    ;;
  run-local)
    run_local
    ;;
  build-run-local)
    build_run_local
    ;;
  pull-from-ecr)
    pull_from_ecr
    ;;
  build-push-ecr)
    build_and_push_to_ecr
    ;;
  run-prod)
    run_prod
    ;;
  build-run-prod)
    build_run_prod
    ;;
  *)
    help_menu
    ;;
esac
