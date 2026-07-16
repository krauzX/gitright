-- Fix profile cache: drop FK constraint on config_id since profile_configs is no longer used
ALTER TABLE generated_profiles DROP CONSTRAINT IF EXISTS generated_profiles_config_id_fkey;
ALTER TABLE generated_profiles ALTER COLUMN config_id DROP NOT NULL;
ALTER TABLE generated_profiles ALTER COLUMN config_id SET DEFAULT 0;
