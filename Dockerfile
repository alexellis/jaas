FROM golang:1.7.3
MAINTAINER alexellis2@gmail.com
RUN go get -v -d github.com/docker/docker/api/types \
    && go get -v -d github.com/docker/docker/api/types/filters \
    && go get -v -d github.com/docker/docker/api/types/swarm \
    && go get -v -d github.com/docker/docker/client \
    && go get -v -d golang.org/x/net/context

RUN mkdir -p /go/src/github.com/alexellis2/jaas
WORKDIR /go/src/github.com/alexellis2/jaas
COPY ./app.go ./

RUN go build

ENTRYPOINT ["./jaas"]