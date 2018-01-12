FROM golang:1.9.2 as build

MAINTAINER alexellis2@gmail.com

RUN mkdir -p /go/src/github.com/alexellis/jaas
WORKDIR /go/src/github.com/alexellis/jaas

COPY vendor     vendor
COPY cmd        cmd
COPY main.go    .

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-s -w" -installsuffix cgo -o /root/jaas

FROM alpine:3.6
WORKDIR /root/
COPY --from=build /root/jaas /root/jaas

ENTRYPOINT ["/root/jaas"]
