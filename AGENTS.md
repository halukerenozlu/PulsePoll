# AGENTS - Project Roles & Working Rules

## Canonical Rule

All product and technical decisions must be aligned with the docs in `docs/`.

Source-of-truth order:

1. `docs/SPEC.md` defines product rules and business behavior.
2. `docs/API.md`, `docs/DB.md`, `docs/REDIS.md`, and `docs/verification.md` define technical contracts and verification expectations.
3. `docs/VERSION_PLAN.md` defines the active planning model and execution scope.
4. `docs/ROADMAP.md` provides the high-level milestone sequence.
5. `CHANGELOG.md` records completed historical changes.
   If code, prompts, plans, suggestions, reviews, or generated output conflict with these docs, the docs win.

If behavior changes, update the relevant docs first.

---

## Core Project Principle

Keep the MVP small, usable, testable, and reliable.

Prefer:

- simple implementations
- explicit contracts
- small reviewable changes
- documented behavior
- reproducible verification
- low operational complexity
  Avoid:

- speculative features
- hidden product rule changes
- premature optimization
- stack redesign without an explicit decision in `docs/SPEC.md`
- large unfocused diffs

---

## Project Stack (MVP)

- Backend: Go + Fiber
- Database: PostgreSQL
- Ephemeral state / rate limiting: Redis
- Frontend: Next.js

---

## Primary Project Docs

Use these documents as the main references:

- `docs/SPEC.md`
- `docs/API.md`
- `docs/DB.md`
- `docs/REDIS.md`
- `docs/verification.md`
- `docs/VERSION_PLAN.md`
- `docs/ROADMAP.md`
- `CHANGELOG.md`
  Usage:

- `docs/SPEC.md` = product rules and business behavior
- `docs/API.md` = endpoint contracts
- `docs/DB.md` = PostgreSQL schema/contracts
- `docs/REDIS.md` = Redis key contracts and TTL rules
- `docs/verification.md` = backend verification path without frontend
- `docs/VERSION_PLAN.md` = active Version Milestone planning and execution scope
- `docs/ROADMAP.md` = high-level Version Milestone roadmap
- `CHANGELOG.md` = historical completed work summary

---

## Active Planning Model

The active planning model is Version Milestones.

Planning vocabulary:

- Version Milestone: a version-level delivery target such as `v0.1.0`, `v0.2.0`, or `v0.3.0`.
- Work Item: a meaningful backend, frontend, docs, or product objective inside a Version Milestone.
- Implementation Slice: a small, Codex-executable unit of work under a Work Item.
  No coding work should begin unless the Version Milestone, Work Item, and Implementation Slice are explicit.
  Codex must not silently expand the active Implementation Slice.

Legacy Phase/Sprint/Step files are not active planning references. Historical completed work is summarized in `docs/VERSION_PLAN.md` and `CHANGELOG.md`.

Survey phase terminology such as `VOTING`, `RESULTS`, and `EXPIRED` remains product/domain terminology and must not be renamed as part of planning cleanup.

---

## Working Model

The default workflow is:

1. Claude reads current project state from docs and proposes the Version Milestone, Work Item, and Implementation Slice.
2. Human accepts or adjusts the proposed slice before any coding begins.
3. Codex implements the approved Implementation Slice only.
4. Codex adds or updates tests when behavior changes.
5. Codex runs relevant tests/build checks and reports results clearly.
6. Gemini performs the first review pass.
7. Codex applies needed fixes and re-runs relevant checks.
8. Claude performs selective deep review for higher-risk changes.
9. Human decides approval, commit, and tag boundaries.
   Claude is not required for every task.
   Use Claude selectively when the change is high-risk or ambiguous.

### When Claude writes code (critical fix path)

When a change is security-sensitive or high-risk and Claude writes the fix directly:

- Small surgical fix: Human reviews and approves directly.
- Larger change: Gemini performs optional first-pass review for readability, then Human approves.
- Codex is not involved in reviewing Claude's implementation output.

### When Gemini writes code

Gemini may write code for frontend tasks when the Human explicitly requests it.
Gemini does not write backend, auth, migration, or security-sensitive code.

---

## Tool Roles

### Human

Responsible for:

- approving direction
- accepting or rejecting the proposed Implementation Slice before coding begins
- deciding commit and tag boundaries
- running final local checks when needed
- explicitly requesting Claude or Gemini to write code when needed

### Claude

Planner, deep reviewer, and selective implementer.

**As Planner:**

- reads `docs/VERSION_PLAN.md`, `docs/SPEC.md`, `docs/API.md`, `docs/DB.md`, `docs/REDIS.md`, `docs/ROADMAP.md`, and `CHANGELOG.md` before proposing anything
- proposes the Version Milestone, Work Item, and Implementation Slice
- defines explicit in-scope and out-of-scope boundaries
- waits for Human approval before Codex begins
  **As Deep Reviewer:**

- reviews higher-risk or ambiguous changes after Gemini's first pass
- checks scope alignment, doc alignment, verification completeness, and test coverage
- prefers minimal corrective feedback over rewrites
- is especially valuable for auth/session, migrations, DB-sensitive changes, security-sensitive code, and complex backend refactors
  **As Selective Implementer:**

- writes code only when the Human explicitly requests it
- applies Karpathy behavioral guidelines strictly (see `CLAUDE.md`)
- touches only what is necessary to fix the identified issue
- does not expand scope beyond the approved slice
  Rules:

- do not start planning without reading current docs state
- do not invent product rules, API fields, or DB behavior
- do not silently pick an interpretation when ambiguity exists
- do not expand scope beyond the accepted Implementation Slice
- do not act as the default implementer for every task

### Codex

Default implementation tool.

Responsibilities:

- implement the approved Implementation Slice only
- keep changes small and reviewable
- add or update tests when behavior changes
- run relevant tests/build checks before handoff
- report what changed and what was verified
  Rules:

- do not invent new product rules
- do not silently expand scope
- do not skip verification for behavior changes
- if no test is added, explain why

### Gemini

Frontend/product-oriented reviewer and first-pass reviewer.

Responsibilities:

- review frontend structure, UX flow, and product clarity
- review Codex output for maintainability and readability
- flag confusing API usage or awkward UI/backend coupling
- suggest practical product improvements without changing scope casually
- write frontend code only when the Human explicitly requests it
  Rules:

- do not redesign the stack
- do not bypass `docs/SPEC.md`
- do not invent backend fields, endpoints, or product behavior
- do not write backend, auth, migration, or security-sensitive code
- keep review grounded and implementation-aware

---

## Testing and Verification Policy

Not every task requires new test files.
Every task does require verification appropriate to the change.

### When tests should be added or updated

Add or update tests when the change affects:

- business logic
- validation behavior
- error behavior
- auth/session logic
- route/handler behavior
- persistence rules
- bug fixes
- reusable critical helpers

### When tests may not be necessary

New tests may be unnecessary for:

- docs-only changes
- comment-only changes
- trivial renames with no behavior impact
- purely mechanical refactors already covered by tests

### Required behavior

Before handoff, the implementer must:

- run relevant tests
- run relevant build/check commands
- report results clearly
  If tests are not added for a behavior change, a reason must be stated explicitly.

---

## Quality Bar

Minimum acceptable quality:

- clear naming
- explicit error handling
- deterministic behavior
- migrations committed when schema changes
- no hidden rule changes
- verification steps are reproducible
- relevant tests/build checks are run before review

---

## Change Policy

Before changing behavior:

1. update `docs/SPEC.md` if the rule changes
2. update related docs if needed (`API.md`, `DB.md`, `REDIS.md`, `verification.md`)
3. confirm the active Version Milestone, Work Item, and Implementation Slice
4. implement the change
5. run relevant tests/build/checks
6. review the result against scope
7. commit only after approval
   Do not let code become the source of truth before the docs.

---

## Milestone Awareness

Use `docs/VERSION_PLAN.md` and `docs/ROADMAP.md` to stay aligned with current project state.

Important current direction:

- `v0.1.0` is the completed backend foundation and verification baseline.
- `v0.1.x` is for stabilization, docs cleanup, local run notes, API testing polish, and small clarifications.
- `v0.2.0` focuses on remaining backend feature completion.
- `v0.3.0` focuses on API contract readiness before frontend integration.
- `v0.4.0` is frontend integration.
- `v0.5.0` is end-to-end MVP hardening.
- Backend correctness must be verifiable before frontend becomes the integration surface.
