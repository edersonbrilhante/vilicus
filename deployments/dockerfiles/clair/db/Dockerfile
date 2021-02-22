FROM vilicus/clairdb:local-update as db

FROM postgres:9.6.21-alpine

COPY --chown=postgres:postgres --from=db /data/ /data
