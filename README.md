[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fcelenium-io%2Fcelestia-indexer.svg?type=shield&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Fcelenium-io%2Fcelestia-indexer?ref=badge_shield&issueType=license)
[![Build Status](https://github.com/celenium-io/celestia-indexer/workflows/Build/badge.svg)](https://github.com/celenium-io/celestia-indexer/actions?query=branch%3Amaster+workflow%3A%22Build%22)
[![made_with golang](https://img.shields.io/badge/made_with-golang-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Coverage](https://github.com/celenium-io/celestia-indexer/wiki/coverage.svg)](https://raw.githack.com/wiki/celenium-io/celestia-indexer/coverage.html)

# Celestia Indexer | Celenium #

This is an indexing layer for Celestia DA written in Golang that operates on top of the [Celestia Full node](https://docs.celestia.org/nodes/consensus-full-node/) and stores data in a Postgres database.

## Run

**Prerequisites:**

- Git
- [Docker](https://docs.docker.com/engine/install/)
- [Go 1.21.2](https://go.dev/doc/install) (for development and testing)

### Local run ###

Clone the repository:

```sh
git clone https://github.com/celenium-io/celestia-indexer.git
cd celestia-indexer
```

Create `.env` file and set up required environment variables:

```sh
cp .env.example .env
vim .env
``` 

> **Required environment variables:**
> 
> - `CELESTIA_DAL_API_URL` - uri for [Celestia Full Storage Node](https://docs.celestia.org/nodes/full-storage-node)
> - `CELESTIA_NODE_AUTH_TOKEN` - token with read access level for full storage node. You can get it from your running node instance by command `celestia full auth read`
> - `CELESTIA_NODE_URL` - uri to [Celestia Consensus Node](https://docs.celestia.org/nodes/consensus-node)
> - `POSTGRES_USER` - username for Postgres
> - `POSTGRES_PASSWORD` - password for Postgres
>

Build the Docker images for the indexer and API:

```sh
docker compose build
```

Start the services using Docker Compose:

```sh
docker compose up -d
```

This will start the indexer and API services as well as a Postgres database instance.
The services will be configured according to the `.env` file and the `docker-compose.yml` file in the repository.

≠
## Features ##

- [x] RPC node client
- [x] Rollbacks are handled
- [x] Database is partitioned for better performance
- [ ] Optional diagnostic mode for consistency checks


## Indexed entities ##

- Blocks
    - Transactions
    - Balance updates (block rewards, gov-triggered issuance/burn, other events)
    - Header
    - Stats
- Transactions
    - Details
    - Balance updates
- Blobs
    - Metadata
- Accounts
    - Balances
    - Stats
- Namespaces
    - Stats
- Summary
    - Stats


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fcelenium-io%2Fcelestia-indexer.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fcelenium-io%2Fcelestia-indexer?ref=badge_large)