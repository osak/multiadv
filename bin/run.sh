#!/usr/bin/env bash

MONGO_NAME=multiadv-mongo

if docker ps | grep multiadv-mongo; then
    docker start ${MONGO_NAME}
else
    docker run -d --name ${MONGO_NAME} mongo:latest
fi
./MultiAdv