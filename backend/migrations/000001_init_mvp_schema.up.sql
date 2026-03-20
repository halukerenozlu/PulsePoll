CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE user_status AS ENUM ('active', 'suspended', 'deleted');
CREATE TYPE survey_visibility AS ENUM ('public', 'unlisted', 'private_pin');
CREATE TYPE survey_results_mode AS ENUM ('open_live', 'closed_hidden_until_end');
CREATE TYPE moderation_status AS ENUM ('approved', 'flagged', 'blocked', 'pending');
CREATE TYPE report_status AS ENUM ('open', 'reviewed', 'dismissed');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    display_name TEXT NOT NULL,
    status user_status NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE auth_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    refresh_token_hash TEXT NOT NULL UNIQUE,
    user_agent TEXT NOT NULL,
    ip TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE surveys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    creator_id UUID NOT NULL REFERENCES users(id),
    title TEXT NOT NULL,
    description TEXT,
    visibility survey_visibility NOT NULL,
    access_pin_hash TEXT,
    results_mode survey_results_mode NOT NULL,
    max_votes_per_user INT NOT NULL DEFAULT 1 CHECK (max_votes_per_user >= 1),
    allow_vote_change_once BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    vote_ends_at TIMESTAMPTZ NOT NULL,
    results_ends_at TIMESTAMPTZ NOT NULL,
    retention_ends_at TIMESTAMPTZ NOT NULL,
    moderation_status moderation_status NOT NULL DEFAULT 'approved',
    moderation_reason TEXT,
    CHECK (NOT allow_vote_change_once OR max_votes_per_user = 1)
);

CREATE TABLE survey_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    survey_id UUID NOT NULL REFERENCES surveys(id),
    text TEXT NOT NULL,
    position INT NOT NULL,
    vote_count BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (survey_id, position)
);

CREATE TABLE reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    survey_id UUID NOT NULL REFERENCES surveys(id),
    reporter_user_id UUID REFERENCES users(id),
    reporter_guest_id TEXT,
    reason TEXT NOT NULL,
    details TEXT,
    status report_status NOT NULL DEFAULT 'open',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    guest_id TEXT,
    message TEXT NOT NULL,
    page TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
