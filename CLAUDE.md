# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# Celestia Indexer — CLAUDE.md

## Project Overview

Go-based blockchain indexer + REST API for the Celestia Data Availability (DA) blockchain. Indexes blocks, transactions, messages, events, blobs, namespaces, validators, governance, IBC, Hyperlane, and rollups into PostgreSQL (TimescaleDB) and exposes them via public Echo HTTP/WebSocket API and a private admin API.

## Architecture

```
cmd/
  indexer/              # Core indexer daemon
  api/                  # Public REST API (port 9876)
    handler/            # Echo handlers (one file per entity)
    handler/responses/  # DTO structs for API responses
  private_api/          # Admin API (port 9877)
  celestials/           # Off-chain Celestials data indexer
pkg/
  indexer/              # Core indexer pipeline
    receiver/           # Fetches blocks from CometBFT RPC/API/WS
    parser/             # Decodes raw blocks, txs, messages, events
    storage/            # Saves parsed data to DB in one DB transaction
    rollback/           # Handles chain reorganizations
    genesis/            # Handles genesis block separately
    decode/context/     # Context object accumulates all parsed entities (passed parser → storage)
    config/             # Indexer config structures
  node/
    rpc/                # CometBFT RPC client
    api/                # Node REST API client
    dal/                # DAL API client (blob retrieval)
  types/                # pkg-level domain types (Level, etc.)
internal/
  storage/              # Domain model structs + storage interfaces (IXxx)
    id.go               # Deterministic ID generation (height<<24 | position)
    transaction.go      # Block transaction: Transaction, domain *Tx interfaces, Insert, TxRepos
    postgres/           # Bun ORM implementations of all interfaces
      scopes.go         # Reusable query filters and pagination helpers
      transaction.go    # Transaction core: BeginTransaction, Insert, RollbackByHeight, NewTxRepos, block writes
      transaction_*.go  # Domain writes/rollbacks: account, namespace, staking, gov, ibc, hyperlane, zkism, rollup
      core.go           # DB init, migrations, hypertables, enums, indexes
      migrations/       # Bun migrations (named by date)
    types/              # Enums (MsgType, EventType, ModuleType, etc.)
  blob/                 # Blob handling utilities
  pool/                 # sync.Pool wrappers for reusing slices
  currency/             # Currency utilities
  stats/                # Statistics calculations
database/
  functions/            # PostgreSQL functions (materialized view refresh)
  views/                # Materialized views for analytics (minute/hour/day/...)
configs/
  dipdup.yml            # YAML config with ${ENV_VAR:-default} substitution
```

**Indexer pipeline:** CometBFT RPC/WS → Receiver → Parser → Storage module → PostgreSQL

**Module wiring** (in `pkg/indexer/indexer.go`): modules connect via named inputs/outputs using `module.AttachTo(source, outputName, inputName)`. Every module has a `StopOutput` that feeds into the stopper.

## Key Libraries

| Purpose | Library |
|---------|---------|
| HTTP | `github.com/labstack/echo/v4` |
| ORM | `github.com/uptrace/bun` + `github.com/jackc/pgx/v5` (via `pgx/v5/stdlib`) |
| Blockchain | `github.com/celestiaorg/celestia-app/v10`, `github.com/cometbft/cometbft` |
| Cosmos | `github.com/cosmos/cosmos-sdk`, `github.com/cosmos/ibc-go/v8` |
| Cache | `github.com/valkey-io/valkey-go` |
| Logging | `github.com/rs/zerolog` |
| Validation | `github.com/go-playground/validator` |
| Errors | `github.com/pkg/errors` |
| Mocks | `go.uber.org/mock/mockgen` |
| Swagger | `github.com/swaggo/swag` |
| Indexer SDK | `github.com/dipdup-net/indexer-sdk` |
| JSON | `github.com/bytedance/sonic` |
| Profiling | `github.com/grafana/pyroscope-go` |
| Sentry | `github.com/getsentry/sentry-go` |

## Commands

```bash
make init         # chmod + run init.dev.sh (first-time dev setup)
make indexer      # go run ./cmd/indexer -c ./configs/dipdup.yml
make api          # go run ./cmd/api -c ./configs/dipdup.yml
make private_api  # go run ./cmd/private_api -c ./configs/dipdup.yml
make celestials   # go run ./cmd/celestials -c ./configs/dipdup.yml
make build        # build all binaries to /bin
make test         # go test -p 8 -timeout 120s ./...
make generate     # go generate (mocks + enums for storage, blob, node, gas packages)
make api-docs     # swag init (regenerate Swagger)
make ga           # generate + api-docs
make lint         # golangci-lint
make gc           # lint → test → commit
make compose      # docker compose up --build
make cover        # generate coverage report
make license-header  # add/update SPDX headers across all files
```

**Run a single test or package:**
```bash
go test ./internal/storage/postgres/... -run TestBlockByHeight -v
go test ./cmd/api/handler/... -timeout 30s
```

## Configuration

YAML config with `${ENV_VAR:-default}` substitution (`configs/dipdup.yml`):

```
# Datasources
CELESTIA_NODE_RPC_URL / CELESTIA_NODE_API_URL / CELESTIA_NODE_WS_URL
CELESTIA_DAL_URL / CELESTIALS_URL / CELENIUM_BLOBS_URL

# Database
POSTGRES_HOST / PORT / USER / PASSWORD / DB / MAX_OPEN_CONNECTIONS

# API
API_HOST / API_PORT / API_RATE_LIMIT / WEBSOCKET_ENABLED
PRIVATE_API_HOST / PRIVATE_API_PORT

# Cache
CACHE_URL / CACHE_TTL

# Indexer
INDEXER_START_LEVEL / INDEXER_SCRIPTS_DIR / NETWORK
```

## Migrations

Migrations live in `internal/storage/postgres/migrations/`, named by date (e.g. `20260320000001_description.go`). Each file registers exactly one migration via `init()`:

```go
func init() {
    Migrations.MustRegister(upXxx, downXxx)
}
```

**Both `up` and `down` functions are required.** Never leave `down` as a stub or `TODO`.

- `up` — applies the change (e.g. `ALTER TYPE ... ADD VALUE`, `CREATE INDEX`, `ALTER TABLE ... ADD COLUMN`)
- `down` — fully reverts it. PostgreSQL-specific notes:
  - Removing an added enum value: `DELETE FROM pg_enum WHERE enumlabel = '...' AND enumtypid = (SELECT oid FROM pg_type WHERE typname = '...')` — only safe if no rows use that value
  - Dropping a column: `ALTER TABLE ... DROP COLUMN IF EXISTS ...`
  - Dropping an index: `DROP INDEX IF EXISTS ...`

## Storage Patterns

All storage files in `internal/storage/postgres/`. Each entity has its own file (`address.go`, `block.go`, `tx.go`, etc.).

**Typical query pattern** — subquery for filters, outer query for JOINs:
```go
func (a *Address) ByHash(ctx context.Context, hash []byte) (address storage.Address, err error) {
    addressQuery := a.DB().NewSelect().
        Model((*storage.Address)(nil)).
        Where("hash = ?", hash)

    err = a.DB().NewSelect().
        TableExpr("(?) as address", addressQuery).
        ColumnExpr("address.*").
        ColumnExpr("celestial.id as celestials__id, celestial.image_url as celestials__image_url").
        ColumnExpr("balance.currency as balance__currency, balance.spendable as balance__spendable").
        Join("left join balance on balance.id = address.id").
        Join("left join celestial on celestial.address_id = address.id and celestial.status = 'PRIMARY'").
        Scan(ctx, &address)
    return
}
```

**Joined relation columns** use `__` separator: `"celestial.id AS celestials__id"` maps to `Address.Celestials.Id`.

**Pagination helpers** (`scopes.go`):
- `limitScope(q, limit)` — clamps 1–100, default 10
- `sortScope(q, field, order)` — single field sort
- `txFilterWithoutLimit(q, fltrs)` — sort by `time, id` (time-series ordering)
- Message type filtering uses bitmask: `bit_count(message_types & ?::bit(115)) > 0`

**Repository constructors** take `bun.IDB` (`NewTable` from indexer-sdk does too), so the same repository runs on the pool (`strg.Connection().DB()`) or inside a transaction (`tx.Tx()`). Query through `x.DB()` / the `bun.IDB` field, never through `*database.Bun`, or the query escapes the transaction.

**Block transaction** (`internal/storage/transaction.go`, implementation in `postgres/transaction*.go`):
```go
tx, _ := postgres.BeginTransaction(ctx, module.storage)
defer tx.Close(ctx)
repos := module.reposFactory(tx) // postgres.NewTxRepos: repositories bound to tx
state, _ := repos.State.ByName(ctx, name) // reads see the tx's own writes
// writes via tx.*, tx.Flush() — then tx.HandleError() on failure
```
- `Transaction` holds **writes and rollbacks only**. It embeds `sdk.Transaction` and the domain interfaces `BlockTx`, `AccountTx`, `NamespaceTx`, `StakingTx`, `GovTx`, `IbcTx`, `HyperlaneTx`, `ZkIsmTx`, `RollupTx`. Each domain interface embeds `Inserter`.
- **Reads** inside the transaction go through `storage.TxRepos` (repository interfaces bound to the tx). Do not add getters to `Transaction`; add the method to the entity repository and, if needed, a field to `TxRepos` **and** to `postgres.NewTxRepos` (a missing field is `nil` at runtime).
- **Plain inserts**: `storage.Insert(ctx, tx, items...)`, a type-safe wrapper over `Inserter.Insert(ctx, *[]T)`. No new `SaveFoo` method is needed unless the insert has `ON CONFLICT`, COPY, counting or per-row logic.
- **Plain rollbacks** (`DELETE ... WHERE height = ?`): add the model to the `tx.RollbackByHeight(ctx, height, (*storage.Foo)(nil), ...)` list in `pkg/indexer/rollback/rollback.go`. The compiler will not remind you. Tables with extra rollback logic (grants, validators, IBC counters, …) keep dedicated `RollbackFoo` methods.
- **Narrow parameters**: `save*`/`rollback*` helpers take the narrowest type they need (`storage.IbcTx`, `storage.IValidator`, …); keep `storage.Transaction` for orchestrators and cross-domain functions.
- **TimescaleDB**: an `UPDATE` of a hypertable (`signal_version`, …) with subqueries in `WHERE` can fail with `variable not found in subplan target list`. Compute those values in separate queries first.
- **Indexer reads must not reuse API list methods blindly**: `limitScope` silently replaces limits above 100 with 10, and API methods add joins/relations. Methods used by the indexer (paging with `sdkSync.Paginate`, id lookups on hot paths) get their own lean query: explicit `ORDER BY`, no clamp, only needed columns (e.g. `IVote.ListByProposal`, `IHLIGP.IdByHash`).

**Deterministic IDs** (`internal/storage/id.go`): `tx` and `message` IDs are computed at parse time as `height<<24 | position` (5 bytes height + 3 bytes position). This removes autoincrement sequences for those tables. Genesis block (height=0) uses `position+1` to avoid zero IDs. Use `idFromHeightAndPosition(height, position)` when assigning IDs in parsers.

**Decode context** (`pkg/indexer/decode/context/`): All parsed entities (Messages, Events, Namespaces, AddressMessages, Grants, IbcClients, IbcChannels, HlMailboxes, BlobLogs, etc.) are accumulated into `*Context` during parsing. The storage module then reads from `dCtx.*` fields instead of reconstructing them. Use `ctx.AddMessage()`, `ctx.AddEvents()`, `ctx.AddSignal()`, etc.

**Storage module split** (`pkg/indexer/storage/`): The monolithic `message.go` logic is now split into domain-specific files: `blob.go`, `forwarding.go`, `grants.go`, `hyperlane.go`, `ibc.go`, `signal.go`, `vesting.go`, `zkism.go`. Each file handles saving its respective entity type from the decode context.

**Bulk COPY** (`transaction.go`): `SaveTransactions`, `SaveMessages`, `SaveBlobLogs`, `SaveMsgAddresses` now use `pg.SaveBulkWithCopy` (PostgreSQL COPY protocol) instead of INSERT with `RETURNING id`, since IDs are pre-computed.

## Signals and Upgrades (x/signal)

The tally mirrors `celestia-app/x/signal`. It is built only from chain data, with no node queries. Code: `pkg/indexer/storage/signal.go`, `SignalVersion.Latest`, `FixSignalsPower`.

- **Power snapshot** (`prepareSignalRound` → `takePowerSnapshot`): the power of bonded validators, read at the start of `processBlockInTransaction` before any validator write. It is taken only if the block can tally: pending versions, a signal above the current version, `MsgTryUpgrade`, or an applied upgrade. `PendingVersions` is read there too, since only `processSignalModule` writes `upgrade`. This is the power x/signal sees during the block (left by the previous EndBlock). A validator created in the block is not in it.
- **Tally** (`signalTally`, in Go): each validator's latest signal since the last applied upgrade (`ISignalVersion.Latest`; `ResetTally` wipes older ones), weighted by the snapshot. The total is the snapshot's sum. The threshold is `signalThreshold` = `ceil(signal.Threshold(v) × total)`, compared with `>=`. `tryUpgrade` takes its candidates from these signals, as `TallyVotingPower` does.
- **Current version** is `block.VersionApp`, the header version x/signal uses, not `state.Version`, which belongs to the previous block.
- **Order in `processSignalModule`** (called at the end of `processBlockInTransaction`, after `saveValidators`, so new validators have ids): `FixSignalsPower` (upgrade applied) → `setUpgradeApplied` → `saveSignals` → `recountUpgrades` → `tryUpgrade`.
- **`recountUpgrades` runs on every block** for `IUpgrade.PendingVersions(currentVersion)`: versions above the current one, not applied, with no `MsgTryUpgrade` (`end_height = 0`). It writes through `UpdateUpgradeTally`. `SaveUpgrades` skips zero fields, so it cannot record a lost quorum.
- **`waiting_upgrade`** means the quorum is reached. It goes back to `processing` if the quorum is lost before `MsgTryUpgrade`.
- **`signal_version.voting_power`**: `NULL` while the round is open (the API shows the current `validator.power`), the counted power after `FixSignalsPower`, or `0` if the signal was not counted. The API multiplies signal power by 10^6, because the field has always been in utia. Upgrade `voting_power`/`voted_power` are consensus power, with no multiplication.
- **`validator.version`** is the version of the validator's latest signal, not the max. Nothing in the tally reads it.
- **`expected_upgrade_height`** = `end_height + appconsts.GetUpgradeHeightDelay(chainID)`, set on `MsgTryUpgrade` and written once. If no `MsgTryUpgrade` was seen, it is `applied_at_level − 1`, set when the upgrade is applied. `NULL` until then (`nullzero`). The actual switch is `applied_at_level` (scheduled height + 1).
- **`MsgSignalVersion` / `MsgTryUpgrade` inside `MsgExec`** are parsed from events in `pkg/indexer/parser/events/exec.go`.
- **Rollback**: `RollbackUpgrades(height)` runs before `RollbackByHeight`, because it needs the block's signals. It lowers `signals_count`, undoes `MsgTryUpgrade` (`end_height = height`) and the applied upgrade (`applied_at_level = height`), and reopens the round (signals back to `NULL`) when that block closed it. `rollbackBlock` also restores `state.Version` from the new last block. Without that, a re-indexed upgrade block would not count as a version change.

## API Handler Pattern

```go
// 1. Handler struct holds injected storage interfaces
type BlockHandler struct {
    block      storage.IBlock
    blockStats storage.IBlockStats
    events     storage.IEvent
    // ...
}

// 2. Request struct with Echo binding tags + validator tags
type getBlockRequest struct {
    Height types.Level `param:"height" validate:"min=0"`
    Stats  bool        `query:"stats"  validate:"omitempty"`
}

// 3. Swagger annotations above every handler
// @Summary     Get block info
// @Tags        block
// @ID          get-block
// @Param       height path integer true "Block height" minimum(1)
// @Produce     json
// @Success     200 {object} responses.Block
// @Failure     400 {object} Error
// @Router      /v1/blocks/{height} [get]
func (h *BlockHandler) Get(c echo.Context) error {
    req, err := bindAndValidate[getBlockRequest](c)
    if err != nil { return badRequestError(c, err) }

    block, err := h.block.ByHeight(c.Request().Context(), req.Height)
    if err != nil { return handleError(c, err, h.block) }

    return c.JSON(http.StatusOK, responses.NewBlock(block))
}
```

**Helper functions** (`handler/` package):
- `bindAndValidate[T](c)` — generic bind + validate
- `badRequestError(c, err)` / `handleError(c, err, storage)` — consistent error responses
- `returnArray(c, arr)` — returns `[]` not `null` for empty slices
- `StringArray` — comma-separated query param `?types=val1,val2`

## Indexer Module Pattern

Each pipeline module embeds `modules.BaseModule`, has named string constants for inputs/outputs:

```go
const (
    InputName  = "data"
    StopOutput = "stop"
)

type Module struct {
    modules.BaseModule
    storage     sdk.Transactable
    constants   storage.IConstant
    validators  storage.IValidator
    // ...
}

func NewModule(pg postgres.Storage, cfg config.Config, ...) (*Module, error) {
    m := &Module{BaseModule: modules.New("storage"), ...}
    m.CreateInput(InputName)
    m.CreateOutput(StopOutput)
    return m, nil
}

func (m *Module) Start(ctx context.Context) {
    m.G.GoCtx(ctx, m.listen)
}

func (m *Module) listen(ctx context.Context) {
    input := m.MustInput(InputName)
    for {
        select {
        case <-ctx.Done():
            return
        case msg, ok := <-input.Listen():
            if !ok {
                m.MustOutput(StopOutput).Push(struct{}{})
                return
            }
            // process msg...
        }
    }
}

func (m *Module) Close() error { m.G.Wait(); return nil }
```

**Module wiring** in `pkg/indexer/indexer.go`:
```go
// receiver → parser → storage (data flow)
p.AttachTo(r, receiver.OutputName, parser.InputName)
s.AttachTo(p, parser.OutputName, storage.InputName)
// All modules → stopper (shutdown flow)
stopperModule.AttachTo(r, receiver.StopOutput, stopper.InputName)
```

## Adding a New Entity (Checklist)

1. `internal/storage/` — add model struct + filter struct + interface `IFoo`
2. `internal/storage/postgres/foo.go` — implement queries using subquery+JOIN pattern
3. `internal/storage/postgres/core.go` — register in `Storage` struct, create hypertable if time-series
4. `internal/storage/postgres/index.go` — add indexes
5. Saving/rollback: plain inserts use `storage.Insert`, plain rollbacks go into the `RollbackByHeight` list in `pkg/indexer/rollback/rollback.go`; only non-trivial ones get a method in the matching domain `*Tx` interface (`internal/storage/transaction.go`) and `postgres/transaction_<domain>.go`. Reads needed by the indexer go to the repository + `TxRepos`/`NewTxRepos`
6. Mock: add `//go:generate` directive, run `make generate`
7. Parser/decode: add parsing logic; assign deterministic ID via `idFromHeightAndPosition`; accumulate into `dCtx` via `ctx.AddXxx()` methods
8. `pkg/indexer/storage/foo.go` — add `saveFoo(ctx, tx, dCtx.Foos, ...)` function, with `tx` typed as the domain interface (e.g. `storage.IbcTx`); call it from `processBlockInTransaction` in `storage.go`
9. `cmd/api/handler/foo.go` — handler with Swagger annotations
10. Register routes in `cmd/api/main.go`
11. Run `make api-docs`

## Key Conventions

- `zerolog` only for logging — never `fmt.Print` in production paths
- `errors.Wrap(err, "context")` from `github.com/pkg/errors`
- Storage interfaces only — don't use concrete postgres types outside `internal/`
- WebSocket notifications are skipped during initial sync (`time.Since(block.Time) > time.Hour`)
- `pool.New(func() []T)` — use `internal/pool` for reusing slices in hot paths
- Message types use bitmask (bit vectors of size 115) for efficient multi-type filtering
- Enum types are code-generated via `make generate` — never edit `*_enum.go` files manually
- Active linters to watch: `zerologlint`, `musttag`, `gosec`, `containedctx`
- SPDX license headers required on all new files: `// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>` + `// SPDX-License-Identifier: MIT` (run `make license-header` to apply automatically)
- JSON marshaling uses `github.com/bytedance/sonic` (faster than `encoding/json`)
- Keep comments short — 1-2 lines max. No multi-paragraph doc comments; explain the non-obvious "why" tersely, not the "what"

## Entity Types Overview

Key entities indexed (57 total storage types):

| Category | Entities |
|----------|---------|
| Blockchain | Block, Tx, Message, Event |
| DA / Blobs | Namespace, BlobLog |
| Validators | Validator, Delegation, Redelegation, Undelegation, Jail |
| Accounts | Address, Balance, Vesting, Grant |
| Governance | Proposal, Vote, Signal |
| IBC | IBC transfers, channels |
| Hyperlane | Hyperlane transfers |
| Rollups | Rollup, RollupProvider |
| Analytics | BlockStats, NamespaceStats, ValidatorStats, etc. |
| Other | Constant, State, Forwarding, Celestial |

## Testing

- **Contexts in tests: prefer `t.Context()`** (`s.T().Context()` in suites, `b.Context()` in benchmarks) over `context.Background()`, including helpers (pass `t` in). The one exception is `t.Cleanup` / `s.T().Cleanup`: the test context is already canceled when cleanup runs, so use `context.Background()` there (e.g. `psqlContainer.Terminate`). A transaction begun on a canceled context is rolled back by `database/sql`, so defer `tx.Rollback`/`tx.Close` inside the test instead of registering them as cleanup.
- Do not run `golangci-lint --fix`: even with `--enable-only`, it applies `usetesting` rewrites across the repo, cleanup contexts included.
- Mocks are auto-generated in `mock/` subdirectories — never edit manually. Exceptions:
  - `internal/storage/transaction.go` uses mockgen package mode (`. Transaction`), so only `MockTransaction` is generated, not the domain `*Tx` interfaces.
  - `internal/storage/mock/tx_repos.go` is hand-written: `mock.NewTxRepos(ctrl)` returns typed repository mocks (`repos.Validators.EXPECT()`), `repos.Repos()` returns `storage.TxRepos`. Keep it in sync with `storage.TxRepos`.
- Mocks do not apply `limitScope` or SQL semantics: a new indexer read method needs a DB test (fixtures reload per test in `TransactionTestSuite`; use `NewTxRepos(tx)` to read uncommitted data)
- DB integration tests spin up a real TimescaleDB Docker container via testcontainers — **Docker must be running**
- `testfixtures` for DB integration tests (`test/` directory); **avoid `0x`-prefixed strings in YAML fixtures** — testfixtures interprets them as hex-encoded bytea, causing invalid UTF-8 errors when inserting into `text` columns
- Newman collection for API tests: `make test-api`
- Run `make test` before committing
- Coverage: `make cover`
