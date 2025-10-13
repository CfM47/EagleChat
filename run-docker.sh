#!/bin/bash

set -e

go work vendor

docker swarm leave --force || true

docker build -t eaglechat-client:latest -f iac/client/Dockerfile .
docker build -t eaglechat-id_manager:latest -f iac/id_manager/Dockerfile .

docker swarm init --advertise-addr 10.37.129.212

docker stack deploy -c  iac/docker-compose.local.yml eaglechat 
