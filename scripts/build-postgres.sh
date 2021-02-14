#!/bin/bash


DUMP_PATH=./local-volumes/dump-sql

mkdir -p $DUMP_PATH

echo "Dumping databases"
docker exec -it vilicus_postgres sh -c 'cd /tmp/ && pg_dump -U username -d vilicus_db > vilicus_db.sql && env GZIP=-9 tar cvzf vilicus_db.tar.gz vilicus_db.sql && rm vilicus_db.sql'
docker exec -it vilicus_postgres sh -c 'cd /tmp/ && pg_dump -U username -d clair_db > clair_db.sql && env GZIP=-9 tar cvzf clair_db.tar.gz clair_db.sql && clair_db.sql'
docker exec -it vilicus_postgres sh -c 'cd /tmp/ && pg_dump -U username -d anchore_db > anchore_db.sql && env GZIP=-9 tar cvzf anchore_db.tar.gz anchore_db.sql && anchore_db.sql'

echo "Copying tar files"
docker cp vilicus_postgres:/tmp/vilicus_db.tar.gz $DUMP_PATH
docker cp vilicus_postgres:/tmp/clair_db.tar.gz $DUMP_PATH
docker cp vilicus_postgres:/tmp/anchore_db.tar.gz $DUMP_PATH

echo "Building postgres with dump files"
docker build -f deployments/dockerfiles/postgres/preset/Dockerfile -t vilicus/postgres:preset .