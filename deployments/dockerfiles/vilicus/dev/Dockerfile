FROM korbin/dockerize:latest as dockerize

FROM golang:1.15-alpine

COPY ./ /go/src/github.com/edersonbrilhante/vilicus
WORKDIR /go/src/github.com/edersonbrilhante/vilicus

COPY --from=dockerize /usr/local/bin/dockerize /usr/local/bin/dockerize

RUN go mod download -x
