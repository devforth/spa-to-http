# spa-to-http

Lightweight zero-configuration SPA HTTP server. Serves SPA bundle on HTTP port which makes it play well with traefik

# Benefits

* Zero-configuration, add it to your pipeline without managing additional configs
* Written in Go, which makes it fast (no overhead on runtime) and tiny (small binary size)
* Supports Brotly compression on original files, you don't need to archivate files by yourself, it does it for you
* Open-Source commertial friendly MIT license
* Plays well with all popular SPA frameworks and libraries: Vue, React, Angular and bundlers: Webpack/Vite.
* Optimal statics caching out of the box: no-cache on index.html file to auto-update caches and infinite max-age for all other resources which have hash-URLs in all SPA frameworks.
* Created and maintained by Devforth 💪🏼



# Example use-cases

Create `Dockerfile` in yoru SPA directory (near `package.json`):

```
FROM node:16-alpine as builder
WORKDIR /code/
ADD package-lock.json .
ADD package.json .
RUN npm ci
ADD . .
RUN npm run build

FROM devforth/spa-to-http:latest
COPY --from=builder /code/dist/ .
```

# Example serving SPA with Traefik and Docker-Compose

```
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

  trfk-vue:
    build: "frontend"
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.trfk-vue.rule=Host(`trfk-vue.localhost`)"
      - "traefik.http.services.trfk-vue.loadbalancer.server.port=8080" # port inside of trfk-vue which should be used
```      


## Available Options:

| Environment Variable   | Command                                 | Description                                                                                                                                                                                                                           | Defaults |
|-------------------|-----------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------|
| ADDRESS           | `-a` or `--address`                     | Address to use                                                                                                                                                                                                                        | 0.0.0.0  |
| PORT              | `-p` or `--port`                        | Port to listen on                                                                                                                                                                                                                      | 8080     |
| GZIP              | `-g` or `--gzip`                        | When enabled it will create .gz files using gzip compression for files which size exceedes threshold and serve it instead of original one if client accepts gzip encoding. If brotli also enabled it will try to serve brotli first   | `false`  |
| BROTLI            | `-b` or `--brotli`                      | When enabled it will create .br files using brotli compression for files which size exceedes threshold and serve it instead of original one if client accepts brotli encoding. If gzip also enabled it will try to serve brotli first | `false`  |
| THRESHOLD         | `--threshold <number>`                  | Threshold in bytes for gzip and brotli compressions                                                                                                                                                                                   | 1024     |
| DIRECTORY         | `-d <string>` or `--directory <string>` | Directory to serve                                                                                                                                                                                                                    | `.`      |
| DIRECTORY_LISTING | `--dir-lising`                          | Whether to show directory listing. SPA mode and directory listing cannot be enabled at the same time                                                                                                                                  | `false`  |
| CACHE_MAX_AGE     | `--cache-max-age <number>`              | Set cache time (in seconds) for cache-control max-age header To disable cache set to -1. `.html` files are not being cached                                                                                                           | 604800   |
| SPA_MODE          | `--spa` or `--spa <bool>`               | When SPA mode if file for requested path does not exists server returns index.html from root of serving directory. SPA mode and directory listing cannot be enabled at the same time                                                  | `true`   |
