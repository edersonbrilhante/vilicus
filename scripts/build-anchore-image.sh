#!/bin/bash

# This script will build anchore and anchoredb images

set -e
set -u

COLOR_RESET="\033[0;39;49m"
COLOR_YELO="\033[38;5;227m"


build_anchoredb () {
    printf $COLOR_YELO"Build preset anchoredb image: Starting\n\n"$COLOR_RESET


    printf $COLOR_YELO"Stop container anchore: Starting\n"$COLOR_RESET
    docker stop anchore
    printf $COLOR_YELO"Stop container anchore: Done\n\n"$COLOR_RESET

    printf $COLOR_YELO"Kill postgres pid: Starting\n"$COLOR_RESET
    docker exec -u postgres anchoredb sh -c 'pg_ctl stop -m smart'
    printf $COLOR_YELO"Kill postgres pid: Done\n\n"$COLOR_RESET

    printf $COLOR_YELO"Run docker commit: Starting\n"$COLOR_RESET
    CID=$(docker inspect --format="{{.Id}}" anchoredb)
    docker commit $CID vilicus/anchoredb:local-update
    printf $COLOR_YELO"Run docker commit: Done\n\n"$COLOR_RESET

    printf $COLOR_YELO"Build preset anchoredb: Starting\n"$COLOR_RESET
    docker build -f deployments/dockerfiles/anchore/db/Dockerfile -t vilicus/anchoredb:latest .
    printf $COLOR_YELO"Build preset anchoredb: Done\n\n"$COLOR_RESET

    printf $COLOR_YELO"Build preset anchoredb image: Done\n\n"$COLOR_RESET
}

run_updater() {
    printf $COLOR_YELO"Run updater: Starting\n\n"$COLOR_RESET
    
    printf $COLOR_YELO"Run compose: Starting\n"$COLOR_RESET    
    docker-compose -f deployments/docker-compose.updater.yml -f deployments/docker-compose.adminer.yml up --build -d --force  --remove-orphans --renew-anon-volumes anchore adminer
    printf $COLOR_YELO"Run compose: Done\n\n"$COLOR_RESET

    printf $COLOR_YELO"Starting postgres\n"$COLOR_RESET
    docker exec anchoredb sh -c 'docker-entrypoint.sh postgres' &

    printf $COLOR_YELO"Test connection with anchore: Starting\n"$COLOR_RESET
    docker run --network container:anchore vilicus/vilicus:latest sh -c "dockerize -wait http://anchore:8228/health -wait-retry-interval 10s -timeout 1000s echo done"
    printf $COLOR_YELO"Test connection with anchore: Done\n\n"$COLOR_RESET
    
    printf $COLOR_YELO"Run sync feeds: Starting\n"$COLOR_RESET
    docker exec anchore sh -c 'anchore-cli system wait'
    printf $COLOR_YELO"Run sync feeds: Done\n\n"$COLOR_RESET    

    printf $COLOR_YELO"Run updater: Done\n\n"$COLOR_RESET
}

build_anchoredb_files() {
    printf $COLOR_YELO"Build anchoredb:files image: Starting\n"$COLOR_RESET
    docker build -f deployments/dockerfiles/anchore/db/files/Dockerfile -t vilicus/anchoredb:files .
    printf $COLOR_YELO"Build anchoredb:files image: Done\n\n"$COLOR_RESET
}

build_anchore() {
    printf $COLOR_YELO"Build anchore image: Starting\n"$COLOR_RESET
    docker build -f deployments/dockerfiles/anchore/Dockerfile -t vilicus/anchore:latest .
    printf $COLOR_YELO"Build anchore image: Done\n\n"$COLOR_RESET
}

build_anchore

build_anchoredb_files

run_updater

build_anchoredb