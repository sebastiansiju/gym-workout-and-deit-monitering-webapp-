# Fly.io Demo Instance — Setup Guide

One-time setup for the public demo at `demo.sebu.app` (or your chosen subdomain).

## Prerequisites

```bash
# Install flyctl
curl -L https://fly.io/install.sh | sh

# Log in
fly auth login
```

## 1. Create the app and volume

```bash
# From the project root
fly apps create sebu-demo            # or your preferred name
fly volumes create sebu_data --app sebu-demo --region fra --size 1
```

## 2. Set secrets

```bash
fly secrets set JWT_SECRET=$(openssl rand -hex 32) --app sebu-demo
```

## 3. Update fly.toml

Edit `fly.toml` — set `app`, `CORS_ORIGIN` and `[env].CORS_ORIGIN` to your final URL
(e.g. `https://sebu-demo.fly.dev` or `https://demo.sebu.app`).

## 4. Deploy

```bash
fly deploy --app sebu-demo
```

The first deploy will:
- Start the Go backend
- Seed the demo user (`demo@sebu.local` / `password123`)
- Begin seeding 800+ exercises async (takes ~30s)
- `DemoData` goroutine waits for exercises then seeds 8 weeks of workouts,
  90 days of weight logs, and 7 days of food logs automatically

## 5. Create the seed snapshot

Once exercises and demo data are fully seeded (~60s after first deploy):

```bash
fly ssh console --app sebu-demo -C "cp /app/data/sebu.db /app/data/sebu.seed.db"
```

From this point the hourly cron (`reset.sh`) will restore this snapshot every hour,
keeping the demo clean regardless of what visitors do.

## 6. Custom domain (optional)

```bash
fly certs add demo.sebu.app --app sebu-demo
```

Then add a CNAME `demo.sebu.app → sebu-demo.fly.dev` in your DNS.
Update `CORS_ORIGIN` in `fly.toml` and redeploy.

## Ongoing operations

```bash
# Check logs
fly logs --app sebu-demo

# Manual reset
fly ssh console --app sebu-demo -C "/app/reset.sh"

# Redeploy after code changes
fly deploy --app sebu-demo

# Update seed snapshot after improving demo data
fly ssh console --app sebu-demo -C "cp /app/data/sebu.db /app/data/sebu.seed.db"
```

## Architecture

Single Fly machine running:
- **nginx** (port 80) — serves React SPA, proxies `/api/` → localhost:3000
- **Go backend** (port 3000) — in a restart loop (so reset.sh can kill and restart it)
- **crond** — runs `reset.sh` at the top of every hour

SQLite DB lives on a persistent Fly volume at `/app/data/sebu.db`.
Seed snapshot lives alongside it at `/app/data/sebu.seed.db`.
