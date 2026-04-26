# Changelog

All notable project changes are documented here in a concise Keep a Changelog style.

## Unreleased

### Changed

- Migrated active planning documentation from the legacy Phase/Sprint/Step model to the Version Milestone model.
- Added `docs/VERSION_PLAN.md` as the active planning source of truth.
- Replaced active roadmap and role guidance with Version Milestone, Work Item, and Implementation Slice terminology.
- Removed legacy active planning files from `docs/phases/`; historical context is preserved in Git history and summarized here.
- Noted that feed endpoint documentation is simplified/aligned with the current MVP contract by keeping only currently documented MVP feed query behavior.
- Aligned backend README references with the Version Milestone planning model.

## v0.1.0 - Backend Foundation and Verification Baseline

### Added

- Local Docker-based development baseline for backend, PostgreSQL, and Redis.
- Health check with dependency status for DB and Redis.
- Initial backend foundation for the Go + Fiber API.
- PostgreSQL MVP schema contract documentation.
- Redis key and TTL contract documentation for vote receipts, PIN state, PIN failure tracking, and vote rate limiting.
- Verification-first backend workflow documentation.
- Manual API testing guide for local backend checks.
- Redis-backed vote endpoint rate limiting using `rl:ip:{ip}:vote` with a 60-second TTL.

### Changed

- Backend stabilization completed for consistent errors, validation, logging, critical test coverage, and small cleanup.
- API testing documentation expanded with vote and vote-change verification guidance.

### Notes

- `v0.1.0` consolidates completed baseline work from the old planning model: foundation work, backend stabilization, verification foundation, API testing foundation, and vote rate limiting.
- This does not mean the entire former backend feature/readiness umbrella is complete. Remaining backend and API readiness work continues in future Version Milestones.
