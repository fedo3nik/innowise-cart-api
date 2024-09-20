.PHONY: all

all:
	docker-compose --env-file ./internal/pkg/config/envs/cfg.env up --build

tests:
	go test -v ./...

docker-tests:
	docker-compose --env-file ./internal/pkg/config/envs/cfg.env up --build
	docker-compose exe -T http go test ./...
	docker-compose down

docker-build-test:
	docker build -t api-tests --target=test
 