# PulsePoll Project Spec

## Purpose

This file is maintained as the main English reference summary for the PulsePoll project.

Purpose:

- help resume the project quickly after a break
- collect the product logic, technical direction, and working model in one place
- clarify agent roles and the project workflow
- provide a high-level summary of the existing English contract documents
  This file is a high-level guide.
  Detailed and binding technical contracts still live under `docs/`.

---

## Project Summary

PulsePoll is an MVP survey platform built around temporary survey lifecycles.

Core idea:

- surveys are short-lived
- the voting window is limited
- result visibility depends on the product phase
- the system is kept small, simple, and verifiable
  This project is currently developed as a public project.

---

## Core Principle

In this project, documents are the source of truth.

If code, plans, prompts, review notes, or any other output conflict with the following files, the documents win:

- `docs/SPEC.md`
- `docs/API.md`
- `docs/DB.md`
- `docs/REDIS.md`
- `docs/verification.md`
- `docs/VERSION_PLAN.md`
- `docs/ROADMAP.md`
- `CHANGELOG.md`
  If behavior changes, the relevant document must be updated first.

---

## Tech Stack

### Backend

- Go
- Fiber

### Database

- PostgreSQL

### Ephemeral State / Rate Limiting / Short-Lived Records

- Redis

### Frontend

- Next.js

### Local Development

- Docker Compose

---

## Current Project Status

The active planning model is now the Version Milestone model.

Current completed baseline:

- `v0.1.0` - Backend Foundation and Verification Baseline
  Next planning area, depending on repository state:

- `v0.1.x` - small stabilization, documentation cleanup, local verification notes, API testing guide polish
- `v0.2.0` - completion of remaining backend MVP endpoint/flow work
  Important note:

- `v0.1.0` consolidates the completed baseline work from the old planning model.
- This does not mean the older, broader backend feature/readiness scope is fully complete.
- Remaining backend and API readiness work moves into the relevant future Version Milestones.

---

## Version Milestone Model

Planning terms:

- Version Milestone: a version-level delivery target such as `v0.1.0` or `v0.2.0`.
- Work Item: a meaningful backend, frontend, docs, or product objective inside a Version Milestone.
- Implementation Slice: a small, clearly scoped unit of work under a Work Item that Codex can implement.
  Before coding work starts, the Version Milestone, Work Item, and Implementation Slice must be explicit.
  Codex must not silently expand the active Implementation Slice scope.

Active planning source:

- `docs/VERSION_PLAN.md`
  High-level sequence:

- `v0.1.0` - Backend Foundation and Verification Baseline
- `v0.1.x` - Stabilization and Docs Cleanup
- `v0.2.0` - Backend Feature Completion
- `v0.3.0` - API Contract Readiness
- `v0.4.0` - Frontend Integration
- `v0.5.0` - End-to-End MVP Hardening

---

## Product Rules Summary

### Survey Phases

Surveys are evaluated in three main states:

- `VOTING`
- `RESULTS`
- `EXPIRED`
  These terms are product/domain terminology and must not be confused with old phase terms from the project planning model.

General logic:

- `now < vote_ends_at` -> `VOTING`
- `vote_ends_at <= now < results_ends_at` -> `RESULTS`
- `now >= results_ends_at` -> `EXPIRED`

### Default Durations

MVP defaults:

- `vote_ends_at = created_at + 24h`
- `results_ends_at = created_at + 48h`
- `retention_ends_at = created_at + 48h`

### Visibility Types

- `PUBLIC`
- `UNLISTED`
- `PRIVATE_PIN`

### Result Visibility

- `OPEN_LIVE`
- `CLOSED_HIDDEN_UNTIL_END`

### User Rules

- registered users can create surveys
- guest users cannot create surveys
- guest users must accept the required cookie to vote
- guest users can still browse and view allowed results without accepting cookies

### Vote Rules

- `max_votes_per_user >= 1`
- `allow_vote_change_once` is meaningful only when `max_votes_per_user == 1`
- vote changes are allowed only during the `VOTING` phase and at most once
- the same registered user or guest can vote for the same option multiple times if the rule allows it

### Moderation

At the MVP level:

- a basic keyword filter exists during survey creation
- creation is rejected if an inappropriate term is detected
- a report endpoint exists

---

## Consent (Guest Voting) Logic

Guest users can always:

- browse the public feed
- view survey details
- view allowed results
  But to vote or change a vote:

- they must accept the required service cookie
- a short-lived `guest_id` is stored in that cookie
  This structure serves these purposes:

- reduce spam
- limit repeat-voting abuse
- enforce vote limits
- track the one-time vote change rule
- remember short-lived PIN verification status

---

## Data Retention Logic

### PostgreSQL

Persistent and core data is stored here.

Example tables:

- `users`
- `auth_sessions`
- `surveys`
- `survey_options`
- `reports`
- `feedback` (optional)
  MVP approach:

- no raw vote event log is stored
- aggregate counts are stored
- `survey_options.vote_count` is the critical counter field

### Redis

Temporary and TTL-focused data is stored here.

Example uses:

- vote receipts
- guest-based vote limit tracking
- one-time vote change state
- PIN verification state
- brute-force prevention counters
- rate limiting

---

## API Summary

Base path:

- `/api/v1`
  Important endpoint groups:

- Auth
- Consent
- Surveys
- Feed
- PIN verify
- Vote / vote change
- Results
- Report
- Feedback (optional)
  Important error classes:

- `400 BAD_REQUEST`
- `401 UNAUTHORIZED`
- `403 FORBIDDEN`
- `404 NOT_FOUND`
- `429 TOO_MANY_REQUESTS`
- `500 INTERNAL_SERVER_ERROR`
  Some specifically tracked error codes:

- `CONSENT_REQUIRED`
- `PIN_REQUIRED`
- `PHASE_NOT_VOTING`
- `VOTE_CHANGE_NOT_ALLOWED`

---

## Verification Approach

In this project, backend correctness is not delegated to the frontend.

`docs/verification.md` exists to:

- verify the backend without a UI
- perform local startup checks
- verify `/health` status
- check endpoint success / failure scenarios
- verify persistence effects
- create a reproducible verification path by Version Milestone, Work Item, and Implementation Slice
  Main idea:

- the frontend should not be the first place where the backend is tested

---

## Agent Roles

### Human

- approves direction
- accepts or rejects the proposed Implementation Slice before coding begins
- defines commit and tag boundaries
- makes the final decision on every approval, commit, and tag
- explicitly requests Claude or Gemini to write code when needed
- can run local verification personally when needed

### Claude

Claude takes on three roles: Planner, Deep Reviewer, and Selective Implementer.

See `CLAUDE.md` for the full behavioral guidelines including the Karpathy principles.

**As Planner:**

- reads `docs/VERSION_PLAN.md`, `docs/SPEC.md`, `docs/API.md`, `docs/DB.md`, `docs/REDIS.md`, `docs/ROADMAP.md`, and `CHANGELOG.md` before proposing anything
- proposes the Version Milestone, Work Item, and Implementation Slice
- waits for Human approval before Codex begins
  **As Deep Reviewer:**

- performs selective deep review for high-risk or ambiguous changes after Gemini's first pass
- is especially valuable for auth/session, migrations, DB-sensitive changes, security-sensitive code, vote rule enforcement, PIN logic, and complex backend refactors
- prefers minimal corrective feedback over rewrites
  **As Selective Implementer:**

- writes code only when the Human explicitly requests it
- for small surgical fixes: Human reviews and approves directly
- for larger changes: Gemini does optional first-pass review for readability, then Human approves
- Codex is not involved in reviewing Claude's implementation output

### Codex

Codex is the default implementer.

Responsibilities:

- implements the approved Implementation Slice only
- adds or updates tests when behavior changes
- runs relevant test/build commands and reports results clearly
- reports what changed and what was verified
  Rules:

- do not invent product rules or silently expand scope
- do not skip verification for behavior changes
- if no test is added, explain why

### Gemini

Gemini is the first reviewer and is especially useful for frontend and product flow.

Responsibilities:

- perform first-pass review of Codex output
- evaluate frontend structure, UX flow, and product clarity
- flag confusing API usage or awkward UI/backend coupling
- check maintainability and readability
- write frontend code only when the Human explicitly requests it
  Rules:

- do not redesign the stack
- do not invent backend fields, endpoints, or product behavior
- do not write backend, auth, migration, or security-sensitive code

---

## Default Workflow

Normal flow:

1. Claude reads current project state from docs and proposes the Version Milestone, Work Item, and Implementation Slice
2. Human accepts or adjusts the proposed slice before any coding begins
3. Codex implements the approved Implementation Slice
4. Codex adds or updates required tests
5. Codex runs relevant test/build steps and reports results
6. Gemini performs first-pass review
7. Codex applies needed fixes and verifies again
8. Claude performs selective deep review for higher-risk changes
9. Human decides approval, commit, and tag boundaries
   Critical fix path (when Claude writes code):

- Small surgical fix: Human reviews and approves directly
- Larger change: Gemini does optional first-pass review, then Human approves
- Codex is not involved in reviewing Claude's implementation output
  Gemini code path (when Gemini writes frontend code):

- Human explicitly requests it
- Gemini applies changes scoped to the approved slice
- Human reviews and approves

---

## Test Policy

Not every change requires a new test file.
Every change does require appropriate verification.

### Cases Where Tests Should Be Added

- business logic changes
- validation changes
- error behavior changes
- auth/session flow
- route/handler behavior
- bug fixes
- critical helper functions
- persistence rules

### Cases Where New Tests May Not Be Required

- docs-only changes
- comment-only changes
- small renames that do not affect behavior
- mechanical refactors already covered by existing tests
  Rule:

- if behavior changes, add a test or clearly explain why no test was added
- before handoff, run the relevant test/build commands

---

## Quality Bar

Minimum acceptable level:

- clear naming
- explicit error handling
- deterministic behavior
- no unnecessary scope expansion
- migrations committed when required
- reproducible verification steps
- relevant test/build commands have been run

---

## Document Map

### Product and Behavior

- `docs/SPEC.md`

### API Contract

- `docs/API.md`

### Database Contract

- `docs/DB.md`

### Redis Contract

- `docs/REDIS.md`

### Verification Flow

- `docs/verification.md`
- `docs/API_TESTING.md`

### Planning and History

- `docs/VERSION_PLAN.md`
- `docs/ROADMAP.md`
- `CHANGELOG.md`

### High-Level Summaries

- `docs/VISION.md`
- `docs/ARCHITECTURE.md`

---

## Local Run Summary

Example command to start services:

```bash
docker compose -p pulsepoll up --build
```

Health check:

```text
GET http://localhost:8080/health
```

Expected healthy response shape:

```json
{
  "db": "up",
  "ok": true,
  "redis": "up"
}
```

---

## How This File Should Be Used

This file should be used as:

- a quick resume guide
- a project memory file
- a high-level public reference
- an agent and workflow summary
  For detailed changes, the primary source should still be the technical contract files.
