# SPEC — PulsePoll (MVP)

## Survey Phases

Let:

- vote_ends_at
- results_ends_at
- retention_ends_at

Phase rules:

- VOTING: now < vote_ends_at
- RESULTS: vote_ends_at <= now < results_ends_at
- EXPIRED: now >= results_ends_at

Defaults (MVP):

- vote_ends_at = created_at + 24h
- results_ends_at = created_at + 48h
- retention_ends_at = created_at + 48h

## Visibility / Access

- PUBLIC: appears in feed
- UNLISTED: accessible by link, not in feed
- PRIVATE_PIN: accessible by link; **voting requires PIN** (results visible to link-holders when allowed)

## Results Mode

- OPEN_LIVE: results visible during voting
- CLOSED_HIDDEN_UNTIL_END: results hidden until RESULTS phase

## Voting Rules

- Registered users can create surveys.
- Guests cannot create surveys.
- Guests can vote only if they accept the “service-required cookie” (guest_id).
- max_votes_per_user: integer >= 1
- allow_vote_change_once: only valid if max_votes_per_user == 1
- vote change allowed only during VOTING phase, at most once.

Additional MVP decisions:

- Guests may vote (with consent); guests can browse and view results without consent.
- Same option can be voted multiple times by the same user/guest until max_votes_per_user is reached.

## Consent Gating (Guest Voting)

Guests can always:

- browse public feed
- view survey details
- view results (when results are allowed by phase)

Guests can vote/change vote **only if** they accept the “service-required cookie”.
This cookie stores a random guest_id (no personal data) to:

- prevent spam/repeat voting abuse
- enforce vote limits per survey
- enforce “one-time vote change”
- remember short-lived PIN verification status

Retention:

- guest_id is short-lived (recommended: until retention_ends_at, default 48h).

Recommended banner copy (MVP):

- Title: “Oy kullanmak için gerekli çerez”
- Text: “Oyları adil tutmak, spam’i engellemek ve ‘1 kez oy değiştirme’ hakkını uygulamak için gerekli bir çerez kullanıyoruz.
  Kişisel bilgi içermez ve 48 saat içinde otomatik olarak geçersiz olur. Kabul etmeden de anketleri ve sonuçları görüntüleyebilirsin.”
- Buttons: “Kabul et & Oy ver” / “Şimdi değil (sadece görüntüle)”

## Data Storage (MVP)

- Postgres stores: users, surveys, options, aggregated counts.
- Redis stores: vote receipts for limits and one-time change tracking (TTL until retention_ends_at) + PIN state.
- No raw vote event log in MVP.

## Moderation (MVP)

- Basic keyword filter on survey title/description/options at creation.
- If blocked terms are detected: hard-block creation (reject request).
- Report endpoint exists for users/guests.
