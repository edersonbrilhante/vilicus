FROM golang:1.15-alpine AS builder

WORKDIR /build/
COPY . /build/

RUN apk add \
    git \
    build-base

RUN make check-deps
RUN make build-linux

FROM golang:1.15-alpine AS waitfori

WORKDIR /build/

RUN apk add \
    curl

ENV WAITFORIT_VERSION="v2.4.1"
RUN curl -o /build/waitforit -sSL https://github.com/maxcnunes/waitforit/releases/download/$WAITFORIT_VERSION/waitforit-linux_amd64 && \
    chmod +x /build/waitforit

FROM golang:1.15-alpine AS final
LABEL vilicus.app.version=dev

WORKDIR /run
EXPOSE 8080
USER nobody:nobody

COPY --from=builder /build/vilicus-api-bin /bin/vilicus-api
COPY --from=builder /build/vilicus-migration-bin /bin/vilicus-migration
COPY --from=builder /build/vilicus-client-bin /bin/vilicus-client

COPY --from=builder /build/configs/conf.docker-compose.yaml /run/conf.yaml

COPY --from=waitfori /build/waitforit /run/waitforit
