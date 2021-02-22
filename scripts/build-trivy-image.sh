#!/bin/bash

# This script will build preset trivydb image

set -e
set -u

COLOR_RESET="\033[0;39;49m"
COLOR_YELO="\033[38;5;227m"


build_trivydb () {
    printf $COLOR_YELO"Build preset trivydb image: Starting\n\n"$COLOR_RESET

    printf $COLOR_YELO"Run docker commit: Starting\n"$COLOR_RESET
    CID=$(docker inspect --format="{{.Id}}" trivy)
    docker commit $CID vilicus/trivydb:local-update
    printf $COLOR_YELO"Run docker commit: Done\n\n"$COLOR_RESET

    printf $COLOR_YELO"Build preset trivydb: Starting\n"$COLOR_RESET
    docker build -f deployments/dockerfiles/trivy/db/Dockerfile -t vilicus/trivydb:latest .
    printf $COLOR_YELO"Build preset trivydb: Done\n\n"$COLOR_RESET

    printf $COLOR_YELO"Build preset trivydb image: Done\n\n"$COLOR_RESET
}

run_updater() {
    printf $COLOR_YELO"Run updater: Starting\n\n"$COLOR_RESET
    
    printf $COLOR_YELO"Run compose: Starting\n"$COLOR_RESET    
    docker-compose -f deployments/docker-compose.updater.yml up --build -d --force  --remove-orphans --renew-anon-volumes trivy
    printf $COLOR_YELO"Run compose: Done\n\n"$COLOR_RESET

    printf $COLOR_YELO"Run trivy: Starting\n"$COLOR_RESET
    docker exec trivy sh -c 'trivy server --listen 0.0.0.0:8080 --download-db-only'
    printf $COLOR_YELO"Run trivy: Done\n\n"$COLOR_RESET
    
    printf $COLOR_YELO"Run updater: Done\n\n"$COLOR_RESET
}

build_trivy() {
    printf $COLOR_YELO"Build trivy image: Starting\n"$COLOR_RESET
    docker build -f deployments/dockerfiles/trivy/Dockerfile -t vilicus/trivy:latest .
    printf $COLOR_YELO"Build trivy image: Done\n\n"$COLOR_RESET
}

build_trivy

run_updater

build_trivydb
