#!/bin/bash

# This script will build preset postgres images

set -e
set -u

COLOR_RESET="\033[0;39;49m"
COLOR_YELO="\033[38;5;227m"

DUMP_PATH=./local-volumes/dump-sql

run_postgres () {
    printf $COLOR_YELO"Running postgres\n"$COLOR_RESET
    docker exec vilicus_postgres sh -c 'docker-entrypoint.sh postgres' &

    printf $COLOR_YELO"Test connection with vilicus: Starting\n"$COLOR_RESET
    docker run --network container:vilicus_postgres vilicus/vilicus:latest sh -c "dockerize -wait tcp://vilicus_postgres:5432 -wait-retry-interval 10s -timeout 10000s echo done"
    printf $COLOR_YELO"Test connection with vilicus: Done\n"$COLOR_RESET
}

preset_volume () {
    
    printf $COLOR_YELO"Kill postgres pid: Starting\n"$COLOR_RESET
    docker exec vilicus_postgres sh -c 'kill -INT `head -1 /data/postmaster.pid`'
    printf $COLOR_YELO"Kill postgres pid: Done\n"$COLOR_RESET

    printf $COLOR_YELO"Run docker commit for postgres: Starting\n"$COLOR_RESET
    CID=$(docker inspect --format="{{.Id}}" vilicus_postgres)
    docker commit $CID vilicus/postgres:local-update
    printf $COLOR_YELO"Run docker commit for postgres: Done\n"$COLOR_RESET

    printf $COLOR_YELO"Build preset postgres: Starting\n"$COLOR_RESET
    docker build -f deployments/dockerfiles/postgres/preset/volume/Dockerfile -t vilicus/postgres:preset-volume-latest .
    printf $COLOR_YELO"Build preset postgres: Done\n"$COLOR_RESET
}

preset_files () {
    run_postgres
    printf $COLOR_YELO"Dump databases: Starting\n"$COLOR_RESET
    for db in vilicus_db clair_db anchore_db
    do
        docker exec vilicus_postgres sh -c "cd /tmp/; pg_dump -U username -d $db > $db.sql; env GZIP=-9 tar cvzf $db.tar.gz $db.sql; rm $db.sql" && docker cp vilicus_postgres:/tmp/$db.tar.gz $DUMP_PATH &
    done
    wait    
    printf $COLOR_YELO"Dump databases: Done\n"$COLOR_RESET

    printf $COLOR_YELO"Build postgres with dump files: Starting\n"$COLOR_RESET
    docker build --build-arg DUMP_PATH=$DUMP_PATH -f deployments/dockerfiles/postgres/preset/Dockerfile -t vilicus/postgres:preset-files-latest .
    printf $COLOR_YELO"Build postgres with dump files: Done\n"$COLOR_RESET
}

run_updater () {
    printf $COLOR_YELO"Run compose: Starting\n"$COLOR_RESET
    docker-compose -f deployments/docker-compose.yml up --build -d --force  --remove-orphans --renew-anon-volumes
    printf $COLOR_YELO"Run compose: Done\n"$COLOR_RESET

    run_postgres

    printf $COLOR_YELO"Sleep for 3 hours: Starting\n"$COLOR_RESET
    sleep 10800 # 3 hours
    printf $COLOR_YELO"Sleep for 3 hours: Done\n"$COLOR_RESET

    printf $COLOR_YELO"Stop app containers: Starting\n"$COLOR_RESET
    docker stop clair vilicus anchore_engine trivy
    printf $COLOR_YELO"Stop app containers: Done\n"$COLOR_RESET
}

no_updater() {
    printf $COLOR_YELO"Run compose: Starting\n"$COLOR_RESET
    docker-compose -f deployments/docker-compose.yml up --build -d --force  --remove-orphans --renew-anon-volumes postgres
    printf $COLOR_YELO"Run compose: Done\n"$COLOR_RESET

}

if [[ $1 == "updater" ]]; then
    echo "run_updater"
else
    echo "no_updater"
fi

printf $COLOR_YELO"Run Preset Volume: Starting\n"$COLOR_RESET
preset_volume
printf $COLOR_YELO"Run Preset Volume: Done\n"$COLOR_RESET

printf $COLOR_YELO"Run Preset Files: Starting\n"$COLOR_RESET
preset_files
printf $COLOR_YELO"Run Preset Files: Done\n"$COLOR_RESET