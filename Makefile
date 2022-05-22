version_api=develop
container_name=azuki774/mawinter-discord
container_id=`docker ps -aqf "name=mawinter-discord"`

.PHONY: build test run stop logs
build:
	docker build -t $(container_name):$(version_api) -f build/Dockerfile .

push:
	docker tag $(container_name):$(version_api) ghcr.io/$(container_name):$(version_api)
	docker push ghcr.io/$(container_name):$(version_api)

test:
	go test -v ./...

run:
	docker-compose -f deploy/docker/docker-compose.yml up -d

stop:
	docker-compose -f deploy/docker/docker-compose.yml down

rebuild:
	make stop && make && make run

logs:
	docker logs $(container_id)
