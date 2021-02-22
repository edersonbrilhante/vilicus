FROM vilicus/anchoredb:dumpsql as dumpsql

FROM postgres:9.6.21-alpine
LABEL vilicus.app.version=9.6.21-alpine

COPY --chown=postgres:postgres --from=dumpsql /opt/vilicus/data/anchore_db.tar.gz /opt/vilicus/data/anchore_db.tar.gz
COPY deployments/dockerfiles/anchore/db/files/restore-dbs.sh /docker-entrypoint-initdb.d/01.restore-dbs.sh