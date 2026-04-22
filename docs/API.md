# API — Endpoints (MVP)

Base path: `/api/v1`

## Common Errors (recommended)
- 400 BAD_REQUEST: validation errors
- 401 UNAUTHORIZED: auth required / invalid tokens
- 403 FORBIDDEN: policy violation
  - code: CONSENT_REQUIRED (guest voting requires service cookie)
  - code: PIN_REQUIRED (PIN not verified for PRIVATE_PIN)
  - code: PHASE_NOT_VOTING (voting not allowed in this phase)
  - code: VOTE_CHANGE_NOT_ALLOWED (change used / not enabled)
- 404 NOT_FOUND
- 429 TOO_MANY_REQUESTS (rate limit)
- 500 INTERNAL_SERVER_ERROR

### Error shape (suggested)
```json
{ "error": { "code": "CONSENT_REQUIRED", "message": "..." } }
```

## Auth
POST /auth/register
POST /auth/login
POST /auth/refresh      # rotates refresh cookie
POST /auth/logout       # revokes refresh token
GET  /me                # current user

## Consent (Guest)
POST /consent/accept
- sets guest_id cookie (HttpOnly) for guest voting
- response: { ok: true }

## Surveys
POST /surveys (auth required)
- runs moderation filter
- sets default timestamps if not provided (24/24/48)

GET /surveys/{id}
- returns survey + options
- includes computed fields: phase, can_vote, results_visible, requires_pin

GET /feed
- query: sort=new, visibility=public
- optional: search

GET /surveys/{id}/results
- returns counts + percentages if results_visible == true

## PIN (for PRIVATE_PIN)
POST /surveys/{id}/pin/verify
- body: { pin }
- on success sets short-lived pin_ok state in Redis for this guest/user

## Voting
POST /surveys/{id}/vote
- body: { option_id, pin? }
- rules:
  - phase must be VOTING
  - if PRIVATE_PIN => require pin_ok
  - guests require consent/guest_id cookie
  - enforce max_votes_per_user (user or guest)
  - increment Postgres survey_options.vote_count atomically

PUT /surveys/{id}/vote
- body: { new_option_id, pin? }
- allowed only if:
  - phase == VOTING
  - max_votes_per_user == 1
  - allow_vote_change_once == true
  - change not used before (user/guest)
  - requires pin_ok if PRIVATE_PIN
  - guests require consent/guest_id cookie

## Moderation / Safety
POST /surveys/{id}/report
- body: { reason, details? }

## Feedback (optional)
POST /feedback
