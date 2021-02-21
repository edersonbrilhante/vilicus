FROM vilicus/trivy:local-update as db

FROM vilicus/trivy:base

COPY --from=db /root/.cache/trivy/ /root/.cache/trivy