FROM vilicus/registry:latest as certs

FROM korbin/dockerize:latest as dockerize

FROM anchore/anchore-engine:v0.9.0
LABEL vilicus.app.version=v0.9.0

COPY deployments/dockerfiles/anchore/config.yaml /config/config.yaml

USER root:root
COPY --from=dockerize /usr/local/bin/dockerize /usr/local/bin/dockerize
RUN chmod 755 /usr/local/bin/dockerize

COPY --from=certs /opt/vilicus/certs/vilicus.crt /etc/pki/ca-trust/source/anchors
RUN update-ca-trust

USER 1000

ENTRYPOINT []