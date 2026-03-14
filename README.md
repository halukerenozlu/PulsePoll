# PulsePoll

This repo is a lightweight skeleton. Core contracts live in `docs/`.

## Start services

```bash
docker compose up --build
```

API health:

- GET http://localhost:8080/health

## Docs (Single Source of Truth)

- docs/SPEC.md
- docs/API.md
- docs/DB.md
- docs/REDIS.md
- docs/TASKS_CLAUDE.md

# Tech Stack

- Backend: Go + Fiber
- DB: Postgres
- Cache/Rate limit/Counter: Redis
- DB access: GORM + (hot-path) sqlc/raw SQL
- Migrations: goose/atlas
- Frontend: Next.js (App Router)
- Client data: React Query / SWR + services/api.ts
- Infra: Docker Compose

# Frontend (Next.js) — placeholder

Plan: use Next.js (App Router) for SEO + share previews + public feed.
For now this folder is a placeholder so we don't create lots of files before decisions are finalized.

When ready:

- `npx create-next-app@latest`
- keep API calls in a single `src/services/api.ts`
