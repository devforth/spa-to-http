[← Architecture](architecture.md) · [Back to README](../README.md)

# Benchmarks

## spa-to-http vs Nginx

| | spa-to-http | Nginx |
|---|---|---|
| Zero-configuration | ✅ No config files, SPA serving works out of the box | ❌ Requires dedicated config file |
| Config via env/CLI | ✅ Yes | ❌ No |
| Docker image size | ✅ 13.2 MiB (v1.0.3) | ❌ 142 MiB (v1.23.1) |
| Brotli out-of-the-box | ✅ Yes | ❌ Requires module |

Performance numbers from the benchmark section of this post:
`https://devforth.io/blog/deploy-react-vue-angular-in-docker-simply-and-efficiently-using-spa-to-http-and-traefik/`

| | spa-to-http | Nginx |
|---|---|---|
| Average time from container start to HTTP port availability (100 startups) | ✅ 1.358 s (11.5% faster) | ❌ 1.514 s |
| Requests-per-second on 0.5 KiB HTML file at localhost | ✅ 80497 (1.6% faster) | ❌ 79214 |
| Transfer speed on 0.5 KiB HTML file at localhost | ❌ 74.16 MiB/sec | ✅ 75.09 MiB/sec (1.3% faster) |
| Requests-per-second on 5 KiB JS file at localhost | ✅ 66126 (5.2% faster) | ❌ 62831 |
| Transfer speed on 5 KiB HTML file at localhost | ✅ 301.32 MiB/sec (4.5% faster) | ❌ 288.4 |

## See Also

- [Getting Started](getting-started.md) — Install, build, and run
- [Deployment](deployment.md) — Docker and reverse proxy setup
- [Configuration](configuration.md) — Environment variables and CLI flags
