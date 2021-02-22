#!/bin/bash

set -e
set -u

function restore_database() {
	echo "Uncompressing database file"
	tar -xf /opt/vilicus/data/anchore_db.tar.gz --directory /opt/vilicus/data
	
	echo "Restoring database"
	psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d postgres < "/opt/vilicus/data/anchore_db.sql"
}

mkdir -p $PGDATA
chown -R postgres:postgres $PGDATA

echo "Database restore requested"
restore_database
echo "Database restored"

rm /opt/vilicus/data/*