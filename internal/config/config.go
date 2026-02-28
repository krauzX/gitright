package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Environment string
	Host        string
	Port        int
	LogLevel    string
	FrontendURL string
	HTTPTimeout time.Duration

	GitHub    GitHubConfig
	GoogleAI  GoogleAIConfig
	Database  DatabaseConfig
	Session   SessionConfig
	CORS      CORSConfig
	RateLimit RateLimitConfig
	Security  SecurityConfig
}

type GitHubConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scopes       []string
}

type GoogleAIConfig struct {
	APIKey       string
	Model        string
	UseGrounding bool
	Timeout      time.Duration
	MaxRetries   int
}

type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type SessionConfig struct {
	Secret        string
	MaxAge        int
	EncryptionKey string
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

type RateLimitConfig struct {
	RequestsPerMinute int
	Burst             int
}

type SecurityConfig struct {
	EnableHTTPS bool
	CertFile    string
	KeyFile     string
}

// Load reads all configuration from environment variables. Returns a joined
// error listing every missing required variable so operators see all problems
// at once rather than fixing them one restart at a time.
func Load() (*Config, error) {
	var missing []error
	requireEnv := func(key string) string {
		v := os.Getenv(key)
		if v == "" {
			missing = append(missing, fmt.Errorf("required environment variable %s is not set", key))
		}
		return v
	}

	githubClientID := requireEnv("GITHUB_CLIENT_ID")
	githubClientSecret := requireEnv("GITHUB_CLIENT_SECRET")
	databaseURL := requireEnv("DATABASE_URL")
	sessionSecret := requireEnv("SESSION_SECRET")
	tokenEncryptionKey := requireEnv("TOKEN_ENCRYPTION_KEY")

	if len(missing) > 0 {
		return nil, fmt.Errorf("configuration errors:\n%w", errors.Join(missing...))
	}

	cfg := &Config{
		Environment: getEnv("ENV", "development"),
		Host:        getEnv("HOST", "0.0.0.0"),
		Port:        getEnvAsInt("PORT", 8080),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		HTTPTimeout: getEnvAsDuration("HTTP_TIMEOUT", 5*time.Minute),

		GitHub: GitHubConfig{
			ClientID:     githubClientID,
			ClientSecret: githubClientSecret,
			RedirectURI:  getEnv("GITHUB_REDIRECT_URI", "http://localhost:3000/auth/callback"),
			Scopes:       strings.Split(getEnv("GITHUB_OAUTH_SCOPES", "repo,user:email"), ","),
		},

		GoogleAI: GoogleAIConfig{
			APIKey:       getEnv("GOOGLE_AI_API_KEY", ""),
			Model:        getEnv("GOOGLE_AI_MODEL", "gemini-2.5-flash-preview-0409-2025"),
			UseGrounding: getEnvAsBool("GOOGLE_AI_USE_GROUNDING", false),
			Timeout:      getEnvAsDuration("GOOGLE_AI_TIMEOUT", 5*time.Minute),
			MaxRetries:   getEnvAsInt("GOOGLE_AI_MAX_RETRIES", 2),
		},

		Database: DatabaseConfig{
			URL:             databaseURL,
			MaxOpenConns:    getEnvAsInt("DATABASE_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DATABASE_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsDuration("DATABASE_CONN_MAX_LIFETIME", 5*time.Minute),
		},

		Session: SessionConfig{
			Secret:        sessionSecret,
			MaxAge:        getEnvAsInt("SESSION_MAX_AGE", 86400),
			EncryptionKey: tokenEncryptionKey,
		},

		CORS: CORSConfig{
			AllowedOrigins: strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"), ","),
			AllowedMethods: strings.Split(getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"), ","),
			AllowedHeaders: strings.Split(getEnv("CORS_ALLOWED_HEADERS", "Content-Type,Authorization"), ","),
		},

		RateLimit: RateLimitConfig{
			RequestsPerMinute: getEnvAsInt("RATE_LIMIT_REQUESTS_PER_MINUTE", 60),
			Burst:             getEnvAsInt("RATE_LIMIT_BURST", 10),
		},

		Security: SecurityConfig{
			EnableHTTPS: getEnvAsBool("ENABLE_HTTPS", false),
			CertFile:    getEnv("CERT_FILE", ""),
			KeyFile:     getEnv("KEY_FILE", ""),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if len(c.Session.EncryptionKey) != 32 {
		return fmt.Errorf("TOKEN_ENCRYPTION_KEY must be exactly 32 bytes for AES-256")
	}

	if c.Security.EnableHTTPS {
		if c.Security.CertFile == "" || c.Security.KeyFile == "" {
			return fmt.Errorf("CERT_FILE and KEY_FILE must be set when HTTPS is enabled")
		}
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
