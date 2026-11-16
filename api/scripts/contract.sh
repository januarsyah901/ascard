#!/usr/bin/env bash
# Run the HTTP contract suite (api/contract) against a real server.
#
#   api/scripts/contract.sh go        # the Go rewrite (default)
#   api/scripts/contract.sh express   # the Express reference in ../server
#
# It starts a throwaway Postgres container, creates the schema the way production has it
# (the Prisma migration), starts the chosen server on a free port and runs the suite.
# Needs: docker, go, node (express only) and a built SPA in ../dist (built if missing).
set -euo pipefail

target="${1:-go}"
api_dir="$(cd "$(dirname "$0")/.." && pwd)"
repo_dir="$(cd "$api_dir/.." && pwd)"
export PATH="$HOME/.local/go/bin:$PATH"

db_name="ascard-contract-db-$$"
db_port="${CONTRACT_DB_PORT:-55441}"
app_port="${CONTRACT_APP_PORT:-55401}"
export DATABASE_URL="postgresql://contract:contract@127.0.0.1:${db_port}/contract"
export JWT_SECRET="contract-jwt-secret-0123456789abcdef0123456789"
export BASE_URL="http://127.0.0.1:${app_port}"

server_pid=""
log_file="$(mktemp)"
cleanup() {
  [ -n "$server_pid" ] && kill "$server_pid" 2>/dev/null || true
  [ -n "$server_pid" ] && wait "$server_pid" 2>/dev/null || true
  docker rm -f "$db_name" >/dev/null 2>&1 || true
  rm -f "$log_file"
}
trap cleanup EXIT

echo "==> Postgres ($db_name on :$db_port)"
docker run -d --rm --name "$db_name" -p "127.0.0.1:${db_port}:5432" \
  -e POSTGRES_USER=contract -e POSTGRES_PASSWORD=contract -e POSTGRES_DB=contract \
  postgres:16-alpine >/dev/null
for _ in $(seq 1 60); do
  docker exec "$db_name" pg_isready -U contract -d contract >/dev/null 2>&1 && break
  sleep 0.5
done
sleep 1
docker exec -i "$db_name" psql -q -v ON_ERROR_STOP=1 -U contract -d contract \
  < "$repo_dir/prisma/migrations/20261111000000_init/migration.sql"

if [ ! -f "$repo_dir/dist/index.html" ]; then
  echo "==> Building SPA (dist/)"
  (cd "$repo_dir" && npm run build >/dev/null)
fi

echo "==> Starting $target server on :$app_port"
case "$target" in
  go)
    (cd "$api_dir" && go build -o "$api_dir/bin/server" ./cmd/server)
    PORT="$app_port" WEB_DIR="$repo_dir/dist" "$api_dir/bin/server" >"$log_file" 2>&1 &
    ;;
  express)
    (cd "$repo_dir" && npx prisma generate >/dev/null 2>&1)
    (cd "$repo_dir" && PORT="$app_port" NODE_ENV=production node server/index.js) >"$log_file" 2>&1 &
    ;;
  *)
    echo "unknown target: $target (use go or express)" >&2
    exit 2
    ;;
esac
server_pid=$!

for _ in $(seq 1 60); do
  curl -fs "$BASE_URL/api/health" >/dev/null 2>&1 && break
  if ! kill -0 "$server_pid" 2>/dev/null; then
    echo "server exited early:" >&2
    cat "$log_file" >&2
    exit 1
  fi
  sleep 0.5
done

echo "==> Contract suite against $target"
status=0
(cd "$api_dir" && go test -tags=contract -count=1 ${CONTRACT_TEST_FLAGS:-} ./contract/) || status=$?
if [ "$status" -ne 0 ]; then
  echo "---- server log (last 40 lines) ----" >&2
  tail -40 "$log_file" >&2
fi
exit "$status"
