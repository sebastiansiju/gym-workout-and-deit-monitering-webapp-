---
title: FAQ
description: Common questions about self-hosting Sebu, the open-source Hevy / Strong alternative.
---

## Is Sebu really free?

Yes. Sebu is open source under the **MIT license**. No subscription, no paywalled features, no
"your export is a Pro feature." You run it yourself.

## How is it different from Hevy, Strong, Wger, or FitNotes?

- **Hevy / Strong** — polished, but cloud-only, increasingly paywalled, and your data lives on
  someone else's server.
- **Wger** — a solid self-hosted option; Sebu's focus is a more modern, mobile-first UI and a
  simpler one-command deploy.
- **FitNotes** — local-only, with no sync or server deployment story.

Sebu is for people who want a modern, mobile-first workout tracker they **fully own** and can run
on a $5 VPS or a Raspberry Pi.

## Where is my data stored?

In a single SQLite file on your server. Nothing is sent to any third party. Backing up is copying
that one file.

## Do I need an API key or account for the exercise library?

No. Sebu seeds 800+ exercises from [free-exercise-db](https://github.com/yuhonas/free-exercise-db)
automatically on first startup. No key, no setup.

## Which platforms are supported?

Android (side-load APK) and web today; iOS is planned via TestFlight / the App Store. The server
runs on Raspberry Pi, x86 VPS, Synology NAS, Proxmox, and local Mac/Linux/Windows.

## Can I import my data from Strong or Hevy?

CSV import is on the roadmap (planned), alongside a PWA build.

## How do I get help?

Open an issue on [GitHub](https://github.com/Cawlumm/sebu) or join the
[Discord](https://discord.gg/hfFWsrebQA).
