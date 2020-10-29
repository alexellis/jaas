FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.13 as build

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG LDFLAGS

ENV GO111MODULE=off
ENV CGO_ENABLED=0

LABEL maintainer alex@openfaas.com

WORKDIR /go/src/github.com/alexellis/jaas

RUN apk --no-cache add git

COPY . .

# Run a gofmt and exclude all vendored code.
RUN test -z "$(gofmt -l $(find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./function/vendor/*"))" || { echo "Run \"gofmt -s -w\" on your Golang code"; exit 1; }

ARG GO111MODULE="off"
ARG GOPROXY=""

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go test ./... -cover

RUN CGO_ENABLED=${CGO_ENABLED} GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o /usr/bin/jaas .

FROM --platform=${TARGETPLATFORM:-linux/amd64} alpine:3.12
# Add non root user and certs
RUN apk --no-cache add ca-certificates \
    && addgroup -S app && adduser -S -g app app \
    && mkdir -p /home/app \
    && chown app /home/app

WORKDIR /home/app

COPY --from=build /go/src/handler/handler    .
COPY --from=build /usr/bin/fwatchdog         .
COPY --from=build /go/src/handler/function/  .

RUN chown -R app /home/app

USER app
USER root
WORKDIR /root/
COPY --from=build /usr/bin/jaas /usr/bin/jaas
ENTRYPOINT ["/usr/bin/jaas"]

FROM golang:1.13 as build


COPY .git       .git
COPY cmd        cmd
COPY vendor     vendor
COPY version    version
COPY main.go    .

RUN VERSION=$(git describe --all --exact-match `git rev-parse HEAD` | grep tags | sed 's/tags\///') \
    && GIT_COMMIT=$(git rev-list -1 HEAD) \
    && CGO_ENABLED=0 GOOS=linux go build --ldflags "-s -w -X github.com/alexellis/jaas/version.GitCommit=${GIT_COMMIT} -X github.com/alexellis/jaas/version.Version=${VERSION}" -a -installsuffix cgo -o /root/jaas

FROM alpine:3.12


