# API - Endpoints (MVP)

Base path: `/api/v1`

## Common Errors

- 400 BAD_REQUEST: validation errors or invalid request body
- 401 UNAUTHORIZED: auth required / invalid tokens
- 403 FORBIDDEN: policy violation
  - code: CONSENT_REQUIRED (guest voting requires service cookie)
  - code: PIN_REQUIRED (PIN not verified for PRIVATE_PIN)
  - code: PHASE_NOT_VOTING (voting not allowed in this phase)
  - code: VOTE_CHANGE_NOT_ALLOWED (change used / not enabled)
  - code: MAX_VOTES_REACHED (user or guest has reached the vote limit)
- 404 NOT_FOUND: requested resource does not exist
- 409 CONFLICT: duplicate resource conflict
- 429 TOO_MANY_REQUESTS: rate limit
- 500 INTERNAL_SERVER_ERROR: unexpected server error

Error shape:

```json
{ "error": { "code": "CONSENT_REQUIRED", "message": "..." } }
```

## Auth

### POST /auth/register

Request:

- email (string, required)
- password (string, required)
- display_name (string, required)

Response: 201 Created

```json
{
  "access_token": "...",
  "user": {
    "id": "...",
    "email": "...",
    "display_name": "..."
  }
}
```

Errors:

- 400 BAD_REQUEST - validation failure or invalid request body
- 409 CONFLICT - email already registered
- 500 INTERNAL_SERVER_ERROR - failed to create user or session

---

### POST /auth/login

Request:

- email (string, required)
- password (string, required)

Response: 200 OK

```json
{
  "access_token": "...",
  "user": {
    "id": "...",
    "email": "...",
    "display_name": "..."
  }
}
```

Errors:

- 400 BAD_REQUEST - validation failure or invalid request body
- 401 UNAUTHORIZED - invalid credentials
- 500 INTERNAL_SERVER_ERROR - failed to load user or create session

---

### POST /auth/refresh

Request:

- No request body fields.

Response: 200 OK

```json
{
  "access_token": "...",
  "user": {
    "id": "...",
    "email": "...",
    "display_name": "..."
  }
}
```

Errors:

- 401 UNAUTHORIZED - missing, invalid, expired, or revoked refresh token
- 500 INTERNAL_SERVER_ERROR - failed to load user or rotate session

---

### POST /auth/logout

Request:

- No request body fields.

Response: 200 OK

```json
{ "ok": true }
```

Errors:

- 500 INTERNAL_SERVER_ERROR - failed to revoke session

---

### GET /me

Request:

- No request body fields.

Response: 200 OK

```json
{
  "id": "...",
  "email": "...",
  "display_name": "..."
}
```

Errors:

- 401 UNAUTHORIZED - missing or invalid access token
- 500 INTERNAL_SERVER_ERROR - failed to load user

---

## Consent (Guest)

### POST /consent/accept

Request:

- No request body fields.

Response: 200 OK

```json
{ "ok": true }
```

Errors:

- 500 INTERNAL_SERVER_ERROR - failed to generate guest id

---

## Surveys

### POST /surveys

Request:

- title (string, required)
- description (string, optional)
- options (array of strings, required)
- visibility (string, required: "public" | "unlisted" | "private_pin")
- access_pin (string, required when visibility is "private_pin", otherwise optional)
- results_mode (string, required: "open_live" | "closed_hidden_until_end")
- max_votes_per_user (number, optional; defaults to 1)
- allow_vote_change_once (boolean, optional; only valid when max_votes_per_user is 1)
- vote_ends_at (string, optional; RFC3339 timestamp)
- results_ends_at (string, optional; RFC3339 timestamp)
- retention_ends_at (string, optional; RFC3339 timestamp)

Response: 201 Created

```json
{
  "id": "...",
  "creator_id": "...",
  "title": "...",
  "description": "...",
  "visibility": "public",
  "results_mode": "open_live",
  "max_votes_per_user": 1,
  "allow_vote_change_once": false,
  "created_at": "2026-05-25T00:00:00Z",
  "vote_ends_at": "2026-05-26T00:00:00Z",
  "results_ends_at": "2026-05-27T00:00:00Z",
  "retention_ends_at": "2026-05-27T00:00:00Z",
  "phase": "VOTING",
  "can_vote": true,
  "results_visible": true,
  "requires_pin": false,
  "options": [
    {
      "id": "...",
      "text": "...",
      "position": 1
    }
  ]
}
```

`description` is omitted when empty.

Errors:

- 400 BAD_REQUEST - validation failure or invalid request body
- 401 UNAUTHORIZED - auth required or invalid access token
- 403 FORBIDDEN - moderation filter blocked survey content
- 404 NOT_FOUND - created survey could not be loaded
- 500 INTERNAL_SERVER_ERROR - failed to hash PIN, create survey, or load survey

---

### GET /surveys/{id}

Request:

- No request body fields.

Response: 200 OK

```json
{
  "id": "...",
  "creator_id": "...",
  "title": "...",
  "description": "...",
  "visibility": "public",
  "results_mode": "open_live",
  "max_votes_per_user": 1,
  "allow_vote_change_once": false,
  "created_at": "2026-05-25T00:00:00Z",
  "vote_ends_at": "2026-05-26T00:00:00Z",
  "results_ends_at": "2026-05-27T00:00:00Z",
  "retention_ends_at": "2026-05-27T00:00:00Z",
  "phase": "VOTING",
  "can_vote": true,
  "results_visible": true,
  "requires_pin": false,
  "options": [
    {
      "id": "...",
      "text": "...",
      "position": 1
    }
  ]
}
```

Computed fields:

- phase (string: "VOTING" | "RESULTS" | "EXPIRED")
- can_vote (boolean: true | false)
- results_visible (boolean: true | false)
- requires_pin (boolean: true | false)

`description` is omitted when empty.

Errors:

- 400 BAD_REQUEST - invalid survey id
- 404 NOT_FOUND - survey not found
- 500 INTERNAL_SERVER_ERROR - failed to load survey

---

### GET /feed

Request:

- sort (string, optional query parameter; only "new" is supported)
- visibility (string, optional query parameter; only "public" is supported)
- search (string, optional query parameter)

Response: 200 OK

```json
{
  "items": [
    {
      "id": "...",
      "title": "...",
      "description": "...",
      "visibility": "public",
      "results_mode": "open_live",
      "created_at": "2026-05-25T00:00:00Z",
      "vote_ends_at": "2026-05-26T00:00:00Z",
      "results_ends_at": "2026-05-27T00:00:00Z",
      "phase": "VOTING",
      "can_vote": true,
      "results_visible": true,
      "requires_pin": false
    }
  ]
}
```

Note: Feed items do not include options. Use GET /surveys/{id} to retrieve full survey details including options.

`description` is omitted when empty.

Errors:

- 400 BAD_REQUEST - unsupported sort or visibility query value
- 500 INTERNAL_SERVER_ERROR - failed to load feed

---

### GET /surveys/{id}/results

Request:

- No request body fields.

Response: 200 OK

```json
{
  "survey_id": "...",
  "total_votes": 0,
  "options": [
    {
      "id": "...",
      "text": "...",
      "vote_count": 0,
      "percentage": 0
    }
  ]
}
```

Errors:

- 400 BAD_REQUEST - invalid survey id
- 403 FORBIDDEN - results are not visible
- 404 NOT_FOUND - survey not found
- 500 INTERNAL_SERVER_ERROR - failed to load survey

---

## PIN (for PRIVATE_PIN)

### POST /surveys/{id}/pin/verify

Request:

- pin (string, required)

Response: 200 OK

```json
{ "ok": true }
```

Errors:

- 400 BAD_REQUEST - validation failure, invalid request body, or PIN not required for survey
- 403 FORBIDDEN (CONSENT_REQUIRED) - guest PIN verification requires service cookie
- 403 FORBIDDEN (PIN_REQUIRED) - invalid PIN
- 403 FORBIDDEN (PHASE_NOT_VOTING) - PIN verification cannot be stored because voting is closed
- 404 NOT_FOUND - survey not found
- 429 TOO_MANY_REQUESTS - too many failed guest PIN attempts
- 500 INTERNAL_SERVER_ERROR - failed to load survey or store PIN verification

---

## Voting

### POST /surveys/{id}/vote

Request:

- option_id (string, required)

Response: 200 OK

```json
{ "ok": true }
```

Rules:

- phase must be VOTING
- if PRIVATE_PIN, successful PIN verification must already exist through POST /surveys/{id}/pin/verify
- guests require consent/guest_id cookie
- enforce max_votes_per_user by user or guest
- increment Postgres survey_options.vote_count atomically

Errors:

- 400 BAD_REQUEST - invalid survey id, invalid request body, validation failure, or invalid option_id
- 403 CONSENT_REQUIRED - guest voting requires service cookie
- 403 PIN_REQUIRED - PIN verification required for PRIVATE_PIN survey
- 403 PHASE_NOT_VOTING - voting is not allowed in this phase
- 403 FORBIDDEN (MAX_VOTES_REACHED) - user or guest has reached the maximum vote limit for this survey
- 404 NOT_FOUND - survey not found
- 429 TOO_MANY_REQUESTS - vote rate limit exceeded
- 500 INTERNAL_SERVER_ERROR - failed to load survey, record vote, or store vote receipt

---

### PUT /surveys/{id}/vote

Request:

- new_option_id (string, required)

Response: 200 OK

```json
{ "ok": true }
```

Allowed only if:

- phase is VOTING
- max_votes_per_user is 1
- allow_vote_change_once is true
- change has not been used before by the user or guest
- if PRIVATE_PIN, successful PIN verification must already exist through POST /surveys/{id}/pin/verify
- guests require consent/guest_id cookie

Errors:

- 400 BAD_REQUEST - invalid survey id, invalid request body, validation failure, invalid new_option_id, or unchanged option
- 403 CONSENT_REQUIRED - guest vote change requires service cookie
- 403 PIN_REQUIRED - PIN verification required for PRIVATE_PIN survey
- 403 PHASE_NOT_VOTING - voting is not allowed in this phase
- 403 VOTE_CHANGE_NOT_ALLOWED - vote change is disabled, already used, or no previous vote exists
- 404 NOT_FOUND - survey not found
- 429 TOO_MANY_REQUESTS - vote rate limit exceeded
- 500 INTERNAL_SERVER_ERROR - failed to load survey, change vote, or store vote receipt

---

## Moderation / Safety

### POST /surveys/{id}/report

Request:

- reason (string, required)
- details (string, optional)

Response: 201 Created

```json
{ "ok": true }
```

Errors:

- 400 BAD_REQUEST - invalid survey id, invalid request body, or validation failure
- 401 UNAUTHORIZED - invalid access token when Authorization header is provided
- 404 NOT_FOUND - survey not found
- 500 INTERNAL_SERVER_ERROR - failed to load survey or create report

---

## Feedback (optional)

Not implemented in the current MVP. No backend route is registered.
This section is reserved for a future version.
