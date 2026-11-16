# SPEC — Go rewrite of the Bongkar Ilusi backend

The Go server in `api/` replaces the Express server in `../server/` **as a drop-in**: same
URLs, status codes, JSON shapes, messages, headers and database. The frontend (`../src`)
does not change. The Express code on branch `legacy/express` (and in `../server/` until the
cutover) is the reference: read it, never edit it.

**Definition of done for the whole rewrite:** `api/scripts/contract.sh go` passes. The same
suite already passes against Express (`api/scripts/contract.sh express`).

## 1. Runtime

| Env | Default | Notes |
|---|---|---|
| `DATABASE_URL` | — required | Prisma-style URL accepted (`?schema=public`) |
| `JWT_SECRET` | — required, ≥ 32 chars | same secret as Express ⇒ old tokens stay valid |
| `PORT` | `5001` | |
| `WEB_DIR` | `dist` | built SPA (Vite output) |
| `CORS_ORIGINS` | empty | comma list; empty = no CORS headers at all |
| `TRUSTED_PROXIES` | loopback + private ranges | peers whose `X-Forwarded-For` is believed |
| `LOG_FORMAT` | `text` | `text` or `json` |

One process serves `/api/*` (JSON) and the SPA on every other path. At startup it runs
`db.EnsureSchema` (idempotent `CREATE ... IF NOT EXISTS`, a no-op on the live database).

## 2. Database (unchanged, created by Prisma)

```
"Admin"   id TEXT PK, username TEXT UNIQUE, password TEXT (bcrypt $2a/$2b), "createdAt" TIMESTAMP(3) default now
"Article" id TEXT PK, title TEXT UNIQUE, content TEXT, category TEXT NULL,
          published BOOLEAN default false, "createdAt" TIMESTAMP(3) default now, "updatedAt" TIMESTAMP(3) NOT NULL
```
- Quote every identifier. `id` has **no default**: insert `gen_random_uuid()::text`.
  `"updatedAt"` has **no default**: set `now()` on insert and update.
- Timestamps are UTC (`db.Open` pins the session time zone).

## 3. HTTP contract

All JSON errors are `{"message": "..."}` (`httpx.Error`). 500s use `httpx.ServerError`
("Internal server error"). Timestamps in responses use `httpx.Time` (`2026-11-11T09:21:00.123Z`).
Request bodies go through `httpx.DecodeObject` (non-JSON content type ⇒ empty object;
malformed JSON ⇒ 400 "Invalid JSON"; > 1 MB ⇒ 413).

Already done in the foundation: `GET /api/health`, JSON 404 "Not found" for unknown `/api/*`,
`auth.JWT` (`Sign`, `Verify`, `Guard`), request id / access log / panic recovery.

### 3.1 Auth — worker `go/auth` (packages `authapi`, `ratelimit`, `cmd/seed`)

**POST /api/auth/login** — rate limited, see below.
1. `username` and `password` must both be JSON strings ⇒ else **400** "Username and password are required".
2. Look up `"Admin"` by username; compare with bcrypt (`golang.org/x/crypto/bcrypt`, accepts the
   bcryptjs `$2b$` hashes). Unknown user or wrong password ⇒ **401** "Invalid credentials"
   (same message, no timing shortcut needed beyond bcrypt itself).
3. **200** `{"token": "<JWT.Sign(id)>", "admin": {"id": "...", "username": "..."}}`.

Rate limit (`ratelimit` package, in-memory, fixed window like express-rate-limit):
**5 requests per 15 minutes per client IP** (`d.ClientIP.Resolve(r)`), counting every login
request (success, failure, 400 alike). The 6th ⇒ **429** "Too many login attempts. Try again in 15 minutes."
Expired windows must be cleaned up (no unbounded map growth).

**GET /api/auth/me** — behind `d.JWT.Guard`. Load the admin by `auth.AdminID(ctx)`:
found ⇒ **200** `{"admin": {"id": "...", "username": "..."}}`; not found ⇒ **401** "Invalid token".

**cmd/seed** — one-off CLI replacing `prisma/seed.js`: reads `DATABASE_URL`, `ADMIN_USERNAME`
(default `admin`), `ADMIN_PASSWORD` (required, ≥ 12 chars, else exit 1). Runs `db.EnsureSchema`,
upserts the admin (insert if missing; never overwrite an existing admin's password), and inserts the
three sample articles from `../prisma/seed.js` by title if missing (`ON CONFLICT ("title") DO NOTHING`).
bcrypt cost 10.

`go/auth` is the only worker allowed to add a dependency: `golang.org/x/crypto`.

### 3.2 Articles — worker `go/articles` (package `articles`)

Article JSON — exactly these keys, in this order:
`{"id","title","content","category","published","createdAt","updatedAt"}`; `category` is `null` when absent.

Validation (`pickArticle` in `../server/routes/articles.js`), used by POST and PUT **before** any DB access:
- `title`: JSON string, `trim()` non-empty, length of the **untrimmed** value ≤ 200 ⇒ stored trimmed.
- `content`: JSON string, `trim()` non-empty, length ≤ 100 000 ⇒ stored **untrimmed**.
- `category`: a string that is non-empty after trim ⇒ trimmed and cut to the first 50 characters; anything else ⇒ `null`.
- `published`: `true` only if the JSON value is literally `true` (`Fields.IsTrue`), otherwise `false`.
- Any failure ⇒ **400** "Title and content are required".
- Length means characters as JavaScript counts them (UTF-16 code units). Counting runes is acceptable.

| Route | Guard | Behaviour |
|---|---|---|
| `GET /api/articles` | public | published only, `ORDER BY "createdAt" DESC`; `[]` when empty, never `null` |
| `GET /api/articles/admin/all` | `JWT.Guard` | all articles incl. drafts, same order |
| `GET /api/articles/{id}` | public | published article ⇒ 200; draft or missing ⇒ **404** "Article not found" |
| `POST /api/articles` | `JWT.Guard` | validate ⇒ insert ⇒ **201** article |
| `PUT /api/articles/{id}` | `JWT.Guard` | validate ⇒ replace title/content/category/published, set `"updatedAt"=now()` ⇒ 200; missing ⇒ 404 |
| `DELETE /api/articles/{id}` | `JWT.Guard` | **200** `{"message":"Article deleted successfully"}`; missing ⇒ 404 "Article not found" |

Duplicate title on POST or PUT (unique violation, `db.IsUniqueViolation`) ⇒ **409**
"An article with this title already exists". Ids are opaque strings (no UUID parsing; any id that
is not in the table is simply 404).

### 3.3 Platform — worker `go/platform` (package `web`, deploy files, CI)

**SPA** (`web.Register` on `"/"`, serving `d.Config.WebDir`):
- GET/HEAD of an existing file (e.g. `/assets/index-abc.js`, `/icon.png`) ⇒ the file, correct
  Content-Type, `Cache-Control: public, max-age=31536000, immutable`.
- GET/HEAD of `/assets/<missing>` ⇒ **404** (deliberate fix: Express returned index.html with 200,
  which breaks stale browser tabs after a deploy with a confusing MIME error).
- Any other GET/HEAD path ⇒ `index.html` with `Cache-Control: no-cache`, `Content-Type: text/html; charset=utf-8`.
- Other methods ⇒ 404 JSON "Not found". Never serve directory listings; never escape `WebDir`
  (`..` traversal). If `WebDir` has no `index.html`, log a warning at startup and answer 404.

**SecurityHeaders** (helmet defaults, as Express sends them today) on **every** response:
```
Content-Security-Policy: default-src 'self';base-uri 'self';font-src 'self' https://fonts.gstatic.com;form-action 'self';frame-ancestors 'self';img-src 'self' data: https:;object-src 'none';script-src 'self';script-src-attr 'none';style-src 'self' 'unsafe-inline' https://fonts.googleapis.com;upgrade-insecure-requests;connect-src 'self'
Cross-Origin-Opener-Policy: same-origin
Cross-Origin-Resource-Policy: same-origin
Origin-Agent-Cluster: ?1
Referrer-Policy: no-referrer
Strict-Transport-Security: max-age=31536000; includeSubDomains
X-Content-Type-Options: nosniff
X-DNS-Prefetch-Control: off
X-Download-Options: noopen
X-Frame-Options: SAMEORIGIN
X-Permitted-Cross-Domain-Policies: none
X-XSS-Protection: 0
```

**CORS(origins)**: no headers unless the request `Origin` is in the list. For a listed origin:
`Access-Control-Allow-Origin: <origin>`, `Vary: Origin`; preflight `OPTIONS` with
`Access-Control-Request-Method` ⇒ 204 with `Access-Control-Allow-Methods: GET,HEAD,PUT,PATCH,POST,DELETE`
and the requested headers echoed in `Access-Control-Allow-Headers`.

**Deploy files** (repo root, used at the cutover — do not touch the live stack):
- `Dockerfile.go`: stage 1 `node:22-alpine` builds the SPA (`npm ci && npm run build`); stage 2
  `golang:1.27-alpine` builds `cmd/server` and `cmd/seed` with `CGO_ENABLED=0 -trimpath -ldflags="-s -w"`;
  final `gcr.io/distroless/static-debian12:nonroot` with `/app/server`, `/app/seed`, `/app/dist`,
  `ENV WEB_DIR=/app/dist PORT=5001 LOG_FORMAT=json APP_ENV=production`, `EXPOSE 5001`.
  Distroless has no shell/wget: add a `-healthcheck` flag to `cmd/server` (GET `http://127.0.0.1:$PORT/api/health`,
  exit 0/1) and use it in `HEALTHCHECK CMD ["/app/server","-healthcheck"]`. You own that flag in `cmd/server/main.go`
  (keep the change minimal).
- `docker-compose.go.yml`: same as `docker-compose.yml` but builds `Dockerfile.go` (image `ascard-app:go`).
  Same container name, network, env and DB volume, so the cutover is `docker compose -f docker-compose.go.yml up -d`.
- `.github/workflows/go.yml`: on push/PR touching `api/**` — gofmt check, `go vet`, `go test ./...`,
  and `api/scripts/contract.sh go` (GitHub runners have Docker and Node).
- `api/README.md`: how to run, test and deploy.

## 4. Work split and rules

See `../AGENTS.md`-style rules in `api/AGENTS.md`. Each worker owns only its packages/files listed
above. Shared foundation (`config`, `db`, `httpx`, `auth`, `app`, `server`, `contract`, `scripts`)
is read-only for workers: describe needed changes in the final report instead.
