VERSION_API=develop
container_name=azuki774/mawinter-discord

.PHONY: build test run
build:
	docker build -t $(container_name):$(VERSION_API) -f build/Dockerfile .

test:
	go test -v ./...
