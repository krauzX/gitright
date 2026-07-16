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

export interface RepositoryAnalysis {
  repository: Repository;
  languages: Record<string, number>;
  files: string[];
  dependencies: Record<string, string[]>;
  key_files: Record<string, string>;
  commit_count: number;
  contributor_count: number;
}

export interface ContactPreferences {
  linkedin: string;
  personal_website: string;
  email: string;
  twitter: string;
  preferred_order: string[];
}

export interface Badge {
  name: string;
  url: string;
  color: string;
}

export interface ContentGenerationRequest {
  target_role: string;
  emphasized_skills: string[];
  tone_of_voice: "professional" | "friendly" | "technical";
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

export interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}

export interface TopRepository {
  name: string;
  stars: number;
  language: string;
  contributions: number;
}

export interface SocialLink {
  provider: string;
  url: string;
  username: string;
}

export interface AutoImportResponse {
  name: string;
  login: string;
  bio: string;
  avatar_url: string;
  company: string;
  location: string;
  email: string;
  website_url: string;
  twitter_username: string;
  followers: number;
  total_contributions: number;
  issue_contributions: number;
  review_contributions: number;
  skills: string[];
  language_breakdown: Record<string, number>;
  top_repositories: TopRepository[];
  social_links: SocialLink[];
}
