# PulsePoll Backend Status

## Current State

The backend MVP implementation is complete through **Step 9**.

Implemented scope includes:

- config, environment loading, and app bootstrap
- PostgreSQL and Redis connectivity
- `/health` endpoint
- database migrations for core MVP tables
- survey domain rules and computed flags
- authentication flow
- guest consent flow
- survey creation, detail, and feed endpoints
- voting and private PIN flow
- results and report endpoints
- worker placeholder skeleton

---

## Completed Steps

### Step 1 — App bootstrap and health

Implemented:

- backend config loading
- PostgreSQL connection
- Redis connection
- `/health` endpoint

Verified health response:

```json
{
  "db": "up",
  "ok": true,
  "redis": "up"
}
```

### Step 2 — Database schema

Implemented migrations for:

- users
- auth_sessions
- surveys
- survey_options
- reports
- optional feedback

### Step 3 — Survey domain rules

Implemented domain logic for:

- phase calculation:
- VOTING
- RESULTS
- EXPIRED
- computed flags:
- can_vote
- results_visible
- requires_pin

### Step 4 — Auth

Implemented:

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`
- `GET /api/v1/me`

Notes:

- refresh token uses HttpOnly (only Http access) cookie
- JWT signing method validation was patched
- login restricted to active users

### Step 5 — Guest consent

Implemented:

- `POST /api/v1/consent/accept`

Behavior:

- sets guest_id cookie
- prepares consent enforcement for guest voting flows

### Step 6 — Surveys

Implemented:

- `POST /api/v1/surveys`
- `GET /api/v1/surveys/:id`
- `GET /api/v1/feed`

Behavior includes:

- moderation keyword filter
- default survey timestamps
- enum validation for visibility and results_mode

### Step 7 — Voting and private PIN

Implemented:

- `POST /api/v1/surveys/:id/pin/verify`
- `POST /api/v1/surveys/:id/vote`
- `PUT /api/v1/surveys/:id/vote`

Behavior includes:

- consent enforcement for guests
- private PIN verification
- Redis-backed vote receipt storage
- Redis-backed PIN verification state
- vote change rules

Important:

- A critical consent/PIN bypass bug was found during review and patched before approval.
  Current implementation uses a safe sentinel error pattern so execution does not continue after an already-written HTTP error response.

### Step 8 — Results and reports

Implemented:

- `GET /api/v1/surveys/:id/results`
- `POST /api/v1/surveys/:id/report`

Behavior includes:

- results returned only when results_visible == true
- total vote counts and percentages
- zero-vote safety
- report persistence into reports table

### Step 9 — Worker skeleton

Implemented:

- `backend/cmd/worker/main.go`

This is a placeholder only:

- no business logic
- no DB logic
- no Redis logic
- no cleanup/stat aggregation yet

### Key Security / Correctness Fixes Applied

The following review findings were fixed during implementation:

- JWT signing method validation added
- login restricted to active users
- survey enum values validated before DB insert
- consent/PIN bypass bug fixed in vote flow
- private PIN is hashed before persistence using bcrypt

### Important MVP Notes

The following items are intentionally not implemented yet:

- brute-force protection for PIN verification
- IP-based rate limiting
- report spam protection
- advanced moderation pipelines
- worker cleanup/stat aggregation logic
- pagination improvements beyond MVP needs
- production hardening beyond current MVP scope

These are known future improvements, not accidental omissions.

### Review Status

All backend steps were reviewed step-by-step and approved after fixes where required.

Final backend state:

- scope-controlled
- MVP-aligned
- review-patched
- ready for frontend integration

#### Backend Entry Points

Main API app:

- `backend/cmd/api/main.go`

Worker placeholder:

- `backend/cmd/worker/main.go`

Main route packages:

- auth routes
- consent routes
- survey routes
- vote routes
- results/report routes

### Suggested Next Phase

Recommended next work after backend MVP completion:

- frontend integration
- API client / collection setup
- end-to-end flow testing
- README/API usage examples
- deployment preparation
