FROM golang:1.15-alpine

COPY ./ /go/src/github.com/edersonbrilhante/vilicus
WORKDIR /go/src/github.com/edersonbrilhante/vilicus

RUN go mod download -x
