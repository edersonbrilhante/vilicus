FROM vilicus/clairdb:local-update as db

FROM postgres:9.6.22-alpine

COPY --chown=postgres:postgres --from=db /data/ /data
