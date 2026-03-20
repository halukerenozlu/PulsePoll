# Project Quick Guide

This file is a short human-readable guide for the project owner.

For full project rules and AI tool behavior, use:

- `AGENTS.md`
- `docs/SPEC.md`
- `docs/API.md`
- `docs/DB.md`
- `docs/REDIS.md`
- `docs/TASKS.md`

## What this project is

PulsePoll is an MVP survey platform.

## Core product direction

- Registered users can create surveys.
- Guests cannot create surveys.
- Guests may vote only after accepting the required service cookie (`guest_id`).
- Guests may still browse surveys and view allowed results without accepting that cookie.

## Main stack

- Backend: Go + Fiber
- Database: PostgreSQL
- Ephemeral state / rate limits: Redis
- Frontend: Next.js

## Important rule

If a product or technical rule changes, update `docs/SPEC.md` first.

## Notes

- Keep the MVP small.
- Prefer simple and reliable solutions.
- Do not treat this file as the source of truth.
- The main working rules for AI tools are in `AGENTS.md`.
