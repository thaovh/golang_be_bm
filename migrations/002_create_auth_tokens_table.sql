-- Migration: Create auth_tokens table
-- Created: 2025-11-23

-- Create auth_tokens table
CREATE TABLE IF NOT EXISTS auth_tokens (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Audit fields
    created_by UUID NULL,
    updated_by UUID NULL,
    
    -- Optimistic locking
    version INTEGER NOT NULL DEFAULT 1,
    
    -- Status: active, inactive, archived, deleted
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    
    -- Token information
    user_id UUID NOT NULL,
    token TEXT NOT NULL,                    -- Access token (for reference)
    refresh_token TEXT NOT NULL,            -- Refresh token (UUID v7)
    expires_at TIMESTAMP NOT NULL,          -- Access token expiry
    refresh_expires_at TIMESTAMP NOT NULL,  -- Refresh token expiry
    ip_address VARCHAR(45) NULL,            -- IP address of login
    user_agent TEXT NULL,                   -- User agent of login
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    revoked_at TIMESTAMP NULL,
    
    -- Foreign key
    CONSTRAINT fk_auth_tokens_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_auth_tokens_deleted_at ON auth_tokens(deleted_at);
CREATE INDEX IF NOT EXISTS idx_auth_tokens_user_id ON auth_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_auth_tokens_refresh_token ON auth_tokens(refresh_token);
CREATE INDEX IF NOT EXISTS idx_auth_tokens_expires_at ON auth_tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_auth_tokens_refresh_expires_at ON auth_tokens(refresh_expires_at);
CREATE INDEX IF NOT EXISTS idx_auth_tokens_revoked ON auth_tokens(revoked);
CREATE INDEX IF NOT EXISTS idx_auth_tokens_status ON auth_tokens(status);

-- Create unique index on refresh_token (should be unique)
CREATE UNIQUE INDEX IF NOT EXISTS idx_auth_tokens_refresh_token_unique ON auth_tokens(refresh_token) WHERE deleted_at IS NULL AND revoked = FALSE;

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_auth_tokens_updated_at BEFORE UPDATE ON auth_tokens
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add comments
COMMENT ON TABLE auth_tokens IS 'Authentication tokens table for JWT refresh tokens';
COMMENT ON COLUMN auth_tokens.id IS 'UUID v7 primary key';
COMMENT ON COLUMN auth_tokens.user_id IS 'Reference to users table';
COMMENT ON COLUMN auth_tokens.token IS 'Access token (JWT) for reference';
COMMENT ON COLUMN auth_tokens.refresh_token IS 'Refresh token (UUID v7) stored in database';
COMMENT ON COLUMN auth_tokens.revoked IS 'Whether token has been revoked';

