# REDIS — Key Schema (MVP)

## TTL Rule
Vote receipts must expire at (retention_ends_at - now) for that survey.

## Consent cookie
Guests must accept consent to receive `guest_id` cookie.
Without guest_id, voting endpoints return 403 CONSENT_REQUIRED.

## Vote receipts (limit enforcement + one-time change)
# Registered
vote:survey:{surveyId}:user:{userId} -> JSON
  { "votes_used": 1, "last_option_id": "...", "change_used": false }
TTL: until retention_ends_at

# Guest
vote:survey:{surveyId}:guest:{guestId} -> JSON
  { "votes_used": 1, "last_option_id": "...", "change_used": false }
TTL: until retention_ends_at

## PIN verification state (short TTL)
pinok:survey:{surveyId}:user:{userId} -> "1"
pinok:survey:{surveyId}:guest:{guestId} -> "1"
TTL: 30 minutes (or until vote_ends_at, whichever is smaller)

## PIN brute-force protection
pinfail:survey:{surveyId}:guest:{guestId} -> integer counter (INCR)
TTL: 15 minutes

## Rate limiting (simple)
rl:ip:{ip}:vote -> integer counter
TTL: 60 seconds
