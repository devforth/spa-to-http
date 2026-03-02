[← Deployment](deployment.md) · [Development](development.md) · [Back to README](../README.md) · [Benchmarks →](benchmarks.md)

# Architecture

`spa-to-http` is a single Go binary focused on static file serving for SPA bundles. It keeps the runtime simple, with a small number of packages and minimal configuration.

## Project Structure

```
./src
./src/main.go   # Entry point and wiring
./src/app       # HTTP server and handlers
./src/param     # CLI and environment parsing
./src/util      # Small shared helpers
```

## Request Flow

1. Parse configuration (CLI and environment variables).
2. Initialize the HTTP server and middleware.
3. Serve static files from the configured directory.
4. Apply caching headers optimized for SPA assets.
5. Optionally serve compressed files (Brotli/Gzip) when enabled.

## Design Goals

- Fast startup and low memory footprint
- Predictable caching for hashed assets and `index.html`
- Simple configuration for containerized use

## See Also

- [Configuration](configuration.md) — Environment variables and CLI flags
- [Getting Started](getting-started.md) — Install, build, and run
- [Benchmarks](benchmarks.md) — `spa-to-http` vs Nginx
