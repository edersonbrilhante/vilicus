#!/bin/bash

# This script will build clair and clairdb images

set -e
set -u

COLOR_RESET="\033[0;39;49m"
COLOR_YELO="\033[38;5;227m"


build_clairdb () {
    printf $COLOR_YELO"Build preset clairdb image: Starting\n\n"$COLOR_RESET


    printf $COLOR_YELO"Stop container clair: Starting\n"$COLOR_RESET
    docker stop clair
    printf $COLOR_YELO"Stop container clair: Done\n\n"$COLOR_RESET

    printf $COLOR_YELO"Kill postgres pid: Starting\n"$COLOR_RESET
    docker exec -u postgres clairdb sh -c 'pg_ctl stop -m smart'
    printf $COLOR_YELO"Kill postgres pid: Done\n\n"$COLOR_RESET


    printf $COLOR_YELO"Run docker commit: Starting\n"$COLOR_RESET
    CID=$(docker inspect --format="{{.Id}}" clairdb)
    docker commit $CID vilicus/clairdb:local-update
    printf $COLOR_YELO"Run docker commit: Done\n\n"$COLOR_RESET
    
    printf $COLOR_YELO"Run cleanup docker image: Starting\n"$COLOR_RESET
    docker image prune -f
    printf $COLOR_YELO"Run cleanup docker image: Done\n\n"$COLOR_RESET
    

    printf $COLOR_YELO"Build preset clairdb: Starting\n"$COLOR_RESET
    docker build -f deployments/dockerfiles/clair/db/Dockerfile -t vilicus/clairdb:latest .
    printf $COLOR_YELO"Build preset clairdb: Done\n\n"$COLOR_RESET

    printf $COLOR_YELO"Build preset clairdb image: Done\n\n"$COLOR_RESET

    printf $COLOR_YELO"Removing preset clairdb: Starting\n"$COLOR_RESET
    docker image rm vilicus/clairdb:local-update
    printf $COLOR_YELO"Removing preset clairdb: Done\n\n"$COLOR_RESET
}

run_updater() {
    printf $COLOR_YELO"Run updater: Starting\n\n"$COLOR_RESET
    
    printf $COLOR_YELO"Run compose: Starting\n"$COLOR_RESET    
    docker-compose -f deployments/docker-compose.updater.yml -f deployments/docker-compose.adminer.yml up --build -d --force  --remove-orphans --renew-anon-volumes clair adminer
    printf $COLOR_YELO"Run compose: Done\n\n"$COLOR_RESET

    printf $COLOR_YELO"Starting postgres\n"$COLOR_RESET
    docker exec clairdb sh -c 'docker-entrypoint.sh postgres' &

    printf $COLOR_YELO"Test connection with clairdb and running clair api: Starting\n"$COLOR_RESET
    docker exec clair sh -c "dockerize -wait tcp://clairdb:5432 -wait-retry-interval 10s -timeout 1000s /bin/clair" &

    printf $COLOR_YELO"Test connection with clair: Starting\n"$COLOR_RESET
    docker run --network container:clair vilicus/vilicus:latest sh -c "dockerize -wait http://clair:6061/healthz -wait-retry-interval 60s -timeout 100000s echo done"
    printf $COLOR_YELO"Test connection with clair: Done\n\n"$COLOR_RESET    
    printf $COLOR_YELO"Test connection with clairdb and running clair api: Done\n\n"$COLOR_RESET

    printf $COLOR_YELO"Run updater: Done\n\n"$COLOR_RESET
}

build_clair() {
    printf $COLOR_YELO"Build clair image: Starting\n"$COLOR_RESET
    docker build -f deployments/dockerfiles/clair/Dockerfile -t vilicus/clair:latest .
    printf $COLOR_YELO"Build clair image: Done\n\n"$COLOR_RESET
}

build_clair

run_updater

build_clairdb