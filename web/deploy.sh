#!/bin/bash

docker compose down
docker rmi fast-storage-service:latest
docker build -t fast-storage-service:latest .
docker compose up -d
docker compose logs -f
