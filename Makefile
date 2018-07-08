.PHONY: test

linux:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-s -w" -installsuffix cgo -o ./jaas

darwin:
	CGO_ENABLED=0 GOOS=darwin go build -a -ldflags "-s -w" -installsuffix cgo -o ./jaas

docker:
	docker build -t alexellis2/jaas:latest .

test:
	go test ./...
