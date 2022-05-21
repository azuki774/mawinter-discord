VERSION_API=develop
container_name=azuki774/mawinter-discord
container_id=`docker ps -aqf "name=mawinter-discord"`

.PHONY: build test run stop logs
build:
	docker build -t $(container_name):$(VERSION_API) -f build/Dockerfile .

test:
	go test -v ./...

run:
	docker-compose -f deploy/docker/docker-compose.yml up -d

stop:
	docker-compose -f deploy/docker/docker-compose.yml down

rebuild:
	make stop && make && make run

logs:
	echo $(container_id)
	docker logs $(container_id)
