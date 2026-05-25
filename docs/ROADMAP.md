# ROADMAP

This roadmap is a high-level Version Milestone sequence. It is not a task tracker.

Use `docs/VERSION_PLAN.md` for active execution details, Work Items, and Implementation Slices.

---

## v0.1.0 - Backend Foundation and Verification Baseline

Status: Completed baseline

Purpose:

- consolidate completed foundation work
- establish local Docker-based development
- establish health checks with DB and Redis dependency status
- establish backend stabilization
- establish verification-first backend workflow
- establish API testing guidance
- add Redis-backed vote endpoint rate limiting

Note:

- This milestone does not mean all former backend feature/readiness goals are complete.
- Remaining backend and API readiness work continues in later Version Milestones.

---

## v0.1.x - Stabilization and Docs Cleanup

Status: Active planning area as needed

Purpose:

- small corrections after the Version Milestone migration
- docs alignment
- local run verification notes
- API testing guide polish
- small bug fixes or contract clarifications that do not change major milestone scope

---

## v0.2.0 - Backend Feature Completion

Status: Planned

Purpose:

- complete remaining backend MVP endpoint/flow work
- close remaining backend behavior gaps
- preserve direct backend verification without frontend
- keep API, DB, Redis, and SPEC behavior aligned

---

## v0.3.0 - API Contract Readiness

Status: Planned

Purpose:

- stabilize request/response/error behavior before frontend integration
- ensure `docs/API.md` and `docs/API_TESTING.md` are clear enough for frontend work
- ensure success, failure, and persistence behavior is directly verifiable
- remove backend ambiguity before frontend integration begins

---

## v0.4.0 - Frontend Integration

Status: Planned

Purpose:

- connect the Next.js frontend to the verified backend
- implement loading, empty, success, and error states against real API behavior
- avoid frontend-invented backend fields or endpoints

---

## v0.5.0 - End-to-End MVP Hardening

Status: Planned

Purpose:

- full-stack validation
- bug fixing
- UX polish
- release readiness
- final MVP stabilization

## v0.6.0 - Deployment Preparation

Objectives:

- Prepare for the production environment
- Configure environment variables
- Set up CORS for production
- Configure PostgreSQL and Redis connections for production
- Conduct a security review
- Migration strategy
- Frontend + backend deployment preparation

## v1.0.0 - Public MVP Release

Objective:

- Validation of all MVP workflows in the deployment environment
- Documentation cleanup
- Launch for the first real users

Translated with DeepL.com (free version)
