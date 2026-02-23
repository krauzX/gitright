// User types
export interface User {
  id: number;
  github_id: number;
  username: string;
  email: string;
  avatar_url: string;
  bio: string;
  location: string;
  company: string;
  blog: string;
  created_at: string;
  updated_at: string;
  last_login_at: string;
}

// Repository types
export interface Repository {
  id: number;
  github_id: number;
  name: string;
  full_name: string;
  description: string;
  private: boolean;
  fork: boolean;
  language: string;
  stargazers_count: number;
  forks_count: number;
  open_issues_count: number;
  default_branch: string;
  topics: string[];
  html_url: string;
  clone_url: string;
  created_at: string;
  updated_at: string;
  pushed_at: string;
}

// Repository analysis types
export interface RepositoryAnalysis {
  repository: Repository;
  languages: Record<string, number>;
  files: string[];
  dependencies: Record<string, string[]>;
  key_files: Record<string, string>;
  commit_count: number;
  contributor_count: number;
}

// Project types
export type FocusTag =
  | "best_performance"
  | "team_project"
  | "personal_favorite"
  | null;

export interface Project {
  id: number;
  user_id: number;
  github_id: number; // Direct GitHub repo ID
  full_name: string; // owner/repo format
  priority: number;
  focus_tag: FocusTag;
  custom_summary: string;
  generated_summary: string;
  include_in_profile: boolean;
  created_at: string;
  updated_at: string;
}

export interface ProjectWithRepository extends Project {
  repository: Repository;
}

// Profile configuration types
export type ToneOfVoice = "professional" | "friendly" | "technical";
export type TemplateID =
  | "technical_deep_dive"
  | "hiring_manager_scan"
  | "community_contributor";

export interface ContactPreferences {
  linkedin: string;
  personal_website: string;
  email: string;
  twitter: string;
  preferred_order: string[];
}

export interface ProfileConfig {
  id: number;
  user_id: number;
  target_role: string;
  skills_emphasis: string[];
  tone_of_voice: ToneOfVoice;
  template_id: TemplateID;
  contact_prefs: ContactPreferences;
  show_private_repos: boolean;
  created_at: string;
  updated_at: string;
}

// Content generation types
export interface Badge {
  name: string;
  url: string;
  color: string;
}

export interface ContentGenerationRequest {
  target_role: string;
  emphasized_skills: string[];
  tone_of_voice: ToneOfVoice;
  contact_prefs: ContactPreferences;
  projects: RepositoryAnalysis[];
  user_api_key: string;
}

export interface ContentGenerationResponse {
  markdown: string;
  extracted_skills: string[];
  suggested_badges: Badge[];
  confidence: number;
}

// Generated profile types
export interface GeneratedProfile {
  id: number;
  user_id: number;
  config_id: number;
  content: string;
  markdown_preview: string;
  deployed: boolean;
  deployed_at: string | null;
  version: number;
  created_at: string;
}

// Template data types
export interface ProfilePitch {
  content: string;
  confidence: number;
}

export interface ProjectSummary {
  project: Project;
  repository: Repository;
  summary: string;
  tech_stack: string[];
  highlights: string[];
}

export interface TemplateData {
  user: User;
  config: ProfileConfig;
  projects: ProjectSummary[];
  skills: string[];
  badges: Badge[];
  profile_pitch: ProfilePitch;
}

// API response types
export interface APIError {
  error: string;
  message: string;
  code?: string;
}

export interface APIResponse<T> {
  data: T;
  message?: string;
}

// Auth types
export interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}
