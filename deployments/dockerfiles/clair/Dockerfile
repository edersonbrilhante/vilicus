FROM jwilder/dockerize::0.6.1 as dockerize

FROM quay.io/projectquay/clair:4.0.0
LABEL vilicus.app.version=4.0.0

COPY deployments/dockerfiles/clair/config.yaml /tmp/config.yaml

USER root:root
COPY --from=dockerize /usr/local/bin/dockerize /usr/local/bin/dockerize
RUN chmod 755 /usr/local/bin/dockerize
USER nobody:nobody

ENTRYPOINT []