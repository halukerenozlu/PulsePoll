# PulsePoll Backend

Backend service for the PulsePoll MVP.

This document is backend-technical only.
Project milestone/status tracking lives outside this file and should be treated as authoritative in the project docs.

---

## Scope

This README exists to summarize:

- backend stack
- backend entry points
- local backend run notes
- health check expectations
- backend verification entry points
- implemented backend route groups / backend surface
- backend-specific safety notes
- backend-specific known limitations

It should **not** be used as a project roadmap, milestone tracker, or approval history log.

---

## Tech

- Go
- Fiber
- PostgreSQL
- Redis

---

## Backend Entry Points

From repo root:

- API app: `backend/cmd/api/main.go`
- Worker skeleton: `backend/cmd/worker/main.go`

---

## Local Run

From repo root:

```bash
docker compose -p pulsepoll up --build
```

Backend base URL (default local):

```text
http://localhost:8080
```

---

## Health Check

Use:

```text
GET http://localhost:8080/health
```

Expected healthy shape:

```json
{
  "db": "up",
  "ok": true,
  "redis": "up"
}
```

---

## Backend Verification Entry Points

Use these docs as the primary backend verification references:

- Verification guide: `docs/verification.md`
- API contract: `docs/API.md`
- Product rules: `docs/SPEC.md`
- DB contract: `docs/DB.md`
- Redis contract: `docs/REDIS.md`
- Active planning reference: `docs/VERSION_PLAN.md`

Useful backend test command:

```bash
cd backend
go test ./...
```

---

## Implemented Backend Surface

The backend codebase is organized around these major route / behavior areas:

- auth
- consent
- surveys
- voting
- results
- reporting

Expected major backend flow areas in the MVP:

- config loading and app bootstrap
- PostgreSQL and Redis connectivity
- `/health` endpoint
- database migrations for core MVP tables
- survey phase rules and computed flags
- auth/session flow
- guest consent flow
- survey create/detail/feed surface
- voting and private PIN flow
- results and report endpoints
- worker placeholder skeleton

This section is intentionally high-level.
For exact behavior and current contracts, use the docs under `docs/`.

---

## Backend-Specific Safety / Correctness Notes

These are the kinds of backend concerns that matter for this project and should remain visible during implementation and review:

- JWT signing method validation should be enforced where token parsing is used
- login/auth flows should respect user status rules
- survey enum/input validation should happen before persistence
- vote-related handlers must not continue execution after already writing an HTTP error response
- private PIN values should be hashed before persistence
- guest consent and PIN checks must be enforced consistently in voting flows

Do not treat this list as an implementation-complete checklist.
It is a backend correctness reminder.

---

## Known Limitations / Intentionally Not Implemented Yet

The following may remain intentionally out of scope until explicitly scheduled:

- brute-force protection for PIN verification
- IP-based rate limiting hardening beyond basic MVP needs
- report spam protection hardening
- advanced moderation pipelines
- worker cleanup/stat aggregation logic
- broader production hardening beyond current MVP scope

These are future improvements, not accidental omissions.

---

## Working Rule

If this README and the main docs ever conflict, the main docs win.

Use these as authoritative:

- `docs/SPEC.md`
- `docs/API.md`
- `docs/DB.md`
- `docs/REDIS.md`
- `docs/verification.md`
- `docs/ROADMAP.md`
- `docs/VERSION_PLAN.md`
