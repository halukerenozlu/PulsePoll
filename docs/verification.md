# verification.md

## Purpose

This document defines how to verify the PulsePoll backend locally without depending on the frontend.

The goal is to make backend behavior testable and repeatable before `v0.4.0` frontend integration begins.

---

## Verification Principles

- Backend verification must not depend on UI work.
- Verification should be reproducible by following the same steps each time.
- Both success paths and failure paths should be checked.
- Documentation must reflect real behavior, not assumptions.
- If implementation and docs conflict, `docs/SPEC.md` is the product behavior source of truth and related docs must be updated.
- Backend work should map to an explicit Version Milestone, Work Item, and Implementation Slice from `docs/VERSION_PLAN.md`.

---

## Scope

This document covers:

- local backend startup verification
- dependency health verification
- direct HTTP verification
- basic error-path checks
- persistence verification guidance
- evidence to collect before approving backend work

This document does **not** define product behavior.
Product behavior belongs in:

- `docs/SPEC.md`
- `docs/API.md`
- `docs/DB.md`
- `docs/REDIS.md`

---

## Local Verification Prerequisites

Before running checks, confirm:

- Docker Desktop is running
- required services can start locally
- backend environment variables are present
- PostgreSQL is available
- Redis is available
- backend can boot without frontend

---

## Default Local Assumptions

Adjust these if local config changes.

- Backend base URL: `http://localhost:8080`
- Health endpoint: `GET /health`

Example:

```bash
BASE_URL="http://localhost:8080"
```

For PowerShell:

```powershell
$BASE_URL = "http://localhost:8080"
```

---

## Step 1 - Start Local Services

Start the project using the repo's local development flow.

Example:

```bash
docker compose -p pulsepoll up --build
```

Expected result:

- backend service starts successfully
- PostgreSQL is reachable
- Redis is reachable
- backend logs do not show immediate fatal startup errors

---

## Step 2 - Verify Service Health

### Request

```bash
curl -i "$BASE_URL/health"
```

PowerShell example:

```powershell
curl.exe -i "$BASE_URL/health"
```

### Expected result

- HTTP status is `200 OK`
- response is JSON
- response includes healthy dependency state

Expected shape:

```json
{
  "db": "up",
  "ok": true,
  "redis": "up"
}
```

If this fails, stop here and fix environment/startup issues before verifying milestone work.

---

## Step 3 - Baseline Backend Boot Check

Confirm the backend is stable after startup:

- no repeated crash/restart loop
- no immediate connection failure to DB
- no immediate connection failure to Redis
- health remains healthy across multiple checks

Optional repeat check:

```bash
curl -i "$BASE_URL/health"
curl -i "$BASE_URL/health"
curl -i "$BASE_URL/health"
```

---

## Step 4 - Milestone Endpoint Verification

Each approved backend Implementation Slice must be verified directly through HTTP requests when it affects endpoint behavior.

For each endpoint, verify:

- success path
- invalid input path
- not-found path where applicable
- conflict path where applicable
- persistence side effects where applicable
- response shape
- status code
- error shape

Use the checklist below for each endpoint added or changed by a Version Milestone.

---

## Endpoint Verification Template

Copy this block per endpoint as backend Work Items progress.

### Scope

- Version Milestone:
- Work Item:
- Implementation Slice:

### Endpoint

- Method:
- Path:
- Purpose:

### Success Case

- Request:
- Expected status:
- Expected response shape:
- Persistence effect:
- Notes:

### Failure Cases

Invalid input:

- Request:
- Expected status:
- Expected response shape:

Not found:

- Request:
- Expected status:
- Expected response shape:

Conflict / edge case:

- Request:
- Expected status:
- Expected response shape:

### Verification Result

- [ ] Success path verified
- [ ] Invalid input verified
- [ ] Not found verified if applicable
- [ ] Conflict/edge case verified if applicable
- [ ] Response shape matches docs
- [ ] Persistence effect verified if applicable

---

## Step 5 - Persistence Verification Guidance

For endpoints that create, update, or delete data, verification should include persistence checks.

Possible verification methods:

- read-back API call
- database inspection
- known side-effect confirmation
- follow-up request confirming state change

Examples of what to verify:

- record was created
- record was updated correctly
- record was not duplicated unexpectedly
- delete/close/expire logic changed state correctly
- related entities remain consistent

If a flow changes DB state, it is not enough to only check the first HTTP response.

---

## Step 6 - Error Handling Verification

For each relevant endpoint, verify backend error behavior is predictable.

Check:

- malformed JSON
- missing required fields
- invalid field values
- invalid route params
- unsupported state transitions
- resource not found
- duplicate/conflict scenarios where applicable

Minimum expectations:

- status code is correct
- error response is consistent
- response does not leak irrelevant internal details

---

## Verification Evidence

Before marking a backend Implementation Slice as complete, collect evidence such as:

- exact request used
- actual status code received
- actual response body
- follow-up verification result
- note on whether docs matched implementation

Acceptable evidence may include:

- terminal output
- saved curl commands
- Bruno requests
- integration test output
- short local verification notes

---

## Approval Checklist

A backend change should not be approved unless all relevant items below are true:

- [ ] Version Milestone, Work Item, and Implementation Slice are explicit
- [ ] local services boot successfully
- [ ] `/health` is healthy
- [ ] changed endpoint(s) are testable without frontend
- [ ] success path works
- [ ] failure path works
- [ ] response shape matches docs
- [ ] persistence effects were verified where applicable
- [ ] docs were updated if behavior changed

---

## Notes for Future Expansion

This document can later be extended with:

- Bruno collection references
- automated smoke test scripts
- integration test mapping
- endpoint-specific verification sections
- `v0.4.0` integration checks
