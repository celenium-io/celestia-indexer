[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fcelenium-io%2Fcelestia-indexer.svg?type=shield&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Fcelenium-io%2Fcelestia-indexer?ref=badge_shield&issueType=license)
[![Build Status](https://github.com/celenium-io/celestia-indexer/workflows/Build/badge.svg)](https://github.com/celenium-io/celestia-indexer/actions?query=branch%3Amaster+workflow%3A%22Build%22)
[![made_with golang](https://img.shields.io/badge/made_with-golang-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Coverage](https://github.com/celenium-io/celestia-indexer/wiki/coverage.svg)](https://raw.githack.com/wiki/celenium-io/celestia-indexer/coverage.html)
[![Latest release](https://img.shields.io/github/v/release/celenium-io/celestia-indexer.svg)](https://github.com/celenium-io/celestia-indexer/releases)

# Celestia Indexer | Celenium

Go-based blockchain indexer and REST API for the [Celestia](https://celestia.org/) Data Availability (DA) network. It reads data from a Celestia full node, decodes blocks, transactions, messages, events, and blobs, and stores them in a PostgreSQL (TimescaleDB) database. A public Echo HTTP/WebSocket API and a private admin API serve the indexed data.

## Architecture

```
cmd/
  indexer/          # Core indexer daemon
  api/              # Public REST API (port 9876)
  private_api/      # Admin REST API (port 9877)
  celestials/       # Off-chain Celestials metadata indexer
pkg/
  indexer/
    receiver/       # Fetches blocks from CometBFT RPC / REST / WebSocket
    parser/         # Decodes raw blocks, txs, messages, events, blobs
    storage/        # Persists parsed data inside a single DB transaction
    rollback/       # Handles chain reorganizations
    genesis/        # Handles the genesis block
internal/
  storage/          # Domain model structs and storage interfaces
    postgres/       # Bun ORM implementations, migrations, hypertables
node/
  rpc/              # CometBFT RPC client
  api/              # Node REST API client
  dal/              # DAL API client (blob retrieval)
```

**Indexer pipeline:** CometBFT RPC/WS → Receiver → Parser → Storage → PostgreSQL

## Prerequisites

- [Docker](https://docs.docker.com/engine/install/) and Docker Compose
- [Go 1.26+](https://go.dev/doc/install) (for local development and testing)

## Quick Start

Clone the repository:

```sh
git clone https://github.com/celenium-io/celestia-indexer.git
cd celestia-indexer
```

Copy the example env file and fill in the required values:

```sh
cp .env.example .env
$EDITOR .env
```

**Required environment variables:**

| Variable | Description |
|---|---|
| `CELESTIA_NODE_URL` | URI to a [Celestia Consensus Node](https://docs.celestia.org/nodes/consensus-node) RPC endpoint |
| `CELESTIA_NODE_API_URL` | URI to the Consensus Node REST API |
| `CELESTIA_DAL_API_URL` | URI to a [Celestia Full Storage Node](https://docs.celestia.org/nodes/full-storage-node) |
| `CELESTIA_NODE_AUTH_TOKEN` | Read-access auth token for the DAL node (`celestia full auth read`) |
| `POSTGRES_USER` | PostgreSQL username |
| `POSTGRES_PASSWORD` | PostgreSQL password |
| `HYPERLANE_NODE_URL` | Hyperlane node URL for cross-chain message indexing |

Build and start all services:

```sh
docker compose up -d
```

This starts the indexer, public API, private API, TimescaleDB, and Valkey cache.

## Services

| Service | Image | Default port |
|---|---|---|
| `indexer` | `ghcr.io/celenium-io/celestia-indexer` | — |
| `api` | `ghcr.io/celenium-io/celestia-indexer-api` | 9876 |
| `private-api` | `ghcr.io/celenium-io/celestia-indexer-private-api` | 9877 |
| `db` | `timescale/timescaledb-ha:pg15` | 5432 |
| `cache` | `valkey/valkey:8` | 6379 |

## Indexed Entities

| Category | Entities |
|---|---|
| Blockchain | Block, BlockStats, BlockSignature, Transaction, Message, Event |
| DA / Blobs | Namespace, NamespaceMessage, BlobLog |
| Validators | Validator, ValidatorStats, Delegation, Redelegation, Undelegation, Jail, StakingLog |
| Accounts | Address, Balance, Grant, Vesting, Forwarding |
| Governance | Proposal, Vote, Signal |
| IBC | IbcClient, IbcConnection, IbcChannel, IbcTransfer |
| Hyperlane | HlMailbox, HlToken, HlTransfer, HlIgp, HlGasPayment |
| Rollups | Rollup, RollupProvider |
| Infrastructure | Constant, State, DenomMetadata, Upgrade |
| Off-chain | Celestial (identity metadata) |

## Development

First-time setup:

```sh
make init
```

Run individual services locally:

```sh
make indexer      # go run ./cmd/indexer -c ./configs/dipdup.yml
make api          # go run ./cmd/api -c ./configs/dipdup.yml
make private_api  # go run ./cmd/private_api -c ./configs/dipdup.yml
make celestials   # go run ./cmd/celestials -c ./configs/dipdup.yml
```

Other useful commands:

```sh
make build        # build all binaries to /bin
make test         # run all tests (requires Docker for DB integration tests)
make generate     # regenerate mocks and enums
make api-docs     # regenerate Swagger docs (swag init)
make lint         # run golangci-lint
make cover        # generate HTML coverage report
make compose      # docker compose up --build
```

Run a specific test:

```sh
go test ./internal/storage/postgres/... -run TestBlockByHeight -v
go test ./cmd/api/handler/... -timeout 30s
```

## Configuration

The YAML config at `configs/dipdup.yml` uses `${ENV_VAR:-default}` substitution. All settings can be overridden via environment variables or the `.env` file. Key options beyond the required variables above:

| Variable | Default | Description |
|---|---|---|
| `INDEXER_START_LEVEL` | `1` | First block to index |
| `INDEXER_BLOCK_PERIOD` | `15` | Polling interval (seconds) |
| `NETWORK` | — | Network identifier |
| `API_RATE_LIMIT` | `20` | Requests per second per IP |
| `API_WEBSOCKET_ENABLED` | `true` | Enable WebSocket notifications |
| `CACHE_URL` | — | Valkey/Redis connection URL |
| `CACHE_TTL` | — | Cache TTL (seconds) |
| `SENTRY_DSN` | — | Optional Sentry DSN for error tracking |

## Features

- [x] Full block, transaction, and message indexing
- [x] Blob and namespace tracking
- [x] Validator, staking, and governance indexing
- [x] IBC transfer and channel indexing
- [x] Hyperlane cross-chain message indexing
- [x] Chain rollback handling
- [x] TimescaleDB hypertables for time-series performance
- [x] WebSocket real-time notifications
- [x] Public REST + WebSocket API with Swagger docs
- [x] Private admin API
- [x] Valkey/Redis response cache
- [x] Deterministic IDs (no autoincrement sequences for tx/messages)

## License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fcelenium-io%2Fcelestia-indexer.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fcelenium-io%2Fcelestia-indexer?ref=badge_large)
