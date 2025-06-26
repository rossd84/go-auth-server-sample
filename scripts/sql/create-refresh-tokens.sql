CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL,
    issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    revoked_at TIMESTAMPTZ,

    -- Client metadata
    ip_address INET,
    user_agent TEXT,
    device_id TEXT,          -- Persistent identifier from the client
    location TEXT,           -- Derived from IP (e.g., "NY, USA")
    platform TEXT,           -- "iOS", "Android", "Web", etc.
    browser TEXT,            -- parsed from user-agent string
    session_id UUID,         -- Session management

    CONSTRAINT token_unique UNIQUE (user_id, token_hash)
)
