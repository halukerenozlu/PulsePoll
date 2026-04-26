# ARCHITECTURE

## Current Architecture

PulsePoll is an MVP built around a small backend-first architecture.

## Backend

- Go + Fiber application
- Exposes health and API routes
- Uses `/api/v1` as the API base path
- Implements product behavior according to `docs/SPEC.md`

## PostgreSQL

PostgreSQL stores persistent data:

- users
- auth sessions
- surveys
- survey options
- reports
- optional feedback

The MVP stores aggregated vote counts and does not store a raw vote event log.

## Redis

Redis stores ephemeral state:

- vote receipts and vote-change tracking
- guest/user PIN verification state
- PIN brute-force counters
- vote endpoint rate limiting

Redis key contracts and TTL rules are defined in `docs/REDIS.md`.

## Frontend

- Next.js frontend exists as the planned integration surface.
- Frontend integration belongs after backend contracts are stable and directly verifiable.
- Frontend work must not invent backend fields, endpoints, or product behavior.

## Local Development

Docker Compose is the local development path for backend dependencies and service startup.

Use:

```bash
docker compose -p pulsepoll up --build
```

Then verify:

```text
GET http://localhost:8080/health
```
