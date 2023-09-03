-include .env
export $(shell sed 's/=.*//' .env)

init:
	chmod +x init.dev.sh && ./init.dev.sh

indexer:
	cd cmd/indexer && go run . -c ../../build/dipdup.yml

api:
	cd cmd/api && go run . -c ../../build/dipdup.yml

build:
	cd cmd/indexer && go build -a -o ../../bin/indexer .
	cd cmd/api && go build -a -o ../../bin/api .

clean:
	rm -rf bin

compose:
	docker-compose up -d --build

lint:
	golangci-lint run

test:
	go test -p 8 -timeout 60s ./...

## adr: Generate ADR from template. Must set NUM and TITLE parameters.
adr:
	@echo "Generating ADR"
	@cp adr/adr-template.md adr/adr-$(NUM)-$(TITLE).md

mock:
	go generate ./internal/storage

api-docs:
	cd cmd/api && swag init --md markdown

.PHONY: init indexer api build clean compose lint test adr mock
