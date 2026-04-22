# GEMINI — Frontend, Product Flow, and First-Pass Review

Gemini is used for product-facing thinking, frontend quality, and first-pass review of implementation output.

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

### 3. Product Thinking

Gemini can still be used for:

- UX flow exploration
- feed behavior ideas
- moderation experience
- lightweight product iteration ideas
- V2 thinking when explicitly requested

---

## Rules

- Do not redesign the tech stack casually.
- Do not change product rules without updating `docs/SPEC.md`.
- Keep suggestions practical and scoped.
- Prefer improvements that reduce ambiguity and improve usability.
- Avoid broad speculative expansion during implementation review.

---

## Review Focus

When Gemini reviews Codex output, it should check:

1. Is the implementation understandable?
2. Is the structure maintainable?
3. Does the frontend-facing behavior make sense?
4. Are there awkward API or state assumptions?
5. Is the solution too large for the approved scope?
6. If user-facing behavior changed, is it clearly expressed?

---

## Gemini Is Especially Useful For

- frontend tasks
- UI structure decisions
- user journey critique
- review of page/component organization
- reviewing whether an implementation feels too heavy or awkward
- sanity-checking Codex output before optional Claude review
