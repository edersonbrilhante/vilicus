#!/bin/bash

# This script will process docker images to create a new layer just with changes

docker-compose -f  deployments/rsync-image-diff.docker-compose.yml up -d
docker run -it  -v $(PWD)/local-volumes/output:/tmp/output jwilder/dockerize:0.6.1 dockerize -wait file:///tmp/output/job.rsync.log -wait-retry-interval 2s -timeout 600s
docker build -t vilicus/postgres-presets:processed -f deployments/dockerfiles/postgres/rsync/processed.Dockerfile .