# PulsePoll Vision

PulsePoll is a lightweight, temporary survey MVP.

The product is intentionally small. It focuses on short-lived surveys, clear voting rules, limited results windows, and simple verification before broader product expansion.

## Product Direction

PulsePoll should make it easy to create and vote on temporary surveys while keeping the rules transparent:

- surveys have explicit voting and results windows
- guests can participate when the required consent flow allows it
- results visibility follows documented phase and mode rules
- private-link voting can be protected by PIN
- abusive repeat voting is limited with short-lived state

## Engineering Direction

The project is documentation-driven and verification-first.

The backend should be directly verifiable before the frontend becomes the main integration surface. Product rules, API behavior, database contracts, Redis contracts, and verification guidance live in `docs/` and should stay aligned as the MVP evolves.

## What PulsePoll Is Not Yet

PulsePoll is not trying to become a broad polling platform immediately.

The MVP should first prove that temporary surveys, consent-aware guest voting, results visibility, basic moderation, reporting, and rate limiting can work reliably with a small operational footprint.

Future expansion should come through documented milestones, not hidden scope creep.
