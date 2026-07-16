package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

type RateLimiterMemoryStoreConfig struct {
	Rate      rate.Limit
	Burst     int
	ExpiresIn time.Duration
}

type rateLimiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiterMemoryStore struct {
	limiters map[string]*rateLimiterEntry
	config   RateLimiterMemoryStoreConfig
	mu       sync.Mutex
	stopCh   chan struct{}
}

func NewRateLimiterMemoryStoreWithConfig(cfg RateLimiterMemoryStoreConfig) *RateLimiterMemoryStore {
	store := &RateLimiterMemoryStore{
		limiters: make(map[string]*rateLimiterEntry),
		config:   cfg,
		stopCh:   make(chan struct{}),
	}
	go store.cleanup()
	return store
}

func (s *RateLimiterMemoryStore) cleanup() {
	ticker := time.NewTicker(s.config.ExpiresIn / 2)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.mu.Lock()
			for key, entry := range s.limiters {
				if time.Since(entry.lastSeen) > s.config.ExpiresIn {
					delete(s.limiters, key)
				}
			}
			s.mu.Unlock()
		case <-s.stopCh:
			return
		}
	}
}

func (s *RateLimiterMemoryStore) getLimiter(key string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.limiters[key]
	if !exists {
		limiter := rate.NewLimiter(s.config.Rate, s.config.Burst)
		entry = &rateLimiterEntry{limiter: limiter, lastSeen: time.Now()}
		s.limiters[key] = entry
	} else {
		entry.lastSeen = time.Now()
	}
	return entry.limiter
}

func RateLimiter(store *RateLimiterMemoryStore) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := c.RealIP()
			limiter := store.getLimiter(key)
			if !limiter.Allow() {
				return echo.NewHTTPError(http.StatusTooManyRequests, "rate limit exceeded")
			}
			return next(c)
		}
	}
}
