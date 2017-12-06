#!/usr/bin/env bash

MONGO_NAME=multiadv-mongo

if docker ps -a | grep multiadv-mongo; then
    docker start ${MONGO_NAME}
else
    docker run -d --name ${MONGO_NAME} -p 27017:27017 mongo:latest
fi
./MultiAdv
