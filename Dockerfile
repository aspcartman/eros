FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY eros /eros

ENTRYPOINT ["/eros"]