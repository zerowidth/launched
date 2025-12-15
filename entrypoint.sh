#!/bin/bash
set -e

if [ ! -f /data/launched.db ]; then
    echo "database not found locally, attempting restore from litestream..."
    if [ -n "$LITESTREAM_REPLICA_URL" ]; then
        litestream restore -if-db-not-exists /data/launched.db
    fi
fi

if [ -n "$LITESTREAM_REPLICA_URL" ]; then
    echo "Starting Litestream replication..."
    exec litestream replicate -exec "launched $*"
else
    echo "no litestream replica URL configured, running without replication..."
    exec launched "$@"
fi
