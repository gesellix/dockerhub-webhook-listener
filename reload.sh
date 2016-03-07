#!/usr/bin/env bash

IMAGE_NAME="$1"

echo "going to reload ${CONTAINER_NAME} with image ${IMAGE_NAME}"

docker pull ${IMAGE_NAME}:latest

docker stop ${CONTAINER_NAME}
docker wait ${CONTAINER_NAME}
docker rm ${CONTAINER_NAME}

docker run --name ${CONTAINER_NAME} -d -p 8000:80 ${IMAGE_NAME}
