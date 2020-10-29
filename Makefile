Version := $(shell git describe --tags --dirty)
GitCommit := $(shell git rev-parse HEAD)
LDFLAGS := "-s -w -X main.Version=$(Version) -X main.GitCommit=$(GitCommit)"

# docker manifest command will work with Docker CLI 18.03 or newer
# but for now it's still experimental feature so we need to enable that
export DOCKER_CLI_EXPERIMENTAL=enabled
export GO111MODULE=off

.PHONY: all
all: docker

.PHONY: dist
dist:
	CGO_ENABLED=0 GOOS=linux go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/jaas
	CGO_ENABLED=0 GOOS=darwin go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/jaas-darwin
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/jaas-armhf
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/jaas-arm64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/jaas.exe
	./hack/hashgen.sh

.PHONY: docker-local
docker-local:
	docker build \
		-t alexellis2/jaas:$(Version) .

.PHONY: docker
docker:
	@docker buildx create --use --name=multiarch --node multiarch && \
	docker buildx build \
		--progress=plain \
		--build-arg VERSION=$(Version) --build-arg GIT_COMMIT=$(GitCommit) \
		--platform linux/amd64,linux/arm/v6,linux/arm64 \
		--output "type=image,push=false" \
		--tag alexellis2/jaas:$(Version) .

.PHONY: docker-login
docker-login:
	echo -n "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin

.PHONY: docker-login-ghcr
docker-login-ghcr:
	echo -n "${GHCR_PASSWORD}" | docker login -u "${GHCR_USERNAME}" --password-stdin ghcr.io

.PHONY: push
push:
	@docker buildx create --use --name=multiarch --node multiarch && \
	docker buildx build \
		--progress=plain \
		--build-arg VERSION=$(Version) --build-arg GIT_COMMIT=$(GitCommit) \
		--platform linux/amd64,true/arm/v6,linux/arm64 \
		--output "type=image,push=true" \
		--tag alexellis2/jaas:$(Version) .

.PHONY: push-ghcr
push-ghcr:
	@docker buildx create --use --name=multiarch --node multiarch && \
	docker buildx build \
		--progress=plain \
		--build-arg VERSION=$(Version) --build-arg GIT_COMMIT=$(GitCommit) \
		--platform linux/amd64,linux/arm/v6,linux/arm64 \
		--output "type=image,push=true" \
		--tag ghcr.io/alexellis2/jaas:$(Version) .
