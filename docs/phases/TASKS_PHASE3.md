# TASKS_PHASE3.md

## Phase 3 Goal

Complete the Phase 3 backend scope and make the backend verifiable without depending on the frontend.

Phase 3 is not only about adding backend code.  
It is also the phase where backend behavior becomes stable enough for frontend integration.

---

## Working Rules

- Stay strictly within Phase 3 scope.
- Keep changes aligned with:
  - `docs/SPEC.md`
  - `docs/API.md`
  - `docs/DB.md`
  - `docs/ROADMAP.md`
- Prefer minimal, reviewable increments.
- Do not introduce frontend-coupled assumptions into backend logic.
- Do not defer backend verification until frontend work begins.

---

## Main Objectives

### 1. Implement Phase 3 backend scope

Implement the backend functionality required for Phase 3 according to the project specification.

This includes, where applicable:

- route/handler implementation
- service/business logic
- request validation
- response shaping
- database reads/writes
- error handling
- configuration updates limited to Phase 3 needs

### 2. Make the backend verifiable without frontend

Backend correctness must be provable before Phase 4 starts.

This means Phase 3 must include:

- reproducible manual API checks
- happy-path verification
- negative-path verification
- persistence verification
- clear example requests and responses

### 3. Stabilize the API contract

Before frontend integration begins, the backend contract must be clear enough that frontend work can proceed without guessing.

This includes:

- stable endpoint paths and methods
- stable request body/query parameter expectations
- stable response shapes
- documented status codes
- documented error response format

### 4. Keep docs in sync

Documentation must reflect actual implementation by the end of Phase 3.

---

## Required Workstreams

### A. Backend Implementation

Complete the Phase 3 backend logic defined by the project docs.

Checklist:

- [ ] Implement required Phase 3 handlers/routes
- [ ] Implement required Phase 3 service logic
- [ ] Implement request validation rules
- [ ] Implement consistent response/error formatting
- [ ] Implement required database interactions
- [ ] Ensure behavior matches `docs/SPEC.md`
- [ ] Ensure API behavior matches `docs/API.md`
- [ ] Ensure persistence/schema assumptions match `docs/DB.md`

### B. Backend Verification Without Frontend

Create a reliable verification path that does not depend on UI work.

Checklist:

- [ ] Define a repeatable local verification flow for all Phase 3 endpoints
- [ ] Verify each Phase 3 endpoint through direct HTTP requests
- [ ] Verify normal success scenarios
- [ ] Verify invalid input scenarios
- [ ] Verify not-found / conflict / edge-case scenarios where applicable
- [ ] Verify database side effects for create/update/delete style operations
- [ ] Verify that health and dependency checks still pass after Phase 3 changes
- [ ] Add example request/response payloads for manual inspection

Accepted formats for verification artifacts may include one or more of:

- documented `curl` commands
- a Bruno collection
- integration tests
- script-based smoke checks

At least one reproducible verification method must exist in-repo or in docs.

### C. API Contract Readiness

Prepare the backend for Phase 4 frontend integration.

Checklist:

- [ ] Remove ambiguity from endpoint docs
- [ ] Confirm request fields are documented
- [ ] Confirm response fields are documented
- [ ] Confirm HTTP status code behavior is documented
- [ ] Confirm validation/error behavior is documented
- [ ] Confirm any pagination/filter/sort behavior if applicable

### D. Documentation Sync

Update docs so implementation and documentation remain aligned.

Checklist:

- [ ] Update `docs/SPEC.md` if Phase 3 behavior refined the spec
- [ ] Update `docs/API.md` with final endpoint behavior and examples
- [ ] Update `docs/DB.md` if schema/constraints/indexes changed
- [ ] Update any setup or verification documentation if needed
- [ ] Ensure `docs/ROADMAP.md` still reflects the real project sequence

---

## Deliverables

Phase 3 is considered complete only when the following exist:

- backend code for the approved Phase 3 scope
- synchronized documentation
- a reproducible backend verification workflow
- example request/response coverage for manual checking
- evidence that backend behavior is ready for frontend integration

---

## Exit Criteria

Phase 3 is complete only if all of the following are true:

- [ ] Application boots successfully in local development
- [ ] `/health` reports healthy dependencies
- [ ] Phase 3 backend flows are directly testable without frontend
- [ ] Success scenarios are verified
- [ ] Failure scenarios are verified
- [ ] Data persistence behavior is verified where applicable
- [ ] Documentation matches real implementation
- [ ] Backend contract is stable enough for Phase 4 integration

---

## Review Focus

When reviewing Phase 3 work, verify the following:

1. Scope control

- Is the work strictly limited to Phase 3?

2. Spec alignment

- Does implementation still match `docs/SPEC.md`?

3. API alignment

- Do actual endpoint behaviors match `docs/API.md`?

4. Database alignment

- Do persistence rules match `docs/DB.md`?

5. Verification quality

- Can the backend be validated without frontend?

6. Integration readiness

- Is Phase 4 able to proceed without guessing backend behavior?

---

## Execution Model

Phase 3 is planned and tracked at phase level in this document.

Implementation should be executed in sprint-sized increments.
Each sprint must have a clearly defined:

- goal
- scope
- out-of-scope items
- implementation checklist
- verification steps
- review focus
- completion criteria

Recommended structure:

- `docs/phases/TASKS_PHASE3.md` defines the Phase 3 umbrella scope
- sprint documents define the execution details for each implementation step

## Sprint Plan

### Sprint 1 — Verification Foundation

**Status:** Completed

Goal:

- establish the Phase 3 verification baseline
- make backend verification possible without frontend
- align Phase 3 planning docs and execution structure

Outcome:

- `docs/verification.md` created/refined
- `docs/phases/phase3/sprints/SPRINT1.md` created
- Phase 3 planning and workflow docs aligned for verification-first execution

### Sprint 2 — Vote Rate Limiting

**Status:** Completed

Goal:

- add Redis-backed IP rate limiting for vote endpoints
- return deterministic `429 TOO_MANY_REQUESTS`
- keep the change small, testable, and verification-friendly

Expected output:

- rate limiting for:
  - `POST /surveys/:id/vote`
  - `PUT /surveys/:id/vote`
- related tests
- manual verification steps for allowed vs rate-limited paths

### Future Sprints

Additional Phase 3 sprints will be defined after Sprint 2 review, based on:

- implementation results
- verification gaps
- docs/code alignment needs
