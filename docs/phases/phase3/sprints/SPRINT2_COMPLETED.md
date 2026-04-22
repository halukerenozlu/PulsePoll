# Phase 3 — Sprint 2: Vote Rate Limiting

## Status

Planned

---

## Sprint Goal

Add Redis-backed IP rate limiting for vote endpoints so that abusive repeated voting requests can be throttled without changing the rest of the voting flow.

This sprint should stay small, backend-only, testable, and verification-friendly.

---

## Why This Sprint Exists

Phase 3 requires backend behavior to be reliable and directly verifiable before frontend integration begins.

The project docs already define:

- `429 TOO_MANY_REQUESTS` as a valid API error class
- `rl:ip:{ip}:vote` as the Redis key contract for vote rate limiting

This sprint closes that docs/code gap with a minimal implementation slice.

---

## In Scope

This sprint includes:

- adding IP-based rate limiting for:
  - `POST /surveys/:id/vote`
  - `PUT /surveys/:id/vote`
- using Redis key contract aligned with `docs/REDIS.md`
- using TTL aligned with `docs/REDIS.md`
- returning deterministic `429 TOO_MANY_REQUESTS` responses when limit is exceeded
- adding/updating tests for allowed and blocked request paths
- documenting manual verification steps in the implementation handoff summary

---

## Out of Scope

This sprint does **not** include:

- app-wide/global rate limiting
- brute-force protection for PIN verification
- auth/session refactors
- new frontend work
- broader Redis redesign
- unrelated route cleanup
- feature expansion beyond vote endpoint rate limiting

---

## Expected Files Touched

Likely files include:

- vote-related handlers
- vote-related services or middleware
- Redis integration code for vote rate limiting
- vote-related tests

Docs should remain unchanged unless implementation reveals a real mismatch.

---

## Roles

- **Codex**
  - implement only the approved rate-limiting scope
  - add/update tests
  - run relevant checks before handoff

- **Gemini**
  - perform first-pass review
  - check clarity, maintainability, and scope discipline

- **Claude**
  - optional selective deep review
  - only if the patch grows riskier than expected

- **Human**
  - approve the scope
  - review the result
  - decide commit boundaries

- **ChatGPT**
  - keep scope narrow
  - provide implementation and review prompts
  - summarize current phase/sprint state

---

## Implementation Checklist

- [x] Add Redis-backed IP rate limiting for `POST /surveys/:id/vote`
- [x] Add Redis-backed IP rate limiting for `PUT /surveys/:id/vote`
- [x] Use Redis key shape `rl:ip:{ip}:vote`
- [x] Use TTL of 60 seconds
- [x] Return deterministic `429 TOO_MANY_REQUESTS` on limit exceed
- [x] Keep existing vote logic intact for non-rate-limited requests
- [x] Add or update tests for allowed path
- [x] Add or update tests for blocked path
- [x] Run relevant tests/build checks
- [x] Report manual verification guidance

---

## Verification Checklist

Sprint 2 is not complete unless the following are true:

- [x] normal vote request still works when under the limit
- [x] repeated vote requests hit the rate limit as expected
- [x] blocked requests return `429 TOO_MANY_REQUESTS`
- [x] Redis key usage matches `docs/REDIS.md`
- [x] tests cover both allowed and blocked cases
- [x] implementation stays limited to vote endpoint rate limiting

---

## Review Focus

When reviewing Sprint 2, verify:

### 1. Scope Control

- Is the patch limited to vote endpoint IP rate limiting only?

### 2. Contract Alignment

- Does Redis usage match `docs/REDIS.md`?
- Does HTTP behavior match `docs/API.md`?

### 3. Safety

- Does the patch avoid breaking existing vote, consent, phase, and PIN behavior?

### 4. Test Quality

- Are both allowed and blocked flows covered?

### 5. Minimalism

- Were unnecessary refactors avoided?

---

## Done Criteria

Sprint 2 is complete only if all of the following are true:

- [x] vote endpoints are rate-limited by IP
- [x] limit-exceeded requests return deterministic `429`
- [x] Redis contract is respected
- [x] tests were added/updated and run
- [x] the patch remained small and verification-friendly

---

## Handoff to Next Sprint

After Sprint 2, the next sprint should be chosen based on:

- implementation quality
- review feedback
- verification gaps
- remaining docs/code mismatches
