#!/bin/sh
# reset.sh — hourly demo reset, invoked by crond inside the container
#
# Restores the live DB from a pre-seeded snapshot so the demo always shows
# realistic data. After the copy, pkill stops the backend; entrypoint.sh's
# restart loop brings it back with the fresh DB.
#
# The seed snapshot (sebu.seed.db) is created once on first deploy:
#   fly ssh console -a sebu-demo
#   cp /app/data/sebu.db /app/data/sebu.seed.db
SEED="/app/data/sebu.seed.db"
LIVE="/app/data/sebu.db"

if [ ! -f "$SEED" ]; then
    echo "[reset] $(date): no seed snapshot at $SEED — skipping"
    exit 0
fi

echo "[reset] $(date): restoring demo DB..."
cp "$SEED" "$LIVE"
# Kill backend so entrypoint restart loop picks up the fresh DB immediately
pkill sebu-api 2>/dev/null || true
echo "[reset] $(date): done — backend will restart automatically"
