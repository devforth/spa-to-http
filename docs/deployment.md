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
- For fixed asset paths (for example, a service worker), use `--ignore-cache-control-paths` to avoid CDN caching issues.
- Add rate limiting at the reverse proxy (Traefik, Nginx, Cloudflare) to mitigate brute-force attempts.

## See Also

- [Configuration](configuration.md) — Environment variables and CLI flags
- [Getting Started](getting-started.md) — Install, build, and run
- [Benchmarks](benchmarks.md) — spa-to-http vs Nginx
