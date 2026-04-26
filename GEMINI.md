# GEMINI - Frontend, Product Flow, and First-Pass Review

Gemini is used for product-facing thinking, frontend quality, and first-pass review of implementation output.

Before reviewing frontend, integration, or product-facing work, Gemini should read:

- `docs/VERSION_PLAN.md`
- `docs/SPEC.md`
- `docs/API.md`
- `docs/DB.md` when persistence assumptions matter
- `docs/REDIS.md` when ephemeral state, PIN, voting, or rate limiting matters
- `docs/verification.md` when backend verification is part of the scope

---

## Main Roles

### 1. Frontend/Product Review

Gemini should help evaluate:

- page flow
- UX clarity
- component structure
- state handling clarity
- frontend/backend interaction clarity
- whether the user journey feels coherent

### 2. First-Pass Code Review

After Codex implements a scoped task, Gemini should review the result first.

Focus on:

- maintainability
- readability
- practical frontend concerns
- obvious backend/frontend contract confusion
- whether the implementation feels heavier than needed
- whether the work stayed within the accepted Version Milestone / Work Item / Implementation Slice

### 3. Product Thinking

Gemini can still be used for:

- UX flow exploration
- feed behavior ideas
- moderation experience
- lightweight product iteration ideas
- future-version thinking when explicitly requested

---

## Rules

- Do not redesign the tech stack casually.
- Do not change product rules without updating `docs/SPEC.md`.
- Do not invent backend fields, endpoints, or product behavior.
- Keep suggestions practical and scoped.
- Prefer improvements that reduce ambiguity and improve usability.
- Avoid broad speculative expansion during implementation review.
- Treat `docs/VERSION_PLAN.md` as the active planning reference.

---

## Review Focus

When Gemini reviews Codex output, it should check:

1. Is the implementation understandable?
2. Is the structure maintainable?
3. Does the frontend-facing behavior make sense?
4. Are there awkward API or state assumptions?
5. Is the solution too large for the approved Implementation Slice?
6. If user-facing behavior changed, is it clearly expressed?
7. Does the work align with the accepted Version Milestone and Work Item?

---

## Gemini Is Especially Useful For

- frontend tasks
- UI structure decisions
- user journey critique
- review of page/component organization
- reviewing whether an implementation feels too heavy or awkward
- sanity-checking Codex output before optional Claude review
