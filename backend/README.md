# PulsePoll Backend

Backend service for PulsePoll MVP.

## Scope

This README is backend-technical only:

- local backend setup
- local backend run notes
- backend verification entry points

Project phase/status tracking lives outside this file.

## Tech

- Go
- Fiber
- PostgreSQL
- Redis

## Backend Entry Points

- API app: `backend/cmd/api/main.go`
- Worker skeleton: `backend/cmd/worker/main.go`

## Local Run

From repo root:

```bash
docker compose -p pulsepoll up --build
```

Backend base URL (default local):

- `http://localhost:8080`

Health check:

- `GET http://localhost:8080/health`

Expected healthy shape:

```json
{
  "db": "up",
  "ok": true,
  "redis": "up"
}
```

## Backend Verification Entry Points

- Verification guide: `docs/verification.md`
- API contract: `docs/API.md`
- Product rules: `docs/SPEC.md`
- DB contract: `docs/DB.md`
- Redis contract: `docs/REDIS.md`

Useful backend check:

```bash
cd backend
go test ./...
```
