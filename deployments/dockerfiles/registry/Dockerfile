FROM alpine as certs
COPY deployments/dockerfiles/registry/scripts /opt/vilicus/scripts

RUN apk add \
    openssl \
    bash

RUN mkdir -p /opt/vilicus/certs

RUN chmod +x /opt/vilicus/scripts/create-certs.sh
RUN /opt/vilicus/scripts/create-certs.sh

FROM registry:2

ENV REGISTRY_HTTP_TLS_CERTIFICATE=/opt/vilicus/certs/vilicus.crt
ENV REGISTRY_HTTP_TLS_KEY=/opt/vilicus/certs/vilicus.key

COPY --from=certs /opt/vilicus/certs /opt/vilicus/certs