# PulsePoll

PulsePoll is a lightweight MVP survey platform built with a verification-first workflow.

The project is private and documentation-driven.
Core product and technical contracts live under `docs/`.

---

## Current Project Status

- Current completed baseline: `v0.1.0` - Backend Foundation and Verification Baseline
- Active planning model: Version Milestones
- Next planning area: `v0.1.x` stabilization/docs cleanup or `v0.2.0` backend feature completion, depending on repo state

Current focus:

- keep backend behavior directly verifiable without frontend
- keep API, DB, Redis, and product docs aligned
- prepare remaining backend work for clear milestone execution

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
- `docs/VERSION_PLAN.md`
- `docs/ROADMAP.md`
- `CHANGELOG.md`

---

## Tech Stack

- Backend: Go + Fiber
- Database: PostgreSQL
- Ephemeral state / rate limiting: Redis
- Frontend: Next.js
- Infra: Docker Compose

Additional notes:

- PostgreSQL stores persistent survey-related data.
- Redis stores temporary voting, PIN, and rate-limit state.
- Frontend integration belongs after backend verification is stable enough.

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

1. Human + ChatGPT define the Version Milestone, Work Item, and Implementation Slice.
2. Codex implements the approved Implementation Slice.
3. Codex adds or updates tests when behavior changes.
4. Codex runs relevant tests/build checks.
5. Gemini performs first-pass review.
6. Codex applies needed fixes and re-runs checks.
7. Claude performs selective deep review only for higher-risk work.
8. Human decides approval, commit, and tag boundaries.

Claude is not required for every task.

---

## Docs Map

### Product and Contracts

- `docs/SPEC.md`
- `docs/API.md`
- `docs/DB.md`
- `docs/REDIS.md`

### Planning and History

- `docs/VERSION_PLAN.md`
- `docs/ROADMAP.md`
- `CHANGELOG.md`

### Verification and Testing

- `docs/verification.md`
- `docs/API_TESTING.md`

### Overview

- `docs/VISION.md`
- `docs/ARCHITECTURE.md`

---

## Frontend

Frontend work is planned after backend contract stability.

`v0.4.0` will connect the frontend to a backend that is already:

- implemented for the accepted backend milestone scope
- documented
- directly testable without UI
- stable enough for integration
