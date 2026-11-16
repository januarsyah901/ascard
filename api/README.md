# api/ — Go backend for Bongkar Ilusi

A drop-in replacement for the Express server in `../server/`: same URLs, status codes, JSON
shapes, error messages, security headers and Postgres database. One process serves the JSON API
under `/api/*` and the built React SPA (`../dist`) on every other path.
The behaviour is specified in [`docs/SPEC.md`](docs/SPEC.md).

```
cmd/server     HTTP server (also: `server -healthcheck` for the container HEALTHCHECK)
cmd/seed       one-off CLI: create the admin user and sample articles
internal/      config, db, httpx, auth (JWT), authapi, articles, ratelimit, web (SPA, headers, CORS), server
contract/      black-box HTTP contract suite, runs against Go or Express
scripts/       contract.sh
```

## Configuration

| Env | Default | |
|---|---|---|
| `DATABASE_URL` | required | Prisma-style URLs (`?schema=public`) are accepted |
| `JWT_SECRET` | required, ≥ 32 chars | same secret as Express ⇒ existing tokens stay valid |
| `PORT` | `5001` | |
| `WEB_DIR` | `dist` | built SPA (Vite output) |
| `CORS_ORIGINS` | empty | comma list; empty ⇒ no CORS headers (same-origin only) |
| `TRUSTED_PROXIES` | loopback + private ranges | peers whose `X-Forwarded-For` is believed |
| `LOG_FORMAT` | `text` | `text` or `json` |

The server runs an idempotent `CREATE ... IF NOT EXISTS` schema check at startup; on the existing
Prisma-created database it is a no-op.

## Run locally

Needs Go (`export PATH=$HOME/.local/go/bin:$PATH`) and a Postgres.

```sh
npm ci && npm run build            # in the repo root, once: builds ../dist
cd api
export DATABASE_URL=postgresql://ascard:ascard@127.0.0.1:5432/ascard
export JWT_SECRET=change-me-to-at-least-32-characters-long
ADMIN_PASSWORD='a-strong-password' go run ./cmd/seed   # first run only
make run                           # http://localhost:5001 (WEB_DIR=../dist)
```

For frontend work, `npm run dev` (Vite) proxies `/api` to `localhost:5001`.

## Test

```sh
cd api
make lint        # gofmt -s check + go vet (incl. the contract build tag)
make test        # unit and handler tests, no database needed
make contract    # contract suite: builds the server, starts a throwaway Postgres in Docker
make contract-express   # the same suite against the Express reference
```

`scripts/contract.sh` takes `CONTRACT_DB_PORT`, `CONTRACT_APP_PORT` and
`CONTRACT_TEST_FLAGS` (e.g. `-run '^TestSPA'`) so several runs can share one machine.
It needs a built SPA in `../dist` (built with npm if missing).

CI (`.github/workflows/go.yml`) runs the gofmt check, `go vet`, `go test -race ./...` and the
contract suite on every push/PR touching `api/**`.

## Deploy (Docker)

`../Dockerfile.go` builds the SPA with Node, the Go binaries with `CGO_ENABLED=0 -trimpath`,
and ships them on `gcr.io/distroless/static-debian12:nonroot` (≈ 23 MB, no shell). The image runs
as non-root with `WEB_DIR=/app/dist PORT=5001 LOG_FORMAT=json APP_ENV=production` and a
`HEALTHCHECK` that runs `/app/server -healthcheck` (GET `/api/health`, healthy only on 200,
which needs the database to be reachable).

```sh
docker build -f Dockerfile.go -t ascard-app:go .
```

`../docker-compose.go.yml` is the production stack with the Go image. It keeps the compose project
name (`ascard`), container names, `gateway-net`/`internal` networks, env vars and the
`ascard_ascard-db-data` volume of `docker-compose.yml`, so Caddy and the database are untouched.

Cutover (from the repo root, with `.env.production` in place):

```sh
docker compose --env-file .env.production -f docker-compose.go.yml up -d --build
docker inspect -f '{{.State.Health.Status}}' ascard-app   # wait for "healthy"
curl -fsS https://<host>/api/health
```

Only `ascard-app` is recreated; `ascard-db` keeps running with the same volume. Rollback is the
same `up -d --build` with `-f docker-compose.yml` (the Express image). Seeding in the container:

```sh
docker compose -f docker-compose.go.yml run --rm --entrypoint /app/seed \
  -e ADMIN_PASSWORD='a-strong-password' ascard-app
```

Note: the Express image applies Prisma migrations at start; the Go image does not run Prisma.
The schema is unchanged, so this only matters for a future migration — apply it before the cutover.
