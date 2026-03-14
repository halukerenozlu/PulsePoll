# DB — Draft Schema (MVP)

This is a draft. Keep MVP minimal; normalize later if needed.

## users
- id (uuid, pk)
- email (text, unique)
- password_hash (text)
- display_name (text)
- status (enum: active, suspended, deleted)
- created_at, updated_at

## auth_sessions (refresh tokens)
- id (uuid, pk)
- user_id (uuid, fk users.id)
- refresh_token_hash (text, unique)
- user_agent (text)
- ip (text)
- expires_at (timestamptz)
- revoked_at (timestamptz, nullable)
- created_at

## surveys
- id (uuid, pk)
- creator_id (uuid, fk users.id)
- title (text)
- description (text, nullable)

Policy:
- visibility (enum: public, unlisted, private_pin)
- access_pin_hash (text, nullable)
- results_mode (enum: open_live, closed_hidden_until_end)
- max_votes_per_user (int, default 1)
- allow_vote_change_once (bool, default false)

Timing:
- created_at, updated_at
- vote_ends_at (timestamptz)
- results_ends_at (timestamptz)
- retention_ends_at (timestamptz)

Moderation:
- moderation_status (enum: approved, flagged, blocked, pending) default approved
- moderation_reason (text, nullable)

## survey_options
- id (uuid, pk)
- survey_id (uuid, fk surveys.id)
- text (text)
- position (int)
- vote_count (bigint, default 0)
- created_at

Constraint:
- unique(survey_id, position)

## reports
- id (uuid, pk)
- survey_id (uuid, fk surveys.id)
- reporter_user_id (uuid, nullable)
- reporter_guest_id (text, nullable)  # from cookie/guest_id
- reason (text)
- details (text, nullable)
- status (enum: open, reviewed, dismissed) default open
- created_at

## feedback (optional, but useful for MVP)
- id (uuid, pk)
- user_id (uuid, nullable)
- guest_id (text, nullable)
- message (text)
- page (text, nullable)
- created_at

## Notes
- MVP stores only aggregated counts (survey_options.vote_count).
- No raw vote events in Postgres for MVP.
