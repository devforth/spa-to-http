# spa-to-http

Lightweight zero-configuration SPA HTTP server. Serves SPA bundle on HTTP port which makes it play well with traefik

# Benefits

* Zero-configuration in Docker without managing additional configs
* 10x times smaller then Nginx, faster startup time, a little bit better performance
* Plays well with all popular SPA frameworks and libraries: Vue, React, Angular and bundlers: Webpack/Vite
* Supports Brotly compression on original files, you don't need to archivate files by yourself, it does it for you
* Written in Go, which makes it fast (no overhead on runtime) and tiny (small binary size)
* Open-Source commercial friendly MIT license
* Optimal statics caching out of the box: no-cache on index.html file to auto-update caches and infinite max-age for all other resources which have hash-URLs in all SPA frameworks.
* Plays well with CDNs caching (e.g. Clouflare/AWS CloudFront), support for ignoring cache of fixed URLs like service worker
* Created and maintained by Devforth 💪🏼

# Spa-to-http vs Nginx

| | Spa-to-http | Nginx |
|---|---|---|
| Zero-configuration | ✅No config files, SPA serving works out of the box with most optimal settings | ❌Need to create a dedicated config file |
| Ability to config settings like host, port, compression using Environment variables or CLI | ✅Yes | ❌No, only text config file |
| Docker image size | ✅13.2 MiB (v1.0.3) | ❌142 MiB (v1.23.1) |
| Brotli compression out-of-the-box | ✅Yes, just set env BROTLI=true | ❌You need a dedicated module like ngx_brotli |

Performence accroding to [Spa-to-http vs Nginx benchmark (End of the post)](https://devforth.io/blog/deploy-react-vue-angular-in-docker-simply-and-efficiently-using-spa-to-http-and-traefik/)

|  | Spa-to-http | Nginx |
|---|---|---|
| Average time from container start to HTTP port availability (100 startups) | ✅1.358 s (10.3% faster) | ❌1.514s |
| Requests-per-second on 0.5 KiB HTML file at localhost * | ✅80497 (1.6% faster) | ❌79214 |
| Transfer speed on 0.5 KiB HTML file * | ❌74.16 MiB/sec | ✅75.09 MiB/sec (1.7% faster) |
| Requests-per-second on 5 KiB JS file at localhost * | ✅66126 (5% faster) | ❌62831 |
| Transfer speed on 5 KiB HTML file * | ✅301.32 MiB/sec (4.3% faster) | ❌288.4 |

# Hello world & ussage

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

So we built our frontend and included it into container based on Spa-to-http. This way gives us great benefits:

* We build frontend in docker build time
* Bundle has only small resulting dist folder, there are no source code and node_modules so countainer is small
* When you start this container it serves SPA on HTTP port automatically with best settings. Because devforth/spa-to-http already has right CMD inside which runs SPA-to-HTTP webserver with right caching


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
    build: "spa" # name of the folder where Dockerfile is located
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.trfk-vue.rule=Host(`trfk-vue.localhost`)"
      - "traefik.http.services.trfk-vue.loadbalancer.server.port=8080" # port inside of trfk-vue which should be used
```      

How to enable Brotli compression:

```diff 
 trfk-vue:
    build: "spa"
++  command: --brotli
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.trfk-vue.rule=Host(`trfk-vue.localhost`)"
      - "traefik.http.services.trfk-vue.loadbalancer.server.port=8080"
```
How to change thresshold of small files which should not be compressed:

```diff 
 trfk-vue:
    build: "spa"
--  command: --brotli
++  command: --brotli --threshold 500
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.trfk-vue.rule=Host(`trfk-vue.localhost`)"
      - "traefik.http.services.trfk-vue.loadbalancer.server.port=8080"
```

How to run container on a custom port:


```diff 
 trfk-vue:
    build: "spa"
++  command: --brotli --port 8082
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.trfk-vue.rule=Host(`trfk-vue.localhost`)"
--    - "traefik.http.services.trfk-vue.loadbalancer.server.port=8080"
++    - "traefik.http.services.trfk-vue.loadbalancer.server.port=8082"
```

Ignore caching for some specific resources, e.g. prevent Service Worker caching on CDNs like Cloudflare:



```diff 
 trfk-vue:
    build: "spa"
++  command: --ignore-cache-control-paths "/sw.js"
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.trfk-vue.rule=Host(`trfk-vue.localhost`)"
      - "traefik.http.services.trfk-vue.loadbalancer.server.port=8080"
```

This is not needed for most of your assets because their filenames should contain file hash (added by default by modern bundlers). So cache naturally invalidated by referencing hashed assets from uncachable html. However some special resources like service worker must be served on fixed URL without file hash in filename



## Available Options:

| Environment Variable       | Command                                 | Description                                                                                                                                                                                                                           | Defaults |
|----------------------------|-----------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------|
| ADDRESS                    | `-a` or `--address`                     | Address to use                                                                                                                                                                                                                        | 0.0.0.0  |
| PORT                       | `-p` or `--port`                        | Port to listen on                                                                                                                                                                                                                     | 8080     |
| GZIP                       | `--gzip`                                | When enabled it will create .gz files using gzip compression for files which size exceedes threshold and serve it instead of original one if client accepts gzip encoding. If brotli also enabled it will try to serve brotli first   | `false`  |
| BROTLI                     | `--brotli`                              | When enabled it will create .br files using brotli compression for files which size exceedes threshold and serve it instead of original one if client accepts brotli encoding. If gzip also enabled it will try to serve brotli first | `false`  |
| THRESHOLD                  | `--threshold <number>`                  | Threshold in bytes for gzip and brotli compressions                                                                                                                                                                                   | 1024     |
| DIRECTORY                  | `-d <string>` or `--directory <string>` | Directory to serve                                                                                                                                                                                                                    | `.`      |
| CACHE_CONTROL_MAX_AGE      | `--cache-control-max-age <number>`      | Set cache time (in seconds) for cache-control max-age header To disable cache set to -1. `.html` files are not being cached                                                                                                           | 604800   |
| IGNORE_CACHE_CONTROL_PATHS | `--ignore-cache-control-paths <string>` | Additional paths to set "Cache-control: no-store" via comma, example "/file1.js,/file2.js"                                                                                                                                            |          |
| SPA_MODE                   | `--spa` or `--spa <bool>`               | When SPA mode if file for requested path does not exists server returns index.html from root of serving directory. SPA mode and directory listing cannot be enabled at the same time                                                  | `true`   |
| CACHE                      | `--cache`                               | When enabled f.Open reads are being cached using Two Queue LRU Cache in bits                                                                                                                                                          | `true`   |
| CACHE_BUFFER               | `--cache-buffer <number>`               | Specifies the maximum size of LRU cache in bytes                                                                                                                                                                                      | `51200`  |
