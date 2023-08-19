-include .env
export $(shell sed 's/=.*//' .env)

indexer:
	cd cmd/indexer && go run . -c ../../build/dipdup.yml

build:
	docker-compose up -d -- build

lint:
	golangci-lint run

test:
	go test ./...
