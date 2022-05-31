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
