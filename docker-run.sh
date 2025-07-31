#!/bin/bash

set -e

echo "Launching Rate Limiter"


cleanup() {
    docker-compose down
}


trap cleanup SIGINT

echo "Building docker-compose"
docker-compose up --build -d

#Esperando para que os contianers se iniciem
sleep 10

docker-compose ps

docker-compose logs -f 
