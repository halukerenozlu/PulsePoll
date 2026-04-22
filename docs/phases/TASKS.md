# TASKS — MVP Implementation Order

## Rules

- Follow `docs/SPEC.md`, `docs/API.md`, `docs/DB.md`, and `docs/REDIS.md` strictly.
- Implement in small, reviewable steps.
- Run relevant tests/build checks at each step.
- Do not add extra features unless they are approved and documented first.
- If behavior changes, add or update tests unless there is a clear reason not to.
- If no test is added for a behavior change, state the reason explicitly.

---

## Default Execution Flow

1. Human + ChatGPT define the exact scope
2. Codex implements the scoped step
3. Codex adds/updates tests when behavior changes
4. Codex runs relevant tests/build checks and reports results
5. Gemini reviews first
6. Codex applies necessary fixes and re-runs checks
7. Claude reviews only when the task is higher-risk
8. Human approves commit boundaries

---

## Step 0 — Repo Skeleton

- Ensure `backend/` and `docs/` structure exists
- Add `docker-compose.yml` (postgres + redis)
- Add `.env.example`

## Step 1 — Config + Connections

- config loader (`.env`)
- postgres connection (gorm)
- redis client
- `GET /health`

## Step 2 — Migrations

- users
- auth_sessions
- surveys
- survey_options
- reports
- feedback (optional)

## Step 3 — Domain: Phase & Computed Flags

- phase calculation (`VOTING` / `RESULTS` / `EXPIRED`)
- computed flags: `can_vote`, `results_visible`, `requires_pin`
- unit tests for boundary times

## Step 4 — Auth

- register / login / refresh / logout
- `/me`
- refresh token cookie (`HttpOnly`)

## Step 5 — Consent (Guest Voting)

- `POST /consent/accept` sets `guest_id` cookie
- voting endpoints require `guest_id` for guests
- return `403 CONSENT_REQUIRED` when missing

## Step 6 — Surveys

- `POST /surveys` (auth required, moderation filter)
- `GET /surveys/{id}`
- `GET /feed` (public, new-first MVP)

## Step 7 — PIN + Voting

- store `access_pin_hash` for `PRIVATE_PIN`
- `POST /surveys/{id}/pin/verify` -> `pinok` in Redis
- `POST /surveys/{id}/vote` -> enforce phase + pin + consent + max votes
- `PUT /surveys/{id}/vote` -> one-time change (only if enabled and `max_votes_per_user == 1`)

## Step 8 — Results + Report

- `GET /surveys/{id}/results` (only if `results_visible`)
- `POST /surveys/{id}/report`

## Step 9 — Worker Skeleton (Optional)

- `cmd/worker` placeholder for later cleanup / stat aggregation

---

## Verification Expectations

For each implemented step:

- relevant tests should pass
- relevant build/check commands should pass
- behavior should match the docs
- verification should be reported clearly
- frontend should not be required to prove backend correctness
