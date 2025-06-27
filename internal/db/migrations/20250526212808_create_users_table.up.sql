CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    password TEXT, -- nullable for OAuth-only users
    full_name TEXT,
    avatar_url TEXT,

    -- Auth/OAuth
    provider TEXT DEFAULT 'local', -- 'google', 'facebook', etc.
    provider_id TEXT,             -- ID from OAuth provider
    email_verified BOOLEAN DEFAULT FALSE,
    verification_token TEXT,

    -- Roles and permissions
    role TEXT NOT NULL DEFAULT 'user', -- 'user', 'admin', etc.
    is_active BOOLEAN DEFAULT TRUE, -- for managers to suspend permissions

    -- Stripe integration
    stripe_customer_id TEXT,
    subscription_status TEXT,          -- 'active', 'trialing', 'past_due', etc.
    subscription_ends_at TIMESTAMPTZ,  -- for trial/expiry

    -- Auditing
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
