ARG GO_VERSION=1.25.5
FROM golang:${GO_VERSION}-trixie AS builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /launched .

FROM debian:trixie-slim AS litestream

ARG LITESTREAM_VERSION=0.5.3
ADD https://github.com/benbjohnson/litestream/releases/download/v${LITESTREAM_VERSION}/litestream-${LITESTREAM_VERSION}-linux-x86_64.tar.gz /tmp/litestream.tar.gz
RUN tar -C /usr/local/bin -xzf /tmp/litestream.tar.gz && \
  rm /tmp/litestream.tar.gz

FROM debian:trixie-slim

# ca-certificates for HTTPS S3 replication
RUN apt-get update && \
  apt-get install -y ca-certificates sqlite3 && \
  rm -rf /var/lib/apt/lists/*

COPY --from=builder /launched /usr/local/bin/
COPY --from=litestream /usr/local/bin/litestream /usr/local/bin/

COPY litestream.yml /etc/litestream.yml

# Copy entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Create data directory for SQLite database
RUN mkdir -p /data

ENTRYPOINT ["/entrypoint.sh"]
