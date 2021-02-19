FROM vilicus/postgres:local-update as db

FROM vilicus/postgres:base

RUN rm /docker-entrypoint-initdb.d/*

COPY --chown=postgres:postgres --from=db /data/ /data