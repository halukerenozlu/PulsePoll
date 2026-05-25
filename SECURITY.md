# Security Policy

PulsePoll is an MVP project, but security-sensitive behavior should still be handled carefully and privately.

## Reporting a Vulnerability

Please do not report security vulnerabilities in public issues, pull requests, or discussions.

Report suspected vulnerabilities privately to the maintainer via GitHub: @halukerenozlu

When reporting, include as much detail as is safe to share:

- affected area or endpoint
- reproduction steps
- expected impact
- whether the issue requires authentication
- any relevant logs, requests, or responses

Do not include secrets, live tokens, private user data, or unrelated personal information in the report.

## Scope

Security-sensitive areas include:

- authentication and session behavior
- guest voting identity and consent handling
- PIN-protected survey access
- vote limits and vote-change enforcement
- rate limiting
- PostgreSQL persistence rules
- Redis key and TTL behavior
- API validation and error handling
- Docker/local configuration that could expose secrets

## Response Expectations

The maintainer will review reports as project availability allows.

Expected handling:

- acknowledge the report privately when possible
- reproduce and assess impact
- prepare a focused fix
- update docs or verification guidance if behavior changes
- publish a public summary only after the issue is understood and safe to disclose

## Supported Versions

PulsePoll is currently pre-release MVP software.

Security fixes target the active development line unless a release branch is explicitly documented later.

## Dependency Security

Contributors should avoid adding new dependencies unless they are needed for the approved implementation slice.

When adding or updating dependencies:

- prefer well-maintained packages
- keep lockfiles committed
- avoid broad dependency churn
- run relevant build and test checks
