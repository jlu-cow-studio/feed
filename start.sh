#!/bin/bash

IMAGE_NAME=cowstudio/feed
CONTAINER_NAME=feed

SERVICE_NAME=cowstudio/feed
SERVICE_PORT=3085
SERVICE_ADDRESS=$(curl -s http://ipecho.net/plain)
SIDECAR_PORT=4085

# 构建镜像
docker build -t $IMAGE_NAME .

# 关闭容器
echo "removing....."
docker stop $CONTAINER_NAME
docker rm $CONTAINER_NAME

# 运行容器
echo "starting....."
docker run --name $CONTAINER_NAME -p $SERVICE_PORT:8080 -p $SIDECAR_PORT:8081 -d -e ENV_SERVICE_NAME=$SERVICE_NAME -e ENV_SERVICE_PORT=$SERVICE_PORT -e ENV_SERVICE_ADDRESS=$SERVICE_ADDRESS -e ENV_SIDECAR_PORT=$SIDECAR_PORT $IMAGE_NAME
