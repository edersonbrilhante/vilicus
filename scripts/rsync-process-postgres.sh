#!/bin/bash

# This script will process docker images to create a new layer just with changes

COMPOSE_IGNORE_ORPHANS=True docker-compose -f  deployments/docker-compose.rsync.yml up -d --force-recreate
docker exec postgres_volume_new /.docker-image-diff/rsync-new.sh
docker build -t vilicus/postgres-presets:processed -f deployments/dockerfiles/postgres/rsync/processed.Dockerfile .