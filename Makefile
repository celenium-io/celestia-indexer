-include .env
export $(shell sed 's/=.*//' .env)

init:
	chmod +x init.dev.sh && ./init.dev.sh

indexer:
	cd cmd/indexer && go run . -c ../../configs/dipdup.yml

api:
	cd cmd/api && go run . -c ../../configs/dipdup.yml

quotes:
	cd cmd/quotes && go run . -c ../../configs/dipdup.yml

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

## Make green comment
gc:
	make lint && make test
	git commit -am "$(m)"

## adr: Generate ADR from template. Must set NUM and TITLE parameters.
adr:
	@echo "Generating ADR"
	@cp adr/adr-template.md adr/adr-$(NUM)-$(TITLE).md

generate:
	go generate -v ./internal/storage ./internal/storage/types ./pkg/node ./internal/binance

api-docs:
	cd cmd/api && swag init --md markdown -parseDependency --parseInternal --parseDepth 1

## Generate all code and docs
ga:
	make generate && make api-docs

check-licenses:
	go-licenses check ./... \
		--include_tests \
    	--allowed_licenses=MIT,Apache-1.0,Apache-1.1,Apache-2.0,BSD-2-Clause-FreeBSD,BSD-2-Clause-NetBSD,BSD-2-Clause,BSD-3-Clause-Attribution,BSD-3-Clause-Clear,BSD-3-Clause-LBNL,BSD-3-Clause,BSD-4-Clause,BSD-4-Clause-UC,BSD-Protection,ISC,LGPL-2.0,LGPL-2.1,LGPL-3.0,LGPLLR,MPL-1.0,MPL-1.1,MPL-2.0,Unlicense \
    	--ignore github.com/ethereum/go-ethereum \
    	--ignore github.com/regen-network/cosmos-proto \
    	--ignore github.com/modern-go/reflect2 \
    	--ignore golang.org/x/sys \
    	--ignore mellium.im/sasl \
    	--ignore github.com/klauspost/compress/zstd/internal/xxhash \
    	--ignore github.com/mattn/go-sqlite3 \
    	--ignore github.com/cespare/xxhash/v2 \
    	--ignore github.com/klauspost/reedsolomon \
    	--ignore github.com/klauspost/cpuid/v2 \
    	--ignore filippo.io/edwards25519/field \
    	--ignore github.com/golang/snappy \
    	--ignore golang.org/x/crypto/chacha20 \
    	--ignore github.com/pierrec/lz4/v4

test-api:
	newman run ./test/newman/tests.json -e ./test/newman/env.json

cover:
	go test ./... -coverpkg=./... -coverprofile ./coverage.out
	go tool cover -func ./coverage.out

license-header:
	update-license -path=./ -license=./HEADER

.PHONY: init indexer api build clean compose lint test adr mock api-docs check-licenses cover license-header
