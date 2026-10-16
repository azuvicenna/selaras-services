package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/azuvicenna/selaras-services/b2b-gateway/internal/response"
)

type tokenBucket struct {
	tokens     float64
	maxTokens  float64
	refillRate float64 // tokens per second
	lastRefill time.Time
	mu         sync.Mutex
}

func (tb *tokenBucket) allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.lastRefill = now

	tb.tokens += elapsed * tb.refillRate
	if tb.tokens > tb.maxTokens {
		tb.tokens = tb.maxTokens
	}

	if tb.tokens >= 1.0 {
		tb.tokens -= 1.0
		return true
	}

	return false
}

type RateLimiter struct {
	rps     int
	burst   int
	buckets sync.Map
}

func NewRateLimiter(rps, burst int) *RateLimiter {
	return &RateLimiter{
		rps:   rps,
		burst: burst,
	}
}

func (rl *RateLimiter) getBucket(key string) *tokenBucket {
	val, ok := rl.buckets.Load(key)
	if ok {
		return val.(*tokenBucket)
	}

	bucket := &tokenBucket{
		tokens:     float64(rl.burst),
		maxTokens:  float64(rl.burst),
		refillRate: float64(rl.rps),
		lastRefill: time.Now(),
	}
	actual, _ := rl.buckets.LoadOrStore(key, bucket)
	return actual.(*tokenBucket)
}

func RateLimitMiddleware(rl *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("X-Partner-Id")
			if key == "" {
				key = r.RemoteAddr
			}

			bucket := rl.getBucket(key)
			if !bucket.allow() {
				w.Header().Set("Retry-After", "1")
				response.WriteError(w, http.StatusTooManyRequests, "Rate limit exceeded", "Too many requests. Please throttle your API calls.")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
