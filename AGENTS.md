# Codex Role: Implementation

Implement strictly from docs/* (SPEC/API/DB/REDIS).
Small PR-style steps. Do not invent new rules without updating docs/SPEC.md first.

MVP priorities:
- Auth (register/login/refresh cookie)
- Survey create/read/feed
- Vote + one-time change (only when max_votes_per_user == 1)
- PIN flow (if enabled)
- Consent gating for guest voting
- Basic moderation filter + report endpoint

Quality:
- Clear folder structure
- Migrations included
- Basic unit tests for phase calculation & vote rules
