FROM vilicus/trivy:local-update as db

FROM vilicus/trivy:latest

COPY --from=db /root/.cache/trivy/db /root/.cache/trivy/db