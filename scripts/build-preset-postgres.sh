#!/bin/bash

# This script will build preset postgres image

COLOR_RESET="\033[0;39;49m"
COLOR_RED="\033[38;5;161m"
COLOR_YELO="\033[38;5;227m"

DUMP_PATH=./local-volumes/dump-sql
mkdir -p $DUMP_PATH


printf $COLOR_YELO"Build preset postgres: Starting\n"$COLOR_RESET

printf $COLOR_YELO"Test connection with vilicus: Starting\n"$COLOR_RESET
OK=$(docker exec -i vilicus_client sh -c "dockerize -wait http://vilicus:8080/healthz -wait-retry-interval 60s -timeout 1000s echo"  2>&1 | grep "Command finished successfully.")

if [[ ! -z "$OK" ]]
then    
    printf $COLOR_YELO"Dump databases: Starting\n"$COLOR_RESET
    for db in vilicus_db clair_db anchore_db
    do
        docker exec vilicus_postgres sh -c "cd /tmp/; pg_dump -U username -d $db > $db.sql; env GZIP=-9 tar cvzf $db.tar.gz $db.sql; rm $db.sql" && docker cp vilicus_postgres:/tmp/$db.tar.gz $DUMP_PATH &
    done
    wait    
    printf $COLOR_YELO"Dump databases: Done\n"$COLOR_RESET

    printf $COLOR_YELO"Build postgres with dump files: Starting\n"$COLOR_RESET
    docker build -f deployments/dockerfiles/postgres/preset/Dockerfile -t vilicus/postgres:preset .
    printf $COLOR_YELO"Build postgres with dump files: Done\n"$COLOR_RESET
    
    printf $COLOR_YELO"Build preset postgres: Done\n"$COLOR_RESET
else 
    printf $COLOR_RED"Test connection with vilicus: Fail\n"$COLOR_RESET
    
    printf $COLOR_RED"Build preset postgres: Fail\n"$COLOR_RESET
    exit 2
fi