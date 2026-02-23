-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    github_id BIGINT NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    avatar_url TEXT,
    bio TEXT,
    location VARCHAR(255),
    company VARCHAR(255),
    blog VARCHAR(500),
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    token_expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_login_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_github_id ON users(github_id);
CREATE INDEX idx_users_username ON users(username);

-- Create repositories table (cache)
CREATE TABLE IF NOT EXISTS repositories (
    id BIGSERIAL PRIMARY KEY,
    github_id BIGINT NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    full_name VARCHAR(500) NOT NULL,
    description TEXT,
    private BOOLEAN DEFAULT FALSE,
    fork BOOLEAN DEFAULT FALSE,
    language VARCHAR(100),
    stargazers_count INTEGER DEFAULT 0,
    forks_count INTEGER DEFAULT 0,
    open_issues_count INTEGER DEFAULT 0,
    default_branch VARCHAR(100),
    topics TEXT[], -- Array of topics
    html_url TEXT NOT NULL,
    clone_url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    pushed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    cached_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_repositories_github_id ON repositories(github_id);
CREATE INDEX idx_repositories_language ON repositories(language);
CREATE INDEX idx_repositories_cached_at ON repositories(cached_at);

-- Create projects table (selected repositories for profile)
CREATE TABLE IF NOT EXISTS projects (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    github_id BIGINT NOT NULL, -- Direct GitHub repo ID, not FK
    full_name VARCHAR(500) NOT NULL, -- owner/repo format
    priority INTEGER NOT NULL DEFAULT 0,
    focus_tag VARCHAR(50), -- "best_performance", "team_project", "personal_favorite"
    custom_summary TEXT,
    generated_summary TEXT,
    include_in_profile BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, github_id)
);

CREATE INDEX idx_projects_user_id ON projects(user_id);
CREATE INDEX idx_projects_github_id ON projects(github_id);
CREATE INDEX idx_projects_user_priority ON projects(user_id, priority);
CREATE INDEX idx_projects_user_include ON projects(user_id, include_in_profile);

-- Create profile_configs table
CREATE TABLE IF NOT EXISTS profile_configs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    target_role VARCHAR(255),
    skills_emphasis TEXT[], -- Array of emphasized skills
    tone_of_voice VARCHAR(50) DEFAULT 'professional', -- "professional", "friendly", "technical"
    template_id VARCHAR(100) DEFAULT 'hiring_manager_scan',
    contact_prefs JSONB, -- JSON for contact preferences
    show_private_repos BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_profile_configs_user_id ON profile_configs(user_id);

-- Create generated_profiles table
CREATE TABLE IF NOT EXISTS generated_profiles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    config_id BIGINT NOT NULL REFERENCES profile_configs(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    markdown_preview TEXT NOT NULL,
    deployed BOOLEAN DEFAULT FALSE,
    deployed_at TIMESTAMP WITH TIME ZONE,
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_generated_profiles_user_id ON generated_profiles(user_id);
CREATE INDEX idx_generated_profiles_deployed ON generated_profiles(deployed);

-- Create repository_analysis_cache table
CREATE TABLE IF NOT EXISTS repository_analysis_cache (
    id BIGSERIAL PRIMARY KEY,
    github_id BIGINT NOT NULL UNIQUE, -- Reference GitHub ID directly, not DB ID
    full_name VARCHAR(500) NOT NULL,
    languages JSONB,
    dependencies JSONB,
    key_files JSONB,
    commit_count INTEGER DEFAULT 0,
    contributor_count INTEGER DEFAULT 0,
    analyzed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() + INTERVAL '7 days'
);

CREATE INDEX idx_repository_analysis_cache_github_id ON repository_analysis_cache(github_id);
CREATE INDEX idx_repository_analysis_cache_full_name ON repository_analysis_cache(full_name);
CREATE INDEX idx_repository_analysis_cache_expires_at ON repository_analysis_cache(expires_at);

-- Create sessions table (for OAuth state management)
CREATE TABLE IF NOT EXISTS sessions (
    id VARCHAR(255) PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    data JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for updated_at columns
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_repositories_updated_at BEFORE UPDATE ON repositories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_projects_updated_at BEFORE UPDATE ON projects
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_profile_configs_updated_at BEFORE UPDATE ON profile_configs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert default template data (optional)
COMMENT ON TABLE users IS 'Stores authenticated GitHub users';
COMMENT ON TABLE repositories IS 'Caches GitHub repository information';
COMMENT ON TABLE projects IS 'Stores user-selected projects for profile generation';
COMMENT ON TABLE profile_configs IS 'Stores user profile customization settings';
COMMENT ON TABLE generated_profiles IS 'Stores generated README content versions';
COMMENT ON TABLE repository_analysis_cache IS 'Caches repository analysis to reduce API calls';
COMMENT ON TABLE sessions IS 'Stores session data for OAuth and user state management';
