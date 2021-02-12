FROM vilicus/postgres-presets:latest as db

FROM postgres:9.6
RUN apt-get update && apt-get install rsync -yq

COPY --chown=postgres:postgres --from=db /postgres /postgres