# TASKS — Claude Code (MVP)

## Rules
- Follow docs/SPEC.md, docs/API.md, docs/DB.md, docs/REDIS.md strictly.
- Implement in small steps; run tests/build at each step.
- Do not add extra features unless in SPEC.

## Step 0 — Repo skeleton
- Ensure backend/ and docs/ structure exists.
- Add docker-compose.yml (postgres + redis) + .env.example.

## Step 1 — Config + connections
- config loader (.env)
- postgres connection (gorm)
- redis client
- GET /health

## Step 2 — Migrations
- users, auth_sessions, surveys, survey_options, reports (+feedback optional)

## Step 3 — Domain: phase & computed flags
- phase calculation (VOTING/RESULTS/EXPIRED)
- computed flags: can_vote, results_visible, requires_pin
- unit tests for boundary times

## Step 4 — Auth
- register/login/refresh/logout + /me
- refresh token cookie (httpOnly)

## Step 5 — Consent (guest voting)
- POST /consent/accept sets guest_id cookie
- voting endpoints require guest_id for guests
- return 403 CONSENT_REQUIRED when missing

## Step 6 — Surveys
- POST /surveys (auth required, moderation filter)
- GET /surveys/{id}
- GET /feed (public new)

## Step 7 — PIN + Voting
- store access_pin_hash for PRIVATE_PIN
- POST /surveys/{id}/pin/verify -> pinok in Redis
- POST /surveys/{id}/vote -> enforce phase + pin + consent + max votes
- PUT /surveys/{id}/vote -> one-time change (only if enabled and max_votes_per_user==1)

## Step 8 — Results + report
- GET /surveys/{id}/results (only if results_visible)
- POST /surveys/{id}/report

## Step 9 — Worker skeleton (optional)
- cmd/worker placeholder for later cleanup/stat aggregation
