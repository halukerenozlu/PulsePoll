# VERSION_PLAN

## Active Planning Model

PulsePoll uses Version Milestones as the active planning model.

Legacy Phase/Sprint/Step files are no longer active planning references. Completed legacy work is summarized here and in `CHANGELOG.md`; detailed deleted files remain available in Git history.

Survey phase terminology remains product/domain terminology. Do not rename `VOTING`, `RESULTS`, or `EXPIRED` because those terms describe product behavior, not project planning.

---

## Versioning Rules

PulsePoll uses `vMAJOR.MINOR.PATCH`.

- `MAJOR`: reserved for future breaking product or architecture changes.
- `MINOR`: meaningful product, backend, frontend, or verification milestone.
- `PATCH`: small stabilization, docs sync, bug fix, contract clarification, or verification polish.

Tags are milestone checkpoints, not routine progress markers.

---

## Planning Vocabulary

Version Milestone:

- a version-level delivery target such as `v0.1.0`, `v0.2.0`, or `v0.3.0`
- has a clear purpose, scope, and completion signal

Work Item:

- a meaningful backend, frontend, docs, or product objective inside a Version Milestone
- should be narrow enough to review and verify

Implementation Slice:

- a small, Codex-executable unit of work under a Work Item
- should have explicit scope, out-of-scope boundaries, verification expectations, and handoff notes

No coding work should begin unless the Version Milestone, Work Item, and Implementation Slice are explicit.

---

## v0.1.0 - Backend Foundation and Verification Baseline

Status: Completed baseline

`v0.1.0` consolidates completed baseline work from the old planning model, including foundation work, backend stabilization, verification foundation, API testing foundation, and vote rate limiting.

Completion note:

- `v0.1.0` consolidates completed baseline work.
- It does not mean the entire old backend feature/readiness umbrella is complete.
- Only the old verification foundation and vote rate limiting increments from that umbrella are treated as completed here.
- Remaining backend feature and API readiness goals carry forward into future Version Milestones where appropriate.

### Product and Documentation Foundation

- Product rules captured in `docs/SPEC.md`.
- API contract captured in `docs/API.md`.
- Database contract captured in `docs/DB.md`.
- Redis contract captured in `docs/REDIS.md`.
- Project workflow and agent responsibilities documented.

### Local Development and Health Baseline

- Docker Compose local development path established.
- Backend, PostgreSQL, and Redis can run together locally.
- Health check reports backend dependency status.
- Expected healthy shape remains:

```json
{
  "db": "up",
  "ok": true,
  "redis": "up"
}
```

### Backend Foundation

- Go + Fiber backend foundation established.
- PostgreSQL connection baseline established.
- Redis connection baseline established.
- Health endpoint available for local dependency verification.

### Backend Stabilization

- Standardized API error handling.
- Request validation expectations established.
- Logging baseline improved for local development.
- Critical backend test coverage introduced where needed.
- Small cleanup/refactor work completed without broad redesign.

### Verification Foundation

- Backend verification must not depend on frontend work.
- Local startup, health, endpoint, error-path, and persistence verification guidance established in `docs/verification.md`.
- Endpoint verification template created for future backend work.

### API Testing Foundation

- Manual API testing guide established in `docs/API_TESTING.md`.
- Local assumptions documented:
  - Backend base URL: `http://localhost:8080`
  - API base path: `/api/v1`
  - Startup: `docker compose -p pulsepoll up --build`
- Auth, survey, vote, vote-change, results, report, and rate-limit testing examples documented.

### Vote Rate Limiting

- Redis-backed IP rate limiting added for vote endpoints:
  - `POST /surveys/:id/vote`
  - `PUT /surveys/:id/vote`
- Limit-exceeded requests return deterministic `429 TOO_MANY_REQUESTS`.
- Rate limiting uses Redis key `rl:ip:{ip}:vote`.
- TTL is 60 seconds.
- Tests and manual verification guidance cover allowed and blocked paths.

### Redis Contract Baseline

Redis contracts preserved for:

- `vote:survey:{surveyId}:user:{userId}`
- `vote:survey:{surveyId}:guest:{guestId}`
- `pinok:survey:{surveyId}:user:{userId}`
- `pinok:survey:{surveyId}:guest:{guestId}`
- `pinfail:survey:{surveyId}:guest:{guestId}`
- `rl:ip:{ip}:vote`

Vote rate limiting TTL remains 60 seconds.

---

## v0.1.x - Stabilization and Docs Cleanup

Purpose:

- small corrections after the Version Milestone migration
- docs alignment
- local run verification notes
- API testing guide polish
- small bug fixes or contract clarifications that do not change major milestone scope

### Work Item: Documentation Alignment

Implementation Slices:

- Review remaining docs for stale active planning terminology.
- Clarify cross-references between `README.md`, `docs/VERSION_PLAN.md`, `docs/ROADMAP.md`, and `CHANGELOG.md`.
- Keep historical references only where they explain completed work.

### Work Item: Local Run Verification Notes

Implementation Slices:

- Confirm Docker Compose startup notes match current repo behavior.
- Confirm `/health` examples match current response shape.
- Update verification notes only if current docs are ambiguous.

### Work Item: API Testing Polish

Implementation Slices:

- Improve manual API testing guide clarity without changing API behavior.
- Keep vote rate limiting verification aligned with `docs/REDIS.md`.
- Add missing verification notes only when backed by current implementation/docs.

---

## v0.2.0 - Backend Feature Completion

Purpose:

- complete remaining backend MVP endpoint/flow work
- close remaining backend behavior gaps
- preserve direct backend verification without frontend
- keep API, DB, Redis, and SPEC behavior aligned

### Work Item: Remaining MVP Endpoint Flows

Implementation Slices:

- Identify the next backend endpoint or flow from `docs/SPEC.md` and `docs/API.md`.
- Implement only one flow per slice.
- Add or update tests when route, validation, error, auth/session, persistence, or Redis behavior changes.
- Verify directly through HTTP without frontend.

### Work Item: Backend Behavior Gaps

Implementation Slices:

- Compare implemented behavior against `docs/SPEC.md`, `docs/API.md`, `docs/DB.md`, and `docs/REDIS.md`.
- Convert each confirmed gap into a small Implementation Slice.
- Update docs first if a behavior rule must change.

### Work Item: Verification Coverage for New Backend Work

Implementation Slices:

- Add endpoint-specific verification notes or API testing examples for each completed flow.
- Include success, failure, and persistence checks where relevant.
- Keep verification reproducible in local development.

---

## v0.3.0 - API Contract Readiness

Purpose:

- stabilize request/response/error behavior before frontend integration
- ensure `docs/API.md` and `docs/API_TESTING.md` are clear enough for frontend work
- ensure success, failure, and persistence behavior is directly verifiable
- remove backend ambiguity before frontend integration begins

### Work Item: API Contract Review

Implementation Slices:

- Review endpoint paths, methods, request fields, response fields, status codes, and error shapes.
- Clarify docs where behavior is already implemented but ambiguous.
- Avoid changing behavior unless explicitly approved as a scoped slice.

### Work Item: Failure-Path Readiness

Implementation Slices:

- Verify documented error classes against current backend behavior.
- Add focused tests for contract-sensitive failure paths when needed.
- Document manual checks for important failure scenarios.

### Work Item: Frontend Readiness Handoff

Implementation Slices:

- Ensure frontend-relevant fields are documented.
- Ensure API testing examples provide enough data for integration work.
- Record remaining known limitations before starting frontend integration.

---

## v0.4.0 - Frontend Integration

Purpose:

- connect the Next.js frontend to the verified backend
- implement loading, empty, success, and error states against real API behavior
- avoid frontend-invented backend fields or endpoints

### Work Item: API Client and Data Wiring

Implementation Slices:

- Connect one frontend flow to one verified backend endpoint group at a time.
- Use documented request/response fields only.
- Surface documented error behavior in the UI.

### Work Item: User-Facing States

Implementation Slices:

- Implement loading, empty, success, and error states for each connected flow.
- Keep state handling consistent with `docs/API.md`.
- Validate integration against local backend responses.

### Work Item: Integration Fixes

Implementation Slices:

- Fix only confirmed frontend/backend integration issues.
- Update docs first if a backend contract needs to change.
- Avoid adding frontend-only assumptions about backend behavior.

---

## v0.5.0 - End-to-End MVP Hardening

Purpose:

- full-stack validation
- bug fixing
- UX polish
- release readiness
- final MVP stabilization

### Work Item: Full-Stack Validation

Implementation Slices:

- Verify major user flows across backend and frontend.
- Confirm persistence and Redis-backed behavior through realistic flows.
- Record any release-blocking defects.

### Work Item: Bug Fixing and UX Polish

Implementation Slices:

- Fix confirmed MVP bugs in small, reviewable slices.
- Polish UX only where it improves clarity or reliability.
- Keep product rules aligned with `docs/SPEC.md`.

### Work Item: Release Readiness

Implementation Slices:

- Confirm docs match shipped behavior.
- Confirm local verification steps remain reproducible.
- Prepare final MVP handoff notes and tag checkpoint.
