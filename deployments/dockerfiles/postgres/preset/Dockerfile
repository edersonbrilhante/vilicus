FROM vilicus/postgres:base

ARG DUMP_PATH=local-volumes/dump-sql

COPY deployments/dockerfiles/postgres/scripts/restore-dbs.sh /docker-entrypoint-initdb.d/02.restore-dbs.sh

COPY --chown=postgres:postgres ${DUMP_PATH}/vilicus_db.tar.gz /opt/vilicus/data/vilicus_db.tar.gz
COPY --chown=postgres:postgres ${DUMP_PATH}/anchore_db.tar.gz /opt/vilicus/data/anchore_db.tar.gz
COPY --chown=postgres:postgres ${DUMP_PATH}/clair_db.tar.gz /opt/vilicus/data/clair_db.tar.gz
