-- Migration: Redis Removal and PostgreSQL-based Caching
-- Purpose: Replace Redis with PostgreSQL for all caching needs
-- Date: 2026-02-07

-- 1. OAuth state management (extend sessions table)
-- Add columns to store OAuth state in sessions table instead of Redis
ALTER TABLE sessions
  ADD COLUMN IF NOT EXISTS state_type VARCHAR(50) DEFAULT 'oauth_state',
  ADD COLUMN IF NOT EXISTS state_value TEXT;

-- Create index for fast OAuth state lookups
CREATE INDEX IF NOT EXISTS idx_sessions_state_type ON sessions(state_type, expires_at)
WHERE state_type = 'oauth_state';

-- 2. Profile generation cache (extend generated_profiles table)
-- Add caching metadata to support 24-hour profile caching
ALTER TABLE generated_profiles
  ADD COLUMN IF NOT EXISTS cache_key VARCHAR(500) UNIQUE,
  ADD COLUMN IF NOT EXISTS expires_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() + INTERVAL '24 hours',
  ADD COLUMN IF NOT EXISTS cache_hit_count INTEGER DEFAULT 0,
  ADD COLUMN IF NOT EXISTS last_accessed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();

-- Create indexes for cache lookups
CREATE INDEX IF NOT EXISTS idx_generated_profiles_cache_key ON generated_profiles(cache_key, expires_at);
CREATE INDEX IF NOT EXISTS idx_generated_profiles_expires_at ON generated_profiles(expires_at);

-- 3. Repository list caching (new table)
-- Cache GitHub repository lists to reduce API calls (5 minute TTL)
CREATE TABLE IF NOT EXISTS repository_list_cache (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    include_private BOOLEAN NOT NULL,
    repositories JSONB NOT NULL,
    cached_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() + INTERVAL '5 minutes',
    UNIQUE(user_id, include_private)
);

CREATE INDEX IF NOT EXISTS idx_repository_list_cache_user_id ON repository_list_cache(user_id, expires_at);
CREATE INDEX IF NOT EXISTS idx_repository_list_cache_expires_at ON repository_list_cache(expires_at);

-- 4. Cleanup function for expired data
-- Automatically clean up expired sessions, profiles, and repository cache
CREATE OR REPLACE FUNCTION cleanup_expired_data()
RETURNS void AS $$
BEGIN
  -- Clean expired sessions (OAuth states and regular sessions)
  DELETE FROM sessions WHERE expires_at < NOW();

  -- Clean expired profile cache (only non-deployed profiles)
  -- Keep deployed profiles forever for history
  DELETE FROM generated_profiles
  WHERE expires_at < NOW()
    AND cache_key IS NOT NULL
    AND NOT deployed;

  -- Clean expired repository list cache
  DELETE FROM repository_list_cache WHERE expires_at < NOW();

  -- Clean expired repository analysis cache (7 day TTL)
  DELETE FROM repository_analysis_cache WHERE expires_at < NOW();
END;
$$ LANGUAGE plpgsql;

-- 5. Analytics tables (for future analytics dashboard feature)
-- Track profile views, deployments, and user engagement
CREATE TABLE IF NOT EXISTS profile_analytics (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL, -- 'view', 'deploy', 'share', 'generate'
    metadata JSONB, -- Additional event data (e.g., IP country, referrer, etc.)
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_profile_analytics_user_id ON profile_analytics(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_profile_analytics_event_type ON profile_analytics(event_type, created_at DESC);

COMMENT ON TABLE profile_analytics IS 'Tracks user events for analytics dashboard';
COMMENT ON COLUMN profile_analytics.event_type IS 'Type of event: view, deploy, share, generate';
COMMENT ON COLUMN profile_analytics.metadata IS 'Additional event data stored as JSON';

-- 6. Collaborative sessions (for future collaborative editing feature)
-- Support multi-user collaboration on profiles
CREATE TABLE IF NOT EXISTS collaborative_sessions (
    id BIGSERIAL PRIMARY KEY,
    profile_id BIGINT NOT NULL REFERENCES generated_profiles(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    cursor_position INTEGER,
    active BOOLEAN DEFAULT TRUE,
    last_active_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_collaborative_sessions_profile ON collaborative_sessions(profile_id, active);
CREATE INDEX IF NOT EXISTS idx_collaborative_sessions_user ON collaborative_sessions(user_id, active);

COMMENT ON TABLE collaborative_sessions IS 'Tracks active collaborative editing sessions';
COMMENT ON COLUMN collaborative_sessions.cursor_position IS 'Cursor position in the markdown editor';
COMMENT ON COLUMN collaborative_sessions.active IS 'Whether the user is currently active in the session';

-- Add comments for new columns
COMMENT ON COLUMN sessions.state_type IS 'Type of session: oauth_state or regular session';
COMMENT ON COLUMN sessions.state_value IS 'OAuth state value for CSRF protection';
COMMENT ON COLUMN generated_profiles.cache_key IS 'Cache key for profile lookup: profile:v3:username:role:tone:projectcount';
COMMENT ON COLUMN generated_profiles.expires_at IS 'Cache expiration timestamp (24 hours for profiles)';
COMMENT ON COLUMN generated_profiles.cache_hit_count IS 'Number of times this cached profile was served';
COMMENT ON COLUMN generated_profiles.last_accessed_at IS 'Last time this cached profile was accessed';
COMMENT ON TABLE repository_list_cache IS 'Caches GitHub repository lists to reduce API calls';

-- Migration complete
-- Total tables created: 2 (repository_list_cache, profile_analytics, collaborative_sessions)
-- Total columns added: 6 (sessions: 2, generated_profiles: 4)
-- Total indexes created: 7
-- Functions created: 1 (cleanup_expired_data)
