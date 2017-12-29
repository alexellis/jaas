FROM golang:1.9.2 as build

MAINTAINER alexellis2@gmail.com

RUN mkdir -p /go/src/github.com/alexellis2/jaas
WORKDIR /go/src/github.com/alexellis2/jaas

COPY app.go         .
COPY show_tasks.go  .
COPY poll.go        .
COPY vendor vendor

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-s -w" -installsuffix cgo -o /root/jaas

FROM alpine:3.6
WORKDIR /root/
COPY --from=build /root/jaas /root/jaas

ENTRYPOINT ["/root/jaas"]
