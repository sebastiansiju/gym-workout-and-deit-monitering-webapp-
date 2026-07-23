---
title: Backups & Updates
description: Back up your Sebu data (a single SQLite file) and update to the latest version safely.
---

All of your Sebu data — workouts, programs, weight, nutrition — lives in **one SQLite file**.
Backing up is copying that file. Updating is pulling the latest image. Your data survives both.

## Where your data lives

The database is stored on a Docker volume and mounted at `./data/sebu.db` next to your
`docker-compose.yml`. That single file *is* your instance.

## Back up

Copy the database file. For a personal instance with little write traffic, a live copy is fine:

```bash
cp ./data/sebu.db ./data/sebu.db.backup
```

For a guaranteed-consistent copy (e.g. right before an update), stop the stack first, then copy:

```bash
docker compose down
cp ./data/sebu.db ./data/sebu.db.backup
docker compose up -d
```

### Automate it (cron)

A nightly backup keeping the last 7 days:

```bash
0 3 * * * cp /path/to/sebu/data/sebu.db /path/to/backups/sebu-$(date +\%F).db && \
  find /path/to/backups -name 'sebu-*.db' -mtime +7 -delete
```

## Restore

Stop the stack, drop the backup in place, start again:

```bash
docker compose down
cp ./data/sebu.db.backup ./data/sebu.db
docker compose up -d
```

## Update to the latest version

```bash
docker compose pull
docker compose up -d
```

Your data volume is preserved across updates.

:::tip[Pin a stable version]
Tracking `main` gets you the newest features but also the newest rough edges. For a stable
self-host target, pin a released tag in your compose file instead of `latest`.
:::

:::note
Back up **before** every update. It's one file — you have no excuse.
:::
