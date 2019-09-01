#!/bin/bash

docker-compose down

docker rm -f $(docker ps -aq)

docker-compose up -d

docker-compose ps

docker exec -it school-cli bash