FROM golang:1.7.3
MAINTAINER alexellis2@gmail.com

RUN mkdir -p /go/src/github.com/alexellis2/jaas
WORKDIR /go/src/github.com/alexellis2/jaas
COPY ./app.go ./

RUN go get -v -d
RUN go build
