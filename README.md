![Coverage Badge](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/LbP22/7a0933f8cba0bddbcc95c8b850e32663/raw/spa-to-http_units_passing__heads_main.json) ![Coverage Badge](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/LbP22/7a0933f8cba0bddbcc95c8b850e32663/raw/spa-to-http_units_coverage__heads_main.json)

<a href="https://devforth.io"><img src="https://raw.githubusercontent.com/devforth/OnLogs/e97944fffc24fec0ce2347b205c9bda3be8de5c5/.assets/df_powered_by.svg" style="height:36px"/></a>

# spa-to-http

> A zero-configuration HTTP server for built SPA bundles.

`spa-to-http` serves your built frontend (`dist/`) with SPA fallback routing, cache-friendly defaults, and optional Brotli/Gzip compression.

If you want to ship a static SPA quickly in Docker without writing Nginx config, this is the fast path.

## Why use spa-to-http

- Fast-to-deploy: run one container, mount your build folder, done.
- Operationally simple: no custom web-server config files.
- Lightweight runtime: small Docker image and fast startup.
- SPA-focused defaults: sensible handling for `index.html`, hashed assets, and optional compression.

## Benchmark Highlights

### spa-to-http vs Nginx

| | spa-to-http | Nginx |
|---|---|---|
| Zero-configuration | ✅ No config files, SPA serving works out of the box | ❌ Requires dedicated config file |
| Config via env/CLI | ✅ Yes | ❌ No |
| Docker image size | ✅ 7.54 MiB (v1.1.1) | ❌ 142 MiB (v1.23.1) |
| Brotli out-of-the-box | ✅ Yes | ❌ Requires module |

Performance numbers and benchmark setup details are from <https://devforth.io/blog/deploy-react-vue-angular-in-docker-simply-and-efficiently-using-spa-to-http-and-traefik/>

| | spa-to-http | Nginx |
|---|---|---|
| Average time from container start to HTTP port availability (100 startups) | ✅ 1.358 s (11.5% faster) | ❌ 1.514 s |
| Requests-per-second on 0.5 KiB HTML file at localhost | ✅ 80497 (1.6% faster) | ❌ 79214 |
| Transfer speed on 0.5 KiB HTML file at localhost | ❌ 74.16 MiB/sec | ✅ 75.09 MiB/sec (1.3% faster) |
| Requests-per-second on 5 KiB JS file at localhost | ✅ 66126 (5.2% faster) | ❌ 62831 |
| Transfer speed on 5 KiB HTML file at localhost | ✅ 301.32 MiB/sec (4.5% faster) | ❌ 288.4 |

## Get Started in 60 Seconds

```bash
# Serve ./dist at http://localhost:8080
docker run --rm -p 8080:8080 -v $(pwd)/dist:/code devforth/spa-to-http:latest
```

Open `http://localhost:8080`.

## Recommended: Build SPA + Runtime in One Dockerfile

In most cases, users prefer building the SPA and serving it in one image:

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

## Common Copy-Paste Commands

### Custom port + Brotli

```bash
docker run --rm -p 8082:8082 \
  -v $(pwd)/dist:/code \
  devforth/spa-to-http:latest \
  --brotli --port 8082
```

### Basic Auth (Docker)

```bash
docker run --rm -p 8080:8080 \
  -e BASIC_AUTH="admin:secret" \
  -e BASIC_AUTH_REALM="SPA Server" \
  -v $(pwd)/dist:/code \
  devforth/spa-to-http:latest
```

### Subpath hosting (`/app`)

```bash
docker run --rm -p 8080:8080 \
  -v $(pwd)/dist:/code \
  devforth/spa-to-http:latest \
  --base-path /app
```

This maps `/app/...` requests to the same build root (for example `/app/assets/main.js` -> `/code/assets/main.js`).

### Compose / reverse proxy setup

For full Docker Compose and Traefik examples, see [`docs/deployment.md`](docs/deployment.md).

## Feature Snapshot

- Zero-configuration Docker usage for SPA bundles
- Optional Brotli/Gzip compression
- Cache-control tuning (`--cache-max-age`, `--ignore-cache-control-paths`)
- Subpath hosting with URL prefixes (`--base-path`)
- SPA mode toggle (`--spa` / `SPA_MODE`)
- In-memory file cache (`--cache`, `--cache-buffer`)
- Optional request logging and basic auth

## Documentation Map

| Need | Go to |
|---|---|
| Fast Docker onboarding | [`docs/getting-started.md`](docs/getting-started.md) |
| Local development and source-based workflows | [`docs/development.md`](docs/development.md) |
| Full flag/env reference and examples | [`docs/configuration.md`](docs/configuration.md) |
| Deployment behind Traefik / reverse proxy | [`docs/deployment.md`](docs/deployment.md) |
| Internal package layout and request flow | [`docs/architecture.md`](docs/architecture.md) |

## License

MIT
