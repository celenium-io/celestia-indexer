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

## adr: Generate ADR from template. Must set NUM and TITLE parameters.
adr:
	@echo "Generating ADR"
	@cp adr/adr-template.md adr/adr-$(NUM)-$(TITLE).md

.PHONY: indexer build lint test adr