[Back to README](../README.md) · [Configuration →](configuration.md)

# Getting Started

## Prerequisites

- Docker (recommended)
- A built SPA bundle (for example, `dist/` from Vite, Webpack, or similar)

## Local Run (Console)

Build from source and run the server directly:

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

### Configure via CLI Flags

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
  --threshold 2048
```

### Configure via Environment Variables

```bash
cd src
ADDRESS=0.0.0.0 PORT=8080 \
GZIP=true BROTLI=true \
SPA_MODE=true LOGGER=true LOG_PRETTY=true \
CACHE_MAX_AGE=3600 THRESHOLD=2048 \
DIRECTORY=../test/frontend/dist \
go run .
```

Full list of options is in [Configuration](configuration.md).

## Serve a Local Build (Docker)

```bash
# Serve ./dist at http://localhost:8080
docker run --rm -p 8080:8080 -v $(pwd)/dist:/code devforth/spa-to-http:latest
```

Open `http://localhost:8080` in your browser.

## Build + Run in One Dockerfile

Use this pattern when you want to build the SPA and ship a small runtime image:

```dockerfile
FROM node:20-alpine as builder
WORKDIR /code/
ADD package-lock.json .
ADD package.json .
RUN npm ci
ADD . .
RUN npm run build

FROM devforth/spa-to-http:latest
COPY --from=builder /code/dist/ .
```

Build and run:

```bash
docker build -q . | xargs docker run --rm -p 8080:8080
```

## Test With `test/frontend/dist`

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

## Next Steps

- Configure compression, cache settings, and ports in [Configuration](configuration.md)
- Deploy behind Traefik or another reverse proxy in [Deployment](deployment.md)

## See Also

- [Configuration](configuration.md) — Environment variables and CLI flags
- [Deployment](deployment.md) — Docker and reverse proxy setup
- [Architecture](architecture.md) — Project structure and request flow
