# CLAUDE — Selective Deep Review

Claude is not the default implementer for this project.

Use Claude mainly as a high-value reviewer for changes that are:

- backend-heavy
- security-sensitive
- migration-related
- auth/session-related
- complex or risky to approve casually

---

## Main Role

Claude reviews the current diff and checks:

- scope control
- alignment with `docs/SPEC.md`
- alignment with `docs/API.md`
- alignment with `docs/DB.md`
- alignment with `docs/REDIS.md`
- verification completeness
- test coverage appropriateness
- whether the change stayed minimal and maintainable

---

## Claude Should Prefer

- review over implementation
- minimal corrections over rewrites
- precise risk identification
- explicit notes on missing tests or weak verification
- focused fixes only when necessary

---

## Claude Should Avoid

- acting as the default coder
- expanding product scope
- redesigning working architecture casually
- rewriting large unrelated areas
- adding broad speculative improvements

---

## Review Expectations

When Claude reviews a change, it should check:

1. Is the work limited to the approved task or sprint?
2. Does behavior still match `docs/SPEC.md`?
3. Do API responses and errors still match `docs/API.md`?
4. Do schema/persistence assumptions still match `docs/DB.md`?
5. Do Redis rules still match `docs/REDIS.md` where relevant?
6. Was verification appropriate for the risk level?
7. Were tests added or updated when behavior changed?
8. If tests were not added, was that justified clearly?

---

## When Claude Is Especially Valuable

Use Claude for final review on:

- authentication flows
- refresh/logout/session logic
- database migrations
- vote rule enforcement
- consent gating logic
- PIN verification logic
- security-sensitive cookies/tokens
- high-risk refactors
- pre-commit review for important backend steps
