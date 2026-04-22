# AGENTS — Project Roles & Working Rules

## Canonical Rule

All product and technical decisions must be aligned with the docs in `docs/`.

If any code, prompt, plan, suggestion, or generated output conflicts with:

- `docs/SPEC.md`
- `docs/API.md`
- `docs/DB.md`
- `docs/REDIS.md`
- `docs/verification.md`
- `docs/ROADMAP.md`
- `docs/phases/TASKS.md`
- `docs/phases/TASKS_PHASE2.md`
- `docs/phases/TASKS_PHASE3.md`

the docs win.

If a behavior or rule changes, update the relevant docs first.

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
- `docs/ROADMAP.md`
- `docs/phases/TASKS.md`
- `docs/phases/TASKS_PHASE2.md`
- `docs/phases/TASKS_PHASE3.md`

Usage:

- `docs/SPEC.md` = product rules and business behavior
- `docs/API.md` = endpoint contracts
- `docs/DB.md` = PostgreSQL schema/contracts
- `docs/REDIS.md` = Redis key contracts and TTL rules
- `docs/verification.md` = backend verification path without frontend
- `docs/ROADMAP.md` = project phase sequence and exit logic
- `docs/phases/TASKS.md` = MVP implementation order
- `docs/phases/TASKS_PHASE2.md` = Phase 2 stabilization scope
- `docs/phases/TASKS_PHASE3.md` = Phase 3 backend + verification scope

---

## Working Model

The default workflow is:

1. Human + ChatGPT define scope
2. Codex implements the approved scope
3. Codex adds or updates tests when behavior changes
4. Codex runs relevant tests/build checks and reports results
5. Gemini performs the first review pass
6. Codex applies needed fixes and re-runs relevant checks
7. Claude performs selective deep review only for higher-risk changes
8. Human decides approval, commit, and tag boundaries

Claude is not required for every task.
Use Claude selectively when the change is high-risk or ambiguous.

---

## Tool Roles

### Human

Responsible for:

- approving direction
- deciding commit boundaries
- running final local checks when needed
- accepting or rejecting changes

### ChatGPT

Responsible for:

- planning
- scope control
- sprint/task shaping
- workflow design
- prompt design
- summarizing state and next steps

### Codex

Default implementation tool.

Responsibilities:

- implement the approved scope only
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
- review Codex output for maintainability and clarity
- flag confusing API usage or awkward UI/backend coupling
- suggest practical product improvements without changing scope casually

Rules:

- do not redesign the stack
- do not bypass `docs/SPEC.md`
- keep review grounded and implementation-aware

### Claude

Selective deep reviewer.

Use mainly for:

- auth/session logic
- migrations / DB-sensitive changes
- security-sensitive code
- complex backend refactors
- higher-risk changes before approval

Rules:

- do not act as the default implementer
- review only the approved scope
- prefer minimal corrective feedback
- avoid broad unsolicited rewrites

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
3. implement the change
4. run relevant tests/build/checks
5. review the result against scope
6. commit only after approval

Do not let code become the source of truth before the docs.

---

## Phase Awareness

Use the roadmap and phase task docs to stay aligned with current project state.

Important current direction:

- Phase 1 is complete
- Phase 2 is complete
- Phase 3 focuses on backend features + verification
- frontend integration belongs to Phase 4
- backend correctness must be verifiable before frontend becomes the integration surface
