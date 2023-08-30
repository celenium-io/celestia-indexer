# Celestia Indexer #

This is an indexing layer for Celestia DA written in Golang that operates on top of the [Celestia Full node](https://docs.celestia.org/nodes/consensus-full-node/) and stores data in a Postgres database.


## Features ##

- [ ] RPC and WebSocket node reader
- [ ] Rollbacks are handled
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
