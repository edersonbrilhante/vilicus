FROM vilicus/postgres-presets:local-update as db

FROM vilicus/postgres-volume:old

COPY --chown=postgres:postgres --from=db /postgres /postgres