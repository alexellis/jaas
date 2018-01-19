FROM golang:1.9.2 as build

MAINTAINER alexellis2@gmail.com

RUN mkdir -p /go/src/github.com/alexellis/jaas
WORKDIR /go/src/github.com/alexellis/jaas

COPY . .

RUN VERSION=$(git describe --all --exact-match `git rev-parse HEAD` | grep tags | sed 's/tags\///') \
    && GIT_COMMIT=$(git rev-list -1 HEAD) \
    && CGO_ENABLED=0 GOOS=linux go build --ldflags "-s -w -X github.com/alexellis/jaas/version.GitCommit=${GIT_COMMIT} -X github.com/alexellis/jaas/version.Version=${VERSION}" -a -installsuffix cgo -o /root/jaas

FROM alpine:3.6
WORKDIR /root/
COPY --from=build /root/jaas /root/jaas

ENTRYPOINT ["/root/jaas"]
