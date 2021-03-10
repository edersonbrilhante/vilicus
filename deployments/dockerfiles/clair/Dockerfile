FROM vilicus/registry:latest as certs

FROM korbin/dockerize:latest as dockerize

FROM quay.io/projectquay/clair:4.0.2
LABEL vilicus.app.version=4.0.2

COPY deployments/dockerfiles/clair/config.yaml /opt/vilicus/data/config.yaml

USER root:root
COPY --from=dockerize /usr/local/bin/dockerize /usr/local/bin/dockerize
RUN chmod 755 /usr/local/bin/dockerize

COPY --from=certs /opt/vilicus/certs/vilicus.crt /etc/pki/ca-trust/source/anchors
RUN update-ca-trust extract
# quick win: https://unix.stackexchange.com/questions/525601/update-ca-trust-extract-not-adding-certificates-to-ca-bundle
RUN cat /etc/pki/ca-trust/source/anchors/vilicus.crt >> /etc/pki/tls/certs/ca-bundle.crt

# USER nobody:nobody

ENTRYPOINT []