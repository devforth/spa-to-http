[← Getting Started](getting-started.md) · [Development](development.md) · [Back to README](../README.md) · [Deployment →](deployment.md)

# Configuration

`spa-to-http` supports both environment variables and CLI flags. Environment variables map directly to the CLI options.

## Options

| Environment Variable | Command | Description | Default |
|---|---|---|---|
| ADDRESS | `-a` or `--address` | Address to use | `0.0.0.0` |
| PORT | `-p` or `--port` | Port to listen on | `8080` |
| GZIP | `--gzip` | Enable gzip compression for files above the threshold | `false` |
| BROTLI | `--brotli` | Enable Brotli compression for files above the threshold | `false` |
| THRESHOLD | `--threshold <number>` | Threshold in bytes for gzip and Brotli | `1024` |
| DIRECTORY | `-d <string>` or `--directory <string>` | Directory to serve | `.` |
| BASE_PATH | `--base-path <string>` | URL prefix to mount the SPA under (normalized path prefix) | `/` |
| CACHE_MAX_AGE | `--cache-max-age <number>` | Cache max-age in seconds; use `-1` to disable | `604800` |
| IGNORE_CACHE_CONTROL_PATHS | `--ignore-cache-control-paths <string>` | Comma-separated paths to force `Cache-Control: no-store` | (empty) |
| SPA_MODE | `--spa` or `--spa <bool>` | Serve `index.html` on missing paths (SPA routing) | `true` |
| CACHE | `--cache` | Enable in-memory file read cache (LRU) | `true` |
| CACHE_BUFFER | `--cache-buffer <number>` | Max size of the LRU cache in bytes | `51200` |
| LOGGER | `--logger` | Enable request logging | `false` |
| LOG_PRETTY | `--log-pretty` | Pretty-print logs instead of JSON | `false` |
| BASIC_AUTH | `--basic-auth <username:password>` | Enable Basic Auth (username:password) | (empty) |
| BASIC_AUTH_REALM | `--basic-auth-realm <string>` | Basic Auth realm name | `Restricted` |

## Security Note

When Basic Auth is enabled, always run behind HTTPS (or a TLS-terminating reverse proxy). Basic Auth over plain HTTP exposes credentials.

## Examples

Enable Brotli:

```yaml
# docker-compose.yml
services:
  spa:
    image: devforth/spa-to-http:latest
    command: --brotli
```

Change compression threshold:

```yaml
services:
  spa:
    image: devforth/spa-to-http:latest
    command: --brotli --threshold 500
```

Serve on a custom port:

```yaml
services:
  spa:
    image: devforth/spa-to-http:latest
    command: --port 8082
    ports:
      - "8082:8082"
```

Ignore cache-control for fixed paths (for example, a service worker):

```yaml
services:
  spa:
    image: devforth/spa-to-http:latest
    command: --ignore-cache-control-paths "/sw.js"
```

## Base Path (Subpath Hosting)

Use `--base-path` when your SPA is exposed behind a URL prefix (for example, `/app`) while files stay in the same build directory.

### Normalization Rules

- Empty value becomes `/`
- Missing leading slash is added (`app` -> `/app`)
- Trailing slash is removed (`/app/` -> `/app`)
- Query strings and fragments are invalid (`/app?a=1`, `/app#x`)

### Request Mapping

With `--base-path /app` and `--directory /code`:

- `/app` -> `/code/index.html`
- `/app/` -> `/code/index.html`
- `/app/assets/main.js` -> `/code/assets/main.js`
- `/app/route1` -> SPA fallback to `/code/index.html` (when `--spa` is enabled)

Requests outside the base path continue to follow normal root behavior:

- `/assets/main.js` -> `/code/assets/main.js`

### Interaction with `IGNORE_CACHE_CONTROL_PATHS`

When `--base-path` is set, ignore paths are matched against both:

1. Raw incoming request path (for example, `/app/sw.js`)
2. Mapped internal path (for example, `/sw.js`)

That means either style works:

```yaml
services:
  spa:
    image: devforth/spa-to-http:latest
    command: --base-path /app --ignore-cache-control-paths "/sw.js,/app/sw.js"
```

## See Also

- [Getting Started](getting-started.md) — Install, build, and run
- [Deployment](deployment.md) — Docker and reverse proxy setup
- [Architecture](architecture.md) — Project structure and request flow
