FROM golang:1.24-alpine3.21 AS builder

WORKDIR /code/

ADD src/ .

RUN go build -o dist/ -ldflags "-s -w"

FROM alpine:3.21

WORKDIR /code/

COPY docker-entrypoint.sh /bin/
ENTRYPOINT ["/bin/docker-entrypoint.sh"]

COPY --from=builder /code/dist/go-http-server /bin/
RUN chmod +x /bin/go-http-server

CMD ["/bin/go-http-server"]
