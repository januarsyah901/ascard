#!/usr/bin/env sh
# Redeploy ascard on this server: pull latest main, rebuild, wait until healthy.
# Usage: ./deploy.sh            (from the repo root)
#        ./deploy.sh --no-pull  (deploy the current working tree as-is)
set -eu

cd "$(dirname "$0")"

if [ ! -f .env.production ]; then
  echo "Missing .env.production (copy .env.production.example and fill it in)" >&2
  exit 1
fi

if [ "${1:-}" != "--no-pull" ]; then
  git pull --ff-only origin main
fi

docker compose --env-file .env.production up -d --build

printf 'Waiting for ascard-app to become healthy'
i=0
until [ "$(docker inspect -f '{{.State.Health.Status}}' ascard-app 2>/dev/null)" = "healthy" ]; do
  i=$((i + 1))
  if [ "$i" -gt 60 ]; then
    echo
    echo "ascard-app did not become healthy; recent logs:" >&2
    docker logs --tail 40 ascard-app >&2
    exit 1
  fi
  printf '.'
  sleep 2
done
echo ' ok'

docker image prune -f >/dev/null
echo "Deployed: https://ascard.hallojanu.me"
