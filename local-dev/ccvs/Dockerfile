FROM golang:1.15-alpine

COPY ./ /go/src/github.com/edersonbrilhnte/ccvs
WORKDIR /go/src/github.com/edersonbrilhnte/ccvs

RUN go mod download -x
