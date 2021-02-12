#!/bin/bash

# This script will build docker images with rsync, and data from postgres used for further processing

CID=$(docker inspect --format="{{.Id}}" vilicus_postgres)
docker commit $CID vilicus/postgres-presets:local-update
docker build -t vilicus/postgres-volume:old -f deployments/dockerfiles/postgres/rsync/old.Dockerfile .
docker build -t vilicus/postgres-volume:new -f deployments/dockerfiles/postgres/rsync/new.Dockerfile .