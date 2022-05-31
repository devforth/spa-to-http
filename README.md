# spa-to-http
Lightweight tiny zero-configuration SPA HTTP server. Serves SPA bundle on port 80 which makes it play well with traefik

# Benefits

* Zero-configuration, add it to your pipeline without managing additional configs
* Written in Go, which makes it fast and tiny
* Supports Brotly compression on original files, you don't need to archivate files by yourself, it does it for you
* Open-Source commertial friendly MIT license
* Plays well with all popular SPA frameworks and libraries: Vue, React, Angular and all bundlers: Webpack/Vite.
* Optimal statics caching out of the box: no-cache on index.html file to auto-update caches and infinite max-age for all other resources which have hash-URLs in all SPA frameworks.
* Created and maintained by Devforth 💪🏼



# Example use-cases

Create `Dockerfile` in yoru SPA directory (near `package.json`):

```
FROM node:16-alpine as builder
WORKDIR /code/
ADD package-lock.json .
RUN npm ci
ADD * ./
RUN npm run build

FROM spa-to-http:latest
COPY --from=builder /code/dist/ static/
```

## Available Options:
| Environment       | Command                                 | Description                                                                                                                                                                                                                            | Defaults |
|-------------------|-----------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------|
| ADDRESS           | `-a` or `--address`                     | Address to use                                                                                                                                                                                                                         | 0.0.0.0  |
| PORT              | `-p` or `--port`                        | Port to use                                                                                                                                                                                                                            | 8080     |
| GZIP              | `-g` or `--gzip`                        | When enabled it will create .gz files using gzip compression for files which size exceedes threshold and serve it instead of original one if client accepts gzip encoding. If brotli also enabled it will try to serve brotli first    | `false`  |
| BROTLI            | `-b` or `--brotli`                      | When enabled it will create .br files using brotli compression for files which size exceedes threshold and serve it instead of original one if client accepts brotli encoding. If gzip also enabled it will try to serve brotli first  | `false`  |
| THRESHOLD         | `--threshold <number>`                  | Threshold in bytes for gzip and brotli compressions                                                                                                                                                                                    | 1024     |
| DIRECTORY         | `-d <string>` or `--directory <string>` | Directory to serve                                                                                                                                                                                                                     | `.`      |
| DIRECTORY_LISTING | `--dir-lising`                          | Whether to show directory listing. SPA mode and directory listing cannot be enabled at the same time                                                                                                                                   | `false`  |
| CACHE_MAX_AGE     | `--cache-max-age <number>`              | Set cache time (in seconds) for cache-control max-age header                                                                                                                                                                           | 604800   |
| SPA_MODE          | `--spa`                                 | Where to enable SPA mode. In SPA mode if file for requested path does not exists server returns index.html from root of serving directory. SPA mode and directory listing cannot be enabled at the same time                           | `false`  |