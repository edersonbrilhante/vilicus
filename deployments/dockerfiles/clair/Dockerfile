FROM korbin/dockerize:latest as dockerize

FROM quay.io/projectquay/clair:4.0.0
LABEL vilicus.app.version=4.0.0

COPY deployments/dockerfiles/clair/scanner.yaml /tmp/scanner.yaml
COPY deployments/dockerfiles/clair/updater.yaml /tmp/updater.yaml

USER root:root
COPY --from=dockerize /usr/local/bin/dockerize /usr/local/bin/dockerize
RUN chmod 755 /usr/local/bin/dockerize
USER nobody:nobody

ENTRYPOINT []