#!/bin/bash

set -e
set -u

function restore_database() {
	local database=$1
	echo "Uncompressing database '$database' file"
	tar -xf /opt/vilicus/data/$database.tar.gz --directory /opt/vilicus/data
	
	echo "Restoring database '$database'"
	psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d $database < "/opt/vilicus/data/$database.sql"
}

if [ -n "$POSTGRES_MULTIPLE_DATABASES" ]; then
	mkdir -p $PGDATA
	chown -R postgres:postgres $PGDATA
	echo "Multiple database restore requested: $POSTGRES_MULTIPLE_DATABASES"
	for db in $(echo $POSTGRES_MULTIPLE_DATABASES | tr ',' ' '); do
		restore_database $db &
	done;
	wait
	rm /opt/vilicus/data/*
	echo "Multiple databases restored"
fi
