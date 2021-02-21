#!/bin/bash

# This script will build preset trivydb image

set -e
set -u

COLOR_RESET="\033[0;39;49m"
COLOR_YELO="\033[38;5;227m"


preset_volume () {
    
    printf $COLOR_YELO"Run docker commit for trivy: Starting\n"$COLOR_RESET
    CID=$(docker inspect --format="{{.Id}}" trivy)
    docker commit $CID vilicus/trivy:local-update
    printf $COLOR_YELO"Run docker commit for trivy: Done\n"$COLOR_RESET

    printf $COLOR_YELO"Build preset trivy: Starting\n"$COLOR_RESET
    docker build -f deployments/dockerfiles/trivy/preset/Dockerfile -t vilicus/trivy:preset-latest .
    printf $COLOR_YELO"Build preset trivy: Done\n"$COLOR_RESET
}

run_updater() {
    printf $COLOR_YELO"Run compose: Starting\n"$COLOR_RESET    
    docker-compose -f deployments/docker-compose.yml -f deployments/docker-compose.updater.yml up --build -d --force  --remove-orphans --renew-anon-volumes trivy
    printf $COLOR_YELO"Run compose: Done\n"$COLOR_RESET
    
    printf $COLOR_YELO"Running trivy\n"$COLOR_RESET
    docker exec trivy sh -c 'trivy server --listen 0.0.0.0:8080 --download-db-only'
}

run_updater

printf $COLOR_YELO"Run Preset Volume: Starting\n"$COLOR_RESET
preset_volume
printf $COLOR_YELO"Run Preset Volume: Done\n"$COLOR_RESET
