[Back to README](../README.md) · [Development →](development.md)

# Getting Started

## Prerequisites

- Docker
- A built SPA bundle (for example, `dist/` from Vite, Webpack, or similar)

## Run a Local Build (Docker)

```bash
# Serve ./dist at http://localhost:8080
docker run --rm -p 8080:8080 -v $(pwd)/dist:/code devforth/spa-to-http:latest
```

Open `http://localhost:8080` in your browser.

## Basic Auth (Docker)

```bash
docker run --rm -p 8080:8080 \
  -e BASIC_AUTH="admin:secret" \
  -e BASIC_AUTH_REALM="SPA Server" \
  -v $(pwd)/dist:/code \
  devforth/spa-to-http:latest
```

## Build + Run in One Dockerfile

Use this pattern when you want to build the SPA and ship a small runtime image:

```dockerfile
FROM node:20-alpine AS builder
WORKDIR /code
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM devforth/spa-to-http:latest
COPY --from=builder /code/dist/ .
```

Build and run:

```bash
docker build -t my-spa .
docker run --rm -p 8080:8080 my-spa
```

## Next Steps

- Configure compression, cache settings, and ports in [Configuration](configuration.md)
- Deploy behind Traefik or another reverse proxy in [Deployment](deployment.md)
- For source-based local runs and contributor workflow, see [Development](development.md)

## See Also

- [Development](development.md) — Local setup, source runs, and fixture usage
- [Configuration](configuration.md) — Environment variables and CLI flags
- [Deployment](deployment.md) — Docker and reverse proxy setup
- [Architecture](architecture.md) — Project structure and request flow
