FROM golang:1.15-alpine

COPY ./ /go/src/github.com/edersonbrilhnte/vilicus
WORKDIR /go/src/github.com/edersonbrilhnte/vilicus

RUN go mod download -x
