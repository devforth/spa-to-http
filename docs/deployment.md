[← Configuration](configuration.md) · [Development](development.md) · [Back to README](../README.md) · [Architecture →](architecture.md)

# Deployment

`spa-to-http` is designed to run as a small container and sit behind a reverse proxy. It works well with Traefik and CDNs like Cloudflare.

## Traefik + Docker Compose Example

```yaml
version: "3.3"

services:
  traefik:
    image: "traefik:v2.7"
    command:
      - "--providers.docker=true"
      - "--providers.docker.exposedbydefault=false"
      - "--entrypoints.web.address=:80"
    ports:
      - "80:80"
    volumes:
      - "/var/run/docker.sock:/var/run/docker.sock:ro"

  spa:
    image: devforth/spa-to-http:latest
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.spa.rule=Host(`spa.localhost`)"
      - "traefik.http.services.spa.loadbalancer.server.port=8080"
```

## Notes

- Use `--port` if you run the container on a non-default port.
- Enable compression (`--brotli` or `--gzip`) when serving large static bundles.
- Use `--base-path` when the SPA is mounted under a subpath (for example, `/app`) behind your proxy.
- For fixed asset paths (for example, a service worker), use `--ignore-cache-control-paths` to avoid CDN caching issues.
- Add rate limiting at the reverse proxy (Traefik, Nginx, Cloudflare) to mitigate brute-force attempts.

## Subpath Deployment (`--base-path`)

When your reverse proxy exposes the app at a subpath such as `/app`, configure:

```yaml
services:
  spa:
    image: devforth/spa-to-http:latest
    command: --base-path /app
```

Behavior:
- `/app` and `/app/` serve `index.html`
- `/app/assets/...` maps to assets from the same dist root
- SPA routes under `/app/...` fall back to `index.html` (with `--spa=true`)

### `--ignore-cache-control-paths` with `--base-path`

Ignore paths are matched against both:
- Raw incoming path (for example, `/app/sw.js`)
- Internal mapped path (for example, `/sw.js`)

So either notation works in deployment config.

## See Also

- [Configuration](configuration.md) — Environment variables and CLI flags
- [Getting Started](getting-started.md) — Install, build, and run
- [Development](development.md) — Includes a local Traefik fixture for testing `/qwerty` base-path routing
- [Benchmarks](benchmarks.md) — spa-to-http vs Nginx
