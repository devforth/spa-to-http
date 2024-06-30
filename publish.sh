docker buildx create --use
docker buildx build --platform=linux/amd64,linux/arm64 --tag "devforth/spa-to-http:latest" --tag "devforth/spa-to-http:1.0.8" --push .
