FROM golang:1.15-alpine AS builder

WORKDIR /build/
COPY . /build/

RUN apk update && \
    apk add \
    git \
    build-base

RUN make check-deps
RUN make build-linux

FROM golang:1.15-alpine AS final
LABEL vilicus.app.version=dev

WORKDIR /run
EXPOSE 8080
USER nobody:nobody

COPY --from=builder /build/vilicus-api-bin /bin/vilicus-api
COPY --from=builder /build/vilicus-migration-bin /bin/vilicus-migration
COPY --from=builder /build/cmd/api/conf.local.yaml /run/conf.yaml