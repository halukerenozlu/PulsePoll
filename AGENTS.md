# AGENTS — Project Roles & Working Rules

## Single Source of Truth

All product and technical decisions must be written into `docs/SPEC.md` first.

If any plan, code, suggestion, or generated output conflicts with:

- `docs/SPEC.md`
- `docs/API.md`
- `docs/DB.md`
- `docs/REDIS.md`

the docs win.

## Core Principle

Keep the MVP small, usable, and reliable.

Prefer:

- simple implementations
- explicit contracts
- small reviewable changes
- clear reasoning
- low operational complexity

Avoid:

- premature optimization
- speculative features
- hidden product rule changes
- stack redesign without an explicit decision in `docs/SPEC.md`

## Project Stack (MVP)

- Backend: Go + Fiber
- Database: PostgreSQL
- Ephemeral state / rate limiting: Redis
- Frontend: Next.js

## Primary Project Docs

Always use these documents as the main references:

- `docs/SPEC.md`
- `docs/API.md`
- `docs/DB.md`
- `docs/REDIS.md`
- `docs/TASKS.md`

Usage:

- `docs/SPEC.md` = product rules and business behavior
- `docs/API.md` = endpoint contracts
- `docs/DB.md` = PostgreSQL schema/contracts
- `docs/REDIS.md` = Redis key contracts and TTL rules
- `docs/TASKS.md` = implementation order

## Documentation Ownership

### Planner / Architecture Role

Responsible for maintaining and refining:

- `docs/SPEC.md`
- `docs/API.md`
- `docs/DB.md`
- `docs/REDIS.md`
- `docs/TASKS.md`

Rules:

1. Always align with `docs/SPEC.md`.
2. Record each new accepted product or technical rule in `docs/SPEC.md` first.
3. Keep scope MVP-small.
4. Prefer contracts, schemas, endpoint definitions, edge cases, and implementation order over large unsolicited code dumps.

### Implementation Role

Implement strictly from:

- `docs/SPEC.md`
- `docs/API.md`
- `docs/DB.md`
- `docs/REDIS.md`

Use `docs/TASKS.md` for implementation order.

Rules:

1. Work in small PR-style steps.
2. Do not invent new business rules without updating `docs/SPEC.md` first.
3. Include migrations when schema changes are required.
4. Add basic tests for core domain logic.
5. Keep folder structure clear and maintainable.

### Product Vision / Growth Role

Explore:

- UX flows
- sharing mechanics
- feed behavior
- moderation principles
- V2+ ideas

Rules:

1. Do not casually change the tech stack.
2. Convert every accepted product rule into `docs/SPEC.md`.
3. Prefer measurable hypotheses and simple experiments over broad feature expansion.

## MVP Product Direction

For exact product behavior, always follow `docs/SPEC.md`.

High-level summary:

- Registered users can create surveys.
- Guests cannot create surveys.
- Guests may vote only after accepting the required service cookie (`guest_id`).
- Guests may still browse surveys and view allowed results without accepting that cookie.
- PostgreSQL stores core persistent data and aggregated counts.
- Redis stores temporary voting-related state, PIN verification state, and rate-limit related ephemeral data.
- Raw vote event logs are out of MVP scope.

## Implementation Priorities

Build in this order unless `docs/TASKS.md` says otherwise:

1. Config + connections
2. Health endpoint
3. Migrations
4. Domain rules (phase calculation, computed flags)
5. Auth
6. Consent flow
7. Survey create/read/feed
8. PIN flow
9. Voting + one-time change
10. Results + reporting

## Quality Bar

Minimum acceptable quality:

- clear naming
- explicit error handling
- deterministic domain logic
- migrations committed
- no dead features outside MVP scope
- basic unit tests for phase logic and vote rules

## Change Policy

Before changing behavior:

1. update `docs/SPEC.md` if the rule changes
2. update related contract docs if needed (`API.md`, `DB.md`, `REDIS.md`)
3. implement the change
4. verify tests/build still pass

Do not let code become the source of truth before the docs.
