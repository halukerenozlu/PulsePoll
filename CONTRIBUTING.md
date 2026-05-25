# Contributing to PulsePoll

Thanks for considering a contribution to PulsePoll.

PulsePoll is a small, documentation-driven MVP. Contributions are welcome when they keep the project focused, verifiable, and aligned with the contracts in `docs/`.

## Maintainer

**Haluk Eren Özlü** is the sole maintainer of PulsePoll and is responsible for all final decisions regarding the codebase, architecture, and user data handling.

All pull requests are reviewed and merged at the maintainer's discretion.

## Restricted Areas

The following areas are under direct maintainer control. Pull requests touching these areas will not be merged without explicit maintainer review and approval. If you are unsure whether your change falls into one of these categories, open an issue first and ask before writing code.

| Area                           | Examples                                                   |
| ------------------------------ | ---------------------------------------------------------- |
| Authentication & authorization | JWT logic, session handling, token refresh                 |
| User data handling             | Registration, login, email storage, password hashing       |
| Database schema                | Migrations, new tables, column changes on sensitive tables |
| Privacy-related endpoints      | Consent cookie logic, vote anonymization                   |
| Security configuration         | CORS, rate limiting rules, input validation                |
| Cookie & token logic           | Guest ID generation, auth cookie attributes                |

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
