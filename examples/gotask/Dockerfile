FROM golang:1.13 as build

ENV GO111MODULE=off
ENV CGO_ENABLED=0

LABEL maintainer alex@openfaas.com

RUN mkdir -p /go/src/github.com/alexellis/jaas/gotask
WORKDIR /go/src/github.com/alexellis/jaas/gotask

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build --ldflags "-s -w" -a -installsuffix cgo -o /bin/gotask

FROM alpine:3.11

WORKDIR /root/
COPY --from=build /bin/gotask /bin/gotask

ENTRYPOINT ["/bin/gotask"]
