# PulsePoll

PulsePoll is a lightweight MVP survey platform built with a verification-first workflow.

The project is private and documentation-driven.
Core product and technical contracts live under `docs/`.

---

## Current Project Status

- Phase 1 — Completed
- Phase 2 — Completed
- Phase 3 — In progress

Current focus:
- backend feature work
- backend verification without frontend
- API contract readiness for future frontend integration

---

## Core Principle

Documentation is the source of truth.

If code, plans, or generated output conflict with project docs, the docs win.

Primary references:

- `docs/SPEC.md`
- `docs/API.md`
- `docs/DB.md`
- `docs/REDIS.md`
- `docs/verification.md`
- `docs/ROADMAP.md`
- `docs/phases/TASKS.md`
- `docs/phases/TASKS_PHASE2.md`
- `docs/phases/TASKS_PHASE3.md`

---

## Tech Stack

- Backend: Go + Fiber
- Database: PostgreSQL
- Ephemeral state / rate limiting: Redis
- Frontend: Next.js
- Infra: Docker Compose

Additional notes:
- PostgreSQL stores persistent survey-related data
- Redis stores temporary voting, PIN, and rate-limit state
- frontend integration belongs after backend verification is stable enough

---

## Local Development

Start services:

```bash
docker compose -p pulsepoll up --build
```

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

---

## Verification

Backend correctness must be verifiable without depending on the frontend.

Use:

- `docs/verification.md`

This document defines:
- local startup checks
- health verification
- endpoint verification expectations
- failure-path verification
- persistence verification guidance

---

## Workflow

Default workflow:

1. Human + ChatGPT define scope
2. Codex implements the approved task
3. Codex adds or updates tests when behavior changes
4. Codex runs relevant tests/build checks
5. Gemini performs first-pass review
6. Codex applies needed fixes and re-runs checks
7. Claude performs selective deep review only for higher-risk work
8. Human decides approval, commit, and tag boundaries

Claude is not required for every task.

---

## Docs Map

### Product and Contracts
- `docs/SPEC.md`
- `docs/API.md`
- `docs/DB.md`
- `docs/REDIS.md`

### Verification and Planning
- `docs/verification.md`
- `docs/ROADMAP.md`

### Phase Execution
- `docs/phases/TASKS.md`
- `docs/phases/TASKS_PHASE2.md`
- `docs/phases/TASKS_PHASE3.md`
- `docs/phases/phase3/sprints/SPRINT1.md`

---

## Frontend

Frontend work is planned, but backend verification comes first.

Phase 4 will connect the frontend to a backend that is already:
- implemented
- documented
- directly testable without UI
- stable enough for integration
