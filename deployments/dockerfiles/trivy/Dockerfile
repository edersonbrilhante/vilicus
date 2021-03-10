FROM vilicus/registry:latest as certs

FROM aquasec/trivy:0.15.0
LABEL vilicus.app.version=0.15.0

COPY --from=certs /opt/vilicus/certs/vilicus.crt /usr/local/share/ca-certificates
RUN update-ca-certificates