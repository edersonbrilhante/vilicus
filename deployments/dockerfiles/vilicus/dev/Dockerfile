FROM vilicus/registry:latest as certs

FROM korbin/dockerize:latest as dockerize

FROM golang:1.15-alpine

COPY ./ /go/src/github.com/edersonbrilhante/vilicus
WORKDIR /go/src/github.com/edersonbrilhante/vilicus

COPY --from=dockerize /usr/local/bin/dockerize /usr/local/bin/dockerize

COPY --from=certs /opt/vilicus/certs/vilicus.crt /usr/local/share/ca-certificates
RUN update-ca-certificates

RUN go mod download -x
