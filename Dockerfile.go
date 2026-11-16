# Go rewrite image: Go API + built SPA in one distroless container.
# Build: docker build -f Dockerfile.go -t ascard-app:go .
# (The Express image is still built from ./Dockerfile.)

# ---- 1. SPA (Vite) ----
FROM node:22-alpine AS web
WORKDIR /src
COPY package.json package-lock.json ./
RUN npm ci
COPY index.html vite.config.js ./
COPY public ./public
COPY src ./src
RUN npm run build

# ---- 2. Go binaries (server + seed) ----
FROM golang:1.27-alpine AS go
WORKDIR /src/api
COPY api/go.mod api/go.sum ./
RUN go mod download
COPY api/ ./
# ./cmd/... builds every command (server, seed) into /out, one binary per command.
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/ ./cmd/...

# ---- 3. runtime ----
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=go /out/ /app/
COPY --from=web /src/dist /app/dist
ENV WEB_DIR=/app/dist PORT=5001 LOG_FORMAT=json APP_ENV=production
EXPOSE 5001
USER nonroot:nonroot
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD ["/app/server", "-healthcheck"]
ENTRYPOINT ["/app/server"]
