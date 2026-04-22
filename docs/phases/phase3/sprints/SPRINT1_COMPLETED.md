# Phase 3 — Sprint 1: Verification Foundation

## Sprint Goal

Establish the verification foundation for Phase 3 so that backend work can be validated without relying on the frontend.

This sprint creates the baseline verification path before deeper Phase 3 backend implementation continues.

---

## Why This Sprint Exists

Phase 3 should not be the phase where backend behavior is only understood later through frontend integration.

Before larger backend changes continue, the project needs:

- a clear backend verification workflow
- a repeatable local verification path
- a stable baseline for review
- documentation that reduces ambiguity
- a common execution model for later Phase 3 sprints

---

## In Scope

This sprint includes:

- creating or refining `docs/verification.md`
- defining the baseline local backend verification flow
- documenting `/health` verification
- defining how Phase 3 endpoint checks should be recorded
- establishing the Phase 3 sprint execution structure
- ensuring docs reflect the verification-first approach

Optional, only if small and low-risk:

- a lightweight reusable smoke-check helper
- example `curl` commands for early manual verification

---

## Out of Scope

This sprint does **not** include:

- frontend implementation
- Phase 4 integration work
- major new product features
- large schema redesign
- deployment expansion
- broad infrastructure redesign
- non-Phase-3 refactors
- large unrelated cleanup

---

## Expected Files Touched

Likely files include:

- `docs/ROADMAP.md`
- `docs/verification.md`
- `docs/phases/TASKS_PHASE3.md`
- `docs/phases/phase3/sprints/SPRINT1.md`

Optional, only if truly needed:

- small helper script files for smoke checks
- tightly related README updates

---

## Roles

- **Codex**
  - implement the approved Sprint 1 scope only
  - keep the work small, clear, and reviewable
  - add or update tests only when behavior actually changes
  - run relevant checks before handoff

- **Gemini**
  - perform the first review pass
  - review clarity, maintainability, and practical usability of the sprint output

- **Claude**
  - optional selective deep reviewer
  - only used when the sprint changes become higher-risk than expected

- **Human**
  - run local verification when needed
  - confirm actual project behavior
  - decide approval and commit boundaries

- **ChatGPT**
  - keep scope clear
  - define the execution shape
  - help coordinate next steps

---

## Implementation Checklist

- [ ] Create or refine `docs/verification.md`
- [ ] Document local backend startup verification
- [ ] Document `/health` verification flow
- [ ] Define a reusable endpoint verification template
- [ ] Define success-path verification expectations
- [ ] Define failure-path verification expectations
- [ ] Define persistence verification expectations
- [ ] Ensure Sprint 1 aligns with `docs/phases/TASKS_PHASE3.md`
- [ ] Ensure `docs/ROADMAP.md` still reflects Phase 3 verification-first logic

Optional:

- [ ] Add a lightweight smoke-check helper if it remains small and low-risk

---

## Verification Checklist

Sprint 1 is not complete unless the following are true:

- [ ] documentation exists in the repo
- [ ] local startup flow is written clearly
- [ ] `/health` check is documented and works locally
- [ ] verification steps can be followed without frontend
- [ ] the endpoint verification template is reusable for later Phase 3 work
- [ ] the docs are understandable without relying on old chat context

---

## Review Focus

When reviewing Sprint 1, verify:

### 1. Scope Control

- Is the work limited to verification foundation and planning support?

### 2. Clarity

- Can a future reader understand how to verify backend work locally?

### 3. Consistency

- Does Sprint 1 align with:
  - `docs/SPEC.md`
  - `docs/verification.md`
  - `docs/ROADMAP.md`
  - `docs/phases/TASKS_PHASE3.md`

### 4. Practicality

- Can the documented verification steps actually be followed?

### 5. Minimalism

- Was unnecessary feature work avoided?

---

## Done Criteria

Sprint 1 is complete only if all of the following are true:

- [ ] `docs/verification.md` exists and is usable
- [ ] `docs/phases/phase3/sprints/SPRINT1.md` exists
- [ ] backend verification can begin without frontend
- [ ] project docs explain how Phase 3 work should be validated
- [ ] the sprint produced a usable base for later Phase 3 implementation and review

---

## Handoff to Next Sprint

Sprint 2 should begin only after Sprint 1 establishes a clear verification baseline.

Sprint 2 can then focus on actual Phase 3 backend implementation while using the Sprint 1 verification structure as the default validation path.
