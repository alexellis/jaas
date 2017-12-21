FROM golang:1.9.2-alpine3.7 AS build
MAINTAINER alexellis2@gmail.com

RUN mkdir -p /go/src/github.com/alexellis2/jaas
WORKDIR /go/src/github.com/alexellis2/jaas
COPY . .

RUN go build

FROM alpine

COPY --from=build /go/src/github.com/alexellis2/jaas/jaas /jaas

ENTRYPOINT ["/jaas"]
