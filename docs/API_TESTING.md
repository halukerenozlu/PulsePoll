# API_TESTING.md

Practical manual API testing guide for PulsePoll backend in local development.

## Local Assumptions

- Backend base URL: `http://localhost:8080`
- API base path: `/api/v1`
- Backend, PostgreSQL, and Redis are running locally

Example startup:

```powershell
docker compose -p pulsepoll up --build
```

## Testing Mindset (Short)

For each endpoint:

1. Prepare test data/prerequisites
2. Call the endpoint
3. Inspect status code + response body
4. Verify side effects when relevant (DB/Redis-backed behavior via follow-up API calls)

## Quick Variables (PowerShell)

```powershell
$BASE_URL = "http://localhost:8080"
$API = "$BASE_URL/api/v1"
```

---

## 1) Health Check

```powershell
curl.exe -i "$BASE_URL/health"
```

Expected:

- `200 OK`
- JSON like:

```json
{ "ok": true, "db": "up", "redis": "up" }
```

---

## 2) Register

```powershell
$registerBody = @{
  email = "tester1@example.com"
  password = "StrongPass123!"
  display_name = "tester1"
} | ConvertTo-Json

curl.exe -i -X POST "$API/auth/register" `
  -H "Content-Type: application/json" `
  -d $registerBody
```

Expected:

- `201 Created`
- `access_token` in response JSON

---

## 3) Login

```powershell
$loginBody = @{
  email = "tester1@example.com"
  password = "StrongPass123!"
} | ConvertTo-Json

$loginResponse = Invoke-RestMethod -Method Post -Uri "$API/auth/login" -ContentType "application/json" -Body $loginBody
$ACCESS_TOKEN = $loginResponse.access_token
$ACCESS_TOKEN
```

Expected:

- `200 OK`
- `access_token` returned

---

## 4) Create Survey

For easy `PUT /surveys/:id/vote` testing, create survey with vote-change enabled.

```powershell
$createSurveyBody = @{
  title = "Vote Change Test Survey"
  description = "Manual API testing"
  options = @("Option A", "Option B")
  visibility = "public"
  results_mode = "open_live"
  max_votes_per_user = 1
  allow_vote_change_once = $true
} | ConvertTo-Json

$createSurveyResponse = Invoke-RestMethod -Method Post -Uri "$API/surveys" `
  -Headers @{ Authorization = "Bearer $ACCESS_TOKEN" } `
  -ContentType "application/json" `
  -Body $createSurveyBody

$SURVEY_ID = $createSurveyResponse.id
$SURVEY_ID
```

Expected:

- `201 Created`
- Survey JSON with `id`, `options`, computed fields (`phase`, `can_vote`, `results_visible`, `requires_pin`)

---

## 5) Get Survey Details

```powershell
$surveyDetail = Invoke-RestMethod -Method Get -Uri "$API/surveys/$SURVEY_ID"
$surveyDetail | ConvertTo-Json -Depth 8
```

Expected:

- `200 OK`
- `phase` should be `VOTING` for a newly created default survey

---

## 6) Obtain Option IDs

```powershell
$OPTION_ID_1 = $surveyDetail.options[0].id
$OPTION_ID_2 = $surveyDetail.options[1].id

"Option1: $OPTION_ID_1"
"Option2: $OPTION_ID_2"
```

---

## 7) Vote (POST /surveys/:id/vote)

```powershell
$voteBody = @{ option_id = $OPTION_ID_1 } | ConvertTo-Json

curl.exe -i -X POST "$API/surveys/$SURVEY_ID/vote" `
  -H "Authorization: Bearer $ACCESS_TOKEN" `
  -H "Content-Type: application/json" `
  -d $voteBody
```

Expected:

- `200 OK`
- `{ "ok": true }`

---

## 8) Vote Change (PUT /surveys/:id/vote)

Prerequisites:

- Survey still in `VOTING`
- `max_votes_per_user = 1`
- `allow_vote_change_once = true`
- You already voted once

```powershell
$changeBody = @{ new_option_id = $OPTION_ID_2 } | ConvertTo-Json

curl.exe -i -X PUT "$API/surveys/$SURVEY_ID/vote" `
  -H "Authorization: Bearer $ACCESS_TOKEN" `
  -H "Content-Type: application/json" `
  -d $changeBody
```

Expected:

- `200 OK`
- `{ "ok": true }`

---

## 9) Results

```powershell
Invoke-RestMethod -Method Get -Uri "$API/surveys/$SURVEY_ID/results" | ConvertTo-Json -Depth 8
```

Expected:

- `200 OK` when results are visible (for `open_live`, visible during voting)
- `total_votes` and per-option percentages

---

## 10) Report

```powershell
$reportBody = @{
  reason = "manual_test"
  details = "report endpoint verification"
} | ConvertTo-Json

curl.exe -i -X POST "$API/surveys/$SURVEY_ID/report" `
  -H "Authorization: Bearer $ACCESS_TOKEN" `
  -H "Content-Type: application/json" `
  -d $reportBody
```

Expected:

- `201 Created`
- `{ "ok": true }`

---

## Additional Endpoint Examples

These examples cover endpoints that are useful during manual verification but are not part of the linear vote flow above.

### Consent Accept

```powershell
curl.exe -i -X POST "$API/consent/accept" `
  -c guest-cookies.txt
```

Expected:

- `200 OK`
- response body:

```json
{ "ok": true }
```

- response headers include a `Set-Cookie` for `guest_id`
- cookie is saved to `guest-cookies.txt`

You can inspect the saved cookie:

```powershell
Get-Content .\guest-cookies.txt
```

Expected:

- a `guest_id` cookie value is present
- the response `Set-Cookie` header includes `HttpOnly`

### Auth Refresh

Log in with a cookie jar so the refresh cookie is saved:

```powershell
$loginBody = @{
  email = "tester1@example.com"
  password = "StrongPass123!"
} | ConvertTo-Json

curl.exe -i -X POST "$API/auth/login" `
  -H "Content-Type: application/json" `
  -c auth-cookies.txt `
  -d $loginBody
```

Refresh using that cookie:

```powershell
curl.exe -i -X POST "$API/auth/refresh" `
  -b auth-cookies.txt `
  -c auth-cookies.txt
```

Expected:

- `200 OK`
- response JSON includes a fresh `access_token`
- refresh cookie remains available in `auth-cookies.txt`

### Auth Logout

```powershell
curl.exe -i -X POST "$API/auth/logout" `
  -b auth-cookies.txt
```

Expected:

- `200 OK`
- response body:

```json
{ "ok": true }
```

### Current User

```powershell
curl.exe -i -X GET "$API/me" `
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

Expected:

- `200 OK`
- response JSON contains the current user's `id`, `email`, and `display_name`

### Feed

Basic feed request:

```powershell
curl.exe -i "$API/feed"
```

Expected:

- `200 OK`
- response body has an `items` array
- only public surveys appear

Explicit newest-first sort:

```powershell
curl.exe -i "$API/feed?sort=new"
```

Expected:

- `200 OK`
- `items` are ordered by `created_at` descending

Search by keyword:

```powershell
$KEYWORD = "Vote"
curl.exe -i "$API/feed?sort=new&search=$KEYWORD"
```

Expected:

- `200 OK`
- `items` only include public surveys whose title or description matches the keyword

### PIN Verify

Create a private PIN survey:

```powershell
$pinSurveyBody = @{
  title = "Private PIN Manual Test"
  options = @("A", "B")
  visibility = "private_pin"
  access_pin = "1234"
  results_mode = "open_live"
  max_votes_per_user = 1
  allow_vote_change_once = $false
} | ConvertTo-Json

$pinSurvey = Invoke-RestMethod -Method Post -Uri "$API/surveys" `
  -Headers @{ Authorization = "Bearer $ACCESS_TOKEN" } `
  -ContentType "application/json" `
  -Body $pinSurveyBody

$PIN_SURVEY_ID = $pinSurvey.id
```

Accept guest consent first if testing as a guest:

```powershell
curl.exe -i -X POST "$API/consent/accept" `
  -c guest-cookies.txt
```

Successful PIN verification:

```powershell
$pinBody = @{ pin = "1234" } | ConvertTo-Json

curl.exe -i -X POST "$API/surveys/$PIN_SURVEY_ID/pin/verify" `
  -H "Content-Type: application/json" `
  -b guest-cookies.txt `
  -d $pinBody
```

Expected:

- `200 OK`
- response body:

```json
{ "ok": true }
```

- Redis receives a short-lived `pinok:survey:{surveyId}:guest:{guestId}` key for guest verification

Wrong PIN:

```powershell
$wrongPinBody = @{ pin = "9999" } | ConvertTo-Json

curl.exe -i -X POST "$API/surveys/$PIN_SURVEY_ID/pin/verify" `
  -H "Content-Type: application/json" `
  -b guest-cookies.txt `
  -d $wrongPinBody
```

Expected:

- `403 Forbidden`
- error code `PIN_REQUIRED`
- Redis increments `pinfail:survey:{surveyId}:guest:{guestId}` with a 15-minute TTL

---

## Vote Rate Limiting Verification

Scope: only vote endpoints

- `POST /surveys/:id/vote`
- `PUT /surveys/:id/vote`

Redis contract:

- key: `rl:ip:{ip}:vote`
- TTL: `60 seconds`

Important:

- Use the same client IP for all calls in this section so requests hit the same `rl:ip:{ip}:vote` key.
- Do not use business-rule-invalid requests to test rate limiting.

Note:

- PUT rate-limit verification should be isolated from preparatory POST requests.
- Prep votes also consume the same `rl:ip:{ip}:vote` bucket, so wait for the 60-second window to reset before starting the isolated PUT burst check.

### A) Verify 429 on repeated POST /vote (business-rule valid path)

Use a separate survey for POST rate-limit testing with a high vote limit so repeated POST requests stay valid.

```powershell
$postRateSurveyBody = @{
  title = "RateLimit POST Survey"
  options = @("A", "B")
  visibility = "public"
  results_mode = "open_live"
  max_votes_per_user = 200
  allow_vote_change_once = $false
} | ConvertTo-Json

$postRateSurvey = Invoke-RestMethod -Method Post -Uri "$API/surveys" `
  -Headers @{ Authorization = "Bearer $ACCESS_TOKEN" } `
  -ContentType "application/json" `
  -Body $postRateSurveyBody

$POST_RATE_SURVEY_ID = $postRateSurvey.id
$POST_RATE_OPTION_ID = $postRateSurvey.options[0].id
$postVoteBody = @{ option_id = $POST_RATE_OPTION_ID } | ConvertTo-Json
```

Now send repeated POST votes quickly:

```powershell
1..40 | ForEach-Object {
  curl.exe -s -o NUL -w "POST vote => HTTP %{http_code}`n" -X POST "$API/surveys/$POST_RATE_SURVEY_ID/vote" `
    -H "Authorization: Bearer $ACCESS_TOKEN" `
    -H "Content-Type: application/json" `
    -d $postVoteBody
}
```

Expected:

- early requests: normal non-429 responses (typically `200`)
- once limit is exceeded in the same 60-second window: deterministic `429 TOO_MANY_REQUESTS`

### B) Verify 429 on repeated PUT /vote (without consuming one-time change on same survey)

Do not spam PUT on one survey after change is already used.
Instead, prepare multiple surveys where each PUT request is still business-rule valid:

- `max_votes_per_user = 1`
- `allow_vote_change_once = true`
- one initial vote already placed
- change not used yet

```powershell
$prepared = @()

1..40 | ForEach-Object {
  $body = @{
    title = "RateLimit PUT Survey $_"
    options = @("A", "B")
    visibility = "public"
    results_mode = "open_live"
    max_votes_per_user = 1
    allow_vote_change_once = $true
  } | ConvertTo-Json

  $s = Invoke-RestMethod -Method Post -Uri "$API/surveys" `
    -Headers @{ Authorization = "Bearer $ACCESS_TOKEN" } `
    -ContentType "application/json" `
    -Body $body

  $optA = $s.options[0].id
  $optB = $s.options[1].id

  $firstVote = @{ option_id = $optA } | ConvertTo-Json
  Invoke-RestMethod -Method Post -Uri "$API/surveys/$($s.id)/vote" `
    -Headers @{ Authorization = "Bearer $ACCESS_TOKEN" } `
    -ContentType "application/json" `
    -Body $firstVote | Out-Null

  $prepared += [PSCustomObject]@{ SurveyId = $s.id; NewOptionId = $optB }
}
```

Then send one PUT change per prepared survey quickly:

```powershell
$prepared | ForEach-Object {
  $changeBody = @{ new_option_id = $_.NewOptionId } | ConvertTo-Json
  curl.exe -s -o NUL -w "PUT vote-change => HTTP %{http_code}`n" -X PUT "$API/surveys/$($_.SurveyId)/vote" `
    -H "Authorization: Bearer $ACCESS_TOKEN" `
    -H "Content-Type: application/json" `
    -d $changeBody
}
```

Expected:

- early PUT requests: normal non-429 responses (typically `200`)
- once limit is exceeded in the same 60-second window: deterministic `429 TOO_MANY_REQUESTS`

### C) Verify 60-second TTL reset

After you observe 429, wait 60+ seconds and retry a valid vote request from the same client IP.

```powershell
Start-Sleep -Seconds 61

curl.exe -i -X POST "$API/surveys/$POST_RATE_SURVEY_ID/vote" `
  -H "Authorization: Bearer $ACCESS_TOKEN" `
  -H "Content-Type: application/json" `
  -d $postVoteBody
```

Expected:

- request is no longer blocked by the previous 60-second rate-limit window

---

## Frontend Integration Reference

### POST /auth/register

```js
const res = await fetch('http://localhost:8080/api/v1/auth/register', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'user@example.com',
    password: 'StrongPass123!',
    display_name: 'user'
  }),
  credentials: 'include'
});
const data = await res.json();
// data.access_token, data.user.id
```

### POST /auth/login

```js
const res = await fetch('http://localhost:8080/api/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'user@example.com',
    password: 'StrongPass123!'
  }),
  credentials: 'include'
});
const data = await res.json();
// data.access_token, data.user.id
```

### GET /me

```js
const res = await fetch('http://localhost:8080/api/v1/me', {
  headers: { Authorization: `Bearer ${accessToken}` }
});
const data = await res.json();
// data.id, data.email, data.display_name
```

### POST /surveys

```js
const res = await fetch('http://localhost:8080/api/v1/surveys', {
  method: 'POST',
  headers: {
    Authorization: `Bearer ${accessToken}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    title: 'Frontend Test Survey',
    options: ['Option A', 'Option B'],
    visibility: 'public',
    results_mode: 'open_live'
  })
});
const data = await res.json();
// data.id, data.options
```

### GET /surveys/{id}

```js
const res = await fetch(`http://localhost:8080/api/v1/surveys/${surveyId}`);
const data = await res.json();
// data.id, data.options, data.phase
```

### GET /feed

```js
const res = await fetch('http://localhost:8080/api/v1/feed?sort=new&visibility=public');
const data = await res.json();
// data.items
```

### POST /surveys/{id}/vote

```js
const res = await fetch(`http://localhost:8080/api/v1/surveys/${surveyId}/vote`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    option_id: optionId
  }),
  credentials: 'include'
});
const data = await res.json();
// data.ok
```

### PUT /surveys/{id}/vote

```js
const res = await fetch(`http://localhost:8080/api/v1/surveys/${surveyId}/vote`, {
  method: 'PUT',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    new_option_id: newOptionId
  }),
  credentials: 'include'
});
const data = await res.json();
// data.ok
```

### GET /surveys/{id}/results

```js
const res = await fetch(`http://localhost:8080/api/v1/surveys/${surveyId}/results`);
const data = await res.json();
// data.total_votes, data.options
```

### POST /consent/accept

```js
const res = await fetch('http://localhost:8080/api/v1/consent/accept', {
  method: 'POST',
  credentials: 'include'
});
const data = await res.json();
// data.ok
```

### POST /surveys/{id}/pin/verify

```js
const res = await fetch(`http://localhost:8080/api/v1/surveys/${surveyId}/pin/verify`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    pin: '1234'
  }),
  credentials: 'include'
});
const data = await res.json();
// data.ok
```

## Known Limitations (v0.3.0)

1. POST /feedback is not implemented. No backend route is registered.
2. pinfail rate limiting blocks after 5 failed PIN attempts per guest. The counter is not reset on correct PIN. TTL is 15 minutes.
3. Feed endpoint does not return options. Use GET /surveys/{id} for full detail.
4. Timestamps follow RFC3339 format (e.g. 2026-05-25T00:00:00Z).
