![Coverage Badge](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/LbP22/7a0933f8cba0bddbcc95c8b850e32663/raw/spa-to-http_units_passing__heads_main.json) ![Coverage Badge](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/LbP22/7a0933f8cba0bddbcc95c8b850e32663/raw/spa-to-http_units_coverage__heads_main.json)

<a href="https://devforth.io"><img src="https://raw.githubusercontent.com/devforth/OnLogs/e97944fffc24fec0ce2347b205c9bda3be8de5c5/.assets/df_powered_by.svg" style="height:36px"/></a>

# spa-to-http

> World's fastest lightweight zero-configuration SPA HTTP server.

Serve a built SPA bundle over HTTP with sensible defaults for caching and optional Brotli/Gzip compression. It’s designed to run cleanly in Docker and behind reverse proxies like Traefik or Cloudflare.

## Quick Start

```bash
# Serve ./dist at http://localhost:8080
docker run --rm -p 8080:8080 -v $(pwd)/dist:/code devforth/spa-to-http:latest
```

## Key Features

- Zero-configuration Docker setup for common SPA outputs
- Small image and fast startup (Go binary)
- Optional Brotli/Gzip compression
- Cache-control optimized for hashed assets and index.html
- Works with popular SPA toolchains (React, Vue, Angular, Svelte, Vite, Webpack)

## Example

```bash
# Enable Brotli and serve on a custom port
docker run --rm -p 8082:8082 -v $(pwd)/dist:/code devforth/spa-to-http:latest --brotli --port 8082
```

---

## Documentation

| Guide | Description |
|-------|-------------|
| [Getting Started](docs/getting-started.md) | Install, build, and run | 
| [Configuration](docs/configuration.md) | Environment variables and CLI flags |
| [Deployment](docs/deployment.md) | Docker and reverse proxy setup |
| [Architecture](docs/architecture.md) | Project structure and request flow |
| [Benchmarks](docs/benchmarks.md) | spa-to-http vs Nginx |

## License

MIT
