[← Getting Started](getting-started.md) · [Back to README](../README.md) · [Configuration →](configuration.md)

# Development

This page is for contributors and local source-based workflows.

## Prerequisites

- Go 1.24+
- Git

## Run from Source

Build and run the server directly:

```bash
cd src
go build -o spa-to-http
./spa-to-http --directory ../test/frontend/dist
```

Or run without building a binary:

```bash
cd src
go run . --directory ../test/frontend/dist
```

Open `http://localhost:8080` in your browser.

## Configure via CLI Flags

```bash
cd src
go run . \
  --directory ../test/frontend/dist \
  --brotli \
  --gzip \
  --spa=true \
  --logger \
  --log-pretty \
  --cache-max-age 3600 \
  --threshold 2048 \
  --base-path /app
```

## Configure via Environment Variables

```bash
cd src
ADDRESS=0.0.0.0 PORT=8080 \
GZIP=true BROTLI=true \
SPA_MODE=true LOGGER=true LOG_PRETTY=true \
CACHE_MAX_AGE=3600 THRESHOLD=2048 \
DIRECTORY=../test/frontend/dist \
go run .
```

## Basic Auth (Console)

```bash
cd src
go run . \
  --directory ../test/frontend/dist \
  --basic-auth "admin:secret" \
  --basic-auth-realm "SPA Server"
```

## Test with `test/frontend/dist`

Use the built-in fixture to verify routing, caching, and compression behavior.

### Console

```bash
cd src
go run . --directory ../test/frontend/dist
```

### Docker

```bash
docker run --rm -p 8080:8080 -v $(pwd)/test/frontend/dist:/code devforth/spa-to-http:latest
```

Open `http://localhost:8080` in your browser.

## Prod-like Local Test: Traefik + Base Path `/qwerty`

Use the local fixture at `test/traefik/docker-compose.base-path.yml` to emulate production routing:
- `spa-to-http` serves only `${SPA_BASE_PATH}` paths (default: `/qwerty`)
- non-base-path requests are routed to a different backend (`whoami`)
- `spa-to-http` image is built locally from this repository

Run from repository root:

```bash
SPA_BASE_PATH=/qwerty SPA_MODE=true \
docker compose -f test/traefik/docker-compose.base-path.yml up --build
```

Verify behavior:

```bash
# SPA (should return your app index)
curl -i http://localhost:8081/qwerty

# SPA asset under base path (should be served by spa-to-http)
curl -i http://localhost:8081/qwerty/vite.svg

# Not SPA (should be routed to whoami, not your app)
curl -i http://localhost:8081/
curl -i http://localhost:8081/anything
```

Stop:

```bash
docker compose -f test/traefik/docker-compose.base-path.yml down
```

## See Also

- [Getting Started](getting-started.md) — Fast Docker onboarding
- [Configuration](configuration.md) — Full options reference
- [Architecture](architecture.md) — Package structure and request flow
