# Contributing to PulsePoll

Thanks for considering a contribution to PulsePoll.

PulsePoll is a small, documentation-driven MVP. Contributions are welcome when they keep the project focused, verifiable, and aligned with the contracts in `docs/`.

## Source of Truth

Product and technical decisions must follow the project docs.

Primary references:

- `docs/SPEC.md` for product rules and business behavior
- `docs/API.md` for endpoint contracts
- `docs/DB.md` for PostgreSQL contracts
- `docs/REDIS.md` for Redis key and TTL contracts
- `docs/verification.md` for verification expectations
- `docs/VERSION_PLAN.md` and `docs/ROADMAP.md` for planning scope
- `CHANGELOG.md` for completed historical changes

If code and docs disagree, the docs win. If a contribution changes behavior, update the relevant docs first.

## Project Scope

The MVP should stay small, usable, testable, and reliable.

Prefer:

- small reviewable changes
- explicit contracts
- simple implementations
- deterministic behavior
- reproducible verification

Avoid:

- speculative features
- hidden product rule changes
- broad refactors without a documented reason
- frontend behavior that invents backend fields or endpoints

## Planning Model

PulsePoll uses Version Milestones.

Before implementation work starts, the change should identify:

- Version Milestone
- Work Item
- Implementation Slice

This keeps contributions narrow and reviewable.

## Local Development

Start the local stack:

```bash
docker compose -p pulsepoll up --build
```

Backend health check:

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

## Verification

Every change needs appropriate verification before handoff.

For backend changes, run relevant Go tests from `backend/`:

```bash
go test ./...
```

For frontend changes, run the relevant checks from `frontend/`:

```bash
pnpm build
```

Use `docs/verification.md` and `docs/API_TESTING.md` for manual backend/API verification flows.

## Pull Requests

A good pull request should include:

- the milestone/work item/slice it belongs to
- a concise summary of what changed
- test or build commands run
- any docs updated
- any known limitations or follow-up work

Keep pull requests focused. Separate unrelated cleanup from behavior changes.

## Security Issues

Please do not open public issues for security vulnerabilities.
See `SECURITY.md` for the preferred reporting process.
