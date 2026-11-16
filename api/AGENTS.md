# Working in api/ (Go rewrite of the Bongkar Ilusi backend)

Read `docs/SPEC.md` first. It is the source of truth. The Express reference is `../server/`
(read-only, never edit it). Behaviour, status codes, JSON shapes and **English** error
messages must match it exactly unless SPEC lists a deliberate change.

## Toolchain
- Go lives at `~/.local/go/bin`: run `export PATH=$HOME/.local/go/bin:$PATH` first.
- `make test` (unit + handler tests), `make vet`, `make fmt`.
- `make contract` builds the server and runs the HTTP contract suite against it with a
  throwaway Postgres in Docker. `make contract-express` runs the same suite against Express.
- Run only your part of the suite while other workers are still busy:
  `CONTRACT_TEST_FLAGS='-run <regexp>' ./scripts/contract.sh go`. Add `CONTRACT_MINT_TOKEN=1` when your
  tests need an admin token but `/api/auth/login` is not implemented in your branch (signs one directly).
- `gofmt -s -l .` must print nothing and `go vet ./...` must be clean before you commit.

## Module layout (one package per domain under internal/)
| File | Content |
|---|---|
| `routes.go` | `func Register(mux *http.ServeMux, d app.Deps)` — already stubbed; fill it in. Build the real repo from `d.Pool` here. |
| `validate.go` | Pure validation, no I/O. |
| `repo.go` | SQL via pgx. Define a `Repo` interface in the package; a `pgRepo` struct implements it. |
| `handler.go` | `Handler` struct holding the `Repo` interface; one method per endpoint. |
| `*_test.go` | Table-driven unit tests for validate.go; handler tests with `httptest` and an in-memory fake Repo. |

## Rules
- Routing: stdlib `net/http` patterns, e.g. `mux.Handle("PUT /api/articles/{id}", d.JWT.Guard(http.HandlerFunc(h.update)))`.
- Bodies: `fields, ok := httpx.DecodeObject(w, r); if !ok { return }`, then `fields.String("title")`,
  `fields.IsTrue("published")`. Do not decode input straight into structs.
- Responses: `httpx.JSON`, `httpx.Error`, `httpx.ServerError(w, r, "label", err)` for 500s.
- JSON models: camelCase tags matching the Prisma fields, **no `omitempty`**, nullable columns as
  pointers so they serialise as `null`, timestamps as `httpx.Time`, lists as `[]` never `null`.
- SQL: double-quote identifiers (`"Article"."createdAt"`); pass `gen_random_uuid()::text` for ids and
  `now()` for `"updatedAt"` yourself; map `pgx.ErrNoRows` with `db.NotFound`; detect duplicates with
  `db.IsUniqueViolation`. Use `r.Context()` for every DB call. No globals, no `init()` side effects.
- Stay inside the packages/files SPEC §3 assigns to you. Do not edit `internal/{app,auth,config,db,httpx,server}`,
  `contract/`, `scripts/`, `go.mod` (except `go/auth` adding `golang.org/x/crypto`), other workers'
  packages, or anything outside `api/` (except the deploy files SPEC §3.3 gives to `go/platform`).
  If you truly need a shared change, describe it in your final report instead of making it.
- Commits: conventional messages (`feat(articles): ...`, `test(auth): ...`), no AI/Claude attribution lines.
