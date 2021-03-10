FROM vilicus/registry:latest as certs

FROM golang:1.15-alpine AS builder

WORKDIR /build/
COPY . /build/

RUN apk add \
    git \
    build-base

RUN make check-deps
RUN make build-linux

FROM korbin/dockerize:latest as dockerize

FROM golang:1.15-alpine AS final
LABEL vilicus.app.version=dev

COPY --from=certs /opt/vilicus/certs/vilicus.crt /usr/local/share/ca-certificates
RUN update-ca-certificates

WORKDIR /run

USER nobody:nobody

COPY --from=builder /build/vilicus-api-bin /bin/vilicus-api
COPY --from=builder /build/vilicus-migration-bin /bin/vilicus-migration
COPY --from=builder /build/vilicus-client-bin /bin/vilicus-client

COPY --from=builder /build/configs/conf.docker-compose.yaml /opt/vilicus/configs/conf.yaml
COPY --from=builder /build/contrib /opt/vilicus/contrib

COPY --from=dockerize /usr/local/bin/dockerize /usr/local/bin/dockerize
