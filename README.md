[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fcelenium-io%2Fcelestia-indexer.svg?type=shield&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Fcelenium-io%2Fcelestia-indexer?ref=badge_shield&issueType=license)
[![Build Status](https://github.com/celenium-io/celestia-indexer/workflows/Build/badge.svg)](https://github.com/celenium-io/celestia-indexer/actions?query=branch%3Amaster+workflow%3A%22Build%22)
[![made_with golang](https://img.shields.io/badge/made_with-golang-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Coverage](https://github.com/celenium-io/celestia-indexer/wiki/coverage.svg)](https://raw.githack.com/wiki/celenium-io/celestia-indexer/coverage.html)

# Celestia Indexer | Celenium #

This is an indexing layer for Celestia DA written in Golang that operates on top of the [Celestia Full node](https://docs.celestia.org/nodes/consensus-full-node/) and stores data in a Postgres database.


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