# Project Agent Guide (Single Source of Truth)

## Golden Rule
All product/tech decisions must be written into `docs/SPEC.md` first.
If any plan/code conflicts with SPEC, SPEC wins.

## MVP Focus
Ship a usable product fast. Prefer simple, reliable solutions over premature optimization.

## Stack (MVP)
- Backend: Go + Fiber
- DB: PostgreSQL
- Ephemeral state/rate-limit: Redis
- Frontend: Next.js (for SEO + share previews)

## Default Survey Lifecycle
- Voting: 24h
- Results-only: 24h
- Expire (retention end): 48h (delete ephemeral state; keep only global stats later)

## Voting & Access (MVP)
- Registered users can create surveys.
- Guests cannot create surveys.
- Guests may vote **only if** they accept the “service-required cookie” (guest_id).
- If cookie not accepted: guests can browse feed, view surveys and results, but cannot vote or change vote.
- Closed surveys may require a PIN to vote; results are visible to link-holders when results are allowed by phase.

## Data Policy (MVP)
- Postgres stores surveys/options + aggregated counts only.
- Redis stores vote receipts (for limits + one-time change) and PIN verification state with TTL.
- No raw vote event log in MVP.
