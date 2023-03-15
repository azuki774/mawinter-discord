SHELL=/bin/bash
version_api=develop
container_name=azuki774/mawinter-discord
container_id=`docker ps -aqf "name=mawinter-discord"`

.PHONY: build test migration-test run stop logs
build:
	docker build -t $(container_name):$(version_api) -f build/Dockerfile .

test:
	go test -v ./...

# migration-test:
# 	docker compose -f deploy/docker/migration-test.yml up --build -d
# 	sleep 20s
# 	go test -v ./... -tags=integration
# 	docker compose -f deploy/docker/migration-test.yml down
	
run:
	docker compose -f deploy/docker/docker-compose.yml up -d

stop:
	docker compose -f deploy/docker/docker-compose.yml down

rebuild:
	make stop && make && make run

logs:
	docker logs $(container_id)
