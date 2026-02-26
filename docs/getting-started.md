[Back to README](../README.md) · [Configuration →](configuration.md)

# Getting Started

## Prerequisites

- Docker (recommended)
- A built SPA bundle (for example, `dist/` from Vite, Webpack, or similar)

## Serve a Local Build

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

## Next Steps

- Configure compression, cache settings, and ports in [Configuration](configuration.md)
- Deploy behind Traefik or another reverse proxy in [Deployment](deployment.md)

## See Also

- [Configuration](configuration.md) — Environment variables and CLI flags
- [Deployment](deployment.md) — Docker and reverse proxy setup
- [Architecture](architecture.md) — Project structure and request flow
